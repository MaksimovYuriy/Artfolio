package entity

import (
	"net/mail"
	"strings"
)

const maxContactMessageLength = 10_000

type ContactMessage struct {
	SenderEmail string
	Message     string
}

// ContactMessageSubmitted is the integration event published after the
// backend has accepted a contact form submission and resolved its recipient.
type ContactMessageSubmitted struct {
	EventID        string `json:"eventId"`
	RecipientEmail string `json:"recipientEmail"`
	SenderEmail    string `json:"senderEmail"`
	Message        string `json:"message"`
}

func (event ContactMessageSubmitted) Validated() (ContactMessageSubmitted, error) {
	event.RecipientEmail = strings.TrimSpace(event.RecipientEmail)
	if event.RecipientEmail == "" {
		return ContactMessageSubmitted{}, NewValidationError("recipientEmail", "is required")
	}

	address, err := mail.ParseAddress(event.RecipientEmail)
	if err != nil || address.Address != event.RecipientEmail || !isASCII(event.RecipientEmail) {
		return ContactMessageSubmitted{}, NewValidationError("recipientEmail", "is invalid")
	}

	return event, nil
}

func (message ContactMessage) Validated() (ContactMessage, error) {
	message.SenderEmail = strings.TrimSpace(message.SenderEmail)
	message.Message = strings.TrimSpace(message.Message)

	address, err := mail.ParseAddress(message.SenderEmail)
	if err != nil || address.Address != message.SenderEmail || !isASCII(message.SenderEmail) {
		return ContactMessage{}, NewValidationError("senderEmail", "is invalid")
	}
	if message.Message == "" {
		return ContactMessage{}, NewValidationError("message", "is required")
	}
	if len([]rune(message.Message)) > maxContactMessageLength {
		return ContactMessage{}, NewValidationError("message", "is too long")
	}

	return message, nil
}
