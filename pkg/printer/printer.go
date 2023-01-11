package printer

import (
	"github.com/stebennett/squad-dashboard/pkg/dashboard/models"
	"github.com/stebennett/squad-dashboard/pkg/jiracalculationsrepository"
)

type Reports struct {
	EscapedDefects       []models.WeekCount
	CycleTimeReports     []models.WeekCount
	ThroughputReports    []models.WeekCount
	UnplannedWorkReports []models.WeekCount
	AllCycleTimes        []jiracalculationsrepository.CycleTimes
}

type Printer interface {
	Print(reports Reports) error
}
