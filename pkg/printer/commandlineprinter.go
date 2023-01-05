package printer

import (
	"log"

	"github.com/stebennett/squad-dashboard/pkg/dashboard/models"
)

type CommandLinePrinter struct{}

func NewCommandLinePrinter() *CommandLinePrinter {
	return &CommandLinePrinter{}
}

func (c *CommandLinePrinter) PrintDefectCounts(defectCounts []models.WeekCount) error {
	log.Printf("---- defects ----")
	for idx, defectCount := range defectCounts {
		log.Printf("%d> weekEnding: %s; defectsCreated: %d", idx, defectCount.WeekEnding, defectCount.Count)
	}

	return nil
}

func (c *CommandLinePrinter) PrintCycleTimes(cycleTimeReports []models.WeekCount) error {
	log.Printf("---- cycle times ----")
	for idx, ct := range cycleTimeReports {
		log.Printf("%d> weekEnding: %s; avgCycleTime: %d", idx, ct.WeekEnding, ct.Count)
	}

	return nil
}

func (c *CommandLinePrinter) PrintThroughput(throughputReports []models.WeekCount) error {
	log.Printf("---- throughput ----")
	for idx, tp := range throughputReports {
		log.Printf("%d> weekEnding: %s; throughput: %d", idx, tp.WeekEnding, tp.Count)
	}

	return nil
}
