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
	"github.com/stebennett/squad-dashboard/pkg/pagerdutyrepository"
	"github.com/stebennett/squad-dashboard/pkg/pagerdutyservice"
)

type Environment struct {
	PagerDutyAuthToken string `env:"PAGERDUTY_AUTH_TOKEN,required=true"`
}

func main() {
	var environment Environment
	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}

	// create a pagerduty service
	pagerduty := createPagerDutyService(environment)

	// create a pagerduty service
	pagerdutyrepository := createPagerDutyRepository()

	since, _ := time.Parse(time.RFC3339, "2022-06-01T00:00:00Z")
	until, _ := time.Parse(time.RFC3339, "2022-06-30T23:59:59Z")

	oncalls, err := pagerduty.GetOnCalls(since, until)
	if err != nil {
		log.Fatalf("Failed to complete on-call requests - %s", err)
	}

	log.Printf("found %d on-calls", len(oncalls))

	for _, oncall := range oncalls {
		_, err = pagerdutyrepository.SaveOnCall(context.Background(), oncall)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("on-call found - start[%s], end[%s], person[%s], escalation[%s], schedule[%s], level[%d]",
			oncall.Start.Format(time.RFC3339),
			oncall.End.Format(time.RFC3339),
			oncall.User.Name,
			oncall.EscalationPolicy.Name,
			oncall.Schedule.Name,
			oncall.EscalationLevel,
		)
	}
}

func createPagerDutyService(environment Environment) *pagerdutyservice.PagerDutyService {
	params := pagerdutyservice.PagerDutyParams{
		AuthToken: environment.PagerDutyAuthToken,
		BaseUrl:   "api.pagerduty.com",
	}

	client := http.Client{
		Timeout: time.Second * 30,
	}

	return pagerdutyservice.NewPagerDutyService(&client, params)
}

func createPagerDutyRepository() pagerdutyrepository.PagerDutyRepository {
	var err error
	var db *sql.DB
	connStr := os.ExpandEnv("postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable") // load from env vars

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database initialised")
	return pagerdutyrepository.NewPostgresPagerDutyRepository(db)
}
