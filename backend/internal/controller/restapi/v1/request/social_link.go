package request

import "github.com/maksimovyuriy/artfolio/backend/internal/entity"

type SocialLink struct {
	Platform entity.SocialPlatform `json:"platform"`
	Handle   string                `json:"handle"`
}

type ReplaceSocialLinks struct {
	SocialLinks []SocialLink `json:"socialLinks"`
}

func (r ReplaceSocialLinks) Entities() []entity.SocialLink {
	links := make([]entity.SocialLink, 0, len(r.SocialLinks))
	for _, link := range r.SocialLinks {
		links = append(links, entity.SocialLink{
			Platform: link.Platform,
			Handle:   link.Handle,
		})
	}
	return links
}
