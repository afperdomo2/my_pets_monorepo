/**
 * Índice de utilidades y helpers de la aplicación.
 * Centraliza todas las exportaciones de funciones utilitarias.
 */

export * from './date'
export * from './formatters'
export * from './healthRecord'
// Pet re-exporta algunas funciones de formatters, evitamos duplicación
export {
  calcAge,
  formatAge,
  isBirthdayToday,
  formatWeight,
  toGrams,
  estimatedBirthDate,
  formatBirthDate,
  lifeStageLabel,
} from './pet'
