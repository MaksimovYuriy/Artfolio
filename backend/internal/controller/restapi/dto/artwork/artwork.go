package artwork

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
)

var ErrInvalidArtworkForm = errors.New("invalid artwork form")

type ArtworkRequest struct {
	Title       string
	Description string
	Technique   string
	Year        *int16
	ImageAlt    string
	Position    int
	IsPublished bool
}

func ArtworkRequestFromValues(values url.Values) (ArtworkRequest, error) {
	allowed := map[string]bool{
		"title": true, "description": true, "technique": true, "year": true,
		"imageAlt": true, "position": true, "isPublished": true,
	}
	for key := range values {
		if !allowed[key] {
			return ArtworkRequest{}, fmt.Errorf("%w: unknown field %q", ErrInvalidArtworkForm, key)
		}
	}

	title, err := singleValue(values, "title")
	if err != nil {
		return ArtworkRequest{}, err
	}
	description, err := singleValue(values, "description")
	if err != nil {
		return ArtworkRequest{}, err
	}
	technique, err := singleValue(values, "technique")
	if err != nil {
		return ArtworkRequest{}, err
	}
	imageAlt, err := singleValue(values, "imageAlt")
	if err != nil {
		return ArtworkRequest{}, err
	}

	yearValue, err := singleValue(values, "year")
	if err != nil {
		return ArtworkRequest{}, err
	}
	var year *int16
	if yearValue != "" {
		parsed, err := strconv.ParseInt(yearValue, 10, 16)
		if err != nil {
			return ArtworkRequest{}, fmt.Errorf("%w: invalid year", ErrInvalidArtworkForm)
		}
		value := int16(parsed)
		year = &value
	}

	positionValue, err := singleValue(values, "position")
	if err != nil {
		return ArtworkRequest{}, err
	}
	position := 0
	if positionValue != "" {
		parsed, err := strconv.ParseInt(positionValue, 10, 32)
		if err != nil {
			return ArtworkRequest{}, fmt.Errorf("%w: invalid position", ErrInvalidArtworkForm)
		}
		position = int(parsed)
	}

	publishedValue, err := singleValue(values, "isPublished")
	if err != nil {
		return ArtworkRequest{}, err
	}
	published := false
	if publishedValue != "" {
		published, err = strconv.ParseBool(publishedValue)
		if err != nil {
			return ArtworkRequest{}, fmt.Errorf("%w: invalid publication status", ErrInvalidArtworkForm)
		}
	}

	return ArtworkRequest{
		Title:       title,
		Description: description,
		Technique:   technique,
		Year:        year,
		ImageAlt:    imageAlt,
		Position:    position,
		IsPublished: published,
	}, nil
}

func (r ArtworkRequest) Artwork(id int64) entity.Artwork {
	return entity.Artwork{
		ID:          id,
		Title:       r.Title,
		Description: r.Description,
		Technique:   r.Technique,
		Year:        r.Year,
		ImageAlt:    r.ImageAlt,
		Position:    r.Position,
		IsPublished: r.IsPublished,
	}
}

type ArtworkResponse struct {
	ID          int64  `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description,omitempty"`
	Technique   string `json:"technique,omitempty"`
	Year        *int16 `json:"year,omitempty"`
	ImageURL    string `json:"imageUrl"`
	ImageAlt    string `json:"imageAlt,omitempty"`
}

type AdminArtworkResponse struct {
	ArtworkResponse
	Position    int       `json:"position"`
	IsPublished bool      `json:"isPublished"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

func ArtworkResponseFromEntity(artwork entity.Artwork, publicURL string) ArtworkResponse {
	return ArtworkResponse{
		ID:          artwork.ID,
		Title:       artwork.Title,
		Description: artwork.Description,
		Technique:   artwork.Technique,
		Year:        artwork.Year,
		ImageURL:    imageURL(publicURL, artwork.ImageKey),
		ImageAlt:    artwork.ImageAlt,
	}
}

func ArtworkResponsesFromEntities(artworks []entity.Artwork, publicURL string) []ArtworkResponse {
	responses := make([]ArtworkResponse, 0, len(artworks))
	for _, artwork := range artworks {
		responses = append(responses, ArtworkResponseFromEntity(artwork, publicURL))
	}
	return responses
}

func AdminArtworkResponseFromEntity(artwork entity.Artwork, publicURL string) AdminArtworkResponse {
	return AdminArtworkResponse{
		ArtworkResponse: ArtworkResponseFromEntity(artwork, publicURL),
		Position:        artwork.Position,
		IsPublished:     artwork.IsPublished,
		CreatedAt:       artwork.CreatedAt,
		UpdatedAt:       artwork.UpdatedAt,
	}
}

func AdminArtworkResponsesFromEntities(artworks []entity.Artwork, publicURL string) []AdminArtworkResponse {
	responses := make([]AdminArtworkResponse, 0, len(artworks))
	for _, artwork := range artworks {
		responses = append(responses, AdminArtworkResponseFromEntity(artwork, publicURL))
	}
	return responses
}

func singleValue(values url.Values, key string) (string, error) {
	items := values[key]
	if len(items) > 1 {
		return "", fmt.Errorf("%w: field %q must occur once", ErrInvalidArtworkForm, key)
	}
	if len(items) == 0 {
		return "", nil
	}
	return items[0], nil
}

func imageURL(publicURL, key string) string {
	return strings.TrimRight(publicURL, "/") + "/" + strings.TrimLeft(key, "/")
}
