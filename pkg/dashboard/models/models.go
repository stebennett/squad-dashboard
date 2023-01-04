package models

import "time"

type EscapedDefectCount struct {
	WeekEnding             time.Time
	NumberOfDefectsCreated int
}

type CycleTimeReport struct {
	WeekEnding       time.Time
	AverageCycleTime int
}

type ThroughputReport struct {
	WeekEnding time.Time
	Throughput int
}
