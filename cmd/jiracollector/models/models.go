package models

import (
	"time"
)

type Timestamp time.Time

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

type JiraResultIssue struct {
	Key    string `json:"key"`
	Fields struct {
		IssueType struct {
			Name string `json:"name"`
		} `json:"issuetype"`
		EpicKey string `json:"epicKey"`
		Created string `json:"created"`
		Updated string `json:"updated"`
	} `json:"fields"`
}

type JiraSearchResults struct {
	StartAt    int               `json:"startAt"`
	MaxResults int               `json:"maxResults"`
	Total      int               `json:"total"`
	Issues     []JiraResultIssue `json:"issues"`
}

func Create(issue JiraResultIssue) JiraIssue {
	return JiraIssue{
		Key:       issue.Key,
		IssueType: issue.Fields.IssueType.Name,
		ParentKey: issue.Fields.EpicKey,
		CreatedAt: time.Now(), // TODO: Update with correct time
		UpdatedAt: time.Now(), // TODO: Update with correct time
	}
}
