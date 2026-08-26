package sociallink

import (
	"context"
	"database/sql"
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

var _ repo.SocialLinkRepository = (*Repo)(nil)

func (r *Repo) List(ctx context.Context, artistProfileID int64) ([]entity.SocialLink, error) {
	const query = `
		SELECT artist_profile_id, platform, handle
		FROM artist_social_links
		WHERE artist_profile_id = $1
		ORDER BY CASE platform
			WHEN 'telegram' THEN 1
			WHEN 'instagram' THEN 2
			WHEN 'vk' THEN 3
			WHEN 'behance' THEN 4
			ELSE 5
		END
	`

	rows, err := r.database.QueryContext(ctx, query, artistProfileID)
	if err != nil {
		return nil, fmt.Errorf("list social links: %w", err)
	}
	defer rows.Close()

	links := make([]entity.SocialLink, 0)
	for rows.Next() {
		var link entity.SocialLink
		if err := rows.Scan(&link.ArtistProfileID, &link.Platform, &link.Handle); err != nil {
			return nil, fmt.Errorf("scan social link: %w", err)
		}
		links = append(links, link)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("iterate social links: %w", err)
	}
	return links, nil
}

func (r *Repo) Replace(ctx context.Context, artistProfileID int64, links []entity.SocialLink) error {
	tx, err := r.database.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin social links replacement: %w", err)
	}
	defer tx.Rollback()

	var exists bool
	if err := tx.QueryRowContext(
		ctx,
		"SELECT EXISTS (SELECT 1 FROM artist_profiles WHERE id = $1)",
		artistProfileID,
	).Scan(&exists); err != nil {
		return fmt.Errorf("check artist profile: %w", err)
	}
	if !exists {
		return repo.ErrNotFound
	}

	if _, err := tx.ExecContext(ctx, "DELETE FROM artist_social_links WHERE artist_profile_id = $1", artistProfileID); err != nil {
		return fmt.Errorf("delete previous social links: %w", err)
	}

	const insertQuery = `
		INSERT INTO artist_social_links (artist_profile_id, platform, handle)
		VALUES ($1, $2, $3)
	`
	for _, link := range links {
		if _, err := tx.ExecContext(ctx, insertQuery, artistProfileID, link.Platform, link.Handle); err != nil {
			return fmt.Errorf("insert social link: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit social links replacement: %w", err)
	}
	return nil
}
