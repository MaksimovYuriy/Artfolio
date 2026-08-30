package contact

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
)

func TestNewPublisherRequiresTopic(t *testing.T) {
	if _, err := NewPublisher(&producerStub{}, ""); err == nil {
		t.Fatal("NewPublisher() accepted empty topic")
	}
}

func TestPublishMapsContactEventToKafkaMessage(t *testing.T) {
	transport := &producerStub{}
	publisher, err := NewPublisher(transport, "contact-message.submitted.v1")
	if err != nil {
		t.Fatalf("NewPublisher() error = %v", err)
	}
	event := entity.ContactMessageSubmitted{
		EventID:        "1c329a48-9402-4a19-86a7-8f67ac06c83b",
		RecipientEmail: "artist@example.com",
		SenderEmail:    "sender@example.com",
		Message:        "Hello",
	}

	if err := publisher.Publish(context.Background(), event); err != nil {
		t.Fatalf("Publish() error = %v", err)
	}
	if transport.topic != "contact-message.submitted.v1" {
		t.Fatalf("topic = %q", transport.topic)
	}
	if string(transport.key) != event.EventID {
		t.Fatalf("key = %q, want %q", transport.key, event.EventID)
	}
	var payload entity.ContactMessageSubmitted
	if err := json.Unmarshal(transport.value, &payload); err != nil {
		t.Fatalf("decode message value: %v", err)
	}
	if payload != event {
		t.Fatalf("payload = %#v, want %#v", payload, event)
	}
}

type producerStub struct {
	topic string
	key   []byte
	value []byte
}

func (p *producerStub) Publish(_ context.Context, topic string, key, value []byte) error {
	p.topic = topic
	p.key = key
	p.value = value
	return nil
}
