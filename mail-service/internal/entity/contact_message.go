package entity

import (
	"errors"
	"net/mail"
	"strings"
)

type ContactMessageSubmitted struct {
	EventID        string `json:"eventId"`
	RecipientEmail string `json:"recipientEmail"`
	SenderEmail    string `json:"senderEmail"`
	Message        string `json:"message"`
}

func (event ContactMessageSubmitted) Validate() error {
	if strings.TrimSpace(event.EventID) == "" {
		return errors.New("eventId is required")
	}
	if !validEmail(event.RecipientEmail) {
		return errors.New("recipientEmail is invalid")
	}
	if !validEmail(event.SenderEmail) {
		return errors.New("senderEmail is invalid")
	}
	if strings.TrimSpace(event.Message) == "" {
		return errors.New("message is required")
	}
	return nil
}

func validEmail(value string) bool {
	value = strings.TrimSpace(value)
	address, err := mail.ParseAddress(value)
	return err == nil && address.Address == value
}
