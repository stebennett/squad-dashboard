package models

import "time"

type Timestamp time.Time

type JiraIssue struct {
	Key       string
	IssueType string
	CreatedAt Timestamp
	UpdateAt  Timestamp
	ParentKey string
}

type JiraTransition struct {
	FromState      string
	ToState        string
	TransitionedAt Timestamp
}
