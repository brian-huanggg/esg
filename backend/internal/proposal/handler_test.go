package proposal

import (
	"context"
	"net/http"
	"testing"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/api"
)

func TestHandler_StubsReturn501(t *testing.T) {
	h := &ProposalHandler{}
	ctx := context.Background()

	created, _ := h.CreateProposal(ctx, api.CreateProposalRequestObject{})
	if r, ok := created.(api.CreateProposaldefaultJSONResponse); !ok || r.StatusCode != http.StatusNotImplemented {
		t.Fatalf("CreateProposal = %#v, want 501", created)
	}

	got, _ := h.GetProposal(ctx, api.GetProposalRequestObject{})
	if r, ok := got.(api.GetProposaldefaultJSONResponse); !ok || r.StatusCode != http.StatusNotImplemented {
		t.Fatalf("GetProposal = %#v, want 501", got)
	}
}
