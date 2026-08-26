package entity

type ArtistProfile struct {
	ID              int64
	Name            string
	Tagline         string
	Biography       string
	ArtistStatement *string
	AvatarURL       *string
	Email           *string
	SocialLinks     []SocialLink
}
