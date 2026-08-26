package entity

type SocialPlatform string

const (
	SocialPlatformTelegram  SocialPlatform = "telegram"
	SocialPlatformInstagram SocialPlatform = "instagram"
	SocialPlatformVK        SocialPlatform = "vk"
	SocialPlatformBehance   SocialPlatform = "behance"
)

type SocialLink struct {
	ArtistProfileID int64
	Platform        SocialPlatform
	Handle          string
}
