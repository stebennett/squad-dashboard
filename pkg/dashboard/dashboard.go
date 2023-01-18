package dashboard

import (
	"context"
	"log"
	"time"

	"github.com/stebennett/squad-dashboard/pkg/dashboard/models"
	"github.com/stebennett/squad-dashboard/pkg/dateutil"
	"github.com/stebennett/squad-dashboard/pkg/jiracalculationsrepository"
	"github.com/stebennett/squad-dashboard/pkg/mathutil"
)

func GenerateEscapedDefects(weekCount int, project string, defectIssueType string, repo jiracalculationsrepository.JiraCalculationsRepository) (models.EscapedDefectReport, error) {
	// 1. Calculate dates of last weekCount fridays
	now := time.Now()

	nearestFriday := dateutil.NearestPreviousDateForDay(dateutil.AsDate(now.Year(), now.Month(), now.Day()), time.Friday)
	weekEndings := dateutil.PreviousWeekDates(nearestFriday, weekCount)

	escapedDefectCounts := []models.WeekCount{}
	var lastWeekEscapedDefect []models.EscapedDefectItem

	maxDate := weekEndings[0]

	// 2. Count the number of defects created for week
	for _, d := range weekEndings {
		startDate := d.AddDate(0, 0, -7)
		escapedDefects, err := repo.GetEscapedDefects(context.Background(), project, defectIssueType, startDate, d)
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
			lastWeekEscapedDefect = make([]models.EscapedDefectItem, len(escapedDefects))
			for i, d := range escapedDefects {
				lastWeekEscapedDefect[i] = models.EscapedDefectItem{
					IssueKey:  d.IssueKey,
					CreatedAt: d.CreatedAt,
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

func GenerateCycleTime(weekCount int, percentile float64, project string, issueTypes []string, repo jiracalculationsrepository.JiraCalculationsRepository) (models.CycleTimeReport, error) {
	// 1. Calculate dates of last weekCount fridays
	now := time.Now()

	nearestFriday := dateutil.NearestPreviousDateForDay(dateutil.AsDate(now.Year(), now.Month(), now.Day()), time.Friday)
	weekEndings := dateutil.PreviousWeekDates(nearestFriday, weekCount)
	maxDate := weekEndings[0]

	cycleTimeReports := []models.WeekCount{}
	cycleTimeValues := []models.CycleTimeItem{}
	var lastWeekCycleTimes []models.CycleTimeItem

	// 2. Get the average cycle time for a week
	for _, d := range weekEndings {
		startDate := d.AddDate(0, 0, -7)
		ct, err := repo.GetCompletedWorkingCycleTimes(context.Background(), project, issueTypes, startDate, d)
		if err != nil {
			return models.CycleTimeReport{
				WeeklyReports:          cycleTimeReports,
				AllCycleTimeItems:      cycleTimeValues,
				LastWeekCycleTimeItems: lastWeekCycleTimes,
			}, err
		}

		cycleTimes := make([]int, len(ct))
		for _, ctitem := range ct {
			cycleTimes = append(cycleTimes, ctitem.Size)
			cycleTimeValues = append(cycleTimeValues, models.CycleTimeItem{
				IssueKey:    ctitem.IssueKey,
				CompletedAt: ctitem.CompletedAt,
				Size:        ctitem.Size,
			})
		}

		cycleTimeReports = append(cycleTimeReports, models.WeekCount{
			WeekEnding: d,
			Count:      mathutil.Percentile(percentile, cycleTimes),
		})

		// get the values for the last week only
		if maxDate.Equal(d) {
			lastWeekCycleTimes = make([]models.CycleTimeItem, len(ct))
			for i, c := range ct {
				lastWeekCycleTimes[i] = models.CycleTimeItem{
					IssueKey:    c.IssueKey,
					CompletedAt: c.CompletedAt,
					Size:        c.Size,
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

func GenerateThroughput(weekCount int, project string, issueTypes []string, repo jiracalculationsrepository.JiraCalculationsRepository) (models.ThroughputReport, error) {
	// 1. Calculate dates of last weekCount fridays
	now := time.Now()

	nearestFriday := dateutil.NearestPreviousDateForDay(dateutil.AsDate(now.Year(), now.Month(), now.Day()), time.Friday)
	weekEndings := dateutil.PreviousWeekDates(nearestFriday, weekCount)
	maxDate := weekEndings[0]

	throughputReports := []models.WeekCount{}
	var lastWeekThroughputItems []models.ThroughputItem

	// 2. Get throughput by week
	for _, d := range weekEndings {
		startDate := d.AddDate(0, 0, -7)
		issues, err := repo.GetThroughput(context.Background(), project, issueTypes, startDate, d)
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
			lastWeekThroughputItems = make([]models.ThroughputItem, len(issues))
			for i, item := range issues {
				lastWeekThroughputItems[i] = models.ThroughputItem{
					IssueKey:    item.IssueKey,
					CompletedAt: item.CompletedAt,
				}
			}
		}
	}

	return models.ThroughputReport{
		WeeklyReports:           throughputReports,
		LastWeekThroughputItems: lastWeekThroughputItems,
	}, nil
}

func GenerateUnplannedWorkReport(weekCount int, project string, issueTypes []string, repo jiracalculationsrepository.JiraCalculationsRepository) (models.UnplannedWorkReport, error) {
	// 1. Calculate dates of last weekCount fridays
	now := time.Now()

	nearestFriday := dateutil.NearestPreviousDateForDay(dateutil.AsDate(now.Year(), now.Month(), now.Day()), time.Friday)
	weekEndings := dateutil.PreviousWeekDates(nearestFriday, weekCount)
	maxDate := weekEndings[0]

	unplannedWorkReports := []models.WeekCount{}
	var lastWeekUnplannedItems []models.UnplannedWorkItem

	for _, d := range weekEndings {
		startDate := d.AddDate(0, 0, -7)
		throughputIssues, err := repo.GetThroughput(context.Background(), project, issueTypes, startDate, d)
		if err != nil {
			return models.UnplannedWorkReport{
				WeeklyReports: unplannedWorkReports,
			}, err
		}

		unplannedIssues, err := repo.GetUnplannedThroughput(context.Background(), project, issueTypes, startDate, d)
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
			lastWeekUnplannedItems = make([]models.UnplannedWorkItem, len(unplannedIssues))
			for i, item := range unplannedIssues {
				lastWeekUnplannedItems[i] = models.UnplannedWorkItem{
					IssueKey:    item.IssueKey,
					CompletedAt: item.CompletedAt,
				}
			}
		}
	}

	return models.UnplannedWorkReport{
		WeeklyReports:              unplannedWorkReports,
		LastWeekUnplannedWorkItems: lastWeekUnplannedItems,
	}, nil
}
