package models

import "time"

type PagerDutyEntitySummary struct {
	Id   string
	Name string
}

type OnCall struct {
	User             PagerDutyEntitySummary
	Schedule         PagerDutyEntitySummary
	EscalationPolicy PagerDutyEntitySummary
	EscalationLevel  int
	Start            time.Time
	End              time.Time
}
