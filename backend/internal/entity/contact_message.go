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
