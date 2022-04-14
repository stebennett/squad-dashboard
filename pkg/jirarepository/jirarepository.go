package jirarepository

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/stebennett/squad-dashboard/pkg/jiramodels"
)

type JiraRepository interface {
	SaveIssue(ctx context.Context, project string, jiraIssue jiramodels.JiraIssue) (int64, error)
	SaveTransition(ctx context.Context, issueKey string, jiraTransition []jiramodels.JiraTransition) (int64, error)
	GetIssuesWithStateTransition(ctx context.Context, toState string) ([]string, error)
	GetTransitionsForIssue(ctx context.Context, issueKey string) ([]jiramodels.JiraTransition, error)
	SaveCreateWeekDate(ctx context.Context, issueKey string, year int, week int) (int64, error)
	SaveStartWeekDate(ctx context.Context, issueKey string, year int, week int) (int64, error)
	SaveCompleteWeekDate(ctx context.Context, issueKey string, year int, week int) (int64, error)
	SaveCycleTime(ctx context.Context, issueKey string, cycleTime int, workingCycleTime int) (int64, error)
	SaveLeadTime(ctx context.Context, issueKey string, leadTime int, workingLeadTime int) (int64, error)
	SaveSystemDelayTime(ctx context.Context, issueKey string, systemDelayTime int, workingSystemDelayTime int) (int64, error)
}

type PostgresJiraRepository struct {
	db *sql.DB
}

func NewPostgresJiraRepository(db *sql.DB) *PostgresJiraRepository {
	return &PostgresJiraRepository{
		db: db,
	}
}

func (p *PostgresJiraRepository) SaveIssue(ctx context.Context, project string, jiraIssue jiramodels.JiraIssue) (int64, error) {
	insertIssueStatement := `
		INSERT INTO jira_issues(issue_key, issue_type, parent_key, created_at, updated_at, project)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (issue_key)
		DO UPDATE
		SET issue_type = $2, parent_key = $3, created_at = $4, updated_at = $5, project = $6
		WHERE jira_issues.issue_key = $1
	`

	result, err := p.db.ExecContext(ctx,
		insertIssueStatement,
		jiraIssue.Key,
		jiraIssue.IssueType,
		jiraIssue.ParentKey,
		jiraIssue.CreatedAt,
		jiraIssue.UpdatedAt,
		project,
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

func (p *PostgresJiraRepository) SaveCreateWeekDate(ctx context.Context, issueKey string, year int, week int) (int64, error) {
	insertStatement := `
		INSERT INTO jira_issues_calculations(issue_key, create_week, create_year)
		VALUES ($1, $2, $3)
		ON CONFLICT (issue_key)
		DO UPDATE
		SET create_week = $2, create_year = $3
		WHERE jira_issues_calculations.issue_key = $1
	`

	result, err := p.db.ExecContext(ctx,
		insertStatement,
		issueKey,
		week,
		year,
	)

	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

func (p *PostgresJiraRepository) SaveStartWeekDate(ctx context.Context, issueKey string, year int, week int) (int64, error) {
	insertStatement := `
		INSERT INTO jira_issues_calculations(issue_key, start_week, start_year)
		VALUES ($1, $2, $3)
		ON CONFLICT (issue_key)
		DO UPDATE
		SET start_week = $2, start_year = $3
		WHERE jira_issues_calculations.issue_key = $1
	`

	result, err := p.db.ExecContext(ctx,
		insertStatement,
		issueKey,
		week,
		year,
	)

	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

func (p *PostgresJiraRepository) SaveCompleteWeekDate(ctx context.Context, issueKey string, year int, week int) (int64, error) {
	insertStatement := `
		INSERT INTO jira_issues_calculations(issue_key, complete_week, complete_year)
		VALUES ($1, $2, $3)
		ON CONFLICT (issue_key)
		DO UPDATE
		SET complete_week = $2, complete_year = $3
		WHERE jira_issues_calculations.issue_key = $1
	`

	result, err := p.db.ExecContext(ctx,
		insertStatement,
		issueKey,
		week,
		year,
	)

	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

func (p *PostgresJiraRepository) SaveCycleTime(ctx context.Context, issueKey string, cycleTime int, workingCycleTime int) (int64, error) {
	insertStatement := `
		INSERT INTO jira_issues_calculations(issue_key, cycle_time, working_cycle_time)
		VALUES ($1, $2, $3)
		ON CONFLICT (issue_key)
		DO UPDATE
		SET cycle_time = $2, working_cycle_time = $3
		WHERE jira_issues_calculations.issue_key = $1
	`

	result, err := p.db.ExecContext(ctx,
		insertStatement,
		issueKey,
		cycleTime,
		workingCycleTime,
	)

	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

func (p *PostgresJiraRepository) SaveLeadTime(ctx context.Context, issueKey string, leadTime int, workingLeadTime int) (int64, error) {
	insertStatement := `
		INSERT INTO jira_issues_calculations(issue_key, lead_time, working_lead_time)
		VALUES ($1, $2, $3)
		ON CONFLICT (issue_key)
		DO UPDATE
		SET lead_time = $2, working_lead_time = $3
		WHERE jira_issues_calculations.issue_key = $1
	`

	result, err := p.db.ExecContext(ctx,
		insertStatement,
		issueKey,
		leadTime,
		workingLeadTime,
	)

	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

func (p *PostgresJiraRepository) SaveSystemDelayTime(ctx context.Context, issueKey string, systemDelayTime int, workingSystemDelayTime int) (int64, error) {
	insertStatement := `
		INSERT INTO jira_issues_calculations(issue_key, system_delay_time, working_system_delay_time)
		VALUES ($1, $2, $3)
		ON CONFLICT (issue_key)
		DO UPDATE
		SET system_delay_time = $2, working_system_delay_time = $3
		WHERE jira_issues_calculations.issue_key = $1
	`

	result, err := p.db.ExecContext(ctx,
		insertStatement,
		issueKey,
		systemDelayTime,
		workingSystemDelayTime,
	)

	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}
