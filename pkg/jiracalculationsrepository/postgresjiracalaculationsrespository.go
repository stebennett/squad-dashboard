package jiracalculationsrepository

import (
	"context"
	"database/sql"
	"time"
)

type PostgresJiraCalculationsRepository struct {
	db *sql.DB
}

func NewPostgresJiraCalculationsRepository(db *sql.DB) *PostgresJiraCalculationsRepository {
	return &PostgresJiraCalculationsRepository{
		db: db,
	}
}

func (p *PostgresJiraCalculationsRepository) SaveCycleTime(ctx context.Context, issueKey string, cycleTime int, workingCycleTime int) (int64, error) {
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

func (p *PostgresJiraCalculationsRepository) SaveLeadTime(ctx context.Context, issueKey string, leadTime int, workingLeadTime int) (int64, error) {
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

func (p *PostgresJiraCalculationsRepository) SaveSystemDelayTime(ctx context.Context, issueKey string, systemDelayTime int, workingSystemDelayTime int) (int64, error) {
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

func (p *PostgresJiraCalculationsRepository) SaveCreateDates(ctx context.Context, issueKey string, year int, week int, createdAt time.Time) (int64, error) {
	insertStatement := `
		INSERT INTO jira_issues_calculations(issue_key, week_create, year_create, issue_created_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (issue_key)
		DO UPDATE
		SET week_create = $2, year_create = $3, issue_created_at = $4
		WHERE jira_issues_calculations.issue_key = $1
	`

	result, err := p.db.ExecContext(ctx,
		insertStatement,
		issueKey,
		week,
		year,
		createdAt,
	)

	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

func (p *PostgresJiraCalculationsRepository) SaveStartDates(ctx context.Context, issueKey string, year int, week int, startedAt time.Time) (int64, error) {
	insertStatement := `
		INSERT INTO jira_issues_calculations(issue_key, year_start, week_start, issue_started_at)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (issue_key)
		DO UPDATE
		SET year_start = $2, week_start = $3, issue_started_at = $4
		WHERE jira_issues_calculations.issue_key = $1
	`

	result, err := p.db.ExecContext(ctx,
		insertStatement,
		issueKey,
		year,
		week,
		startedAt,
	)

	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

func (p *PostgresJiraCalculationsRepository) SaveCompleteDates(ctx context.Context, issueKey string, year int, week int, completedAt time.Time, endState string) (int64, error) {
	insertStatement := `
		INSERT INTO jira_issues_calculations(issue_key, year_complete, week_complete, issue_completed_at, issue_end_state)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (issue_key)
		DO UPDATE
		SET year_complete = $2, week_complete = $3, issue_completed_at = $4, issue_end_state = $5
		WHERE jira_issues_calculations.issue_key = $1
	`

	result, err := p.db.ExecContext(ctx,
		insertStatement,
		issueKey,
		year,
		week,
		completedAt,
		endState,
	)

	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}
