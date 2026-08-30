package entity

import (
	"errors"
	"testing"
)

func TestContactMessageValidated(t *testing.T) {
	message, err := (ContactMessage{SenderEmail: " sender@example.com ", Message: " Hello "}).Validated()
	if err != nil {
		t.Fatalf("Validated() error = %v", err)
	}
	if message.SenderEmail != "sender@example.com" || message.Message != "Hello" {
		t.Fatalf("Validated() = %#v", message)
	}
}

func TestContactMessageRejectsInvalidInput(t *testing.T) {
	tests := []ContactMessage{
		{SenderEmail: "invalid", Message: "Hello"},
		{SenderEmail: "sender@example.com", Message: "   "},
	}
	for _, message := range tests {
		if _, err := message.Validated(); !errors.Is(err, ErrValidation) {
			t.Fatalf("Validated(%#v) error = %v, want ErrValidation", message, err)
		}
	}
}

func TestContactMessageSubmittedValidatesRecipientEmail(t *testing.T) {
	event, err := (ContactMessageSubmitted{RecipientEmail: " artist@example.com "}).Validated()
	if err != nil {
		t.Fatalf("Validated() error = %v", err)
	}
	if event.RecipientEmail != "artist@example.com" {
		t.Fatalf("RecipientEmail = %q", event.RecipientEmail)
	}
}

func TestContactMessageSubmittedRequiresValidRecipientEmail(t *testing.T) {
	for _, recipientEmail := range []string{"", "   ", "invalid"} {
		event := ContactMessageSubmitted{RecipientEmail: recipientEmail}
		if _, err := event.Validated(); !errors.Is(err, ErrValidation) {
			t.Fatalf("Validated(%q) error = %v, want ErrValidation", recipientEmail, err)
		}
	}
}
