package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Netflix/go-env"
	"github.com/stebennett/squad-dashboard/pkg/jirarepository"
)

type Environment struct {
	JiraProject        string `env:"JIRA_PROJECT,required=true"`
	WorkToDoStates     string `env:"WORK_TODO_STATES,required=true"`
	WorkStartStates    string `env:"WORK_START_STATES,required=true"`
	WorkCompleteStates string `env:"WORK_COMPLETE_STATES,required=true"`
}

func main() {
	// create a new database to store calculations
	issueRepo := createIssueRepository()

	// load environment vars
	var environment Environment
	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("running updates for project %s", environment.JiraProject)

	// fetch all issues and set create year-week
	_, err = setCreateDates(issueRepo, environment.JiraProject)
	if err != nil {
		log.Fatalf("Failed to set created year-week for issues. %s", err)
	}
	log.Println("Completed update of created year-week for issues")

	// fetch all issues started and set started year-week
	_, err = setStartDates(issueRepo, environment.JiraProject, strings.Split(environment.WorkStartStates, ","), strings.Split(environment.WorkToDoStates, ","))
	if err != nil {
		log.Fatalf("Failed to set started year-week for issues. %s", err)
	}
	log.Println("Completed update of started year-week for issues")

	// fetch all issues completed and set complete year-week
	_, err = setCompleteDates(issueRepo, environment.JiraProject, strings.Split(environment.WorkCompleteStates, ","))
	if err != nil {
		log.Fatalf("Failed to set completed year-week for issues. %s", err)
	}
	log.Println("Completed update of completed year-week for issues")

	// fetch all issues completed and set cycle time
	// _, err = setCycleTimeForCompletedIssues(issueRepo, environment.JiraProject)
	// if err != nil {
	// 	log.Fatalf("Failed to set cycle time")
	// }
	// log.Println("Completed updating cycle time for completed issues")

	// fetch all issues completed and set lead time

	log.Fatal("All calculations complete")
}

func createIssueRepository() jirarepository.JiraRepository {
	var err error
	var db *sql.DB
	connStr := os.ExpandEnv("postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable") // load from env vars

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database initialised")
	return jirarepository.NewPostgresJiraRepository(db)
}

func setCreateDates(repo jirarepository.JiraRepository, project string) (int64, error) {
	issues, err := repo.GetIssues(context.Background(), project)
	if err != nil {
		return -1, err
	}

	updatedCount := int64(0)

	for _, issue := range issues {
		year, week := issue.CreatedAt.UTC().ISOWeek()

		rowsChanged, err := repo.SaveCreateDates(context.Background(), issue.Key, year, week, issue.CreatedAt)
		if err != nil {
			return updatedCount, err
		}

		updatedCount = updatedCount + rowsChanged
	}

	return updatedCount, nil
}

func setStartDates(repo jirarepository.JiraRepository, project string, workStartStates []string, workToDoStates []string) (int64, error) {
	transitions, err := repo.GetTransitionTimeByStateChanges(context.Background(), project, workToDoStates, workStartStates)
	if err != nil {
		return -1, err
	}

	updatedCount := int64(0)

	for issueKey, transitionTime := range transitions {
		year, week := transitionTime.UTC().ISOWeek()

		rowsChanged, err := repo.SaveStartDates(context.Background(), issueKey, year, week, transitionTime)
		if err != nil {
			return updatedCount, err
		}

		updatedCount = updatedCount + rowsChanged
	}

	return updatedCount, nil
}

func setCompleteDates(repo jirarepository.JiraRepository, project string, workCompleteStates []string) (int64, error) {
	transitions, err := repo.GetTransitionTimeByToState(context.Background(), project, workCompleteStates)
	if err != nil {
		return -1, err
	}

	updatedCount := int64(0)

	for issueKey, transitionTime := range transitions {
		year, week := transitionTime.UTC().ISOWeek()

		rowsChanged, err := repo.SaveCompleteDates(context.Background(), issueKey, year, week, transitionTime)
		if err != nil {
			return updatedCount, err
		}

		updatedCount = updatedCount + rowsChanged
	}

	return updatedCount, nil
}

// func setCycleTimeForCompletedIssues(repo jirarepository.JiraRepository, project string) (int64, error) {
// 	repo.GetCompletedIssues(context.Background(), project)
// }
