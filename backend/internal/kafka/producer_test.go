package kafka

import (
	"context"
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/config"
)

func TestNewProducerValidatesBrokers(t *testing.T) {
	if _, err := NewProducer(config.KafkaConfig{}); err == nil {
		t.Fatal("NewProducer() accepted empty brokers")
	}
	if _, err := NewProducer(config.KafkaConfig{Brokers: []string{""}}); err == nil {
		t.Fatal("NewProducer() accepted empty broker")
	}
}

func TestPublishRequiresTopic(t *testing.T) {
	producer := &Producer{}
	if err := producer.Publish(context.Background(), "", nil, nil); err == nil {
		t.Fatal("Publish() accepted empty topic")
	}
}
