package jiramodels

import "time"

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
