/**
 * Constantes para registros de salud (health records).
 * Estos valores deben coincidir con los definidos en el backend:
 * apps/api/internal/domain/health_record/payload.go
 */

// ==================== STATUS ====================

/**
 * Estados posibles para un registro de salud.
 * Nota: 'overdue' es calculado en runtime por el backend, no se persiste en BD.
 * Uso: HealthRecordStatus.Pending, HealthRecordStatus.Applied, HealthRecordStatus.Overdue
 */
export const HealthRecordStatus = {
  Pending: 'pending',
  Applied: 'applied',
  Overdue: 'overdue',
} as const

export type HealthRecordStatusType = (typeof HealthRecordStatus)[keyof typeof HealthRecordStatus]

/**
 * Labels en español para los estados de registros de salud.
 */
export const HEALTH_RECORD_STATUS_LABELS: Record<HealthRecordStatusType, string> = {
  [HealthRecordStatus.Pending]: 'Pendiente',
  [HealthRecordStatus.Applied]: 'Aplicado',
  [HealthRecordStatus.Overdue]: 'Vencido',
}

/**
 * Descripciones para cada estado.
 */
export const HEALTH_RECORD_STATUS_DESCRIPTIONS: Record<HealthRecordStatusType, string> = {
  [HealthRecordStatus.Pending]: 'El registro está pendiente de aplicación',
  [HealthRecordStatus.Applied]: 'El registro ha sido aplicado correctamente',
  [HealthRecordStatus.Overdue]: 'El registro ha vencido y no ha sido aplicado',
}

/**
 * Colores sugeridos para cada estado (para UI).
 */
export const HEALTH_RECORD_STATUS_COLORS: Record<HealthRecordStatusType, string> = {
  [HealthRecordStatus.Pending]: 'yellow',
  [HealthRecordStatus.Applied]: 'green',
  [HealthRecordStatus.Overdue]: 'red',
}

// ==================== CATEGORÍAS ====================

/**
 * Categorías válidas para registros de salud.
 * Debe coincidir con el binding del backend: oneof=vaccine deworming exam
 * Uso: HealthCatalogCategory.Vaccine, HealthCatalogCategory.Deworming, HealthCatalogCategory.Exam
 */
export const HealthCatalogCategory = {
  Vaccine: 'vaccine',
  Deworming: 'deworming',
  Exam: 'exam',
} as const

export type HealthCatalogCategoryType = (typeof HealthCatalogCategory)[keyof typeof HealthCatalogCategory]

/**
 * Labels en español para las categorías.
 */
export const HEALTH_CATALOG_CATEGORY_LABELS: Record<HealthCatalogCategoryType, string> = {
  [HealthCatalogCategory.Vaccine]: 'Vacuna',
  [HealthCatalogCategory.Deworming]: 'Desparasitación',
  [HealthCatalogCategory.Exam]: 'Examen',
}

/**
 * Descripciones para cada categoría.
 */
export const HEALTH_CATALOG_CATEGORY_DESCRIPTIONS: Record<HealthCatalogCategoryType, string> = {
  [HealthCatalogCategory.Vaccine]: 'Vacunas preventivas y terapéuticas',
  [HealthCatalogCategory.Deworming]: 'Desparasitación interna y externa',
  [HealthCatalogCategory.Exam]: 'Exámenes y chequeos veterinarios',
}

/**
 * Iconos/emojis para cada categoría (para UI).
 */
export const HEALTH_CATALOG_CATEGORY_ICONS: Record<HealthCatalogCategoryType, string> = {
  [HealthCatalogCategory.Vaccine]: '💉',
  [HealthCatalogCategory.Deworming]: '💊',
  [HealthCatalogCategory.Exam]: '🩺',
}

// ==================== HELPERS ====================

/**
 * Obtiene el label en español para un estado de registro de salud.
 */
export function getHealthRecordStatusLabel(status: string): string {
  return HEALTH_RECORD_STATUS_LABELS[status as HealthRecordStatusType] ?? status
}

/**
 * Obtiene el label en español para una categoría de registro de salud.
 */
export function getHealthCatalogCategoryLabel(category: string): string {
  return HEALTH_CATALOG_CATEGORY_LABELS[category as HealthCatalogCategoryType] ?? category
}

/**
 * Valida si un string es un estado válido de registro de salud.
 */
export function isValidHealthRecordStatus(status: string): status is HealthRecordStatusType {
  return Object.values(HealthRecordStatus).includes(status as HealthRecordStatusType)
}

/**
 * Valida si un string es una categoría válida de registro de salud.
 */
export function isValidHealthCatalogCategory(category: string): category is HealthCatalogCategoryType {
  return Object.values(HealthCatalogCategory).includes(category as HealthCatalogCategoryType)
}
