/**
 * Pet utility helpers — pure functions, no Vue dependencies.
 * Imported by components, composables and views that display pet data.
 */

export interface AgeResult {
  years: number
  months: number
}

/**
 * Calculates years and months elapsed from a birth date ISO string to today.
 * Uses UTC dates to avoid timezone edge cases.
 */
export function calcAge(birthDateIso: string): AgeResult {
  const birth = new Date(birthDateIso)
  const now = new Date()

  let years = now.getUTCFullYear() - birth.getUTCFullYear()
  let months = now.getUTCMonth() - birth.getUTCMonth()

  if (months < 0) {
    years--
    months += 12
  }

  // If the day of the month hasn't been reached yet this month, subtract one month
  if (now.getUTCDate() < birth.getUTCDate()) {
    months--
    if (months < 0) {
      years--
      months += 12
    }
  }

  return { years: Math.max(0, years), months: Math.max(0, months) }
}

/**
 * Returns a human-readable age string.
 * Examples: "3 años, 2 meses" | "5 meses" | "1 año" | "Recién nacido"
 */
export function formatAge(age: AgeResult): string {
  const { years, months } = age

  if (years === 0 && months === 0) return 'Recién nacido'
  if (years === 0) return `${months} ${months === 1 ? 'mes' : 'meses'}`
  if (months === 0) return `${years} ${years === 1 ? 'año' : 'años'}`
  return `${years} ${years === 1 ? 'año' : 'años'}, ${months} ${months === 1 ? 'mes' : 'meses'}`
}

/**
 * Checks whether today is the pet's birthday (same day and month).
 * Should ONLY be called when birth_date_exact === true.
 */
export function isBirthdayToday(birthDateIso: string): boolean {
  const birth = new Date(birthDateIso)
  const now = new Date()
  return birth.getUTCDate() === now.getUTCDate() && birth.getUTCMonth() === now.getUTCMonth()
}

/**
 * Formats a weight in grams to a human-readable string.
 * - < 1000 g → "Xg"   (e.g. "850 g")
 * - >= 1000 g → "X kg" with up to 2 decimal places (e.g. "3.5 kg")
 */
export function formatWeight(grams: number): string {
  if (grams < 1000) return `${grams} g`
  const kg = grams / 1000
  // Remove trailing zeros: 3.50 → "3.5", 4.00 → "4"
  return `${parseFloat(kg.toFixed(2))} kg`
}

/**
 * Converts a user input weight to grams based on the selected unit.
 */
export function toGrams(value: number, unit: 'kg' | 'g'): number {
  return unit === 'kg' ? Math.round(value * 1000) : Math.round(value)
}

/**
 * Calculates an estimated birth date string ("YYYY-MM-DD") from years and months.
 * Subtracts the given period from today's date.
 */
export function estimatedBirthDate(years: number, months: number): string {
  const now = new Date()
  let y = now.getUTCFullYear() - years
  let m = now.getUTCMonth() - months

  if (m < 0) {
    y--
    m += 12
  }

  // Use the same day-of-month as today, clamped to valid range for the target month
  const maxDay = new Date(Date.UTC(y, m + 1, 0)).getUTCDate()
  const d = Math.min(now.getUTCDate(), maxDay)

  return `${y}-${String(m + 1).padStart(2, '0')}-${String(d).padStart(2, '0')}`
}

/**
 * Formats an ISO birth date for display (exact date only — not for estimated dates).
 * Example: "12 de enero de 2022"
 */
export function formatBirthDate(birthDateIso: string): string {
  return new Date(birthDateIso).toLocaleDateString('es-ES', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    timeZone: 'UTC',
  })
}

const LIFE_STAGE_LABEL: Record<string, string> = {
  // Perro
  puppy: 'Cachorro',
  junior: 'Joven',
  adult: 'Adulto',
  senior: 'Senior',
  geriatric: 'Geriátrico',
  // Gato
  kitten: 'Gatito',
  young_adult: 'Joven Adulto',
  mature_adult: 'Adulto Maduro',
  end_of_life: 'Fin de Vida',
  // Conejo
  infant: 'Infancia',
  juvenile: 'Juvenil',
  teenager: 'Adolescente',
}

export function lifeStageLabel(stage: string): string {
  return LIFE_STAGE_LABEL[stage] ?? stage
}
