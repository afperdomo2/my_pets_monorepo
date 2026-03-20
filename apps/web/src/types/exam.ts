import type { ApiResponse, PaginatedResponse } from '@/types/shared'

/**
 * Estados posibles para un examen.
 */
export const ExamStatus = {
  Scheduled: 'scheduled',
  Completed: 'completed',
} as const

export type ExamStatusType = (typeof ExamStatus)[keyof typeof ExamStatus]

/**
 * Examen veterinario completo devuelto por la API.
 */
export interface Exam {
  id: string
  pet_id: string
  user_id: string
  name: string
  reason: string | null
  status: ExamStatusType
  scheduled_date: string | null
  completed_date: string | null
  notes: string | null
  created_at: string
  updated_at: string
}

/**
 * Resultado individual de un examen.
 */
export interface ExamResult {
  id: string
  exam_id: string
  parameter_name: string
  value: string
  unit: string | null
  created_at: string
}

/**
 * Examen con sus resultados.
 */
export interface ExamWithResults extends Exam {
  results?: ExamResult[]
}

/**
 * Payload para crear un examen.
 */
export interface CreateExamPayload {
  pet_id: string
  name: string
  reason?: string
  status?: ExamStatusType
  scheduled_date?: string
  completed_date?: string
  notes?: string
  results?: Array<{
    parameter_name: string
    value: string
    unit?: string
  }>
}

/**
 * Payload para actualizar un examen.
 */
export interface UpdateExamPayload {
  name?: string
  reason?: string
  scheduled_date?: string
  notes?: string
}

/**
 * Payload para programar un examen.
 */
export interface ScheduleExamPayload {
  scheduled_date: string
}

/**
 * Payload para completar un examen.
 */
export interface CompleteExamPayload {
  completed_date?: string
  results?: Array<{
    parameter_name: string
    value: string
    unit?: string
  }>
}

// Re-exportar tipos compartidos para compatibilidad
export type { ApiResponse, PaginatedResponse }
