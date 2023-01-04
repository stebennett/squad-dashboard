package dashboard

import (
	"context"
	"log"
	"time"

	"github.com/stebennett/squad-dashboard/pkg/dashboard/models"
	"github.com/stebennett/squad-dashboard/pkg/dateutil"
	"github.com/stebennett/squad-dashboard/pkg/jirarepository"
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
		log.Printf("found: %s", escapedDefects)
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
