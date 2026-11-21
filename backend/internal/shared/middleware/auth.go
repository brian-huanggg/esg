package middleware

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/session"
	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/api"
)

const SessionCookieName = "session"

type sessionFinder interface {
	FindActive(ctx context.Context, tokenHash []byte) (session.Active, error)
	Touch(ctx context.Context, tokenHash []byte) error
}

type accountIDKey struct{}
type tokenHashKey struct{}

func AccountIDFrom(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(accountIDKey{}).(int64)
	return id, ok
}

func TokenHashFrom(ctx context.Context) ([]byte, bool) {
	h, ok := ctx.Value(tokenHashKey{}).([]byte)
	return h, ok
}

const csrfHeader = "X-CSRF-Token"

func isSafeMethod(m string) bool {
	return m == http.MethodGet || m == http.MethodHead || m == http.MethodOptions
}

func csrfValid(r *http.Request, wantHash []byte) bool {
	token := r.Header.Get(csrfHeader)
	if token == "" {
		return false
	}
	sum := sha256.Sum256([]byte(token))
	// Constant-time compare so response timing leaks nothing about the stored hash.
	return subtle.ConstantTimeCompare(sum[:], wantHash) == 1
}

func Authenticate(sessions sessionFinder, logger *slog.Logger, protected map[string]bool) api.StrictMiddlewareFunc {
	return func(next api.StrictHandlerFunc, operationID string) api.StrictHandlerFunc {
		if !protected[operationID] {
			return next
		}
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, req any) (any, error) {
			cookie, err := r.Cookie(SessionCookieName)
			if err != nil {
				writeError(ctx, w, http.StatusUnauthorized, "UNAUTHORIZED", "login required")
				return nil, nil
			}
			sum := sha256.Sum256([]byte(cookie.Value))
			hash := sum[:]

			active, err := sessions.FindActive(ctx, hash)
			if errors.Is(err, session.ErrNotFound) {
				writeError(ctx, w, http.StatusUnauthorized, "UNAUTHORIZED", "login required")
				return nil, nil
			}
			if err != nil {
				logger.ErrorContext(ctx, "find session failed", "error", err)
				writeError(ctx, w, http.StatusInternalServerError, "INTERNAL", "internal error")
				return nil, nil
			}

			// Order matters (PRD 5.2): 401 first, then CSRF.
			if !isSafeMethod(r.Method) && !csrfValid(r, active.CSRFTokenHash) {
				writeError(ctx, w, http.StatusForbidden, "CSRF_FAILED", "invalid csrf token")
				return nil, nil
			}

			// Touch only after CSRF passes, so a forged request cannot keep a session alive.
			if err := sessions.Touch(ctx, hash); err != nil {
				logger.ErrorContext(ctx, "touch session failed", "error", err)
			}

			ctx = context.WithValue(ctx, accountIDKey{}, active.AccountID)
			ctx = context.WithValue(ctx, tokenHashKey{}, hash)
			return next(ctx, w, r, req)
		}
	}
}

func writeError(ctx context.Context, w http.ResponseWriter, status int, code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(api.Error{Code: code, Message: msg, RequestId: RequestIDFrom(ctx)})
}

func WithAccountID(ctx context.Context, id int64) context.Context {
	return context.WithValue(ctx, accountIDKey{}, id)
}

func WithTokenHash(ctx context.Context, h []byte) context.Context {
	return context.WithValue(ctx, tokenHashKey{}, h)
}
