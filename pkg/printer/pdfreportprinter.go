package printer

import (
	"fmt"

	"github.com/stebennett/squad-dashboard/pkg/dashboard/models"
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

const (
	DateCreated   = "Date Created"
	DateCompleted = "Date Completed"
	Link          = "Link"
	Size          = "Size"
)

var (
	EscapedDefectsColumns = []string{DateCreated, Link}
	QuantityColumns       = []string{DateCompleted, Link}
	SpeedColumns          = []string{DateCompleted, Size, Link}
	UnplannedWorkColumns  = []string{DateCompleted, Link}
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
			BackgroundColor: pickColorLowerBetter(reports.EscapedDefects.WeeklyReports),
			Chart:           p.EscapedDefectsChart,
			InfoTable:       createEscapedDefectsTable(reports.EscapedDefects),
		},
		Quantity: report.ReportDashboardItem{
			BackgroundColor: pickColorHigherBetter(reports.ThroughputReports.WeeklyReports),
			Chart:           p.ThroughputChart,
			InfoTable:       createQuantityTable(reports.ThroughputReports),
		},
		Speed: report.ReportDashboardItem{
			BackgroundColor: pickColorLowerBetter(reports.CycleTimeReports.WeeklyReports),
			Chart:           p.CycleTimeChart,
			InfoTable:       createSpeedTable(reports.CycleTimeReports),
		},
		UnplannedWork: report.ReportDashboardItem{
			BackgroundColor: pickColorLowerBetter(reports.UnplannedWorkReports.WeeklyReports),
			Chart:           p.UnplannedWorkChart,
			InfoTable:       createUnplannedTable(reports.UnplannedWorkReports),
		},
	}

	reportData := report.ReportData{
		Dashboards: reportDashboards,
	}

	err := report.GeneratePdfReport(reportData, "/output/report-"+p.JiraProject+".pdf")
	return err
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

func createEscapedDefectsTable(edr models.EscapedDefectReport) report.Table {
	tableData := make(report.TableData, len(edr.LastWeekEscapedDefectItems))

	for i, v := range edr.LastWeekEscapedDefectItems {
		row := make(map[string]string, len(EscapedDefectsColumns))
		row[DateCreated] = v.CreatedAt.Format("2006-01-02")
		row[Link] = v.IssueKey

		tableData[i] = row
	}

	return report.Table{
		Cols: EscapedDefectsColumns,
		Data: tableData,
	}
}

func createQuantityTable(tr models.ThroughputReport) report.Table {
	tableData := make(report.TableData, len(tr.LastWeekThroughputItems))

	for i, v := range tr.LastWeekThroughputItems {
		row := make(map[string]string, len(QuantityColumns))
		row[DateCompleted] = v.CompletedAt.Format("2006-01-02")
		row[Link] = v.IssueKey

		tableData[i] = row
	}

	return report.Table{
		Cols: QuantityColumns,
		Data: tableData,
	}
}

func createSpeedTable(ctr models.CycleTimeReport) report.Table {
	tableData := make(report.TableData, len(ctr.LastWeekCycleTimeItems))

	for i, v := range ctr.LastWeekCycleTimeItems {
		row := make(map[string]string, len(SpeedColumns))
		row[DateCompleted] = v.CompletedAt.Format("2006-01-02")
		row[Link] = v.IssueKey
		row[Size] = fmt.Sprintf("%d", v.Size)

		tableData[i] = row
	}

	return report.Table{
		Cols: SpeedColumns,
		Data: tableData,
	}
}

func createUnplannedTable(uwr models.UnplannedWorkReport) report.Table {
	tableData := make(report.TableData, len(uwr.LastWeekUnplannedWorkItems))

	for i, v := range uwr.LastWeekUnplannedWorkItems {
		row := make(map[string]string, len(UnplannedWorkColumns))
		row[DateCompleted] = v.CompletedAt.Format("2006-01-02")
		row[Link] = v.IssueKey

		tableData[i] = row
	}

	return report.Table{
		Cols: UnplannedWorkColumns,
		Data: tableData,
	}
}
