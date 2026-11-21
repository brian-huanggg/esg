package account

import (
	"context"
	"testing"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/api"
)

func TestCreateAccountStubReturn501(t *testing.T) {
	h := &Handler{}

	resp, err := h.CreateAccount(context.Background(), api.CreateAccountRequestObject{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	got, ok := resp.(api.CreateAccountdefaultJSONResponse)
	if !ok {
		t.Fatalf("unexpected response type %T", resp)
	}
	if got.StatusCode != 501 {
		t.Errorf("status = %d, want 501", got.StatusCode)
	}
}
