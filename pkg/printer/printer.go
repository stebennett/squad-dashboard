package printer

import (
	"github.com/stebennett/squad-dashboard/pkg/models"
)

type Reports struct {
	EscapedDefects       models.EscapedDefectReport
	CycleTimeReports     models.CycleTimeReport
	ThroughputReports    models.ThroughputReport
	UnplannedWorkReports models.UnplannedWorkReport
}

type Printer interface {
	Print(reports Reports) error
}
