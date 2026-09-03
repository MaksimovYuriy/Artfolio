package config

import "testing"

func TestLoadKafkaConfig(t *testing.T) {
	t.Setenv("KAFKA_BROKERS", "kafka-1:9092,kafka-2:9092")
	t.Setenv("KAFKA_CONTACT_MESSAGE_SUBMITTED_TOPIC", "contact.custom.v1")
	t.Setenv("KAFKA_CONSUMER_GROUP", "mail-custom")
	t.Setenv("RESEND_API_KEY", "test-api-key")
	t.Setenv("RESEND_FROM", "Artfolio <notifications@example.com>")

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}
	if len(cfg.Kafka.Brokers) != 2 || cfg.Kafka.Brokers[1] != "kafka-2:9092" {
		t.Fatalf("brokers = %#v", cfg.Kafka.Brokers)
	}
	if cfg.Kafka.Topic != "contact.custom.v1" || cfg.Kafka.GroupID != "mail-custom" {
		t.Fatalf("Kafka config = %#v", cfg.Kafka)
	}
	if cfg.Resend.APIKey != "test-api-key" || cfg.Resend.From != "Artfolio <notifications@example.com>" {
		t.Fatalf("Resend config = %#v", cfg.Resend)
	}
}
