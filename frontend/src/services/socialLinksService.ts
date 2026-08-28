import {
  emptySocialHandles,
  socialPlatforms,
  type AdminSocialLink,
  type SocialHandles,
  type SocialPlatform,
} from '../types/socialLink'
import { APIClientError, apiRequest, apiRequestJSON } from './apiClient'

export class SocialLinksServiceError extends Error {
  unauthorized: boolean

  constructor(message: string, unauthorized = false) {
    super(message)
    this.name = 'SocialLinksServiceError'
    this.unauthorized = unauthorized
  }
}

export async function getSocialHandles(): Promise<SocialHandles> {
  const links = await requestJSON<AdminSocialLink[]>('/api/v1/admin/social_links', 'Не удалось загрузить социальные сети.')
  const handles: SocialHandles = { ...emptySocialHandles }

  for (const link of links) {
    if (isSocialPlatform(link.platform)) handles[link.platform] = link.handle
  }
  return handles
}

export async function replaceSocialHandles(handles: SocialHandles): Promise<void> {
  await request('/api/v1/admin/social_links', 'Не удалось сохранить социальные сети.', {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      socialLinks: socialPlatforms.map(({ id }) => ({ platform: id, handle: handles[id] })),
    }),
  })
}

function isSocialPlatform(value: string): value is SocialPlatform {
  return socialPlatforms.some(({ id }) => id === value)
}

async function request(url: string, fallbackMessage: string, init?: RequestInit): Promise<void> {
  try {
    await apiRequest(url, init)
  } catch (error) {
    throw socialLinksError(error, fallbackMessage)
  }
}

async function requestJSON<T>(url: string, fallbackMessage: string, init?: RequestInit): Promise<T> {
  try {
    return await apiRequestJSON<T>(url, init)
  } catch (error) {
    throw socialLinksError(error, fallbackMessage)
  }
}

function socialLinksError(error: unknown, fallbackMessage: string): SocialLinksServiceError {
  if (!(error instanceof APIClientError)) return new SocialLinksServiceError(fallbackMessage)
  if (error.status === null) return new SocialLinksServiceError('Не удалось связаться с сервером. Проверьте соединение.')
  if (error.status === 401) return new SocialLinksServiceError('Сессия завершилась. Войдите снова.', true)
  if (error.status === 400) return new SocialLinksServiceError('Проверьте адреса социальных сетей.')
  return new SocialLinksServiceError(fallbackMessage)
}
