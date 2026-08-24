package restapi

import (
	"net/http"

	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

type ArtworkController struct {
	usecase usecase.ArtworkUseCase
}

func NewArtworkController(usecase usecase.ArtworkUseCase) *ArtworkController {
	return &ArtworkController{usecase: usecase}
}

func (c *ArtworkController) ListPublished(w http.ResponseWriter, r *http.Request) {
	// Implement
}

func (c *ArtworkController) ListAll(w http.ResponseWriter, r *http.Request) {
	// Implement
}

func (c *ArtworkController) Create(w http.ResponseWriter, r *http.Request) {
	// Implement
}

func (c *ArtworkController) Update(w http.ResponseWriter, r *http.Request) {
	// Implement
}

func (c *ArtworkController) Delete(w http.ResponseWriter, r *http.Request) {
	// Implement
}
