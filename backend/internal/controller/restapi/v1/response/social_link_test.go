package response

import (
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
)

func TestSocialLinksFromEntitiesBuildsPublicURLs(t *testing.T) {
	links := SocialLinksFromEntities([]entity.SocialLink{
		{Platform: entity.SocialPlatformTelegram, Handle: "anna_art"},
		{Platform: entity.SocialPlatformBehance, Handle: "anna-art"},
	})

	if len(links) != 2 {
		t.Fatalf("links = %#v", links)
	}
	if links[0].Label != "Telegram" || links[0].URL != "https://t.me/anna_art" {
		t.Fatalf("telegram link = %#v", links[0])
	}
	if links[1].Label != "Behance" || links[1].URL != "https://behance.net/anna-art" {
		t.Fatalf("behance link = %#v", links[1])
	}
}
