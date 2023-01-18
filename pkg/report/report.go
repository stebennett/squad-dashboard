package report

import (
	"time"

	"github.com/jung-kurt/gofpdf"
)

type Direction int

const (
	Up Direction = iota
	Down
	Flat
)

type CellColor struct {
	R, G, B int
}

type Columns []string
type TableData []map[string]string

type Table struct {
	Cols Columns
	Data TableData
}

type ReportDashboardItem struct {
	BackgroundColor CellColor
	Chart           string
	InfoTable       Table
}

type SpeedAnomaly struct {
	IssueKey      string
	Size          int
	CompletedDate time.Time
	Link          string
}

type ReportDashboard struct {
	Quality       ReportDashboardItem
	Quantity      ReportDashboardItem
	Speed         ReportDashboardItem
	UnplannedWork ReportDashboardItem

	SpeedAnomalies []SpeedAnomaly
}

type ReportData struct {
	Dashboards map[string]ReportDashboard
}

const (
	font                 = "Arial"
	margin               = 10.0
	headingFontHeight    = 16.0
	subheadingFontHeight = 14.0
	fontHeight           = 10.0
	cellMargin           = 2.0
	pageWidth            = 210
)

func GeneratePdfReport(reportData ReportData, outputFile string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(margin, margin, margin)
	pdf.AddPage()
	pdf.SetFont(font, "B", headingFontHeight)
	pdf.SetCellMargin(cellMargin)
	pdf.Cell(0, pdf.PointToUnitConvert(headingFontHeight)+2*cellMargin, "Squad Dashboard Report")

	pdf.SetY(pdf.PointToUnitConvert(headingFontHeight) + 2*cellMargin + 2*margin)
	pdf.SetFont(font, "", fontHeight)

	addDashboardView(pdf, reportData)

	pdf.SetFont(font, "B", subheadingFontHeight)
	for k, v := range reportData.Dashboards {
		addQualityPage(pdf, k, v.Quality)
		addSpeedPage(pdf, k, v.Speed)
		addQuantityPage(pdf, k, v.Quantity)
		addConsistencyPage(pdf, k, v.UnplannedWork)
	}

	return pdf.OutputFileAndClose(outputFile)
}

func addDashboardView(pdf gofpdf.Pdf, reportData ReportData) {
	columns := [...]string{"Squad", "Quality", "Speed", "Quantity", "Consistency", "Resilience"}

	cellWidth := float64((pageWidth - (2 * margin) - (len(columns) * 2 * cellMargin)) / len(columns))
	cellHeight := float64(pdf.PointConvert(fontHeight) + 2*cellMargin)

	pdf.SetFont(font, "B", fontHeight)

	// headings
	for _, c := range columns {
		pdf.CellFormat(cellWidth, cellHeight, c, gofpdf.BorderFull, 0, gofpdf.AlignLeft, false, 0, "")
	}

	pdf.SetY(pdf.GetY() + cellHeight)
	pdf.SetFont(font, "B", fontHeight)

	for k, v := range reportData.Dashboards {
		pdf.CellFormat(cellWidth, cellHeight, k, gofpdf.BorderFull, 0, gofpdf.AlignLeft, false, 0, "")

		pdf.SetFillColor(v.Quality.BackgroundColor.R, v.Quality.BackgroundColor.G, v.Quality.BackgroundColor.B)
		pdf.CellFormat(cellWidth, cellHeight, "", gofpdf.BorderFull, 0, gofpdf.AlignLeft, true, 0, "")

		pdf.SetFillColor(v.Speed.BackgroundColor.R, v.Speed.BackgroundColor.G, v.Speed.BackgroundColor.B)
		pdf.CellFormat(cellWidth, cellHeight, "", gofpdf.BorderFull, 0, gofpdf.AlignLeft, true, 0, "")

		pdf.SetFillColor(v.Quantity.BackgroundColor.R, v.Quantity.BackgroundColor.G, v.Quantity.BackgroundColor.B)
		pdf.CellFormat(cellWidth, cellHeight, "", gofpdf.BorderFull, 0, gofpdf.AlignLeft, true, 0, "")

		pdf.SetFillColor(v.UnplannedWork.BackgroundColor.R, v.UnplannedWork.BackgroundColor.G, v.UnplannedWork.BackgroundColor.B)
		pdf.CellFormat(cellWidth, cellHeight, "", gofpdf.BorderFull, 0, gofpdf.AlignLeft, true, 0, "")

		pdf.SetFillColor(163, 163, 163)
		pdf.CellFormat(cellWidth, cellHeight, "", gofpdf.BorderFull, 0, gofpdf.AlignLeft, true, 0, "")

		pdf.SetY(pdf.GetY() + cellHeight)
	}
}

func addQualityPage(pdf gofpdf.Pdf, squad string, quality ReportDashboardItem) {
	addChartPage(pdf, squad+" - Quality", quality.Chart)
	addTable(pdf, quality.InfoTable)
}

func addSpeedPage(pdf gofpdf.Pdf, squad string, speed ReportDashboardItem) {
	addChartPage(pdf, squad+" - Speed", speed.Chart)
	addTable(pdf, speed.InfoTable)
}

func addQuantityPage(pdf gofpdf.Pdf, squad string, quantity ReportDashboardItem) {
	addChartPage(pdf, squad+" - Quantity", quantity.Chart)
	addTable(pdf, quantity.InfoTable)
}

func addConsistencyPage(pdf gofpdf.Pdf, squad string, consistency ReportDashboardItem) {
	addChartPage(pdf, squad+" - Consistency", consistency.Chart)
	addTable(pdf, consistency.InfoTable)
}

func addChartPage(pdf gofpdf.Pdf, chartTitle string, chartLocation string) {
	opt := gofpdf.ImageOptions{
		ImageType: "png",
	}

	pdf.AddPage()
	pdf.Cell(0, pdf.PointToUnitConvert(subheadingFontHeight), chartTitle)
	pdf.SetY(pdf.PointToUnitConvert(subheadingFontHeight) + 2*margin)
	pdf.ImageOptions(chartLocation, margin, margin, pageWidth-(2*margin), 0, true, opt, 0, "")
}

func addTable(pdf gofpdf.Pdf, table Table) {
	colCount := len(table.Cols)
	if colCount == 0 {
		return
	}

	pdf.SetY(pdf.GetY() + margin)

	cellWidth := float64((pageWidth - (2 * margin) - (colCount * 2 * cellMargin)) / colCount)
	cellHeight := pdf.PointToUnitConvert(fontHeight) + 2*cellMargin

	pdf.SetFont(font, "B", fontHeight)

	// headings
	for _, heading := range table.Cols {
		pdf.CellFormat(cellWidth, cellHeight, heading, gofpdf.BorderFull, 0, gofpdf.AlignLeft, false, 0, "")
	}

	pdf.SetFont(font, "", fontHeight)
	pdf.SetY(pdf.GetY() + cellHeight)

	for _, row := range table.Data {
		for _, col := range table.Cols {
			pdf.CellFormat(cellWidth, cellHeight, row[col], gofpdf.BorderFull, 0, gofpdf.AlignLeft, false, 0, "")
		}
		pdf.SetY(pdf.GetY() + cellHeight)
	}
}
