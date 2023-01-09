package report

import (
	"image/color"

	"github.com/jung-kurt/gofpdf"
)

type Direction int

const (
	Up Direction = iota
	Down
	Flat
)

type ReportDashboardItem struct {
	BackgroundColor color.Color
}

type ReportDashboard struct {
	Quality  ReportDashboardItem
	Quantity ReportDashboardItem
	Speed    ReportDashboardItem

	SpeedAnomalies []string
}

type ReportData struct {
	Dashboards map[string]ReportDashboard
}

func GeneratePdfReport(reportData ReportData, outputFile string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "Squad Dashboard Report")

	return pdf.OutputFileAndClose(outputFile)
}
