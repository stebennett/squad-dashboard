package jiramodels

import "time"

type JiraIssue struct {
	Key       string
	IssueType string
	CreatedAt time.Time
	UpdatedAt time.Time
	ParentKey string
}

type JiraTransition struct {
	FromState      string
	ToState        string
	TransitionedAt time.Time
}
