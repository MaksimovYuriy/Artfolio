import type { ArtistProfile } from '../types/artist'
import { APIClientError, apiRequestJSON } from './apiClient'

interface ArtistProfileResponse {
  name: string
  tagline: string
  biography: string
  artistStatement?: string | null
  email?: string | null
  socialLinks?: Array<{
    label: string
    url: string
  }>
}

export class ArtistServiceError extends Error {
  constructor(message: string) {
    super(message)
    this.name = 'ArtistServiceError'
  }
}

export async function getArtist(): Promise<ArtistProfile> {
  let profile: ArtistProfileResponse
  try {
    profile = await apiRequestJSON<ArtistProfileResponse>('/api/v1/artist_profile')
  } catch (error) {
    if (error instanceof APIClientError) {
      if (error.status === null) throw new ArtistServiceError('Не удалось связаться с сервером.')
      if (error.status === 404) throw new ArtistServiceError('Профиль художницы ещё не заполнен.')
      if (error.code === 'invalid_response') {
        throw new ArtistServiceError('Сервер вернул некорректные данные профиля.')
      }
    }
    throw new ArtistServiceError('Не удалось загрузить профиль.')
  }

  return {
    brandName: 'Artfolio',
    name: profile.name,
    tagline: profile.tagline,
    biography: profile.biography,
    artistStatement: profile.artistStatement ?? undefined,
    email: profile.email ?? undefined,
    socialLinks: profile.socialLinks ?? [],
  }
}
