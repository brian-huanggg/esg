package account

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost     = 10
	minPasswordLen = 8
	maxPasswordLen = 72
	maxNameLen     = 100
)

var employeeIDPattern = regexp.MustCompile(`^[0-9]{6}$`)

type FieldError struct {
	Field   string
	Message string
}

type ValidationError struct {
	Fields []FieldError
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %d field(s)", len(e.Fields))
}

type accountCreator interface {
	Create(ctx context.Context, a *Account) error
}

type Service struct {
	repo accountCreator
}

func NewService(repo accountCreator) *Service {
	return &Service{repo: repo}
}

func (s *Service) Register(ctx context.Context, employeeID, name, password string) (*Account, error) {
	name = strings.TrimSpace(name)
	if fields := validate(employeeID, name, password); len(fields) > 0 {
		return nil, &ValidationError{Fields: fields}
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	a := &Account{EmployeeID: employeeID, Name: name, PasswordHash: string(hash)}
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	return a, nil
}

func validate(employeeID, name, password string) []FieldError {
	var fs []FieldError
	if !employeeIDPattern.MatchString(employeeID) {
		fs = append(fs, FieldError{"employeeId", "must be 6 digits"})
	}
	if name == "" {
		fs = append(fs, FieldError{"name", "must not be empty"})
	} else if utf8.RuneCountInString(name) > maxNameLen {
		fs = append(fs, FieldError{"name", "must be at most 100 characters"})
	}
	if len(password) < minPasswordLen || len(password) > maxPasswordLen {
		fs = append(fs, FieldError{"password", "must be 8 to 72 bytes"})
	}
	return fs
}
