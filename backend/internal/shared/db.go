package shared

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Host, Port, Name, User, Password string
}

func LoadDBConfig() (DBConfig, error) {
	cfg := DBConfig{
		Host:     getenv("DB_HOST", "127.0.0.1"),
		Port:     getenv("DB_PORT", "3306"),
		Name:     os.Getenv("DB_NAME"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
	}
	if cfg.Name == "" || cfg.User == "" || cfg.Password == "" {
		return DBConfig{}, errors.New("DB_NAME, DB_USER, DB_PASSWORD are required.")
	}
	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func (c DBConfig) DSN() string {
	m := mysql.NewConfig()
	m.User = c.User
	m.Passwd = c.Password
	m.Net = "tcp"
	m.Addr = net.JoinHostPort(c.Host, c.Port)
	m.DBName = c.Name
	m.ParseTime = true // scan DATETIME into time.Time
	m.Loc = time.UTC
	return m.FormatDSN()
}

func OpenDB(ctx context.Context, cfg DBConfig) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.DSN())
	if err != nil {
		return nil, fmt.Errorf("open db: %w", err)
	}
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)
	db.SetConnMaxLifetime(5 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := db.PingContext(pingCtx); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("ping db: %w", err)
	}
	return db, nil
}
