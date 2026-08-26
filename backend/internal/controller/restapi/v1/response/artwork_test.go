package response

import (
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
)

func TestArtworkFromEntityBuildsPublicURL(t *testing.T) {
	mapper := NewArtworkMapper("/media/")
	response := mapper.FromEntity(entity.Artwork{ImageKey: "artworks/image.jpg"})
	if response.ImageURL != "/media/artworks/image.jpg" {
		t.Fatalf("ImageURL = %q", response.ImageURL)
	}
}
