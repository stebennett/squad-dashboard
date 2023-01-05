package printer

import "github.com/stebennett/squad-dashboard/pkg/dashboard/models"

type Printer interface {
	PrintDefectCounts(defectCounts []models.WeekCount) error
	PrintCycleTimes(cycleTimeReports []models.WeekCount) error
	PrintThroughput(throughputReports []models.WeekCount) error
}
