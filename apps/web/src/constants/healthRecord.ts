/**
 * Constantes para registros de salud (health records).
 * Estos valores deben coincidir con los definidos en el backend:
 * apps/api/internal/domain/health_record/payload.go
 * 
 * Nota: El campo status fue eliminado de health_records.
 * Ahora solo se manejan vacunas (vaccine) y desparasitaciones (deworming).
 */

// ==================== CATEGORÍAS ====================

/**
 * Categorías válidas para registros de salud.
 * Debe coincidir con el binding del backend: oneof=vaccine deworming
 * Uso: HealthCatalogCategory.Vaccine, HealthCatalogCategory.Deworming
 */
export const HealthCatalogCategory = {
  Vaccine: 'vaccine',
  Deworming: 'deworming',
} as const

export type HealthCatalogCategoryType = (typeof HealthCatalogCategory)[keyof typeof HealthCatalogCategory]

/**
 * Labels en español para las categorías.
 */
export const HEALTH_CATALOG_CATEGORY_LABELS: Record<HealthCatalogCategoryType, string> = {
  [HealthCatalogCategory.Vaccine]: 'Vacuna',
  [HealthCatalogCategory.Deworming]: 'Desparasitación',
}

/**
 * Descripciones para cada categoría.
 */
export const HEALTH_CATALOG_CATEGORY_DESCRIPTIONS: Record<HealthCatalogCategoryType, string> = {
  [HealthCatalogCategory.Vaccine]: 'Vacunas preventivas y terapéuticas',
  [HealthCatalogCategory.Deworming]: 'Desparasitación interna y externa',
}

/**
 * Iconos/emojis para cada categoría (para UI).
 */
export const HEALTH_CATALOG_CATEGORY_ICONS: Record<HealthCatalogCategoryType, string> = {
  [HealthCatalogCategory.Vaccine]: '💉',
  [HealthCatalogCategory.Deworming]: '💊',
}

// ==================== HELPERS ====================

/**
 * Obtiene el label en español para una categoría de registro de salud.
 */
export function getHealthCatalogCategoryLabel(category: string): string {
  return HEALTH_CATALOG_CATEGORY_LABELS[category as HealthCatalogCategoryType] ?? category
}

/**
 * Valida si un string es una categoría válida de registro de salud.
 */
export function isValidHealthCatalogCategory(category: string): category is HealthCatalogCategoryType {
  return Object.values(HealthCatalogCategory).includes(category as HealthCatalogCategoryType)
}
