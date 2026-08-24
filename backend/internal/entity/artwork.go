package entity

import "time"

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
