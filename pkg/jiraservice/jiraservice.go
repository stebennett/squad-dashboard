package jiraservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

func MakeJiraSearchRequest(jiraSearchQuery *JiraSearchQuery, jiraBaseUrl string, jiraClient *http.Client, jiraUser string, jiraAuthToken string) (*JiraSearchResults, error) {

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
