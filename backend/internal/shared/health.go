// Package shared holds handlers that do not belong to any business module.
package shared

import (
	"context"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/api"
)

// HealthHandler answers the liveness check defined in api/openapi.yaml
type HealthHandler struct{}

func (h *HealthHandler) GetHealth(_ context.Context, _ api.GetHealthRequestObject) (api.GetHealthResponseObject, error) {
	return api.GetHealth200JSONResponse{Status: api.Ok}, nil
}
