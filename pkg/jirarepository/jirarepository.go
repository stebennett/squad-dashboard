package jirarepository

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/stebennett/squad-dashboard/pkg/jiramodels"
)

type JiraRepository interface {
	SaveIssue(ctx context.Context, jiraIssue jiramodels.JiraIssue) (int64, error)
	SaveTransition(ctx context.Context, issueKey string, jiraTransition []jiramodels.JiraTransition) (int64, error)
}

type PostgresJiraRepository struct {
	db *sql.DB
}

func NewPostgresJiraRepository(db *sql.DB) *PostgresJiraRepository {
	return &PostgresJiraRepository{
		db: db,
	}
}

func (p *PostgresJiraRepository) SaveIssue(ctx context.Context, jiraIssue jiramodels.JiraIssue) (int64, error) {
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

func (p *PostgresJiraRepository) SaveTransition(ctx context.Context, issueKey string, jiraTransitions []jiramodels.JiraTransition) (int64, error) {
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

func (p *PostgresJiraRepository) GetIssuesWithStateTransition(ctx context.Context, toState string) ([]string, error) {
	selectStatement := `
		SELECT DISTINCT issue_key
		FROM jira_transitions
		WHERE to_state = $1
	`

	rows, err := p.db.QueryContext(ctx, selectStatement, toState)
	if err != nil {
		return []string{}, err
	}

	var result = []string{}

	for rows.Next() {
		var nextIssueKey string
		err = rows.Scan(&nextIssueKey)
		if err != nil {
			return result, nil
		}

		result = append(result, nextIssueKey)
	}

	return result, nil
}

func (p *PostgresJiraRepository) GetTransitionsForIssue(ctx context.Context, issueKey string) ([]jiramodels.JiraTransition, error) {
	selectStatement := `
		SELECT from_state, to_state, created_at
		FROM jira_transitions
		WHERE issue_key = $1
	`
	rows, err := p.db.QueryContext(ctx, selectStatement, issueKey)
	if err != nil {
		return []jiramodels.JiraTransition{}, err
	}

	var result = []jiramodels.JiraTransition{}

	for rows.Next() {
		transition := jiramodels.JiraTransition{}

		err = rows.Scan(&transition.FromState, &transition.ToState, &transition.TransitionedAt)
		if err != nil {
			return result, err
		}

		result = append(result, transition)
	}

	return result, nil
}
