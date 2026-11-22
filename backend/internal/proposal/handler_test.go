package proposal

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"testing"

	"github.com/google/uuid"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/api"
	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/middleware"
)

type fakeStore struct {
	createErr error
	found     *Proposal
	findErr   error
}

func (f *fakeStore) Create(_ context.Context, p *Proposal) error {
	if f.createErr != nil {
		return f.createErr
	}
	p.UUID = uuid.New()
	p.Version = 1
	return nil
}

func (f *fakeStore) FindByUUID(_ context.Context, _ int64, _ uuid.UUID) (*Proposal, error) {
	return f.found, f.findErr
}

func newTestHandler(store *fakeStore) *ProposalHandler {
	logger := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewHandler(NewService(store), logger)
}

func authedCtx() context.Context {
	return middleware.WithAccountID(context.Background(), 7)
}

func TestCreateProposal_Created(t *testing.T) {
	h := newTestHandler(&fakeStore{})
	resp, _ := h.CreateProposal(authedCtx(), api.CreateProposalRequestObject{
		Body: &api.CreateProposalRequest{Title: " Solar roof ", Description: "panels"},
	})
	r, ok := resp.(api.CreateProposal201JSONResponse)
	if !ok {
		t.Fatalf("resp = %#v, want 201", resp)
	}
	if r.Body.Title != "Solar roof" {
		t.Errorf("title = %q, want trimmed", r.Body.Title)
	}
	if r.Headers.ETag != `"1"` {
		t.Errorf("ETag = %q, want \"1\" with quotes", r.Headers.ETag)
	}
	if want := "/api/proposals/" + r.Body.Id.String(); r.Headers.Location != want {
		t.Errorf("Location = %q, want %q", r.Headers.Location, want)
	}
}

func TestCreateProposal_Invalid(t *testing.T) {
	h := newTestHandler(&fakeStore{})
	resp, _ := h.CreateProposal(authedCtx(), api.CreateProposalRequestObject{
		Body: &api.CreateProposalRequest{Title: "   "},
	})
	r, ok := resp.(api.CreateProposal422JSONResponse)
	if !ok {
		t.Fatalf("resp = %#v, want 422", resp)
	}
	if r.Fields == nil || (*r.Fields)[0].Field != "title" {
		t.Errorf("fields = %+v, want title error", r.Fields)
	}
}

func TestCreateProposal_DuplicateTitle(t *testing.T) {
	h := newTestHandler(&fakeStore{createErr: ErrDuplicateTitle})
	resp, _ := h.CreateProposal(authedCtx(), api.CreateProposalRequestObject{
		Body: &api.CreateProposalRequest{Title: "Solar roof"},
	})
	if _, ok := resp.(api.CreateProposal409JSONResponse); !ok {
		t.Fatalf("resp = %#v, want 409", resp)
	}
}

func TestCreateProposal_InternalError(t *testing.T) {
	h := newTestHandler(&fakeStore{createErr: errors.New("db down")})
	resp, _ := h.CreateProposal(authedCtx(), api.CreateProposalRequestObject{
		Body: &api.CreateProposalRequest{Title: "Solar roof"},
	})
	r, ok := resp.(api.CreateProposaldefaultJSONResponse)
	if !ok || r.StatusCode != http.StatusInternalServerError {
		t.Fatalf("resp = %#v, want 500", resp)
	}
}

func TestGetProposal_Found(t *testing.T) {
	p := &Proposal{UUID: uuid.New(), Title: "Solar roof", Version: 3}
	h := newTestHandler(&fakeStore{found: p})
	resp, _ := h.GetProposal(authedCtx(), api.GetProposalRequestObject{Id: p.UUID})
	r, ok := resp.(api.GetProposal200JSONResponse)
	if !ok {
		t.Fatalf("resp = %#v, want 200", resp)
	}
	if r.Headers.ETag != `"3"` {
		t.Errorf("ETag = %q, want \"3\"", r.Headers.ETag)
	}
}

func TestGetProposal_NotFound(t *testing.T) {
	h := newTestHandler(&fakeStore{findErr: ErrNotFound})
	resp, _ := h.GetProposal(authedCtx(), api.GetProposalRequestObject{Id: uuid.New()})
	if _, ok := resp.(api.GetProposal404JSONResponse); !ok {
		t.Fatalf("resp = %#v, want 404", resp)
	}
}
