import type { EditableArtistProfile } from '../types/artist'
import { APIClientError, apiRequest, apiRequestJSON } from './apiClient'

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
  let profile: Partial<EditableArtistProfile>
  try {
    profile = await apiRequestJSON<Partial<EditableArtistProfile>>('/api/v1/artist_profile')
  } catch (error) {
    if (error instanceof APIClientError) {
      if (error.status === 404) return { ...emptyArtistProfile }
      if (error.status === null) {
        throw new AdminProfileError('Не удалось загрузить профиль. Проверьте соединение.')
      }
    }
    throw new AdminProfileError('Не удалось загрузить профиль. Попробуйте позже.')
  }

  return {
    name: profile.name ?? '',
    tagline: profile.tagline ?? '',
    biography: profile.biography ?? '',
    artistStatement: profile.artistStatement ?? '',
    email: profile.email ?? '',
  }
}

export async function updateArtistProfile(profile: EditableArtistProfile): Promise<void> {
  try {
    await apiRequest('/api/v1/admin/artist_profile', {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        name: profile.name,
        tagline: profile.tagline,
        biography: profile.biography,
        artistStatement: profile.artistStatement || null,
        email: profile.email || null,
      }),
    })
  } catch (error) {
    if (error instanceof APIClientError) {
      if (error.status === null) {
        throw new AdminProfileError('Не удалось сохранить профиль. Проверьте соединение.')
      }
      if (error.status === 401) throw new AdminProfileError('Сессия завершилась. Войдите снова.', true)
      if (error.status === 400) throw new AdminProfileError('Проверьте заполнение полей.')
    }
    throw new AdminProfileError('Не удалось сохранить профиль. Попробуйте позже.')
  }
}

export async function createAdminSession(accessKey: string): Promise<void> {
  try {
    await apiRequest('/api/v1/admin/session', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ accessKey }),
    })
  } catch (error) {
    if (error instanceof APIClientError) {
      if (error.status === null) {
        throw new AdminAuthError('Не удалось связаться с сервером. Попробуйте ещё раз.')
      }
      if (error.status === 401) {
        throw new AdminAuthError('Ключ не подошёл. Проверьте его и попробуйте снова.')
      }
    }
    throw new AdminAuthError('Сервис временно недоступен. Попробуйте позже.')
  }
}

export async function verifyAdminSession(): Promise<boolean> {
  try {
    await apiRequest('/api/v1/admin/session')
    return true
  } catch (error) {
    if (error instanceof APIClientError && error.status === 401) return false
    if (error instanceof APIClientError && error.status === null) {
      throw new AdminAuthError('Не удалось проверить сессию. Попробуйте обновить страницу.')
    }
    throw new AdminAuthError('Сервис временно недоступен. Попробуйте позже.')
  }
}

export async function revokeAdminSession(): Promise<void> {
  try {
    await apiRequest('/api/v1/admin/session', {
      method: 'DELETE',
    })
  } catch (error) {
    if (error instanceof APIClientError && error.status === null) {
      throw new AdminAuthError('Не удалось завершить сессию. Попробуйте ещё раз.')
    }
    throw new AdminAuthError('Не удалось завершить сессию. Попробуйте позже.')
  }
}
