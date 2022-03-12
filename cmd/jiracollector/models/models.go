package models

import (
	"time"

	"github.com/stebennett/squad-dashboard/pkg/jiraservice"
)

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

func Create(issue jiraservice.JiraIssue) JiraIssue {
	return JiraIssue{
		Key:       issue.Key,
		IssueType: issue.Fields.IssueType.Name,
		ParentKey: "",                    // TODO: Add parent key
		CreatedAt: Timestamp(time.Now()), // TODO: Update with correct time
		UpdateAt:  Timestamp(time.Now()), // TODO: Update with correct time
	}
}
