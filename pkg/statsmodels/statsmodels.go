package statsmodels

import "time"

type WeeklyTimeItem struct {
	WeekStarting  time.Time
	NumberOfItems int
}

type TrendlineItem struct {
	WeekStarting time.Time
	TrendPoint   float64
}

type ThrouputResult struct {
	Project         string
	ThroughputItems []WeeklyTimeItem
	Trendline       []TrendlineItem
}

type ProjectWeeklyTimeItem struct {
	Project  string
	TimeItem WeeklyTimeItem
}
