package printer

import (
	"image/color"
	"sort"

	"github.com/stebennett/squad-dashboard/pkg/dashboard/models"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

type PlotPrinter struct {
	OutputDirectory string
	JiraProject     string
}

func NewPlotPrinter(outputDirectory string, project string) *PlotPrinter {
	return &PlotPrinter{OutputDirectory: outputDirectory, JiraProject: project}
}

func (pp *PlotPrinter) PrintDefectCounts(defectCounts []models.WeekCount) error {
	err := pp.printChart(defectCounts, "Escaped Defects", "escaped-defects", color.NRGBA{R: 190, G: 0, B: 0, A: 100}, color.NRGBA{R: 190, G: 0, B: 0, A: 255})
	if err != nil {
		return err
	}

	return pp.printChart(pp.normalizeData(defectCounts), "Escaped Defects - Normalized", "escaped-defects-normalized", color.NRGBA{R: 190, G: 0, B: 0, A: 100}, color.NRGBA{R: 190, G: 0, B: 0, A: 255})
}

func (pp *PlotPrinter) PrintCycleTimes(cycleTimeReports []models.WeekCount) error {
	err := pp.printChart(cycleTimeReports, "Cycle Time", "cycle-time", color.NRGBA{R: 0, G: 0, B: 190, A: 100}, color.NRGBA{R: 0, G: 0, B: 190, A: 255})
	if err != nil {
		return err
	}

	return pp.printChart(pp.normalizeData(cycleTimeReports), "Cycle Time - Normalized", "cycle-time-normalized", color.NRGBA{R: 0, G: 0, B: 190, A: 100}, color.NRGBA{R: 0, G: 0, B: 190, A: 255})
}

func (pp *PlotPrinter) PrintThroughput(throughputReports []models.WeekCount) error {
	err := pp.printChart(throughputReports, "Throughput", "throughput", color.NRGBA{R: 0, G: 190, B: 0, A: 100}, color.NRGBA{R: 0, G: 190, B: 0, A: 255})
	if err != nil {
		return err
	}

	return pp.printChart(pp.normalizeData(throughputReports), "Throughput - Normalized", "throughput-normalized", color.NRGBA{R: 0, G: 190, B: 0, A: 100}, color.NRGBA{R: 0, G: 190, B: 0, A: 255})
}

func (pp *PlotPrinter) printChart(weekCounts []models.WeekCount, title string, filename string, plotColor color.Color, trendlineColor color.Color) error {
	p := plot.New()

	xticks := plot.TimeTicks{Format: "2006-01-02"}

	p.Title.Text = title + " - " + pp.JiraProject
	p.X.Label.Text = "Week Ending"
	p.X.Tick.Marker = xticks
	p.Y.Label.Text = "Number"

	data := make(plotter.XYs, len(weekCounts))
	for i, d := range weekCounts {
		data[i].X = float64(d.WeekEnding.Unix())
		data[i].Y = float64(d.Count)
	}

	p.Add(plotter.NewGrid())

	line, points, err := plotter.NewLinePoints(data)
	if err != nil {
		return err
	}

	line.Color = plotColor
	line.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	points.Shape = draw.BoxGlyph{}
	points.Color = plotColor

	p.Add(line, points)

	linearRegression := pp.createLinearRegression(data)

	linearRegressionLine, linearRegressionPoints, err := plotter.NewLinePoints(linearRegression)
	if err != nil {
		return err
	}

	linearRegressionLine.Color = trendlineColor
	linearRegressionPoints.Color = trendlineColor

	p.Add(linearRegressionLine, linearRegressionPoints)

	err = p.Save(20*vg.Centimeter, 10*vg.Centimeter, pp.OutputDirectory+"/"+filename+"-"+pp.JiraProject+".png")
	return err
}

func (pp *PlotPrinter) createLinearRegression(inData plotter.XYs) plotter.XYs {
	q := len(inData)

	if q == 0 {
		return make(plotter.XYs, 0)
	}

	p := float64(q)

	sum_x, sum_y, sum_xx, sum_xy := 0.0, 0.0, 0.0, 0.0

	for _, p := range inData {
		sum_x += p.X
		sum_y += p.Y
		sum_xx += p.X * p.X
		sum_xy += p.X * p.Y
	}

	m := (p*sum_xy - sum_x*sum_y) / (p*sum_xx - sum_x*sum_x)
	b := (sum_y / p) - (m * sum_x / p)

	r := make(plotter.XYs, q)
	for i, p := range inData {
		r[i].X = p.X
		r[i].Y = p.X*m + b
	}

	return r
}

func (pp *PlotPrinter) normalizeData(data []models.WeekCount) []models.WeekCount {
	sort.Slice(data, func(i, j int) bool {
		return data[i].Count < data[j].Count
	})

	size := len(data)
	min, max := data[0].Count, data[size-1].Count

	result := make([]models.WeekCount, size)
	for i, d := range data {
		result[i] = models.WeekCount{
			WeekEnding: d.WeekEnding,
			Count:      int((float64(d.Count-min) / float64(max-min)) * 100.0),
		}
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i].WeekEnding.Before(result[j].WeekEnding)
	})

	return result
}
