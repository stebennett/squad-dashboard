package adapters

import (
	"github.com/stebennett/squad-dashboard/cmd/jiracollector/models"
	"github.com/stebennett/squad-dashboard/pkg/jiraservice"
)

func AdaptIssue(jiraIssue jiraservice.JiraIssue) models.JiraIssue {
	return models.JiraIssue{}
}
