package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-chi/chi"
)

const RequestIDHeader = "X-Request-ID"

type requestStateKey struct{}

type requestState struct {
	id        string
	actorID   int64
	sessionID int64
	err       error
}

type responseRecorder struct {
	http.ResponseWriter
	status int
	bytes  int
}

func (w *responseRecorder) WriteHeader(status int) {
	if w.status != 0 {
		return
	}
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}

func (w *responseRecorder) Write(body []byte) (int, error) {
	if w.status == 0 {
		w.WriteHeader(http.StatusOK)
	}
	n, err := w.ResponseWriter.Write(body)
	w.bytes += n
	return n, err
}

// RecordError attaches an internal error to the current request. It is logged by
// RequestLogger without exposing its details in the HTTP response.
func RecordError(ctx context.Context, err error) {
	if state, ok := ctx.Value(requestStateKey{}).(*requestState); ok {
		state.err = err
	}
}

// RecordAuthentication adds non-sensitive authentication identifiers to the
// request log. Tokens and cookies must never be passed here.
func RecordAuthentication(ctx context.Context, actorID, sessionID int64) {
	if state, ok := ctx.Value(requestStateKey{}).(*requestState); ok {
		state.actorID = actorID
		state.sessionID = sessionID
	}
}

func RequestLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			state := &requestState{id: requestID(r)}
			r = r.WithContext(context.WithValue(r.Context(), requestStateKey{}, state))
			w.Header().Set(RequestIDHeader, state.id)
			recorder := &responseRecorder{ResponseWriter: w}

			defer func() {
				if recovered := recover(); recovered != nil {
					state.err = fmt.Errorf("panic: %v", recovered)
					log.Error("HTTP handler panicked",
						slog.String("request_id", state.id),
						slog.Any("panic", recovered),
						slog.String("stack", string(debug.Stack())),
					)
					if recorder.status == 0 {
						http.Error(recorder, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					}
				}

				status := recorder.status
				if status == 0 {
					status = http.StatusOK
				}
				attrs := []any{
					slog.String("request_id", state.id),
					slog.String("method", r.Method),
					slog.String("path", r.URL.Path),
					slog.Int("status", status),
					slog.Int("response_bytes", recorder.bytes),
					slog.Duration("duration", time.Since(started)),
				}
				if routeContext := chi.RouteContext(r.Context()); routeContext != nil {
					if pattern := routeContext.RoutePattern(); pattern != "" {
						attrs = append(attrs, slog.String("route", pattern))
					}
				}
				if state.err != nil {
					attrs = append(attrs, slog.Any("error", state.err))
				}
				if state.actorID > 0 {
					attrs = append(attrs, slog.Int64("actor_id", state.actorID))
				}
				if state.sessionID > 0 {
					attrs = append(attrs, slog.Int64("session_id", state.sessionID))
				}

				switch {
				case status >= http.StatusInternalServerError:
					log.Error("HTTP request completed", attrs...)
				case status >= http.StatusBadRequest:
					log.Warn("HTTP request completed", attrs...)
				default:
					log.Info("HTTP request completed", attrs...)
				}
			}()

			next.ServeHTTP(recorder, r)
		})
	}
}

func requestID(r *http.Request) string {
	if id := r.Header.Get(RequestIDHeader); validRequestID(id) {
		return id
	}
	var value [16]byte
	if _, err := rand.Read(value[:]); err == nil {
		return hex.EncodeToString(value[:])
	}
	return fmt.Sprintf("fallback-%d", time.Now().UnixNano())
}

func validRequestID(id string) bool {
	if id == "" || len(id) > 128 {
		return false
	}
	return strings.IndexFunc(id, func(char rune) bool {
		return char != '-' && char != '_' &&
			(char < '0' || char > '9') &&
			(char < 'A' || char > 'Z') &&
			(char < 'a' || char > 'z')
	}) == -1
}
