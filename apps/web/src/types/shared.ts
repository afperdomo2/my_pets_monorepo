/**
 * Tipos compartidos y reutilizables en toda la aplicación.
 * Estos tipos reflejan las respuestas estándar de la API.
 */

/**
 * Respuesta estándar de la API para un único recurso.
 * @template T - Tipo del recurso devuelto (ej: Pet, User, HealthRecord)
 */
export interface ApiResponse<T> {
  data: T
  total?: number
}

/**
 * Respuesta estándar de la API para respuestas paginadas.
 * Refleja exactamente el formato devuelto por el backend Go.
 * @template T - Tipo de cada elemento en el array de datos
 */
export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  per_page: number
  total_pages: number
}

/**
 * Parámetros de paginación para requests a la API.
 */
export interface PaginationParams {
  page?: number
  perPage?: number
}

/**
 * Estado de paginación para mantener en stores/composables.
 */
export interface PaginationState {
  currentPage: number
  perPage: number
  totalItems: number
  totalPages: number
}
