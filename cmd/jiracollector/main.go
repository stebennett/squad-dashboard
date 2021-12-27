package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type JiraSearchQuery struct {
	jql    string
	fields []string
	expand []string
}

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

	url := fmt.Sprintf("https://%s/rest/api/2/search", jiraBaseUrl)

	query := &JiraSearchQuery{
		jql:    jiraQuery,
		fields: []string{"summary", "issuetype", jiraEpicField},
		expand: []string{"changelog"},
	}

	queryJSON, err := json.Marshal(query)
	if err != nil {
		log.Fatal("Failed to create Jira Query in JSON")
	}

	jiraClient := http.Client{
		Timeout: time.Second * 30,
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(queryJSON))
	if err != nil {
		log.Fatal(err)
	}

	req.SetBasicAuth(jiraUser, jiraAuthToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := jiraClient.Do(req)
	if err != nil {
		log.Fatalf("Failed to make request. Got error %s", err.Error())
	}

	defer resp.Body.Close()
}
