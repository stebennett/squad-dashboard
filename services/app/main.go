package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/Netflix/go-env"
	"github.com/gin-gonic/gin"
	"github.com/stebennett/squad-dashboard/pkg/jira/repo/calculationsrepository"
	"github.com/stebennett/squad-dashboard/pkg/statsservice"
	"github.com/stebennett/squad-dashboard/services/app/routes"
)

type Environment struct {
	JiraDefectIssueType  string `env:"JIRA_DEFECT_ISSUE_TYPE,required=true"`
	JiraReportIssueTypes string `env:"JIRA_REPORT_ISSUE_TYPES,required=true"`
	CycleTimePercentile  string `env:"REPORT_CYCLE_TIME_PERCENTILE,required=true"`
}

func main() {
	var environment Environment
	_, err := env.UnmarshalFromEnviron(&environment)
	if err != nil {
		log.Fatal(err)
	}

	repo := createJiraRepository()
	ss := statsservice.NewStatsService(repo)
	api := routes.Api{
		StatsService:         *ss,
		CycleTimeIssueTypes:  strings.Split(environment.JiraReportIssueTypes, ","),
		ThroughputIssueTypes: strings.Split(environment.JiraReportIssueTypes, ","),
		DefectIssueType:      environment.JiraDefectIssueType,
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.GET("/api/v1/:project/workingcycletime", api.WorkingCycleTimeGET)
	r.GET("/api/v1/:project/throughput", api.ThroughputGET)
	r.GET("/api/v1/:project/escapeddefects", api.EscapedDefectsGET)
	r.GET("/api/v1/:project/unplannedwork", api.UnplannedWorkGET)

	r.Run() // listen and serve on 0.0.0.0:8080
}

func createJiraRepository() calculationsrepository.JiraCalculationsRepository {
	var err error
	var db *sql.DB
	connStr := os.ExpandEnv("postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable") // load from env vars

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database initialised")
	return calculationsrepository.NewPostgresJiraCalculationsRepository(db)
}
