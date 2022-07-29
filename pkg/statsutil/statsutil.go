package statsutil

import (
	"time"

	"github.com/montanaflynn/stats"
	"github.com/stebennett/squad-dashboard/pkg/statsmodels"
)

func CalculateTrendline(items []statsmodels.WeeklyTimeItem) ([]statsmodels.TrendlineItem, error) {
	series := stats.Series{}
	for _, item := range items {
		coord := stats.Coordinate{
			X: float64(item.WeekStarting.Unix()),
			Y: float64(item.NumberOfItems),
		}

		series = append(series, coord)
	}

	linearRegression, err := stats.LinearRegression(series)
	if err != nil {
		return []statsmodels.TrendlineItem{}, err
	}

	trendline := []statsmodels.TrendlineItem{}
	for _, lr := range linearRegression {
		wti := statsmodels.TrendlineItem{
			WeekStarting: time.Unix(int64(lr.X), 0),
			TrendPoint:   lr.Y,
		}

		trendline = append(trendline, wti)
	}

	return trendline, nil
}
