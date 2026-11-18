package middleware

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRequestID(t *testing.T) {
	tests := []struct {
		name     string
		incoming string // empty means 'send no header'
		wantSame bool   // true: the incoming id must be reused as-is
	}{
		{"no header generates one", "", false},
		{"valid header is reused", "abc-123", true},
		{"newline is rejected", "abc\nevil", false},
		{"too long is rejected", strings.Repeat("a", 65), false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var seen string
			h := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				seen = RequestIDFrom(r.Context())
			}))
			req := httptest.NewRequest(http.MethodGet, "/", nil)
			if tt.incoming != "" {
				req.Header.Set(HeaderRequestID, tt.incoming)
			}
			rec := httptest.NewRecorder()

			h.ServeHTTP(rec, req)

			got := rec.Header().Get(HeaderRequestID)
			if got == "" {
				t.Fatal("response has no X-Request-Id")
			}
			if seen != got {
				t.Errorf("context id = %q, header id = %q, want equal", seen, got)
			}
			if (got == tt.incoming) != tt.wantSame {
				t.Errorf("id = %q, incoming = %q, wantSame = %v", got, tt.incoming, tt.wantSame)
			}
		})
	}
}
