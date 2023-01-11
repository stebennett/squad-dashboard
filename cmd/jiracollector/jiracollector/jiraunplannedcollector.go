package jiracollector

import (
	"context"
	"encoding/json"

	"github.com/stebennett/squad-dashboard/cmd/jiracollector/models"
	"github.com/stebennett/squad-dashboard/pkg/jirarepository"
	"github.com/stebennett/squad-dashboard/pkg/jiraservice"
	"github.com/stebennett/squad-dashboard/pkg/paginator"
)

type JiraUnplannedCollector struct {
	repo jirarepository.JiraRepository
	jira *jiraservice.JiraService
}

func NewJiraUnplannedCollector(jira *jiraservice.JiraService, repo jirarepository.JiraRepository) *JiraUnplannedCollector {
	return &JiraUnplannedCollector{
		repo: repo,
		jira: jira,
	}
}

func (jc *JiraUnplannedCollector) Execute(project string, jql string) error {
	jc.repo.ClearUnplannedIssuesForProject(context.Background(), project)
	return jc.execute(project, 0, 100, jql, jc.repo.SaveUnplannedIssue)
}

func (jc *JiraUnplannedCollector) execute(project string, startAt int, maxResults int, jql string, saveUnplannedIssues func(ctx context.Context, issueKey string) (int64, error)) error {
	query := jiraservice.JiraSearchQuery{
		Jql:        jql,
		StartAt:    startAt,
		MaxResults: maxResults,
	}

	searchResult, err := jc.jira.MakeJiraSearchRequest(&query)
	if err != nil {
		return err
	}

	var jiraResult models.JiraSearchResults
	jsonErr := json.Unmarshal([]byte(searchResult), &jiraResult)
	if jsonErr != nil {
		return jsonErr
	}

	for _, issue := range jiraResult.Issues {
		_, err = saveUnplannedIssues(context.Background(), issue.Key)
		if err != nil {
			return err
		}
	}

	var nextPageStartAt = paginator.NextPaginationArgs(startAt, maxResults, len(jiraResult.Issues), jiraResult.Total)
	if nextPageStartAt == -1 {
		return nil
	}

	return nil
}
