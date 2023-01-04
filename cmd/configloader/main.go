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
	NonWorkingDays     string `env:"NON_WORKING_DAYS,required=true"`
}

func main() {
	repo := createJiraRepository()

	var environment Environment
	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}

	_, err = repo.SaveJiraToDoStates(context.Background(), environment.JiraProject, strings.Split(environment.WorkToDoStates, ","))
	if err != nil {
		log.Fatalf("Failed to save Jira Work ToDo States. %s", err)
	}

	_, err = repo.SaveJiraInProgressStates(context.Background(), environment.JiraProject, strings.Split(environment.WorkStartStates, ","))
	if err != nil {
		log.Fatalf("Failed to save Jira Work In Progress States. %s", err)
	}

	_, err = repo.SaveJiraDoneStates(context.Background(), environment.JiraProject, strings.Split(environment.WorkCompleteStates, ","))
	if err != nil {
		log.Fatalf("Failed to save Jira Work Done States. %s", err)
	}

	_, err = repo.SaveNonWorkingDays(context.Background(), environment.JiraProject, strings.Split(environment.NonWorkingDays, ","))
	if err != nil {
		log.Fatalf("Failed to save non working days. %s", err)
	}

	log.Println("Config successfully loaded")
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
