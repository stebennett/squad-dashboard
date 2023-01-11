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

func GenerateEscapedDefects(weekCount int, project string, defectIssueType string, repo jiracalculationsrepository.JiraCalculationsRepository) ([]models.WeekCount, error) {
	// 1. Calculate dates of last weekCount fridays
	now := time.Now()

	nearestFriday := dateutil.NearestPreviousDateForDay(dateutil.AsDate(now.Year(), now.Month(), now.Day()), time.Friday)
	weekEndings := dateutil.PreviousWeekDates(nearestFriday, weekCount)

	escapedDefectCounts := []models.WeekCount{}

	// 2. Count the number of defects created for week
	for _, d := range weekEndings {
		startDate := d.AddDate(0, 0, -7)
		escapedDefects, err := repo.GetEscapedDefects(context.Background(), project, defectIssueType, startDate, d)
		if err != nil {
			return escapedDefectCounts, err
		}

		escapedDefectCounts = append(escapedDefectCounts, models.WeekCount{
			WeekEnding: d,
			Count:      len(escapedDefects),
		})
	}

	// 3. Return the data
	return escapedDefectCounts, nil
}

func GenerateCycleTime(weekCount int, percentile float64, project string, issueTypes []string, repo jiracalculationsrepository.JiraCalculationsRepository) ([]models.WeekCount, []jiracalculationsrepository.CycleTimes, error) {
	// 1. Calculate dates of last weekCount fridays
	now := time.Now()

	nearestFriday := dateutil.NearestPreviousDateForDay(dateutil.AsDate(now.Year(), now.Month(), now.Day()), time.Friday)
	weekEndings := dateutil.PreviousWeekDates(nearestFriday, weekCount)

	cycleTimeReports := []models.WeekCount{}
	cycleTimeValues := []jiracalculationsrepository.CycleTimes{}

	// 2. Get the average cycle time for a week
	for _, d := range weekEndings {
		startDate := d.AddDate(0, 0, -7)
		ct, err := repo.GetCompletedWorkingCycleTimes(context.Background(), project, issueTypes, startDate, d)
		if err != nil {
			return cycleTimeReports, cycleTimeValues, err
		}

		cycleTimes := make([]int, len(ct))
		for _, ctitem := range ct {
			cycleTimes = append(cycleTimes, ctitem.Size)
		}

		cycleTimeReports = append(cycleTimeReports, models.WeekCount{
			WeekEnding: d,
			Count:      mathutil.Percentile(percentile, cycleTimes),
		})

		cycleTimeValues = append(cycleTimeValues, ct...)
	}

	return cycleTimeReports, cycleTimeValues, nil
}

func GenerateThroughput(weekCount int, project string, issueTypes []string, repo jiracalculationsrepository.JiraCalculationsRepository) ([]models.WeekCount, error) {
	// 1. Calculate dates of last weekCount fridays
	now := time.Now()

	nearestFriday := dateutil.NearestPreviousDateForDay(dateutil.AsDate(now.Year(), now.Month(), now.Day()), time.Friday)
	weekEndings := dateutil.PreviousWeekDates(nearestFriday, weekCount)

	throughputReports := []models.WeekCount{}

	// 2. Get throughput by week
	for _, d := range weekEndings {
		startDate := d.AddDate(0, 0, -7)
		issues, err := repo.GetThroughput(context.Background(), project, issueTypes, startDate, d)
		if err != nil {
			return throughputReports, err
		}

		throughputReports = append(throughputReports, models.WeekCount{
			WeekEnding: d,
			Count:      len(issues),
		})
	}

	return throughputReports, nil
}

func GenerateUnplannedWorkReport(weekCount int, project string, issueTypes []string, repo jiracalculationsrepository.JiraCalculationsRepository) ([]models.WeekCount, error) {
	// 1. Calculate dates of last weekCount fridays
	now := time.Now()

	nearestFriday := dateutil.NearestPreviousDateForDay(dateutil.AsDate(now.Year(), now.Month(), now.Day()), time.Friday)
	weekEndings := dateutil.PreviousWeekDates(nearestFriday, weekCount)

	unplannedWorkReports := []models.WeekCount{}

	for _, d := range weekEndings {
		startDate := d.AddDate(0, 0, -7)
		throughputIssues, err := repo.GetThroughput(context.Background(), project, issueTypes, startDate, d)
		if err != nil {
			return unplannedWorkReports, err
		}

		unplannedIssues, err := repo.GetUnplannedThroughput(context.Background(), project, issueTypes, startDate, d)
		if err != nil {
			return unplannedWorkReports, err
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
	}

	return unplannedWorkReports, nil
}
