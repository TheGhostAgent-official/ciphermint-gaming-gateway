package models

// PlayerID is a stable identifier coming from the game / platform.
type PlayerID string

// TokenSymbol is the ticker / symbol for the token (e.g. RACKDOG, STAKE).
type TokenSymbol string

// Player represents a unique player in a game or app.
type Player struct {
	ID       PlayerID              `json:"id"`
	Aliases  []string              `json:"aliases,omitempty"`
	Balances map[TokenSymbol]int64 `json:"balances"`
}

// EarnRequest represents a request to reward tokens to a player.
type EarnRequest struct {
	Token    TokenSymbol       `json:"token"`
	Amount   int64             `json:"amount"`
	Source   string            `json:"source"`             // e.g. "quest", "win", "achievement"
	Metadata map[string]string `json:"metadata,omitempty"` // optional, for studios to extend
}

// SpendRequest represents a request to spend tokens (for items, skins, etc.).
type SpendRequest struct {
	Token    TokenSymbol       `json:"token"`
	Amount   int64             `json:"amount"`
	Reason   string            `json:"reason"`             // e.g. "skin_purchase"
	Metadata map[string]string `json:"metadata,omitempty"` // optional
}

// CreatePlayerRequest is used to register a new player.
type CreatePlayerRequest struct {
	PlayerID PlayerID `json:"player_id"`
	Alias    string   `json:"alias,omitempty"`
}
