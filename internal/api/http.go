package api

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	"ciphermint-gaming-gateway/internal/models"
	"ciphermint-gaming-gateway/internal/sqlstore"
)

// Basic response envelopes

type ErrorResponse struct {
	Error string `json:"error"`
}

type HealthResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

// Request payloads

type RegisterIntegrationRequest struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CompanyID string `json:"company_id"`
}

type CreatePlayerRequest struct {
	PlayerID string `json:"player_id"`
	Alias    string `json:"alias"`
}

type EarnOrSpendRequest struct {
	Token  string `json:"token"`
	Amount int64  `json:"amount"`
	Source string `json:"source"`
}

// Handler ties everything to the SQL store.
type Handler struct {
	store *sqlstore.Store
}

// NewRouter wires up all routes and returns the HTTP handler.
func NewRouter(store *sqlstore.Store) http.Handler {
	h := &Handler{store: store}
	r := mux.NewRouter()

	// Health
	r.HandleFunc("/health", h.HealthHandler).Methods("GET")

	// API group: /v1/game/...
	api := r.PathPrefix("/v1/game").Subrouter()
	api.HandleFunc("", h.RegisterIntegrationHandler).Methods("POST")
	api.HandleFunc("/", h.RegisterIntegrationHandler).Methods("POST")

	api.HandleFunc("/{integration_id}/player", h.RegisterPlayerHandler).Methods("POST")
	api.HandleFunc("/{integration_id}/player/{player_id}/earn", h.EarnHandler).Methods("POST")
	api.HandleFunc("/{integration_id}/player/{player_id}", h.GetPlayerHandler).Methods("GET")

	return r
}

// Helpers

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}

// checkAPIKey enforces the RACKDOG_API_KEY header when configured.
func checkAPIKey(w http.ResponseWriter, r *http.Request) bool {
	expected := os.Getenv("RACKDOG_API_KEY")
	if expected == "" {
		// Dev mode: no key configured, allow all
		return true
	}

	got := r.Header.Get("X-RACKDOG-API-KEY")
	if got == "" || got != expected {
		writeError(w, http.StatusUnauthorized, "invalid or missing API key")
		return false
	}
	return true
}

// Handlers

func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{
		Service: "CipherMint Gaming Gateway",
		Status:  "ok",
	}
	writeJSON(w, http.StatusOK, resp)
}

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

func (h *Handler) RegisterPlayerHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAPIKey(w, r) {
		return
	}

	vars := mux.Vars(r)
	integrationID := vars["integration_id"]
	if integrationID == "" {
		writeError(w, http.StatusBadRequest, "integration_id is required")
		return
	}

	var req CreatePlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.PlayerID == "" {
		writeError(w, http.StatusBadRequest, "player_id is required")
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

func (h *Handler) EarnHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAPIKey(w, r) {
		return
	}

	vars := mux.Vars(r)
	integrationID := vars["integration_id"]
	playerID := vars["player_id"]

	if integrationID == "" || playerID == "" {
		writeError(w, http.StatusBadRequest, "integration_id and player_id are required")
		return
	}

	var req EarnOrSpendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if req.Token == "" || req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "token and positive amount are required")
		return
	}

	if err := h.store.UpdateBalance(
		r.Context(),
		integrationID,
		playerID,
		req.Token,
		req.Amount,
	); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"token":  req.Token,
	})
}

func (h *Handler) GetPlayerHandler(w http.ResponseWriter, r *http.Request) {
	if !checkAPIKey(w, r) {
		return
	}

	vars := mux.Vars(r)
	integrationID := vars["integration_id"]
	playerID := vars["player_id"]

	if integrationID == "" || playerID == "" {
		writeError(w, http.StatusBadRequest, "integration_id and player_id are required")
		return
	}

	player, err := h.store.GetPlayer(r.Context(), integrationID, playerID)
	if err != nil {
		writeError(w, http.StatusNotFound, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, player)
}