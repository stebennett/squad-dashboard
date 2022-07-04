package pagerdutyrepository

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/stebennett/squad-dashboard/pkg/pagerdutymodels"
)

type PagerDutyRepository interface {
	SaveOnCall(ctx context.Context, oncall pagerdutymodels.OnCall) (int64, error)
}

type PostgresPagerDutyRepository struct {
	db *sql.DB
}

func NewPostgresPagerDutyRepository(db *sql.DB) *PostgresPagerDutyRepository {
	return &PostgresPagerDutyRepository{
		db: db,
	}
}

func (p *PostgresPagerDutyRepository) SaveOnCall(ctx context.Context, oncall pagerdutymodels.OnCall) (int64, error) {
	insertStatement := `
		INSERT INTO pagerduty_oncalls(pd_user_id, pd_user_name, schedule_id, schedule_name, escalation_policy_id, escalation_policy_name, escalation_level, on_call_start, on_call_end)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		ON CONFLICT (pd_user_id, schedule_id, escalation_policy_id, escalation_level, on_call_start, on_call_end)
		DO NOTHING
	`

	result, err := p.db.ExecContext(ctx,
		insertStatement,
		oncall.User.Id,
		oncall.User.Name,
		oncall.Schedule.Id,
		oncall.Schedule.Name,
		oncall.EscalationPolicy.Id,
		oncall.EscalationPolicy.Name,
		oncall.EscalationLevel,
		oncall.Start,
		oncall.End,
	)

	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}
