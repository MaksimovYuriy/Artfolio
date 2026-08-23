package usecase

import "context"

type KeyCreator interface {
	Create(ctx context.Context) (string, error)
}
