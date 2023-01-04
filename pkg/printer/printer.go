package printer

import "github.com/stebennett/squad-dashboard/pkg/dashboard/models"

type Printer interface {
	Print(defectCounts []models.EscapedDefectCount) error
}
