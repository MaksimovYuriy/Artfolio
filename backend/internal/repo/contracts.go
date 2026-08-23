package repo

import "context"

type KeyRepository interface {
	Create(ctx context.Context, keyHash []byte) error
}
