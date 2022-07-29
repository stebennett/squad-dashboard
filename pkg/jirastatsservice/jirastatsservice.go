package jirastatsservice

import (
	"context"
	"time"

	"github.com/stebennett/squad-dashboard/pkg/dateutil"
	"github.com/stebennett/squad-dashboard/pkg/jirarepository"
	"github.com/stebennett/squad-dashboard/pkg/statsmodels"
	"github.com/stebennett/squad-dashboard/pkg/statsutil"
)

type JiraStatsService struct {
	JiraRepository jirarepository.JiraRepository
}

func (jss JiraStatsService) FetchThrougputDataForProject(project string) (statsmodels.ThrouputResult, error) {
	throughputItems, err := jss.JiraRepository.GetWeeklyThroughputByProject(context.Background(), project, dateutil.NearestPreviousDateForDay(time.Now(), time.Friday), 12)
	if err != nil {
		return statsmodels.ThrouputResult{}, err
	}

	trendline, err := statsutil.CalculateTrendline(throughputItems)
	if err != nil {
		return statsmodels.ThrouputResult{}, err
	}

	return statsmodels.ThrouputResult{
		Project:         project,
		ThroughputItems: throughputItems,
		Trendline:       trendline,
	}, nil
}

func (jss JiraStatsService) FetchThrougputDataForAllProjects() ([]statsmodels.ThrouputResult, error) {
	projectThroughputItems, err := jss.JiraRepository.GetWeeklyThroughputAllProjects(context.Background(), dateutil.NearestPreviousDateForDay(time.Now(), time.Friday), 12)
	if err != nil {
		return []statsmodels.ThrouputResult{}, err
	}

	throughtputMap := map[string][]statsmodels.WeeklyTimeItem{}
	for _, pti := range projectThroughputItems {
		items := throughtputMap[pti.Project]
		items = append(items, pti.TimeItem)
		throughtputMap[pti.Project] = items
	}

	results := []statsmodels.ThrouputResult{}
	for k, v := range throughtputMap {
		trendline, err := statsutil.CalculateTrendline(v)
		if err != nil {
			return []statsmodels.ThrouputResult{}, err
		}
		tr := statsmodels.ThrouputResult{
			Project:         k,
			ThroughputItems: v,
			Trendline:       trendline,
		}
		results = append(results, tr)
	}
	return results, nil
}
