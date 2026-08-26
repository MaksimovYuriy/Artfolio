package artistprofile

import (
	"context"
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
)

func TestGetIncludesSocialLinks(t *testing.T) {
	profileRepo := &profileRepositoryStub{profile: entity.ArtistProfile{ID: 12, Name: "Анна"}}
	linkRepo := &socialLinkRepositoryStub{links: []entity.SocialLink{
		{ArtistProfileID: 12, Platform: entity.SocialPlatformTelegram, Handle: "anna_art"},
	}}
	uc := NewUseCase(profileRepo, linkRepo)

	profile, err := uc.Get(context.Background())
	if err != nil {
		t.Fatalf("Get() error = %v", err)
	}
	if linkRepo.profileID != 12 || len(profile.SocialLinks) != 1 {
		t.Fatalf("Get() profile = %#v, social profile ID = %d", profile, linkRepo.profileID)
	}
}

type profileRepositoryStub struct {
	profile entity.ArtistProfile
}

func (r *profileRepositoryStub) Get(context.Context) (entity.ArtistProfile, error) {
	return r.profile, nil
}

func (r *profileRepositoryStub) Update(context.Context, entity.ArtistProfile) error {
	return nil
}

type socialLinkRepositoryStub struct {
	links     []entity.SocialLink
	profileID int64
}

func (r *socialLinkRepositoryStub) List(_ context.Context, artistProfileID int64) ([]entity.SocialLink, error) {
	r.profileID = artistProfileID
	return r.links, nil
}

func (r *socialLinkRepositoryStub) Replace(context.Context, int64, []entity.SocialLink) error {
	return nil
}
