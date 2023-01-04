package printer

import (
	"log"

	"github.com/stebennett/squad-dashboard/pkg/dashboard/models"
)

type CommandLinePrinter struct{}

func NewCommandLinePrinter() *CommandLinePrinter {
	return &CommandLinePrinter{}
}

func (c *CommandLinePrinter) PrintDefectCounts(defectCounts []models.EscapedDefectCount) error {
	log.Printf("---- defects ----")
	for idx, defectCount := range defectCounts {
		log.Printf("%d> weekEnding: %s; defectsCreated: %d", idx, defectCount.WeekEnding, defectCount.NumberOfDefectsCreated)
	}

	return nil
}

func (c *CommandLinePrinter) PrintCycleTimes(cycleTimeReports []models.CycleTimeReport) error {
	log.Printf("---- cycle times ----")
	for idx, ct := range cycleTimeReports {
		log.Printf("%d> weekEnding: %s; avgCycleTime: %d", idx, ct.WeekEnding, ct.AverageCycleTime)
	}

	return nil
}

func (c *CommandLinePrinter) PrintThroughput(throughputReports []models.ThroughputReport) error {
	log.Printf("---- throughput ----")
	for idx, tp := range throughputReports {
		log.Printf("%d> weekEnding: %s; throughput: %d", idx, tp.WeekEnding, tp.Throughput)
	}

	return nil
}
