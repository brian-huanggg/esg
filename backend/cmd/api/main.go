package main

import (
	"log"
	"net/http"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared"
	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/api"
	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/middleware"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log.Println("listening on :8080")
	return http.ListenAndServe(":8080", newHandler())
}

func newHandler() http.Handler {
	mux := http.NewServeMux()
	strict := api.NewStrictHandler(&shared.HealthHandler{}, nil)
	api.HandlerFromMux(strict, mux)
	return middleware.RequestID(mux)
}
