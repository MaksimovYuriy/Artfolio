package config

import "testing"

func TestLoadKafkaTopics(t *testing.T) {
	t.Setenv("KAFKA_CONTACT_MESSAGE_SUBMITTED_TOPIC", "contact-message.custom.v1")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if cfg.Kafka.Topics.ContactMessageSubmitted != "contact-message.custom.v1" {
		t.Fatalf("contact message topic = %q", cfg.Kafka.Topics.ContactMessageSubmitted)
	}
}
