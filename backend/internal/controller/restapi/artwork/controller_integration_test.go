//go:build integration

package artwork

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/maksimovyuriy/artfolio/backend/internal/config"
	artworkdto "github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/dto/artwork"
	artworkstorage "github.com/maksimovyuriy/artfolio/backend/internal/lib/storage/artwork"
	artworkrepo "github.com/maksimovyuriy/artfolio/backend/internal/repo/artwork"
	artworkusecase "github.com/maksimovyuriy/artfolio/backend/internal/usecase/artwork"
	"github.com/pressly/goose"
)

const integrationDatabaseURL = "postgres://artfolio_test:artfolio_test@127.0.0.1:55432/artfolio_test?sslmode=disable"

func TestArtworkHTTPFlow(t *testing.T) {
	database := openIntegrationDatabase(t)
	mediaRoot := t.TempDir()
	storageConfig := config.StorageConfig{
		Path:        mediaRoot,
		PublicURL:   "/media",
		MaxFileSize: 1 << 20,
		MaxPixels:   1_000_000,
	}
	files, err := artworkstorage.New(storageConfig)
	if err != nil {
		t.Fatalf("create storage: %v", err)
	}
	repository := artworkrepo.NewRepo(database)
	log := slog.New(slog.NewTextHandler(io.Discard, nil))
	uc := artworkusecase.NewUseCase(repository, files, log)
	controller := New(uc, storageConfig)

	created := createArtworkThroughHTTP(t, controller, jpegBytes(t), map[string]string{
		"title":       "Первая работа",
		"year":        "2026",
		"position":    "1",
		"isPublished": "true",
	})
	oldKey := strings.TrimPrefix(created.ImageURL, "/media/")
	assertFileExists(t, mediaRoot, oldKey)

	t.Run("published list", func(t *testing.T) {
		response := httptest.NewRecorder()
		controller.ListPublished(response, httptest.NewRequest(http.MethodGet, "/artworks", nil))
		if response.Code != http.StatusOK {
			t.Fatalf("ListPublished() status = %d, body = %q", response.Code, response.Body.String())
		}
		var artworks []artworkdto.ArtworkResponse
		if err := json.Unmarshal(response.Body.Bytes(), &artworks); err != nil {
			t.Fatalf("decode published list: %v", err)
		}
		if len(artworks) != 1 || artworks[0].ID != created.ID || artworks[0].ImageURL != created.ImageURL {
			t.Fatalf("published artworks = %#v", artworks)
		}
	})

	t.Run("replace image", func(t *testing.T) {
		request := multipartRequest(t, http.MethodPut, "/admin/artworks/"+strconv.FormatInt(created.ID, 10), pngBytes(t), map[string]string{
			"title":       "Обновлённая работа",
			"year":        "2025",
			"position":    "2",
			"isPublished": "true",
		})
		request = withArtworkID(request, created.ID)
		response := httptest.NewRecorder()
		controller.Update(response, request)
		if response.Code != http.StatusNoContent {
			t.Fatalf("Update() status = %d, body = %q", response.Code, response.Body.String())
		}

		if _, err := os.Stat(filepath.Join(mediaRoot, filepath.FromSlash(oldKey))); !os.IsNotExist(err) {
			t.Fatalf("old image still exists or cannot be checked: %v", err)
		}
		updated, err := repository.GetByID(context.Background(), created.ID)
		if err != nil {
			t.Fatalf("GetByID() after HTTP update: %v", err)
		}
		if updated.Title != "Обновлённая работа" || updated.ImageKey == oldKey {
			t.Fatalf("updated artwork = %#v", updated)
		}
		assertFileExists(t, mediaRoot, updated.ImageKey)
		created.ImageURL = "/media/" + updated.ImageKey
	})

	t.Run("delete artwork and image", func(t *testing.T) {
		key := strings.TrimPrefix(created.ImageURL, "/media/")
		request := withArtworkID(
			httptest.NewRequest(http.MethodDelete, "/admin/artworks/"+strconv.FormatInt(created.ID, 10), nil),
			created.ID,
		)
		response := httptest.NewRecorder()
		controller.Delete(response, request)
		if response.Code != http.StatusNoContent {
			t.Fatalf("Delete() status = %d, body = %q", response.Code, response.Body.String())
		}
		if _, err := os.Stat(filepath.Join(mediaRoot, filepath.FromSlash(key))); !os.IsNotExist(err) {
			t.Fatalf("deleted image still exists or cannot be checked: %v", err)
		}
		artworks, err := repository.ListAll(context.Background())
		if err != nil || len(artworks) != 0 {
			t.Fatalf("ListAll() after delete = %#v, error = %v", artworks, err)
		}
	})
}

