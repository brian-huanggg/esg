package account

import (
	"context"
	"errors"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

type fakeRepo struct {
	called bool
	got    *Account
	err    error
}

func (f *fakeRepo) Create(_ context.Context, a *Account) error {
	f.called = true
	f.got = a
	return f.err
}

func TestRegister_Invalid(t *testing.T) {
	tests := []struct {
		name, id, uname, pw, field string
	}{
		{"short id", "12345", "Alice", "password1", "employeeId"},
		{"letters in id", "12345a", "Alice", "password1", "employeeId"},
		{"blank name", "123456", "   ", "password1", "name"},
		{"long name", "123456", strings.Repeat("a", 101), "password1", "name"},
		{"short password", "123456", "Alice", "1234567", "password"},
		{"long password", "123456", "Alice", strings.Repeat("a", 73), "password"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repo := &fakeRepo{}
			_, err := NewService(repo).Register(context.Background(), tt.id, tt.uname, tt.pw)

			var ve *ValidationError
			if !errors.As(err, &ve) {
				t.Fatalf("err = %v, want *ValidationError", err)
			}
			if ve.Fields[0].Field != tt.field {
				t.Errorf("field = %q, want %q", ve.Fields[0].Field, tt.field)
			}
			if repo.called {
				t.Error("repo must not be called on invalid input")
			}
		})
	}
}

func TestRegister_HashesPassword(t *testing.T) {
	repo := &fakeRepo{}
	a, err := NewService(repo).Register(context.Background(), "123456", " Alice ", "password1")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(a.PasswordHash, "$2a$10$") {
		t.Errorf("hash = %q, want $2a$10$ prefix", a.PasswordHash)
	}
	if err := bcrypt.CompareHashAndPassword([]byte(a.PasswordHash), []byte("password1")); err != nil {
		t.Errorf("hash does not match password: %v", err)
	}
	if repo.got.Name != "Alice" {
		t.Errorf("name = %q, want trimmed Alice", repo.got.Name)
	}
}

func TestRegister_DuplicatePassesThrough(t *testing.T) {
	repo := &fakeRepo{err: ErrEmployeeIDTaken}
	_, err := NewService(repo).Register(context.Background(), "123456", "Alice", "password1")
	if !errors.Is(err, ErrEmployeeIDTaken) {
		t.Fatalf("err = %v, want ErrEmployeeIDTaken", err)
	}
}
