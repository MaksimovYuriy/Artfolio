//go:build integration

package integrationtests

import (
	"context"
	"errors"
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
	artworkrepo "github.com/maksimovyuriy/artfolio/backend/internal/repo/artwork"
)

func TestArtworkRepository(t *testing.T) {
	database := openTestDatabase(t, "TRUNCATE artworks RESTART IDENTITY")
	repository := artworkrepo.NewRepo(database)
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
		if len(artworks) != 2 || artworks[0].ID != first.ID || artworks[1].ID != second.ID {
			t.Fatalf("ListPublished() = %#v", artworks)
		}
	})

	t.Run("admin list contains drafts", func(t *testing.T) {
		artworks, err := repository.ListAll(ctx)
		if err != nil {
			t.Fatalf("ListAll() error = %v", err)
		}
		if len(artworks) != 3 || artworks[0].ID != draft.ID {
			t.Fatalf("ListAll() = %#v", artworks)
		}
	})

	t.Run("reorder all artworks", func(t *testing.T) {
		if err := repository.Reorder(ctx, []int64{second.ID, draft.ID, first.ID}); err != nil {
			t.Fatalf("Reorder() error = %v", err)
		}
		artworks, err := repository.ListAll(ctx)
		if err != nil {
			t.Fatalf("ListAll() after reorder error = %v", err)
		}
		if len(artworks) != 3 || artworks[0].ID != second.ID || artworks[1].ID != draft.ID || artworks[2].ID != first.ID {
			t.Fatalf("ListAll() after reorder = %#v", artworks)
		}
		if err := repository.Reorder(ctx, []int64{second.ID, draft.ID, 999999}); !errors.Is(err, repo.ErrConflict) {
			t.Fatalf("Reorder() missing id error = %v, want repo.ErrConflict", err)
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

func createArtwork(t *testing.T, ctx context.Context, repository *artworkrepo.Repo, artwork entity.Artwork) entity.Artwork {
	t.Helper()
	created, err := repository.Create(ctx, artwork)
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	return created
}
