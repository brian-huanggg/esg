package auth

import (
	"context"
	"testing"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/api"
)

func TestStubsReturn501(t *testing.T) {
	h := &SessionHandler{}
	ctx := context.Background()

	create, _ := h.CreateSession(ctx, api.CreateSessionRequestObject{})
	if r, ok := create.(api.CreateSessiondefaultJSONResponse); !ok || r.StatusCode != 501 {
		t.Errorf("CreateSession = %#v, want 501 default response", create)
	}

	del, _ := h.DeleteCurrentSession(ctx, api.DeleteCurrentSessionRequestObject{})
	if r, ok := del.(api.DeleteCurrentSessiondefaultJSONResponse); !ok || r.StatusCode != 501 {
		t.Errorf("DeleteCurrentSession = %#v, want 501 default response", del)
	}

	me, _ := h.GetCurrentAccount(ctx, api.GetCurrentAccountRequestObject{})
	if r, ok := me.(api.GetCurrentAccountdefaultJSONResponse); !ok || r.StatusCode != 501 {
		t.Errorf("GetCurrentAccount = %#v, want 501 default response", me)
	}
}
