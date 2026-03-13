// Categorías válidas para los registros de la guía de salud
export type HealthCatalogCategory = 'vaccine' | 'deworming' | 'exam'

export interface HealthCatalog {
  id: string
  name: string
  category: HealthCatalogCategory
  description: string
  species: string[]
  frequency_months: number | null
  is_mandatory: boolean
  created_at: string
  updated_at: string
}

export interface CreateHealthCatalogPayload {
  name: string
  category: HealthCatalogCategory
  description?: string
  species: string[]
  frequency_months: number | null
  is_mandatory?: boolean
}

export interface UpdateHealthCatalogPayload {
  name: string
  category: HealthCatalogCategory
  description?: string
  species: string[]
  frequency_months: number | null
  is_mandatory?: boolean
}
