package statsservice

import (
	"context"
	"log"
	"time"

	"github.com/stebennett/squad-dashboard/pkg/dateutil"
	"github.com/stebennett/squad-dashboard/pkg/jira/repo/calculationsrepository"
	"github.com/stebennett/squad-dashboard/pkg/mathutil"
	"github.com/stebennett/squad-dashboard/pkg/models"
)

type StatsService struct {
	calculationsRepository calculationsrepository.JiraCalculationsRepository
}

func NewStatsService(repo calculationsrepository.JiraCalculationsRepository) *StatsService {
	return &StatsService{
		calculationsRepository: repo,
	}
}

func (ss *StatsService) GenerateEscapedDefects(weekCount int, project string, defectIssueType string, startTime time.Time, endOfWeekDay time.Weekday) (models.EscapedDefectReport, error) {
	// 1. Calculate dates of last weekCount fridays
	nearestFriday := dateutil.NearestPreviousDateForDay(dateutil.AsDate(startTime.Year(), startTime.Month(), startTime.Day()), endOfWeekDay)
	weekEndings := dateutil.PreviousWeekDates(nearestFriday, weekCount)

	escapedDefectCounts := []models.WeekCount{}
	var lastWeekEscapedDefect []models.WorkItem

	maxDate := weekEndings[0]

	// 2. Count the number of defects created for week
	for _, d := range weekEndings {
		startDate := d.AddDate(0, 0, -7)
		escapedDefects, err := ss.calculationsRepository.GetEscapedDefects(context.Background(), project, defectIssueType, startDate, d)
		if err != nil {
			return models.EscapedDefectReport{
				WeeklyReports: escapedDefectCounts,
			}, err
		}

		escapedDefectCounts = append(escapedDefectCounts, models.WeekCount{
			WeekEnding: d,
			Count:      len(escapedDefects),
		})

		// get the values for the last week only
		if maxDate.Equal(d) {
			lastWeekEscapedDefect = make([]models.WorkItem, len(escapedDefects))
			for i, d := range escapedDefects {
				lastWeekEscapedDefect[i] = models.WorkItem{
					IssueKey:         d.IssueKey,
					CreatedAt:        d.IssueCreatedAt.Time,
					CompletedAt:      d.IssueCompletedAt.Time,
					WorkingCycleTime: d.WorkingCycleTime,
				}
			}
		}
	}

	// 3. Return the data
	return models.EscapedDefectReport{
		WeeklyReports:              escapedDefectCounts,
		LastWeekEscapedDefectItems: lastWeekEscapedDefect,
	}, nil
}

func (ss *StatsService) GenerateCycleTime(weekCount int, percentile float64, project string, issueTypes []string, startTime time.Time, endOfWeekDay time.Weekday) (models.CycleTimeReport, error) {
	// 1. Calculate dates of last weekCount fridays
	nearestFriday := dateutil.NearestPreviousDateForDay(dateutil.AsDate(startTime.Year(), startTime.Month(), startTime.Day()), endOfWeekDay)
	weekEndings := dateutil.PreviousWeekDates(nearestFriday, weekCount)
	maxDate := weekEndings[0]

	cycleTimeReports := []models.WeekCount{}
	cycleTimeValues := []models.WorkItem{}
	var lastWeekCycleTimes []models.WorkItem

	// 2. Get the average cycle time for a week
	for _, d := range weekEndings {
		startDate := d.AddDate(0, 0, -7)
		ct, err := ss.calculationsRepository.GetCompletedWorkingCycleTimes(context.Background(), project, issueTypes, startDate, d)
		if err != nil {
			return models.CycleTimeReport{
				WeeklyReports:          cycleTimeReports,
				AllCycleTimeItems:      cycleTimeValues,
				LastWeekCycleTimeItems: lastWeekCycleTimes,
			}, err
		}

		cycleTimes := make([]int, len(ct))
		for _, ctitem := range ct {
			cycleTimes = append(cycleTimes, ctitem.WorkingCycleTime)
			cycleTimeValues = append(cycleTimeValues, models.WorkItem{
				IssueKey:         ctitem.IssueKey,
				CreatedAt:        ctitem.IssueCreatedAt.Time,
				CompletedAt:      ctitem.IssueCompletedAt.Time,
				WorkingCycleTime: ctitem.WorkingCycleTime,
			})
		}

		cycleTimeReports = append(cycleTimeReports, models.WeekCount{
			WeekEnding: d,
			Count:      mathutil.Percentile(percentile, cycleTimes),
		})

		// get the values for the last week only
		if maxDate.Equal(d) {
			lastWeekCycleTimes = make([]models.WorkItem, len(ct))
			for i, c := range ct {
				lastWeekCycleTimes[i] = models.WorkItem{
					IssueKey:         c.IssueKey,
					CreatedAt:        c.IssueCreatedAt.Time,
					CompletedAt:      c.IssueCompletedAt.Time,
					WorkingCycleTime: c.WorkingCycleTime,
				}
			}
		}
	}

	return models.CycleTimeReport{
		WeeklyReports:          cycleTimeReports,
		AllCycleTimeItems:      cycleTimeValues,
		LastWeekCycleTimeItems: lastWeekCycleTimes,
	}, nil
}

