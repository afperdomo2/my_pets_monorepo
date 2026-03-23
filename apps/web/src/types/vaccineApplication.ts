import type { ApiResponse } from '@/types/shared'

/**
 * Aplicación de vacuna/desparasitación.
 * Representa una dosis aplicada de un tratamiento.
 */
export interface VaccineApplication {
  id: string
  health_record_id: string
  application_date: string
  notes: string | null
  created_at: string
}

/**
 * Payload para crear una aplicación de vacuna.
 */
export interface CreateVaccineApplicationPayload {
  health_record_id: string
  application_date: string
  notes?: string
  next_dose_date?: string | null
}

/**
 * Payload para actualizar una aplicación de vacuna.
 */
export interface UpdateVaccineApplicationPayload {
  application_date?: string
  notes?: string
}

// Re-exportar ApiResponse para compatibilidad
export type { ApiResponse }
