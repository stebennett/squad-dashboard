package models

import (
	"time"
)

type WeekCount struct {
	WeekEnding time.Time
	Count      int
}

type WorkItem struct {
	IssueKey         string
	CreatedAt        time.Time
	CompletedAt      time.Time
	WorkingCycleTime int
}

type TrendLineItem struct {
	X time.Time
	Y float64
}

type TrendDetails struct {
	XYs           []TrendLineItem
	Slope         float64
	CrossingPoint float64
}

type CycleTimeReport struct {
	WeeklyReports          []WeekCount
	AllCycleTimeItems      []WorkItem
	LastWeekCycleTimeItems []WorkItem
	Trend                  TrendDetails
}

type ThroughputReport struct {
	WeeklyReports           []WeekCount
	LastWeekThroughputItems []WorkItem
	Trend                   TrendDetails
}

type EscapedDefectReport struct {
	WeeklyReports              []WeekCount
	LastWeekEscapedDefectItems []WorkItem
	Trend                      TrendDetails
}

type UnplannedWorkReport struct {
	WeeklyReports              []WeekCount
	LastWeekUnplannedWorkItems []WorkItem
	Trend                      TrendDetails
}
