package main

import (
	"log"
	"net/http"
	"time"

	"github.com/Netflix/go-env"
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

	since, _ := time.Parse(time.RFC3339, "2022-06-01T00:00:00Z")
	until, _ := time.Parse(time.RFC3339, "2022-06-30T23:59:59Z")

	oncalls, err := pagerduty.GetOnCalls(since, until)
	if err != nil {
		log.Fatalf("Failed to complete on-call requests - %s", err)
	}

	log.Printf("found %d on-calls", len(oncalls))

	for _, oncall := range oncalls {
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
