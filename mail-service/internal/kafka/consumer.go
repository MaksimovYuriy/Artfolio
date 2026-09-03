package kafka

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/maksimovyuriy/artfolio/mail-service/internal/config"
	"github.com/maksimovyuriy/artfolio/mail-service/internal/entity"
	kafkago "github.com/segmentio/kafka-go"
)

type Handler func(context.Context, entity.Message) error

type Consumer struct {
	reader *kafkago.Reader
}

func NewConsumer(cfg config.KafkaConfig) (*Consumer, error) {
	if len(cfg.Brokers) == 0 {
		return nil, errors.New("kafka brokers are required")
	}
	for _, broker := range cfg.Brokers {
		if strings.TrimSpace(broker) == "" {
			return nil, errors.New("kafka broker must not be empty")
		}
	}
	if strings.TrimSpace(cfg.Topic) == "" {
		return nil, errors.New("kafka topic is required")
	}
	if strings.TrimSpace(cfg.GroupID) == "" {
		return nil, errors.New("kafka consumer group is required")
	}

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: cfg.Brokers,
		Topic:   cfg.Topic,
		GroupID: cfg.GroupID,
	})
	return &Consumer{reader: reader}, nil
}

func (c *Consumer) Consume(ctx context.Context, handler Handler) error {
	if handler == nil {
		return errors.New("kafka message handler is required")
	}

	for {
		message, err := c.reader.FetchMessage(ctx)
		if err != nil {
			return fmt.Errorf("fetch Kafka message: %w", err)
		}
		if err := handler(ctx, entity.Message{Key: message.Key, Value: message.Value}); err != nil {
			return fmt.Errorf("handle Kafka message: %w", err)
		}
		if err := c.reader.CommitMessages(ctx, message); err != nil {
			return fmt.Errorf("commit Kafka message: %w", err)
		}
	}
}

func (c *Consumer) Close() error {
	if err := c.reader.Close(); err != nil {
		return fmt.Errorf("close Kafka reader: %w", err)
	}
	return nil
}
