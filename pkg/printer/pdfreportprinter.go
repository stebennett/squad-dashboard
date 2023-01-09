package printer

import (
	"image/color"

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
		},
		Quantity: report.ReportDashboardItem{
			BackgroundColor: color.Black,
		},
		Speed: report.ReportDashboardItem{
			BackgroundColor: color.Black,
		},
	}

	reportData := report.ReportData{
		Dashboards: reportDashboards,
	}

	err := report.GeneratePdfReport(reportData, "/output/report.pdf")
	return err
}
