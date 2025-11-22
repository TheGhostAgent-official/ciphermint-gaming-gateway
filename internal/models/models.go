package models

// Integration represents a single game/company integration into the CipherMint
// Gaming Gateway. This is NOT a game itself; it's the bridge between a
// provider (Activision, EA, 2K, etc.) and our token logic.
type Integration struct {
	ID        string `json:"id"`         // e.g. "ghostops_cod"
	Name      string `json:"name"`       // e.g. "Ghost Ops -- CoD Integration"
	CompanyID string `json:"company_id"` // optional, can be empty for now
}

// Player represents a player identity inside a specific integration.
// Balances are per-token (RACKDOG, STAKE, etc).
type Player struct {
	ID            string           `json:"id"`
	Alias         string           `json:"alias"`
	IntegrationID string           `json:"integration_id"`
	Balances      map[string]int64 `json:"balances"` // token symbol -> balance
}
