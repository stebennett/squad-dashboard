package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type JiraSearchQuery struct {
	jql        string
	fields     []string
	expand     []string
	startAt    int
	maxResults int
}

type JiraSearchResults struct {
	StartAt    int `json:"startAt"`
	MaxResults int `json:"maxResults"`
	Total      int `json:"total"`
	Issues     []struct {
		Key    string `json:"key"`
		Fields struct {
			IssueType struct {
				Name string `json:"name"`
			} `json:"issuetype"`
		} `json:"fields"`
	} `json:"issues"`
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

	query := JiraSearchQuery{
		jql:        jiraQuery,
		fields:     []string{"summary", "issuetype", jiraEpicField},
		expand:     []string{"changelog"},
		startAt:    0,
		maxResults: 100,
	}

	jiraClient := http.Client{
		Timeout: time.Second * 30,
	}

	searchResult, err := makeJiraSearchRequest(&query, jiraBaseUrl, &jiraClient, jiraUser, jiraAuthToken)
	if err != nil {
		log.Fatalf("Failed to make request %s", err)
	}

	log.Printf("Query exectuted. Returned %d results", searchResult.MaxResults)
}

func makeJiraSearchRequest(jiraSearchQuery *JiraSearchQuery, jiraBaseUrl string, jiraClient *http.Client, jiraUser string, jiraAuthToken string) (*JiraSearchResults, error) {

	url := fmt.Sprintf("https://%s/rest/api/2/search", jiraBaseUrl)

	queryJSON, err := json.Marshal(jiraSearchQuery)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(queryJSON))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(jiraUser, jiraAuthToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := jiraClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	log.Printf("response status: %s", resp.Status)
	log.Printf("response headers: %s", resp.Header)

	body, _ := ioutil.ReadAll(resp.Body)

	var jiraResult JiraSearchResults
	jsonErr := json.Unmarshal(body, &jiraResult)
	if jsonErr != nil {
		return nil, jsonErr
	}

	return &jiraResult, nil
}
