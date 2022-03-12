package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	env "github.com/Netflix/go-env"

	"github.com/stebennett/squad-dashboard/cmd/jiracollector/models"
	"github.com/stebennett/squad-dashboard/cmd/jiracollector/repository"
	"github.com/stebennett/squad-dashboard/pkg/jiraservice"
	"github.com/stebennett/squad-dashboard/pkg/util"
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

	db := initDb()
	issueRepo := repository.NewPostgresIssueRepository(db)

	jiraParams := jiraservice.JiraParams{
		BaseUrl:   environment.JiraBaseUrl,
		User:      environment.JiraUser,
		AuthToken: environment.JiraAuthToken,
	}

	jiraClient := http.Client{
		Timeout: time.Second * 30,
	}

	execute(issueRepo, &jiraParams, environment.JiraQuery, environment.JiraEpicField, &jiraClient)
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

func execute(repo repository.IssueRepository, jiraParams *jiraservice.JiraParams, jql string, epicField string, jiraClient *http.Client) {

	query := jiraservice.JiraSearchQuery{
		Jql:        jql,
		Fields:     []string{"summary", "issuetype", epicField},
		Expand:     []string{"changelog"},
		StartAt:    0,
		MaxResults: 100,
	}

	log.Printf("Querying Jira for startAt: %d; maxResults: %d", query.StartAt, query.MaxResults)
	searchResult, err := jiraservice.MakeJiraSearchRequest(&query, jiraParams, jiraClient)
	if err != nil {
		log.Fatalf("Failed to make request %s", err)
	}

	for _, issue := range searchResult.Issues {
		saveableIssue := models.Create(issue)
		go repo.StoreIssue(context.Background(), saveableIssue)
	}

	var nextPageStartAt = util.NextPaginationArgs(0, 100, len(searchResult.Issues), searchResult.Total)
	for {
		if nextPageStartAt == -1 {
			log.Println("No new pages to fetch.")
			break
		}

		query.StartAt = nextPageStartAt

		log.Printf("Querying Jira for startAt: %d; maxResults: %d; total: %d", query.StartAt, query.MaxResults, searchResult.Total)
		searchResult, err := jiraservice.MakeJiraSearchRequest(&query, jiraParams, jiraClient)
		if err != nil {
			log.Fatalf("Failed to make request %s; startAt: %d", err, nextPageStartAt)
		}

		nextPageStartAt = util.NextPaginationArgs(nextPageStartAt, 100, len(searchResult.Issues), searchResult.Total)
	}
}
