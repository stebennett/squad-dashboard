package models

import (
	"time"

	"github.com/lib/pq"
)

type JiraIssue struct {
	Key       string
	IssueType string
	CreatedAt time.Time
	UpdatedAt time.Time
	ParentKey string
	Labels    []string
}

type JiraTransition struct {
	Key            string
	FromState      string
	ToState        string
	TransitionedAt time.Time
}

type IssueCalculations struct {
	IssueKey         string
	CycleTime        int
	WorkingCycleTime int
	LeadTime         int
	SystemDelayTime  int
	IssueCreatedAt   pq.NullTime
	IssueStartedAt   pq.NullTime
	IssueCompletedAt pq.NullTime
}
