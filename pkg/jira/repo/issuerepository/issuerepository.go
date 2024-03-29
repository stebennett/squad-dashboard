package issuerepository

import (
	"context"
	"time"

	jiramodels "github.com/stebennett/squad-dashboard/pkg/jira/models"
)

type IssueRepository interface {
	GetIssues(ctx context.Context, project string) ([]jiramodels.JiraIssue, error)
	SaveIssue(ctx context.Context, project string, jiraIssue jiramodels.JiraIssue) (int64, error)

	SaveTransition(ctx context.Context, issueKey string, jiraTransition []jiramodels.JiraTransition) (int64, error)
	SaveIssueLabels(ctx context.Context, issueKey string, label []string) (int64, error)

	GetTransitionsForIssue(ctx context.Context, issueKey string) ([]jiramodels.JiraTransition, error)
	GetTransitionTimeByStateChanges(ctx context.Context, project string, fromStates []string, toStates []string) (map[string]time.Time, error)
	GetTransitionTimeByToState(ctx context.Context, project string, toStates []string) (map[string]time.Time, error)

	GetEndStateForIssue(ctx context.Context, issueKey string, transitionDate time.Time) (string, error)

	SetIssuesStartedInWeekStarting(ctx context.Context, project string, startDate time.Time, count int) (int64, error)
	SetIssuesCompletedInWeekStarting(ctx context.Context, project string, startDate time.Time, count int) (int64, error)
	GetProjects(ctx context.Context) ([]string, error)

	ClearUnplannedIssuesForProject(ctx context.Context, project string) (int64, error)
	SaveUnplannedIssue(ctx context.Context, issueKey string) (int64, error)
}
