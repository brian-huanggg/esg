package account

import (
	"context"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/api"
)

type Handler struct{}

func (h *Handler) CreateAccount(ctx context.Context, _ api.CreateAccountRequestObject) (api.CreateAccountResponseObject, error) {
	return api.CreateAccountdefaultJSONResponse{
		StatusCode: 501,
		Body:       api.Error{Code: "NOT_IMPLEMENTED", Message: "not implemented"},
	}, nil
}
