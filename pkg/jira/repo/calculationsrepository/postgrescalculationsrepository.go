package calculationsrepository

import (
	"context"
	"database/sql"
	"time"

	"github.com/lib/pq"
	"github.com/stebennett/squad-dashboard/pkg/jira/models"
)

type PostgresJiraCalculationsRepository struct {
	db *sql.DB
}

func NewPostgresJiraCalculationsRepository(db *sql.DB) *PostgresJiraCalculationsRepository {
	return &PostgresJiraCalculationsRepository{
		db: db,
	}
}

func (p *PostgresJiraCalculationsRepository) DropAllCalculations(ctx context.Context, project string) (int64, error) {
	insertStatement := `
		DELETE FROM jira_issues_calculations jic
		USING jira_issues ji
		WHERE ji.project=$1 AND jic.issue_key = ji.issue_key
	`

	result, err := p.db.ExecContext(ctx, insertStatement, project)

	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
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

func (p *PostgresJiraCalculationsRepository) GetEscapedDefects(ctx context.Context, project string, issueType string, startDate time.Time, endDate time.Time) ([]models.IssueCalculations, error) {
	selectStatement := `
		SELECT jira_issues_calculations.issue_key, 
			jira_issues_calculations.cycle_time,
			jira_issues_calculations.working_cycle_time,
			jira_issues_calculations.lead_time,
			jira_issues_calculations.system_delay_time,
			jira_issues_calculations.issue_created_at,
			jira_issues_calculations.issue_started_at,
			jira_issues_calculations.issue_completed_at
		FROM jira_issues_calculations
		JOIN jira_issues ON jira_issues_calculations.issue_key = jira_issues.issue_key
		WHERE jira_issues_calculations.issue_created_at > $3
		AND jira_issues_calculations.issue_created_at <= $4
		AND jira_issues.issue_type = $2
		AND jira_issues.project = $1
	`
	var result = []models.IssueCalculations{}

	rows, err := p.db.QueryContext(ctx, selectStatement, project, issueType, startDate, endDate)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var item models.IssueCalculations

		err = rows.Scan(&item.IssueKey,
			&item.CycleTime,
			&item.WorkingCycleTime,
			&item.LeadTime,
			&item.SystemDelayTime,
			&item.IssueCreatedAt,
			&item.IssueStartedAt,
			&item.IssueCompletedAt,
		)
		if err != nil {
			return result, nil
		}

		result = append(result, item)
	}

	return result, nil
}

func (p *PostgresJiraCalculationsRepository) GetCompletedWorkingCycleTimes(ctx context.Context, project string, issueTypes []string, startDate time.Time, endDate time.Time) ([]models.IssueCalculations, error) {
	selectStatement := `
		SELECT jira_issues_calculations.issue_key, 
			jira_issues_calculations.cycle_time,
			jira_issues_calculations.working_cycle_time,
			jira_issues_calculations.lead_time,
			jira_issues_calculations.system_delay_time,
			jira_issues_calculations.issue_created_at,
			jira_issues_calculations.issue_started_at,
			jira_issues_calculations.issue_completed_at
		FROM jira_issues_calculations
		JOIN jira_issues ON jira_issues_calculations.issue_key = jira_issues.issue_key
		WHERE jira_issues_calculations.issue_completed_at > $3
		AND jira_issues_calculations.issue_completed_at <= $4
		AND jira_issues_calculations.working_cycle_time > -1
		AND jira_issues.issue_type = ANY($2)
		AND jira_issues.project = $1
	`
	var result = []models.IssueCalculations{}

	rows, err := p.db.QueryContext(ctx, selectStatement, project, pq.Array(issueTypes), startDate, endDate)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var item models.IssueCalculations

		err = rows.Scan(&item.IssueKey,
			&item.CycleTime,
			&item.WorkingCycleTime,
			&item.LeadTime,
			&item.SystemDelayTime,
			&item.IssueCreatedAt,
			&item.IssueStartedAt,
			&item.IssueCompletedAt,
		)
		if err != nil {
			return result, nil
		}

		result = append(result, item)
	}

	return result, nil
}

func (p *PostgresJiraCalculationsRepository) GetThroughput(ctx context.Context, project string, issueTypes []string, startDate time.Time, endDate time.Time) ([]models.IssueCalculations, error) {
	selectStatement := `
		SELECT jira_issues_calculations.issue_key, 
			jira_issues_calculations.cycle_time,
			jira_issues_calculations.working_cycle_time,
			jira_issues_calculations.lead_time,
			jira_issues_calculations.system_delay_time,
			jira_issues_calculations.issue_created_at,
			jira_issues_calculations.issue_started_at,
			jira_issues_calculations.issue_completed_at
		FROM jira_issues_calculations
		JOIN jira_issues ON jira_issues_calculations.issue_key = jira_issues.issue_key
		WHERE jira_issues_calculations.issue_completed_at > $3
		AND jira_issues_calculations.issue_completed_at <= $4
		AND jira_issues_calculations.working_cycle_time > -1
		AND jira_issues.issue_type = ANY($2)
		AND jira_issues.project = $1
	`
	var result = []models.IssueCalculations{}

	rows, err := p.db.QueryContext(ctx, selectStatement, project, pq.Array(issueTypes), startDate, endDate)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var item models.IssueCalculations

		err = rows.Scan(&item.IssueKey,
			&item.CycleTime,
			&item.WorkingCycleTime,
			&item.LeadTime,
			&item.SystemDelayTime,
			&item.IssueCreatedAt,
			&item.IssueStartedAt,
			&item.IssueCompletedAt,
		)
		if err != nil {
			return result, nil
		}

		result = append(result, item)
	}

	return result, nil
}

