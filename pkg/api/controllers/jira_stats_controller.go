package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/stebennett/squad-dashboard/pkg/jirastatsservice"
)

type StatsContoller struct {
	StatsService jirastatsservice.JiraStatsService
}

func (s StatsContoller) ThroughputByProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// get project parameter
	project := vars["project"]

	// fetch throughput data for project - last 12 weeks based on previous friday
	results := s.StatsService.FetchThrougputDataForProject(project)

	// return to api
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(results)
}
