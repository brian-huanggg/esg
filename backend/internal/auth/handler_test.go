package auth

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/account"
	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/api"
	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/middleware"
)

type fakeAccountByID struct {
	acc *account.Account
	err error
}

func (f fakeAccountByID) FindByID(context.Context, int64) (*account.Account, error) {
	return f.acc, f.err
}

type fakeLogin struct {
	tokens Tokens
	acc    *account.Account
	err    error
}

type fakeSessionDeleter struct {
	err     error
	deleted []byte
}

func (f *fakeSessionDeleter) Delete(_ context.Context, h []byte) error {
	f.deleted = h
	return f.err
}

func newTestHandler(f fakeLogin) *SessionHandler {
	return NewSessionHandler(f, fakeAccountByID{}, &fakeSessionDeleter{}, slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func (f fakeLogin) Login(context.Context, string, string) (Tokens, *account.Account, error) {
	return f.tokens, f.acc, f.err
}

func loginReq(employeeID, password string) api.CreateSessionRequestObject {
	return api.CreateSessionRequestObject{
		Body: &api.CreateSessionRequest{EmployeeId: employeeID, Password: password},
	}
}

func TestCreateSessionSuccessSetsCookies(t *testing.T) {
	h := newTestHandler(fakeLogin{
		tokens: Tokens{Session: "tok123", CSRF: "csrf456"},
		acc:    &account.Account{EmployeeID: "123456", Name: "Amy"},
	})

	resp, _ := h.CreateSession(context.Background(), loginReq("123456", "password1"))

	r, ok := resp.(createSessionResponse)
	if !ok {
		t.Fatalf("resp = %#v, want createSessionResponse", resp)
	}
	rec := httptest.NewRecorder()
	if err := r.VisitCreateSessionResponse(rec); err != nil {
		t.Fatal(err)
	}
	if rec.Code != 201 {
		t.Errorf("status = %d, want 201", rec.Code)
	}

	cookies := rec.Header().Values("Set-Cookie")
	if len(cookies) != 2 {
		t.Fatalf("got %d Set-Cookie headers, want 2: %q", len(cookies), cookies)
	}
	wantSession := []string{"session=tok123", "Path=/", "HttpOnly", "Secure", "SameSite=Lax"}
	for _, want := range wantSession {
		if !strings.Contains(cookies[0], want) {
			t.Errorf("session cookie %q missing %q", cookies[0], want)
		}
	}
	for _, want := range []string{"csrf_token=csrf456", "Path=/", "Secure", "SameSite=Lax"} {
		if !strings.Contains(cookies[1], want) {
			t.Errorf("csrf cookie %q missing %q", cookies[1], want)
		}
	}
	if strings.Contains(cookies[1], "HttpOnly") {
		t.Errorf("csrf cookie %q must be readable by JS", cookies[1])
	}
}

func TestCreateSessionErrors(t *testing.T) {
	tests := []struct {
		name     string
		login    fakeLogin
		employee string
		password string
		check    func(api.CreateSessionResponseObject) bool
	}{
		{"wrong credentials", fakeLogin{err: ErrInvalidCredentials}, "123456", "bad",
			func(r api.CreateSessionResponseObject) bool { _, ok := r.(api.CreateSession401JSONResponse); return ok }},
		{"empty password", fakeLogin{}, "123456", "",
			func(r api.CreateSessionResponseObject) bool { _, ok := r.(api.CreateSession422JSONResponse); return ok }},
		{"db failure", fakeLogin{err: errors.New("db down")}, "123456", "password1",
			func(r api.CreateSessionResponseObject) bool {
				d, ok := r.(api.CreateSessiondefaultJSONResponse)
				return ok && d.StatusCode == 500
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, _ := newTestHandler(tt.login).CreateSession(context.Background(), loginReq(tt.employee, tt.password))
			if !tt.check(resp) {
				t.Errorf("unexpected response %#v", resp)
			}
		})
	}
}

func TestGetCurrentAccount(t *testing.T) {
	found := &account.Account{EmployeeID: "123456", Name: "Amy"}
	tests := []struct {
		name  string
		ctx   context.Context
		accts fakeAccountByID
		check func(api.GetCurrentAccountResponseObject) bool
	}{
		{"found", middleware.WithAccountID(context.Background(), 7), fakeAccountByID{acc: found},
			func(r api.GetCurrentAccountResponseObject) bool {
				ok200, ok := r.(api.GetCurrentAccount200JSONResponse)
				return ok && ok200.EmployeeId == "123456"
			}},
		{"no account id in ctx", context.Background(), fakeAccountByID{},
			func(r api.GetCurrentAccountResponseObject) bool {
				_, ok := r.(api.GetCurrentAccount401JSONResponse)
				return ok
			}},
		{"account deleted", middleware.WithAccountID(context.Background(), 7), fakeAccountByID{err: account.ErrNotFound},
			func(r api.GetCurrentAccountResponseObject) bool {
				_, ok := r.(api.GetCurrentAccount401JSONResponse)
				return ok
			}},
		{"db failure", middleware.WithAccountID(context.Background(), 7), fakeAccountByID{err: errors.New("db down")},
			func(r api.GetCurrentAccountResponseObject) bool {
				d, ok := r.(api.GetCurrentAccountdefaultJSONResponse)
				return ok && d.StatusCode == 500
			}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := NewSessionHandler(fakeLogin{}, tt.accts, &fakeSessionDeleter{}, slog.New(slog.NewTextHandler(io.Discard, nil)))
			resp, _ := h.GetCurrentAccount(tt.ctx, api.GetCurrentAccountRequestObject{})
			if !tt.check(resp) {
				t.Errorf("unexpected response %#v", resp)
			}
		})
	}
}

func TestDeleteCurrentSession(t *testing.T) {
	hash := []byte("hash")
	tests := []struct {
		name     string
		ctx      context.Context
		err      error
		wantType string
	}{
		{"ok", middleware.WithTokenHash(context.Background(), hash), nil, "204"},
		{"no token hash in ctx", context.Background(), nil, "401"},
		{"db failure", middleware.WithTokenHash(context.Background(), hash), errors.New("db down"), "500"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fs := &fakeSessionDeleter{err: tt.err}
			h := NewSessionHandler(fakeLogin{}, fakeAccountByID{}, fs, slog.New(slog.NewTextHandler(io.Discard, nil)))

			resp, err := h.DeleteCurrentSession(tt.ctx, api.DeleteCurrentSessionRequestObject{})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			switch r := resp.(type) {
			case deleteSessionResponse:
				if tt.wantType != "204" {
					t.Fatalf("unexpected 204")
				}
				rec := httptest.NewRecorder()
				if err := r.VisitDeleteCurrentSessionResponse(rec); err != nil {
					t.Fatal(err)
				}
				cookies := rec.Result().Cookies()
				if rec.Code != http.StatusNoContent || len(cookies) != 2 {
					t.Fatalf("code = %d, cookies = %d, want 204 and 2", rec.Code, len(cookies))
				}
				for _, c := range cookies {
					if c.MaxAge != -1 { // net/http parses Max-Age=0 as -1
						t.Errorf("cookie %s MaxAge = %d, want expired", c.Name, c.MaxAge)
					}
				}
				if string(fs.deleted) != "hash" {
					t.Errorf("deleted %q, want hash", fs.deleted)
				}
			case api.DeleteCurrentSession401JSONResponse:
				if tt.wantType != "401" {
					t.Errorf("unexpected 401")
				}
			case api.DeleteCurrentSessiondefaultJSONResponse:
				if tt.wantType != "500" || r.StatusCode != 500 {
					t.Errorf("unexpected default response %d", r.StatusCode)
				}
			default:
				t.Fatalf("unexpected response type %T", resp)
			}
		})
	}
}
