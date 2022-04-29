package githubservice

import (
	"errors"
	"net/http"

	"github.com/stebennett/squad-dashboard/pkg/githubmodels"
)

type GithubParams struct {
	User                string
	PersonalAccessToken string
}

type GithubService struct {
	githubParams GithubParams
	httpClient   *http.Client
}

func NewGithubService(httpClient *http.Client, githubParams GithubParams) *GithubService {
	return &GithubService{
		githubParams: githubParams,
		httpClient:   httpClient,
	}
}

func (g *GithubService) GetPullRequestsForRepo(organisation string, repository string) ([]githubmodels.GithubPullRequest, error) {
	return []githubmodels.GithubPullRequest{}, errors.New("not yet implemented")
}

func (g *GithubService) GetRepositoriesForOrganistion(organisation string) ([]string, error) {
	return []string{}, errors.New("not yet implemented")
}
