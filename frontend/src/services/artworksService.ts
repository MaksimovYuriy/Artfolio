import type { AdminArtwork, Artwork, ArtworkInput } from '../types/artwork'

export async function getArtworks(): Promise<Artwork[]> {
  return requestJSON<Artwork[]>('/api/v1/artworks', 'Не удалось загрузить работы.')
}

export async function getAdminArtworks(): Promise<AdminArtwork[]> {
  return requestJSON<AdminArtwork[]>('/api/v1/admin/artworks', 'Не удалось загрузить список работ.')
}

export async function createArtwork(input: ArtworkInput, image: File, isPublished: boolean): Promise<AdminArtwork> {
  return requestJSON<AdminArtwork>('/api/v1/admin/artworks', 'Не удалось сохранить работу.', {
    method: 'POST',
    body: artworkFormData(input, isPublished, image),
  })
}

export async function updateArtwork(id: number, input: ArtworkInput, image: File | null, isPublished: boolean): Promise<void> {
  await request(`/api/v1/admin/artworks/${id}`, 'Не удалось обновить работу.', {
    method: 'PUT',
    body: artworkFormData(input, isPublished, image),
  })
}

export async function deleteArtwork(id: number): Promise<void> {
  await request(`/api/v1/admin/artworks/${id}`, 'Не удалось удалить работу.', { method: 'DELETE' })
}

export async function reorderArtworks(artworkIds: number[]): Promise<void> {
  await request('/api/v1/admin/artworks/order', 'Не удалось сохранить порядок работ.', {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ artworkIds }),
  })
}

export class ArtworkServiceError extends Error {
  unauthorized: boolean

  constructor(message: string, unauthorized = false) {
    super(message)
    this.name = 'ArtworkServiceError'
    this.unauthorized = unauthorized
  }
}

function artworkFormData(input: ArtworkInput, isPublished: boolean, image: File | null): FormData {
  const form = new FormData()
  form.set('title', input.title)
  form.set('description', input.description)
  form.set('technique', input.technique)
  form.set('year', input.year)
  form.set('imageAlt', input.imageAlt)
  form.set('isPublished', String(isPublished))
  if (image) form.set('image', image)
  return form
}

async function requestJSON<T>(url: string, fallbackMessage: string, init?: RequestInit): Promise<T> {
  const response = await request(url, fallbackMessage, init)
  return response.json() as Promise<T>
}

async function request(url: string, fallbackMessage: string, init?: RequestInit): Promise<Response> {
  let response: Response
  try {
    response = await fetch(url, { credentials: 'same-origin', ...init })
  } catch {
    throw new ArtworkServiceError('Не удалось связаться с сервером. Проверьте соединение.')
  }

  if (response.status === 401) throw new ArtworkServiceError('Сессия завершилась. Войдите снова.', true)
  if (response.status === 413) throw new ArtworkServiceError('Файл слишком большой. Максимальный размер — 12 МБ.')
  if (response.status === 415) throw new ArtworkServiceError('Выберите корректное изображение JPEG или PNG.')
  if (response.status === 422) throw new ArtworkServiceError('У изображения слишком большое разрешение.')
  if (response.status === 400) throw new ArtworkServiceError('Проверьте заполнение полей.')
  if (response.status === 404) throw new ArtworkServiceError('Работа не найдена. Возможно, она уже удалена.')
  if (response.status === 409) throw new ArtworkServiceError('Список работ изменился. Обновите страницу и повторите сортировку.')
  if (!response.ok) throw new ArtworkServiceError(fallbackMessage)
  return response
}
