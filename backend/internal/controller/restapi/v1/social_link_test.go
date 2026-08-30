package v1

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/maksimovyuriy/artfolio/backend/internal/controller/restapi/v1/response"
	"github.com/maksimovyuriy/artfolio/backend/internal/entity"
)

func TestSocialLinkHandlers(t *testing.T) {
	uc := &socialLinkUseCaseStub{links: []entity.SocialLink{
		{Platform: entity.SocialPlatformTelegram, Handle: "anna_art"},
	}}
	controller := NewController(nil, nil, nil, uc, nil, response.NewArtworkMapper("/media"))

	t.Run("list", func(t *testing.T) {
		result := httptest.NewRecorder()
		controller.listSocialLinks(result, httptest.NewRequest(http.MethodGet, "/admin/social_links", nil))
		if result.Code != http.StatusOK {
			t.Fatalf("status = %d, body = %q", result.Code, result.Body.String())
		}
		var links []response.AdminSocialLink
		if err := json.Unmarshal(result.Body.Bytes(), &links); err != nil {
			t.Fatalf("decode response: %v", err)
		}
		if len(links) != 1 || links[0].Handle != "anna_art" {
			t.Fatalf("links = %#v", links)
		}
	})

	t.Run("replace", func(t *testing.T) {
		body := bytes.NewBufferString(`{"socialLinks":[{"platform":"vk","handle":"anna"}]}`)
		result := httptest.NewRecorder()
		controller.replaceSocialLinks(result, httptest.NewRequest(http.MethodPut, "/admin/social_links", body))
		if result.Code != http.StatusNoContent {
			t.Fatalf("status = %d, body = %q", result.Code, result.Body.String())
		}
		if len(uc.replaced) != 1 || uc.replaced[0].Platform != entity.SocialPlatformVK {
			t.Fatalf("replaced = %#v", uc.replaced)
		}
	})
}

type socialLinkUseCaseStub struct {
	links    []entity.SocialLink
	replaced []entity.SocialLink
}

func (u *socialLinkUseCaseStub) List(context.Context) ([]entity.SocialLink, error) {
	return u.links, nil
}

func (u *socialLinkUseCaseStub) Replace(_ context.Context, links []entity.SocialLink) error {
	u.replaced = append([]entity.SocialLink(nil), links...)
	return nil
}
