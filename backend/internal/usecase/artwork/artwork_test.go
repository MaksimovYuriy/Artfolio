package artwork

import (
	"bytes"
	"context"
	"errors"
	"io"
	"log/slog"
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/repo"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

func TestCreateSavesImageAndCreatesArtwork(t *testing.T) {
	repository := &fakeRepository{}
	files := &fakeStorage{saved: entity.StoredArtworkImage{Key: "artworks/new.jpg"}}
	uc := newTestUseCase(repository, files)

	created, err := uc.Create(context.Background(), entity.Artwork{Title: "  Работа  "}, bytes.NewReader([]byte("image")))
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if repository.created.Title != "Работа" || repository.created.ImageKey != "artworks/new.jpg" {
		t.Fatalf("Create() repository artwork = %#v", repository.created)
	}
	if created.ImageKey != "artworks/new.jpg" {
		t.Fatalf("Create() image key = %q", created.ImageKey)
	}
}

func TestCreateCleansImageWhenRepositoryFails(t *testing.T) {
	repository := &fakeRepository{createErr: errors.New("database unavailable")}
	files := &fakeStorage{saved: entity.StoredArtworkImage{Key: "artworks/new.jpg"}}
	uc := newTestUseCase(repository, files)

	_, err := uc.Create(context.Background(), entity.Artwork{Title: "Работа"}, bytes.NewReader([]byte("image")))
	if err == nil {
		t.Fatal("Create() error = nil")
	}
	if len(files.deleted) != 1 || files.deleted[0] != "artworks/new.jpg" {
		t.Fatalf("Create() deleted keys = %v", files.deleted)
	}
}

func TestCreateValidatesBeforeSavingImage(t *testing.T) {
	files := &fakeStorage{}
	uc := newTestUseCase(&fakeRepository{}, files)

	_, err := uc.Create(context.Background(), entity.Artwork{Title: "   "}, bytes.NewReader([]byte("image")))
	if !errors.Is(err, entity.ErrValidation) {
		t.Fatalf("Create() error = %v, want ErrValidation", err)
	}
	if files.saveCalls != 0 {
		t.Fatalf("Create() storage save calls = %d, want 0", files.saveCalls)
	}
}

func TestUpdateWithoutImageDoesNotReplaceStoredImage(t *testing.T) {
	repository := &fakeRepository{}
	files := &fakeStorage{}
	uc := newTestUseCase(repository, files)

	err := uc.Update(context.Background(), entity.Artwork{ID: 1, Title: "Работа"}, nil)
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if repository.updateWithImage {
		t.Fatal("Update() replaced image")
	}
	if files.saveCalls != 0 || len(files.deleted) != 0 {
		t.Fatalf("Update() storage calls: save=%d delete=%v", files.saveCalls, files.deleted)
	}
}

func TestUpdateWithImageDeletesPreviousImage(t *testing.T) {
	repository := &fakeRepository{previous: entity.Artwork{ImageKey: "artworks/old.jpg"}}
	files := &fakeStorage{saved: entity.StoredArtworkImage{Key: "artworks/new.jpg"}}
	uc := newTestUseCase(repository, files)

	err := uc.Update(context.Background(), entity.Artwork{ID: 1, Title: "Работа"}, bytes.NewReader([]byte("image")))
	if err != nil {
		t.Fatalf("Update() error = %v", err)
	}
	if !repository.updateWithImage || repository.updated.ImageKey != "artworks/new.jpg" {
		t.Fatalf("Update() repository artwork = %#v", repository.updated)
	}
	if len(files.deleted) != 1 || files.deleted[0] != "artworks/old.jpg" {
		t.Fatalf("Update() deleted keys = %v", files.deleted)
	}
}

func TestUpdateCleansNewImageWhenArtworkDoesNotExist(t *testing.T) {
	repository := &fakeRepository{updateErr: repo.ErrNotFound}
	files := &fakeStorage{saved: entity.StoredArtworkImage{Key: "artworks/new.jpg"}}
	uc := newTestUseCase(repository, files)

	err := uc.Update(context.Background(), entity.Artwork{ID: 1, Title: "Работа"}, bytes.NewReader([]byte("image")))
	if !errors.Is(err, usecase.ErrArtworkNotFound) {
		t.Fatalf("Update() error = %v, want ErrArtworkNotFound", err)
	}
	if len(files.deleted) != 1 || files.deleted[0] != "artworks/new.jpg" {
		t.Fatalf("Update() deleted keys = %v", files.deleted)
	}
}

func TestDeleteMapsNotFound(t *testing.T) {
	uc := newTestUseCase(&fakeRepository{deleteErr: repo.ErrNotFound}, &fakeStorage{})

	err := uc.Delete(context.Background(), 1)
	if !errors.Is(err, usecase.ErrArtworkNotFound) {
		t.Fatalf("Delete() error = %v, want ErrArtworkNotFound", err)
	}
}

func newTestUseCase(repository *fakeRepository, files *fakeStorage) *UseCase {
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	return NewUseCase(repository, files, log)
}

type fakeRepository struct {
	created         entity.Artwork
	updated         entity.Artwork
	previous        entity.Artwork
	createErr       error
	updateErr       error
	deleteErr       error
	updateWithImage bool
}

func (r *fakeRepository) ListPublished(context.Context) ([]entity.Artwork, error) {
	return []entity.Artwork{}, nil
}

func (r *fakeRepository) ListAll(context.Context) ([]entity.Artwork, error) {
	return []entity.Artwork{}, nil
}

func (r *fakeRepository) GetByID(context.Context, int64) (entity.Artwork, error) {
	return entity.Artwork{}, nil
}

func (r *fakeRepository) Create(_ context.Context, artwork entity.Artwork) (entity.Artwork, error) {
	r.created = artwork
	return artwork, r.createErr
}

func (r *fakeRepository) Update(_ context.Context, artwork entity.Artwork) (entity.Artwork, error) {
	r.updated = artwork
	return r.previous, r.updateErr
}

func (r *fakeRepository) UpdateWithImage(_ context.Context, artwork entity.Artwork) (entity.Artwork, error) {
	r.updated = artwork
	r.updateWithImage = true
	return r.previous, r.updateErr
}

func (r *fakeRepository) Delete(context.Context, int64) (entity.Artwork, error) {
	return r.previous, r.deleteErr
}

type fakeStorage struct {
	saved     entity.StoredArtworkImage
	saveErr   error
	deleteErr error
	saveCalls int
	deleted   []string
}

func (s *fakeStorage) Save(context.Context, io.Reader) (entity.StoredArtworkImage, error) {
	s.saveCalls++
	return s.saved, s.saveErr
}

func (s *fakeStorage) Delete(_ context.Context, key string) error {
	s.deleted = append(s.deleted, key)
	return s.deleteErr
}
