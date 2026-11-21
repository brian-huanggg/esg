package auth

import (
	"context"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/api"
)

// SessionHandler is named explicitly: embedding it next to account.Handler in
// main.go would give two fields called "Handler", which does not compile.
type SessionHandler struct{}

func notImplemented() api.Error {
	return api.Error{Code: "NOT_IMPLEMENTED", Message: "not implemented"}
}

func (h *SessionHandler) CreateSession(_ context.Context, _ api.CreateSessionRequestObject) (api.CreateSessionResponseObject, error) {
	return api.CreateSessiondefaultJSONResponse{StatusCode: 501, Body: notImplemented()}, nil
}

func (h *SessionHandler) DeleteCurrentSession(_ context.Context, _ api.DeleteCurrentSessionRequestObject) (api.DeleteCurrentSessionResponseObject, error) {
	return api.DeleteCurrentSessiondefaultJSONResponse{StatusCode: 501, Body: notImplemented()}, nil
}

func (h *SessionHandler) GetCurrentAccount(_ context.Context, _ api.GetCurrentAccountRequestObject) (api.GetCurrentAccountResponseObject, error) {
	return api.GetCurrentAccountdefaultJSONResponse{StatusCode: 501, Body: notImplemented()}, nil
}
