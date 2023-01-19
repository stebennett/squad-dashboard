package calculationsrepository

import (
	"context"
	"time"

	"github.com/stebennett/squad-dashboard/pkg/jira/models"
)

type CycleTimes struct {
	IssueKey    string
	CompletedAt time.Time
	Size        int
}

type EscapedDefect struct {
	IssueKey  string
	CreatedAt time.Time
}

type ThroughputIssue struct {
	IssueKey    string
	CompletedAt time.Time
}

type JiraCalculationsRepository interface {
	DropAllCalculations(ctx context.Context, project string) (int64, error)

	SaveCreateDates(ctx context.Context, issueKey string, year int, week int, createdAt time.Time) (int64, error)
	SaveStartDates(ctx context.Context, issueKey string, year int, week int, startedAt time.Time) (int64, error)
	SaveCompleteDates(ctx context.Context, issueKey string, year int, week int, completedAt time.Time, endState string) (int64, error)

	SaveCycleTime(ctx context.Context, issueKey string, cycleTime int, workingCycleTime int) (int64, error)
	SaveLeadTime(ctx context.Context, issueKey string, leadTime int, workingLeadTime int) (int64, error)
	SaveSystemDelayTime(ctx context.Context, issueKey string, systemDelayTime int, workingSystemDelayTime int) (int64, error)

	GetEscapedDefects(ctx context.Context, project string, issueType string, startDate time.Time, endDate time.Time) ([]EscapedDefect, error)
	GetCompletedWorkingCycleTimes(ctx context.Context, project string, issueTypes []string, startDate time.Time, endDate time.Time) ([]CycleTimes, error)
	GetThroughput(ctx context.Context, project string, issueTypes []string, startDate time.Time, endDate time.Time) ([]ThroughputIssue, error)
	GetUnplannedThroughput(ctx context.Context, project string, issueTypes []string, startDate time.Time, endDate time.Time) ([]ThroughputIssue, error)
	GetCompletedIssues(ctx context.Context, project string) (map[string]models.IssueCalculations, error)

	GetIssuesStartedBetweenDates(ctx context.Context, project string, startDate time.Time, endDate time.Time, issueTypes []string) ([]string, error)
	GetIssuesCompletedBetweenDates(ctx context.Context, project string, startDate time.Time, endDate time.Time, issueTypes []string, endStates []string) ([]string, error)
}
