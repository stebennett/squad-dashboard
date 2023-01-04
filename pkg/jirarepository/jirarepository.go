package jirarepository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/stebennett/squad-dashboard/pkg/calculatormodels"
	"github.com/stebennett/squad-dashboard/pkg/jiramodels"
	"github.com/stebennett/squad-dashboard/pkg/statsmodels"
)

type JiraRepository interface {
	GetIssues(ctx context.Context, project string) ([]jiramodels.JiraIssue, error)
	SaveIssue(ctx context.Context, project string, jiraIssue jiramodels.JiraIssue) (int64, error)
	SaveTransition(ctx context.Context, issueKey string, jiraTransition []jiramodels.JiraTransition) (int64, error)
	GetTransitionTimeByStateChanges(ctx context.Context, project string, fromStates []string, toStates []string) (map[string]time.Time, error)
	GetTransitionTimeByToState(ctx context.Context, project string, toStates []string) (map[string]time.Time, error)
	GetTransitionsForIssue(ctx context.Context, issueKey string) ([]jiramodels.JiraTransition, error)
	SaveCreateDates(ctx context.Context, issueKey string, year int, week int, createdAt time.Time) (int64, error)
	SaveStartDates(ctx context.Context, issueKey string, year int, week int, startedAt time.Time) (int64, error)
	SaveCompleteDates(ctx context.Context, issueKey string, year int, week int, completedAt time.Time, endState string) (int64, error)
	SaveCycleTime(ctx context.Context, issueKey string, cycleTime int, workingCycleTime int) (int64, error)
	SaveLeadTime(ctx context.Context, issueKey string, leadTime int, workingLeadTime int) (int64, error)
	SaveSystemDelayTime(ctx context.Context, issueKey string, systemDelayTime int, workingSystemDelayTime int) (int64, error)
	GetCompletedIssues(ctx context.Context, project string) (map[string]calculatormodels.IssueCalculations, error)
	GetIssuesStartedBetweenDates(ctx context.Context, project string, startDate time.Time, endDate time.Time, issueTypes []string) ([]string, error)
	SetIssuesStartedInWeekStarting(ctx context.Context, project string, startDate time.Time, count int) (int64, error)
	GetIssuesCompletedBetweenDates(ctx context.Context, project string, startDate time.Time, endDate time.Time, issueTypes []string, endStates []string) ([]string, error)
	SetIssuesCompletedInWeekStarting(ctx context.Context, project string, startDate time.Time, count int) (int64, error)
	GetEndStateForIssue(ctx context.Context, issueKey string, transitionDate time.Time) (string, error)
	SaveIssueLabels(ctx context.Context, issueKey string, label []string) (int64, error)
	SaveJiraToDoStates(ctx context.Context, project string, states []string) (int64, error)
	SaveJiraInProgressStates(ctx context.Context, project string, states []string) (int64, error)
	SaveJiraDoneStates(ctx context.Context, project string, states []string) (int64, error)
	GetWeeklyThroughputByProject(ctx context.Context, project string, endDate time.Time, numberOfWeeks int) ([]statsmodels.WeeklyTimeItem, error)
	GetWeeklyThroughputAllProjects(ctx context.Context, endDate time.Time, numberOfWeeks int) ([]statsmodels.ProjectWeeklyTimeItem, error)
	GetWeeklyCycleTimeByProject(ctx context.Context, project string, endDate time.Time, numberOfWeeks int) ([]statsmodels.WeeklyCycleTimeItem, error)
	GetWeeklyCycleTimeAllProjects(ctx context.Context, endDate time.Time, numberOfWeeks int) ([]statsmodels.ProjectWeeklyCycleTimeItem, error)
	GetProjects(ctx context.Context) ([]string, error)
	SaveNonWorkingDays(ctx context.Context, project string, nonWorkingDays []string) (int64, error)
	GetNonWorkingDays(ctx context.Context, project string) ([]time.Time, error)
}

type PostgresJiraRepository struct {
	db *sql.DB
}

func NewPostgresJiraRepository(db *sql.DB) *PostgresJiraRepository {
	return &PostgresJiraRepository{
		db: db,
	}
}

