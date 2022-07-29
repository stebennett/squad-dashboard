package routes

import (
	"github.com/gorilla/mux"
	"github.com/stebennett/squad-dashboard/pkg/api/controllers"
)

func JiraDataRoutes(dataController controllers.JiraDataController, router *mux.Router) {
	router.HandleFunc("/jira/projects", dataController.ListProjects).Methods("GET")
}
