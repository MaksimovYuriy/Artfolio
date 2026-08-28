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
	tx, err := r.database.BeginTx(ctx, nil)
	if err != nil {
		return entity.Artwork{}, fmt.Errorf("begin artwork creation: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "LOCK TABLE artworks IN SHARE ROW EXCLUSIVE MODE"); err != nil {
		return entity.Artwork{}, fmt.Errorf("lock artworks for creation: %w", err)
	}
	if err := tx.QueryRowContext(ctx, "SELECT COALESCE(MAX(position) + 1, 0) FROM artworks").Scan(&artwork.Position); err != nil {
		return entity.Artwork{}, fmt.Errorf("select next artwork position: %w", err)
	}

	const query = `
		INSERT INTO artworks (
			title, description, technique, year, image_key, image_alt, position, is_published
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, title, description, technique, year, image_key, image_alt,
			position, is_published, created_at, updated_at
	`

	created, err := scanArtwork(tx.QueryRowContext(
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
	if err := tx.Commit(); err != nil {
		return entity.Artwork{}, fmt.Errorf("commit artwork creation: %w", err)
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
			is_published = $7,
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
				is_published = $8,
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

func (r *Repo) Reorder(ctx context.Context, artworkIDs []int64) error {
	tx, err := r.database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin artwork reorder: %w", err)
	}
	defer tx.Rollback()

	if _, err := tx.ExecContext(ctx, "LOCK TABLE artworks IN SHARE ROW EXCLUSIVE MODE"); err != nil {
		return fmt.Errorf("lock artworks for reorder: %w", err)
	}

	var count int
	if err := tx.QueryRowContext(ctx, "SELECT COUNT(*) FROM artworks").Scan(&count); err != nil {
		return fmt.Errorf("count artworks for reorder: %w", err)
	}
	if count != len(artworkIDs) {
		return fmt.Errorf("reorder artwork count mismatch: %w", repo.ErrConflict)
	}

	const query = `
		WITH desired_order AS (
			SELECT id, ordinality - 1 AS position
			FROM UNNEST($1::BIGINT[]) WITH ORDINALITY AS item(id, ordinality)
		)
		UPDATE artworks AS artwork
		SET position = desired_order.position,
			updated_at = CURRENT_TIMESTAMP
		FROM desired_order
		WHERE artwork.id = desired_order.id
	`
	result, err := tx.ExecContext(ctx, query, artworkIDs)
	if err != nil {
		return fmt.Errorf("reorder artworks: %w", err)
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("count reordered artworks: %w", err)
	}
	if affected != int64(len(artworkIDs)) {
		return fmt.Errorf("reorder artwork ids mismatch: %w", repo.ErrConflict)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit artwork reorder: %w", err)
	}
	return nil
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
