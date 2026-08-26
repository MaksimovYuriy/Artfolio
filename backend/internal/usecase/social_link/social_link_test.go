package sociallink

import (
	"context"
	"errors"
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
	"github.com/maksimovyuriy/artfolio/backend/internal/usecase"
)

func TestReplaceNormalizesAndSkipsEmptyHandles(t *testing.T) {
	profileRepo := &fakeProfileRepository{profile: entity.ArtistProfile{ID: 7}}
	linkRepo := &fakeSocialLinkRepository{}
	uc := NewUseCase(profileRepo, linkRepo)

	err := uc.Replace(context.Background(), []entity.SocialLink{
		{Platform: entity.SocialPlatformTelegram, Handle: " https://t.me/anna_art "},
		{Platform: entity.SocialPlatformInstagram, Handle: "  @anna.art  "},
		{Platform: entity.SocialPlatformVK, Handle: "   "},
	})
	if err != nil {
		t.Fatalf("Replace() error = %v", err)
	}
	if linkRepo.replacedProfileID != 7 {
		t.Fatalf("Replace() profile ID = %d", linkRepo.replacedProfileID)
	}
	if len(linkRepo.replaced) != 2 {
		t.Fatalf("Replace() links = %#v", linkRepo.replaced)
	}
	if linkRepo.replaced[0].Handle != "anna_art" || linkRepo.replaced[1].Handle != "anna.art" {
		t.Fatalf("Replace() normalized links = %#v", linkRepo.replaced)
	}
}

func TestReplaceRejectsInvalidAndDuplicatePlatforms(t *testing.T) {
	tests := []struct {
		name  string
		links []entity.SocialLink
	}{
		{
			name:  "unsupported platform",
			links: []entity.SocialLink{{Platform: "unknown", Handle: "artist"}},
		},
		{
			name: "duplicate platform",
			links: []entity.SocialLink{
				{Platform: entity.SocialPlatformTelegram, Handle: "artist_one"},
				{Platform: entity.SocialPlatformTelegram, Handle: "artist_two"},
			},
		},
		{
			name:  "foreign URL",
			links: []entity.SocialLink{{Platform: entity.SocialPlatformTelegram, Handle: "https://example.com/artist"}},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uc := NewUseCase(&fakeProfileRepository{profile: entity.ArtistProfile{ID: 1}}, &fakeSocialLinkRepository{})
			if err := uc.Replace(context.Background(), test.links); !errors.Is(err, usecase.ErrInvalidSocialLinks) {
				t.Fatalf("Replace() error = %v", err)
			}
		})
	}
}

func TestListUsesSingletonProfileID(t *testing.T) {
	expected := []entity.SocialLink{{ArtistProfileID: 4, Platform: entity.SocialPlatformVK, Handle: "artist"}}
	linkRepo := &fakeSocialLinkRepository{links: expected}
	uc := NewUseCase(&fakeProfileRepository{profile: entity.ArtistProfile{ID: 4}}, linkRepo)

	links, err := uc.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if linkRepo.listedProfileID != 4 || len(links) != 1 || links[0] != expected[0] {
		t.Fatalf("List() = %#v, profile ID = %d", links, linkRepo.listedProfileID)
	}
}

type fakeProfileRepository struct {
	profile entity.ArtistProfile
	err     error
}

func (r *fakeProfileRepository) Get(context.Context) (entity.ArtistProfile, error) {
	return r.profile, r.err
}

func (r *fakeProfileRepository) Update(context.Context, entity.ArtistProfile) error {
	return nil
}

type fakeSocialLinkRepository struct {
	links             []entity.SocialLink
	listedProfileID   int64
	replacedProfileID int64
	replaced          []entity.SocialLink
	err               error
}

func (r *fakeSocialLinkRepository) List(_ context.Context, artistProfileID int64) ([]entity.SocialLink, error) {
	r.listedProfileID = artistProfileID
	return r.links, r.err
}

func (r *fakeSocialLinkRepository) Replace(_ context.Context, artistProfileID int64, links []entity.SocialLink) error {
	r.replacedProfileID = artistProfileID
	r.replaced = append([]entity.SocialLink(nil), links...)
	return r.err
}
