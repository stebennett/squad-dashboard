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

type CycleTimeReport struct {
	WeeklyReports          []WeekCount
	AllCycleTimeItems      []WorkItem
	LastWeekCycleTimeItems []WorkItem
}

type ThroughputReport struct {
	WeeklyReports           []WeekCount
	LastWeekThroughputItems []WorkItem
}

type EscapedDefectReport struct {
	WeeklyReports              []WeekCount
	LastWeekEscapedDefectItems []WorkItem
}

type UnplannedWorkReport struct {
	WeeklyReports              []WeekCount
	LastWeekUnplannedWorkItems []WorkItem
}
