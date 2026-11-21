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

func TestFindByEmployeeID_Found(t *testing.T) {
	dbtest.Reset(t, testDB)
	ctx := context.Background()
	repo := NewRepository(testDB)

	created := &Account{EmployeeID: "123456", Name: "Alice", PasswordHash: "hash"}
	if err := repo.Create(ctx, created); err != nil {
		t.Fatal(err)
	}

	got, err := repo.FindByEmployeeID(ctx, "123456")
	if err != nil {
		t.Fatal(err)
	}
	if got.ID == 0 {
		t.Error("ID not filled in")
	}
	if got.UUID != created.UUID || got.Name != "Alice" || got.PasswordHash != "hash" {
		t.Errorf("got %+v, want fields of %+v", got, created)
	}
}

func TestFindByEmployeeID_NotFound(t *testing.T) {
	dbtest.Reset(t, testDB)

	_, err := NewRepository(testDB).FindByEmployeeID(context.Background(), "999999")
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestFindByID_Found(t *testing.T) {
	dbtest.Reset(t, testDB)
	ctx := context.Background()
	repo := NewRepository(testDB)

	if err := repo.Create(ctx, &Account{EmployeeID: "123456", Name: "Alice", PasswordHash: "hash"}); err != nil {
		t.Fatal(err)
	}
	// Create does not return the BIGINT id, so look it up the same way login does.
	byEmp, err := repo.FindByEmployeeID(ctx, "123456")
	if err != nil {
		t.Fatal(err)
	}

	got, err := repo.FindByID(ctx, byEmp.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.UUID != byEmp.UUID || got.Name != "Alice" {
		t.Errorf("got %+v, want %+v", got, byEmp)
	}
}

func TestFindByID_NotFound(t *testing.T) {
	dbtest.Reset(t, testDB)

	_, err := NewRepository(testDB).FindByID(context.Background(), 999999)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}
