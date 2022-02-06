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

	"github.com/stebennett/squad-dashboard/pkg/util"
)

type JiraSearchQuery struct {
	Jql        string   `json:"jql"`
	Fields     []string `json:"fields"`
	Expand     []string `json:"expand"`
	StartAt    int      `json:"startAt"`
	MaxResults int      `json:"maxResults"`
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
	searchResult, err := makeJiraSearchRequest(&query, jiraBaseUrl, &jiraClient, jiraUser, jiraAuthToken)
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
		searchResult, err := makeJiraSearchRequest(&query, jiraBaseUrl, &jiraClient, jiraUser, jiraAuthToken)
		if err != nil {
			log.Fatalf("Failed to make request %s; startAt: %d", err, nextPageStartAt)
		}

		nextPageStartAt = util.NextPaginationArgs(nextPageStartAt, 100, len(searchResult.Issues), searchResult.Total)
	}
}

func makeJiraSearchRequest(jiraSearchQuery *JiraSearchQuery, jiraBaseUrl string, jiraClient *http.Client, jiraUser string, jiraAuthToken string) (*JiraSearchResults, error) {

	url := fmt.Sprintf("https://%s/rest/api/2/search", jiraBaseUrl)

	queryJSON, err := json.Marshal(jiraSearchQuery)
	if err != nil {
		return nil, err
	}
	log.Printf("Making query: %s", queryJSON)

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
