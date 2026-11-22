package proposal

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"strings"
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

func defaultList() ListParams {
	return ListParams{Page: 1, PageSize: 20, Sort: "createdAt", Order: "desc"}
}

func seedTitles(t *testing.T, repo *Repository, leader int64, titles ...string) {
	t.Helper()
	for _, title := range titles {
		if err := repo.Create(context.Background(), &Proposal{LeaderAccountID: leader, Title: title}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestList_TieBreakerKeepsPagesStable(t *testing.T) {
	dbtest.Reset(t, testDB)
	ctx := context.Background()
	repo := NewRepository(testDB)
	leader := seedAccount(t, "100001")
	seedTitles(t, repo, leader, "a", "b", "c", "d", "e")
	// Force identical sort values so only the id tie-breaker decides the order.
	if _, err := testDB.ExecContext(ctx, `UPDATE proposals SET created_at = '2026-01-01 00:00:00.000'`); err != nil {
		t.Fatal(err)
	}

	var got []string
	for page := 1; page <= 3; page++ {
		p := defaultList()
		p.Page, p.PageSize, p.Order = page, 2, "asc"
		items, total, err := repo.List(ctx, leader, p)
		if err != nil {
			t.Fatal(err)
		}
		if total != 5 {
			t.Errorf("total = %d, want 5", total)
		}
		for _, it := range items {
			got = append(got, it.Title)
		}
	}
	if want := "a b c d e"; strings.Join(got, " ") != want {
		t.Errorf("got %v, want %s", got, want)
	}
}

func TestList_PageSizeCappedAt100(t *testing.T) {
	dbtest.Reset(t, testDB)
	ctx := context.Background()
	repo := NewRepository(testDB)
	leader := seedAccount(t, "100001")
	for i := 0; i < 101; i++ {
		seedTitles(t, repo, leader, fmt.Sprintf("p%03d", i))
	}

	p := defaultList()
	p.PageSize = 1000
	items, total, err := repo.List(ctx, leader, p)
	if err != nil {
		t.Fatal(err)
	}
	if len(items) != 100 || total != 101 {
		t.Errorf("len = %d, total = %d, want 100 and 101", len(items), total)
	}
}

func TestList_SearchEscapesWildcards(t *testing.T) {
	dbtest.Reset(t, testDB)
	repo := NewRepository(testDB)
	leader := seedAccount(t, "100001")
	seedTitles(t, repo, leader, "Save 50% energy", "Save 500 energy", "a_b", "axb")

	for q, want := range map[string]string{"50%": "Save 50% energy", "a_b": "a_b"} {
		p := defaultList()
		p.Q = q
		items, total, err := repo.List(context.Background(), leader, p)
		if err != nil {
			t.Fatal(err)
		}
		if total != 1 || len(items) != 1 || items[0].Title != want {
			t.Errorf("q=%q: got %d items (total %d), want only %q", q, len(items), total, want)
		}
	}
}

func TestList_OnlyOwnProposals(t *testing.T) {
	dbtest.Reset(t, testDB)
	repo := NewRepository(testDB)
	owner := seedAccount(t, "100001")
	stranger := seedAccount(t, "100002")
	seedTitles(t, repo, owner, "mine")
	seedTitles(t, repo, stranger, "theirs")

	items, total, err := repo.List(context.Background(), owner, defaultList())
	if err != nil {
		t.Fatal(err)
	}
	if total != 1 || len(items) != 1 || items[0].Title != "mine" {
		t.Errorf("got %+v (total %d), want only \"mine\"", items, total)
	}
}

func TestList_RejectsUnknownSortAndOrder(t *testing.T) {
	dbtest.Reset(t, testDB)
	repo := NewRepository(testDB)
	leader := seedAccount(t, "100001")

	for _, p := range []ListParams{
		{Page: 1, PageSize: 20, Sort: "title; DROP TABLE proposals", Order: "asc"},
		{Page: 1, PageSize: 20, Sort: "title", Order: "sideways"},
	} {
		if _, _, err := repo.List(context.Background(), leader, p); !errors.Is(err, ErrInvalidSort) {
			t.Errorf("%+v: err = %v, want ErrInvalidSort", p, err)
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
