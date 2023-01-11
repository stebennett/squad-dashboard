package jiracollector

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/stebennett/squad-dashboard/cmd/jiracollector/models"
	"github.com/stebennett/squad-dashboard/pkg/jiramodels"
	"github.com/stebennett/squad-dashboard/pkg/jirarepository"
	"github.com/stebennett/squad-dashboard/pkg/jiraservice"
	"github.com/stebennett/squad-dashboard/pkg/paginator"
)

type JiraIssueCollector struct {
	repo      jirarepository.JiraRepository
	jira      *jiraservice.JiraService
	epicField string
}

func NewJiraIssueCollector(jira *jiraservice.JiraService, repo jirarepository.JiraRepository, epicField string) *JiraIssueCollector {
	return &JiraIssueCollector{
		repo:      repo,
		jira:      jira,
		epicField: epicField,
	}
}

func (jc *JiraIssueCollector) Execute(project string, jql string) error {
	return jc.execute(project, 0, 100, jql, jc.epicField, jc.repo.SaveIssue, jc.repo.SaveTransition, jc.repo.SaveIssueLabels)
}

func (jc *JiraIssueCollector) execute(project string, startAt int, maxResults int, jql string, epicField string,
	saveIssue func(ctx context.Context, project string, jiraIssue jiramodels.JiraIssue) (int64, error),
	saveTransition func(ctx context.Context, issueKey string, jiraTransitions []jiramodels.JiraTransition) (int64, error),
	saveIssueLabels func(ctx context.Context, issueKey string, labels []string) (int64, error)) error {

	query := jiraservice.JiraSearchQuery{
		Jql:        jql,
		Fields:     []string{"summary", "issuetype", epicField, "created", "updated", "labels"},
		Expand:     []string{"changelog"},
		StartAt:    startAt,
		MaxResults: maxResults,
	}

	log.Printf("Querying Jira for startAt: %d; maxResults: %d", query.StartAt, query.MaxResults)
	searchResult, err := jc.jira.MakeJiraSearchRequest(&query)
	if err != nil {
		return err
	}

	var jiraResult models.JiraSearchResults
	jsonErr := json.Unmarshal([]byte(strings.ReplaceAll(searchResult, epicField, "epicKey")), &jiraResult)
	if jsonErr != nil {
		return jsonErr
	}

	for _, issue := range jiraResult.Issues {
		saveableIssue, err := models.Create(issue)
		if err != nil {
			return fmt.Errorf("error: failed to covert issue %s - %s", issue.Key, err)
		}

		_, err = saveIssue(context.Background(), project, *saveableIssue)
		if err != nil {
			return fmt.Errorf("error: failed to save issue %s - %s", saveableIssue.Key, err)
		}

		_, err = saveIssueLabels(context.Background(), saveableIssue.Key, saveableIssue.Labels)
		if err != nil {
			return fmt.Errorf("error: failed to save issue labels %s, %s - %s", saveableIssue.Key, saveableIssue.Labels, err)
		}

		transitions, err := jc.fetchTransitions(issue)
		if err != nil {
			return fmt.Errorf("error: failed to get issue transitions for %s - %s", saveableIssue.Key, err)
		}

		_, err = saveTransition(context.Background(), saveableIssue.Key, transitions)
		if err != nil {
			return fmt.Errorf("error: failed to save issue transitions - %s", err)
		}
	}

	var nextPageStartAt = paginator.NextPaginationArgs(startAt, maxResults, len(jiraResult.Issues), jiraResult.Total)
	if nextPageStartAt == -1 {
		log.Println("No new pages to fetch.")
		return nil
	}

	return jc.execute(project, nextPageStartAt, maxResults, jql, epicField, saveIssue, saveTransition, saveIssueLabels)
}

func (jc *JiraIssueCollector) fetchTransitions(jiraIssue models.JiraResultIssue) ([]jiramodels.JiraTransition, error) {
	if jiraIssue.ChangeLog.Total > jiraIssue.ChangeLog.MaxResults {
		issueHistory, err := jc.fetchTransitionsFromIssue(jiraIssue.Key, 0, 100)
		if err != nil {
			return []jiramodels.JiraTransition{}, err
		}

		return models.CreateTransitions(issueHistory), nil
	}

	return models.CreateTransitions(jiraIssue.ChangeLog.Histories), nil
}

func (jc *JiraIssueCollector) fetchTransitionsFromIssue(issueKey string, startAt int, maxResults int) ([]models.JiraHistory, error) {
	if startAt == -1 {
		return []models.JiraHistory{}, nil
	}

	pagedHistory, err := jc.jira.MakeJiraGetHistoryRequest(issueKey, startAt, maxResults)
	if err != nil {
		return []models.JiraHistory{}, err
	}

	var jiraResult models.JiraChangeLogHistoryResults
	err = json.Unmarshal([]byte(pagedHistory), &jiraResult)
	if err != nil {
		return []models.JiraHistory{}, err
	}

	nextStartAt := paginator.NextPaginationArgs(jiraResult.StartAt, maxResults, len(jiraResult.Histories), jiraResult.Total)

	nextResults, err := jc.fetchTransitionsFromIssue(issueKey, nextStartAt, maxResults)
	if err != nil {
		return []models.JiraHistory{}, err
	}

	return append(nextResults, jiraResult.Histories...), nil
}
