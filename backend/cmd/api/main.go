package main

import (
	"log"
	"net/http"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared"
	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/api"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	log.Println("listening on :8080")
	return http.ListenAndServe(":8080", newMux())
}

func newMux() *http.ServeMux {
	mux := http.NewServeMux()
	strict := api.NewStrictHandler(&shared.HealthHandler{}, nil)
	api.HandlerFromMux(strict, mux)
	return mux
}
