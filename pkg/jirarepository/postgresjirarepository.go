package jirarepository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	"github.com/stebennett/squad-dashboard/pkg/calculatormodels"
	jiramodels "github.com/stebennett/squad-dashboard/pkg/jira/models"
)

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

func (p *PostgresJiraRepository) ClearUnplannedIssuesForProject(ctx context.Context, project string) (int64, error) {
	updateStatement := `
		UPDATE jira_issues SET unplanned = FALSE WHERE project = $1
	`

	result, err := p.db.ExecContext(ctx, updateStatement, project)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}

func (p *PostgresJiraRepository) SaveUnplannedIssue(ctx context.Context, issueKey string) (int64, error) {
	updateStatement := `
		UPDATE jira_issues SET unplanned = TRUE WHERE issue_key = $1
	`

	result, err := p.db.ExecContext(ctx, updateStatement, issueKey)
	if err != nil {
		return -1, err
	}
	return result.RowsAffected()
}
