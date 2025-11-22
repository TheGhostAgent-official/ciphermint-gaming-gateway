package models

// Integration represents a single game/company integration into the CipherMint
// Gaming Gateway. This is NOT a game itself; it's the bridge between a
// provider (Activision, EA, 2K, etc.) and our token logic.
type Integration struct {
	ID         string `json:"id"`          // e.g. "ghostops_cod"
	Name       string `json:"name"`        // e.g. "Ghost Ops -- CoD Integration"
	CompanyID  string `json:"company_id"`  // optional, can be empty for now
}

// Player represents a player identity inside a specific integration.
// Balances are per-token (RACKDAWG, STAKE, etc).
type Player struct {
	ID            string           `json:"id"`
	Alias         string           `json:"alias"`
	IntegrationID string           `json:"integration_id"`
	Balances      map[string]int64 `json:"balances"` // token symbol -> balance
}

// CreatePlayerRequest is the payload for creating/attaching a player
// to an integration.
type CreatePlayerRequest struct {
	PlayerID string `json:"player_id"`
	Alias    string `json:"alias"`
}

// EarnOrSpendRequest is used for both earning and spending tokens.
type EarnOrSpendRequest struct {
	Token  string `json:"token"`  // e.g. "RACKDAWG"
	Amount int64  `json:"amount"` // must be > 0
	Source string `json:"source"` // e.g. "login", "match_win", "skin_purchase"
}

// IntegrationRequest is what a game company sends to register their
// integration shell with us.
type IntegrationRequest struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	CompanyID    string `json:"company_id"`
	Provider     string `json:"provider,omitempty"`      // optional, for future use
	GameTitle    string `json:"game_title,omitempty"`    // optional, for clarity
	Integration  string `json:"integration_name,omitempty"` // optional alias
}