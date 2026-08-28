import type { AdminArtwork, Artwork, ArtworkInput } from '../types/artwork'
import { APIClientError, apiRequest, apiRequestJSON } from './apiClient'

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
  try {
    return await apiRequestJSON<T>(url, init)
  } catch (error) {
    throw artworkError(error, fallbackMessage)
  }
}

async function request(url: string, fallbackMessage: string, init?: RequestInit): Promise<void> {
  try {
    await apiRequest(url, init)
  } catch (error) {
    throw artworkError(error, fallbackMessage)
  }
}

function artworkError(error: unknown, fallbackMessage: string): ArtworkServiceError {
  if (!(error instanceof APIClientError)) return new ArtworkServiceError(fallbackMessage)
  if (error.status === null) return new ArtworkServiceError('Не удалось связаться с сервером. Проверьте соединение.')
  if (error.status === 401) return new ArtworkServiceError('Сессия завершилась. Войдите снова.', true)
  if (error.status === 413) return new ArtworkServiceError('Файл слишком большой. Максимальный размер — 12 МБ.')
  if (error.status === 415) return new ArtworkServiceError('Выберите корректное изображение JPEG или PNG.')
  if (error.status === 422) return new ArtworkServiceError('У изображения слишком большое разрешение.')
  if (error.status === 400) return new ArtworkServiceError('Проверьте заполнение полей.')
  if (error.status === 404) return new ArtworkServiceError('Работа не найдена. Возможно, она уже удалена.')
  if (error.status === 409) return new ArtworkServiceError('Список работ изменился. Обновите страницу и повторите сортировку.')
  return new ArtworkServiceError(fallbackMessage)
}
