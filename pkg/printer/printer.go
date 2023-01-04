package printer

import "github.com/stebennett/squad-dashboard/pkg/dashboard/models"

type Printer interface {
	PrintDefectCounts(defectCounts []models.EscapedDefectCount) error
	PrintCycleTimes(cycleTimeReports []models.CycleTimeReport) error
}
