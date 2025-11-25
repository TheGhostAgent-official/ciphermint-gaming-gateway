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

const DefaultDBPath = "./ciphermint_gateway.db?_fk=1"

// OpenDefault opens (or creates) the SQLite DB file.
func OpenDefault() (*Store, error) {
	db, err := sql.Open("sqlite3", DefaultDBPath)
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

// Close closes the underlying DB safely.
func (s *Store) Close() error {
	if s == nil || s.db == nil {
		return nil
	}
	return s.db.Close()
}

// migrate ensures the tables we need exist.
func (s *Store) migrate() error {
	schema := `
CREATE TABLE IF NOT EXISTS integrations (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL,
	company_id TEXT
);

CREATE TABLE IF NOT EXISTS players (
	id TEXT PRIMARY KEY,
	alias TEXT,
	integration_id TEXT NOT NULL,
	FOREIGN KEY (integration_id) REFERENCES integrations(id)
);

CREATE TABLE IF NOT EXISTS balances (
	player_id TEXT NOT NULL,
	integration_id TEXT NOT NULL,
	token TEXT NOT NULL,
	amount INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (player_id, integration_id, token),
	FOREIGN KEY (player_id) REFERENCES players(id),
	FOREIGN KEY (integration_id) REFERENCES integrations(id)
);
`
	_, err := s.db.Exec(schema)
	return err
}

// RegisterIntegration inserts or updates an integration.
func (s *Store) RegisterIntegration(ctx context.Context, integ *models.Integration) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO integrations (id, name, company_id)
		 VALUES (?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET
		   name = excluded.name,
		   company_id = excluded.company_id`,
		integ.ID, integ.Name, integ.CompanyID,
	)
	return err
}

// RegisterPlayer inserts or updates a player.
func (s *Store) RegisterPlayer(ctx context.Context, p *models.Player) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO players (id, alias, integration_id)
		 VALUES (?, ?, ?)
		 ON CONFLICT(id) DO UPDATE SET
		   alias = excluded.alias,
		   integration_id = excluded.integration_id`,
		p.ID, p.Alias, p.IntegrationID,
	)
	return err
}

// UpdateBalance upserts a balance row and adds the amount.
func (s *Store) UpdateBalance(ctx context.Context, integrationID, playerID, token string, amount int64, source string) error {
	_, err := s.db.ExecContext(ctx,
		`INSERT INTO balances (player_id, integration_id, token, amount)
		 VALUES (?, ?, ?, ?)
		 ON CONFLICT(player_id, integration_id, token)
		 DO UPDATE SET amount = balances.amount + excluded.amount`,
		playerID, integrationID, token, amount,
	)
	return err
}

// GetPlayer loads a player and all balances.
func (s *Store) GetPlayer(ctx context.Context, integrationID, playerID string) (*models.Player, error) {
	row := s.db.QueryRowContext(ctx,
		`SELECT id, alias, integration_id
		   FROM players
		  WHERE id = ? AND integration_id = ?`,
		playerID, integrationID,
	)

	var p models.Player
	if err := row.Scan(&p.ID, &p.Alias, &p.IntegrationID); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	rows, err := s.db.QueryContext(ctx,
		`SELECT token, amount
		   FROM balances
		  WHERE player_id = ? AND integration_id = ?`,
		playerID, integrationID,
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return &p, nil
}
