package jirastatsservice

import (
	"context"
	"time"

	"github.com/stebennett/squad-dashboard/pkg/dateutil"
	"github.com/stebennett/squad-dashboard/pkg/jiramodels"
	"github.com/stebennett/squad-dashboard/pkg/jirarepository"
)

type JiraStatsService struct {
	JiraRepository jirarepository.JiraRepository
}

func (jss JiraStatsService) FetchThrougputDataForProject(project string) ([]jiramodels.ThroughputItem, error) {
	return jss.JiraRepository.GetWeeklyThroughputByProject(context.Background(), project, dateutil.NearestPreviousDateForDay(time.Now(), time.Friday), 12)
}
