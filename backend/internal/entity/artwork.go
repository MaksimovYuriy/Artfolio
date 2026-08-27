package entity

import (
	"strings"
	"time"
)

type Artwork struct {
	ID          int64
	Title       string
	Description string
	Technique   string
	Year        *int16
	ImageKey    string
	ImageAlt    string
	Position    int
	IsPublished bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type StoredArtworkImage struct {
	Key    string
	Width  int
	Height int
}

func (artwork Artwork) Validated() (Artwork, error) {
	artwork.normalize()

	if artwork.Title == "" {
		return Artwork{}, NewValidationError("title", "is required")
	}
	if artwork.Year != nil && (*artwork.Year < 0 || *artwork.Year > 9999) {
		return Artwork{}, NewValidationError("year", "must be between 0 and 9999")
	}
	if artwork.Position < 0 {
		return Artwork{}, NewValidationError("position", "must not be negative")
	}
	return artwork, nil
}

func (artwork *Artwork) normalize() {
	artwork.Title = strings.TrimSpace(artwork.Title)
	artwork.Description = strings.TrimSpace(artwork.Description)
	artwork.Technique = strings.TrimSpace(artwork.Technique)
	artwork.ImageAlt = strings.TrimSpace(artwork.ImageAlt)
}
