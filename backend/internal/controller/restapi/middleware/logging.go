package middleware

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/apierror"
)

const RequestIDHeader = "X-Request-ID"

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

func RequestLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			started := time.Now()
			requestID := requestID(r)
			r.Header.Set(RequestIDHeader, requestID)
			w.Header().Set(RequestIDHeader, requestID)
			recorder := &responseRecorder{ResponseWriter: w}

			defer func() {
				if recovered := recover(); recovered != nil {
					panicErr := fmt.Errorf("panic: %v", recovered)
					log.Error("HTTP handler panicked",
						slog.String("request_id", requestID),
						slog.Any("panic", recovered),
						slog.String("stack", string(debug.Stack())),
					)
					if recorder.status == 0 {
						apierror.Write(recorder, r, panicErr)
					}
				}

				status := recorder.status
				if status == 0 {
					status = http.StatusOK
				}
				attrs := []any{
					slog.String("request_id", requestID),
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
