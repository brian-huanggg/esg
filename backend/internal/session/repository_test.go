package session

import (
	"context"
	"crypto/sha256"
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
	stop()
	os.Exit(code)
}

func setup(t *testing.T) (*Repository, []byte, int64) {
	t.Helper()
	dbtest.Reset(t, testDB)
	ctx := context.Background()

	res, err := testDB.ExecContext(ctx,
		`INSERT INTO accounts (uuid, employee_id, name, password_hash) VALUES (UUID(), '123456', 'Alice', 'h')`)
	if err != nil {
		t.Fatal(err)
	}
	accountID, _ := res.LastInsertId()

	sum := sha256.Sum256([]byte("token"))
	hash := sum[:]
	csrfSum := sha256.Sum256([]byte("csrf"))
	repo := NewRepository(testDB)
	if err := repo.Create(ctx, hash, csrfSum[:], accountID); err != nil {
		t.Fatal(err)
	}
	return repo, hash, accountID
}

func backdate(t *testing.T, col string, hash []byte, minutes int) {
	t.Helper()
	_, err := testDB.Exec(
		"UPDATE sessions SET "+col+" = UTC_TIMESTAMP(3) - INTERVAL ? MINUTE WHERE token_hash = ?",
		minutes, hash)
	if err != nil {
		t.Fatal(err)
	}
}

func TestFindActive_Valid(t *testing.T) {
	repo, hash, accountID := setup(t)
	got, err := repo.FindActive(context.Background(), hash)
	if err != nil {
		t.Fatal(err)
	}
	if got != accountID {
		t.Errorf("accountID = %d, want %d", got, accountID)
	}
}

func TestFindActive_IdleTimeout(t *testing.T) {
	repo, hash, _ := setup(t)
	backdate(t, "last_seen_at", hash, 31)
	_, err := repo.FindActive(context.Background(), hash)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestFindActive_AbsoluteTimeout(t *testing.T) {
	repo, hash, _ := setup(t)
	// Recently active, but past its 8 hour lifetime.
	backdate(t, "expires_at", hash, 1)
	_, err := repo.FindActive(context.Background(), hash)
	if !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
}

func TestFindActive_Delete(t *testing.T) {
	repo, hash, _ := setup(t)
	ctx := context.Background()
	if err := repo.Delete(ctx, hash); err != nil {
		t.Fatal(err)
	}
	if _, err := repo.FindActive(ctx, hash); !errors.Is(err, ErrNotFound) {
		t.Fatalf("err = %v, want ErrNotFound", err)
	}
	// Second delete must still succeed (idempotent).
	if err := repo.Delete(ctx, hash); err != nil {
		t.Fatal(err)
	}
}

func TestTouch_SlidesIdleWindow(t *testing.T) {
	repo, hash, _ := setup(t)
	backdate(t, "last_seen_at", hash, 20)
	if err := repo.Touch(context.Background(), hash); err != nil {
		t.Fatal(err)
	}
	var fresh bool
	err := testDB.QueryRow(
		`SELECT last_seen_at > UTC_TIMESTAMP(3) - INTERVAL 1 MINUTE FROM sessions WHERE token_hash = ?`,
		hash).Scan(&fresh)
	if err != nil {
		t.Fatal(err)
	}
	if !fresh {
		t.Error("last_seen_at not refreshed")
	}
}
