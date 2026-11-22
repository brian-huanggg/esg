// Package proposal serves the /api/proposals endpoints.
package proposal

import (
	"context"
	"net/http"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/api"
	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/middleware"
)

// ProposalHandler is a 501 stub so the server satisfies the generated interface
// before the real implementation lands.
type ProposalHandler struct{}

func notImplemented(ctx context.Context) api.Error {
	return api.Error{
		Code:      "NOT_IMPLEMENTED",
		Message:   "not implemented",
		RequestId: middleware.RequestIDFrom(ctx),
	}
}

func (h *ProposalHandler) CreateProposal(ctx context.Context, _ api.CreateProposalRequestObject) (api.CreateProposalResponseObject, error) {
	// 501 goes through `default` so it never has to appear in the spec.
	return api.CreateProposaldefaultJSONResponse{StatusCode: http.StatusNotImplemented, Body: notImplemented(ctx)}, nil
}

func (h *ProposalHandler) GetProposal(ctx context.Context, _ api.GetProposalRequestObject) (api.GetProposalResponseObject, error) {
	return api.GetProposaldefaultJSONResponse{StatusCode: http.StatusNotImplemented, Body: notImplemented(ctx)}, nil
}
