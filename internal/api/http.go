package api

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

// Handler is the main HTTP handler for the CipherMint Gaming Gateway.
// It holds in-memory state for integrations ("games") and their players.
type Handler struct {
	mu    sync.RWMutex
	games map[string]*Game
}

// Game represents a single integration (not a literal video game, but a
// tokenized integration layer for an existing title or studio).
type Game struct {
	ID        string             `json:"id"`
	Name      string             `json:"name"`
	CompanyID string             `json:"company_id"`
	Players   map[string]*Player `json:"-"`
}

// Player represents a player profile inside a specific integration.
type Player struct {
	ID       string           `json:"id"`
	GameID   string           `json:"game_id"`
	Aliases  []string         `json:"aliases"`
	Balances map[string]int64 `json:"balances"`
}

// NewHandler constructs a new Handler with empty in-memory state.
func NewHandler() *Handler {
	return &Handler{
		games: make(map[string]*Game),
	}
}

// NewRouter wires all HTTP routes for the gateway.
func NewRouter(h *Handler) http.Handler {
	r := mux.NewRouter()

	// Health + manifest
	r.HandleFunc("/health", h.Health).Methods("GET")
	r.HandleFunc("/v1/manifest", h.Manifest).Methods("GET")

	// Integration (game) registration
	r.HandleFunc("/v1/game", h.RegisterGame).Methods("POST")

	// Player + token routes
	r.HandleFunc("/v1/game/{gameID}/player", h.CreatePlayer).Methods("POST")
	r.HandleFunc("/v1/game/{gameID}/player/{playerID}", h.GetPlayer).Methods("GET")
	r.HandleFunc("/v1/game/{gameID}/player/{playerID}/earn", h.Earn).Methods("POST")
	r.HandleFunc("/v1/game/{gameID}/player/{playerID}/spend", h.Spend).Methods("POST")

	return r
}

// Health is a simple health check so studios can verify the gateway is alive.
func (h *Handler) Health(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{
		"service": "CipherMint Gaming Gateway",
		"status":  "ok",
	}
	writeJSON(w, http.StatusOK, resp)
}

// Manifest returns a machine-readable description of the core CipherMint
// Gaming Gateway API. Studios hit this first to understand what they can hook into.
func (h *Handler) Manifest(w http.ResponseWriter, r *http.Request) {
	type Route struct {
		Method      string `json:"method"`
		Path        string `json:"path"`
		Description string `json:"description"`
	}

	manifest := struct {
		Service  string  `json:"service"`
		Version  string  `json:"version"`
		Category string  `json:"category"`
		Routes   []Route `json:"routes"`
	}{
		Service:  "CipherMint Gaming Gateway",
		Version:  "v1",
		Category: "tokenized-game-integrations",
		Routes: []Route{
			{
				Method:      "GET",
				Path:        "/health",
				Description: "Basic health check for the gateway.",
			},
			{
				Method:      "GET",
				Path:        "/v1/manifest",
				Description: "Returns this API manifest: routes and purpose.",
			},
			{
				Method:      "POST",
				Path:        "/v1/game",
				Description: "Register an integration (e.g., 'Ghost Ops – CoD Integration').",
			},
			{
				Method:      "POST",
				Path:        "/v1/game/{gameID}/player",
				Description: "Create/link a player profile under an integration.",
			},
			{
				Method:      "GET",
				Path:        "/v1/game/{gameID}/player/{playerID}",
				Description: "Fetch a player's profile and token balances.",
			},
			{
				Method:      "POST",
				Path:        "/v1/game/{gameID}/player/{playerID}/earn",
				Description: "Credit tokens to a player (sign-in rewards, streaks, match wins, etc.).",
			},
			{
				Method:      "POST",
				Path:        "/v1/game/{gameID}/player/{playerID}/spend",
				Description: "Spend tokens for in-game items, skins, boosts, etc.",
			},
		},
	}

	writeJSON(w, http.StatusOK, manifest)
}

// ------- Core request/response models -------

type registerGameRequest struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CompanyID string `json:"company_id"`
}

type createPlayerRequest struct {
	PlayerID string `json:"player_id"`
	Alias    string `json:"alias"`
}

type earnRequest struct {
	Token  string `json:"token"`
	Amount int64  `json:"amount"`
	Source string `json:"source"`
}

type spendRequest struct {
	Token  string `json:"token"`
	Amount int64  `json:"amount"`
	Reason string `json:"reason"`
}

