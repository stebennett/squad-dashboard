package printer

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/stebennett/squad-dashboard/pkg/dashboard/models"
	"github.com/stebennett/squad-dashboard/pkg/jiracalculationsrepository"
	"github.com/stebennett/squad-dashboard/pkg/report"
)

type PdfReportPrinter struct {
	CycleTimeChart      string
	ThroughputChart     string
	EscapedDefectsChart string
	JiraProject         string
}

func NewPdfReportPrinter(cycleTimeChart string, throughputChart string, escapedDefectsChart string, project string) *PdfReportPrinter {
	return &PdfReportPrinter{
		CycleTimeChart:      cycleTimeChart,
		ThroughputChart:     throughputChart,
		EscapedDefectsChart: escapedDefectsChart,
		JiraProject:         project,
	}
}

func (p *PdfReportPrinter) Print(reports Reports) error {
	reportDashboards := make(map[string]report.ReportDashboard)

	reportDashboards[p.JiraProject] = report.ReportDashboard{
		Quality: report.ReportDashboardItem{
			BackgroundColor: color.Black,
			Chart:           p.EscapedDefectsChart,
		},
		Quantity: report.ReportDashboardItem{
			BackgroundColor: color.Black,
			Chart:           p.ThroughputChart,
		},
		Speed: report.ReportDashboardItem{
			BackgroundColor: color.Black,
			Chart:           p.CycleTimeChart,
		},
		SpeedAnomalies: getSpeedAnomalies(reports.CycleTimeReports, reports.AllCycleTimes),
	}

	reportData := report.ReportData{
		Dashboards: reportDashboards,
	}

	err := report.GeneratePdfReport(reportData, "/output/report.pdf")
	return err
}

func getSpeedAnomalies(cycleTimesReports []models.WeekCount, allCycleTimes []jiracalculationsrepository.CycleTimes) []report.SpeedAnomaly {
	speedAnomalies := []report.SpeedAnomaly{}

	for _, ctavg := range cycleTimesReports {
		cycleTimesInWeek := filterCycleTimes(allCycleTimes, ctavg.WeekEnding.AddDate(0, 0, -7), ctavg.WeekEnding)
		for _, ct := range cycleTimesInWeek {
			if ct.Size > ctavg.Count {
				speedAnomalies = append(speedAnomalies, report.SpeedAnomaly{
					IssueKey:      ct.IssueKey,
					Size:          ct.Size,
					CompletedDate: ct.Completed,
					Link:          fmt.Sprintf("https://jira/%s", ct.IssueKey),
				})
			}
		}
	}

	return speedAnomalies
}

func filterCycleTimes(cycleTimes []jiracalculationsrepository.CycleTimes, startDate time.Time, endDate time.Time) []jiracalculationsrepository.CycleTimes {
	filteredCycleTimes := []jiracalculationsrepository.CycleTimes{}

	compareBeforeYear, compareBeforeMonth, compareBeforeDate := endDate.Date()
	compareAfterYear, compareAfterMonth, compareAfterDate := startDate.Date()

	log.Printf("filtering issues between [%d, %d, %d] and [%d, %d, %d]", compareAfterYear, compareAfterMonth, compareAfterDate, compareBeforeYear, compareBeforeMonth, compareBeforeDate)

	for _, ct := range cycleTimes {
		completedYear, completedMonth, completedDate := ct.Completed.Date()

		if completedYear > compareAfterYear && completedYear <= compareBeforeYear &&
			completedMonth > compareAfterMonth && completedMonth <= compareBeforeMonth &&
			completedDate > compareAfterDate && completedDate <= compareBeforeDate {
			log.Println("Adding item")
			filteredCycleTimes = append(filteredCycleTimes, ct)
		}
	}

	return filteredCycleTimes
}
