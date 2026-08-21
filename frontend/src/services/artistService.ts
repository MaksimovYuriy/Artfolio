import { mockArtist } from '../data/artist'
import type { ArtistProfile } from '../types/artist'

// Keep this async contract when replacing the implementation with GET /api/artist.
export async function getArtist(): Promise<ArtistProfile> {
  return mockArtist
}
