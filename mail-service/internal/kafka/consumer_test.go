package kafka

import (
	"testing"

	"github.com/maksimovyuriy/artfolio/mail-service/internal/config"
)

func TestNewConsumerValidatesConfig(t *testing.T) {
	tests := []config.KafkaConfig{
		{},
		{Brokers: []string{""}, Topic: "topic", GroupID: "group"},
		{Brokers: []string{"broker"}, GroupID: "group"},
		{Brokers: []string{"broker"}, Topic: "topic"},
	}

	for _, cfg := range tests {
		if _, err := NewConsumer(cfg); err == nil {
			t.Fatalf("NewConsumer(%#v) accepted invalid config", cfg)
		}
	}
}
