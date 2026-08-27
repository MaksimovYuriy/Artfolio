export interface Artwork {
  id: number
  title: string
  imageUrl?: string
  imageAlt?: string
  year?: number
  technique?: string
  description?: string
}

export interface AdminArtwork extends Artwork {
  position: number
  isPublished: boolean
  createdAt: string
  updatedAt: string
}

export interface ArtworkInput {
  title: string
  description: string
  technique: string
  year: string
  imageAlt: string
  position: number
}
