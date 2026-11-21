package account

import (
	"context"
	"errors"
	"log/slog"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/api"
	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/middleware"
)

type Handler struct {
	svc    *Service
	logger *slog.Logger
}

func NewHandler(svc *Service, logger *slog.Logger) *Handler {
	return &Handler{svc: svc, logger: logger}
}

func (h *Handler) CreateAccount(ctx context.Context, req api.CreateAccountRequestObject) (api.CreateAccountResponseObject, error) {
	a, err := h.svc.Register(ctx, req.Body.EmployeeId, req.Body.Name, req.Body.Password)
	if err != nil {
		return h.errorResponse(ctx, err), nil
	}
	return api.CreateAccount201JSONResponse{
		Id:         a.UUID,
		EmployeeId: a.EmployeeID,
		Name:       a.Name,
		CreatedAt:  a.CreatedAt,
	}, nil
}

func (h *Handler) errorResponse(ctx context.Context, err error) api.CreateAccountResponseObject {
	requestID := middleware.RequestIDFrom(ctx)

	var verr *ValidationError
	switch {
	case errors.As(err, &verr):
		fields := make([]api.FieldError, len(verr.Fields))
		for i, f := range verr.Fields {
			fields[i] = api.FieldError{Field: f.Field, Message: f.Message}
		}
		return api.CreateAccount422JSONResponse{
			Code: "VALIDATION_FAILED", Message: "validation failed", Fields: &fields, RequestId: requestID,
		}
	case errors.Is(err, ErrEmployeeIDTaken):
		return api.CreateAccount409JSONResponse{
			Code: "EMPLOYEE_ID_TAKEN", Message: "employee id already taken", RequestId: requestID,
		}
	default:
		// Log the real cause; never leak it to the client.
		h.logger.ErrorContext(ctx, "create account failed", "error", err)
		return api.CreateAccountdefaultJSONResponse{
			StatusCode: 500,
			Body:       api.Error{Code: "INTERNAL", Message: "internal error", RequestId: requestID},
		}
	}
}
