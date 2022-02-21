package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/stebennett/squad-dashboard/cmd/jiracollector/models"
)

type DBIssueRepository struct {
	dbPool *pgxpool.Pool
}

func NewDBIssueRepository(dbPool *pgxpool.Pool) *DBIssueRepository {
	return &DBIssueRepository{
		dbPool: dbPool,
	}
}

func (p DBIssueRepository) SaveIssue(ctx context.Context, jiraIssue models.JiraIssue) (*models.JiraIssue, error) {
	// TODO: Implement storing of issue
	return nil, errors.New("not yet implemented")
}

func (p DBIssueRepository) SaveTransition(ctx context.Context, jiraTransition models.JiraTransition) (*models.JiraTransition, error) {
	// TODO: Implement storage of a transition
	return nil, errors.New("not yet implemented")
}
