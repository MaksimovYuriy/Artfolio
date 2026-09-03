package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/maksimovyuriy/artfolio/mail-service/internal/entity"
)

type Email struct {
	log *slog.Logger
}

func NewEmail(log *slog.Logger) *Email {
	return &Email{log: log}
}

func (h *Email) Handle(_ context.Context, message entity.Message) error {
	var event entity.ContactMessageSubmitted
	if err := json.Unmarshal(message.Value, &event); err != nil {
		return fmt.Errorf("decode contact message submitted event: %w", err)
	}
	if err := event.Validate(); err != nil {
		return fmt.Errorf("validate contact message submitted event: %w", err)
	}

	h.log.Info("Contact message received", slog.String("event_id", event.EventID))
	return nil
}
