package session

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// ErrNotFound covers unknown, idle-expired and absolute-expired sessions alike;
// callers must not be able to tell them apart.
var ErrNotFound = errors.New("session not found")

const (
	idleTimeoutMinutes   = 30
	absoluteTimeoutHours = 8
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, tokenHash []byte, csrfHash []byte, accountID int64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO sessions (token_hash, csrf_token_hash, account_id, expires_at)
                 VALUES (?, ?,  ?, UTC_TIMESTAMP(3) + INTERVAL ? HOUR)`,
		tokenHash, csrfHash, accountID, absoluteTimeoutHours,
	)
	if err != nil {
		return fmt.Errorf("insert session: %w", err)
	}
	return nil
}

func (r *Repository) FindActive(ctx context.Context, tokenHash []byte) (int64, error) {
	var accountID int64
	err := r.db.QueryRowContext(ctx,
		`SELECT account_id FROM sessions
                 WHERE token_hash = ?
                   AND last_seen_at > UTC_TIMESTAMP(3) - INTERVAL ? MINUTE
                   AND expires_at > UTC_TIMESTAMP(3)`,
		tokenHash, idleTimeoutMinutes,
	).Scan(&accountID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, ErrNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("find session: %w", err)
	}
	return accountID, nil
}

func (r *Repository) Touch(ctx context.Context, tokenHash []byte) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE sessions SET last_seen_at = UTC_TIMESTAMP(3)
                 WHERE token_hash = ?
                   AND last_seen_at > UTC_TIMESTAMP(3) - INTERVAL ? MINUTE
                   AND expires_at > UTC_TIMESTAMP(3)`,
		tokenHash, idleTimeoutMinutes,
	)
	if err != nil {
		return fmt.Errorf("touch session: %w", err)
	}
	return nil
}

func (r *Repository) Delete(ctx context.Context, tokenHash []byte) error {
	_, err := r.db.ExecContext(ctx,
		`DELETE FROM sessions WHERE token_hash = ?`, tokenHash)
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	return nil
}
