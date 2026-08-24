package artwork

import (
	"errors"
	"net/url"
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
)

func TestArtworkRequestFromValues(t *testing.T) {
	request, err := ArtworkRequestFromValues(url.Values{
		"title":       {"Работа"},
		"year":        {"2026"},
		"position":    {"3"},
		"isPublished": {"true"},
	})
	if err != nil {
		t.Fatalf("ArtworkRequestFromValues() error = %v", err)
	}
	if request.Year == nil || *request.Year != 2026 || request.Position != 3 || !request.IsPublished {
		t.Fatalf("ArtworkRequestFromValues() = %#v", request)
	}
}

func TestArtworkRequestFromValuesRejectsUnknownAndRepeatedFields(t *testing.T) {
	tests := []url.Values{
		{"unknown": {"value"}},
		{"title": {"one", "two"}},
	}
	for _, values := range tests {
		if _, err := ArtworkRequestFromValues(values); !errors.Is(err, ErrInvalidArtworkForm) {
			t.Fatalf("ArtworkRequestFromValues(%v) error = %v", values, err)
		}
	}
}

func TestArtworkResponseBuildsPublicURL(t *testing.T) {
	response := ArtworkResponseFromEntity(entity.Artwork{ImageKey: "artworks/image.jpg"}, "/media/")
	if response.ImageURL != "/media/artworks/image.jpg" {
		t.Fatalf("ImageURL = %q", response.ImageURL)
	}
}
