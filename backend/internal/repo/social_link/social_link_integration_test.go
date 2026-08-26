//go:build integration

package sociallink

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
	"github.com/pressly/goose"
)

const defaultTestDatabaseURL = "postgres://artfolio_test:artfolio_test@127.0.0.1:55432/artfolio_test?sslmode=disable"

func TestSocialLinkRepositoryReplace(t *testing.T) {
	database := openTestDatabase(t)
	ctx := context.Background()
	profileID := createArtistProfile(t, ctx, database)
	repository := NewRepo(database)

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

func openTestDatabase(t *testing.T) *sql.DB {
	t.Helper()
	databaseURL := os.Getenv("ARTFOLIO_TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = defaultTestDatabaseURL
	}
	database, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatalf("open test database: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := database.PingContext(ctx); err != nil {
		t.Fatalf("connect to test database: %v", err)
	}
	lockConnection, err := database.Conn(ctx)
	if err != nil {
		t.Fatalf("reserve integration test connection: %v", err)
	}
	if _, err := lockConnection.ExecContext(ctx, "SELECT pg_advisory_lock(91724001)"); err != nil {
		_ = lockConnection.Close()
		t.Fatalf("lock integration database: %v", err)
	}
	t.Cleanup(func() {
		_, _ = lockConnection.ExecContext(context.Background(), "SELECT pg_advisory_unlock(91724001)")
		_ = lockConnection.Close()
	})

	if err := goose.SetDialect("postgres"); err != nil {
		t.Fatalf("set migration dialect: %v", err)
	}
	migrations := filepath.Join("..", "..", "..", "migrations")
	if err := goose.Up(database, migrations); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	if _, err := database.ExecContext(ctx, "TRUNCATE artist_profiles RESTART IDENTITY CASCADE"); err != nil {
		t.Fatalf("truncate artist profiles: %v", err)
	}
	return database
}

func createArtistProfile(t *testing.T, ctx context.Context, database *sql.DB) int64 {
	t.Helper()
	var id int64
	if err := database.QueryRowContext(ctx, "INSERT INTO artist_profiles (name) VALUES ('Анна') RETURNING id").Scan(&id); err != nil {
		t.Fatalf("create artist profile: %v", err)
	}
	return id
}
