package kafka

import (
	"context"
	"errors"
	"fmt"

	"github.com/maksimovyuriy/artfolio/backend/internal/config"
	kafkago "github.com/segmentio/kafka-go"
)

type Producer struct {
	writer *kafkago.Writer
}

func NewProducer(cfg config.KafkaConfig) (*Producer, error) {
	if len(cfg.Brokers) == 0 {
		return nil, errors.New("kafka brokers are required")
	}
	for _, broker := range cfg.Brokers {
		if broker == "" {
			return nil, errors.New("kafka broker must not be empty")
		}
	}

	return &Producer{writer: &kafkago.Writer{
		Addr:         kafkago.TCP(cfg.Brokers...),
		Balancer:     &kafkago.Hash{},
		RequiredAcks: kafkago.RequireAll,
		Async:        false,
	}}, nil
}

func (p *Producer) Publish(ctx context.Context, topic string, key, value []byte) error {
	if topic == "" {
		return errors.New("kafka topic is required")
	}

	message := kafkago.Message{Topic: topic, Key: key, Value: value}
	if err := p.writer.WriteMessages(ctx, message); err != nil {
		return fmt.Errorf("write kafka message: %w", err)
	}
	return nil
}

func (p *Producer) Close() error {
	if err := p.writer.Close(); err != nil {
		return fmt.Errorf("close kafka writer: %w", err)
	}
	return nil
}
