package contact

import (
	"context"
	"crypto/rand"
	"encoding/hex"
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
	recipientEmail := ""
	if profile.Email != nil {
		recipientEmail = *profile.Email
	}

	eventID, err := u.newEventID()
	if err != nil {
		return fmt.Errorf("generate contact event id: %w", err)
	}
	event := entity.ContactMessageSubmitted{
		EventID:        eventID,
		RecipientEmail: recipientEmail,
		SenderEmail:    message.SenderEmail,
		Message:        message.Message,
	}
	event, err = event.Validated()
	if err != nil {
		return usecase.ErrContactRecipientAbsent
	}
	if err := u.publisher.Publish(ctx, event); err != nil {
		return fmt.Errorf("publish contact message: %w", err)
	}
	return nil
}

func (u *UseCase) newEventID() (string, error) {
	var value [16]byte
	if _, err := rand.Read(value[:]); err != nil {
		return "", fmt.Errorf("read random bytes: %w", err)
	}

	// UUID version 4 and RFC 4122 variant.
	value[6] = value[6]&0x0f | 0x40
	value[8] = value[8]&0x3f | 0x80

	var encoded [32]byte
	hex.Encode(encoded[:], value[:])
	return string(encoded[0:8]) + "-" +
		string(encoded[8:12]) + "-" +
		string(encoded[12:16]) + "-" +
		string(encoded[16:20]) + "-" +
		string(encoded[20:32]), nil
}
