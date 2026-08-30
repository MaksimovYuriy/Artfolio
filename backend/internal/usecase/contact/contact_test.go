package contact

import (
	"context"
	"errors"
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

func TestSendBuildsEmailMessage(t *testing.T) {
	recipient := "artist@example.com"
	producer := &producerStub{}
	uc := NewUseCase(&profileRepositoryStub{profile: entity.ArtistProfile{Email: &recipient}}, producer)

	err := uc.Send(context.Background(), entity.ContactMessage{
		SenderEmail: "sender@example.com",
		Message:     "Здравствуйте",
	})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if producer.email.Recipient != recipient || producer.email.ReplyTo != "sender@example.com" || producer.email.Body != "Здравствуйте" {
		t.Fatalf("produced email = %#v", producer.email)
	}
	if producer.email.Subject == "" {
		t.Fatal("produced email has empty subject")
	}
}

func TestSendRequiresConfiguredRecipient(t *testing.T) {
	uc := NewUseCase(&profileRepositoryStub{}, &producerStub{})
	err := uc.Send(context.Background(), entity.ContactMessage{SenderEmail: "sender@example.com", Message: "Hello"})
	if !errors.Is(err, usecase.ErrEmailRecipientAbsent) {
		t.Fatalf("Send() error = %v, want ErrEmailRecipientAbsent", err)
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

type producerStub struct {
	email entity.EmailMessage
}

func (p *producerStub) Send(_ context.Context, email entity.EmailMessage) error {
	p.email = email
	return nil
}
