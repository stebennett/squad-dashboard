package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/stebennett/squad-dashboard/pkg/jirarepository"
)

func main() {
	// create a new database to store calculations
	issueRepo := createIssueRepository()
	project, set := os.LookupEnv("JIRA_PROJECT")
	if !set {
		log.Fatal("No JIRA_PROJECT set.")
	}

	// fetch all issues and set create year-week
	_, err := setCreateYearWeek(issueRepo, project)
	if err != nil {
		log.Fatalf("Failed to set created year-week for issues. %s", err)
	}

	// fetch all issues started and set started year-week

	// fetch all issues completed and set cycle time

	// fetch all issues completed and set lead time

	// fetch all issues completed and set complete year-week

	log.Fatal("Not yet implemented")
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

func setCreateYearWeek(repo jirarepository.JiraRepository, project string) (int64, error) {
	issues, err := repo.GetIssues(context.Background(), project)
	if err != nil {
		return -1, err
	}

	updatedCount := int64(0)

	for _, issue := range issues {
		year, week := issue.CreatedAt.UTC().ISOWeek()

		rowsChanged, err := repo.SaveCreateWeekDate(context.Background(), issue.Key, year, week)
		if err != nil {
			return updatedCount, err
		}

		updatedCount = updatedCount + rowsChanged
	}

	return updatedCount, nil
}
