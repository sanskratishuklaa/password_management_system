package handler

import (
	"encoding/json"
	"net/http"
)


type HealthResponse struct {
	Status string `json:"status"`
}

func Health(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status: "ok",
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to create Response", http.StatusInternalServerError)
		return

	}
}
