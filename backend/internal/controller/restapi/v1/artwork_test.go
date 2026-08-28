package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/response"
	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

func TestArtworkControllerCreate(t *testing.T) {
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	_ = w.WriteField("title", "Работа")
	_ = w.WriteField("isPublished", "true")
	file, err := w.CreateFormFile("image", "work.jpg")
	if err != nil {
		t.Fatalf("CreateFormFile() error = %v", err)
	}
	_, _ = file.Write([]byte("image content"))
	_ = w.Close()

	uc := &fakeArtworkUseCase{}
	controller := testArtworkController(uc)
	request := httptest.NewRequest(http.MethodPost, "/admin/artworks", &body)
	request.Header.Set("Content-Type", w.FormDataContentType())
	response := httptest.NewRecorder()

	controller.createArtwork(response, request)

	if response.Code != http.StatusCreated {
		t.Fatalf("Create() status = %d, body = %q", response.Code, response.Body.String())
	}
	if uc.created.Title != "Работа" || !uc.created.IsPublished || uc.image == "" {
		t.Fatalf("Create() use case input = %#v, image = %q", uc.created, uc.image)
	}
}

func TestArtworkControllerCreateRequiresMultipartAndImage(t *testing.T) {
	controller := testArtworkController(&fakeArtworkUseCase{})

	plainRequest := httptest.NewRequest(http.MethodPost, "/admin/artworks", strings.NewReader("title=work"))
	plainRequest.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	plainResponse := httptest.NewRecorder()
	controller.createArtwork(plainResponse, plainRequest)
	if plainResponse.Code != http.StatusUnsupportedMediaType {
		t.Fatalf("Create() plain status = %d", plainResponse.Code)
	}

	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	_ = w.WriteField("title", "Работа")
	_ = w.Close()
	request := httptest.NewRequest(http.MethodPost, "/admin/artworks", &body)
	request.Header.Set("Content-Type", w.FormDataContentType())
	response := httptest.NewRecorder()
	controller.createArtwork(response, request)
	if response.Code != http.StatusBadRequest {
		t.Fatalf("Create() missing image status = %d", response.Code)
	}
}

func TestArtworkControllerListPublishedUsesPublicDTO(t *testing.T) {
	uc := &fakeArtworkUseCase{published: []entity.Artwork{{
		ID: 1, Title: "Работа", ImageKey: "artworks/work.jpg", IsPublished: true,
	}}}
	controller := testArtworkController(uc)
	response := httptest.NewRecorder()

	controller.listPublishedArtworks(response, httptest.NewRequest(http.MethodGet, "/artworks", nil))

	if response.Code != http.StatusOK {
		t.Fatalf("ListPublished() status = %d", response.Code)
	}
	var payload []map[string]any
	if err := json.Unmarshal(response.Body.Bytes(), &payload); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if payload[0]["imageUrl"] != "/media/artworks/work.jpg" {
		t.Fatalf("imageUrl = %#v", payload[0]["imageUrl"])
	}
	if _, exists := payload[0]["isPublished"]; exists {
		t.Fatal("public response contains isPublished")
	}
}

func TestArtworkControllerReorder(t *testing.T) {
	uc := &fakeArtworkUseCase{}
	controller := testArtworkController(uc)
	request := httptest.NewRequest(http.MethodPut, "/admin/artworks/order", strings.NewReader(`{"artworkIds":[3,1,2]}`))
	response := httptest.NewRecorder()

	controller.reorderArtworks(response, request)

	if response.Code != http.StatusNoContent {
		t.Fatalf("Reorder() status = %d, body = %q", response.Code, response.Body.String())
	}
	if len(uc.reordered) != 3 || uc.reordered[0] != 3 || uc.reordered[2] != 2 {
		t.Fatalf("Reorder() ids = %v", uc.reordered)
	}
}

func TestArtworkControllerReorderConflict(t *testing.T) {
	uc := &fakeArtworkUseCase{reorderErr: usecase.ErrArtworkOrderConflict}
	controller := testArtworkController(uc)
	request := httptest.NewRequest(http.MethodPut, "/admin/artworks/order", strings.NewReader(`{"artworkIds":[1]}`))
	response := httptest.NewRecorder()

	controller.reorderArtworks(response, request)

	if response.Code != http.StatusConflict {
		t.Fatalf("Reorder() status = %d, want 409", response.Code)
	}
}

func testArtworkController(uc *fakeArtworkUseCase) *Controller {
	return NewController(nil, nil, uc, nil, response.NewArtworkMapper("/media"))
}

type fakeArtworkUseCase struct {
	published  []entity.Artwork
	created    entity.Artwork
	image      string
	reordered  []int64
	reorderErr error
}

func (u *fakeArtworkUseCase) ListPublished(context.Context) ([]entity.Artwork, error) {
	return u.published, nil
}

func (u *fakeArtworkUseCase) ListAll(context.Context) ([]entity.Artwork, error) {
	return []entity.Artwork{}, nil
}

func (u *fakeArtworkUseCase) Create(_ context.Context, artwork entity.Artwork, image io.Reader) (entity.Artwork, error) {
	u.created = artwork
	content, _ := io.ReadAll(image)
	u.image = string(content)
	artwork.ID = 1
	artwork.ImageKey = "artworks/work.jpg"
	return artwork, nil
}

func (u *fakeArtworkUseCase) Update(context.Context, entity.Artwork, io.Reader) error {
	return nil
}

func (u *fakeArtworkUseCase) Delete(context.Context, int64) error {
	return nil
}

func (u *fakeArtworkUseCase) Reorder(_ context.Context, ids []int64) error {
	u.reordered = ids
	return u.reorderErr
}
