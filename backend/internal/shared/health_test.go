package shared

import (
	"context"
	"testing"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/api"
)

func TestHealthHandlerGetHealth(t *testing.T) {
	h := &HealthHandler{}

	resp, err := h.GetHealth(context.Background(), api.GetHealthRequestObject{})
	if err != nil {
		t.Fatalf("GetHealth returned error: %v", err)
	}

	got, ok := resp.(api.GetHealth200JSONResponse)
	if !ok {
		t.Fatalf("response type = %T, want api.GetHealth200JSONResponse", resp)
	}
	if got.Status != api.Ok {
		t.Errorf("Status = %q, want %q", got.Status, api.Ok)
	}
}
