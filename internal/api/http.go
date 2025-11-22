package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"ciphermint-gaming-gateway/internal/models"
	"ciphermint-gaming-gateway/internal/sqlstore"
)

//
// ===== Shared types =====
//

// ErrorResponse is used for all error JSON.
type ErrorResponse struct {
	Error string `json:"error"`
}

// HealthResponse is returned from /health.
type HealthResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

// RegisterIntegrationRequest is the body for POST /v1/game.
type RegisterIntegrationRequest struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CompanyID string `json:"company_id"`
}

// CreatePlayerRequest is the body for POST /v1/game/{integration_id}/player.
type CreatePlayerRequest struct {
	PlayerID string `json:"player_id"`
	Alias    string `json:"alias"`
}

// EarnOrSpendRequest is the body for earning/spending a token.
type EarnOrSpendRequest struct {
	Token  string `json:"token"`  // e.g. "RACKDOG"
	Amount int64  `json:"amount"` // must be > 0
	Source string `json:"source"` // e.g. "signup_bonus", "match_win"
}

//
// ===== Handler + router =====
//

// Handler ties everything to the SQL store.
type Handler struct {
	store *sqlstore.Store
}

// NewRouter wires up all routes and returns the HTTP handler the gateway uses.
func NewRouter(store *sqlstore.Store) http.Handler {
	h := &Handler{store: store}
	r := mux.NewRouter()

	// Health check (no API key required)
	r.HandleFunc("/health", h.HealthHandler).Methods("GET")

	// Integration + player + balances (API key required)
	r.HandleFunc("/v1/game", h.RegisterIntegrationHandler).Methods("POST")
	r.HandleFunc("/v1/game/{integration_id}/player", h.RegisterPlayerHandler).Methods("POST")
	r.HandleFunc("/v1/game/{integration_id}/player/{player_id}/earn", h.EarnHandler).Methods("POST")
	r.HandleFunc("/v1/game/{integration_id}/player/{player_id}", h.GetPlayerHandler).Methods("GET")

	return r
}

//
// ===== Helpers =====
//

// writeJSON writes any object as a proper JSON response.
func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

// writeError sends a JSON error using ErrorResponse.
func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}

// checkAPIKey enforces the RACKDOG_API_KEY header when set in the env.
func checkAPIKey(w http.ResponseWriter, r *http.Request) bool {
	expected := os.Getenv("RACKDOG_API_KEY")
	if expected == "" {
		// dev mode -- no key configured, allow all
		return true
	}

	got := r.Header.Get("X-RACKDOG-API-KEY")
	if got == "" || got != expected {
		writeError(w, http.StatusUnauthorized, "invalid or missing API key")
		return false
	}
	return true
}

//
// ===== Handlers =====
//

// HealthHandler returns a simple OK payload.
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{
		Service: "CipherMint Gaming Gateway",
		Status:  "ok",
	}
	writeJSON(w, http.StatusOK, resp)
}

// RegisterIntegrationHandler handles POST /v1/game.
func (h *Handler) RegisterIntegrationHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAPIKey(w, r) {
		return
	}

	var req RegisterIntegrationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.ID == "" || req.Name == "" {
		writeError(w, http.StatusBadRequest, "id and name are required")
		return
	}

	integ := &models.Integration{
		ID:        req.ID,
		Name:      req.Name,
		CompanyID: req.CompanyID,
	}

	if err := h.store.RegisterIntegration(r.Context(), integ); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, integ)
}

// RegisterPlayerHandler handles POST /v1/game/{integration_id}/player.
func (h *Handler) RegisterPlayerHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAPIKey(w, r) {
		return
	}

	vars := mux.Vars(r)
	integrationID := vars["integration_id"]

	var req CreatePlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.PlayerID == "" || req.Alias == "" {
		writeError(w, http.StatusBadRequest, "player_id and alias are required")
		return
	}

	player := &models.Player{
		ID:            req.PlayerID,
		Alias:         req.Alias,
		IntegrationID: integrationID,
		Balances:      map[string]int64{},
	}

	if err := h.store.RegisterPlayer(r.Context(), player); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, player)
}

// EarnHandler handles POST /v1/game/{integration_id}/player/{player_id}/earn.
func (h *Handler) EarnHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAPIKey(w, r) {
		return
	}

	vars := mux.Vars(r)
	integrationID := vars["integration_id"]
	playerID := vars["player_id"]

	var req EarnOrSpendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Token == "" || req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "token and positive amount are required")
		return
	}

	if err := h.store.UpdateBalance(r.Context(), integrationID, playerID, req.Token, req.Amount); err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"token":  req.Token,
		"amount": req.Amount,
	})
}

// GetPlayerHandler handles GET /v1/game/{integration_id}/player/{player_id}.
func (h *Handler) GetPlayerHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAPIKey(w, r) {
		return
	}

	vars := mux.Vars(r)
	integrationID := vars["integration_id"]
	playerID := vars["player_id"]

	player, err := h.store.GetPlayer(r.Context(), integrationID, playerID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, player)
}