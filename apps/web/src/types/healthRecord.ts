import type { VaccineApplication } from '@/types/vaccineApplication'
import type { ApiResponse, PaginatedResponse } from '@/types/shared'

/**
 * Registro de salud completo devuelto por la API.
 * Solo maneja vacunas (vaccine) y desparasitaciones (deworming).
 */
export interface HealthRecord {
  id: string
  pet_id: string
  user_id: string
  health_catalog_id: string | null
  category: 'vaccine' | 'deworming'
  name: string
  application_date: string | null
  next_dose_date: string | null
  notes: string | null
  created_at: string
  updated_at: string
  pet: {
    id: string
    name: string
    species: string
    breed: string
  }
  // Relación con vaccine_applications (opcional, se puede precargar)
  vaccine_applications?: VaccineApplication[]
}

/**
 * Payload para crear un registro de salud.
 */
export interface CreateHealthRecordPayload {
  pet_id: string
  health_catalog_id?: string
  category?: 'vaccine' | 'deworming'
  name?: string
  application_date?: string
  next_dose_date?: string
  notes?: string
}

/**
 * Payload para actualizar un registro de salud completo.
 * Pet_id y health_catalog_id no son actualizables.
 */
export interface UpdateHealthRecordPayload {
  category: 'vaccine' | 'deworming'
  name: string
  application_date?: string
  next_dose_date?: string
  notes?: string
}

// Re-exportar tipos compartidos para compatibilidad
export type { ApiResponse, PaginatedResponse }
