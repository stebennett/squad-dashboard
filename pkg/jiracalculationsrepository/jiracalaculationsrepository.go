package jiracalculationsrepository

import (
	"context"
	"time"
)

type JiraCalculationsRepository interface {
	SaveCreateDates(ctx context.Context, issueKey string, year int, week int, createdAt time.Time) (int64, error)
	SaveStartDates(ctx context.Context, issueKey string, year int, week int, startedAt time.Time) (int64, error)
	SaveCompleteDates(ctx context.Context, issueKey string, year int, week int, completedAt time.Time, endState string) (int64, error)

	SaveCycleTime(ctx context.Context, issueKey string, cycleTime int, workingCycleTime int) (int64, error)
	SaveLeadTime(ctx context.Context, issueKey string, leadTime int, workingLeadTime int) (int64, error)
	SaveSystemDelayTime(ctx context.Context, issueKey string, systemDelayTime int, workingSystemDelayTime int) (int64, error)

	GetEscapedDefects(ctx context.Context, project string, issueType string, startDate time.Time, endDate time.Time) ([]string, error)
	GetCompletedWorkingCycleTimes(ctx context.Context, project string, issueTypes []string, startDate time.Time, endDate time.Time) ([]int, error)
	GetThroughput(ctx context.Context, project string, issueTypes []string, startDate time.Time, endDate time.Time) ([]string, error)
}
