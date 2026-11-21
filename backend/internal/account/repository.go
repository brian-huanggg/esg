package account

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

var ErrEmployeeIDTaken = errors.New("employee id already taken")
var ErrNotFound = errors.New("account not found")

// mysqlDupEntry is MariaDB's error number for a UNIQUE violation.
const mysqlDupEntry = 1062

type Account struct {
	ID           int64
	UUID         uuid.UUID
	EmployeeID   string
	Name         string
	PasswordHash string
	CreatedAt    time.Time
}

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, a *Account) error {
	a.UUID = uuid.New()

	_, err := r.db.ExecContext(ctx,
		`INSERT INTO accounts (uuid, employee_id, name, password_hash) VALUES (?, ?, ?, ?)`,
		a.UUID.String(), a.EmployeeID, a.Name, a.PasswordHash,
	)
	if err != nil {
		var me *mysql.MySQLError
		// The uuid column is UNIQUE too; only the employee_id key means "taken".
		if errors.As(err, &me) && me.Number == mysqlDupEntry &&
			strings.Contains(me.Message, "uq_accounts_employee_id") {
			return ErrEmployeeIDTaken
		}
		return fmt.Errorf("insert account: %w", err)
	}
	err = r.db.QueryRowContext(ctx,
		`SELECT created_at FROM accounts WHERE uuid = ?`, a.UUID.String(),
	).Scan(&a.CreatedAt)
	if err != nil {
		return fmt.Errorf("read back created_at: %w", err)
	}
	return nil
}

func (r *Repository) FindByEmployeeID(ctx context.Context, employeeID string) (*Account, error) {
	a := &Account{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, uuid, employee_id, name, password_hash, created_at
                 FROM accounts WHERE employee_id = ?`, employeeID,
	).Scan(&a.ID, &a.UUID, &a.EmployeeID, &a.Name, &a.PasswordHash, &a.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find account: %w", err)
	}
	return a, nil
}

func (r *Repository) FindByID(ctx context.Context, id int64) (*Account, error) {
	a := &Account{}
	err := r.db.QueryRowContext(ctx,
		`SELECT id, uuid, employee_id, name, password_hash, created_at
                 FROM accounts WHERE id = ?`, id,
	).Scan(&a.ID, &a.UUID, &a.EmployeeID, &a.Name, &a.PasswordHash, &a.CreatedAt)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("find account by id: %w", err)
	}
	return a, nil
}
