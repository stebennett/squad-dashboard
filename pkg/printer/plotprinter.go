package printer

import (
	"image/color"

	"github.com/stebennett/squad-dashboard/pkg/dashboard/models"
	"github.com/stebennett/squad-dashboard/pkg/mathutil"

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

func (pp *PlotPrinter) Print(reports Reports) error {
	err := pp.printCycleTimes(reports.CycleTimeReports)
	if err != nil {
		return err
	}

	err = pp.printDefectCounts(reports.EscapedDefects)
	if err != nil {
		return err
	}

	err = pp.printThroughput(reports.ThroughputReports)
	if err != nil {
		return err
	}

	err = pp.printUnplannedWork(reports.UnplannedWorkReports)
	if err != nil {
		return err
	}

	return nil
}

func (pp *PlotPrinter) printDefectCounts(defectCounts models.EscapedDefectReport) error {
	p, err := pp.printChart(defectCounts.WeeklyReports, "Escaped Defects", color.NRGBA{R: 190, G: 0, B: 0, A: 100}, color.NRGBA{R: 190, G: 0, B: 0, A: 255})
	if err != nil {
		return err
	}

	return p.Save(20*vg.Centimeter, 10*vg.Centimeter, pp.GetEscapedDefectsChartLocation())
}

func (pp *PlotPrinter) printCycleTimes(cycleTimeReport models.CycleTimeReport) error {
	p, err := pp.printChart(cycleTimeReport.WeeklyReports, "Cycle Time", color.NRGBA{R: 0, G: 0, B: 190, A: 100}, color.NRGBA{R: 0, G: 0, B: 190, A: 255})
	if err != nil {
		return err
	}

	p, err = pp.addPoints(p, cycleTimeReport.AllCycleTimeItems, color.NRGBA{R: 0, G: 0, B: 190, A: 110})
	if err != nil {
		return err
	}

	return p.Save(20*vg.Centimeter, 10*vg.Centimeter, pp.GetCycleTimeChartLocation())
}

func (pp *PlotPrinter) printThroughput(throughputReports models.ThroughputReport) error {
	p, err := pp.printChart(throughputReports.WeeklyReports, "Throughput", color.NRGBA{R: 0, G: 190, B: 0, A: 100}, color.NRGBA{R: 0, G: 190, B: 0, A: 255})
	if err != nil {
		return err
	}

	return p.Save(20*vg.Centimeter, 10*vg.Centimeter, pp.GetThroughputChartLocation())
}

func (pp *PlotPrinter) printUnplannedWork(unplannedWorkReports models.UnplannedWorkReport) error {
	p, err := pp.printChart(unplannedWorkReports.WeeklyReports, "Unplanned Work", color.NRGBA{R: 100, G: 0, B: 100, A: 100}, color.NRGBA{R: 100, G: 0, B: 100, A: 255})
	if err != nil {
		return err
	}

	return p.Save(20*vg.Centimeter, 10*vg.Centimeter, pp.GetUnplannedWorkChartLocation())
}

func (pp *PlotPrinter) printChart(weekCounts []models.WeekCount, title string, plotColor color.Color, trendlineColor color.Color) (p *plot.Plot, err error) {
	p = plot.New()

	xticks := plot.TimeTicks{Format: "2006-01-02"}

	p.Title.Text = title + " - " + pp.JiraProject
	p.X.Label.Text = "Week Ending"
	p.X.Tick.Marker = xticks
	p.Y.Tick.Length = 0
	p.Y.Tick.Label.Font.Size = 0

	data := make(plotter.XYs, len(weekCounts))
	xys := make([]mathutil.XY, len(weekCounts))
	for i, d := range weekCounts {
		data[i].X = float64(d.WeekEnding.Unix())
		data[i].Y = float64(d.Count)

		xys[i].X = float64(d.WeekEnding.Unix())
		xys[i].Y = float64(d.Count)
	}

	p.Add(plotter.NewGrid())

	line, points, err := plotter.NewLinePoints(data)
	if err != nil {
		return p, err
	}

	line.Color = plotColor
	line.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	points.Shape = draw.BoxGlyph{}
	points.Color = plotColor

	p.Add(line, points)

	linearRegression, _, _ := mathutil.LinearRegression(xys)

	linearRegressionLine, linearRegressionPoints, err := plotter.NewLinePoints(asPlotPoints(linearRegression))
	if err != nil {
		return p, err
	}

	linearRegressionLine.Color = trendlineColor
	linearRegressionPoints.Color = trendlineColor

	p.Add(linearRegressionLine, linearRegressionPoints)

	return p, nil
}

func (pp *PlotPrinter) addPoints(p *plot.Plot, cycleTimes []models.CycleTimeItem, plotColor color.Color) (*plot.Plot, error) {
	data := make(plotter.XYs, len(cycleTimes))
	for i, d := range cycleTimes {
		data[i].X = float64(d.CompletedAt.Unix())
		data[i].Y = float64(d.Size)
	}

	scatter, err := plotter.NewScatter(data)
	if err != nil {
		return p, err
	}

	scatter.Color = plotColor
	scatter.Shape = draw.CircleGlyph{}

	p.Add(scatter)

	return p, nil
}

func (pp *PlotPrinter) GetCycleTimeChartLocation() string {
	return pp.OutputDirectory + "/cycle-time-" + pp.JiraProject + ".png"
}

func (pp *PlotPrinter) GetThroughputChartLocation() string {
	return pp.OutputDirectory + "/throughput-" + pp.JiraProject + ".png"
}

func (pp *PlotPrinter) GetEscapedDefectsChartLocation() string {
	return pp.OutputDirectory + "/escaped-defects-" + pp.JiraProject + ".png"
}

func (pp *PlotPrinter) GetUnplannedWorkChartLocation() string {
	return pp.OutputDirectory + "/unplanned-" + pp.JiraProject + ".png"
}

func asPlotPoints(points []mathutil.XY) plotter.XYs {
	plotterPoints := make(plotter.XYs, len(points))
	for i, v := range points {
		plotterPoints[i].X = v.X
		plotterPoints[i].Y = v.Y
	}
	return plotterPoints
}
