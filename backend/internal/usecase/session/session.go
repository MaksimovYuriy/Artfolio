package session

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"time"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

const sessionTTL = 24 * time.Hour

type UseCase struct {
	keyRepo     repo.KeyRepository
	sessionRepo repo.SessionRepository
}

func NewUseCase(keyRepo repo.KeyRepository, sessionRepo repo.SessionRepository) *UseCase {
	return &UseCase{
		keyRepo:     keyRepo,
		sessionRepo: sessionRepo,
	}
}

var _ usecase.SessionUseCase = (*UseCase)(nil)

func (uc *UseCase) Create(ctx context.Context, accessKey string) (entity.Session, error) {
	accessKeyHash := sha256.Sum256([]byte(accessKey))

	adminKeyID, err := uc.keyRepo.Find(ctx, accessKeyHash[:])
	if errors.Is(err, repo.ErrNotFound) {
		return entity.Session{}, usecase.ErrInvalidSession
	}
	if err != nil {
		return entity.Session{}, err
	}

	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		return entity.Session{}, err
	}

	token := base64.RawURLEncoding.EncodeToString(randomBytes)
	tokenHash := sha256.Sum256([]byte(token))
	expiresAt := time.Now().Add(sessionTTL)

	if err := uc.sessionRepo.Create(ctx, adminKeyID, tokenHash[:], expiresAt); err != nil {
		return entity.Session{}, err
	}

	return entity.Session{
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (uc *UseCase) Verify(ctx context.Context, token string) (bool, error) {
	tokenHash := sha256.Sum256([]byte(token))

	exists, err := uc.sessionRepo.ExistsActive(ctx, tokenHash[:], time.Now())
	if err != nil {
		return false, err
	}

	return exists, nil
}

func (uc *UseCase) Revoke(ctx context.Context, token string) error {
	tokenHash := sha256.Sum256([]byte(token))

	if err := uc.sessionRepo.Revoke(ctx, tokenHash[:], time.Now()); err != nil {
		return err
	}

	return nil
}
