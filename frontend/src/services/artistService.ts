import type { ArtistProfile } from '../types/artist'

interface ArtistProfileResponse {
  name: string
  tagline: string
  biography: string
  artistStatement?: string | null
  avatarUrl?: string | null
  email?: string | null
}

export class ArtistServiceError extends Error {
  constructor(message: string) {
    super(message)
    this.name = 'ArtistServiceError'
  }
}

export async function getArtist(): Promise<ArtistProfile> {
  let response: Response

  try {
    response = await fetch('/api/v1/artist_profile')
  } catch {
    throw new ArtistServiceError('Не удалось связаться с сервером.')
  }

  if (response.status === 404) {
    throw new ArtistServiceError('Профиль художницы ещё не заполнен.')
  }
  if (!response.ok) {
    throw new ArtistServiceError('Не удалось загрузить профиль.')
  }

  let profile: ArtistProfileResponse
  try {
    profile = await response.json() as ArtistProfileResponse
  } catch {
    throw new ArtistServiceError('Сервер вернул некорректные данные профиля.')
  }

  return {
    brandName: 'Artfolio',
    name: profile.name,
    tagline: profile.tagline,
    biography: profile.biography,
    artistStatement: profile.artistStatement ?? undefined,
    avatarUrl: profile.avatarUrl ?? undefined,
    email: profile.email ?? undefined,
  }
}
