package api

import (
    "ciphermint-gaming-gateway/internal/models"
    "ciphermint-gaming-gateway/internal/store"
    "encoding/json"
    "net/http"

    "github.com/gorilla/mux"
)

// Handler holds the store implementation so our HTTP handlers
// can read/write games, players, and balances.
type Handler struct {
    store store.Store
}

// NewRouter builds the HTTP router and wires all routes to methods on Handler.
func NewRouter(s store.Store) *mux.Router {
    h := &Handler{store: s}
    r := mux.NewRouter()

    // Basic health check
    r.HandleFunc("/health", h.Health).Methods("GET")

    // Game-level routes
    r.HandleFunc("/v1/game", h.RegisterGame).Methods("POST")

    // Player + wallet routes (all scoped by gameID)
    r.HandleFunc("/v1/game/{gameID}/player", h.CreatePlayer).Methods("POST")
    r.HandleFunc("/v1/game/{gameID}/player/{playerID}", h.GetPlayer).Methods("GET")
    r.HandleFunc("/v1/game/{gameID}/player/{playerID}/earn", h.Earn).Methods("POST")
    r.HandleFunc("/v1/game/{gameID}/player/{playerID}/spend", h.Spend).Methods("POST")

    return r
}

// Health is a simple readiness endpoint.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
    writeJSON(w, http.StatusOK, map[string]string{
        "status":  "ok",
        "service": "CipherMint Gaming Gateway",
    })
}

// RegisterGame registers a new game in the store.
// Body JSON should look like:
// {
//   "id": "game-1",
//   "name": "Example Game",
//   "company_id": "company-123"
// }
func (h *Handler) RegisterGame(w http.ResponseWriter, r *http.Request) {
    var game models.Game
    if err := json.NewDecoder(r.Body).Decode(&game); err != nil {
        writeError(w, http.StatusBadRequest, "invalid JSON body")
        return
    }

    saved, err := h.store.RegisterGame(&game)
    if err != nil {
        writeError(w, http.StatusBadRequest, err.Error())
        return
    }

    writeJSON(w, http.StatusCreated, saved)
}

// CreatePlayer registers a player under a specific game.
// URL:   /v1/game/{gameID}/player
// Body:  { "player_id": "player123", "alias": "GhostPlayer" }
func (h *Handler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    gameID := models.GameID(vars["gameID"])

    var req models.CreatePlayerRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid JSON body")
        return
    }

    // Path gameID is the source of truth.
    player, err := h.store.CreatePlayer(gameID, req.PlayerID, req.Alias)
    if err != nil {
        if err == store.ErrGameNotFound {
            writeError(w, http.StatusNotFound, "game not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "could not create player")
        return
    }

    writeJSON(w, http.StatusCreated, player)
}

// GetPlayer fetches a player's balances for a specific game.
// URL: /v1/game/{gameID}/player/{playerID}
func (h *Handler) GetPlayer(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    gameID := models.GameID(vars["gameID"])
    playerID := models.PlayerID(vars["playerID"])

    player, err := h.store.GetPlayer(gameID, playerID)
    if err != nil {
        if err == store.ErrGameNotFound {
            writeError(w, http.StatusNotFound, "game not found")
            return
        }
        if err == store.ErrPlayerNotFound {
            writeError(w, http.StatusNotFound, "player not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "could not fetch player")
        return
    }

    writeJSON(w, http.StatusOK, player)
}

// Earn credits tokens to a player's in-game wallet.
// URL:  /v1/game/{gameID}/player/{playerID}/earn
// Body: { "token": "STAKE", "amount": 100, "source": "match_win" }
func (h *Handler) Earn(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    gameID := models.GameID(vars["gameID"])
    playerID := models.PlayerID(vars["playerID"])

    var req models.EarnRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid JSON body")
        return
    }

    player, err := h.store.Earn(gameID, playerID, req.Token, req.Amount)
    if err != nil {
        if err == store.ErrGameNotFound {
            writeError(w, http.StatusNotFound, "game not found")
            return
        }
        if err == store.ErrPlayerNotFound {
            writeError(w, http.StatusNotFound, "player not found")
            return
        }
        writeError(w, http.StatusInternalServerError, "could not apply earn")
        return
    }

    writeJSON(w, http.StatusOK, player)
}

// Spend debits tokens from a player's in-game wallet.
// URL:  /v1/game/{gameID}/player/{playerID}/spend
// Body: { "token": "STAKE", "amount": 40, "reason": "skin_purchase" }
func (h *Handler) Spend(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    gameID := models.GameID(vars["gameID"])
    playerID := models.PlayerID(vars["playerID"])

    var req models.SpendRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, http.StatusBadRequest, "invalid JSON body")
        return
    }

    player, err := h.store.Spend(gameID, playerID, req.Token, req.Amount)
    if err != nil {
        if err == store.ErrGameNotFound {
            writeError(w, http.StatusNotFound, "game not found")
            return
        }
        if err == store.ErrPlayerNotFound {
            writeError(w, http.StatusNotFound, "player not found")
            return
        }
        if err == store.ErrInsufficientFunds {
            writeError(w, http.StatusBadRequest, "insufficient funds")
            return
        }
        writeError(w, http.StatusInternalServerError, "could not apply spend")
        return
    }

    writeJSON(w, http.StatusOK, player)
}

// Helper to write JSON responses consistently.
func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    _ = json.NewEncoder(w).Encode(v)
}

// Helper to write error responses consistently.
func writeError(w http.ResponseWriter, status int, message string) {
    writeJSON(w, status, map[string]string{
        "error": message,
    })
}