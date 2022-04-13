package repository

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/stebennett/squad-dashboard/cmd/jiracollector/models"
)

type IssueRepository interface {
	SaveIssue(ctx context.Context, jiraIssue models.JiraIssue) (int64, error)
	SaveTransition(ctx context.Context, issueKey string, jiraTransition []models.JiraTransition) (int64, error)
}

type PostgresIssueRepository struct {
	db *sql.DB
}

func NewPostgresIssueRepository(db *sql.DB) *PostgresIssueRepository {
	return &PostgresIssueRepository{
		db: db,
	}
}

func (p *PostgresIssueRepository) SaveIssue(ctx context.Context, jiraIssue models.JiraIssue) (int64, error) {
	insertIssueStatement := `
		INSERT INTO jira_issues(issue_key, issue_type, parent_key, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (issue_key)
		DO UPDATE
		SET issue_type = $2, parent_key = $3, created_at = $4, updated_at = $5
		WHERE jira_issues.issue_key = $1
	`

	result, err := p.db.ExecContext(ctx,
		insertIssueStatement,
		jiraIssue.Key,
		jiraIssue.IssueType,
		jiraIssue.ParentKey,
		jiraIssue.CreatedAt,
		jiraIssue.UpdatedAt,
	)

	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

func (p *PostgresIssueRepository) SaveTransition(ctx context.Context, issueKey string, jiraTransitions []models.JiraTransition) (int64, error) {
	var inserted int64 = 0

	for _, transition := range jiraTransitions {
		insertTransitionStatement := `
			INSERT INTO jira_transitions(issue_key, from_state, to_state, created_at)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (issue_key, created_at)
			DO NOTHING
		`

		result, err := p.db.ExecContext(ctx,
			insertTransitionStatement,
			issueKey,
			transition.FromState,
			transition.ToState,
			transition.TransitionedAt,
		)

		if err != nil {
			return -1, err
		}

		rowsAffected, _ := result.RowsAffected()
		inserted = inserted + rowsAffected
	}

	return inserted, nil
}
