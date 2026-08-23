package key

import "github.com/maksimovyuriy/artfolio/backend/internal/repo/key"

type UseCase struct {
	repo *key.Repo
}

func NewUseCase(repo *key.Repo) *UseCase {
	return &UseCase{repo: repo}
}

// Методы
