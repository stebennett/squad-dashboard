package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/Netflix/go-env"
	"github.com/stebennett/squad-dashboard/cmd/jiraissuecalculator/calculator"
	"github.com/stebennett/squad-dashboard/pkg/jirarepository"
)

type Environment struct {
	JiraProject        string `env:"JIRA_PROJECT,required=true"`
	WorkToDoStates     string `env:"WORK_TODO_STATES,required=true"`
	WorkStartStates    string `env:"WORK_START_STATES,required=true"`
	WorkCompleteStates string `env:"WORK_COMPLETE_STATES,required=true"`
	ReportStartDate    string `env:"JIRA_REPORT_START_DATE,required=true"`
	ReportIssueTypes   string `env:"JIRA_REPORT_ISSUE_TYPES,required=true"`
	ReportEndStates    string `env:"JIRA_REPORT_END_STATES,required=true"`
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
	_, err = setCycleTimeForCompletedIssues(issueRepo, environment.JiraProject)
	if err != nil {
		log.Fatalf("Failed to set cycle time. %s", err)
	}
	log.Println("Completed updating cycle time for completed issues")

	// set number items completed for a given week
	_, err = setNumberOfItemsCompletedByWeek(issueRepo, environment.ReportStartDate, environment.JiraProject, strings.Split(environment.ReportIssueTypes, ","), strings.Split(environment.ReportEndStates, ","))
	if err != nil {
		log.Fatalf("Failed to set number of items completed by week. %s", err)
	}
	log.Println("Completed number of items completed reports")

	// set number items started for a given week
	_, err = setNumberOfItemsStartedByWeek(issueRepo, environment.ReportStartDate, environment.JiraProject, strings.Split(environment.ReportIssueTypes, ","))
	if err != nil {
		log.Fatalf("Failed to set number of items completed by week. %s", err)
	}
	log.Println("Completed number of items started reports")

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

		endState, err := repo.GetEndStateForIssue(context.Background(), issueKey, transitionTime)
		if err != nil {
			return updatedCount, err
		}

		rowsChanged, err := repo.SaveCompleteDates(context.Background(), issueKey, year, week, transitionTime, endState)
		if err != nil {
			return updatedCount, err
		}

		updatedCount = updatedCount + rowsChanged
	}

	return updatedCount, nil
}

func setCycleTimeForCompletedIssues(repo jirarepository.JiraRepository, project string) (int64, error) {
	calculations, err := repo.GetCompletedIssues(context.Background(), project)
	if err != nil {
		return -1, err
	}

	updatedCount := int64(0)

	for issueKey, calculations := range calculations {
		if !calculations.IssueStartedAt.Valid {
			return updatedCount, fmt.Errorf("null start time for issue %s", issueKey)
		}

		if !calculations.IssueCompletedAt.Valid {
			return updatedCount, fmt.Errorf("null completed time for issue %s", issueKey)
		}

		cycleTime, err := calculator.CalculateCycleTime(calculations.IssueStartedAt.Time, calculations.IssueCompletedAt.Time)

		if err != nil {
			return updatedCount, err
		}

		rowsChanged, err := repo.SaveCycleTime(context.Background(), issueKey, cycleTime, -1)
		if err != nil {
			return updatedCount, err
		}

		updatedCount = updatedCount + rowsChanged
	}

	return updatedCount, nil
}

func setNumberOfItemsCompletedByWeek(repo jirarepository.JiraRepository, startDateStr string, project string, issueTypes []string, endStates []string) (int64, error) {
	startDate, err := time.Parse("2006-01-02T15:04:05Z", startDateStr)
	if err != nil {
		return -1, err
	}

	totalUpdates := int64(0)

	for {
		if startDate.After(time.Now()) {
			break
		}

		endDate := startDate.AddDate(0, 0, 7)
		issues, err := repo.GetIssuesCompletedBetweenDates(context.Background(), project, startDate, endDate, issueTypes, endStates)
		if err != nil {
			return totalUpdates, err
		}

		rowsUpdated, err := repo.SetIssuesCompletedInWeekStarting(context.Background(), project, startDate, len(issues))
		if err != nil {
			return totalUpdates, err
		}

		totalUpdates = totalUpdates + rowsUpdated
		startDate = endDate
	}

	return totalUpdates, nil
}

func setNumberOfItemsStartedByWeek(repo jirarepository.JiraRepository, startDateStr string, project string, issueTypes []string) (int64, error) {
	startDate, err := time.Parse("2006-01-02T15:04:05Z", startDateStr)
	if err != nil {
		return -1, err
	}

	totalUpdates := int64(0)

	for {
		if startDate.After(time.Now()) {
			break
		}

		endDate := startDate.AddDate(0, 0, 7)

		issues, err := repo.GetIssuesStartedBetweenDates(context.Background(), project, startDate, endDate, issueTypes)
		if err != nil {
			return totalUpdates, err
		}

		rowsUpdated, err := repo.SetIssuesStartedInWeekStarting(context.Background(), project, startDate, len(issues))
		if err != nil {
			return totalUpdates, err
		}

		totalUpdates = totalUpdates + rowsUpdated
		startDate = endDate
	}

	return totalUpdates, nil
}
