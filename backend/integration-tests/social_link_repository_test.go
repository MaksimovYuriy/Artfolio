//go:build integration

package integrationtests

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
	sociallinkrepo "github.com/maksimovyuriy/artfolio/backend/internal/repo/social_link"
)

func TestSocialLinkRepositoryReplace(t *testing.T) {
	database := openTestDatabase(t, "TRUNCATE artist_profiles RESTART IDENTITY CASCADE")
	ctx := context.Background()
	profileID := createArtistProfile(t, ctx, database)
	repository := sociallinkrepo.NewRepo(database)

	first := []entity.SocialLink{
		{Platform: entity.SocialPlatformTelegram, Handle: "anna_art"},
		{Platform: entity.SocialPlatformInstagram, Handle: "anna.art"},
	}
	if err := repository.Replace(ctx, profileID, first); err != nil {
		t.Fatalf("Replace() initial error = %v", err)
	}
	links, err := repository.List(ctx, profileID)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(links) != 2 || links[0].Platform != entity.SocialPlatformTelegram || links[1].Platform != entity.SocialPlatformInstagram {
		t.Fatalf("List() = %#v", links)
	}

	second := []entity.SocialLink{{Platform: entity.SocialPlatformVK, Handle: "anna_art"}}
	if err := repository.Replace(ctx, profileID, second); err != nil {
		t.Fatalf("Replace() second error = %v", err)
	}
	links, err = repository.List(ctx, profileID)
	if err != nil {
		t.Fatalf("List() after replacement error = %v", err)
	}
	if len(links) != 1 || links[0].Platform != entity.SocialPlatformVK {
		t.Fatalf("List() after replacement = %#v", links)
	}

	if err := repository.Replace(ctx, profileID, nil); err != nil {
		t.Fatalf("Replace() empty error = %v", err)
	}
	links, err = repository.List(ctx, profileID)
	if err != nil || len(links) != 0 {
		t.Fatalf("List() after empty replacement = %#v, error = %v", links, err)
	}

	if err := repository.Replace(ctx, 999999, first); !errors.Is(err, repo.ErrNotFound) {
		t.Fatalf("Replace() missing profile error = %v", err)
	}
}

func createArtistProfile(t *testing.T, ctx context.Context, database *sql.DB) int64 {
	t.Helper()
	var id int64
	if err := database.QueryRowContext(ctx, "INSERT INTO artist_profiles (name) VALUES ('Анна') RETURNING id").Scan(&id); err != nil {
		t.Fatalf("create artist profile: %v", err)
	}
	return id
}