func (p *PostgresJiraCalculationsRepository) GetUnplannedThroughput(ctx context.Context, project string, issueTypes []string, startDate time.Time, endDate time.Time) ([]models.IssueCalculations, error) {
	selectStatement := `
		SELECT jira_issues_calculations.issue_key, 
			jira_issues_calculations.cycle_time,
			jira_issues_calculations.working_cycle_time,
			jira_issues_calculations.lead_time,
			jira_issues_calculations.system_delay_time,
			jira_issues_calculations.issue_created_at,
			jira_issues_calculations.issue_started_at,
			jira_issues_calculations.issue_completed_at
		FROM jira_issues_calculations
		JOIN jira_issues ON jira_issues_calculations.issue_key = jira_issues.issue_key
		WHERE jira_issues_calculations.issue_completed_at > $3
		AND jira_issues_calculations.issue_completed_at <= $4
		AND jira_issues_calculations.working_cycle_time > -1
		AND jira_issues.issue_type = ANY($2)
		AND jira_issues.project = $1
		AND jira_issues.unplanned = TRUE
	`
	var result = []models.IssueCalculations{}

	rows, err := p.db.QueryContext(ctx, selectStatement, project, pq.Array(issueTypes), startDate, endDate)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var item models.IssueCalculations

		err = rows.Scan(&item.IssueKey,
			&item.CycleTime,
			&item.WorkingCycleTime,
			&item.LeadTime,
			&item.SystemDelayTime,
			&item.IssueCreatedAt,
			&item.IssueStartedAt,
			&item.IssueCompletedAt,
		)
		if err != nil {
			return result, nil
		}

		result = append(result, item)
	}

	return result, nil
}

func (p *PostgresJiraCalculationsRepository) GetCompletedIssues(ctx context.Context, project string) (map[string]models.IssueCalculations, error) {
	selectStatement := `
		SELECT jira_issues_calculations.issue_key, 
			jira_issues_calculations.cycle_time,
			jira_issues_calculations.working_cycle_time,
			jira_issues_calculations.lead_time,
			jira_issues_calculations.system_delay_time,
			jira_issues_calculations.issue_created_at,
			jira_issues_calculations.issue_started_at,
			jira_issues_calculations.issue_completed_at
		FROM jira_issues_calculations
		INNER JOIN jira_issues ON jira_issues_calculations.issue_key = jira_issues.issue_key
		WHERE jira_issues_calculations.issue_completed_at IS NOT NULL
		AND jira_issues.project = $1
	`
	var result = make(map[string]models.IssueCalculations)

	rows, err := p.db.QueryContext(ctx, selectStatement, project)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var item models.IssueCalculations

		err = rows.Scan(&item.IssueKey,
			&item.CycleTime,
			&item.WorkingCycleTime,
			&item.LeadTime,
			&item.SystemDelayTime,
			&item.IssueCreatedAt,
			&item.IssueStartedAt,
			&item.IssueCompletedAt,
		)
		if err != nil {
			return result, nil
		}

		result[item.IssueKey] = item
	}

	return result, nil
}

func (p *PostgresJiraCalculationsRepository) GetIssuesStartedBetweenDates(ctx context.Context, project string, startDate time.Time, endDate time.Time, issueTypes []string) ([]string, error) {
	selectStatement := `
		SELECT jira_issues_calculations.issue_key from jira_issues_calculations
		INNER JOIN jira_issues ON jira_issues_calculations.issue_key = jira_issues.issue_key
		WHERE jira_issues.project = $1
		AND jira_issues_calculations.issue_started_at >= $2
		AND jira_issues_calculations.issue_started_at < $3
		AND jira_issues.issue_type = ANY($4)
	`

	var result []string

	rows, err := p.db.QueryContext(ctx,
		selectStatement,
		project,
		startDate,
		endDate,
		pq.Array(issueTypes),
	)

	if err != nil {
		return result, err
	}

	for rows.Next() {
		var issueKey string

		err = rows.Scan(&issueKey)
		if err != nil {
			return result, err
		}

		result = append(result, issueKey)
	}

	return result, nil
}

func (p *PostgresJiraCalculationsRepository) GetIssuesCompletedBetweenDates(ctx context.Context, project string, startDate time.Time, endDate time.Time, issueTypes []string, endStates []string) ([]string, error) {
	selectStatement := `
		SELECT jira_issues_calculations.issue_key from jira_issues_calculations
		INNER JOIN jira_issues ON jira_issues_calculations.issue_key = jira_issues.issue_key
		WHERE jira_issues.project = $1
		AND jira_issues_calculations.issue_completed_at >= $2
		AND jira_issues_calculations.issue_completed_at < $3
		AND jira_issues_calculations.issue_end_state = ANY($5)
		AND jira_issues.issue_type = ANY($4)
	`

	var result []string

	rows, err := p.db.QueryContext(ctx,
		selectStatement,
		project,
		startDate,
		endDate,
		pq.Array(issueTypes),
		pq.Array(endStates),
	)

	if err != nil {
		return result, err
	}

	for rows.Next() {
		var issueKey string

		err = rows.Scan(&issueKey)
		if err != nil {
			return result, err
		}

		result = append(result, issueKey)
	}

	return result, nil
}
