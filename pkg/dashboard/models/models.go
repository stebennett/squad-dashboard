package models

import (
	"time"
)

type WeekCount struct {
	WeekEnding time.Time
	Count      int
}

type CycleTimeItem struct {
	IssueKey    string
	CompletedAt time.Time
	Size        int
}

type ThroughputItem struct {
	IssueKey    string
	CompletedAt time.Time
}

type EscapedDefectItem struct {
	IssueKey  string
	CreatedAt time.Time
}

type UnplannedWorkItem ThroughputItem

type CycleTimeReport struct {
	WeeklyReports          []WeekCount
	AllCycleTimeItems      []CycleTimeItem
	LastWeekCycleTimeItems []CycleTimeItem
}

type ThroughputReport struct {
	WeeklyReports           []WeekCount
	LastWeekThroughputItems []ThroughputItem
}

type EscapedDefectReport struct {
	WeeklyReports              []WeekCount
	LastWeekEscapedDefectItems []EscapedDefectItem
}

type UnplannedWorkReport struct {
	WeeklyReports              []WeekCount
	LastWeekUnplannedWorkItems []UnplannedWorkItem
}
