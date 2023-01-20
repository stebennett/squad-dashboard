package repo

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
)

type PostgresConfigRepository struct {
	db *sql.DB
}

func NewPostgresConfigRepository(db *sql.DB) *PostgresConfigRepository {
	return &PostgresConfigRepository{
		db: db,
	}
}

func (p *PostgresConfigRepository) SaveJiraToDoStates(ctx context.Context, project string, states []string) (int64, error) {
	return p.saveJiraWorkStates(ctx, project, states, "todo")
}

func (p *PostgresConfigRepository) SaveJiraInProgressStates(ctx context.Context, project string, states []string) (int64, error) {
	return p.saveJiraWorkStates(ctx, project, states, "in progress")
}

func (p *PostgresConfigRepository) SaveJiraDoneStates(ctx context.Context, project string, states []string) (int64, error) {
	return p.saveJiraWorkStates(ctx, project, states, "done")
}

func (p *PostgresConfigRepository) SaveNonWorkingDays(ctx context.Context, project string, nonWorkingDays []string) (int64, error) {
	insertStatement := `
		INSERT INTO non_working_days(project, non_working_day)
		VALUES ($1, $2)
		ON CONFLICT (project, non_working_day)
		DO NOTHING
	`
	var inserted int64 = 0

	for _, nwd := range nonWorkingDays {
		// convert to a date
		layout := "2006-01-02"
		nonWorkingDate, err := time.Parse(layout, nwd)
		if err != nil {
			return inserted, err
		}

		result, err := p.db.ExecContext(ctx,
			insertStatement,
			project,
			nonWorkingDate,
		)

		if err != nil {
			return inserted, err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return inserted, err
		}

		inserted = inserted + rowsAffected
	}

	return inserted, nil
}

func (p *PostgresConfigRepository) GetNonWorkingDays(ctx context.Context, project string) ([]time.Time, error) {
	selectStatement := `
		SELECT non_working_day
		FROM non_working_days
		WHERE project = $1
		ORDER BY non_working_day ASC
	`
	var result = []time.Time{}

	rows, err := p.db.QueryContext(ctx, selectStatement, project)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var nonWorkingDate time.Time

		err = rows.Scan(&nonWorkingDate)
		if err != nil {
			return result, nil
		}

		result = append(result, nonWorkingDate)
	}

	return result, nil
}

func (p *PostgresConfigRepository) saveJiraWorkStates(ctx context.Context, project string, states []string, workState string) (int64, error) {
	insertStatement := `
		INSERT INTO jira_work_states(project, state_type, state_name)
		VALUES ($1, $2, $3)
		ON CONFLICT (project, state_name)
		DO UPDATE
		SET state_type=$2
		WHERE jira_work_states.project=$1 AND jira_work_states.state_name=$3
	`
	var inserted int64 = 0

	for _, stateName := range states {
		result, err := p.db.ExecContext(ctx,
			insertStatement,
			project,
			workState,
			stateName,
		)

		if err != nil {
			return inserted, err
		}

		rowsAffected, err := result.RowsAffected()
		if err != nil {
			return inserted, err
		}

		inserted = inserted + rowsAffected
	}

	return inserted, nil
}
