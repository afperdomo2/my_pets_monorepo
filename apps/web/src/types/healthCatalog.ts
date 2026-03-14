import type { HealthCatalogCategoryType } from '@/constants/healthRecord'

// Re-exportar HealthCatalogCategoryType para compatibilidad
export type { HealthCatalogCategoryType }

/**
 * Registro del catálogo de salud devuelto por la API.
 */
export interface HealthCatalog {
  id: string
  name: string
  category: HealthCatalogCategoryType
  description: string
  species: string[]
  frequency_months: number | null
  is_mandatory: boolean
  created_at: string
  updated_at: string
}

/**
 * Payload para crear un registro del catálogo de salud.
 * Solo usuarios sistema pueden crear/actualizar/eliminar.
 */
export interface CreateHealthCatalogPayload {
  name: string
  category: HealthCatalogCategoryType
  description?: string
  species: string[]
  frequency_months: number | null
  is_mandatory?: boolean
}

/**
 * Payload para actualizar un registro del catálogo de salud.
 */
export interface UpdateHealthCatalogPayload {
  name: string
  category: HealthCatalogCategoryType
  description?: string
  species: string[]
  frequency_months: number | null
  is_mandatory?: boolean
}
