package proposal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/google/uuid"

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

func seedAccount(t *testing.T, employeeID string) int64 {
	t.Helper()
	res, err := testDB.ExecContext(context.Background(),
		`INSERT INTO accounts (uuid, employee_id, name, password_hash) VALUES (?, ?, ?, ?)`,
		uuid.NewString(), employeeID, "user-"+employeeID, "hash")
	if err != nil {
		t.Fatal(err)
	}
	id, err := res.LastInsertId()
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func TestCreate_Success(t *testing.T) {
	dbtest.Reset(t, testDB)
	leader := seedAccount(t, "100001")

	p := &Proposal{LeaderAccountID: leader, Title: "Solar roof", Description: "panels"}
	if err := NewRepository(testDB).Create(context.Background(), p); err != nil {
		t.Fatal(err)
	}
	if p.ID == 0 || p.UUID == uuid.Nil {
		t.Errorf("ID/UUID not filled in: %+v", p)
	}
	if p.Version != 1 {
		t.Errorf("Version = %d, want 1", p.Version)
	}
	if p.CreatedAt.IsZero() || p.UpdatedAt.IsZero() {
		t.Error("timestamps not filled in")
	}
}

func TestCreate_DuplicateTitleSameLeader(t *testing.T) {
	dbtest.Reset(t, testDB)
	ctx := context.Background()
	repo := NewRepository(testDB)
	leader := seedAccount(t, "100001")

	if err := repo.Create(ctx, &Proposal{LeaderAccountID: leader, Title: "Solar roof"}); err != nil {
		t.Fatal(err)
	}
	// The column collation is case-insensitive, so this still collides.
	err := repo.Create(ctx, &Proposal{LeaderAccountID: leader, Title: "SOLAR ROOF"})
	if !errors.Is(err, ErrDuplicateTitle) {
		t.Fatalf("err = %v, want ErrDuplicateTitle", err)
	}
}

func TestCreate_SameTitleDifferentLeader(t *testing.T) {
	dbtest.Reset(t, testDB)
	ctx := context.Background()
	repo := NewRepository(testDB)

	for _, emp := range []string{"100001", "100002"} {
		p := &Proposal{LeaderAccountID: seedAccount(t, emp), Title: "Solar roof"}
		if err := repo.Create(ctx, p); err != nil {
			t.Fatalf("leader %s: %v", emp, err)
		}
	}
}

func TestFindByUUID_Owner(t *testing.T) {
	dbtest.Reset(t, testDB)
	ctx := context.Background()
	repo := NewRepository(testDB)
	leader := seedAccount(t, "100001")

	created := &Proposal{LeaderAccountID: leader, Title: "Solar roof", Description: "panels"}
	if err := repo.Create(ctx, created); err != nil {
		t.Fatal(err)
	}

	got, err := repo.FindByUUID(ctx, leader, created.UUID)
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != created.ID || got.Title != "Solar roof" || got.Description != "panels" || got.Version != 1 {
		t.Errorf("got %+v, want fields of %+v", got, created)
	}
}

func TestFindByUUID_OtherLeader(t *testing.T) {
	dbtest.Reset(t, testDB)
	ctx := context.Background()
	repo := NewRepository(testDB)
	owner := seedAccount(t, "100001")
	stranger := seedAccount(t, "100002")

	p := &Proposal{LeaderAccountID: owner, Title: "Solar roof"}
	if err := repo.Create(ctx, p); err != nil {
		t.Fatal(err)
	}

	_, err := repo.FindByUUID(ctx, stranger, p.UUID)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestFindByUUID_Unknown(t *testing.T) {
	dbtest.Reset(t, testDB)
	leader := seedAccount(t, "100001")

	_, err := NewRepository(testDB).FindByUUID(context.Background(), leader, uuid.New())
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}
