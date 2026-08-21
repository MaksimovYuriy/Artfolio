export interface SocialLink {
  label: string
  url: string
}

export interface ArtistProfile {
  brandName: string
  name: string
  tagline: string
  biography: string
  artistStatement?: string
  avatarUrl?: string
  email?: string
  telegram?: string
  instagram?: string
  socialLinks?: SocialLink[]
}
