package routes

import (
	"github.com/gorilla/mux"
	"github.com/stebennett/squad-dashboard/pkg/api/controllers"
)

func HealthRoutes(router *mux.Router) {
	router.HandleFunc("/health", controllers.HealthHandler).Methods("GET")
}
