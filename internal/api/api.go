package api

import (
	"encoding/json"
	"net/http"
)

// WriteJSON writes any value as JSON with the given HTTP status.
func WriteJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// WriteError sends an ErrorResponse JSON with the given status.
func WriteError(w http.ResponseWriter, status int, msg string) {
	WriteJSON(w, status, ErrorResponse{Error: msg})
}
