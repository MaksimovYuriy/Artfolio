package contact

import (
	"context"
	"fmt"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

type eventPublisher interface {
	Publish(ctx context.Context, event entity.ContactMessageSubmitted) error
}

type UseCase struct {
	profileRepo repo.ArtistProfileRepository
	publisher   eventPublisher
}

func NewUseCase(profileRepo repo.ArtistProfileRepository, publisher eventPublisher) *UseCase {
	return &UseCase{profileRepo: profileRepo, publisher: publisher}
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
		return usecase.ErrContactRecipientAbsent
	}

	eventID, err := newEventID()
	if err != nil {
		return fmt.Errorf("generate contact event id: %w", err)
	}
	event := entity.ContactMessageSubmitted{
		EventID:        eventID,
		RecipientEmail: *profile.Email,
		SenderEmail:    message.SenderEmail,
		Message:        message.Message,
	}
	if err := u.publisher.Publish(ctx, event); err != nil {
		return fmt.Errorf("publish contact message: %w", err)
	}
	return nil
}
