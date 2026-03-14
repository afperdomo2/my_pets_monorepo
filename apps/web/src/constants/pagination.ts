/**
 * Constantes de paginación para toda la aplicación.
 * Estos valores deben coincidir con los defaults del backend.
 */

/**
 * Cantidad de registros por página por defecto.
 * Coincide con el default del backend (ver handlers en internal/domain/)
 */
export const PER_PAGE_DEFAULT = 10

/**
 * Opciones disponibles para el usuario seleccionar.
 */
export const PER_PAGE_OPTIONS = [5, 10, 20, 50] as const

/**
 * Máxima cantidad de registros por página permitida.
 */
export const PER_PAGE_MAX = 100

/**
 * Mínima cantidad de registros por página permitida.
 */
export const PER_PAGE_MIN = 1

/**
 * Tipo derivado de las opciones válidas.
 */
export type PerPageValue = (typeof PER_PAGE_OPTIONS)[number]

/**
 * Valida que un valor de per_page sea válido.
 */
export function isValidPerPage(value: number): value is PerPageValue {
  return PER_PAGE_OPTIONS.includes(value as PerPageValue)
}
