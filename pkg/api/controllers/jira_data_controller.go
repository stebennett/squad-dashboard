package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stebennett/squad-dashboard/pkg/jirarepository"
)

type JiraDataController struct {
	JiraRepository jirarepository.JiraRepository
}

func (jdc JiraDataController) ListProjects(w http.ResponseWriter, r *http.Request) {
	results, err := jdc.JiraRepository.GetProjects(context.Background())
	if err != nil {
		fmt.Println(fmt.Errorf("failed to retrieve project, %w", err))
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(results)
}
