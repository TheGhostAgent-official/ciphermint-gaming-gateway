package models

// Integration represents a single game/company integration into the CipherMint
// Gaming Gateway.
type Integration struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	CompanyID string `json:"company_id"`
}

// Player represents a player identity inside a specific integration.
// Balances are per-token (RACKDOG, etc).
type Player struct {
	ID            string           `json:"id"`
	Alias         string           `json:"alias"`
	IntegrationID string           `json:"integration_id"`
	Balances      map[string]int64 `json:"balances"` // token symbol -> balance
}