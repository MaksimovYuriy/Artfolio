package usecase

import (
	"context"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
)

type KeyCreator interface {
	Create(ctx context.Context) (string, error)
}

type SessionUseCase interface {
	Create(ctx context.Context, accessKey string) (entity.Session, error)
	Verify(ctx context.Context, token string) (bool, error)
	Revoke(ctx context.Context, token string) error
}
