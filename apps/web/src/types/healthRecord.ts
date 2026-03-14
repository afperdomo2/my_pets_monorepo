import type { HealthRecordStatusType } from '@/constants/healthRecord'
import type { ApiResponse, PaginatedResponse } from '@/types/shared'

/**
 * Registro de salud completo devuelto por la API.
 */
export interface HealthRecord {
  id: string
  pet_id: string
  user_id: string
  health_catalog_id: string | null
  category: string
  name: string
  status: HealthRecordStatusType
  application_date: string | null
  due_date: string
  notes: string | null
  created_at: string
  updated_at: string
  pet: {
    id: string
    name: string
    species: string
    breed: string
  }
}

/**
 * Payload para crear un registro de salud.
 */
export interface CreateHealthRecordPayload {
  pet_id: string
  health_catalog_id?: string
  category?: string
  name?: string
  status?: string
  application_date?: string
  due_date: string
  notes?: string
}

/**
 * Payload para actualizar un registro de salud completo.
 * Pet_id y health_catalog_id no son actualizables.
 */
export interface UpdateHealthRecordPayload {
  category: string
  name: string
  status?: string
  application_date?: string
  due_date: string
  notes?: string
}

/**
 * Payload para actualizar solo el status de un registro.
 */
export interface UpdateStatusPayload {
  status: string
  application_date?: string
}

// Re-exportar tipos compartidos para compatibilidad
export type { ApiResponse, PaginatedResponse }
