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
	Login(ctx context.Context, employeeID, password string) (string, *account.Account, error)
}

type accountByIDFinder interface {
	FindByID(ctx context.Context, id int64) (*account.Account, error)
}

type SessionHandler struct {
	svc      loginService
	accounts accountByIDFinder
	logger   *slog.Logger
}

func NewSessionHandler(svc loginService, accounts accountByIDFinder, logger *slog.Logger) *SessionHandler {
	return &SessionHandler{svc: svc, accounts: accounts, logger: logger}
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

func notImplemented() api.Error {
	return api.Error{Code: "NOT_IMPLEMENTED", Message: "not implemented"}
}

func (h *SessionHandler) CreateSession(ctx context.Context, req api.CreateSessionRequestObject) (api.CreateSessionResponseObject, error) {
	requestID := middleware.RequestIDFrom(ctx)

	if req.Body.EmployeeId == "" || req.Body.Password == "" {
		return api.CreateSession422JSONResponse{
			Code: "VALIDATION_FAILED", Message: "validation failed", RequestId: requestID,
		}, nil
	}

	token, acc, err := h.svc.Login(ctx, req.Body.EmployeeId, req.Body.Password)
	switch {
	case err == nil:
		return api.CreateSession201JSONResponse{
			Body: api.Account{
				Id: acc.UUID, EmployeeId: acc.EmployeeID, Name: acc.Name, CreatedAt: acc.CreatedAt,
			},
			Headers: api.CreateSession201ResponseHeaders{SetCookie: sessionCookie(token)},
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

func (h *SessionHandler) DeleteCurrentSession(_ context.Context, _ api.DeleteCurrentSessionRequestObject) (api.DeleteCurrentSessionResponseObject, error) {
	return api.DeleteCurrentSessiondefaultJSONResponse{StatusCode: 501, Body: notImplemented()}, nil
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
