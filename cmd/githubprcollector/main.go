package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Netflix/go-env"
	"github.com/stebennett/squad-dashboard/pkg/githubservice"
)

type Environment struct {
	GithubUser         string `env:"GITHUB_USERNAME,required=true"`
	GithubAccessToken  string `env:"GITHUB_ACCESS_TOKEN,required=true"`
	GithubOrganisation string `env:"GITHUB_ORGANISATION,required=true"`
}

func main() {
	var environment Environment
	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}

	// create a new github service
	github := createGithubService(environment)

	repositoryNames, err := github.GetRepositoriesForOrganisation(environment.GithubOrganisation)
	if err != nil {
		log.Fatal(err)
	}

	for _, name := range repositoryNames {
		pullrequests, err := github.GetPullRequestsForRepo(environment.GithubOrganisation, name)
		if err != nil {
			log.Fatal(err)
		}

		for _, pr := range pullrequests {
			log.Printf("[%s-%d] %s", name, pr.Number, pr.Title)
		}
	}
}

func createGithubService(environment Environment) *githubservice.GithubService {
	githubParams := githubservice.GithubParams{
		User:                environment.GithubUser,
		PersonalAccessToken: environment.GithubAccessToken,
	}

	githubClient := http.Client{
		Timeout: time.Second * 30,
	}

	return githubservice.NewGithubService(&githubClient, githubParams)
}
