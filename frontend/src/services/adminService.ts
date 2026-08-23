export class AdminAuthError extends Error {
  constructor(message: string) {
    super(message)
    this.name = 'AdminAuthError'
  }
}

export async function createAdminSession(accessKey: string): Promise<void> {
  let response: Response

  try {
    response = await fetch('/api/admin/session', {
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
    response = await fetch('/api/admin/session', {
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
    response = await fetch('/api/admin/session', {
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
