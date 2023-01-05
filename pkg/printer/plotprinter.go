package printer

import (
	"image/color"

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

func (pp *PlotPrinter) PrintDefectCounts(defectCounts []models.EscapedDefectCount) error {
	p := plot.New()

	xticks := plot.TimeTicks{Format: "2006-01-02"}

	p.Title.Text = "Escaped Defects - " + pp.JiraProject
	p.X.Label.Text = "Week Ending"
	p.X.Tick.Marker = xticks
	p.Y.Label.Text = "Number"

	data := make(plotter.XYs, len(defectCounts))
	for i, d := range defectCounts {
		data[i].X = float64(d.WeekEnding.Unix())
		data[i].Y = float64(d.NumberOfDefectsCreated)
	}

	p.Add(plotter.NewGrid())

	line, points, err := plotter.NewLinePoints(data)
	if err != nil {
		return err
	}

	line.Color = color.RGBA{R: 255, A: 255}
	points.Shape = draw.BoxGlyph{}
	points.Color = color.RGBA{R: 255, A: 255}

	p.Add(line, points)

	err = p.Save(20*vg.Centimeter, 10*vg.Centimeter, pp.OutputDirectory+"/escapeddefects-"+pp.JiraProject+".png")
	return err
}

func (pp *PlotPrinter) PrintCycleTimes(cycleTimeReports []models.CycleTimeReport) error {
	return nil
}

func (pp *PlotPrinter) PrintThroughput(throughputReports []models.ThroughputReport) error {
	return nil
}
