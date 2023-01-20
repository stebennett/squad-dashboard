package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/Netflix/go-env"
	"github.com/stebennett/squad-dashboard/pkg/github/repo"
	"github.com/stebennett/squad-dashboard/pkg/github/service"
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

	// create a github db storage layer
	githubrepository := createGithubRepository()

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
			_, err = githubrepository.SavePullRequest(context.Background(), environment.GithubOrganisation, name, pr)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("[%s-%d] %s", name, pr.Number, pr.Title)
		}
	}
}

func createGithubService(environment Environment) *service.GithubService {
	githubParams := service.GithubParams{
		User:                environment.GithubUser,
		PersonalAccessToken: environment.GithubAccessToken,
	}

	githubClient := http.Client{
		Timeout: time.Second * 30,
	}

	return service.NewGithubService(&githubClient, githubParams)
}

func createGithubRepository() repo.GithubRespository {
	var err error
	var db *sql.DB
	connStr := os.ExpandEnv("postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable") // load from env vars

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database initialised")
	return repo.NewPostgresGithubRepository(db)
}