// ------- Handlers: integration registration & lookup -------

// RegisterGame creates a new integration shell for an existing title/studio.
func (h *Handler) RegisterGame(w http.ResponseWriter, r *http.Request) {
	var req registerGameRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if req.ID == "" {
		writeError(w, http.StatusBadRequest, "game id is required")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "game name is required")
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.games[req.ID]; exists {
		writeError(w, http.StatusConflict, "game already exists")
		return
	}

	game := &Game{
		ID:        req.ID,
		Name:      req.Name,
		CompanyID: req.CompanyID,
		Players:   make(map[string]*Player),
	}

	h.games[req.ID] = game

	writeJSON(w, http.StatusCreated, game)
}

// ------- Handlers: player lifecycle -------

// CreatePlayer creates or links a player profile under a given integration.
func (h *Handler) CreatePlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]

	if gameID == "" {
		writeError(w, http.StatusBadRequest, "game id is required")
		return
	}

	var req createPlayerRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.PlayerID == "" {
		writeError(w, http.StatusBadRequest, "player id is required")
		return
	}
	if req.Alias == "" {
		writeError(w, http.StatusBadRequest, "alias is required")
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	game, ok := h.games[gameID]
	if !ok {
		writeError(w, http.StatusNotFound, "game not found")
		return
	}

	player, ok := game.Players[req.PlayerID]
	if !ok {
		player = &Player{
			ID:       req.PlayerID,
			GameID:   gameID,
			Aliases:  []string{},
			Balances: make(map[string]int64),
		}
		game.Players[req.PlayerID] = player
	}

	// Add alias if it's new
	aliasExists := false
	for _, a := range player.Aliases {
		if a == req.Alias {
			aliasExists = true
			break
		}
	}
	if !aliasExists {
		player.Aliases = append(player.Aliases, req.Alias)
	}

	writeJSON(w, http.StatusCreated, player)
}

// GetPlayer returns the current state of a player's profile.
func (h *Handler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	playerID := vars["playerID"]

	if gameID == "" || playerID == "" {
		writeError(w, http.StatusBadRequest, "game id and player id are required")
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	game, ok := h.games[gameID]
	if !ok {
		writeError(w, http.StatusNotFound, "game not found")
		return
	}

	player, ok := game.Players[playerID]
	if !ok {
		writeError(w, http.StatusNotFound, "player not found")
		return
	}

	writeJSON(w, http.StatusOK, player)
}

// ------- Handlers: earning & spending tokens -------

// Earn credits tokens to a player inside a specific integration.
func (h *Handler) Earn(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	playerID := vars["playerID"]

	if gameID == "" || playerID == "" {
		writeError(w, http.StatusBadRequest, "game id and player id are required")
		return
	}

	var req earnRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Token == "" {
		writeError(w, http.StatusBadRequest, "token is required")
		return
	}
	if req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "amount must be positive")
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	game, ok := h.games[gameID]
	if !ok {
		writeError(w, http.StatusNotFound, "game not found")
		return
	}

	player, ok := game.Players[playerID]
	if !ok {
		writeError(w, http.StatusNotFound, "player not found")
		return
	}

	if player.Balances == nil {
		player.Balances = make(map[string]int64)
	}

	player.Balances[req.Token] += req.Amount

	writeJSON(w, http.StatusOK, player)
}

// Spend debits tokens from a player when they buy in-game value.
func (h *Handler) Spend(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	gameID := vars["gameID"]
	playerID := vars["playerID"]

	if gameID == "" || playerID == "" {
		writeError(w, http.StatusBadRequest, "game id and player id are required")
		return
	}

	var req spendRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Token == "" {
		writeError(w, http.StatusBadRequest, "token is required")
		return
	}
	if req.Amount <= 0 {
		writeError(w, http.StatusBadRequest, "amount must be positive")
		return
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	game, ok := h.games[gameID]
	if !ok {
		writeError(w, http.StatusNotFound, "game not found")
		return
	}

	player, ok := game.Players[playerID]
	if !ok {
		writeError(w, http.StatusNotFound, "player not found")
		return
	}

	current := player.Balances[req.Token]
	if current < req.Amount {
		writeError(w, http.StatusBadRequest, "insufficient balance")
		return
	}

	player.Balances[req.Token] = current - req.Amount

	writeJSON(w, http.StatusOK, player)
}

// ------- small helpers -------

func writeJSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}