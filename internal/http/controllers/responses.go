package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/arman300s/uni-portal/internal/core/contracts"
)

type ErrorResponse struct {
	Error   string                     `json:"error"`
	Details contracts.ValidationErrors `json:"details,omitempty"`
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if payload == nil {
		return
	}
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, message string, details contracts.ValidationErrors) {
	writeJSON(w, status, ErrorResponse{Error: message, Details: details})
}
