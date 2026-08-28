interface APIErrorEnvelope {
  error?: {
    code?: string
    message?: string
    requestId?: string
  }
}

export class APIClientError extends Error {
  status: number | null
  code: string
  requestId?: string

  constructor(message: string, status: number | null, code: string, requestId?: string) {
    super(message)
    this.name = 'APIClientError'
    this.status = status
    this.code = code
    this.requestId = requestId
  }
}

export async function apiRequest(url: string, init?: RequestInit): Promise<Response> {
  let response: Response
  try {
    response = await fetch(url, { credentials: 'same-origin', ...init })
  } catch {
    throw new APIClientError('Network request failed', null, 'network_error')
  }

  if (!response.ok) {
    const payload = await readErrorEnvelope(response)
    throw new APIClientError(
      payload.error?.message ?? `HTTP request failed with status ${response.status}`,
      response.status,
      payload.error?.code ?? 'http_error',
      payload.error?.requestId ?? response.headers.get('X-Request-ID') ?? undefined,
    )
  }

  return response
}

export async function apiRequestJSON<T>(url: string, init?: RequestInit): Promise<T> {
  const response = await apiRequest(url, init)
  try {
    return await response.json() as T
  } catch {
    throw new APIClientError(
      'Server returned invalid JSON',
      response.status,
      'invalid_response',
      response.headers.get('X-Request-ID') ?? undefined,
    )
  }
}

async function readErrorEnvelope(response: Response): Promise<APIErrorEnvelope> {
  try {
    return await response.json() as APIErrorEnvelope
  } catch {
    return {}
  }
}
