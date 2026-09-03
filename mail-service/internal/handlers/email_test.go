package handlers

import (
	"context"
	"io"
	"log/slog"
	"testing"

	"github.com/maksimovyuriy/artfolio/mail-service/internal/entity"
)

func TestHandleContactMessageSubmitted(t *testing.T) {
	handler := NewEmail(slog.New(slog.NewTextHandler(io.Discard, nil)))
	message := entity.Message{Value: []byte(`{
		"eventId":"d193ddb5-b5d1-425c-b516-30687cb09980",
		"recipientEmail":"artist@example.com",
		"senderEmail":"sender@example.com",
		"message":"Hello"
	}`)}

	if err := handler.Handle(context.Background(), message); err != nil {
		t.Fatalf("Handle() error = %v", err)
	}
}

func TestHandleRejectsInvalidMessage(t *testing.T) {
	handler := NewEmail(slog.New(slog.NewTextHandler(io.Discard, nil)))
	for _, value := range [][]byte{[]byte(`{`), []byte(`{"eventId":"event"}`)} {
		if err := handler.Handle(context.Background(), entity.Message{Value: value}); err == nil {
			t.Fatalf("Handle(%q) accepted invalid message", value)
		}
	}
}
