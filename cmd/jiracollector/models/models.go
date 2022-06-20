package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/stebennett/squad-dashboard/pkg/jiramodels"
)

type JiraTimestamp struct {
	time.Time
}

type JiraHistory struct {
	Created JiraTimestamp `json:"created"`
	Items   []struct {
		Field      string `json:"field"`
		FromString string `json:"fromString"`
		ToString   string `json:"toString"`
	} `json:"items"`
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
		Labels  []string      `json:"labels"`
	} `json:"fields"`
	ChangeLog struct {
		StartAt    int           `json:"startAt"`
		MaxResults int           `json:"maxResults"`
		Total      int           `json:"total"`
		Histories  []JiraHistory `json:"histories"`
	} `json:"changelog"`
}

type JiraSearchResults struct {
	StartAt    int               `json:"startAt"`
	MaxResults int               `json:"maxResults"`
	Total      int               `json:"total"`
	Issues     []JiraResultIssue `json:"issues"`
}

type JiraChangeLogHistoryResults struct {
	StartAt    int           `json:"startAt"`
	MaxResults int           `json:"maxResults"`
	Total      int           `json:"total"`
	Histories  []JiraHistory `json:"values"`
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

func Create(issue JiraResultIssue) (*jiramodels.JiraIssue, error) {
	return &jiramodels.JiraIssue{
		Key:       issue.Key,
		IssueType: issue.Fields.IssueType.Name,
		ParentKey: issue.Fields.EpicKey,
		CreatedAt: issue.Fields.Created.Time.UTC(),
		UpdatedAt: issue.Fields.Updated.Time.UTC(),
		Labels:    issue.Fields.Labels,
	}, nil
}

func CreateTransitions(histories []JiraHistory) (result []jiramodels.JiraTransition) {
	for _, h := range histories {
		for _, i := range h.Items {
			if i.Field == "status" {
				jt := jiramodels.JiraTransition{
					FromState:      i.FromString,
					ToState:        i.ToString,
					TransitionedAt: h.Created.Time,
				}
				result = append(result, jt)
			}
		}
	}

	return result
}
