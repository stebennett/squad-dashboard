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

	trendline, err := statsutil.CalculateTrendlineForWeeklyTimeItems(throughputItems)
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
		trendline, err := statsutil.CalculateTrendlineForWeeklyTimeItems(v)
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

func (jss JiraStatsService) FetchCycleTimeDataForProject(project string) (statsmodels.CycleTimeResult, error) {
	cycleTimeItems, err := jss.JiraRepository.GetWeeklyCycleTimeByProject(context.Background(), project, dateutil.NearestPreviousDateForDay(time.Now(), time.Friday), 12)
	if err != nil {
		return statsmodels.CycleTimeResult{}, err
	}

	trendline, err := statsutil.CalculateTrendlineForCycleTimes(cycleTimeItems)
	if err != nil {
		return statsmodels.CycleTimeResult{}, err
	}

	return statsmodels.CycleTimeResult{
		Project:        project,
		CycleTimeItems: cycleTimeItems,
		Trendline:      trendline,
	}, nil
}

func (jss JiraStatsService) FetchCycleTimeDataForAllProjects() ([]statsmodels.CycleTimeResult, error) {
	projectCycleTimeItems, err := jss.JiraRepository.GetWeeklyCycleTimeAllProjects(context.Background(), dateutil.NearestPreviousDateForDay(time.Now(), time.Friday), 12)
	if err != nil {
		return []statsmodels.CycleTimeResult{}, err
	}

	cycleTimeMap := map[string][]statsmodels.WeeklyCycleTimeItem{}
	for _, cti := range projectCycleTimeItems {
		items := cycleTimeMap[cti.Project]
		items = append(items, cti.TimeItem)
		cycleTimeMap[cti.Project] = items
	}

	results := []statsmodels.CycleTimeResult{}
	for k, v := range cycleTimeMap {
		trendline, err := statsutil.CalculateTrendlineForCycleTimes(v)
		if err != nil {
			return []statsmodels.CycleTimeResult{}, err
		}
		tr := statsmodels.CycleTimeResult{
			Project:        k,
			CycleTimeItems: v,
			Trendline:      trendline,
		}
		results = append(results, tr)
	}
	return results, nil
}
