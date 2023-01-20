package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/stebennett/squad-dashboard/pkg/github/models"
	"github.com/tomnomnom/linkheader"
)

type GithubParams struct {
	User                string
	PersonalAccessToken string
}

type GithubService struct {
	githubParams GithubParams
	httpClient   *http.Client
}

type GithubTimestamp struct {
	time.Time
}

type GithubPullRequestResponse struct {
	Id     int64 `json:"id"`
	Number int   `json:"number"`
	User   struct {
		Login string `json:"login"`
	}
	Title     string          `json:"title"`
	CreatedAt GithubTimestamp `json:"created_at"`
	UpdatedAt GithubTimestamp `json:"updated_at"`
	ClosedAt  GithubTimestamp `json:"closed_at"`
	MergedAt  GithubTimestamp `json:"merged_at"`
}

type GithubReposRequestResponse struct {
	Id   int64  `json:"id"`
	Name string `json:"name"`
}

func NewGithubService(httpClient *http.Client, githubParams GithubParams) *GithubService {
	return &GithubService{
		githubParams: githubParams,
		httpClient:   httpClient,
	}
}

func (g *GithubService) GetPullRequestsForRepo(organisation string, repository string) ([]models.GithubPullRequest, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/pulls", organisation, repository)

	result, err := g.getPullRequestsForRepo(url, []models.GithubPullRequest{})

	return result, err
}

func (g *GithubService) getPullRequestsForRepo(url string, result []models.GithubPullRequest) ([]models.GithubPullRequest, error) {
	log.Printf("fetching pull requests for repo from %s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return result, err
	}

	req.SetBasicAuth(g.githubParams.User, g.githubParams.PersonalAccessToken)
	req.Header.Set("Content-Type", "application/vnd.github.v3.raw+json")

	query := req.URL.Query()
	query.Add("state", "all")
	query.Add("per_page", "100")
	req.URL.RawQuery = query.Encode()

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return result, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return result, err
	}

	var pullRequestsResponse []GithubPullRequestResponse
	err = json.Unmarshal(body, &pullRequestsResponse)
	if err != nil {
		return []models.GithubPullRequest{}, err
	}

	for _, pr := range pullRequestsResponse {
		githubPr := models.GithubPullRequest{
			User:      pr.User.Login,
			Title:     pr.Title,
			Id:        pr.Id,
			Number:    pr.Number,
			CreatedAt: pr.CreatedAt.Time,
			UpdatedAt: pr.UpdatedAt.Time,
			ClosedAt:  pr.ClosedAt.Time,
			MergedAt:  pr.MergedAt.Time,
		}
		result = append(result, githubPr)
	}

	linkHeader := resp.Header["Link"]
	if len(linkHeader) != 1 {
		return result, nil
	}

	links := linkheader.Parse(linkHeader[0])
	for _, link := range links {
		if link.Rel == "next" {
			result, err = g.getPullRequestsForRepo(link.URL, result)
			if err != nil {
				return result, err
			}
		}
	}

	return result, nil
}

func (g *GithubService) GetRepositoriesForOrganisation(organisation string) ([]string, error) {
	url := fmt.Sprintf("https://api.github.com/orgs/%s/repos", organisation)

	result, err := g.getRepositoriesForOrganisation(url, []string{})
	return result, err
}

func (g *GithubService) getRepositoriesForOrganisation(url string, result []string) ([]string, error) {
	log.Printf("fetching repos for org from %s", url)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return []string{}, err
	}

	req.SetBasicAuth(g.githubParams.User, g.githubParams.PersonalAccessToken)
	req.Header.Set("Content-Type", "application/vnd.github.v3+json")

	query := req.URL.Query()
	query.Add("sort", "full_name")
	query.Add("per_page", "100")
	req.URL.RawQuery = query.Encode()

	resp, err := g.httpClient.Do(req)
	if err != nil {
		return []string{}, err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []string{}, err
	}

	var reposResponse []GithubReposRequestResponse
	err = json.Unmarshal(body, &reposResponse)
	if err != nil {
		return []string{}, err
	}

	for _, repo := range reposResponse {
		result = append(result, repo.Name)
	}

	linkHeader := resp.Header["Link"]
	if len(linkHeader) != 1 {
		return result, nil
	}

	links := linkheader.Parse(linkHeader[0])
	for _, link := range links {
		if link.Rel == "next" {
			result, err = g.getRepositoriesForOrganisation(link.URL, result)
			if err != nil {
				return result, err
			}
		}
	}

	return result, nil
}

func (p *GithubTimestamp) UnmarshalJSON(bytes []byte) error {
	var raw string
	err := json.Unmarshal(bytes, &raw)

	if err != nil {
		fmt.Printf("Failed to unmarshal timestamp - %s", err)
		return err
	}

	if len(raw) == 0 {
		return err
	}

	p.Time, err = time.Parse("2006-01-02T15:04:05Z", raw)
	return err
}
