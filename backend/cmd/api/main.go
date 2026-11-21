package main

import (
	"context"
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/account"
	"github.com/brian-huanggg/esg-award-registration/backend/internal/auth"
	"github.com/brian-huanggg/esg-award-registration/backend/internal/session"
	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared"
	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/api"
	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/middleware"
)

func main() {
	logger := slog.New(middleware.ContextHandler{Handler: slog.NewJSONHandler(os.Stdout, nil)})
	slog.SetDefault(logger)
	if err := run(logger); err != nil {
		logger.Error("server stopped", "err", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	cfg, err := shared.LoadDBConfig()
	if err != nil {
		return err
	}
	db, err := shared.OpenDB(context.Background(), cfg)
	if err != nil {
		return err
	}
	defer db.Close()

	logger.Info("listening", "addr", ":8080")
	return http.ListenAndServe(":8080", newHandler(logger, db))
}

type server struct {
	*shared.HealthHandler
	*account.Handler
	*auth.SessionHandler
}

func newHandler(logger *slog.Logger, db *sql.DB) http.Handler {
	repo := account.NewRepository(db)
	svc := account.NewService(repo)
	accountHandler := account.NewHandler(svc, logger)

	sessionRepo := session.NewRepository(db)
	authSvc := auth.NewService(repo, sessionRepo)
	sessionHandler := auth.NewSessionHandler(authSvc, repo, logger)

	authMW := middleware.Authenticate(sessionRepo, logger, map[string]bool{
		"GetCurrentAccount": true,
	})

	mux := http.NewServeMux()
	strict := api.NewStrictHandler(
		server{&shared.HealthHandler{}, accountHandler, sessionHandler},
		[]api.StrictMiddlewareFunc{authMW},
	)
	api.HandlerFromMux(strict, mux)
	return middleware.RequestID(middleware.Logging(logger, mux))

}