func createArtworkThroughHTTP(t *testing.T, controller *Controller, content []byte, fields map[string]string) artworkdto.AdminArtworkResponse {
	t.Helper()
	request := multipartRequest(t, http.MethodPost, "/admin/artworks", content, fields)
	response := httptest.NewRecorder()
	controller.Create(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("Create() status = %d, body = %q", response.Code, response.Body.String())
	}
	var created artworkdto.AdminArtworkResponse
	if err := json.Unmarshal(response.Body.Bytes(), &created); err != nil {
		t.Fatalf("decode create response: %v", err)
	}
	return created
}

func multipartRequest(t *testing.T, method, target string, content []byte, fields map[string]string) *http.Request {
	t.Helper()
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	for key, value := range fields {
		if err := w.WriteField(key, value); err != nil {
			t.Fatalf("write multipart field: %v", err)
		}
	}
	file, err := w.CreateFormFile("image", "artwork")
	if err != nil {
		t.Fatalf("create multipart image: %v", err)
	}
	if _, err := file.Write(content); err != nil {
		t.Fatalf("write multipart image: %v", err)
	}
	if err := w.Close(); err != nil {
		t.Fatalf("close multipart writer: %v", err)
	}
	request := httptest.NewRequest(method, target, &body)
	request.Header.Set("Content-Type", w.FormDataContentType())
	return request
}

func withArtworkID(request *http.Request, id int64) *http.Request {
	routeContext := chi.NewRouteContext()
	routeContext.URLParams.Add("id", strconv.FormatInt(id, 10))
	return request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeContext))
}

func openIntegrationDatabase(t *testing.T) *sql.DB {
	t.Helper()
	databaseURL := os.Getenv("ARTFOLIO_TEST_DATABASE_URL")
	if databaseURL == "" {
		databaseURL = integrationDatabaseURL
	}
	database, err := sql.Open("pgx", databaseURL)
	if err != nil {
		t.Fatalf("open integration database: %v", err)
	}
	t.Cleanup(func() { _ = database.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := database.PingContext(ctx); err != nil {
		t.Fatalf("connect to integration database: %v", err)
	}
	lockConnection, err := database.Conn(ctx)
	if err != nil {
		t.Fatalf("reserve integration connection: %v", err)
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
	migrations := filepath.Join("..", "..", "..", "..", "migrations")
	if err := goose.Up(database, migrations); err != nil {
		t.Fatalf("apply migrations: %v", err)
	}
	if _, err := database.ExecContext(ctx, "TRUNCATE artworks RESTART IDENTITY"); err != nil {
		t.Fatalf("truncate artworks: %v", err)
	}
	return database
}

func jpegBytes(t *testing.T) []byte {
	t.Helper()
	return encodedImage(t, func(w io.Writer, img image.Image) error {
		return jpeg.Encode(w, img, nil)
	})
}

func pngBytes(t *testing.T) []byte {
	t.Helper()
	return encodedImage(t, png.Encode)
}

func encodedImage(t *testing.T, encode func(io.Writer, image.Image) error) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 8, 8))
	img.Set(0, 0, color.White)
	var buffer bytes.Buffer
	if err := encode(&buffer, img); err != nil {
		t.Fatalf("encode test image: %v", err)
	}
	return buffer.Bytes()
}

func assertFileExists(t *testing.T, root, key string) {
	t.Helper()
	info, err := os.Stat(filepath.Join(root, filepath.FromSlash(key)))
	if err != nil {
		t.Fatalf("stored image %q: %v", key, err)
	}
	if info.IsDir() {
		t.Fatalf("stored image %q is a directory", key)
	}
}
