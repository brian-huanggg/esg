package middleware

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/session"
	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/api"
)

// SessionCookieName is shared with the login handler that sets the cookie.
const SessionCookieName = "session"

type sessionFinder interface {
	FindActive(ctx context.Context, tokenHash []byte) (int64, error)
	Touch(ctx context.Context, tokenHash []byte) error
}

// Distinct unexported types keep these keys from colliding with ctxKey{}.
type accountIDKey struct{}
type tokenHashKey struct{}

// AccountIDFrom returns the logged-in account's BIGINT id; ok is false on
// routes the middleware did not protect.
func AccountIDFrom(ctx context.Context) (int64, bool) {
	id, ok := ctx.Value(accountIDKey{}).(int64)
	return id, ok
}

// TokenHashFrom lets logout delete the current session (Step 3.5b).
func TokenHashFrom(ctx context.Context) ([]byte, bool) {
	h, ok := ctx.Value(tokenHashKey{}).([]byte)
	return h, ok
}

// Authenticate guards only the operations listed in protected.
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

			accountID, err := sessions.FindActive(ctx, hash)
			if errors.Is(err, session.ErrNotFound) {
				writeError(ctx, w, http.StatusUnauthorized, "UNAUTHORIZED", "login required")
				return nil, nil
			}
			if err != nil {
				logger.ErrorContext(ctx, "find session failed", "error", err)
				writeError(ctx, w, http.StatusInternalServerError, "INTERNAL", "internal error")
				return nil, nil
			}
			if err := sessions.Touch(ctx, hash); err != nil {
				logger.ErrorContext(ctx, "touch session failed", "error", err)
			}

			ctx = context.WithValue(ctx, accountIDKey{}, accountID)
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
