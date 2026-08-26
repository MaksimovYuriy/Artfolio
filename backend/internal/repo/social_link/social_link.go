package sociallink

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

var _ repo.SocialLinkRepository = (*Repo)(nil)

// Create implements [repo.SocialLinkRepository].
func (r *Repo) Create(ctx context.Context, link entity.SocialLink) (entity.SocialLink, error) {
	panic("unimplemented")
}

// Delete implements [repo.SocialLinkRepository].
func (r *Repo) Delete(ctx context.Context, id int64) (entity.SocialLink, error) {
	panic("unimplemented")
}

// List implements [repo.SocialLinkRepository].
func (r *Repo) List(ctx context.Context, artistId int64) ([]entity.SocialLink, error) {
	panic("unimplemented")
}

// Update implements [repo.SocialLinkRepository].
func (r *Repo) Update(ctx context.Context, link entity.SocialLink) (entity.SocialLink, error) {
	panic("unimplemented")
}
