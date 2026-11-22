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

const maxPageSize = 100

var ErrInvalidSort = errors.New("invalid sort field or order")

var sortColumns = map[string]string{
	"createdAt": "created_at",
	"updatedAt": "updated_at",
	"title":     "title",
}

type ListParams struct {
	Page     int
	PageSize int
	Sort     string
	Order    string
	Q        string
}

func (r *Repository) List(ctx context.Context, leaderID int64, p ListParams) ([]Proposal, int, error) {
	col, ok := sortColumns[p.Sort]
	if !ok {
		return nil, 0, ErrInvalidSort
	}
	var dir string

	switch p.Order {
	case "asc":
		dir = "ASC"
	case "desc":
		dir = "DESC"
	default:
		return nil, 0, ErrInvalidSort
	}
	pageSize := min(max(p.PageSize, 1), maxPageSize)
	offset := (max(p.Page, 1) - 1) * pageSize

	where := `leader_account_id = ? AND title LIKE ?`
	args := []any{leaderID, "%" + escapeLike(p.Q) + "%"}

	var total int
	if err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM proposals WHERE `+where, args...,
	).Scan(&total); err != nil {
		return nil, 0, fmt.Errorf("count proposals: %w", err)
	}

	// id is the tie-breaker: without it, rows with equal sort values may swap
	// between pages.
	rows, err := r.db.QueryContext(ctx,
		fmt.Sprintf(`SELECT id, uuid, leader_account_id, title, description, version, created_at, updated_at
                 FROM proposals WHERE %s ORDER BY %s %s, id %s LIMIT ? OFFSET ?`, where, col, dir, dir),
		append(args, pageSize, offset)...,
	)
	if err != nil {
		return nil, 0, fmt.Errorf("list proposals: %w", err)
	}
	defer rows.Close()

	items := []Proposal{}
	for rows.Next() {
		var it Proposal
		if err := rows.Scan(&it.ID, &it.UUID, &it.LeaderAccountID, &it.Title, &it.Description,
			&it.Version, &it.CreatedAt, &it.UpdatedAt); err != nil {
			return nil, 0, fmt.Errorf("scan proposal: %w", err)
		}
		items = append(items, it)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate proposals: %w", err)
	}
	return items, total, nil
}

// escapeLike makes %, _ and \ match literally; backslash is MariaDB's default LIKE escape.
func escapeLike(s string) string {
	return strings.NewReplacer(`\`, `\\`, `%`, `\%`, `_`, `\_`).Replace(s)
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
