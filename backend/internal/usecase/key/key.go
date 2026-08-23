package key

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"

	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
	"github.com/maksimovyuriy/artfolio/backend/internal/repo/key"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

type UseCase struct {
	repo repo.KeyRepository
}

func NewUseCase(repo *key.Repo) *UseCase {
	return &UseCase{repo: repo}
}

var _ usecase.KeyCreator = (*UseCase)(nil)

func (uc UseCase) Create(ctx context.Context) (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	accessKey := "artfolio_" + base64.RawURLEncoding.EncodeToString(bytes)
	keyHash := sha256.Sum256([]byte(accessKey))

	if err := uc.repo.Create(ctx, keyHash[:]); err != nil {
		return "", err
	}

	return accessKey, nil
}
