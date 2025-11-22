package store

import (
    "errors"
    "sync"

    "ciphermint-gaming-gateway/internal/models"
)

// Domain errors we can reuse in the API layer.
var (
    ErrGameNotFound        = errors.New("game not found")
    ErrPlayerNotFound      = errors.New("player not found")
    ErrInsufficientFunds   = errors.New("insufficient funds")
)

// Store is the interface our API layer will depend on.
// This keeps things flexible if we later move from in-memory to a real database.
type Store interface {
    RegisterGame(game *models.Game) (*models.Game, error)
    GetGame(id models.GameID) (*models.Game, error)

    CreatePlayer(gameID models.GameID, playerID models.PlayerID, alias string) (*models.Player, error)
    GetPlayer(gameID models.GameID, playerID models.PlayerID) (*models.Player, error)

    Earn(gameID models.GameID, playerID models.PlayerID, token models.TokenSymbol, amount int64) (*models.Player, error)
    Spend(gameID models.GameID, playerID models.PlayerID, token models.TokenSymbol, amount int64) (*models.Player, error)
}

// MemoryStore is our in-memory implementation of Store.
// Internally, it organizes data like:
// games[gameID] = *Game
// players[gameID][playerID] = *Player
type MemoryStore struct {
    mu      sync.RWMutex
    games   map[models.GameID]*models.Game
    players map[models.GameID]map[models.PlayerID]*models.Player
}

// NewMemoryStore creates a new, empty in-memory store.
func NewMemoryStore() *MemoryStore {
    return &MemoryStore{
        games:   make(map[models.GameID]*models.Game),
        players: make(map[models.GameID]map[models.PlayerID]*models.Player),
    }
}

// RegisterGame registers or updates a game in the store.
func (s *MemoryStore) RegisterGame(game *models.Game) (*models.Game, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    if game.ID == "" {
        return nil, errors.New("game id is required")
    }
    if game.Name == "" {
        return nil, errors.New("game name is required")
    }

    s.games[game.ID] = game

    // Ensure players map exists for this game.
    if _, ok := s.players[game.ID]; !ok {
        s.players[game.ID] = make(map[models.PlayerID]*models.Player)
    }

    return game, nil
}

// GetGame retrieves a game by its ID.
func (s *MemoryStore) GetGame(id models.GameID) (*models.Game, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    g, ok := s.games[id]
    if !ok {
        return nil, ErrGameNotFound
    }
    return g, nil
}

// CreatePlayer registers a new player under a specific game.
func (s *MemoryStore) CreatePlayer(gameID models.GameID, playerID models.PlayerID, alias string) (*models.Player, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    // Ensure the game exists first.
    if _, ok := s.games[gameID]; !ok {
        return nil, ErrGameNotFound
    }

    // Ensure we have a players map for this game.
    if _, ok := s.players[gameID]; !ok {
        s.players[gameID] = make(map[models.PlayerID]*models.Player)
    }

    // Create the player object.
    p := &models.Player{
        ID:       playerID,
        GameID:   gameID,
        Aliases:  []string{},
        Balances: make(map[models.TokenSymbol]int64),
    }

    if alias != "" {
        p.Aliases = append(p.Aliases, alias)
    }

    s.players[gameID][playerID] = p
    return p, nil
}

// GetPlayer fetches a player under a specific game.
func (s *MemoryStore) GetPlayer(gameID models.GameID, playerID models.PlayerID) (*models.Player, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    gamePlayers, ok := s.players[gameID]
    if !ok {
        return nil, ErrGameNotFound
    }

    p, ok := gamePlayers[playerID]
    if !ok {
        return nil, ErrPlayerNotFound
    }

    return p, nil
}

// Earn credits tokens to a player's balance for a given game.
func (s *MemoryStore) Earn(gameID models.GameID, playerID models.PlayerID, token models.TokenSymbol, amount int64) (*models.Player, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    gamePlayers, ok := s.players[gameID]
    if !ok {
        return nil, ErrGameNotFound
    }

    p, ok := gamePlayers[playerID]
    if !ok {
        return nil, ErrPlayerNotFound
    }

    if p.Balances == nil {
        p.Balances = make(map[models.TokenSymbol]int64)
    }

    p.Balances[token] += amount
    return p, nil
}

// Spend debits tokens from a player's balance for a given game.
func (s *MemoryStore) Spend(gameID models.GameID, playerID models.PlayerID, token models.TokenSymbol, amount int64) (*models.Player, error) {
    s.mu.Lock()
    defer s.mu.Unlock()

    gamePlayers, ok := s.players[gameID]
    if !ok {
        return nil, ErrGameNotFound
    }

    p, ok := gamePlayers[playerID]
    if !ok {
        return nil, ErrPlayerNotFound
    }

    current := p.Balances[token]
    if current < amount {
        return nil, ErrInsufficientFunds
    }

    p.Balances[token] = current - amount
    return p, nil
}