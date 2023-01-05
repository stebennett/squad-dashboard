package models

import "time"

type WeekCount struct {
	WeekEnding time.Time
	Count      int
}
