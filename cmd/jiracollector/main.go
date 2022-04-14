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
	"github.com/stebennett/squad-dashboard/pkg/jirarepository"
	"github.com/stebennett/squad-dashboard/pkg/jiraservice"
)

type Environment struct {
	JiraProject   string `env:"JIRA_PROJECT,required=true"`
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
	issueRepo := createIssueRepository()

	// create a new connection to jira
	jira := createJiraService(environment)

	// create a new collector job
	jiracollector := jiracollector.NewJiraCollector(jira, issueRepo)

	// execute the job
	err = jiracollector.Execute(environment.JiraProject, environment.JiraQuery, environment.JiraEpicField)
	log.Fatal(err)
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
