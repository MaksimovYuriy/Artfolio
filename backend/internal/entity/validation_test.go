package entity

import (
	"errors"
	"testing"
)

func TestValidatedNormalizesEntities(t *testing.T) {
	profile, err := (ArtistProfile{Name: "  Анна  "}).Validated()
	if err != nil || profile.Name != "Анна" {
		t.Fatalf("ArtistProfile.Validated() = %#v, %v", profile, err)
	}

	artwork, err := (Artwork{Title: "  Работа  "}).Validated()
	if err != nil || artwork.Title != "Работа" {
		t.Fatalf("Artwork.Validated() = %#v, %v", artwork, err)
	}

	link, err := (SocialLink{Platform: SocialPlatformTelegram, Handle: " https://t.me/anna_art "}).Validated()
	if err != nil || link.Handle != "anna_art" {
		t.Fatalf("SocialLink.Validated() = %#v, %v", link, err)
	}
}

func TestValidatedRejectsBusinessRuleViolation(t *testing.T) {
	_, err := (Artwork{Title: "   "}).Validated()
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("Artwork.Validated() error = %v, want ErrValidation", err)
	}

	email := "художник@example.com"
	_, err = (ArtistProfile{Name: "Анна", Email: &email}).Validated()
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("ArtistProfile.Validated() error = %v, want ErrValidation", err)
	}
}
