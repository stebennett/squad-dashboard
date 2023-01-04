package models

import "time"

type EscapedDefectCount struct {
	WeekEnding             time.Time
	NumberOfDefectsCreated int
}
