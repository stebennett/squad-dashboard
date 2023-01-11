package printer

import (
	"image/color"
	"log"

	"github.com/stebennett/squad-dashboard/pkg/dashboard/models"
	"github.com/stebennett/squad-dashboard/pkg/jiracalculationsrepository"
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
	err := pp.printCycleTimes(reports.CycleTimeReports, reports.AllCycleTimes)
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

func (pp *PlotPrinter) printDefectCounts(defectCounts []models.WeekCount) error {
	trend, p, err := pp.printChart(defectCounts, "Escaped Defects", color.NRGBA{R: 190, G: 0, B: 0, A: 100}, color.NRGBA{R: 190, G: 0, B: 0, A: 255})
	if err != nil {
		return err
	}

	err = p.Save(20*vg.Centimeter, 10*vg.Centimeter, pp.GetEscapedDefectsChartLocation())
	if err != nil {
		return err
	}

	var color string
	switch {
	case trend < 1.0:
		color = "green"
	case trend >= 1.0 && trend < 1.5:
		color = "amber"
	case trend >= 1.5:
		color = "red"
	}

	var lastMove string
	switch {
	case len(defectCounts) < 2:
		lastMove = "static"
	case defectCounts[0].Count > defectCounts[1].Count:
		lastMove = "up"
	case defectCounts[0].Count < defectCounts[1].Count:
		lastMove = "down"
	default:
		lastMove = "static"
	}

	log.Printf("> Escaped Defects: Trend -> %s; Last Move -> %s", color, lastMove)
	return nil
}

func (pp *PlotPrinter) printCycleTimes(cycleTimeReports []models.WeekCount, allCycleTimes []jiracalculationsrepository.CycleTimes) error {
	trend, p, err := pp.printChart(cycleTimeReports, "Cycle Time", color.NRGBA{R: 0, G: 0, B: 190, A: 100}, color.NRGBA{R: 0, G: 0, B: 190, A: 255})
	if err != nil {
		return err
	}

	p, err = pp.addPoints(p, allCycleTimes, color.NRGBA{R: 0, G: 0, B: 190, A: 110})
	if err != nil {
		return err
	}

	err = p.Save(20*vg.Centimeter, 10*vg.Centimeter, pp.GetCycleTimeChartLocation())
	if err != nil {
		return err
	}

	var color string
	switch {
	case trend < 1.0:
		color = "green"
	case trend >= 1.0 && trend < 1.5:
		color = "amber"
	case trend >= 1.5:
		color = "red"
	}

	var lastMove string
	switch {
	case len(cycleTimeReports) < 2:
		lastMove = "static"
	case cycleTimeReports[0].Count > cycleTimeReports[1].Count:
		lastMove = "up"
	case cycleTimeReports[0].Count < cycleTimeReports[1].Count:
		lastMove = "down"
	default:
		lastMove = "static"
	}

	log.Printf("> Cycle Time: Trend -> %s; Last Move -> %s", color, lastMove)
	return nil
}

func (pp *PlotPrinter) printThroughput(throughputReports []models.WeekCount) error {
	trend, p, err := pp.printChart(throughputReports, "Throughput", color.NRGBA{R: 0, G: 190, B: 0, A: 100}, color.NRGBA{R: 0, G: 190, B: 0, A: 255})
	if err != nil {
		return err
	}

	err = p.Save(20*vg.Centimeter, 10*vg.Centimeter, pp.GetThroughputChartLocation())
	if err != nil {
		return err
	}

	var color string
	switch {
	case trend >= 0.0:
		color = "green"
	case trend < 0.0 && trend > -1.5:
		color = "amber"
	case trend <= 1.5:
		color = "red"
	}

	var lastMove string
	switch {
	case len(throughputReports) < 2:
		lastMove = "static"
	case throughputReports[0].Count > throughputReports[1].Count:
		lastMove = "up"
	case throughputReports[0].Count < throughputReports[1].Count:
		lastMove = "down"
	default:
		lastMove = "static"
	}

	log.Printf("> Throughput: Trend -> %s; Last Move -> %s", color, lastMove)
	return nil
}

func (pp *PlotPrinter) printUnplannedWork(unplannedWorkReports []models.WeekCount) error {
	trend, p, err := pp.printChart(unplannedWorkReports, "Unplanned Work", color.NRGBA{R: 100, G: 0, B: 100, A: 100}, color.NRGBA{R: 100, G: 0, B: 100, A: 255})
	if err != nil {
		return err
	}

	err = p.Save(20*vg.Centimeter, 10*vg.Centimeter, pp.GetUnplannedWorkChartLocation())
	if err != nil {
		return err
	}

	var color string
	switch {
	case trend < 1.0:
		color = "green"
	case trend >= 1.0 && trend < 1.5:
		color = "amber"
	case trend >= 1.5:
		color = "red"
	}

	var lastMove string
	switch {
	case len(unplannedWorkReports) < 2:
		lastMove = "static"
	case unplannedWorkReports[0].Count > unplannedWorkReports[1].Count:
		lastMove = "up"
	case unplannedWorkReports[0].Count < unplannedWorkReports[1].Count:
		lastMove = "down"
	default:
		lastMove = "static"
	}

	log.Printf("> Unplanned Work: Trend -> %s; Last Move -> %s", color, lastMove)
	return nil
}

func (pp *PlotPrinter) printChart(weekCounts []models.WeekCount, title string, plotColor color.Color, trendlineColor color.Color) (trend float64, p *plot.Plot, err error) {
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
		return 0.0, p, err
	}

	line.Color = plotColor
	line.LineStyle.Dashes = []vg.Length{vg.Points(5), vg.Points(5)}
	points.Shape = draw.BoxGlyph{}
	points.Color = plotColor

	p.Add(line, points)

	linearRegression, gradiant, _ := mathutil.LinearRegression(xys)
	trend = gradiant * (7 * 24 * 60 * 60)

	linearRegressionLine, linearRegressionPoints, err := plotter.NewLinePoints(asPlotPoints(linearRegression))
	if err != nil {
		return trend, p, err
	}

	linearRegressionLine.Color = trendlineColor
	linearRegressionPoints.Color = trendlineColor

	p.Add(linearRegressionLine, linearRegressionPoints)

	return trend, p, nil
}

func (pp *PlotPrinter) addPoints(p *plot.Plot, cycleTimes []jiracalculationsrepository.CycleTimes, plotColor color.Color) (*plot.Plot, error) {
	data := make(plotter.XYs, len(cycleTimes))
	for i, d := range cycleTimes {
		data[i].X = float64(d.Completed.Unix())
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
