package account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
)

var ErrEmployeeIDTaken = errors.New("employee id already taken")

// mysqlDupEntry is MariaDB's error number for a UNIQUE violation.
const mysqlDupEntry = 1062

type Account struct {
	UUID         uuid.UUID
	EmployeeID   string
	Name         string
	PasswordHash string
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
	return nil
}
