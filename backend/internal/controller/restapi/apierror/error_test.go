package apierror

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

func TestWriteMapsErrors(t *testing.T) {
	tests := []struct {
		name       string
		err        error
		wantStatus int
		wantCode   string
	}{
		{name: "invalid request", err: InvalidRequest(errors.New("bad json")), wantStatus: 400, wantCode: "invalid_request"},
		{name: "validation", err: entity.ErrValidation, wantStatus: 400, wantCode: "validation_failed"},
		{name: "unauthorized", err: usecase.ErrInvalidSession, wantStatus: 401, wantCode: "unauthorized"},
		{name: "not found", err: usecase.ErrArtworkNotFound, wantStatus: 404, wantCode: "artwork_not_found"},
		{name: "conflict", err: usecase.ErrArtworkOrderConflict, wantStatus: 409, wantCode: "artwork_order_conflict"},
		{name: "internal", err: errors.New("database unavailable"), wantStatus: 500, wantCode: "internal_error"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			request := httptest.NewRequest(http.MethodGet, "/", nil)
			request.Header.Set(requestIDHeader, "request-123")
			result := httptest.NewRecorder()

			Write(result, request, test.err)

			if result.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", result.Code, test.wantStatus)
			}
			var payload response
			if err := json.Unmarshal(result.Body.Bytes(), &payload); err != nil {
				t.Fatalf("decode response: %v", err)
			}
			if payload.Error.Code != test.wantCode || payload.Error.RequestID != "request-123" {
				t.Fatalf("error = %+v, want code %q and request ID", payload.Error, test.wantCode)
			}
		})
	}
}
