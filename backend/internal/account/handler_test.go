package account

import (
	"context"
	"errors"
	"log/slog"
	"reflect"
	"testing"

	"github.com/brian-huanggg/esg-award-registration/backend/internal/shared/api"
)

func TestCreateAccount(t *testing.T) {
	valid := &api.CreateAccountJSONRequestBody{EmployeeId: "123456", Name: "Alice", Password: "password1"}
	badID := &api.CreateAccountJSONRequestBody{EmployeeId: "12", Name: "Alice", Password: "password1"}

	tests := []struct {
		name    string
		body    *api.CreateAccountJSONRequestBody
		repoErr error
		want    api.CreateAccountResponseObject // only its type is compared
	}{
		{"created", valid, nil, api.CreateAccount201JSONResponse{}},
		{"invalid input", badID, nil, api.CreateAccount422JSONResponse{}},
		{"duplicate employee id", valid, ErrEmployeeIDTaken, api.CreateAccount409JSONResponse{}},
		{"unexpected error", valid, errors.New("db down"), api.CreateAccountdefaultJSONResponse{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Real service on top of a fake repo: no DB, but the real error mapping.
			h := NewHandler(NewService(&fakeRepo{err: tt.repoErr}), slog.New(slog.DiscardHandler))

			resp, err := h.CreateAccount(context.Background(), api.CreateAccountRequestObject{Body: tt.body})
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if reflect.TypeOf(resp) != reflect.TypeOf(tt.want) {
				t.Fatalf("response = %T, want %T", resp, tt.want)
			}
		})
	}
}

func TestCreateAccount_UnexpectedErrorIs500(t *testing.T) {
	h := NewHandler(NewService(&fakeRepo{err: errors.New("db down")}), slog.New(slog.DiscardHandler))

	resp, _ := h.CreateAccount(context.Background(), api.CreateAccountRequestObject{
		Body: &api.CreateAccountJSONRequestBody{EmployeeId: "123456", Name: "Alice", Password: "password1"},
	})

	got := resp.(api.CreateAccountdefaultJSONResponse)
	if got.StatusCode != 500 {
		t.Errorf("status = %d, want 500", got.StatusCode)
	}
	if got.Body.Message == "db down" {
		t.Error("internal error message must not leak to the client")
	}
}
