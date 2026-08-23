package key

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
)

type Repo struct {
	database *sql.DB
}

func NewRepo(database *sql.DB) *Repo {
	return &Repo{database: database}
}

var _ repo.KeyRepository = (*Repo)(nil)

func (r *Repo) Create(ctx context.Context, keyHash []byte) error {
	query := `
		INSERT INTO admin_keys (key_hash)
		VALUES ($1);
	`

	_, err := r.database.ExecContext(ctx, query, keyHash)
	if err != nil {
		return fmt.Errorf("Create admin key: %w", err)
	}

	return nil
}
