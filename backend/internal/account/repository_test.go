package account

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/dbtest"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	db, stop, err := dbtest.Start(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	testDB = db
	code := m.Run()
	// os.Exit skips defers, so stop the container explicitly first.
	stop()
	os.Exit(code)
}

func TestCreate_Success(t *testing.T) {
	dbtest.Reset(t, testDB)
	ctx := context.Background()

	a := &Account{EmployeeID: "123456", Name: "Alice", PasswordHash: "hash"}
	if err := NewRepository(testDB).Create(ctx, a); err != nil {
		t.Fatal(err)
	}

	var name string
	err := testDB.QueryRowContext(ctx,
		`SELECT name FROM accounts WHERE uuid = ?`, a.UUID.String()).Scan(&name)
	if err != nil {
		t.Fatal(err)
	}
	if a.CreatedAt.IsZero() {
		t.Error("CreatedAt not filled in")
	}
	if name != "Alice" {
		t.Errorf("name = %q, want Alice", name)
	}
}

func TestCreate_DuplicateEmployeeID(t *testing.T) {
	dbtest.Reset(t, testDB)
	ctx := context.Background()
	repo := NewRepository(testDB)

	if err := repo.Create(ctx, &Account{EmployeeID: "123456", Name: "A", PasswordHash: "h"}); err != nil {
		t.Fatal(err)
	}
	err := repo.Create(ctx, &Account{EmployeeID: "123456", Name: "B", PasswordHash: "h"})
	if !errors.Is(err, ErrEmployeeIDTaken) {
		t.Fatalf("err = %v, want ErrEmployeeIDTaken", err)
	}
}
