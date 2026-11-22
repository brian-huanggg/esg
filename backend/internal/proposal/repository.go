package proposal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

var ErrDuplicateTitle = errors.New("proposal title already used by this leader")

var ErrNotFound = errors.New("proposal not found")

// mysqlDupEntry is MariaDB's error number for a UNIQUE violation.
const mysqlDupEntry = 1062

type Proposal struct {
	ID              int64
	UUID            uuid.UUID
	LeaderAccountID int64
	Title           string
	Description     string
	Version         int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, p *Proposal) error {
	p.UUID = uuid.New()

	res, err := r.db.ExecContext(ctx,
		`INSERT INTO proposals (uuid, leader_account_id, title, description) VALUES (?, ?, ?, ?)`,
		p.UUID.String(), p.LeaderAccountID, p.Title, p.Description,
	)
	if err != nil {
		var me *mysql.MySQLError
		// The uuid column is UNIQUE too; only the (leader, title) key is a user error.
		if errors.As(err, &me) && me.Number == mysqlDupEntry &&
			strings.Contains(me.Message, "uq_proposals_leader_title") {
			return ErrDuplicateTitle
		}
		return fmt.Errorf("insert proposal: %w", err)
	}
	p.ID, err = res.LastInsertId()
	if err != nil {
		return fmt.Errorf("last insert id: %w", err)
	}
	err = r.db.QueryRowContext(ctx,
		`SELECT version, created_at, updated_at FROM proposals WHERE id = ?`, p.ID,
	).Scan(&p.Version, &p.CreatedAt, &p.UpdatedAt)
	if err != nil {
		return fmt.Errorf("read back proposal defaults: %w", err)
	}
	return nil
}

// FindByUUID returns the proposal only if leaderID owns it.
func (r *Repository) FindByUUID(ctx context.Context, leaderID int64, id uuid.UUID) (*Proposal, error) {
	p := &Proposal{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, uuid, leader_account_id, title, description, version, created_at, updated_at
                 FROM proposals WHERE uuid = ? AND leader_account_id = ?`, id.String(), leaderID,
	).Scan(&p.ID, &p.UUID, &p.LeaderAccountID, &p.Title, &p.Description,
		&p.Version, &p.CreatedAt, &p.UpdatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find proposal: %w", err)
	}
	return p, nil
}
