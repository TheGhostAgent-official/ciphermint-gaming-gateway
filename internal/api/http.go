package api

import (
	"encoding/json"
	"log"
	"net/http"

	"ciphermint-gaming-gateway/internal/models"
	"ciphermint-gaming-gateway/internal/store"
	"github.com/gorilla/mux"
)

// Handler wires HTTP handlers to the backing store.
type Handler struct {
	store store.Store
}

// NewHandler builds a Handler from a Store.
func NewHandler(s store.Store) *Handler {
	return &Handler{store: s}
}

// NewRouter builds the HTTP router with all routes and middlewares.
func NewRouter(h *Handler) *mux.Router {
	r := mux.NewRouter()

	// Simple health check for studios / infra
	r.HandleFunc("/health", h.Health).Methods(http.MethodGet)

	// Versioned API
	v1 := r.PathPrefix("/v1").Subrouter()
	v1.HandleFunc("/player", h.CreatePlayer).Methods(http.MethodPost)
	v1.HandleFunc("/player/{id}", h.GetPlayer).Methods(http.MethodGet)
	v1.HandleFunc("/player/{id}/earn", h.EarnTokens).Methods(http.MethodPost)
	v1.HandleFunc("/player/{id}/spend", h.SpendTokens).Methods(http.MethodPost)

	return r
}

// --- Handlers ---

func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{
		"status":  "ok",
		"service": "ciphermint-gaming-gateway",
	})
}

func (h *Handler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	var req models.CreatePlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.PlayerID == "" {
		writeError(w, http.StatusBadRequest, "player_id is required")
		return
	}

	player, err := h.store.CreatePlayer(req.PlayerID, req.Alias)
	if err != nil {
		log.Printf("CreatePlayer error: %v", err)
		writeError(w, http.StatusInternalServerError, "unable to create player")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]*models.Player{
		"player": player,
	})
}

func (h *Handler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing player id in path")
		return
	}

	player, err := h.store.GetPlayer(models.PlayerID(id))
	if err != nil {
		if err == store.ErrPlayerNotFound {
			writeError(w, http.StatusNotFound, "player not found")
			return
		}
		log.Printf("GetPlayer error: %v", err)
		writeError(w, http.StatusInternalServerError, "unable to fetch player")
		return
	}

	writeJSON(w, http.StatusOK, map[string]*models.Player{
		"player": player,
	})
}

func (h *Handler) EarnTokens(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing player id in path")
		return
	}

	var req models.EarnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Token == "" || req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "token and positive amount are required")
		return
	}

	player, err := h.store.Earn(models.PlayerID(id), req.Token, req.Amount, req.Source)
	if err != nil {
		if err == store.ErrPlayerNotFound {
			writeError(w, http.StatusNotFound, "player not found")
			return
		}
		log.Printf("EarnTokens error: %v", err)
		writeError(w, http.StatusInternalServerError, "unable to apply earn event")
		return
	}

	writeJSON(w, http.StatusOK, map[string]*models.Player{
		"player": player,
	})
}

func (h *Handler) SpendTokens(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]
	if id == "" {
		writeError(w, http.StatusBadRequest, "missing player id in path")
		return
	}

	var req models.SpendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Token == "" || req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "token and positive amount are required")
		return
	}

	player, err := h.store.Spend(models.PlayerID(id), req.Token, req.Amount, req.Reason)
	if err != nil {
		if err == store.ErrPlayerNotFound {
			writeError(w, http.StatusNotFound, "player not found")
			return
		}
		if err == store.ErrInsufficientFunds {
			writeError(w, http.StatusConflict, "insufficient funds")
			return
		}
		log.Printf("SpendTokens error: %v", err)
		writeError(w, http.StatusInternalServerError, "unable to apply spend event")
		return
	}

	writeJSON(w, http.StatusOK, map[string]*models.Player{
		"player": player,
	})
}

// --- helpers ---

func writeJSON(w http.ResponseWriter, status int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("writeJSON error: %v", err)
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{
		"error": msg,
	})
}