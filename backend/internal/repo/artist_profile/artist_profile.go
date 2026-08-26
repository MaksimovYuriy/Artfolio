package artistprofile

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

var _ repo.ArtistProfileRepository = (*Repo)(nil)

func (r *Repo) Get(ctx context.Context) (entity.ArtistProfile, error) {
	const query = `
		SELECT id, name, tagline, biography, artist_statement, avatar_url, email
		FROM artist_profiles
		LIMIT 1
	`

	var profile entity.ArtistProfile
	err := r.database.QueryRowContext(ctx, query).Scan(
		&profile.ID,
		&profile.Name,
		&profile.Tagline,
		&profile.Biography,
		&profile.ArtistStatement,
		&profile.AvatarURL,
		&profile.Email,
	)
	if err != nil {
		return entity.ArtistProfile{}, err
	}

	return profile, nil
}

func (r *Repo) Update(ctx context.Context, profile entity.ArtistProfile) error {
	const query = `
		INSERT INTO artist_profiles (
			name,
			tagline,
			biography,
			artist_statement,
			avatar_url,
			email
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT ((true)) DO UPDATE SET
			name = EXCLUDED.name,
			tagline = EXCLUDED.tagline,
			biography = EXCLUDED.biography,
			artist_statement = EXCLUDED.artist_statement,
			avatar_url = EXCLUDED.avatar_url,
			email = EXCLUDED.email,
			updated_at = CURRENT_TIMESTAMP
	`

	_, err := r.database.ExecContext(
		ctx,
		query,
		profile.Name,
		profile.Tagline,
		profile.Biography,
		profile.ArtistStatement,
		profile.AvatarURL,
		profile.Email,
	)
	return err
}
