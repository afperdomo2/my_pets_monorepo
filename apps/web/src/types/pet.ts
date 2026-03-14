import type { PetSize } from '@/constants/petSize'
import type { LifeStage } from '@/constants/lifeStage'
import type { ApiResponse, PaginatedResponse } from '@/types/shared'

/**
 * Mascota completa devuelta por la API.
 */
export interface Pet {
  id: string
  user_id: string
  name: string
  species: string
  breed: string
  birth_date: string        // ISO "YYYY-MM-DDT..." from Go time.Time
  birth_date_exact: boolean
  weight_grams: number | null
  life_stage: LifeStage | null // null for non-dog/cat species
  size: PetSize | null      // null for non-dog species
  created_at: string
  updated_at: string
}

/**
 * Payload para crear una mascota.
 */
export interface CreatePetPayload {
  name: string
  species: string
  breed?: string
  birth_date: string        // "YYYY-MM-DD"
  birth_date_exact: boolean
  weight_grams?: number     // grams, only at creation
  life_stage?: LifeStage    // only at creation, only dog/cat
  size?: PetSize            // required when species === 'dog'
}

/**
 * Payload para actualizar una mascota.
 * Species se excluye intencionalmente — no se puede cambiar después de creación.
 * Weight_grams y life_stage se excluyen — se calculan automáticamente.
 */
export interface UpdatePetPayload {
  name: string
  breed?: string
  birth_date: string
  birth_date_exact: boolean
  size?: PetSize            // required when species === 'dog'; null/absent for others
}

// Re-exportar tipos compartidos para compatibilidad
export type { ApiResponse, PaginatedResponse }
