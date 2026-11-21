package auth

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/account"
	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/api"
	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/middleware"
)

const cookieName = middleware.SessionCookieName

type loginService interface {
	Login(ctx context.Context, employeeID, password string) (Tokens, *account.Account, error)
}

type createSessionResponse struct {
	api.CreateSession201JSONResponse
	tokens Tokens
}

type deleteSessionResponse struct{}

func (deleteSessionResponse) VisitDeleteCurrentSessionResponse(w http.ResponseWriter) error {
	w.Header().Add("Set-Cookie", clearedSessionCookie())
	w.Header().Add("Set-Cookie", clearedCSRFCookie())
	w.WriteHeader(http.StatusNoContent)
	return nil
}

func (r createSessionResponse) VisitCreateSessionResponse(w http.ResponseWriter) error {
	w.Header().Add("Set-Cookie", sessionCookie(r.tokens.Session))
	w.Header().Add("Set-Cookie", csrfCookie(r.tokens.CSRF))
	return r.CreateSession201JSONResponse.VisitCreateSessionResponse(w)
}

func csrfCookie(token string) string {
	c := http.Cookie{
		Name:     "csrf_token",
		Value:    token,
		Path:     "/",
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	return c.String()
}

type accountByIDFinder interface {
	FindByID(ctx context.Context, id int64) (*account.Account, error)
}

type sessionDeleter interface {
	Delete(ctx context.Context, tokenHash []byte) error
}

type SessionHandler struct {
	svc      loginService
	accounts accountByIDFinder
	sessions sessionDeleter
	logger   *slog.Logger
}

func NewSessionHandler(svc loginService, accounts accountByIDFinder, sessions sessionDeleter, logger *slog.Logger) *SessionHandler {
	return &SessionHandler{svc: svc, accounts: accounts, sessions: sessions, logger: logger}
}

func sessionCookie(token string) string {
	c := http.Cookie{
		Name:     cookieName,
		Value:    token,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	return c.String()
}

func clearedSessionCookie() string {
	c := http.Cookie{
		Name:     cookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1, // serialized as Max-Age=0
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	return c.String()
}

func clearedCSRFCookie() string {
	c := http.Cookie{
		Name:     "csrf_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
	}
	return c.String()
}

func (h *SessionHandler) CreateSession(ctx context.Context, req api.CreateSessionRequestObject) (api.CreateSessionResponseObject, error) {
	requestID := middleware.RequestIDFrom(ctx)

	if req.Body.EmployeeId == "" || req.Body.Password == "" {
		return api.CreateSession422JSONResponse{
			Code: "VALIDATION_FAILED", Message: "validation failed", RequestId: requestID,
		}, nil
	}

	tokens, acc, err := h.svc.Login(ctx, req.Body.EmployeeId, req.Body.Password)
	switch {
	case err == nil:
		return createSessionResponse{
			CreateSession201JSONResponse: api.CreateSession201JSONResponse(api.Account{
				Id: acc.UUID, EmployeeId: acc.EmployeeID, Name: acc.Name, CreatedAt: acc.CreatedAt,
			}),
			tokens: tokens,
		}, nil
	case errors.Is(err, ErrInvalidCredentials):
		return api.CreateSession401JSONResponse{
			Code: "INVALID_CREDENTIALS", Message: "wrong employee id or password", RequestId: requestID,
		}, nil
	default:
		h.logger.ErrorContext(ctx, "login failed", "error", err)
		return api.CreateSessiondefaultJSONResponse{
			StatusCode: 500,
			Body:       api.Error{Code: "INTERNAL", Message: "internal error", RequestId: requestID},
		}, nil
	}
}

func (h *SessionHandler) DeleteCurrentSession(ctx context.Context, _ api.DeleteCurrentSessionRequestObject) (api.DeleteCurrentSessionResponseObject, error) {
	requestID := middleware.RequestIDFrom(ctx)

	hash, ok := middleware.TokenHashFrom(ctx)
	if !ok {
		return api.DeleteCurrentSession401JSONResponse{UnauthorizedJSONResponse: api.UnauthorizedJSONResponse{
			Code: "UNAUTHORIZED", Message: "login required", RequestId: requestID,
		}}, nil
	}

	if err := h.sessions.Delete(ctx, hash); err != nil {
		h.logger.ErrorContext(ctx, "delete session failed", "error", err)
		return api.DeleteCurrentSessiondefaultJSONResponse{
			StatusCode: 500,
			Body:       api.Error{Code: "INTERNAL", Message: "internal error", RequestId: requestID},
		}, nil
	}
	return deleteSessionResponse{}, nil
}

func (h *SessionHandler) GetCurrentAccount(ctx context.Context, _ api.GetCurrentAccountRequestObject) (api.GetCurrentAccountResponseObject, error) {
	requestID := middleware.RequestIDFrom(ctx)
	unauthorized := api.GetCurrentAccount401JSONResponse{UnauthorizedJSONResponse: api.UnauthorizedJSONResponse{
		Code: "UNAUTHORIZED", Message: "login required", RequestId: requestID,
	}}

	id, ok := middleware.AccountIDFrom(ctx)
	if !ok {
		return unauthorized, nil
	}

	acc, err := h.accounts.FindByID(ctx, id)
	switch {
	case err == nil:
		return api.GetCurrentAccount200JSONResponse{
			Id: acc.UUID, EmployeeId: acc.EmployeeID, Name: acc.Name, CreatedAt: acc.CreatedAt,
		}, nil
	case errors.Is(err, account.ErrNotFound):
		return unauthorized, nil
	default:
		h.logger.ErrorContext(ctx, "load current account failed", "error", err)
		return api.GetCurrentAccountdefaultJSONResponse{
			StatusCode: 500,
			Body:       api.Error{Code: "INTERNAL", Message: "internal error", RequestId: requestID},
		}, nil
	}
}
