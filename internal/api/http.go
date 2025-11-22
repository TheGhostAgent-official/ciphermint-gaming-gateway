package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
)

// Game represents a game integration (e.g. "Ghost Ops — CoD Integration").
type Game struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CompanyID string `json:"company_id"`
}

// Player represents a single gamer account inside one game's economy.
// Example: one Xbox Live account inside a specific title.
type Player struct {
	ID       string         `json:"player_id"`
	GameID   string         `json:"game_id"`
	Aliases  []string       `json:"aliases"`
	Balances map[string]int `json:"balances"`
}

// RegisterRoutes wires all HTTP routes for the CipherMint Gaming Gateway.
func RegisterRoutes(r *mux.Router) {
	r.HandleFunc("/health", HealthHandler).Methods(http.MethodGet)

	r.HandleFunc("/v1/game", CreateGameHandler).Methods(http.MethodPost)

	r.HandleFunc("/v1/game/{gameID}/player", CreatePlayerHandler).Methods(http.MethodPost)

	r.HandleFunc("/v1/game/{gameID}/player/{playerID}", GetPlayerHandler).Methods(http.MethodGet)

	r.HandleFunc("/v1/game/{gameID}/player/{playerID}/earn", EarnTokenHandler).Methods(http.MethodPost)
}

// HealthHandler returns a simple service health indicator.
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	writeJSON(w, http.StatusOK, healthResponse{
		Service: "CipherMint Gaming Gateway",
		Status:  "ok",
	})
}

// CreateGameHandler registers a new game integration.
func CreateGameHandler(w http.ResponseWriter, r *http.Request) {
	type createGameRequest struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		CompanyID string `json:"company_id"`
	}

	var body createGameRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.ID == "" || body.Name == "" {
		writeError(w, http.StatusBadRequest, "game id and name are required")
		return
	}

	game := Game{
		ID:        body.ID,
		Name:      body.Name,
		CompanyID: body.CompanyID,
	}

	writeJSON(w, http.StatusCreated, game)
}

// CreatePlayerHandler creates a player profile under a given game.
// This is where your "Xbox account signs into Ghost Ops and gets tokens"
// flow will eventually plug into real storage.
func CreatePlayerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]

	if gameID == "" {
		writeError(w, http.StatusBadRequest, "gameID path param is required")
		return
	}

	type createPlayerRequest struct {
		PlayerID string   `json:"player_id"`
		Aliases  []string `json:"aliases"`
	}

	var body createPlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.PlayerID == "" {
		writeError(w, http.StatusBadRequest, "player_id is required")
		return
	}

	p := Player{
		ID:       body.PlayerID,
		GameID:   gameID,
		Aliases:  body.Aliases,
		Balances: map[string]int{},
	}

	writeJSON(w, http.StatusCreated, p)
}

// GetPlayerHandler returns a minimal player snapshot.
// Phase 3: this will query storage; for now it just echoes a basic struct.
func GetPlayerHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	playerID := vars["playerID"]

	if gameID == "" || playerID == "" {
		writeError(w, http.StatusBadRequest, "gameID and playerID path params are required")
		return
	}

	p := Player{
		ID:       playerID,
		GameID:   gameID,
		Aliases:  []string{},
		Balances: map[string]int{},
	}

	writeJSON(w, http.StatusOK, p)
}

// EarnTokenHandler credits a player with tokens for some in-game action
// (login bonus, streak, match win, etc.).
func EarnTokenHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	playerID := vars["playerID"]

	if gameID == "" || playerID == "" {
		writeError(w, http.StatusBadRequest, "gameID and playerID path params are required")
		return
	}

	type earnTokenRequest struct {
		Token  string `json:"token"`
		Amount int    `json:"amount"`
		Source string `json:"source"`
	}

	var body earnTokenRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		writeError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if body.Token == "" || body.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "token and positive amount are required")
		return
	}

	response := map[string]interface{}{
		"status":    "ok",
		"game_id":   gameID,
		"player_id": playerID,
		"token":     body.Token,
		"amount":    body.Amount,
		"source":    body.Source,
	}

	writeJSON(w, http.StatusOK, response)
}
