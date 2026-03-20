/**
 * Constantes para exámenes veterinarios.
 * Estos valores deben coincidir con los definidos en el backend:
 * apps/api/internal/domain/exam/payload.go
 */

// ==================== ESTADOS ====================

/**
 * Estados posibles para un examen.
 */
export const ExamStatus = {
  Scheduled: 'scheduled',
  Completed: 'completed',
} as const

export type ExamStatusType = (typeof ExamStatus)[keyof typeof ExamStatus]

/**
 * Labels en español para los estados de exámenes.
 */
export const EXAM_STATUS_LABELS: Record<ExamStatusType, string> = {
  [ExamStatus.Scheduled]: 'Programado',
  [ExamStatus.Completed]: 'Completado',
}

/**
 * Colores sugeridos para cada estado (para UI).
 */
export const EXAM_STATUS_COLORS: Record<ExamStatusType, string> = {
  [ExamStatus.Scheduled]: 'blue',
  [ExamStatus.Completed]: 'green',
}

// ==================== HELPERS ====================

/**
 * Obtiene el label en español para un estado de examen.
 */
export function getExamStatusLabel(status: string): string {
  return EXAM_STATUS_LABELS[status as ExamStatusType] ?? status
}

/**
 * Valida si un string es un estado válido de examen.
 */
export function isValidExamStatus(status: string): status is ExamStatusType {
  return Object.values(ExamStatus).includes(status as ExamStatusType)
}
