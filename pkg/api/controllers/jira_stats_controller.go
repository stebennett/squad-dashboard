package controllers

import (
	"encoding/json"
	"fmt"
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
	results, err := s.StatsService.FetchThrougputDataForProject(project)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to get throughput, %w", err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	// return to api
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(results)
}

func (s StatsContoller) ThroughputAllProjects(w http.ResponseWriter, r *http.Request) {
	// fetch throughput data for project - last 12 weeks based on previous friday
	results, err := s.StatsService.FetchThrougputDataForAllProjects()

	if err != nil {
		fmt.Println(fmt.Errorf("failed to get throughput for all projects, %w", err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	// return to api
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(results)
}

func (s StatsContoller) CycleTimeByProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	// get project parameter
	project := vars["project"]

	// fetch throughput data for project - last 12 weeks based on previous friday
	results, err := s.StatsService.FetchCycleTimeDataForProject(project)
	if err != nil {
		fmt.Println(fmt.Errorf("failed to get cycletime, %w", err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	// return to api
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(results)
}
