//go:build integration

package artwork

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

func TestArtworkRepository(t *testing.T) {
	database := openTestDatabase(t)
	repository := NewRepo(database)
	ctx := context.Background()

	draft := createArtwork(t, ctx, repository, entity.Artwork{
		Title: "Черновик", ImageKey: "artworks/draft.jpg", Position: 2,
	})
	first := createArtwork(t, ctx, repository, entity.Artwork{
		Title: "Первая", ImageKey: "artworks/first.jpg", Position: 1, IsPublished: true,
	})
	second := createArtwork(t, ctx, repository, entity.Artwork{
		Title: "Вторая", ImageKey: "artworks/second.jpg", Position: 1, IsPublished: true,
	})

	t.Run("get by id", func(t *testing.T) {
		actual, err := repository.GetByID(ctx, first.ID)
		if err != nil {
			t.Fatalf("GetByID() error = %v", err)
		}
		if actual.Title != first.Title || actual.ImageKey != first.ImageKey {
			t.Fatalf("GetByID() = %#v", actual)
		}
	})

	t.Run("published list filters and orders", func(t *testing.T) {
		artworks, err := repository.ListPublished(ctx)
		if err != nil {
			t.Fatalf("ListPublished() error = %v", err)
		}
		if len(artworks) != 2 || artworks[0].ID != second.ID || artworks[1].ID != first.ID {
			t.Fatalf("ListPublished() = %#v", artworks)
		}
	})

	t.Run("admin list contains drafts", func(t *testing.T) {
		artworks, err := repository.ListAll(ctx)
		if err != nil {
			t.Fatalf("ListAll() error = %v", err)
		}
		if len(artworks) != 3 || artworks[2].ID != draft.ID {
			t.Fatalf("ListAll() = %#v", artworks)
		}
	})

	t.Run("metadata update preserves image", func(t *testing.T) {
		updated := first
		updated.Title = "Обновлённая"
		updated.ImageKey = "artworks/must-not-be-used.jpg"
		previous, err := repository.Update(ctx, updated)
		if err != nil {
			t.Fatalf("Update() error = %v", err)
		}
		actual, err := repository.GetByID(ctx, first.ID)
		if err != nil {
			t.Fatalf("GetByID() after update error = %v", err)
		}
		if previous.ImageKey != first.ImageKey || actual.ImageKey != first.ImageKey {
			t.Fatalf("Update() previous=%#v actual=%#v", previous, actual)
		}
	})

	t.Run("image update returns previous image", func(t *testing.T) {
		updated, err := repository.GetByID(ctx, first.ID)
		if err != nil {
			t.Fatalf("GetByID() before image update error = %v", err)
		}
		updated.ImageKey = "artworks/replacement.jpg"
		previous, err := repository.UpdateWithImage(ctx, updated)
		if err != nil {
			t.Fatalf("UpdateWithImage() error = %v", err)
		}
		actual, err := repository.GetByID(ctx, first.ID)
		if err != nil {
			t.Fatalf("GetByID() after image update error = %v", err)
		}
		if previous.ImageKey != first.ImageKey || actual.ImageKey != updated.ImageKey {
			t.Fatalf("UpdateWithImage() previous=%#v actual=%#v", previous, actual)
		}
	})

	t.Run("delete returns deleted artwork", func(t *testing.T) {
		deleted, err := repository.Delete(ctx, draft.ID)
		if err != nil {
			t.Fatalf("Delete() error = %v", err)
		}
		if deleted.ID != draft.ID || deleted.ImageKey != draft.ImageKey {
			t.Fatalf("Delete() = %#v", deleted)
		}
		_, err = repository.GetByID(ctx, draft.ID)
		if !errors.Is(err, repo.ErrNotFound) {
			t.Fatalf("GetByID() after delete error = %v, want repo.ErrNotFound", err)
		}
	})

	t.Run("missing rows use repository error", func(t *testing.T) {
		_, err := repository.Delete(ctx, 999999)
		if !errors.Is(err, repo.ErrNotFound) {
			t.Fatalf("Delete() error = %v, want repo.ErrNotFound", err)
		}
	})
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
	if _, err := database.ExecContext(ctx, "TRUNCATE artworks RESTART IDENTITY"); err != nil {
		t.Fatalf("truncate artworks: %v", err)
	}
	return database
}

func createArtwork(t *testing.T, ctx context.Context, repository *Repo, artwork entity.Artwork) entity.Artwork {
	t.Helper()
	created, err := repository.Create(ctx, artwork)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	return created
}
