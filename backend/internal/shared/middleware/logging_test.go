package middleware

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestLogging(t *testing.T) {
	tests := []struct {
		name       string
		inner      http.HandlerFunc
		wantStatus float64 // JSON numbers decode to float64
	}{
		{"explicit status is recorded", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusTeapot)
		}, http.StatusTeapot},
		{"implicit 200 is recorded", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("ok"))
		}, http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			logger := slog.New(ContextHandler{Handler: slog.NewJSONHandler(&buf, nil)})
			h := RequestID(Logging(logger, tt.inner))

			req := httptest.NewRequest(http.MethodGet, "/x", nil)
			req.Header.Set(HeaderRequestID, "trace-1")
			h.ServeHTTP(httptest.NewRecorder(), req)

			var entry map[string]any
			if err := json.Unmarshal(buf.Bytes(), &entry); err != nil {
				t.Fatalf("log line is not JSON: %v (%q)", err, buf.String())
			}
			if entry["request_id"] != "trace-1" {
				t.Errorf("request_id = %v, want trace-1", entry["request_id"])
			}
			if entry["status"] != tt.wantStatus {
				t.Errorf("status = %v, want %v", entry["status"], tt.wantStatus)
			}
			if entry["path"] != "/x" {
				t.Errorf("path = %v, want /x", entry["path"])
			}

		})
	}
}
