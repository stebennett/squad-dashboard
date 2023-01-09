package printer

import (
	"log"

	"github.com/stebennett/squad-dashboard/pkg/dashboard/models"
	"github.com/stebennett/squad-dashboard/pkg/jiracalculationsrepository"
	"github.com/stebennett/squad-dashboard/pkg/mathutil"
)

type CommandLinePrinter struct{}

func NewCommandLinePrinter() *CommandLinePrinter {
	return &CommandLinePrinter{}
}

func (c *CommandLinePrinter) Print(reports Reports) error {
	err := c.PrintDefectCounts(reports.EscapedDefects)
	if err != nil {
		return err
	}
	err = c.PrintCycleTimes(reports.CycleTimeReports, reports.AllCycleTimes)
	if err != nil {
		return err
	}
	err = c.PrintThroughput(reports.ThroughputReports)
	if err != nil {
		return err
	}
	return err
}

func (c *CommandLinePrinter) PrintDefectCounts(defectCounts []models.WeekCount) error {
	log.Printf("---- defects ----")
	for idx, defectCount := range defectCounts {
		log.Printf("%d> weekEnding: %s; defectsCreated: %d", idx, defectCount.WeekEnding, defectCount.Count)
	}

	return nil
}

func (c *CommandLinePrinter) PrintCycleTimes(cycleTimeReports []models.WeekCount, allCycleTime []jiracalculationsrepository.CycleTimes) error {
	log.Printf("---- cycle times ----")
	for idx, ct := range cycleTimeReports {
		log.Printf("%d> weekEnding: %s; avgCycleTime: %d", idx, ct.WeekEnding, ct.Count)
	}

	log.Printf("-- anomaly issues --")

	ctValues := make([]int, len(allCycleTime))
	for i, c := range allCycleTime {
		ctValues[i] = c.Size
	}

	percentile75th := mathutil.Percentile(0.75, ctValues)
	log.Printf("> Percentile: %d", percentile75th)

	for _, c := range allCycleTime {
		if c.Size >= percentile75th {
			log.Printf("Issue outside percentile: %s; Completed: %s; CycleTime: %d", c.IssueKey, c.Completed, c.Size)
		}
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
