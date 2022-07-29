package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	mux "github.com/gorilla/mux"
	"github.com/stebennett/squad-dashboard/pkg/api/controllers"
	"github.com/stebennett/squad-dashboard/pkg/api/routes"
	"github.com/stebennett/squad-dashboard/pkg/jirarepository"
	"github.com/stebennett/squad-dashboard/pkg/jirastatsservice"
)

func main() {
	jiraRepo := createJiraRepository()

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api/v1").Subrouter()

	routes.HealthRoutes(apiRouter)

	statsController := controllers.StatsContoller{
		StatsService: jirastatsservice.JiraStatsService{
			JiraRepository: jiraRepo,
		},
	}
	jiraDataController := controllers.JiraDataController{
		JiraRepository: jiraRepo,
	}

	routes.JiraStatsRoutes(statsController, apiRouter)
	routes.JiraDataRoutes(jiraDataController, apiRouter)

	srv := &http.Server{
		Handler:      router,
		Addr:         os.ExpandEnv(":$PORT"),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  30 * time.Second,
	}

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	log.Printf("Server Started")

	<-done
	log.Print("Server Stopped")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		// add extra close handling here
		cancel()
	}()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed: %+v", err)
	}

	log.Print("Server Shutdown")
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
