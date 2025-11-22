package player

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// errorResponse is a simple envelope for error messages.
type errorResponse struct {
	Error string `json:"error"`
}

// writeJSON is a small helper to send JSON responses consistently.
func writeJSON(w http.ResponseWriter, status int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(body)
}

// -----------------------
// CreatePlayerHandler
// -----------------------
// POST /v1/game/{gameID}/player
// This is the "account created" hook. Think:
// "Xbox Live profile connecting into a specific game"
// We just echo back a basic player object for now.
func CreatePlayerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]

	if gameID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "gameID path param is required",
		})
		return
	}

	type createReq struct {
		PlayerID string   `json:"player_id"`
		Alias    string   `json:"alias"`
		Aliases  []string `json:"aliases,omitempty"`
	}

	var body createReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "invalid request body",
		})
		return
	}

	if body.PlayerID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "player_id is required",
		})
		return
	}

	// For now we just build a Player struct in memory.
	p := Player{
		ID:       body.PlayerID,
		GameID:   gameID,
		Aliases:  append([]string{}, body.Alias),
		Balances: map[string]int{},
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "created",
		"player": p,
	})
}

// -----------------------
// GetPlayerHandler
// -----------------------
// GET /v1/game/{gameID}/player/{playerID}
// For now this simply returns a minimal stubbed player object,
// proving the route shape and handler wiring work.
func GetPlayerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	playerID := vars["playerID"]

	if gameID == "" || playerID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "gameID and playerID path params are required",
		})
		return
	}

	p := Player{
		ID:       playerID,
		GameID:   gameID,
		Aliases:  []string{},
		Balances: map[string]int{},
	}

	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status": "ok",
		"player": p,
	})
}

// -----------------------
// EarnTokenHandler
// -----------------------
// POST /v1/game/{gameID}/player/{playerID}/earn
// This is your "earn RackDawg when user does something in-game" hook.
func EarnTokenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	playerID := vars["playerID"]

	if gameID == "" || playerID == "" {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "gameID and playerID path params are required",
		})
		return
	}

	type earnReq struct {
		Token  string `json:"token"`
		Amount int    `json:"amount"`
		Source string `json:"source"`
	}

	var body earnReq
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "invalid request body",
		})
		return
	}

	if body.Token == "" || body.Amount <= 0 {
		writeJSON(w, http.StatusBadRequest, errorResponse{
			Error: "token and positive amount are required",
		})
		return
	}

	// Stub response – later this will hit storage and mint logic.
	writeJSON(w, http.StatusOK, map[string]interface{}{
		"status":    "ok",
		"game_id":   gameID,
		"player_id": playerID,
		"token":     body.Token,
		"amount":    body.Amount,
		"source":    body.Source,
	})
}
