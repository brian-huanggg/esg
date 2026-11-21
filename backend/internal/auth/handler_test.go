package auth

import (
	"context"
	"errors"
	"io"
	"log/slog"
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
	token string
	acc   *account.Account
	err   error
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

func (f fakeLogin) Login(context.Context, string, string) (string, *account.Account, error) {
	return f.token, f.acc, f.err
}

func loginReq(employeeID, password string) api.CreateSessionRequestObject {
	return api.CreateSessionRequestObject{
		Body: &api.CreateSessionRequest{EmployeeId: employeeID, Password: password},
	}
}

func TestCreateSessionSuccessSetsCookie(t *testing.T) {
	h := newTestHandler(fakeLogin{token: "tok123", acc: &account.Account{EmployeeID: "123456", Name: "Amy"}})

	resp, _ := h.CreateSession(context.Background(), loginReq("123456", "password1"))

	r, ok := resp.(api.CreateSession201JSONResponse)
	if !ok {
		t.Fatalf("resp = %#v, want 201", resp)
	}
	if r.Body.EmployeeId != "123456" {
		t.Errorf("body employeeId = %q", r.Body.EmployeeId)
	}
	cookie := r.Headers.SetCookie
	for _, want := range []string{"session=tok123", "Path=/", "HttpOnly", "Secure", "SameSite=Lax"} {
		if !strings.Contains(cookie, want) {
			t.Errorf("cookie %q missing %q", cookie, want)
		}
	}
	if strings.Contains(cookie, "Max-Age") {
		t.Errorf("cookie %q must not set Max-Age", cookie)
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
			case api.DeleteCurrentSession204Response:
				if tt.wantType != "204" || !strings.Contains(r.Headers.SetCookie, "Max-Age=0") {
					t.Errorf("got 204 with cookie %q", r.Headers.SetCookie)
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
