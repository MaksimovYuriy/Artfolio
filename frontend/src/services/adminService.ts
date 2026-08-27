import type { EditableArtistProfile } from '../types/artist'

export class AdminAuthError extends Error {
  constructor(message: string) {
    super(message)
    this.name = 'AdminAuthError'
  }
}

export class AdminProfileError extends Error {
  unauthorized: boolean

  constructor(message: string, unauthorized = false) {
    super(message)
    this.name = 'AdminProfileError'
    this.unauthorized = unauthorized
  }
}

export const emptyArtistProfile: EditableArtistProfile = {
  name: '',
  tagline: '',
  biography: '',
  artistStatement: '',
  email: '',
}

export async function getEditableArtistProfile(): Promise<EditableArtistProfile> {
  let response: Response

  try {
    response = await fetch('/api/v1/artist_profile', { credentials: 'same-origin' })
  } catch {
    throw new AdminProfileError('Не удалось загрузить профиль. Проверьте соединение.')
  }

  if (response.status === 404) {
    return { ...emptyArtistProfile }
  }
  if (!response.ok) {
    throw new AdminProfileError('Не удалось загрузить прфиль. Попробуйте позже.')
  }

  const profile = await response.json() as Partial<EditableArtistProfile>
  return {
    name: profile.name ?? '',
    tagline: profile.tagline ?? '',
    biography: profile.biography ?? '',
    artistStatement: profile.artistStatement ?? '',
    email: profile.email ?? '',
  }
}

export async function updateArtistProfile(profile: EditableArtistProfile): Promise<void> {
  let response: Response

  try {
    response = await fetch('/api/v1/admin/artist_profile', {
      method: 'PUT',
      credentials: 'same-origin',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: profile.name,
        tagline: profile.tagline,
        biography: profile.biography,
        artistStatement: profile.artistStatement || null,
        email: profile.email || null,
      }),
    })
  } catch {
    throw new AdminProfileError('Не удалось сохранить профиль. Проверьте соединение.')
  }

  if (response.status === 401) {
    throw new AdminProfileError('Сессия завершилась. Войдите снова.', true)
  }
  if (response.status === 400) {
    throw new AdminProfileError('Проверьте заполнение полей.')
  }
  if (!response.ok) {
    throw new AdminProfileError('Не удалось сохранить профиль. Попробуйте позже.')
  }
}

export async function createAdminSession(accessKey: string): Promise<void> {
  let response: Response

  try {
    response = await fetch('/api/v1/admin/session', {
      method: 'POST',
      credentials: 'same-origin',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ accessKey }),
    })
  } catch {
    throw new AdminAuthError('Не удалось связаться с сервером. Попробуйте ещё раз.')
  }

  if (response.status === 401) {
    throw new AdminAuthError('Ключ не подошёл. Проверьте его и попробуйте снова.')
  }

  if (!response.ok) {
    throw new AdminAuthError('Сервис временно недоступен. Попробуйте позже.')
  }
}

export async function verifyAdminSession(): Promise<boolean> {
  let response: Response

  try {
    response = await fetch('/api/v1/admin/session', {
      method: 'GET',
      credentials: 'same-origin',
    })
  } catch {
    throw new AdminAuthError('Не удалось проверить сессию. Попробуйте обновить страницу.')
  }

  if (response.status === 401) {
    return false
  }

  if (!response.ok) {
    throw new AdminAuthError('Сервис временно недоступен. Попробуйте позже.')
  }

  return true
}

export async function revokeAdminSession(): Promise<void> {
  let response: Response

  try {
    response = await fetch('/api/v1/admin/session', {
      method: 'DELETE',
      credentials: 'same-origin',
    })
  } catch {
    throw new AdminAuthError('Не удалось завершить сессию. Попробуйте ещё раз.')
  }

  if (!response.ok) {
    throw new AdminAuthError('Не удалось завершить сессию. Попробуйте позже.')
  }
}
