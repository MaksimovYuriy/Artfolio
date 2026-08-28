package request

import (
	"errors"
	"net/url"
	"testing"
)

func TestArtworkFromValues(t *testing.T) {
	request, err := ArtworkFromValues(url.Values{
		"title":       {"Работа"},
		"year":        {"2026"},
		"isPublished": {"true"},
	})
	if err != nil {
		t.Fatalf("ArtworkFromValues() error = %v", err)
	}
	if request.Year == nil || *request.Year != 2026 || !request.IsPublished {
		t.Fatalf("ArtworkFromValues() = %#v", request)
	}
}

func TestArtworkFromValuesRejectsUnknownAndRepeatedFields(t *testing.T) {
	tests := []url.Values{
		{"unknown": {"value"}},
		{"title": {"one", "two"}},
	}
	for _, values := range tests {
		if _, err := ArtworkFromValues(values); !errors.Is(err, ErrInvalidArtworkForm) {
			t.Fatalf("ArtworkFromValues(%v) error = %v", values, err)
		}
	}
}
