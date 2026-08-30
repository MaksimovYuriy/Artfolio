package kafka

import (
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/config"
)

func TestNewProducerValidatesConfig(t *testing.T) {
	if _, err := NewProducer(config.KafkaConfig{Topic: "email.send"}); err == nil {
		t.Fatal("NewProducer() accepted empty brokers")
	}
	if _, err := NewProducer(config.KafkaConfig{Brokers: []string{"localhost:9092"}}); err == nil {
		t.Fatal("NewProducer() accepted empty topic")
	}
}
