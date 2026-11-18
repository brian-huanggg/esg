package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthRoutes(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string // empty means 'don't check the body'
	}{
		{"get health return ok", http.MethodGet, "/api/health", http.StatusOK, `{"status":"ok"}` + "\n"},
		{"post health is not allowed", http.MethodPost, "/api/health", http.StatusMethodNotAllowed, ""},
		{"unknown path is not found", http.MethodGet, "/api/nope", http.StatusNotFound, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rec := httptest.NewRecorder()

			newMux().ServeHTTP(rec, req)

			if rec.Code != tt.wantStatus {
				t.Errorf("status = %d, want %d", rec.Code, tt.wantStatus)
			}
			if tt.wantBody == "" {
				return
			}
			if got := rec.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q, want application/json", got)
			}
			if got := rec.Body.String(); got != tt.wantBody {
				t.Errorf("Body = %q, want %q", got, tt.wantBody)
			}
		})
	}

}
