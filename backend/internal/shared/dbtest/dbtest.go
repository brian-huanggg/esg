package dbtest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"path/filepath"
	"runtime"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mariadb"
)

func Start(ctx context.Context) (*sql.DB, func(), error) {
	ctr, err := mariadb.Run(ctx, "mariadb:11.4")
	if err != nil {
		_ = testcontainers.TerminateContainer(ctr)
		return nil, nil, fmt.Errorf("start mariadb: %w", err)
	}
	stop := func() { _ = testcontainers.TerminateContainer(ctr) }

	// The host port is random, so only the container knows the real address.
	base, err := ctr.ConnectionString(ctx)
	if err != nil {
		stop()
		return nil, nil, fmt.Errorf("connection string: %w", err)
	}

	m, err := migrate.New("file://"+migrationsDir(), "mysql://"+base)
	if err != nil {
		stop()
		return nil, nil, fmt.Errorf("new migrate: %w", err)
	}
	// ErrNoChange is not a failure: the schema is already current.
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		stop()
		return nil, nil, fmt.Errorf("migrate up: %w", err)
	}
	_, _ = m.Close()

	db, err := sql.Open("mysql", base+"?parseTime=true")
	if err != nil {
		stop()
		return nil, nil, fmt.Errorf("open db: %w", err)
	}
	return db, func() { _ = db.Close(); stop() }, nil
}

func migrationsDir() string {
	_, file, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(file), "..", "..", "..", "migrations")
}

// Reset empties every table so each test starts from a clean state.
func Reset(t *testing.T, db *sql.DB) {
	t.Helper()
	ctx := context.Background()

	// SET FOREIGN_KEY_CHECKS is per-session, and *sql.DB is a pool,
	// so pin one connection for the whole reset.
	conn, err := db.Conn(ctx)
	if err != nil {
		t.Fatal(err)
	}
	defer conn.Close()

	rows, err := conn.QueryContext(ctx, `SELECT table_name FROM information_schema.tables
                WHERE table_schema = DATABASE() AND table_type = 'BASE TABLE'
                AND table_name <> 'schema_migrations'`)
	if err != nil {
		t.Fatal(err)
	}
	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			t.Fatal(err)
		}
		tables = append(tables, name)
	}
	// Close before reusing conn: a connection can't run a new statement mid-result.
	rows.Close()
	if err := rows.Err(); err != nil {
		t.Fatal(err)
	}

	if _, err := conn.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 0"); err != nil {
		t.Fatal(err)
	}
	for _, name := range tables {
		// Identifiers can't be bound as parameters; names come from information_schema, not user input.
		if _, err := conn.ExecContext(ctx, "TRUNCATE TABLE `"+name+"`"); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := conn.ExecContext(ctx, "SET FOREIGN_KEY_CHECKS = 1"); err != nil {
		t.Fatal(err)
	}
}
