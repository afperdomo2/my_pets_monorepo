/**
 * Pet utility helpers — pure functions, no Vue dependencies.
 * Imported by components, composables and views that display pet data.
 *
 * Nota: Las funciones de uso general (formatWeight, toGrams, formatAge, etc.)
 * están centralizadas en @/utils/formatters. Este archivo mantiene compatibilidad
 * y agrega funciones específicas de pets.
 */

import { calculateAge, isBirthday, estimateBirthDate, formatDate } from '@/utils/date'
import { formatWeight, toGrams, formatAge as formatAgeString } from '@/utils/formatters'
import { getLifeStageLabel } from '@/constants/lifeStage'

export interface AgeResult {
  years: number
  months: number
}

/**
 * Calculates years and months elapsed from a birth date ISO string to today.
 * Uses UTC dates to avoid timezone edge cases.
 * @deprecated Usar calculateAge desde @/utils/date
 */
export function calcAge(birthDateIso: string): AgeResult {
  return calculateAge(birthDateIso)
}

/**
 * Returns a human-readable age string.
 * Examples: "3 años, 2 meses" | "5 meses" | "1 año" | "Recién nacido"
 * @deprecated Usar formatAge desde @/utils/formatters
 */
export function formatAge(age: AgeResult): string {
  return formatAgeString(age.years, age.months)
}

/**
 * Checks whether today is the pet's birthday (same day and month).
 * Should ONLY be called when birth_date_exact === true.
 * @deprecated Usar isBirthday desde @/utils/date
 */
export function isBirthdayToday(birthDateIso: string): boolean {
  return isBirthday(birthDateIso)
}

/**
 * Formats a weight in grams to a human-readable string.
 * @deprecated Usar formatWeight desde @/utils/formatters
 */
export function formatWeightLegacy(grams: number): string {
  return formatWeight(grams)
}

/**
 * Converts a user input weight to grams based on the selected unit.
 * @deprecated Usar toGrams desde @/utils/formatters
 */
export function toGramsLegacy(value: number, unit: 'kg' | 'g'): number {
  return toGrams(value, unit)
}

/**
 * Calculates an estimated birth date string ("YYYY-MM-DD") from years and months.
 * @deprecated Usar estimateBirthDate desde @/utils/date
 */
export function estimatedBirthDate(years: number, months: number): string {
  return estimateBirthDate(years, months)
}

/**
 * Formats an ISO birth date for display (exact date only — not for estimated dates).
 * Example: "12 de enero de 2022"
 */
export function formatBirthDate(birthDateIso: string): string {
  return formatDate(birthDateIso, 'es-ES', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    timeZone: 'UTC',
  })
}

/**
 * Obtiene el label en español para una etapa de vida.
 * Centralizado en @/constants/lifeStage.
 */
export function lifeStageLabel(stage: string): string {
  return getLifeStageLabel(stage)
}

// Re-exportar funciones centralizadas para compatibilidad
export {
  calculateAge as calcAgeNew,
  formatWeight,
  toGrams,
  formatAge as formatAgeNew,
  getLifeStageLabel,
}
