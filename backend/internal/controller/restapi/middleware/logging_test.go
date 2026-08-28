package middleware

import (
	"bytes"
	"errors"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/apierror"
)

func TestRequestLoggerRecordsRequest(t *testing.T) {
	var output bytes.Buffer
	log := slog.New(slog.NewJSONHandler(&output, nil))
	previousLog := slog.Default()
	slog.SetDefault(log)
	t.Cleanup(func() { slog.SetDefault(previousLog) })
	handler := RequestLogger(log)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apierror.Write(w, r, errors.New("database unavailable"))
	}))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/artworks", nil)
	request.Header.Set(RequestIDHeader, "request-123")
	handler.ServeHTTP(response, request)

	if got := response.Header().Get(RequestIDHeader); got != "request-123" {
		t.Fatalf("request ID = %q, want request-123", got)
	}
	for _, expected := range []string{
		`"level":"ERROR"`, `"request_id":"request-123"`,
		`"method":"GET"`, `"path":"/artworks"`, `"status":500`,
		`"error":"database unavailable"`,
	} {
		if !strings.Contains(output.String(), expected) {
			t.Errorf("log does not contain %s: %s", expected, output.String())
		}
	}
}

func TestRequestLoggerRecoversPanic(t *testing.T) {
	var output bytes.Buffer
	log := slog.New(slog.NewJSONHandler(&output, nil))
	handler := RequestLogger(log)(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {
		panic("unexpected")
	}))

	response := httptest.NewRecorder()
	handler.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/panic", nil))

	if response.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusInternalServerError)
	}
	if !strings.Contains(output.String(), "HTTP handler panicked") {
		t.Errorf("panic was not logged: %s", output.String())
	}
}

func TestRequestLoggerRejectsUnsafeRequestID(t *testing.T) {
	log := slog.New(slog.NewTextHandler(&bytes.Buffer{}, nil))
	handler := RequestLogger(log)(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))

	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	request.Header.Set(RequestIDHeader, "unsafe\nvalue")
	handler.ServeHTTP(response, request)

	if got := response.Header().Get(RequestIDHeader); got == "unsafe\nvalue" || got == "" {
		t.Fatalf("generated request ID = %q", got)
	}
}
