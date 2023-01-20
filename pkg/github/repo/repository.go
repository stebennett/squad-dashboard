package repo

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"
	"github.com/stebennett/squad-dashboard/pkg/github/models"
)

type GithubRespository interface {
	SavePullRequest(ctx context.Context, organisation string, repository string, pullRequest models.GithubPullRequest) (int64, error)
}

type PosgresGithubRepository struct {
	db *sql.DB
}

func NewPostgresGithubRepository(db *sql.DB) *PosgresGithubRepository {
	return &PosgresGithubRepository{
		db: db,
	}
}

func (p *PosgresGithubRepository) SavePullRequest(ctx context.Context, organisation string, repository string, pullRequest models.GithubPullRequest) (int64, error) {
	insertStatement := `
		INSERT INTO github_pull_requests(organisation, repository, gh_user, title, github_id, pr_number, created_at, updated_at, closed_at, merged_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		ON CONFLICT (organisation, repository, pr_number)
		DO UPDATE
		SET gh_user = $3, title = $4, github_id = $5, created_at = $7, updated_at = $8, closed_at = $9, merged_at = $10
		WHERE github_pull_requests.organisation = $1 AND github_pull_requests.repository = $2 AND github_pull_requests.pr_number = $6
	`

	result, err := p.db.ExecContext(ctx,
		insertStatement,
		organisation,
		repository,
		pullRequest.User,
		pullRequest.Title,
		pullRequest.Id,
		pullRequest.Number,
		pullRequest.CreatedAt,
		pullRequest.UpdatedAt,
		pullRequest.ClosedAt,
		pullRequest.MergedAt,
	)

	if err != nil {
		return -1, err
	}

	return result.RowsAffected()
}
