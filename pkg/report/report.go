package report

import (
	"fmt"
	"image/color"
	"sort"
	"time"

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
	Chart           string
}

type SpeedAnomaly struct {
	IssueKey      string
	Size          int
	CompletedDate time.Time
	Link          string
}

type ReportDashboard struct {
	Quality  ReportDashboardItem
	Quantity ReportDashboardItem
	Speed    ReportDashboardItem

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

	pdf.SetY(pdf.PointToUnitConvert(headingFontHeight) + margin)
	pdf.SetFont(font, "", fontHeight)

	pdf.SetFont(font, "B", subheadingFontHeight)
	for k, v := range reportData.Dashboards {
		addChartPage(pdf, k+" - Quality", v.Quality.Chart)
		addChartPage(pdf, k+" - Quantity", v.Quantity.Chart)
		addChartPage(pdf, k+" - Speed", v.Speed.Chart)

		addSpeedAnomaliesTable(pdf, v.SpeedAnomalies)
	}

	return pdf.OutputFileAndClose(outputFile)
}

func addChartPage(pdf gofpdf.Pdf, chartTitle string, chartLocation string) *gofpdf.Pdf {
	opt := gofpdf.ImageOptions{
		ImageType: "png",
	}

	pdf.AddPage()
	pdf.Cell(0, pdf.PointToUnitConvert(subheadingFontHeight), chartTitle)
	pdf.SetY(pdf.PointToUnitConvert(subheadingFontHeight) + 2*margin)
	pdf.ImageOptions(chartLocation, margin, margin, pageWidth-(2*margin), 0, true, opt, 0, "")
	return &pdf
}

func addSpeedAnomaliesTable(pdf gofpdf.Pdf, anomalies []SpeedAnomaly) *gofpdf.Pdf {
	const (
		colCount = 4
	)
	pdf.Cell(0, pdf.PointToUnitConvert(subheadingFontHeight), "Anomalies")

	pdf.SetY(pdf.GetY() + margin)

	cellWidth := (pageWidth - (2 * margin) - (colCount * 2 * cellMargin)) / colCount
	cellHeight := pdf.PointToUnitConvert(fontHeight) + 2*cellMargin

	pdf.SetFont(font, "B", fontHeight)

	// headings
	pdf.CellFormat(cellWidth, cellHeight, "Date Completed", gofpdf.BorderFull, 0, gofpdf.AlignLeft, false, 0, "")
	pdf.CellFormat(cellWidth, cellHeight, "Issue Key", gofpdf.BorderFull, 0, gofpdf.AlignLeft, false, 0, "")
	pdf.CellFormat(cellWidth, cellHeight, "Cycle Time", gofpdf.BorderFull, 0, gofpdf.AlignLeft, false, 0, "")
	pdf.CellFormat(cellWidth, cellHeight, "Link", gofpdf.BorderFull, 0, gofpdf.AlignLeft, false, 0, "")

	pdf.SetFont(font, "B", fontHeight)
	pdf.SetY(pdf.GetY() + cellHeight)

	sort.Slice(anomalies, func(i, j int) bool {
		return anomalies[i].CompletedDate.After(anomalies[j].CompletedDate)
	})

	for _, anomaly := range anomalies {
		pdf.CellFormat(cellWidth, cellHeight, anomaly.CompletedDate.Format("2006-01-02"), gofpdf.BorderFull, 0, gofpdf.AlignLeft, false, 0, "")
		pdf.CellFormat(cellWidth, cellHeight, anomaly.IssueKey, gofpdf.BorderFull, 0, gofpdf.AlignLeft, false, 0, "")
		pdf.CellFormat(cellWidth, cellHeight, fmt.Sprintf("%d", anomaly.Size), gofpdf.BorderFull, 0, gofpdf.AlignLeft, false, 0, "")
		pdf.CellFormat(cellWidth, cellHeight, anomaly.Link, gofpdf.BorderFull, 0, gofpdf.AlignLeft, false, 0, anomaly.Link)

		pdf.SetY(pdf.GetY() + cellHeight)
	}

	return &pdf
}
