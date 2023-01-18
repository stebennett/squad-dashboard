package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Netflix/go-env"
	"github.com/stebennett/squad-dashboard/pkg/dashboard"
	"github.com/stebennett/squad-dashboard/pkg/jiracalculationsrepository"
	"github.com/stebennett/squad-dashboard/pkg/printer"
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

	escapedDefects, err := dashboard.GenerateEscapedDefects(12, environment.JiraProject, environment.JiraDefectIssueType, repo)
	if err != nil {
		log.Fatal(err)
	}

	cycleTimeReports, err := dashboard.GenerateCycleTime(12, percentile, environment.JiraProject, strings.Split(environment.JiraReportIssueTypes, ","), repo)
	if err != nil {
		log.Fatal(err)
	}

	throughputReports, err := dashboard.GenerateThroughput(12, environment.JiraProject, strings.Split(environment.JiraReportIssueTypes, ","), repo)
	if err != nil {
		log.Fatal(err)
	}

	unplannedWorkReports, err := dashboard.GenerateUnplannedWorkReport(12, environment.JiraProject, strings.Split(environment.JiraReportIssueTypes, ","), repo)
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

func createJiraRepository() jiracalculationsrepository.JiraCalculationsRepository {
	var err error
	var db *sql.DB
	connStr := os.ExpandEnv("postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable") // load from env vars

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database initialised")
	return jiracalculationsrepository.NewPostgresJiraCalculationsRepository(db)
}
