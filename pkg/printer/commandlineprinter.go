package printer

import (
	"log"

	"github.com/stebennett/squad-dashboard/pkg/mathutil"
	"github.com/stebennett/squad-dashboard/pkg/models"
)

type CommandLinePrinter struct{}

func NewCommandLinePrinter() *CommandLinePrinter {
	return &CommandLinePrinter{}
}

func (c *CommandLinePrinter) Print(reports Reports) error {
	err := c.printDefectCounts(reports.EscapedDefects)
	if err != nil {
		return err
	}
	err = c.printCycleTimes(reports.CycleTimeReports)
	if err != nil {
		return err
	}
	err = c.printThroughput(reports.ThroughputReports)
	if err != nil {
		return err
	}
	return err
}

func (c *CommandLinePrinter) printDefectCounts(defectCounts models.EscapedDefectReport) error {
	log.Printf("---- defects ----")
	for idx, defectCount := range defectCounts.WeeklyReports {
		log.Printf("%d> weekEnding: %s; defectsCreated: %d", idx, defectCount.WeekEnding, defectCount.Count)
	}

	return nil
}

func (c *CommandLinePrinter) printCycleTimes(cycleTimeReport models.CycleTimeReport) error {
	log.Printf("---- cycle times ----")
	for idx, ct := range cycleTimeReport.WeeklyReports {
		log.Printf("%d> weekEnding: %s; avgCycleTime: %d", idx, ct.WeekEnding, ct.Count)
	}

	log.Printf("-- anomaly issues --")

	ctValues := make([]int, len(cycleTimeReport.AllCycleTimeItems))
	for i, c := range cycleTimeReport.AllCycleTimeItems {
		ctValues[i] = c.WorkingCycleTime
	}

	percentile75th := mathutil.Percentile(0.75, ctValues)
	log.Printf("> Percentile: %d", percentile75th)

	for _, c := range cycleTimeReport.AllCycleTimeItems {
		if c.WorkingCycleTime >= percentile75th {
			log.Printf("Issue outside percentile: %s; Completed: %s; CycleTime: %d", c.IssueKey, c.CompletedAt, c.WorkingCycleTime)
		}
	}

	return nil
}

func (c *CommandLinePrinter) printThroughput(throughputReports models.ThroughputReport) error {
	log.Printf("---- throughput ----")
	for idx, tp := range throughputReports.WeeklyReports {
		log.Printf("%d> weekEnding: %s; throughput: %d", idx, tp.WeekEnding, tp.Count)
	}

	return nil
}

func (c *CommandLinePrinter) printUnplannedWorkPercent(unplannedWorkReports models.UnplannedWorkReport) error {
	log.Printf("---- throughput ----")
	for idx, tp := range unplannedWorkReports.WeeklyReports {
		log.Printf("%d> weekEnding: %s; unplannedWork: %d%%", idx, tp.WeekEnding, tp.Count)
	}

	return nil
}
