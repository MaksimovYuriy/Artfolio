package kafka

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/maksimovyuriy/artfolio/backend/internal/config"
	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
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
	if cfg.Topic == "" {
		return nil, errors.New("kafka topic is required")
	}

	return &Producer{writer: &kafkago.Writer{
		Addr:         kafkago.TCP(cfg.Brokers...),
		Topic:        cfg.Topic,
		Balancer:     &kafkago.Hash{},
		RequiredAcks: kafkago.RequireAll,
		Async:        false,
	}}, nil
}

func (p *Producer) Send(ctx context.Context, email entity.EmailMessage) error {
	payload, err := json.Marshal(email)
	if err != nil {
		return fmt.Errorf("marshal email message: %w", err)
	}

	if err := p.writer.WriteMessages(ctx, kafkago.Message{Key: []byte(email.ReplyTo), Value: payload}); err != nil {
		return fmt.Errorf("write email message: %w", err)
	}
	return nil
}

func (p *Producer) Close() error {
	if err := p.writer.Close(); err != nil {
		return fmt.Errorf("close kafka writer: %w", err)
	}
	return nil
}