func (p *PostgresJiraRepository) GetIssues(ctx context.Context, project string) ([]jiramodels.JiraIssue, error) {
	selectStatement := `
		SELECT issue_key, parent_key, created_at, updated_at, issue_type
		FROM jira_issues
		WHERE project = $1
	`
	rows, err := p.db.QueryContext(ctx, selectStatement, project)
	if err != nil {
		return []jiramodels.JiraIssue{}, err
	}

	var result = []jiramodels.JiraIssue{}

	for rows.Next() {
		issue := jiramodels.JiraIssue{}

		err = rows.Scan(&issue.Key, &issue.ParentKey, &issue.CreatedAt, &issue.UpdatedAt, &issue.IssueType)
		if err != nil {
			return result, err
		}

		result = append(result, issue)
	}

	return result, nil
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

func (p *PostgresJiraRepository) GetTransitionTimeByStateChanges(ctx context.Context, project string, fromStates []string, toStates []string) (map[string]time.Time, error) {
	selectStatement := `
		SELECT jira_transitions.issue_key, MAX(jira_transitions.created_at)
		FROM jira_transitions
		INNER JOIN jira_issues ON jira_transitions.issue_key = jira_issues.issue_key
		WHERE jira_transitions.from_state = ANY($1) AND jira_transitions.to_state = ANY($2) 
		AND jira_issues.project = $3
		GROUP BY jira_transitions.issue_key
	`
	var result = make(map[string]time.Time)

	rows, err := p.db.QueryContext(ctx, selectStatement, pq.Array(fromStates), pq.Array(toStates), project)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var issueKey string
		var transitionTime time.Time

		err = rows.Scan(&issueKey, &transitionTime)
		if err != nil {
			return result, nil
		}

		result[issueKey] = transitionTime
	}

	return result, nil
}

