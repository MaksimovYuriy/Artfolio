package request

import (
	"errors"
	"fmt"
	"net/url"
	"strconv"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
)

var ErrInvalidArtworkForm = errors.New("invalid artwork form")

type Artwork struct {
	Title       string
	Description string
	Technique   string
	Year        *int16
	ImageAlt    string
	IsPublished bool
}

func ArtworkFromValues(values url.Values) (Artwork, error) {
	allowed := map[string]bool{
		"title": true, "description": true, "technique": true, "year": true,
		"imageAlt": true, "isPublished": true,
	}
	for key := range values {
		if !allowed[key] {
			return Artwork{}, fmt.Errorf("%w: unknown field %q", ErrInvalidArtworkForm, key)
		}
	}

	title, err := singleValue(values, "title")
	if err != nil {
		return Artwork{}, err
	}
	description, err := singleValue(values, "description")
	if err != nil {
		return Artwork{}, err
	}
	technique, err := singleValue(values, "technique")
	if err != nil {
		return Artwork{}, err
	}
	imageAlt, err := singleValue(values, "imageAlt")
	if err != nil {
		return Artwork{}, err
	}

	yearValue, err := singleValue(values, "year")
	if err != nil {
		return Artwork{}, err
	}
	var year *int16
	if yearValue != "" {
		parsed, err := strconv.ParseInt(yearValue, 10, 16)
		if err != nil {
			return Artwork{}, fmt.Errorf("%w: invalid year", ErrInvalidArtworkForm)
		}
		value := int16(parsed)
		year = &value
	}

	publishedValue, err := singleValue(values, "isPublished")
	if err != nil {
		return Artwork{}, err
	}
	published := false
	if publishedValue != "" {
		published, err = strconv.ParseBool(publishedValue)
		if err != nil {
			return Artwork{}, fmt.Errorf("%w: invalid publication status", ErrInvalidArtworkForm)
		}
	}

	return Artwork{
		Title:       title,
		Description: description,
		Technique:   technique,
		Year:        year,
		ImageAlt:    imageAlt,
		IsPublished: published,
	}, nil
}

func (r Artwork) Artwork(id int64) entity.Artwork {
	return entity.Artwork{
		ID:          id,
		Title:       r.Title,
		Description: r.Description,
		Technique:   r.Technique,
		Year:        r.Year,
		ImageAlt:    r.ImageAlt,
		IsPublished: r.IsPublished,
	}
}

type ReorderArtworks struct {
	ArtworkIDs []int64 `json:"artworkIds"`
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
