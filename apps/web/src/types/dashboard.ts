/**
 * Tipos para el dominio de dashboard.
 */

/**
 * Resumen del dashboard con los datos principales para el usuario.
 * Refleja exactamente el formato devuelto por el backend Go.
 */
export interface DashboardSummary {
  total_pets: number
  healthy_pets: number
  pending_tasks: number
  overdue_tasks: number
}

/**
 * Respuesta de la API para el resumen del dashboard.
 */
export interface DashboardSummaryResponse {
  data: DashboardSummary
}
