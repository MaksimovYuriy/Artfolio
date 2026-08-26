package response

import (
	"strings"
	"time"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
)

type Artwork struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Technique   string `json:"technique,omitempty"`
	Year        *int16 `json:"year,omitempty"`
	ImageURL    string `json:"imageUrl"`
	ImageAlt    string `json:"imageAlt,omitempty"`
}

type AdminArtwork struct {
	Artwork
	Position    int       `json:"position"`
	IsPublished bool      `json:"isPublished"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ArtworkMapper struct {
	mediaPublicURL string
}

func NewArtworkMapper(mediaPublicURL string) ArtworkMapper {
	return ArtworkMapper{mediaPublicURL: mediaPublicURL}
}

func (m ArtworkMapper) FromEntity(artwork entity.Artwork) Artwork {
	return Artwork{
		ID:          artwork.ID,
		Title:       artwork.Title,
		Description: artwork.Description,
		Technique:   artwork.Technique,
		Year:        artwork.Year,
		ImageURL:    imageURL(m.mediaPublicURL, artwork.ImageKey),
		ImageAlt:    artwork.ImageAlt,
	}
}

func (m ArtworkMapper) FromEntities(artworks []entity.Artwork) []Artwork {
	responses := make([]Artwork, 0, len(artworks))
	for _, artwork := range artworks {
		responses = append(responses, m.FromEntity(artwork))
	}
	return responses
}

func (m ArtworkMapper) AdminFromEntity(artwork entity.Artwork) AdminArtwork {
	return AdminArtwork{
		Artwork:     m.FromEntity(artwork),
		Position:    artwork.Position,
		IsPublished: artwork.IsPublished,
		CreatedAt:   artwork.CreatedAt,
		UpdatedAt:   artwork.UpdatedAt,
	}
}

func (m ArtworkMapper) AdminFromEntities(artworks []entity.Artwork) []AdminArtwork {
	responses := make([]AdminArtwork, 0, len(artworks))
	for _, artwork := range artworks {
		responses = append(responses, m.AdminFromEntity(artwork))
	}
	return responses
}

func imageURL(publicURL, key string) string {
	return strings.TrimRight(publicURL, "/") + "/" + strings.TrimLeft(key, "/")
}
