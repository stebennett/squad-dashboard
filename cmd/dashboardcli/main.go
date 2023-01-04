package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/stebennett/squad-dashboard/pkg/dashboard"
	"github.com/stebennett/squad-dashboard/pkg/jirarepository"
	"github.com/stebennett/squad-dashboard/pkg/printer"
)

func main() {
	repo := createJiraRepository()
	printer := createPrinter()

	escapedDefects, err := dashboard.GenerateEscapedDefects(12, repo)
	if err != nil {
		log.Fatal(err)
	}

	printer.Print(escapedDefects)
}

func createJiraRepository() jirarepository.JiraRepository {
	var err error
	var db *sql.DB
	connStr := os.ExpandEnv("postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=disable") // load from env vars

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	fmt.Println("Database initialised")
	return jirarepository.NewPostgresJiraRepository(db)
}

func createPrinter() printer.Printer {
	return printer.NewCommandLinePrinter()
}
