package proposal

import (
	"context"
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/google/uuid"
)

const (
	maxTitleLen       = 100
	maxDescriptionLen = 2000
)

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

type proposalStore interface {
	Create(ctx context.Context, p *Proposal) error
	FindByUUID(ctx context.Context, leaderID int64, id uuid.UUID) (*Proposal, error)
}

type Service struct {
	repo proposalStore
}

func NewService(repo proposalStore) *Service {
	return &Service{repo: repo}
}

func (s *Service) Create(ctx context.Context, leaderID int64, title, description string) (*Proposal, error) {
	title = strings.TrimSpace(title)
	if fields := validate(title, description); len(fields) > 0 {
		return nil, &ValidationError{Fields: fields}
	}
	p := &Proposal{LeaderAccountID: leaderID, Title: title, Description: description}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *Service) Get(ctx context.Context, leaderID int64, id uuid.UUID) (*Proposal, error) {
	return s.repo.FindByUUID(ctx, leaderID, id)
}

func validate(title, description string) []FieldError {
	var fs []FieldError
	if title == "" {
		fs = append(fs, FieldError{"title", "must not be empty"})
	} else if utf8.RuneCountInString(title) > maxTitleLen {
		fs = append(fs, FieldError{"title", "must be at most 100 characters"})
	}
	if utf8.RuneCountInString(description) > maxDescriptionLen {
		fs = append(fs, FieldError{"description", "must be at most 2000 characters"})
	}
	return fs
}
