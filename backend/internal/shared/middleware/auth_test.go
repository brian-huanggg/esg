package middleware

import (
	"context"
	"crypto/sha256"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/session"
)

type fakeSessions struct {
	accountID int64
	findErr   error
	touched   bool
}

func (f *fakeSessions) FindActive(context.Context, []byte) (int64, error) {
	return f.accountID, f.findErr
}

func (f *fakeSessions) Touch(context.Context, []byte) error {
	f.touched = true
	return nil
}

// run sends one request through the middleware and reports whether the
// wrapped handler was reached and what context it saw.
func run(fake *fakeSessions, operationID, cookie string) (code int, called bool, seen context.Context) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	mw := Authenticate(fake, logger, map[string]bool{"Protected": true})
	h := mw(func(ctx context.Context, _ http.ResponseWriter, _ *http.Request, _ any) (any, error) {
		called, seen = true, ctx
		return nil, nil
	}, operationID)

	r := httptest.NewRequest(http.MethodGet, "/", nil)
	if cookie != "" {
		r.AddCookie(&http.Cookie{Name: SessionCookieName, Value: cookie})
	}
	w := httptest.NewRecorder()
	_, _ = h(r.Context(), w, r, nil)
	return w.Code, called, seen
}

func TestAuthenticate_Rejects(t *testing.T) {
	tests := []struct {
		name    string
		cookie  string
		findErr error
		want    int
	}{
		{"no cookie", "", nil, http.StatusUnauthorized},
		{"unknown or expired", "t", session.ErrNotFound, http.StatusUnauthorized},
		{"db failure", "t", errors.New("boom"), http.StatusInternalServerError},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, called, _ := run(&fakeSessions{findErr: tt.findErr}, "Protected", tt.cookie)
			if code != tt.want || called {
				t.Errorf("code = %d called = %v, want %d and not called", code, called, tt.want)
			}
		})
	}
}

func TestAuthenticate_Valid(t *testing.T) {
	fake := &fakeSessions{accountID: 7}
	_, called, ctx := run(fake, "Protected", "token")
	if !called || !fake.touched {
		t.Fatalf("called = %v touched = %v, want both true", called, fake.touched)
	}
	if id, ok := AccountIDFrom(ctx); !ok || id != 7 {
		t.Errorf("account id = %d, %v, want 7, true", id, ok)
	}
	want := sha256.Sum256([]byte("token"))
	if h, _ := TokenHashFrom(ctx); string(h) != string(want[:]) {
		t.Error("token hash in context does not match sha256 of the cookie value")
	}
}

func TestAuthenticate_UnprotectedSkipsAuth(t *testing.T) {
	_, called, ctx := run(&fakeSessions{}, "Public", "")
	if !called {
		t.Fatal("public operation must reach the handler")
	}
	if _, ok := AccountIDFrom(ctx); ok {
		t.Error("public operation must not get an account id")
	}
}