func (p *PostgresJiraRepository) GetTransitionTimeByToState(ctx context.Context, project string, toStates []string) (map[string]time.Time, error) {
	selectStatement := `
		SELECT jira_transitions.issue_key, MAX(jira_transitions.created_at)
		FROM jira_transitions
		INNER JOIN jira_issues ON jira_transitions.issue_key = jira_issues.issue_key
		WHERE jira_transitions.to_state = ANY($1) 
		AND jira_issues.project = $2
		GROUP BY jira_transitions.issue_key
	`
	var result = make(map[string]time.Time)

	rows, err := p.db.QueryContext(ctx, selectStatement, pq.Array(toStates), project)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var issueKey string
		var transitionTime time.Time

		err = rows.Scan(&issueKey, &transitionTime)
		if err != nil {
			return result, nil
		}

		result[issueKey] = transitionTime
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

func (p *PostgresJiraRepository) SaveCreateDates(ctx context.Context, issueKey string, year int, week int, createdAt time.Time) (int64, error) {
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

func (p *PostgresJiraRepository) SaveStartDates(ctx context.Context, issueKey string, year int, week int, startedAt time.Time) (int64, error) {
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

func (p *PostgresJiraRepository) SaveCompleteDates(ctx context.Context, issueKey string, year int, week int, completedAt time.Time, endState string) (int64, error) {
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

func (p *PostgresJiraRepository) GetCompletedIssues(ctx context.Context, project string) (map[string]calculatormodels.IssueCalculations, error) {
	selectStatement := `
	SELECT jira_issues_calculations.issue_key, 
		jira_issues_calculations.cycle_time, 
		jira_issues_calculations.lead_time, 
		jira_issues_calculations.system_delay_time, 
		jira_issues_calculations.issue_created_at, 
		jira_issues_calculations.issue_started_at, 
		jira_issues_calculations.issue_completed_at
	FROM jira_issues_calculations
	INNER JOIN jira_issues ON jira_issues_calculations.issue_key = jira_issues.issue_key
	WHERE jira_issues_calculations.issue_completed_at IS NOT NULL
	AND jira_issues_calculations.issue_started_at IS NOT NULL
	AND jira_issues.project = $1
`
	var result = make(map[string]calculatormodels.IssueCalculations)

	rows, err := p.db.QueryContext(ctx, selectStatement, project)
	if err != nil {
		return result, err
	}

	for rows.Next() {
		var issueKey string
		var calculations calculatormodels.IssueCalculations

		err = rows.Scan(&issueKey,
			&calculations.CycleTime,
			&calculations.LeadTime,
			&calculations.LeadTime,
			&calculations.IssueCreatedAt,
			&calculations.IssueStartedAt,
			&calculations.IssueCompletedAt,
		)
		if err != nil {
			return result, err
		}

		result[issueKey] = calculations
	}

	return result, nil
}

func (p *PostgresJiraRepository) GetIssuesStartedBetweenDates(ctx context.Context, project string, startDate time.Time, endDate time.Time, issueTypes []string) ([]string, error) {
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

func (p *PostgresJiraRepository) SetIssuesStartedInWeekStarting(ctx context.Context, project string, startDate time.Time, count int) (int64, error) {
	insertStatement := `
		INSERT INTO jira_issues_reports(project, week_start, number_of_items_started)
		VALUES ($1, $2, $3)
		ON CONFLICT (project, week_start)
		DO UPDATE
		SET number_of_items_started = $3
		WHERE jira_issues_reports.project = $1 AND jira_issues_reports.week_start = $2
	`

	result, err := p.db.ExecContext(ctx,
		insertStatement,
		project,
		startDate,
		count,
	)

	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

func (p *PostgresJiraRepository) GetIssuesCompletedBetweenDates(ctx context.Context, project string, startDate time.Time, endDate time.Time, issueTypes []string, endStates []string) ([]string, error) {
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

func (p *PostgresJiraRepository) SetIssuesCompletedInWeekStarting(ctx context.Context, project string, startDate time.Time, count int) (int64, error) {
	insertStatement := `
		INSERT INTO jira_issues_reports(project, week_start, number_of_items_completed)
		VALUES ($1, $2, $3)
		ON CONFLICT (project, week_start)
		DO UPDATE
		SET number_of_items_completed = $3
		WHERE jira_issues_reports.project = $1 AND jira_issues_reports.week_start = $2
	`

	result, err := p.db.ExecContext(ctx,
		insertStatement,
		project,
		startDate,
		count,
	)

	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}

func (p *PostgresJiraRepository) GetEndStateForIssue(ctx context.Context, issueKey string, transitionDate time.Time) (string, error) {
	selectStatement := `
		SELECT jira_transitions.to_state from jira_transitions
		WHERE jira_transitions.issue_key = $1
		AND jira_transitions.created_at = $2
	`

	rows, err := p.db.QueryContext(ctx,
		selectStatement,
		issueKey,
		transitionDate,
	)

	if err != nil {
		return "", err
	}

	var result []string
	for rows.Next() {
		var issueKey string

		err = rows.Scan(&issueKey)
		if err != nil {
			return "", err
		}

		result = append(result, issueKey)
	}

	if len(result) != 1 {
		return "", fmt.Errorf("unexpected length of transitions: %d", len(result))
	}

	return result[0], nil
}

func (p *PostgresJiraRepository) SaveIssueLabels(ctx context.Context, issueKey string, labels []string) (int64, error) {
	insertStatement := `
		INSERT INTO jira_issue_labels(issue_key, label)
		VALUES ($1, $2)
		ON CONFLICT (issue_key, label)
		DO NOTHING
	`
	var inserted int64 = 0

	for _, label := range labels {
		result, err := p.db.ExecContext(ctx,
			insertStatement,
			issueKey,
			label,
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

func (p *PostgresJiraRepository) SaveJiraToDoStates(ctx context.Context, project string, states []string) (int64, error) {
	return p.saveJiraWorkStates(ctx, project, states, "todo")
}

func (p *PostgresJiraRepository) SaveJiraInProgressStates(ctx context.Context, project string, states []string) (int64, error) {
	return p.saveJiraWorkStates(ctx, project, states, "in progress")
}

func (p *PostgresJiraRepository) SaveJiraDoneStates(ctx context.Context, project string, states []string) (int64, error) {
	return p.saveJiraWorkStates(ctx, project, states, "done")
}

func (p *PostgresJiraRepository) saveJiraWorkStates(ctx context.Context, project string, states []string, workState string) (int64, error) {
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

func (p *PostgresJiraRepository) GetWeeklyThroughputByProject(ctx context.Context, project string, endDate time.Time, numberOfWeeks int) ([]statsmodels.WeeklyTimeItem, error) {
	selectStatement := fmt.Sprintf(`
		SELECT count(jira_issues_calculations.issue_key), 
			date_trunc('week', jira_issues_calculations.issue_completed_at) AS "completed_week"
		FROM jira_issues_calculations
		INNER JOIN jira_issues ON jira_issues_calculations.issue_key = jira_issues.issue_key
		WHERE jira_issues_calculations.cycle_time > 0
		AND jira_issues_calculations.issue_completed_at > date_trunc('week', $1::timestamp - interval '%d weeks') 
		AND jira_issues.issue_type IN ('Task', 'Story')
		AND lower(jira_issues.project) = lower($2)
		AND jira_issues_calculations.issue_end_state IN (
			SELECT state_name FROM jira_work_states WHERE lower(jira_work_states.project) = lower($2)
		)
		GROUP BY date_trunc('week', "public"."jira_issues_calculations"."issue_completed_at")
		ORDER BY date_trunc('week', "public"."jira_issues_calculations"."issue_completed_at") ASC
	`, numberOfWeeks)

	rows, err := p.db.QueryContext(ctx, selectStatement, endDate, project)

	if err != nil {
		return []statsmodels.WeeklyTimeItem{}, err
	}

	var result = []statsmodels.WeeklyTimeItem{}

	for rows.Next() {
		tp := statsmodels.WeeklyTimeItem{}

		err = rows.Scan(&tp.NumberOfItems, &tp.WeekStarting)
		if err != nil {
			return result, err
		}

		result = append(result, tp)
	}

	return result, nil
}

func (p *PostgresJiraRepository) GetWeeklyThroughputAllProjects(ctx context.Context, endDate time.Time, numberOfWeeks int) ([]statsmodels.ProjectWeeklyTimeItem, error) {
	selectStatement := fmt.Sprintf(`
		SELECT count(jira_issues_calculations.issue_key), 
			date_trunc('week', jira_issues_calculations.issue_completed_at) AS "completed_week",
			jira_issues.project
		FROM jira_issues_calculations
		INNER JOIN jira_issues ON jira_issues_calculations.issue_key = jira_issues.issue_key
		WHERE jira_issues_calculations.cycle_time > 0
		AND jira_issues_calculations.issue_completed_at > date_trunc('week', $1::timestamp - interval '%d weeks') 
		AND jira_issues.issue_type IN ('Task', 'Story')
		AND jira_issues_calculations.issue_end_state IN (
			SELECT state_name FROM jira_work_states WHERE jira_work_states.project = jira_issues.project
		)
		GROUP BY date_trunc('week', "public"."jira_issues_calculations"."issue_completed_at"), jira_issues.project
		ORDER BY date_trunc('week', "public"."jira_issues_calculations"."issue_completed_at") ASC
	`, numberOfWeeks)

	rows, err := p.db.QueryContext(ctx, selectStatement, endDate)

	if err != nil {
		return []statsmodels.ProjectWeeklyTimeItem{}, err
	}

	var result = []statsmodels.ProjectWeeklyTimeItem{}

	for rows.Next() {
		tp := statsmodels.ProjectWeeklyTimeItem{}

		err = rows.Scan(&tp.TimeItem.NumberOfItems, &tp.TimeItem.WeekStarting, &tp.Project)
		if err != nil {
			return result, err
		}

		result = append(result, tp)
	}

	return result, nil
}

func (p *PostgresJiraRepository) GetWeeklyCycleTimeByProject(ctx context.Context, project string, endDate time.Time, numberOfWeeks int) ([]statsmodels.WeeklyCycleTimeItem, error) {
	selectStatement := fmt.Sprintf(`
		SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY jira_issues_calculations.cycle_time) as cycle_time, 
			date_trunc('week', jira_issues_calculations.issue_completed_at) AS "completed_week"
		FROM jira_issues_calculations
		INNER JOIN jira_issues ON jira_issues_calculations.issue_key = jira_issues.issue_key
		WHERE jira_issues_calculations.cycle_time > 0
		AND jira_issues_calculations.issue_completed_at > date_trunc('week', $1::timestamp - interval '%d weeks') 
		AND jira_issues.issue_type IN ('Task', 'Story')
		AND lower(jira_issues.project) = lower($2)
		AND jira_issues_calculations.issue_end_state IN (
			SELECT state_name FROM jira_work_states WHERE jira_work_states.state_type = 'done' AND lower(jira_work_states.project) = lower($2)
		)
		GROUP BY date_trunc('week', jira_issues_calculations.issue_completed_at)
		ORDER BY date_trunc('week', jira_issues_calculations.issue_completed_at) ASC
	`, numberOfWeeks)

	rows, err := p.db.QueryContext(ctx, selectStatement, endDate, project)

	if err != nil {
		return []statsmodels.WeeklyCycleTimeItem{}, err
	}

	var result = []statsmodels.WeeklyCycleTimeItem{}

	for rows.Next() {
		ct := statsmodels.WeeklyCycleTimeItem{}

		err = rows.Scan(&ct.CycleTime, &ct.WeekStarting)
		if err != nil {
			return result, err
		}

		result = append(result, ct)
	}

	return result, nil
}

func (p *PostgresJiraRepository) GetWeeklyCycleTimeAllProjects(ctx context.Context, endDate time.Time, numberOfWeeks int) ([]statsmodels.ProjectWeeklyCycleTimeItem, error) {
	selectStatement := fmt.Sprintf(`
		SELECT percentile_cont(0.5) WITHIN GROUP (ORDER BY jira_issues_calculations.cycle_time) as cycle_time, 
			date_trunc('week', jira_issues_calculations.issue_completed_at) AS "completed_week",
			jira_issues.project
		FROM jira_issues_calculations
		INNER JOIN jira_issues ON jira_issues_calculations.issue_key = jira_issues.issue_key
		WHERE jira_issues_calculations.cycle_time > 0
		AND jira_issues_calculations.issue_completed_at > date_trunc('week', $1::timestamp - interval '%d weeks') 
		AND jira_issues.issue_type IN ('Task', 'Story')
		AND jira_issues_calculations.issue_end_state IN (
			SELECT state_name FROM jira_work_states WHERE jira_work_states.state_type = 'done' AND jira_work_states.project = jira_issues.project
		)
		GROUP BY date_trunc('week', jira_issues_calculations.issue_completed_at), jira_issues.project
		ORDER BY date_trunc('week', jira_issues_calculations.issue_completed_at) ASC
	`, numberOfWeeks)
	rows, err := p.db.QueryContext(ctx, selectStatement, endDate)

	if err != nil {
		return []statsmodels.ProjectWeeklyCycleTimeItem{}, err
	}

	var result = []statsmodels.ProjectWeeklyCycleTimeItem{}

	for rows.Next() {
		ct := statsmodels.ProjectWeeklyCycleTimeItem{}

		err = rows.Scan(&ct.TimeItem.CycleTime, &ct.TimeItem.WeekStarting, &ct.Project)
		if err != nil {
			return result, err
		}

		result = append(result, ct)
	}

	return result, nil
}

func (p *PostgresJiraRepository) GetProjects(ctx context.Context) ([]string, error) {
	selectStatement := `
		SELECT DISTINCT project from jira_issues
	`
	rows, err := p.db.QueryContext(ctx, selectStatement)
	if err != nil {
		return []string{}, err
	}

	result := []string{}
	for rows.Next() {
		var project string
		rows.Scan(&project)
		result = append(result, project)
	}

	return result, nil
}

func (p *PostgresJiraRepository) SaveNonWorkingDays(ctx context.Context, project string, nonWorkingDays []string) (int64, error) {
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

func (p *PostgresJiraRepository) GetNonWorkingDays(ctx context.Context, project string) ([]time.Time, error) {
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
