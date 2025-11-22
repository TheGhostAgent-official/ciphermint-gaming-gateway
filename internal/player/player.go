package player

// Player represents a single gamer account inside one game's economy.
// Example: one Xbox Live account inside a specific title.
type Player struct {
	ID       string         `json:"player_id"`
	GameID   string         `json:"game_id"`
	Aliases  []string       `json:"aliases"`
	Balances map[string]int `json:"balances"`
}
