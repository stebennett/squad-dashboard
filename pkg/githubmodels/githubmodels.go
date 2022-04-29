package githubmodels

import "time"

type GithubPullRequest struct {
	User      string
	Title     string
	Id        int64
	Number    int
	CreatedAt time.Time
	UpdatedAt time.Time
	ClosedAt  time.Time
	MergedAt  time.Time
}
