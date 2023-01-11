package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Netflix/go-env"
	"github.com/stebennett/squad-dashboard/pkg/configloader"
	"github.com/stebennett/squad-dashboard/pkg/configrepository"
)

type Environment struct {
	JiraProject        string `env:"JIRA_PROJECT,required=true"`
	WorkToDoStates     string `env:"WORK_TODO_STATES,required=true"`
	WorkStartStates    string `env:"WORK_START_STATES,required=true"`
	WorkCompleteStates string `env:"WORK_COMPLETE_STATES,required=true"`
	NonWorkingDays     string `env:"NON_WORKING_DAYS,required=true"`
}

func main() {
	repo := createConfigRepository()

	var environment Environment
	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}

	config := configloader.Config{
		JiraProject:          environment.JiraProject,
		JiraToDoStates:       strings.Split(environment.WorkToDoStates, ","),
		JiraInProgressStates: strings.Split(environment.WorkStartStates, ","),
		JiraDoneStates:       strings.Split(environment.WorkCompleteStates, ","),
		NonWorkingDays:       strings.Split(environment.NonWorkingDays, ","),
	}

	err = configloader.Load(context.Background(), repo, config)
	if err != nil {
		log.Fatal(err)
	}

	log.Print("Config loaded successfully")
}

func createConfigRepository() configrepository.ConfigRepository {
	var err error
	var db *sql.DB
	connStr := os.ExpandEnv("postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable") // load from env vars

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database initialised")
	return configrepository.NewPostgresConfigRepository(db)
}
