package routes

import (
	"github.com/gorilla/mux"
	"github.com/stebennett/squad-dashboard/pkg/api/controllers"
)

func JiraStatsRoutes(statsController controllers.StatsContoller, router *mux.Router) {
	router.HandleFunc("/stats/throughput", statsController.ThroughputAllProjects).Methods("GET")
	router.HandleFunc("/stats/{project}/throughput", statsController.ThroughputByProject).Methods("GET")
}