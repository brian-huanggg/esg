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
	active  session.Active
	findErr error
	touched bool
}

func (f *fakeSessions) FindActive(context.Context, []byte) (session.Active, error) {
	return f.active, f.findErr
}

func (f *fakeSessions) Touch(context.Context, []byte) error {
	f.touched = true
	return nil
}

func run(fake *fakeSessions, operationID, cookie, method, csrf string) (code int, called bool, seen context.Context) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	mw := Authenticate(fake, logger, map[string]bool{"Protected": true})
	h := mw(func(ctx context.Context, _ http.ResponseWriter, _ *http.Request, _ any) (any, error) {
		called, seen = true, ctx
		return nil, nil
	}, operationID)

	r := httptest.NewRequest(method, "/", nil)
	if cookie != "" {
		r.AddCookie(&http.Cookie{Name: SessionCookieName, Value: cookie})
	}
	if csrf != "" {
		r.Header.Set(csrfHeader, csrf)
	}
	w := httptest.NewRecorder()
	_, _ = h(r.Context(), w, r, nil)
	return w.Code, called, seen
}

func TestAuthenticate_Rejects(t *testing.T) {
	tests := []struct {
		name    string
		method  string
		cookie  string
		findErr error
		want    int
	}{
		{"no cookie", http.MethodGet, "", nil, http.StatusUnauthorized},
		{"unknown or expired", http.MethodGet, "t", session.ErrNotFound, http.StatusUnauthorized},
		{"db failure", http.MethodGet, "t", errors.New("boom"), http.StatusInternalServerError},
		{"no cookie on DELETE", http.MethodDelete, "", nil, http.StatusUnauthorized},
		{"expired session on DELETE", http.MethodDelete, "t", session.ErrNotFound, http.StatusUnauthorized},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, called, _ := run(&fakeSessions{findErr: tt.findErr}, "Protected", tt.cookie, tt.method, "")
			if code != tt.want || called {
				t.Errorf("code = %d called = %v, want %d and not called", code, called, tt.want)
			}
		})
	}
}

func TestAuthenticate_Valid(t *testing.T) {
	fake := &fakeSessions{active: session.Active{AccountID: 7}}
	_, called, ctx := run(fake, "Protected", "token", http.MethodGet, "")
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
	_, called, ctx := run(&fakeSessions{}, "Public", "", http.MethodGet, "")
	if !called {
		t.Fatal("public operation must reach the handler")
	}
	if _, ok := AccountIDFrom(ctx); ok {
		t.Error("public operation must not get an account id")
	}
}

func TestAuthenticate_CSRF(t *testing.T) {
	good := sha256.Sum256([]byte("good"))
	newFake := func() *fakeSessions {
		return &fakeSessions{active: session.Active{AccountID: 7, CSRFTokenHash: good[:]}}
	}
	tests := []struct {
		name       string
		method     string
		csrf       string
		wantCode   int
		wantCalled bool
	}{
		{"DELETE without header", http.MethodDelete, "", http.StatusForbidden, false},
		{"DELETE wrong token", http.MethodDelete, "bad", http.StatusForbidden, false},
		{"DELETE good token", http.MethodDelete, "good", http.StatusOK, true},
		{"GET needs no token", http.MethodGet, "", http.StatusOK, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := newFake()
			code, called, _ := run(f, "Protected", "token", tt.method, tt.csrf)
			if code != tt.wantCode || called != tt.wantCalled {
				t.Errorf("code = %d called = %v, want %d and %v", code, called, tt.wantCode, tt.wantCalled)
			}
			// A rejected request must not slide the idle window.
			if !tt.wantCalled && f.touched {
				t.Error("rejected request must not touch the session")
			}
		})
	}
}
