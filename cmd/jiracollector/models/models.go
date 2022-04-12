package models

import (
	"encoding/json"
	"fmt"
	"time"
)

type JiraTimestamp struct {
	time.Time
}

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
		EpicKey string        `json:"epicKey"`
		Created JiraTimestamp `json:"created"`
		Updated JiraTimestamp `json:"updated"`
	} `json:"fields"`
	ChangeLog struct {
		StartAt    int `json:"startAt"`
		MaxResults int `json:"maxResults"`
		Total      int `json:"total"`
		Histories  []struct {
			Created JiraTimestamp `json:"created"`
			Items   []struct {
				Field      string `json:"field"`
				FromString string `json:"fromString"`
				ToString   string `json:"toString"`
			} `json:"items"`
		} `json:"histories"`
	} `json:"changelog"`
}

type JiraSearchResults struct {
	StartAt    int               `json:"startAt"`
	MaxResults int               `json:"maxResults"`
	Total      int               `json:"total"`
	Issues     []JiraResultIssue `json:"issues"`
}

func (p *JiraTimestamp) UnmarshalJSON(bytes []byte) error {
	var raw string
	err := json.Unmarshal(bytes, &raw)

	if err != nil {
		fmt.Printf("Failed to marshal timestamp - %s", err)
		return err
	}

	p.Time, err = time.Parse("2006-01-02T15:04:05.000-0700", raw)
	return err
}

func Create(issue JiraResultIssue) (*JiraIssue, error) {
	return &JiraIssue{
		Key:       issue.Key,
		IssueType: issue.Fields.IssueType.Name,
		ParentKey: issue.Fields.EpicKey,
		CreatedAt: issue.Fields.Created.Time.UTC(),
		UpdatedAt: issue.Fields.Updated.Time.UTC(),
	}, nil
}
