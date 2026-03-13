export interface VaccineCatalog {
  id: string
  name: string
  species: string[]
  frequency_months: number
  is_mandatory: boolean
  created_at: string
  updated_at: string
}

export interface CreateVaccineCatalogPayload {
  name: string
  species: string[]
  frequency_months: number
  is_mandatory?: boolean
}

export interface UpdateVaccineCatalogPayload {
  name: string
  species: string[]
  frequency_months: number
  is_mandatory?: boolean
}
