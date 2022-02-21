package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	env "github.com/Netflix/go-env"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/stebennett/squad-dashboard/cmd/jiracollector/adapters"
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

	dbPool := initDb()
	defer dbPool.Close()
	repo := repository.NewDBIssueRepository(dbPool)

	jiraParams := jiraservice.JiraParams{
		BaseUrl:   environment.JiraBaseUrl,
		User:      environment.JiraUser,
		AuthToken: environment.JiraAuthToken,
	}

	query := jiraservice.JiraSearchQuery{
		Jql:        environment.JiraQuery,
		Fields:     []string{"summary", "issuetype", environment.JiraEpicField},
		Expand:     []string{"changelog"},
		StartAt:    0,
		MaxResults: 100,
	}

	jiraClient := http.Client{
		Timeout: time.Second * 30,
	}

	storeJiraIssues(repo, &query, &jiraParams, &jiraClient)
}

func storeJiraIssues(repo repository.IssueRepository, query *jiraservice.JiraSearchQuery, params *jiraservice.JiraParams, client *http.Client) {
	if query.StartAt == -1 {
		log.Println("No new pages to fetch.")
		return
	}

	log.Printf("Querying Jira for startAt: %d; maxResults: %d", query.StartAt, query.MaxResults)
	results, err := jiraservice.MakeJiraSearchRequest(query, params, client)
	if err != nil {
		log.Fatalf("Failed to make request %s", err)
	}

	for _, jiraIssue := range results.Issues {
		issue := adapters.AdaptIssue(jiraIssue)
		go repo.SaveIssue(context.Background(), issue)
	}

	query.StartAt = util.NextPaginationArgs(query.StartAt, query.MaxResults, len(results.Issues), results.Total)
	storeJiraIssues(repo, query, params, client)
}

func initDb() *pgxpool.Pool {
	connStr := os.ExpandEnv("postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$DB_SSLMODE") // load from env vars

	dbPool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database initialised")
	return dbPool
}
