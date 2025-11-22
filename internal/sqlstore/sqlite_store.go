package sqlstore

import (
	"context"
	"database/sql"
	"fmt"

	"ciphermint-gaming-gateway/internal/models"
	_ "github.com/mattn/go-sqlite3"
)

// Store wraps the SQLite DB.
type Store struct {
	db *sql.DB
}

// OpenDefault opens (or creates) the SQLite DB file.
func OpenDefault() (*Store, error) {
	db, err := sql.Open("sqlite3", "./ciphermint_gateway.db?_fk=1")
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping sqlite: %w", err)
	}

	s := &Store{db: db}

	if err := s.migrate(); err != nil {
		db.Close()
		return nil, err
	}

	return s, nil
}

// Close safely closes the DB.
func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

// migrate ensures tables exist.
func (s *Store) migrate() error {
	schema := `
CREATE TABLE IF NOT EXISTS games (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    company_id TEXT
);

CREATE TABLE IF NOT EXISTS players (
    id TEXT PRIMARY KEY,
    alias TEXT,
    integration_id TEXT
);

CREATE TABLE IF NOT EXISTS balances (
    player_id TEXT,
    game_id TEXT,
    token TEXT,
    amount INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY(player_id, game_id, token)
);
`
	_, err := s.db.Exec(schema)
	return err
}

// RegisterGame inserts or updates a game.
func (s *Store) RegisterIntegration(ctx context.Context, game *models.Integration) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO games (id, name, company_id)
         VALUES (?, ?, ?)
         ON CONFLICT(id) DO UPDATE SET
            name=excluded.name,
            company_id=excluded.company_id`,
		game.ID, game.Name, game.CompanyID,
	)
	return err
}

// RegisterPlayer inserts a new player.
func (s *Store) RegisterPlayer(ctx context.Context, p *models.Player) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO players (id, alias, integration_id)
         VALUES (?, ?, ?)
         ON CONFLICT(id) DO NOTHING`,
		p.ID, p.Alias, p.IntegrationID,
	)
	return err
}

// AddBalance increases tokens by amount.
func (s *Store) AddBalance(ctx context.Context, gameID, playerID, token string, amt int64) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO balances (player_id, game_id, token, amount)
         VALUES (?, ?, ?, ?)
         ON CONFLICT(player_id, game_id, token)
         DO UPDATE SET amount = balances.amount + excluded.amount`,
		playerID, gameID, token, amt,
	)
	return err
}

// SpendTokens enforces non-negative balances.
func (s *Store) SpendTokens(ctx context.Context, gameID, playerID, token string, amount int64) error {
	if amount <= 0 {
		return fmt.Errorf("amount must be > 0")
	}

	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	defer func() { _ = tx.Rollback() }()

	var current int64
	err = tx.QueryRowContext(
		ctx,
		`SELECT amount FROM balances
         WHERE player_id = ? AND game_id = ? AND token = ?`,
		playerID, gameID, token,
	).Scan(&current)

	if err == sql.ErrNoRows {
		return fmt.Errorf("insufficient balance")
	}
	if err != nil {
		return fmt.Errorf("load balance: %w", err)
	}

	if current < amount {
		return fmt.Errorf("insufficient balance")
	}

	_, err = tx.ExecContext(
		ctx,
		`UPDATE balances
         SET amount = amount - ?
         WHERE player_id = ? AND game_id = ? AND token = ?`,
		amount, playerID, gameID, token,
	)
	if err != nil {
		return fmt.Errorf("update balance: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit: %w", err)
	}

	return nil
}

// LoadPlayer returns player + balances.
func (s *Store) LoadPlayer(ctx context.Context, gameID, playerID string) (*models.Player, error) {
	p := &models.Player{ID: playerID}

	err := s.db.QueryRowContext(
		ctx,
		`SELECT alias, integration_id FROM players WHERE id = ?`,
		playerID,
	).Scan(&p.Alias, &p.IntegrationID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("player not found")
		}
		return nil, err
	}

	rows, err := s.db.QueryContext(
		ctx,
		`SELECT token, amount FROM balances
         WHERE player_id = ? AND game_id = ?`,
		playerID, gameID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	p.Balances = make(map[string]int64)

	for rows.Next() {
		var token string
		var amt int64
		if err := rows.Scan(&token, &amt); err != nil {
			return nil, err
		}
		p.Balances[token] = amt
	}

	return p, nil
}
