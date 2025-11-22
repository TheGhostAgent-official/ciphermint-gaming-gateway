package api

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	"ciphermint-gaming-gateway/internal/models"
	"ciphermint-gaming-gateway/internal/sqlstore"
)

// Basic response types
type ErrorResponse struct {
	Error string `json:"error"`
}

type HealthResponse struct {
	Service string `json:"service"`
	Status  string `json:"status"`
}

// Request payloads
type RegisterGameRequest struct {
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

// NewRouter wires up all routes and returns an http.Handler.
func NewRouter(store *sqlstore.Store) http.Handler {
	h := &Handler{store: store}
	r := mux.NewRouter()

	r.HandleFunc("/health", h.HealthHandler).Methods("GET")

	// Game/integration endpoints
	r.HandleFunc("/v1/game", h.RegisterGameHandler).Methods("POST")
	r.HandleFunc("/v1/game/{gameID}/player", h.CreatePlayerHandler).Methods("POST")
	r.HandleFunc("/v1/game/{gameID}/player/{playerID}", h.GetPlayerHandler).Methods("GET")
	r.HandleFunc("/v1/game/{gameID}/player/{playerID}/earn", h.EarnTokenHandler).Methods("POST")

	return r
}

// --- helpers ---

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, ErrorResponse{Error: msg})
}

func withCtx(r *http.Request) context.Context {
	return r.Context()
}

// --- handlers ---

// HealthHandler returns a simple health response.
func (h *Handler) HealthHandler(w http.ResponseWriter, r *http.Request) {
	resp := HealthResponse{
		Service: "CipherMint Gaming Gateway",
		Status:  "ok",
	}
	writeJSON(w, http.StatusOK, resp)
}

// RegisterGameHandler registers or updates a game/integration.
func (h *Handler) RegisterGameHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if req.ID == "" || req.Name == "" {
		writeError(w, http.StatusBadRequest, "game id and name are required")
		return
	}

	game := &models.Integration{
		ID:        req.ID,
		Name:      req.Name,
		CompanyID: req.CompanyID,
	}

	if err := h.store.RegisterIntegration(withCtx(r), game); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"id":     game.ID,
	})
}

// CreatePlayerHandler attaches/creates a player for a game.
func (h *Handler) CreatePlayerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	if gameID == "" {
		writeError(w, http.StatusBadRequest, "game id is required")
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

	p := &models.Player{
		ID:            req.PlayerID,
		Alias:         req.Alias,
		IntegrationID: gameID,
		Balances:      make(map[string]int64),
	}

	if err := h.store.RegisterPlayer(withCtx(r), p); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":      "ok",
		"player_id":   p.ID,
		"integration": p.IntegrationID,
		"alias":       p.Alias,
		"balances":    p.Balances,
	})
}

// GetPlayerHandler returns a player with balances.
func (h *Handler) GetPlayerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	playerID := vars["playerID"]
	if gameID == "" || playerID == "" {
		writeError(w, http.StatusBadRequest, "game id and player id are required")
		return
	}

	p, err := h.store.LoadPlayer(withCtx(r), gameID, playerID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if p == nil {
		writeError(w, http.StatusNotFound, "player not found")
		return
	}

	writeJSON(w, http.StatusOK, p)
}

// EarnTokenHandler adds token balance for a player.
func (h *Handler) EarnTokenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	playerID := vars["playerID"]
	if gameID == "" || playerID == "" {
		writeError(w, http.StatusBadRequest, "game id and player id are required")
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

	if err := h.store.AddBalance(withCtx(r), gameID, playerID, req.Token, req.Amount); err != nil {
		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"token":  req.Token,
		"amount": req.Amount,
		"source": req.Source,
	})
}
