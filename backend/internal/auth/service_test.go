package auth

import (
	"context"
	"crypto/sha256"
	"errors"
	"testing"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/account"
	"golang.org/x/crypto/bcrypt"
)

type fakeAccounts struct {
	acc *account.Account
	err error
}

func (f *fakeAccounts) FindByEmployeeID(_ context.Context, _ string) (*account.Account, error) {
	return f.acc, f.err
}

type fakeSessions struct {
	called    bool
	tokenHash []byte
	accountID int64
}

func (f *fakeSessions) Create(_ context.Context, tokenHash []byte, accountID int64) error {
	f.called = true
	f.tokenHash = tokenHash
	f.accountID = accountID
	return nil
}

func newAccount(t *testing.T) *account.Account {
	t.Helper()
	h, err := bcrypt.GenerateFromPassword([]byte("password1"), bcrypt.MinCost)
	if err != nil {
		t.Fatal(err)
	}
	return &account.Account{ID: 7, EmployeeID: "123456", PasswordHash: string(h)}
}

func TestLogin_Success(t *testing.T) {
	sessions := &fakeSessions{}
	svc := NewService(&fakeAccounts{acc: newAccount(t)}, sessions)

	token, _, err := svc.Login(context.Background(), "123456", "password1")
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte(token))
	if string(sessions.tokenHash) != string(sum[:]) {
		t.Error("stored hash is not SHA-256 of the returned token")
	}
	if sessions.accountID != 7 {
		t.Errorf("accountID = %d, want 7", sessions.accountID)
	}
}

func TestLogin_Failures(t *testing.T) {
	tests := []struct {
		name     string
		accounts *fakeAccounts
		password string
	}{
		{"wrong password", &fakeAccounts{acc: newAccount(t)}, "wrong-password"},
		{"unknown employee", &fakeAccounts{err: account.ErrNotFound}, "password1"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sessions := &fakeSessions{}
			_, _, err := NewService(tt.accounts, sessions).Login(context.Background(), "123456", tt.password)
			if !errors.Is(err, ErrInvalidCredentials) {
				t.Fatalf("err = %v, want ErrInvalidCredentials", err)
			}
			if sessions.called {
				t.Error("session must not be created on failure")
			}
		})
	}
}

func TestLogin_DBErrorIsNotCredentialsError(t *testing.T) {
	boom := errors.New("db down")
	svc := NewService(&fakeAccounts{err: boom}, &fakeSessions{})

	_, _, err := svc.Login(context.Background(), "123456", "password1")
	if !errors.Is(err, boom) || errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("err = %v, want wrapped db error", err)
	}
}
