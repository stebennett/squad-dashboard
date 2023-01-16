package printer

import (
	"fmt"
	"time"

	"github.com/stebennett/squad-dashboard/pkg/dashboard/models"
	"github.com/stebennett/squad-dashboard/pkg/jiracalculationsrepository"
	"github.com/stebennett/squad-dashboard/pkg/mathutil"
	"github.com/stebennett/squad-dashboard/pkg/report"
)

type PdfReportPrinter struct {
	CycleTimeChart      string
	ThroughputChart     string
	EscapedDefectsChart string
	UnplannedWorkChart  string
	JiraProject         string
}

var (
	green = report.CellColor{R: 47, G: 247, B: 7}
	amber = report.CellColor{R: 247, G: 147, B: 7}
	red   = report.CellColor{R: 247, G: 7, B: 7}
)

func NewPdfReportPrinter(cycleTimeChart string, throughputChart string, escapedDefectsChart string, unplannedWorkChart string, project string) *PdfReportPrinter {
	return &PdfReportPrinter{
		CycleTimeChart:      cycleTimeChart,
		ThroughputChart:     throughputChart,
		EscapedDefectsChart: escapedDefectsChart,
		UnplannedWorkChart:  unplannedWorkChart,
		JiraProject:         project,
	}
}

func (p *PdfReportPrinter) Print(reports Reports) error {
	reportDashboards := make(map[string]report.ReportDashboard)

	reportDashboards[p.JiraProject] = report.ReportDashboard{
		Quality: report.ReportDashboardItem{
			BackgroundColor: pickColorLowerBetter(reports.EscapedDefects),
			Chart:           p.EscapedDefectsChart,
		},
		Quantity: report.ReportDashboardItem{
			BackgroundColor: pickColorHigherBetter(reports.ThroughputReports),
			Chart:           p.ThroughputChart,
		},
		Speed: report.ReportDashboardItem{
			BackgroundColor: pickColorLowerBetter(reports.CycleTimeReports),
			Chart:           p.CycleTimeChart,
		},
		UnplannedWork: report.ReportDashboardItem{
			BackgroundColor: pickColorLowerBetter(reports.UnplannedWorkReports),
			Chart:           p.UnplannedWorkChart,
		},
		SpeedAnomalies: getSpeedAnomalies(reports.CycleTimeReports, reports.AllCycleTimes),
	}

	reportData := report.ReportData{
		Dashboards: reportDashboards,
	}

	err := report.GeneratePdfReport(reportData, "/output/report-"+p.JiraProject+".pdf")
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

	for _, ct := range cycleTimes {
		completedYear, completedMonth, completedDate := ct.Completed.Date()

		if completedYear > compareAfterYear && completedYear <= compareBeforeYear &&
			completedMonth > compareAfterMonth && completedMonth <= compareBeforeMonth &&
			completedDate > compareAfterDate && completedDate <= compareBeforeDate {
			filteredCycleTimes = append(filteredCycleTimes, ct)
		}
	}

	return filteredCycleTimes
}

func pickColorLowerBetter(reports []models.WeekCount) (trendColor report.CellColor) {
	xys := make([]mathutil.XY, len(reports))
	for i, v := range reports {
		xys[i].X = float64(v.WeekEnding.Unix())
		xys[i].Y = float64(v.Count)
	}

	_, gradiant, _ := mathutil.LinearRegression(xys)
	trend := gradiant * (7 * 24 * 60 * 60) // weekly ticks

	switch {
	case trend < 0.15: // downward trend
		trendColor = green
	case trend >= 0.15 && trend < 1.0: // slight upward (1 new item per week)
		trendColor = amber
	case trend >= 1.0:
		trendColor = red // more than 1 item per week
	}
	return trendColor
}

func pickColorHigherBetter(reports []models.WeekCount) (trendColor report.CellColor) {
	xys := make([]mathutil.XY, len(reports))
	for i, v := range reports {
		xys[i].X = float64(v.WeekEnding.Unix())
		xys[i].Y = float64(v.Count)
	}

	_, gradiant, _ := mathutil.LinearRegression(xys)
	trend := gradiant * (7 * 24 * 60 * 60)

	switch {
	case trend < -1.0:
		trendColor = red
	case trend >= -1.0 && trend < 0.0:
		trendColor = amber
	case trend >= 0.0:
		trendColor = green
	}

	return trendColor
}