func (ss *StatsService) GenerateThroughput(weekCount int, project string, issueTypes []string, startTime time.Time, endOfWeekDay time.Weekday) (models.ThroughputReport, error) {
	// 1. Calculate dates of last weekCount fridays
	nearestFriday := dateutil.NearestPreviousDateForDay(dateutil.AsDate(startTime.Year(), startTime.Month(), startTime.Day()), time.Friday)
	weekEndings := dateutil.PreviousWeekDates(nearestFriday, weekCount)
	maxDate := weekEndings[0]

	throughputReports := []models.WeekCount{}
	var lastWeekThroughputItems []models.WorkItem

	// 2. Get throughput by week
	for _, d := range weekEndings {
		startDate := d.AddDate(0, 0, -7)
		issues, err := ss.calculationsRepository.GetThroughput(context.Background(), project, issueTypes, startDate, d)
		if err != nil {
			return models.ThroughputReport{
				WeeklyReports: throughputReports,
			}, nil
		}

		throughputReports = append(throughputReports, models.WeekCount{
			WeekEnding: d,
			Count:      len(issues),
		})

		// get the values for the last week only
		if maxDate.Equal(d) {
			lastWeekThroughputItems = make([]models.WorkItem, len(issues))
			for i, item := range issues {
				lastWeekThroughputItems[i] = models.WorkItem{
					IssueKey:         item.IssueKey,
					CreatedAt:        item.IssueCreatedAt.Time,
					CompletedAt:      item.IssueCompletedAt.Time,
					WorkingCycleTime: item.WorkingCycleTime,
				}
			}
		}
	}

	return models.ThroughputReport{
		WeeklyReports:           throughputReports,
		LastWeekThroughputItems: lastWeekThroughputItems,
	}, nil
}

func (ss *StatsService) GenerateUnplannedWorkReport(weekCount int, project string, issueTypes []string, startTime time.Time, endOfWeekDay time.Weekday) (models.UnplannedWorkReport, error) {
	// 1. Calculate dates of last weekCount fridays
	nearestFriday := dateutil.NearestPreviousDateForDay(dateutil.AsDate(startTime.Year(), startTime.Month(), startTime.Day()), endOfWeekDay)
	weekEndings := dateutil.PreviousWeekDates(nearestFriday, weekCount)
	maxDate := weekEndings[0]

	unplannedWorkReports := []models.WeekCount{}
	var lastWeekUnplannedItems []models.WorkItem

	for _, d := range weekEndings {
		startDate := d.AddDate(0, 0, -7)
		throughputIssues, err := ss.calculationsRepository.GetThroughput(context.Background(), project, issueTypes, startDate, d)
		if err != nil {
			return models.UnplannedWorkReport{
				WeeklyReports: unplannedWorkReports,
			}, err
		}

		unplannedIssues, err := ss.calculationsRepository.GetUnplannedThroughput(context.Background(), project, issueTypes, startDate, d)
		if err != nil {
			return models.UnplannedWorkReport{
				WeeklyReports: unplannedWorkReports,
			}, err
		}

		unplannedPercent := 0
		if len(throughputIssues) > 0 && len(unplannedIssues) > 0 {
			unplannedPercent = int((float64(len(unplannedIssues)) / float64(len(throughputIssues))) * 100.0)
		}

		log.Printf("unplanned: %d; all: %d; percent: %d", len(unplannedIssues), len(throughputIssues), unplannedPercent)

		unplannedWorkReports = append(unplannedWorkReports, models.WeekCount{
			WeekEnding: d,
			Count:      unplannedPercent,
		})

		// get the values for the last week only
		if maxDate.Equal(d) {
			lastWeekUnplannedItems = make([]models.WorkItem, len(unplannedIssues))
			for i, item := range unplannedIssues {
				lastWeekUnplannedItems[i] = models.WorkItem{
					IssueKey:         item.IssueKey,
					CreatedAt:        item.IssueCreatedAt.Time,
					CompletedAt:      item.IssueCompletedAt.Time,
					WorkingCycleTime: item.WorkingCycleTime,
				}
			}
		}
	}

	return models.UnplannedWorkReport{
		WeeklyReports:              unplannedWorkReports,
		LastWeekUnplannedWorkItems: lastWeekUnplannedItems,
	}, nil
}
