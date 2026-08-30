package contact

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
)

type producer interface {
	Publish(ctx context.Context, topic string, key, value []byte) error
}

type Publisher struct {
	producer producer
	topic    string
}

func NewPublisher(producer producer, topic string) (*Publisher, error) {
	if topic == "" {
		return nil, errors.New("contact message topic is required")
	}
	return &Publisher{producer: producer, topic: topic}, nil
}

func (p *Publisher) Publish(ctx context.Context, event entity.ContactMessageSubmitted) error {
	payload, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal contact message submitted event: %w", err)
	}
	if err := p.producer.Publish(ctx, p.topic, []byte(event.EventID), payload); err != nil {
		return fmt.Errorf("publish contact message submitted event: %w", err)
	}
	return nil
}
