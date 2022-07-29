package statsmodels

import "time"

type WeeklyTimeItem struct {
	WeekStarting  time.Time
	NumberOfItems int
}

type WeeklyCycleTimeItem struct {
	WeekStarting time.Time
	CycleTime    float64
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

type CycleTimeResult struct {
	Project        string
	CycleTimeItems []WeeklyCycleTimeItem
	Trendline      []TrendlineItem
}

type ProjectWeeklyTimeItem struct {
	Project  string
	TimeItem WeeklyTimeItem
}
