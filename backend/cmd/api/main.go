package main

import (
	"log/slog"
	"net/http"
	"os"

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
	logger.Info("listening", "addr", ":8080")
	return http.ListenAndServe(":8080", newHandler(logger))
}

func newHandler(logger *slog.Logger) http.Handler {
	mux := http.NewServeMux()
	strict := api.NewStrictHandler(&shared.HealthHandler{}, nil)
	api.HandlerFromMux(strict, mux)
	return middleware.RequestID(middleware.Logging(logger, mux))
}
