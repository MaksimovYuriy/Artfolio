package response

import "github.com/maksimovyuriy/artfolio/backend/internal/entity"

type SocialLink struct {
	Platform entity.SocialPlatform `json:"platform"`
	Label    string                `json:"label"`
	URL      string                `json:"url"`
}

type AdminSocialLink struct {
	Platform entity.SocialPlatform `json:"platform"`
	Handle   string                `json:"handle"`
}

func SocialLinksFromEntities(links []entity.SocialLink) []SocialLink {
	responses := make([]SocialLink, 0, len(links))
	for _, link := range links {
		responses = append(responses, SocialLink{
			Platform: link.Platform,
			Label:    socialPlatformLabel(link.Platform),
			URL:      socialPlatformBaseURL(link.Platform) + link.Handle,
		})
	}
	return responses
}

func AdminSocialLinksFromEntities(links []entity.SocialLink) []AdminSocialLink {
	responses := make([]AdminSocialLink, 0, len(links))
	for _, link := range links {
		responses = append(responses, AdminSocialLink{
			Platform: link.Platform,
			Handle:   link.Handle,
		})
	}
	return responses
}

func socialPlatformLabel(platform entity.SocialPlatform) string {
	switch platform {
	case entity.SocialPlatformTelegram:
		return "Telegram"
	case entity.SocialPlatformInstagram:
		return "Instagram"
	case entity.SocialPlatformVK:
		return "VK"
	case entity.SocialPlatformBehance:
		return "Behance"
	default:
		return ""
	}
}

func socialPlatformBaseURL(platform entity.SocialPlatform) string {
	switch platform {
	case entity.SocialPlatformTelegram:
		return "https://t.me/"
	case entity.SocialPlatformInstagram:
		return "https://instagram.com/"
	case entity.SocialPlatformVK:
		return "https://vk.com/"
	case entity.SocialPlatformBehance:
		return "https://behance.net/"
	default:
		return ""
	}
}
