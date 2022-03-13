package repository

import (
	"context"
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
	"github.com/stebennett/squad-dashboard/cmd/jiracollector/models"
)

type IssueRepository interface {
	StoreIssue(ctx context.Context, jiraIssue models.JiraIssue) error
	StoreTransition(ctx context.Context, jiraTransition models.JiraTransition) error
}

type PostgresIssueRepository struct {
	db *sql.DB
}

func NewPostgresIssueRepository(db *sql.DB) *PostgresIssueRepository {
	return &PostgresIssueRepository{
		db: db,
	}
}

func (p *PostgresIssueRepository) StoreIssue(ctx context.Context, jiraIssue models.JiraIssue) error {
	// TODO: Implement storing of issue
	return errors.New("not yet implemented")
}

func (p *PostgresIssueRepository) StoreTransition(ctx context.Context, jiraTransition models.JiraTransition) error {
	// TODO: Implement storage of a transition
	return errors.New("not yet implemented")
}
