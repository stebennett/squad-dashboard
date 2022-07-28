package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

func ThroughputByProject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"project": vars["project"]})
}
