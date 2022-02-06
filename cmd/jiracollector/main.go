package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/stebennett/squad-dashboard/pkg/jiraservice"
	"github.com/stebennett/squad-dashboard/pkg/util"
)

func main() {
	jiraBaseUrl := os.Getenv("JIRA_HOST")
	if jiraBaseUrl == "" {
		log.Fatal("JIRA_HOST env var is not set")
	}

	jiraUser := os.Getenv("JIRA_USER")
	if jiraUser == "" {
		log.Fatal("JIRA_USER env var is not set")
	}

	jiraAuthToken := os.Getenv("JIRA_AUTH_TOKEN")
	if jiraAuthToken == "" {
		log.Fatal("JIRA_AUTH_TOKEN env var is not set")
	}

	jiraQuery := os.Getenv("JIRA_QUERY")
	if jiraQuery == "" {
		log.Fatal("JIRA_QUERY env var is not set")
	}

	jiraEpicField := os.Getenv("JIRA_EPIC_FIELD")
	if jiraEpicField == "" {
		log.Fatal("JIRA_EPIC_FIELD env var is not set")
	}

	query := jiraservice.JiraSearchQuery{
		Jql:        jiraQuery,
		Fields:     []string{"summary", "issuetype", jiraEpicField},
		Expand:     []string{"changelog"},
		StartAt:    0,
		MaxResults: 100,
	}

	jiraClient := http.Client{
		Timeout: time.Second * 30,
	}

	log.Printf("Querying Jira for startAt: %d; maxResults: %d", query.StartAt, query.MaxResults)
	searchResult, err := jiraservice.MakeJiraSearchRequest(&query, jiraBaseUrl, &jiraClient, jiraUser, jiraAuthToken)
	if err != nil {
		log.Fatalf("Failed to make request %s", err)
	}

	var nextPageStartAt = util.NextPaginationArgs(0, 100, len(searchResult.Issues), searchResult.Total)
	for {
		if nextPageStartAt == -1 {
			log.Println("No new pages to fetch.")
			break
		}

		query.StartAt = nextPageStartAt

		log.Printf("Querying Jira for startAt: %d; maxResults: %d; total: %d", query.StartAt, query.MaxResults, searchResult.Total)
		searchResult, err := jiraservice.MakeJiraSearchRequest(&query, jiraBaseUrl, &jiraClient, jiraUser, jiraAuthToken)
		if err != nil {
			log.Fatalf("Failed to make request %s; startAt: %d", err, nextPageStartAt)
		}

		nextPageStartAt = util.NextPaginationArgs(nextPageStartAt, 100, len(searchResult.Issues), searchResult.Total)
	}
}
