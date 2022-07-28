package routes

import (
	"github.com/gorilla/mux"
	"github.com/stebennett/squad-dashboard/pkg/api/controllers"
)

func JiraStatsRoutes(router *mux.Router) {
	router.HandleFunc("/stats/{project}/throughput", controllers.ThroughputByProject).Methods("GET")
}
