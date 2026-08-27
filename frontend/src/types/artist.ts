import type { SocialLink } from './socialLink'

export interface ArtistProfile {
  brandName: string
  name: string
  tagline: string
  biography: string
  artistStatement?: string
  email?: string
  socialLinks?: SocialLink[]
}

export interface EditableArtistProfile {
  name: string
  tagline: string
  biography: string
  artistStatement: string
  email: string
}
