package kafka

import (
	"encoding/json"
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/config"
	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
)

func TestNewProducerValidatesConfig(t *testing.T) {
	if _, err := NewProducer(config.KafkaConfig{Topic: "contact-message.submitted"}); err == nil {
		t.Fatal("NewProducer() accepted empty brokers")
	}
	if _, err := NewProducer(config.KafkaConfig{Brokers: []string{"localhost:9092"}}); err == nil {
		t.Fatal("NewProducer() accepted empty topic")
	}
}

func TestMessageForUsesEventIDAsKey(t *testing.T) {
	event := entity.ContactMessageSubmitted{
		EventID:        "1c329a48-9402-4a19-86a7-8f67ac06c83b",
		RecipientEmail: "artist@example.com",
		SenderEmail:    "sender@example.com",
		Message:        "Hello",
	}

	message, err := messageFor(event)
	if err != nil {
		t.Fatalf("messageFor() error = %v", err)
	}
	if string(message.Key) != event.EventID {
		t.Fatalf("message key = %q, want %q", message.Key, event.EventID)
	}
	var payload entity.ContactMessageSubmitted
	if err := json.Unmarshal(message.Value, &payload); err != nil {
		t.Fatalf("decode message value: %v", err)
	}
	if payload != event {
		t.Fatalf("message payload = %#v, want %#v", payload, event)
	}
}
