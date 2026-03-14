/**
 * Cliente HTTP compartido para todos los servicios.
 * Centraliza la lógica de fetch, manejo de errores y autenticación.
 */

const BASE_URL = '/api/v1'

/**
 * Opciones por defecto para todas las requests.
 */
const DEFAULT_OPTIONS: RequestInit = {
  headers: {
    'Content-Type': 'application/json',
  },
  credentials: 'include', // Incluir cookies JWT
}

/**
 * Error personalizado para errores de la API.
 */
export class ApiError extends Error {
  status: number
  data?: unknown

  constructor(message: string, status: number, data?: unknown) {
    super(message)
    this.name = 'ApiError'
    this.status = status
    this.data = data
  }
}

/**
 * Realiza una request HTTP a la API.
 * @param url - URL relativa (sin BASE_URL)
 * @param options - Opciones de fetch
 * @returns Promise con la respuesta parseada como JSON
 * @throws ApiError si la request falla
 */
export async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${url}`, {
    ...DEFAULT_OPTIONS,
    ...options,
    headers: {
      ...DEFAULT_OPTIONS.headers,
      ...options?.headers,
    },
  })

  if (!res.ok) {
    let errorData: unknown
    let errorMessage = `HTTP ${res.status}`

    try {
      errorData = await res.json()
      if (errorData && typeof errorData === 'object' && 'error' in errorData) {
        errorMessage = String(errorData.error)
      }
    } catch {
      // Si no se puede parsear el JSON, usar mensaje genérico
    }

    throw new ApiError(errorMessage, res.status, errorData)
  }

  return res.json()
}

/**
 * Realiza una request GET.
 */
export function get<T>(url: string, options?: RequestInit): Promise<T> {
  return request<T>(url, { ...options, method: 'GET' })
}

/**
 * Realiza una request POST.
 */
export function post<T>(url: string, body?: unknown, options?: RequestInit): Promise<T> {
  return request<T>(url, {
    ...options,
    method: 'POST',
    body: body ? JSON.stringify(body) : undefined,
  })
}

/**
 * Realiza una request PUT.
 */
export function put<T>(url: string, body?: unknown, options?: RequestInit): Promise<T> {
  return request<T>(url, {
    ...options,
    method: 'PUT',
    body: body ? JSON.stringify(body) : undefined,
  })
}

/**
 * Realiza una request PATCH.
 */
export function patch<T>(url: string, body?: unknown, options?: RequestInit): Promise<T> {
  return request<T>(url, {
    ...options,
    method: 'PATCH',
    body: body ? JSON.stringify(body) : undefined,
  })
}

/**
 * Realiza una request DELETE.
 */
export function del<T>(url: string, options?: RequestInit): Promise<T> {
  return request<T>(url, { ...options, method: 'DELETE' })
}
