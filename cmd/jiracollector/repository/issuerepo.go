package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/stebennett/squad-dashboard/cmd/jiracollector/models"
)

var db *sql.DB

func init() {
	var err error
	connStr := os.ExpandEnv("postgres://$DB_USERNAME:$DB_PASSWORD@DB_HOST:$DB_PORT/$DB_NAME") // load from env vars

	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	if err = db.Ping(); err != nil {
		panic(err)
	}

	fmt.Println("Database initialised")
}

func StoreIssue(jiraIssue models.JiraIssue) error {
	// TODO: Implement storing of issue
	return errors.New("not yet implemented")
}

func StoreTransition(jiraTransition models.JiraTransition) error {
	// TODO: Implement storage of a transition
	return errors.New("not yet implemented")
}
