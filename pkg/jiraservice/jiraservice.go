package jiraservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

type JiraService struct {
	jiraParams JiraParams
	httpClient *http.Client
}

func NewJiraService(httpClient *http.Client, jiraParams JiraParams) *JiraService {
	return &JiraService{
		jiraParams: jiraParams,
		httpClient: httpClient,
	}
}

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

func (js *JiraService) MakeJiraSearchRequest(jiraSearchQuery *JiraSearchQuery) (string, error) {

	url := fmt.Sprintf("https://%s/rest/api/2/search", js.jiraParams.BaseUrl)

	queryJSON, err := json.Marshal(jiraSearchQuery)
	if err != nil {
		return "", err
	}
	log.Printf("Making query: %s", queryJSON)

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(queryJSON))
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(js.jiraParams.User, js.jiraParams.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := js.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	log.Printf("response status: %s", resp.Status)
	log.Printf("response headers: %s", resp.Header)

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body), nil
}

func (js *JiraService) MakeJiraGetHistoryRequest(issueKey string, startAt int, maxResults int) (string, error) {
	url := fmt.Sprintf("https://%s/rest/api/2/issue/%s/changelog", js.jiraParams.BaseUrl, issueKey)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	req.SetBasicAuth(js.jiraParams.User, js.jiraParams.AuthToken)
	req.Header.Set("Content-Type", "application/json")

	query := req.URL.Query()
	query.Add("startAt", strconv.Itoa(startAt))
	query.Add("maxResults", strconv.Itoa(maxResults))

	req.URL.RawQuery = query.Encode()

	resp, err := js.httpClient.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	log.Printf("response status: %s", resp.Status)
	log.Printf("response headers: %s", resp.Header)

	body, _ := ioutil.ReadAll(resp.Body)

	return string(body), nil
}
