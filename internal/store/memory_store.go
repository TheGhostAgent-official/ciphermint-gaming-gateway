package store

import (
	"errors"
	"sync"

	"ciphermint-gaming-gateway/internal/models"
)

var (
	// ErrPlayerNotFound is returned when a player ID is unknown.
	ErrPlayerNotFound = errors.New("player not found")

	// ErrInsufficientFunds is returned when a spend would drop below zero.
	ErrInsufficientFunds = errors.New("insufficient funds")
)

// Store defines the behaviors the API layer needs from a backing store.
type Store interface {
	CreatePlayer(id models.PlayerID, alias string) (*models.Player, error)
	GetPlayer(id models.PlayerID) (*models.Player, error)
	Earn(id models.PlayerID, token models.TokenSymbol, amount int64, source string) (*models.Player, error)
	Spend(id models.PlayerID, token models.TokenSymbol, amount int64, reason string) (*models.Player, error)
}

// MemoryStore is a simple in-memory implementation (perfect for dev / demos).
type MemoryStore struct {
	mu      sync.RWMutex
	players map[models.PlayerID]*models.Player
}

// NewMemoryStore constructs a fresh empty store.
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		players: make(map[models.PlayerID]*models.Player),
	}
}

func (s *MemoryStore) CreatePlayer(id models.PlayerID, alias string) (*models.Player, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// If the player already exists, just merge alias and return.
	if existing, ok := s.players[id]; ok {
		if alias != "" {
			found := false
			for _, a := range existing.Aliases {
				if a == alias {
					found = true
					break
				}
			}
			if !found {
				existing.Aliases = append(existing.Aliases, alias)
			}
		}
		return clonePlayer(existing), nil
	}

	player := &models.Player{
		ID:       id,
		Aliases:  nil,
		Balances: make(map[models.TokenSymbol]int64),
	}
	if alias != "" {
		player.Aliases = []string{alias}
	}

	s.players[id] = player
	return clonePlayer(player), nil
}

func (s *MemoryStore) GetPlayer(id models.PlayerID) (*models.Player, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	player, ok := s.players[id]
	if !ok {
		return nil, ErrPlayerNotFound
	}
	return clonePlayer(player), nil
}

func (s *MemoryStore) Earn(id models.PlayerID, token models.TokenSymbol, amount int64, source string) (*models.Player, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	player, ok := s.players[id]
	if !ok {
		return nil, ErrPlayerNotFound
	}

	if player.Balances == nil {
		player.Balances = make(map[models.TokenSymbol]int64)
	}
	player.Balances[token] += amount

	return clonePlayer(player), nil
}

func (s *MemoryStore) Spend(id models.PlayerID, token models.TokenSymbol, amount int64, reason string) (*models.Player, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	player, ok := s.players[id]
	if !ok {
		return nil, ErrPlayerNotFound
	}

	if player.Balances == nil || player.Balances[token] < amount {
		return nil, ErrInsufficientFunds
	}

	player.Balances[token] -= amount
	return clonePlayer(player), nil
}

// clonePlayer returns a deep copy so callers can’t accidentally mutate store state.
func clonePlayer(p *models.Player) *models.Player {
	if p == nil {
		return nil
	}

	cp := &models.Player{
		ID:      p.ID,
		Aliases: append([]string(nil), p.Aliases...),
	}

	cp.Balances = make(map[models.TokenSymbol]int64, len(p.Balances))
	for k, v := range p.Balances {
		cp.Balances[k] = v
	}

	return cp
}