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

	sessionID, err := uc.sessionRepo.Create(ctx, adminKeyID, tokenHash[:], expiresAt)
	if err != nil {
		return entity.Session{}, err
	}

	return entity.Session{
		ID:        sessionID,
		ActorID:   adminKeyID,
		Token:     token,
		ExpiresAt: expiresAt,
	}, nil
}

func (uc *UseCase) Authenticate(ctx context.Context, token string) (entity.AuthenticatedSession, error) {
	tokenHash := sha256.Sum256([]byte(token))

	session, err := uc.sessionRepo.FindActive(ctx, tokenHash[:], time.Now())
	if errors.Is(err, repo.ErrNotFound) {
		return entity.AuthenticatedSession{}, usecase.ErrInvalidSession
	}
	if err != nil {
		return entity.AuthenticatedSession{}, err
	}

	return session, nil
}

func (uc *UseCase) Revoke(ctx context.Context, token string) (entity.AuthenticatedSession, error) {
	tokenHash := sha256.Sum256([]byte(token))

	session, err := uc.sessionRepo.Revoke(ctx, tokenHash[:], time.Now())
	if errors.Is(err, repo.ErrNotFound) {
		return entity.AuthenticatedSession{}, usecase.ErrInvalidSession
	}
	if err != nil {
		return entity.AuthenticatedSession{}, err
	}

	return session, nil
}
