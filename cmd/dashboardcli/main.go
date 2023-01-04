package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/Netflix/go-env"
	"github.com/stebennett/squad-dashboard/pkg/dashboard"
	"github.com/stebennett/squad-dashboard/pkg/jirarepository"
	"github.com/stebennett/squad-dashboard/pkg/printer"
)

type Environment struct {
	JiraProject         string `env:"JIRA_PROJECT,required=true"`
	JiraDefectIssueType string `env:"JIRA_DEFECT_ISSUE_TYPE,required=true"`
}

func main() {
	var environment Environment
	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}

	repo := createJiraRepository()
	printer := createPrinter()

	escapedDefects, err := dashboard.GenerateEscapedDefects(12, environment.JiraProject, environment.JiraDefectIssueType, repo)
	if err != nil {
		log.Fatal(err)
	}

	printer.Print(escapedDefects)
}

func createJiraRepository() jirarepository.JiraRepository {
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

func createPrinter() printer.Printer {
	return printer.NewCommandLinePrinter()
}
