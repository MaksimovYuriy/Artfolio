import {
  emptySocialHandles,
  socialPlatforms,
  type AdminSocialLink,
  type SocialHandles,
  type SocialPlatform,
} from '../types/socialLink'

export class SocialLinksServiceError extends Error {
  unauthorized: boolean

  constructor(message: string, unauthorized = false) {
    super(message)
    this.name = 'SocialLinksServiceError'
    this.unauthorized = unauthorized
  }
}

export async function getSocialHandles(): Promise<SocialHandles> {
  const response = await request('/api/v1/admin/social_links', 'Не удалось загрузить социальные сети.')
  const links = await response.json() as AdminSocialLink[]
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

async function request(url: string, fallbackMessage: string, init?: RequestInit): Promise<Response> {
  let response: Response
  try {
    response = await fetch(url, { credentials: 'same-origin', ...init })
  } catch {
    throw new SocialLinksServiceError('Не удалось связаться с сервером. Проверьте соединение.')
  }

  if (response.status === 401) {
    throw new SocialLinksServiceError('Сессия завершилась. Войдите снова.', true)
  }
  if (response.status === 400) {
    throw new SocialLinksServiceError('Проверьте адреса социальных сетей.')
  }
  if (!response.ok) throw new SocialLinksServiceError(fallbackMessage)
  return response
}
