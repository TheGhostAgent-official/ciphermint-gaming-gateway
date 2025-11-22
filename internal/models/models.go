package models

type PlayerID string
type TokenSymbol string

type Player struct {
    ID       PlayerID                `json:"id"`
    Aliases  []string                `json:"aliases,omitempty"`
    Balances map[TokenSymbol]int64   `json:"balances"`
}

type EarnRequest struct {
    Token    TokenSymbol          `json:"token"`
    Amount   int64                `json:"amount"`
    Source   string               `json:"source"`
    Metadata map[string]string    `json:"metadata,omitempty"`
}

type SpendRequest struct {
    Token    TokenSymbol          `json:"token"`
    Amount   int64                `json:"amount"`
    Reason   string               `json:"reason"`
    Metadata map[string]string    `json:"metadata,omitempty"`
}

type CreatePlayerRequest struct {
    PlayerID PlayerID  `json:"player_id"`
    Alias    string    `json:"alias,omitempty"`
}