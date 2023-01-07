package jiraconfigrepository

import (
	"context"
	"time"
)

type ConfigRepository interface {
	SaveJiraToDoStates(ctx context.Context, project string, states []string) (int64, error)
	SaveJiraInProgressStates(ctx context.Context, project string, states []string) (int64, error)
	SaveJiraDoneStates(ctx context.Context, project string, states []string) (int64, error)

	SaveNonWorkingDays(ctx context.Context, project string, nonWorkingDays []string) (int64, error)
	GetNonWorkingDays(ctx context.Context, project string) ([]time.Time, error)
}
