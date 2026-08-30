package contact

import (
	"context"
	"errors"
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

func TestSendPublishesContactMessage(t *testing.T) {
	recipient := "artist@example.com"
	publisher := &publisherStub{}
	uc := NewUseCase(&profileRepositoryStub{profile: entity.ArtistProfile{Email: &recipient}}, publisher)

	err := uc.Send(context.Background(), entity.ContactMessage{
		SenderEmail: "sender@example.com",
		Message:     "Здравствуйте",
	})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if publisher.event.RecipientEmail != recipient || publisher.event.SenderEmail != "sender@example.com" || publisher.event.Message != "Здравствуйте" {
		t.Fatalf("published event = %#v", publisher.event)
	}
	if publisher.event.EventID == "" {
		t.Fatal("published event has empty event id")
	}
}

func TestSendRequiresConfiguredRecipient(t *testing.T) {
	uc := NewUseCase(&profileRepositoryStub{}, &publisherStub{})
	err := uc.Send(context.Background(), entity.ContactMessage{SenderEmail: "sender@example.com", Message: "Hello"})
	if !errors.Is(err, usecase.ErrContactRecipientAbsent) {
		t.Fatalf("Send() error = %v, want ErrContactRecipientAbsent", err)
	}
}

type profileRepositoryStub struct {
	profile entity.ArtistProfile
}

func (r *profileRepositoryStub) Get(context.Context) (entity.ArtistProfile, error) {
	return r.profile, nil
}

func (r *profileRepositoryStub) Update(context.Context, entity.ArtistProfile) error {
	return nil
}

type publisherStub struct {
	event entity.ContactMessageSubmitted
}

func (p *publisherStub) Publish(_ context.Context, event entity.ContactMessageSubmitted) error {
	p.event = event
	return nil
}
