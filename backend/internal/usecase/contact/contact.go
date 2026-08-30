package contact

import (
	"context"
	"fmt"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

const emailSubject = "Новое сообщение из портфолио"

type emailProducer interface {
	Send(ctx context.Context, email entity.EmailMessage) error
}

type UseCase struct {
	profileRepo repo.ArtistProfileRepository
	producer    emailProducer
}

func NewUseCase(profileRepo repo.ArtistProfileRepository, producer emailProducer) *UseCase {
	return &UseCase{profileRepo: profileRepo, producer: producer}
}

var _ usecase.ContactUseCase = (*UseCase)(nil)

func (u *UseCase) Send(ctx context.Context, message entity.ContactMessage) error {
	message, err := message.Validated()
	if err != nil {
		return err
	}

	profile, err := u.profileRepo.Get(ctx)
	if err != nil {
		return fmt.Errorf("get contact email recipient: %w", err)
	}
	if profile.Email == nil || *profile.Email == "" {
		return usecase.ErrEmailRecipientAbsent
	}

	email := entity.EmailMessage{Recipient: *profile.Email, ReplyTo: message.SenderEmail, Subject: emailSubject, Body: message.Message}
	if err := u.producer.Send(ctx, email); err != nil {
		return fmt.Errorf("enqueue contact email: %w", err)
	}
	return nil
}
