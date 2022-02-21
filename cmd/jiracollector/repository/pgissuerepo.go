package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/stebennett/squad-dashboard/cmd/jiracollector/models"
)

type PostgresIssueRepository struct {
	db *sql.DB
}

func NewPostgresIssueRepository(db *sql.DB) *PostgresIssueRepository {
	return &PostgresIssueRepository{
		db: db,
	}
}

func (p PostgresIssueRepository) SaveIssue(ctx context.Context, jiraIssue models.JiraIssue) (*models.JiraIssue, error) {
	// TODO: Implement storing of issue
	return nil, errors.New("not yet implemented")
}

func (p PostgresIssueRepository) StoreTransition(ctx context.Context, jiraTransition models.JiraTransition) (*models.JiraTransition, error) {
	// TODO: Implement storage of a transition
	return nil, errors.New("not yet implemented")
}
