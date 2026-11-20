package dbtest

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/go-sql-driver/mysql"
)

var testDB *sql.DB

func TestMain(m *testing.M) {
	db, stop, err := Start(context.Background())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	testDB = db
	code := m.Run()
	stop() // must stop before os.Exit
	os.Exit(code)
}

const insertAccount = `INSERT INTO accounts (uuid, employee_id, name, password_hash) VALUES (UUID(), ?, 'Test', 'x')`

func TestAccounts_EmployeeIDUnique(t *testing.T) {
	Reset(t, testDB)
	if _, err := testDB.Exec(insertAccount, "123456"); err != nil {
		t.Fatal(err)
	}
	_, err := testDB.Exec(insertAccount, "123456")

	var me *mysql.MySQLError
	if !errors.As(err, &me) || me.Number != 1062 {
		t.Fatalf("want MYSQL error 1062 (duplicate entry), got %v", err)
	}
}

func TestReset_ClearRows(t *testing.T) {
	if _, err := testDB.Exec(insertAccount, "654321"); err != nil {
		t.Fatal(err)
	}
	Reset(t, testDB)

	var n int
	if err := testDB.QueryRow("SELECT COUNT(*) FROM accounts").Scan(&n); err != nil {
		t.Fatal(err)
	}
	if n != 0 {
		t.Fatalf("want 0 rows after Reset, got %d", n)
	}
}
