package key

import (
	"context"
	"database/sql"
	"errors"
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

func (r *Repo) Find(ctx context.Context, keyHash []byte) (int64, error) {
	const query = `
		SELECT id
		FROM admin_keys
		WHERE key_hash = $1
	`

	var id int64
	if err := r.database.QueryRowContext(ctx, query, keyHash).Scan(&id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, repo.ErrNotFound
		}
		return 0, fmt.Errorf("find admin key: %w", err)
	}

	return id, nil
}
