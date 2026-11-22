// Package proposal serves the /api/proposals endpoints.
package proposal

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/api"
	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/middleware"
)

type ProposalHandler struct {
	svc    *Service
	logger *slog.Logger
}

func NewHandler(svc *Service, logger *slog.Logger) *ProposalHandler {
	return &ProposalHandler{svc: svc, logger: logger}
}

func (h *ProposalHandler) CreateProposal(ctx context.Context, req api.CreateProposalRequestObject) (api.CreateProposalResponseObject, error) {
	leaderID, ok := middleware.AccountIDFrom(ctx)
	if !ok {
		// The Authenticate middleware always sets this; missing means a wiring bug.
		h.logger.ErrorContext(ctx, "account id missing from context")
		return h.createInternal(ctx), nil
	}
	p, err := h.svc.Create(ctx, leaderID, req.Body.Title, req.Body.Description)
	if err != nil {
		return h.createErrorResponse(ctx, err), nil
	}
	return api.CreateProposal201JSONResponse{
		Body: toAPI(p),
		Headers: api.CreateProposal201ResponseHeaders{
			ETag:     etag(p),
			Location: "/api/proposals/" + p.UUID.String(),
		},
	}, nil
}

func (h *ProposalHandler) GetProposal(ctx context.Context, req api.GetProposalRequestObject) (api.GetProposalResponseObject, error) {
	requestID := middleware.RequestIDFrom(ctx)
	leaderID, ok := middleware.AccountIDFrom(ctx)
	if !ok {
		h.logger.ErrorContext(ctx, "account id missing from context")
		return api.GetProposaldefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: internalError(requestID)}, nil
	}
	p, err := h.svc.Get(ctx, leaderID, req.Id)
	if errors.Is(err, ErrNotFound) {
		return api.GetProposal404JSONResponse{NotFoundJSONResponse: api.NotFoundJSONResponse{
			Code: "NOT_FOUND", Message: "proposal not found", RequestId: requestID,
		}}, nil
	}
	if err != nil {
		h.logger.ErrorContext(ctx, "get proposal failed", "error", err)
		return api.GetProposaldefaultJSONResponse{StatusCode: http.StatusInternalServerError, Body: internalError(requestID)}, nil
	}
	return api.GetProposal200JSONResponse{
		Body:    toAPI(p),
		Headers: api.GetProposal200ResponseHeaders{ETag: etag(p)},
	}, nil
}

func (h *ProposalHandler) createErrorResponse(ctx context.Context, err error) api.CreateProposalResponseObject {
	requestID := middleware.RequestIDFrom(ctx)

	var verr *ValidationError
	switch {
	case errors.As(err, &verr):
		fields := make([]api.FieldError, len(verr.Fields))
		for i, f := range verr.Fields {
			fields[i] = api.FieldError{Field: f.Field, Message: f.Message}
		}
		return api.CreateProposal422JSONResponse{
			Code: "VALIDATION_FAILED", Message: "validation failed", Fields: &fields, RequestId: requestID,
		}
	case errors.Is(err, ErrDuplicateTitle):
		return api.CreateProposal409JSONResponse{
			Code: "DUPLICATE_TITLE", Message: "you already have a proposal with this title", RequestId: requestID,
		}
	default:
		h.logger.ErrorContext(ctx, "create proposal failed", "error", err)
		return h.createInternal(ctx)
	}
}

func (h *ProposalHandler) createInternal(ctx context.Context) api.CreateProposalResponseObject {
	return api.CreateProposaldefaultJSONResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       internalError(middleware.RequestIDFrom(ctx)),
	}
}

func internalError(requestID string) api.Error {
	return api.Error{Code: "INTERNAL", Message: "internal error", RequestId: requestID}
}

// etag is a strong validator, so it carries quotes: "1".
func etag(p *Proposal) string {
	return fmt.Sprintf(`"%d"`, p.Version)
}

func toAPI(p *Proposal) api.Proposal {
	return api.Proposal{
		Id:          p.UUID,
		Title:       p.Title,
		Description: p.Description,
		CreatedAt:   p.CreatedAt,
		UpdatedAt:   p.UpdatedAt,
	}
}
