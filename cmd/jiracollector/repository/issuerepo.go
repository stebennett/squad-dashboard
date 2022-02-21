package repository

import (
	"context"

	_ "github.com/lib/pq"
	"github.com/stebennett/squad-dashboard/cmd/jiracollector/models"
)

type IssueRepository interface {
	SaveIssue(ctx context.Context, jiraIssue models.JiraIssue) (*models.JiraIssue, error)
	SaveTransition(ctx context.Context, jiraTransition models.JiraTransition) (*models.JiraTransition, error)
}
