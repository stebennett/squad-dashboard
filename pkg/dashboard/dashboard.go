package dashboard

import (
	"context"
	"time"

	"github.com/stebennett/squad-dashboard/pkg/dashboard/models"
	"github.com/stebennett/squad-dashboard/pkg/dateutil"
	"github.com/stebennett/squad-dashboard/pkg/jirarepository"
	"github.com/stebennett/squad-dashboard/pkg/mathutil"
)

func GenerateEscapedDefects(weekCount int, project string, defectIssueType string, repo jirarepository.JiraRepository) ([]models.EscapedDefectCount, error) {
	// 1. Calculate dates of last weekCount fridays
	now := time.Now()

	nearestFriday := dateutil.NearestPreviousDateForDay(dateutil.AsDate(now.Year(), now.Month(), now.Day()), time.Friday)
	weekEndings := dateutil.PreviousWeekDates(nearestFriday, weekCount)

	escapedDefectCounts := []models.EscapedDefectCount{}

	// 2. Count the number of defects created for week
	for _, d := range weekEndings {
		startDate := d.AddDate(0, 0, -7)
		escapedDefects, err := repo.GetEscapedDefects(context.Background(), project, defectIssueType, startDate, d)
		if err != nil {
			return escapedDefectCounts, err
		}

		escapedDefectCounts = append(escapedDefectCounts, models.EscapedDefectCount{
			WeekEnding:             d,
			NumberOfDefectsCreated: len(escapedDefects),
		})
	}

	// 3. Return the data
	return escapedDefectCounts, nil
}

func GenerateCycleTime(weekCount int, project string, issueTypes []string, repo jirarepository.JiraRepository) ([]models.CycleTimeReport, error) {
	// 1. Calculate dates of last weekCount fridays
	now := time.Now()

	nearestFriday := dateutil.NearestPreviousDateForDay(dateutil.AsDate(now.Year(), now.Month(), now.Day()), time.Friday)
	weekEndings := dateutil.PreviousWeekDates(nearestFriday, weekCount)

	cycleTimeReports := []models.CycleTimeReport{}

	// 2. Get the average cycle time for a week
	for _, d := range weekEndings {
		startDate := d.AddDate(0, 0, -7)
		cycleTimes, err := repo.GetCompletedWorkingCycleTimes(context.Background(), project, issueTypes, startDate, d)
		if err != nil {
			return cycleTimeReports, err
		}

		cycleTimeReports = append(cycleTimeReports, models.CycleTimeReport{
			WeekEnding:       d,
			AverageCycleTime: mathutil.Percentile(0.75, cycleTimes),
		})
	}

	return cycleTimeReports, nil
}
