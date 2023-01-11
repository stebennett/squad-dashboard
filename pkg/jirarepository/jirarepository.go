package jirarepository

import (
	"context"
	"time"

	"github.com/stebennett/squad-dashboard/pkg/calculatormodels"
	"github.com/stebennett/squad-dashboard/pkg/jiramodels"
)

type JiraRepository interface {
	GetIssues(ctx context.Context, project string) ([]jiramodels.JiraIssue, error)
	SaveIssue(ctx context.Context, project string, jiraIssue jiramodels.JiraIssue) (int64, error)
	SaveTransition(ctx context.Context, issueKey string, jiraTransition []jiramodels.JiraTransition) (int64, error)
	GetTransitionTimeByStateChanges(ctx context.Context, project string, fromStates []string, toStates []string) (map[string]time.Time, error)
	GetTransitionTimeByToState(ctx context.Context, project string, toStates []string) (map[string]time.Time, error)
	GetTransitionsForIssue(ctx context.Context, issueKey string) ([]jiramodels.JiraTransition, error)
	GetCompletedIssues(ctx context.Context, project string) (map[string]calculatormodels.IssueCalculations, error)
	GetIssuesStartedBetweenDates(ctx context.Context, project string, startDate time.Time, endDate time.Time, issueTypes []string) ([]string, error)
	SetIssuesStartedInWeekStarting(ctx context.Context, project string, startDate time.Time, count int) (int64, error)
	GetIssuesCompletedBetweenDates(ctx context.Context, project string, startDate time.Time, endDate time.Time, issueTypes []string, endStates []string) ([]string, error)
	SetIssuesCompletedInWeekStarting(ctx context.Context, project string, startDate time.Time, count int) (int64, error)
	GetEndStateForIssue(ctx context.Context, issueKey string, transitionDate time.Time) (string, error)
	SaveIssueLabels(ctx context.Context, issueKey string, label []string) (int64, error)
	GetProjects(ctx context.Context) ([]string, error)

	ClearUnplannedIssuesForProject(ctx context.Context, project string) (int64, error)
	SaveUnplannedIssue(ctx context.Context, issueKey string) (int64, error)
}
