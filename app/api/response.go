package api

import (
	"encoding/json"
	"net/http"
)

// Error represents a standard JSON error response
type Error struct {
	Error string `json:"error"`
}

func OKResponse(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If encoding fails, fall back to a 500 error
		ErrorResponse(w, http.StatusInternalServerError, err.Error())
	}
}

func ErrorResponse(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	json.NewEncoder(w).Encode(Error{
		Error: message,
	})
}
