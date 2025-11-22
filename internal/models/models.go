package models

// CompanyID represents a studio, publisher, or platform that integrates with CipherMint.
type CompanyID string

// GameID represents a specific game or digital experience owned by a company.
type GameID string

// PlayerID represents a unique player within a specific game.
type PlayerID string

// TokenSymbol represents the symbol/ticker of a game's token (e.g., STAKE, XP, etc.).
type TokenSymbol string

// Company is a high-level record for a partner integrating with CipherMint.
// (We won't fully use this yet, but it's important for the long-term vision.)
type Company struct {
	ID   CompanyID `json:"id"`
	Name string    `json:"name"`
}

// Game represents a single game that plugs into the CipherMint Gaming Gateway.
type Game struct {
	ID        GameID    `json:"id"`
	Name      string    `json:"name"`
	CompanyID CompanyID `json:"company_id"`
}

// Player represents a player inside a specific game, with on-chain style balances.
type Player struct {
	ID       PlayerID               `json:"id"`
	GameID   GameID                 `json:"game_id"`             // which game this player belongs to
	Aliases  []string               `json:"aliases,omitempty"`   // e.g., in-game usernames
	Balances map[TokenSymbol]int64  `json:"balances"`            // token → amount for this game/player
}

// EarnRequest represents a request to credit/earn tokens for a player.
type EarnRequest struct {
	Token    TokenSymbol          `json:"token"`
	Amount   int64                `json:"amount"`
	Source   string               `json:"source"`               // e.g. "match_win", "quest", "achievement"
	Metadata map[string]string    `json:"metadata,omitempty"`   // optional extra info (mode, map, etc.)
}

// SpendRequest represents a request to spend tokens for a player.
type SpendRequest struct {
	Token    TokenSymbol          `json:"token"`
	Amount   int64                `json:"amount"`
	Reason   string               `json:"reason"`               // e.g. "skin_purchase", "battle_pass"
	Metadata map[string]string    `json:"metadata,omitempty"`   // optional extra info
}

// CreatePlayerRequest is used when a new player is registered for a specific game.
type CreatePlayerRequest struct {
	GameID   GameID   `json:"game_id"`            // which game this player is for
	PlayerID PlayerID `json:"player_id"`          // unique ID from the game side
	Alias    string   `json:"alias,omitempty"`    // optional display name / gamertag
}