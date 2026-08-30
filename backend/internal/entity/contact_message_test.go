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
