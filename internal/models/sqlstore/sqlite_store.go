package sqlstore

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"

	"ciphermint-gaming-gateway/internal/models"
)

// Store wraps the SQLite DB.
type Store struct {
	db *sql.DB
}

// OpenDefault opens (or creates) the SQLite DB file and runs migrations.
func OpenDefault() (*Store, error) {
	db, err := sql.Open("sqlite3", "./ciphermint_gateway.db?_fk=1")
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}

	if err := db.Ping(); err != nil {
		db.Close()
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

// migrate ensures the tables we need exist.
func (s *Store) migrate() error {
	schema := `
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS games (
	id          TEXT PRIMARY KEY,
	name        TEXT NOT NULL,
	company_id  TEXT
);

CREATE TABLE IF NOT EXISTS players (
	id             TEXT PRIMARY KEY,
	alias          TEXT,
	integration_id TEXT NOT NULL,
	FOREIGN KEY (integration_id) REFERENCES games(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS balances (
	player_id      TEXT NOT NULL,
	integration_id TEXT NOT NULL,
	token          TEXT NOT NULL,
	amount         INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (player_id, integration_id, token),
	FOREIGN KEY (player_id)      REFERENCES players(id) ON DELETE CASCADE,
	FOREIGN KEY (integration_id) REFERENCES games(id)   ON DELETE CASCADE
);
`
	_, err := s.db.Exec(schema)
	return err
}

// RegisterIntegration inserts or updates a game/integration.
func (s *Store) RegisterIntegration(ctx context.Context, game *models.Integration) error {
	_, err := s.db.ExecContext(
		ctx,
		`INSERT INTO games (id, name, company_id)
         VALUES (?, ?, ?)
         ON CONFLICT(id) DO UPDATE
         SET name = excluded.name,
             company_id = excluded.company_id`,
		game.ID,
		game.Name,
		game.CompanyID,
	)
	return err
}

// RegisterPlayer inserts or updates a player.
func (s *Store) RegisterPlayer(ctx context.Context, p *models.Player) error {
	_, err := s.db.ExecContext(
		ctx,
		`INSERT INTO players (id, alias, integration_id)
         VALUES (?, ?, ?)
         ON CONFLICT(id) DO UPDATE
         SET alias = excluded.alias,
             integration_id = excluded.integration_id`,
		p.ID,
		p.Alias,
		p.IntegrationID,
	)
	return err
}

// UpdateBalance adds a positive amount to a player's balance for a given token.
func (s *Store) UpdateBalance(
	ctx context.Context,
	integrationID string,
	playerID string,
	token string,
	amount int64,
) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		}
	}()

	// Ensure the balance row exists
	if _, err = tx.ExecContext(
		ctx,
		`INSERT INTO balances (player_id, integration_id, token, amount)
         VALUES (?, ?, ?, 0)
         ON CONFLICT(player_id, integration_id, token) DO NOTHING`,
		playerID,
		integrationID,
		token,
	); err != nil {
		return err
	}

	// Increment the amount
	if _, err = tx.ExecContext(
		ctx,
		`UPDATE balances
         SET amount = amount + ?
         WHERE player_id = ? AND integration_id = ? AND token = ?`,
		amount,
		playerID,
		integrationID,
		token,
	); err != nil {
		return err
	}

	return tx.Commit()
}

// GetPlayer loads a player and all of their balances for a given integration.
func (s *Store) GetPlayer(
	ctx context.Context,
	integrationID string,
	playerID string,
) (*models.Player, error) {
	var (
		id   string
		alias string
		integ string
	)

	row := s.db.QueryRowContext(
		ctx,
		`SELECT id, alias, integration_id
         FROM players
         WHERE id = ? AND integration_id = ?`,
		playerID,
		integrationID,
	)

	if err := row.Scan(&id, &alias, &integ); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("player not found")
		}
		return nil, err
	}

	// Load balances
	rows, err := s.db.QueryContext(
		ctx,
		`SELECT token, amount
         FROM balances
         WHERE player_id = ? AND integration_id = ?`,
		playerID,
		integrationID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	balances := make(map[string]int64)
	for rows.Next() {
		var token string
		var amt int64
		if err := rows.Scan(&token, &amt); err != nil {
			return nil, err
		}
		balances[token] = amt
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &models.Player{
		ID:            id,
		Alias:         alias,
		IntegrationID: integ,
		Balances:      balances,
	}, nil
}