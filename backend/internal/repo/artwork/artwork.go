package artwork

import (
	"context"
	"database/sql"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
)

type Repo struct {
	database *sql.DB
}

func NewRepo(database *sql.DB) *Repo {
	return &Repo{database: database}
}

var _ repo.ArtworkRepository = (*Repo)(nil)

// Delete implements [repo.ArtworkRepository].
func (r *Repo) Delete(ctx context.Context, id int64) error {
	panic("unimplemented")
}

// ListPublished implements [repo.ArtworkRepository].
func (r *Repo) ListPublished(ctx context.Context) ([]entity.Artwork, error) {
	panic("unimplemented")
}

// ListAll implements [repo.ArtworkRepository].
func (r *Repo) ListAll(ctx context.Context) ([]entity.Artwork, error) {
	panic("unimplemented")
}

// Create implements [repo.ArtworkRepository].
func (r *Repo) Create(ctx context.Context, artwork entity.Artwork) (entity.Artwork, error) {
	panic("unimplemented")
}

// Update implements [repo.ArtworkRepository].
func (r *Repo) Update(ctx context.Context, artwork entity.Artwork) error {
	panic("unimplemented")
}
