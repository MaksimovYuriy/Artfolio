import { mockArtworks } from '../data/artworks'
import type { Artwork } from '../types/artwork'

// Keep this async contract when replacing the implementation with GET /api/artworks.
export async function getArtworks(): Promise<Artwork[]> {
  return mockArtworks
}
