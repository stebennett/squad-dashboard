package jiracollector

import (
	"context"
	"encoding/json"
	"log"
	"strings"

	"github.com/stebennett/squad-dashboard/cmd/jiracollector/models"
	"github.com/stebennett/squad-dashboard/cmd/jiracollector/repository"
	"github.com/stebennett/squad-dashboard/pkg/jiraservice"
	"github.com/stebennett/squad-dashboard/pkg/util"
)

type JiraCollector struct {
	repo repository.IssueRepository
	jira *jiraservice.JiraService
}

func NewJiraCollector(jira *jiraservice.JiraService, repo repository.IssueRepository) *JiraCollector {
	return &JiraCollector{
		repo: repo,
		jira: jira,
	}
}

func (jc *JiraCollector) Execute(jql string, epicField string) error {
	return jc.execute(0, 100, jql, epicField, jc.repo.SaveIssue)
}

func (jc *JiraCollector) execute(startAt int, maxResults int, jql string, epicField string, saveIssue func(ctx context.Context, jiraIssue models.JiraIssue) (int64, error)) error {
	query := jiraservice.JiraSearchQuery{
		Jql:        jql,
		Fields:     []string{"summary", "issuetype", epicField, "created", "updated"},
		Expand:     []string{"changelog"},
		StartAt:    startAt,
		MaxResults: maxResults,
	}

	log.Printf("Querying Jira for startAt: %d; maxResults: %d", query.StartAt, query.MaxResults)
	searchResult, err := jc.jira.MakeJiraSearchRequest(&query)
	if err != nil {
		log.Fatalf("Failed to make request %s", err)
	}

	var jiraResult models.JiraSearchResults
	jsonErr := json.Unmarshal([]byte(strings.ReplaceAll(searchResult, epicField, "epicKey")), &jiraResult)
	if jsonErr != nil {
		return jsonErr
	}

	for _, issue := range jiraResult.Issues {
		saveableIssue, err := models.Create(issue)
		if err != nil {
			log.Fatalf("Error: Failed to covert issue %s - %s", issue.Key, err)
			continue
		}

		_, err = saveIssue(context.Background(), *saveableIssue)
		if err != nil {
			log.Fatalf("Error: Failed to save issue %s - %s", saveableIssue.Key, err)
		}
	}

	var nextPageStartAt = util.NextPaginationArgs(startAt, maxResults, len(jiraResult.Issues), jiraResult.Total)
	if nextPageStartAt == -1 {
		log.Println("No new pages to fetch.")
		return nil
	}

	return jc.execute(nextPageStartAt, maxResults, jql, epicField, saveIssue)
}
