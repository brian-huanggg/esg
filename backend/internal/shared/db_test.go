package shared

import (
	"context"
	"testing"
	"time"

	"github.com/go-sql-driver/mysql"
)

func TestLoadDBConfig_MissingPassword(t *testing.T) {
	t.Setenv("DB_NAME", "esg")
	t.Setenv("DB_USER", "esg")
	t.Setenv("DB_PASSWORD", "")
	if _, err := LoadDBConfig(); err == nil {
		t.Fatal("expected error when DB_PASSWORD is empty")
	}
}

func TestLoadDBConfig_Defaults(t *testing.T) {
	t.Setenv("DB_NAME", "esg")
	t.Setenv("DB_USER", "u")
	t.Setenv("DB_PASSWORD", "p")
	t.Setenv("DB_HOST", "")
	t.Setenv("DB_PORT", "")
	cfg, err := LoadDBConfig()
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Host != "127.0.0.1" || cfg.Port != "3306" {
		t.Fatalf("unexpected defaults: %+v", cfg)
	}
}

func TestDSN_HasUTCAndParseTime(t *testing.T) {
	dsn := DBConfig{Host: "h", Port: "1", Name: "n", User: "u", Password: "p@:/x"}.DSN()
	// Parse instead of substring-matching: FormatDSN omits loc when it equals the driver default (UTC).
	m, err := mysql.ParseDSN(dsn)
	if err != nil {
		t.Fatal(err)
	}
	if !m.ParseTime || m.Loc != time.UTC {
		t.Fatalf("want parseTime and UTC, got %+v", m)
	}
	if m.Passwd != "p@:/x" {
		t.Fatalf("password not round-tripped: %q", m.Passwd)
	}
}

func TestOpenDB_Unreachable(t *testing.T) {
	cfg := DBConfig{Host: "127.0.0.1", Port: "1", Name: "n", User: "u", Password: "p"}
	if _, err := OpenDB(context.Background(), cfg); err == nil {
		t.Fatal("expected ping error")
	}
}
