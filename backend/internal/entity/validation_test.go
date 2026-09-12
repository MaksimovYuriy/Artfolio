package entity

import (
	"errors"
	"testing"
)

func TestValidatedNormalizesEntities(t *testing.T) {
	profile, err := (ArtistProfile{Name: "  Тестовый автор  "}).Validated()
	if err != nil || profile.Name != "Тестовый автор" {
		t.Fatalf("ArtistProfile.Validated() = %#v, %v", profile, err)
	}

	artwork, err := (Artwork{Title: "  Работа  "}).Validated()
	if err != nil || artwork.Title != "Работа" {
		t.Fatalf("Artwork.Validated() = %#v, %v", artwork, err)
	}

	link, err := (SocialLink{Platform: SocialPlatformTelegram, Handle: " https://t.me/demo_author "}).Validated()
	if err != nil || link.Handle != "demo_author" {
		t.Fatalf("SocialLink.Validated() = %#v, %v", link, err)
	}
}

func TestValidatedRejectsBusinessRuleViolation(t *testing.T) {
	_, err := (Artwork{Title: "   "}).Validated()
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("Artwork.Validated() error = %v, want ErrValidation", err)
	}

	email := "художник@example.com"
	_, err = (ArtistProfile{Name: "Тестовый автор", Email: &email}).Validated()
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("ArtistProfile.Validated() error = %v, want ErrValidation", err)
	}
}
