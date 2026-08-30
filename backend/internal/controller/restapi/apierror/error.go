package apierror

import (
	"database/sql"
	"errors"
	"log/slog"
	"net/http"

	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/jsonutil"
	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/lib/filestorage"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

const requestIDHeader = "X-Request-ID"

type Error struct {
	Status  int
	Code    string
	Message string
	Cause   error
}

func (e *Error) Error() string {
	if e.Cause != nil {
		return e.Cause.Error()
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Cause
}

type response struct {
	Error detail `json:"error"`
}

type detail struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	RequestID string `json:"requestId,omitempty"`
}

func New(status int, code, message string, cause error) error {
	return &Error{Status: status, Code: code, Message: message, Cause: cause}
}

func InvalidRequest(cause error) error {
	return New(http.StatusBadRequest, "invalid_request", "Request is invalid", cause)
}

func UnsupportedMediaType(cause error) error {
	return New(http.StatusUnsupportedMediaType, "unsupported_media_type", "Unsupported media type", cause)
}

func WithStatus(status int, cause error) error {
	switch status {
	case http.StatusBadRequest:
		return InvalidRequest(cause)
	case http.StatusUnauthorized:
		return New(status, "unauthorized", "Authentication required", cause)
	case http.StatusNotFound:
		return New(status, "not_found", "Resource not found", cause)
	case http.StatusConflict:
		return New(status, "conflict", "Request conflicts with current state", cause)
	case http.StatusRequestEntityTooLarge:
		return New(status, "file_too_large", "Uploaded file is too large", cause)
	case http.StatusUnsupportedMediaType:
		return UnsupportedMediaType(cause)
	case http.StatusUnprocessableEntity:
		return New(status, "unprocessable_entity", "Request cannot be processed", cause)
	default:
		return New(http.StatusInternalServerError, "internal_error", "Internal server error", cause)
	}
}

func Write(w http.ResponseWriter, r *http.Request, err error) {
	mapped := mapError(err)
	if mapped.Status >= http.StatusInternalServerError {
		slog.ErrorContext(r.Context(), "HTTP request failed",
			slog.String("request_id", r.Header.Get(requestIDHeader)),
			slog.Any("error", err),
		)
	}

	if encodeErr := jsonutil.Encode(w, mapped.Status, response{Error: detail{
		Code:      mapped.Code,
		Message:   mapped.Message,
		RequestID: r.Header.Get(requestIDHeader),
	}}); encodeErr != nil {
		slog.ErrorContext(r.Context(), "Failed to encode API error",
			slog.String("request_id", r.Header.Get(requestIDHeader)),
			slog.Any("error", encodeErr),
		)
	}
}

func mapError(err error) *Error {
	var apiErr *Error
	if errors.As(err, &apiErr) {
		return apiErr
	}

	switch {
	case errors.Is(err, usecase.ErrInvalidSession):
		return &Error{Status: http.StatusUnauthorized, Code: "unauthorized", Message: "Authentication required"}
	case errors.Is(err, usecase.ErrArtworkNotFound):
		return &Error{Status: http.StatusNotFound, Code: "artwork_not_found", Message: "Artwork not found"}
	case errors.Is(err, sql.ErrNoRows):
		return &Error{Status: http.StatusNotFound, Code: "artist_profile_not_found", Message: "Artist profile not found"}
	case errors.Is(err, entity.ErrValidation), errors.Is(err, usecase.ErrArtworkImageRequired):
		return &Error{Status: http.StatusBadRequest, Code: "validation_failed", Message: "Request validation failed"}
	case errors.Is(err, usecase.ErrArtworkOrderConflict):
		return &Error{Status: http.StatusConflict, Code: "artwork_order_conflict", Message: "Artwork order has changed"}
	case errors.Is(err, usecase.ErrContactRecipientAbsent):
		return &Error{Status: http.StatusServiceUnavailable, Code: "contact_unavailable", Message: "Contact form is unavailable"}
	case errors.Is(err, filestorage.ErrFileTooLarge):
		return &Error{Status: http.StatusRequestEntityTooLarge, Code: "file_too_large", Message: "Uploaded file is too large"}
	case errors.Is(err, filestorage.ErrInvalidImage):
		return &Error{Status: http.StatusUnsupportedMediaType, Code: "invalid_image_type", Message: "Unsupported image type"}
	case errors.Is(err, filestorage.ErrImageTooManyPixels):
		return &Error{Status: http.StatusUnprocessableEntity, Code: "image_too_large", Message: "Image resolution is too large"}
	default:
		return &Error{Status: http.StatusInternalServerError, Code: "internal_error", Message: "Internal server error"}
	}
}
