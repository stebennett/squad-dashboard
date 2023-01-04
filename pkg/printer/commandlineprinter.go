package printer

import (
	"log"

	"github.com/stebennett/squad-dashboard/pkg/dashboard/models"
)

type CommandLinePrinter struct{}

func NewCommandLinePrinter() *CommandLinePrinter {
	return &CommandLinePrinter{}
}

func (c *CommandLinePrinter) Print(defectCounts []models.EscapedDefectCount) error {
	for idx, defectCount := range defectCounts {
		log.Printf("%d> weekEnding: %s; defectsCreated: %d", idx, defectCount.WeekEnding, defectCount.NumberOfDefectsCreated)
	}

	return nil
}
