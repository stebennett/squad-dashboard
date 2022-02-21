package jiraservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

type JiraParams struct {
	BaseUrl   string
	User      string
	AuthToken string
}

type JiraSearchQuery struct {
	Jql        string   `json:"jql"`
	Fields     []string `json:"fields"`
	Expand     []string `json:"expand"`
	StartAt    int      `json:"startAt"`
	MaxResults int      `json:"maxResults"`
}

type JiraIssue struct {
	Key    string `json:"key"`
	Fields struct {
		IssueType struct {
			Name string `json:"name"`
		} `json:"issuetype"`
	} `json:"fields"`
}

type JiraSearchResults struct {
	StartAt    int         `json:"startAt"`
	MaxResults int         `json:"maxResults"`
	Total      int         `json:"total"`
	Issues     []JiraIssue `json:"issues"`
}

func MakeJiraSearchRequest(jiraSearchQuery *JiraSearchQuery, jiraParams *JiraParams, httpClient *http.Client) (*JiraSearchResults, error) {

	url := fmt.Sprintf("https://%s/rest/api/2/search", jiraParams.BaseUrl)

	queryJSON, err := json.Marshal(jiraSearchQuery)
	if err != nil {
		return nil, err
	}
	log.Printf("Making query: %s", queryJSON)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(queryJSON))
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(jiraParams.User, jiraParams.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := httpClient.Do(req)
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

func MakeJiraGetHistoryRequest(issueKey string, jiraParams *JiraParams, httpClient *http.Client) {
	// TODO: Write code to return a page from history for a given issue - used for getting transitions
}
