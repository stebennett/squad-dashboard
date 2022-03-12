package jiracollector

import (
	"context"
	"log"

	"github.com/stebennett/squad-dashboard/cmd/jiracollector/models"
	"github.com/stebennett/squad-dashboard/cmd/jiracollector/repository"
	"github.com/stebennett/squad-dashboard/pkg/jiraservice"
	"github.com/stebennett/squad-dashboard/pkg/util"
)

type JiraCollector struct {
	repo repository.IssueRepository
	jira *jiraservice.JiraService
}

func NewJiraCollector(repo repository.IssueRepository, jira *jiraservice.JiraService) *JiraCollector {
	return &JiraCollector{
		repo: repo,
		jira: jira,
	}
}

func (jc *JiraCollector) Execute(jql string, epicField string) {

	query := jiraservice.JiraSearchQuery{
		Jql:        jql,
		Fields:     []string{"summary", "issuetype", epicField},
		Expand:     []string{"changelog"},
		StartAt:    0,
		MaxResults: 100,
	}

	log.Printf("Querying Jira for startAt: %d; maxResults: %d", query.StartAt, query.MaxResults)
	searchResult, err := jc.jira.MakeJiraSearchRequest(&query)
	if err != nil {
		log.Fatalf("Failed to make request %s", err)
	}

	for _, issue := range searchResult.Issues {
		saveableIssue := models.Create(issue)
		go jc.repo.StoreIssue(context.Background(), saveableIssue)
	}

	var nextPageStartAt = util.NextPaginationArgs(0, 100, len(searchResult.Issues), searchResult.Total)
	for {
		if nextPageStartAt == -1 {
			log.Println("No new pages to fetch.")
			break
		}

		query.StartAt = nextPageStartAt

		log.Printf("Querying Jira for startAt: %d; maxResults: %d; total: %d", query.StartAt, query.MaxResults, searchResult.Total)
		searchResult, err := jc.jira.MakeJiraSearchRequest(&query)
		if err != nil {
			log.Fatalf("Failed to make request %s; startAt: %d", err, nextPageStartAt)
		}

		nextPageStartAt = util.NextPaginationArgs(nextPageStartAt, 100, len(searchResult.Issues), searchResult.Total)
	}
}

func (jc *JiraCollector) execute(startAt int, maxResults int, jql string, epicField string) {
	query := jiraservice.JiraSearchQuery{
		Jql:        jql,
		Fields:     []string{"summary", "issuetype", epicField},
		Expand:     []string{"changelog"},
		StartAt:    startAt,
		MaxResults: maxResults,
	}

	log.Printf("Querying Jira for startAt: %d; maxResults: %d", query.StartAt, query.MaxResults)
	searchResult, err := jc.jira.MakeJiraSearchRequest(&query)
	if err != nil {
		log.Fatalf("Failed to make request %s", err)
	}

	for _, issue := range searchResult.Issues {
		saveableIssue := models.Create(issue)
		go jc.repo.StoreIssue(context.Background(), saveableIssue)
	}

	var nextPageStartAt = util.NextPaginationArgs(startAt, maxResults, len(searchResult.Issues), searchResult.Total)
	if nextPageStartAt == -1 {
		log.Println("No new pages to fetch.")
		return
	}

	jc.execute(nextPageStartAt, maxResults, jql, epicField)
}
