package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	env "github.com/Netflix/go-env"

	"github.com/stebennett/squad-dashboard/cmd/jiracollector/jiracollector"
	"github.com/stebennett/squad-dashboard/cmd/jiracollector/repository"
	"github.com/stebennett/squad-dashboard/pkg/jiraservice"
)

type Environment struct {
	JiraBaseUrl   string `env:"JIRA_HOST,required=true"`
	JiraUser      string `env:"JIRA_USER,required=true"`
	JiraAuthToken string `env:"JIRA_AUTH_TOKEN,required=true"`
	JiraQuery     string `env:"JIRA_QUERY,required=true"`
	JiraEpicField string `env:"JIRA_EPIC_FIELD,required=true"`
}

func main() {

	var environment Environment
	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}

	// create a new database to store jira issues
	db := initDb()
	issueRepo := repository.NewPostgresIssueRepository(db)

	// create a new connection to jira
	jira := createJiraService(environment)

	// create a new collector job
	jiracollector := jiracollector.NewJiraCollector(issueRepo, jira)

	// execute the job
	jiracollector.Execute(environment.JiraQuery, environment.JiraEpicField)
}

func initDb() *sql.DB {
	var err error
	var db *sql.DB
	connStr := os.ExpandEnv("postgres://$DB_USERNAME:$DB_PASSWORD@DB_HOST:$DB_PORT/$DB_NAME") // load from env vars

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("Database initialised")
	return db
}

func createJiraService(environment Environment) *jiraservice.JiraService {
	jiraParams := jiraservice.JiraParams{
		BaseUrl:   environment.JiraBaseUrl,
		User:      environment.JiraUser,
		AuthToken: environment.JiraAuthToken,
	}

	jiraClient := http.Client{
		Timeout: time.Second * 30,
	}

	return jiraservice.NewJiraService(&jiraClient, jiraParams)
}
