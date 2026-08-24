package artwork

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

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

func (r *Repo) ListPublished(ctx context.Context) ([]entity.Artwork, error) {
	const query = `
		SELECT id, title, description, technique, year, image_key, image_alt,
			position, is_published, created_at, updated_at
		FROM artworks
		WHERE is_published = TRUE
		ORDER BY position ASC, created_at DESC, id DESC
	`

	return r.list(ctx, query)
}

func (r *Repo) ListAll(ctx context.Context) ([]entity.Artwork, error) {
	const query = `
		SELECT id, title, description, technique, year, image_key, image_alt,
			position, is_published, created_at, updated_at
		FROM artworks
		ORDER BY position ASC, created_at DESC, id DESC
	`

	return r.list(ctx, query)
}

func (r *Repo) GetByID(ctx context.Context, id int64) (entity.Artwork, error) {
	const query = `
		SELECT id, title, description, technique, year, image_key, image_alt,
			position, is_published, created_at, updated_at
		FROM artworks
		WHERE id = $1
	`

	artwork, err := scanArtwork(r.database.QueryRowContext(ctx, query, id))
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Artwork{}, fmt.Errorf("get artwork by id: %w", repo.ErrNotFound)
	}
	if err != nil {
		return entity.Artwork{}, fmt.Errorf("get artwork by id: %w", err)
	}
	return artwork, nil
}

func (r *Repo) Create(ctx context.Context, artwork entity.Artwork) (entity.Artwork, error) {
	const query = `
		INSERT INTO artworks (
			title, description, technique, year, image_key, image_alt, position, is_published
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, title, description, technique, year, image_key, image_alt,
			position, is_published, created_at, updated_at
	`

	created, err := scanArtwork(r.database.QueryRowContext(
		ctx,
		query,
		artwork.Title,
		artwork.Description,
		artwork.Technique,
		artwork.Year,
		artwork.ImageKey,
		artwork.ImageAlt,
		artwork.Position,
		artwork.IsPublished,
	))
	if err != nil {
		return entity.Artwork{}, fmt.Errorf("create artwork: %w", err)
	}
	return created, nil
}

func (r *Repo) Update(ctx context.Context, artwork entity.Artwork) (entity.Artwork, error) {
	return r.update(ctx, artwork, false)
}

func (r *Repo) UpdateWithImage(ctx context.Context, artwork entity.Artwork) (entity.Artwork, error) {
	return r.update(ctx, artwork, true)
}

func (r *Repo) update(ctx context.Context, artwork entity.Artwork, updateImage bool) (entity.Artwork, error) {
	tx, err := r.database.BeginTx(ctx, nil)
	if err != nil {
		return entity.Artwork{}, fmt.Errorf("begin artwork update: %w", err)
	}
	defer tx.Rollback()

	const selectQuery = `
		SELECT id, title, description, technique, year, image_key, image_alt,
			position, is_published, created_at, updated_at
		FROM artworks
		WHERE id = $1
		FOR UPDATE
	`
	previous, err := scanArtwork(tx.QueryRowContext(ctx, selectQuery, artwork.ID))
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Artwork{}, fmt.Errorf("get artwork for update: %w", repo.ErrNotFound)
	}
	if err != nil {
		return entity.Artwork{}, fmt.Errorf("get artwork for update: %w", err)
	}

	query := `
		UPDATE artworks
		SET title = $2,
			description = $3,
			technique = $4,
			year = $5,
			image_alt = $6,
			position = $7,
			is_published = $8,
			updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`
	arguments := []any{
		artwork.ID,
		artwork.Title,
		artwork.Description,
		artwork.Technique,
		artwork.Year,
		artwork.ImageAlt,
		artwork.Position,
		artwork.IsPublished,
	}
	if updateImage {
		query = `
			UPDATE artworks
			SET title = $2,
				description = $3,
				technique = $4,
				year = $5,
				image_key = $6,
				image_alt = $7,
				position = $8,
				is_published = $9,
				updated_at = CURRENT_TIMESTAMP
			WHERE id = $1
		`
		arguments = []any{
			artwork.ID,
			artwork.Title,
			artwork.Description,
			artwork.Technique,
			artwork.Year,
			artwork.ImageKey,
			artwork.ImageAlt,
			artwork.Position,
			artwork.IsPublished,
		}
	}

	_, err = tx.ExecContext(ctx, query, arguments...)
	if err != nil {
		return entity.Artwork{}, fmt.Errorf("update artwork: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return entity.Artwork{}, fmt.Errorf("commit artwork update: %w", err)
	}
	return previous, nil
}

func (r *Repo) Delete(ctx context.Context, id int64) (entity.Artwork, error) {
	const query = `
		DELETE FROM artworks
		WHERE id = $1
		RETURNING id, title, description, technique, year, image_key, image_alt,
			position, is_published, created_at, updated_at
	`

	deleted, err := scanArtwork(r.database.QueryRowContext(ctx, query, id))
	if errors.Is(err, sql.ErrNoRows) {
		return entity.Artwork{}, fmt.Errorf("delete artwork: %w", repo.ErrNotFound)
	}
	if err != nil {
		return entity.Artwork{}, fmt.Errorf("delete artwork: %w", err)
	}
	return deleted, nil
}

func (r *Repo) list(ctx context.Context, query string) ([]entity.Artwork, error) {
	rows, err := r.database.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("list artworks: %w", err)
	}
	defer rows.Close()

	artworks := make([]entity.Artwork, 0)
	for rows.Next() {
		artwork, err := scanArtwork(rows)
		if err != nil {
			return nil, fmt.Errorf("scan artwork: %w", err)
		}
		artworks = append(artworks, artwork)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate artworks: %w", err)
	}
	return artworks, nil
}

type scanner interface {
	Scan(dest ...any) error
}

func scanArtwork(row scanner) (entity.Artwork, error) {
	var artwork entity.Artwork
	err := row.Scan(
		&artwork.ID,
		&artwork.Title,
		&artwork.Description,
		&artwork.Technique,
		&artwork.Year,
		&artwork.ImageKey,
		&artwork.ImageAlt,
		&artwork.Position,
		&artwork.IsPublished,
		&artwork.CreatedAt,
		&artwork.UpdatedAt,
	)
	return artwork, err
}
