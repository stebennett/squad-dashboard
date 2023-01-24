package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Netflix/go-env"
	"github.com/stebennett/squad-dashboard/pkg/jira/repo/calculationsrepository"
	"github.com/stebennett/squad-dashboard/pkg/printer"
	"github.com/stebennett/squad-dashboard/pkg/statsservice"
)

type Environment struct {
	JiraProject          string `env:"JIRA_PROJECT,required=true"`
	JiraDefectIssueType  string `env:"JIRA_DEFECT_ISSUE_TYPE,required=true"`
	JiraReportIssueTypes string `env:"JIRA_REPORT_ISSUE_TYPES,required=true"`
	OutputDirectory      string `env:"OUTPUT_DIRECTORY,required=true"`
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
	plotprinter := printer.NewPlotPrinter(environment.OutputDirectory, environment.JiraProject)
	cliprinter := printer.NewCommandLinePrinter()
	pdfprinter := printer.NewPdfReportPrinter(
		plotprinter.GetCycleTimeChartLocation(),
		plotprinter.GetThroughputChartLocation(),
		plotprinter.GetEscapedDefectsChartLocation(),
		plotprinter.GetUnplannedWorkChartLocation(),
		environment.JiraProject,
	)

	percentile, err := strconv.ParseFloat(environment.CycleTimePercentile, 64)
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()

	escapedDefects, err := ss.GenerateEscapedDefects(12, environment.JiraProject, environment.JiraDefectIssueType, now, time.Friday)
	if err != nil {
		log.Fatal(err)
	}

	cycleTimeReports, err := ss.GenerateCycleTime(12, percentile, environment.JiraProject, strings.Split(environment.JiraReportIssueTypes, ","), now, time.Friday)
	if err != nil {
		log.Fatal(err)
	}

	throughputReports, err := ss.GenerateThroughput(12, environment.JiraProject, strings.Split(environment.JiraReportIssueTypes, ","), now, time.Friday)
	if err != nil {
		log.Fatal(err)
	}

	unplannedWorkReports, err := ss.GenerateUnplannedWorkReport(12, environment.JiraProject, strings.Split(environment.JiraReportIssueTypes, ","), now, time.Friday)
	if err != nil {
		log.Fatal(err)
	}

	reports := printer.Reports{
		EscapedDefects:       escapedDefects,
		CycleTimeReports:     cycleTimeReports,
		ThroughputReports:    throughputReports,
		UnplannedWorkReports: unplannedWorkReports,
	}
	plotprinter.Print(reports)
	cliprinter.Print(reports)
	pdfprinter.Print(reports)
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
