package auth

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/account"
	"golang.org/x/crypto/bcrypt"
)

const bcryptCost = 10

const tokenBytes = 32

var ErrInvalidCredentials = errors.New("invalid credentials")

var dummyHash = mustHash("dummy-password")

func mustHash(password string) []byte {
	h, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		// Only fails for passwords over 72 bytes; a startup panic is fine here.
		panic(err)
	}
	return h
}

type accountFinder interface {
	FindByEmployeeID(ctx context.Context, employeeID string) (*account.Account, error)
}

type sessionCreator interface {
	Create(ctx context.Context, tokenHash, csrfHash []byte, accountID int64) error
}

type Tokens struct {
	Session string
	CSRF    string
}

type Service struct {
	accounts accountFinder
	sessions sessionCreator
}

func newToken() (string, []byte, error) {
	raw := make([]byte, tokenBytes)
	if _, err := rand.Read(raw); err != nil {
		return "", nil, fmt.Errorf("generate token: %w", err)
	}
	token := base64.RawURLEncoding.EncodeToString(raw)
	sum := sha256.Sum256([]byte(token))
	return token, sum[:], nil
}

func NewService(accounts accountFinder, sessions sessionCreator) *Service {
	return &Service{accounts: accounts, sessions: sessions}
}

func (s *Service) Login(ctx context.Context, employeeID, password string) (Tokens, *account.Account, error) {
	hash := dummyHash
	acc, err := s.accounts.FindByEmployeeID(ctx, employeeID)
	switch {
	case err == nil:
		hash = []byte(acc.PasswordHash)
	case errors.Is(err, account.ErrNotFound):
		// Keep going with dummyHash.
	default:
		// A DB failure is not a credentials problem; do not disguise it.
		return Tokens{}, nil, fmt.Errorf("find account: %w", err)
	}

	cmpErr := bcrypt.CompareHashAndPassword(hash, []byte(password))
	if acc == nil || cmpErr != nil {
		return Tokens{}, nil, ErrInvalidCredentials
	}

	sessionToken, sessionHash, err := newToken()
	if err != nil {
		return Tokens{}, nil, err
	}
	csrfToken, csrfHash, err := newToken()
	if err != nil {
		return Tokens{}, nil, err
	}

	if err := s.sessions.Create(ctx, sessionHash, csrfHash, acc.ID); err != nil {
		return Tokens{}, nil, err
	}
	return Tokens{Session: sessionToken, CSRF: csrfToken}, acc, nil
}
