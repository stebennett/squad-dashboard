package statsmodels

import "time"

type ThroughputItem struct {
	WeekStarting  time.Time
	NumberOfItems int
}
