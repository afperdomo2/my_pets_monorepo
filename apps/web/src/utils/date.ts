/**
 * Utilidades para manejo y formateo de fechas.
 * Funciones puras, sin dependencias de Vue.
 */

/**
 * Opciones de formateo predefinidas.
 */
export const DATE_FORMATS = {
  /** Ej: "14 de marzo de 2026" */
  LONG: { day: 'numeric', month: 'long', year: 'numeric' } as Intl.DateTimeFormatOptions,
  /** Ej: "14/03/2026" */
  SHORT: { day: '2-digit', month: '2-digit', year: 'numeric' } as Intl.DateTimeFormatOptions,
  /** Ej: "14/03/2026 15:30" */
  DATETIME: { day: '2-digit', month: '2-digit', year: 'numeric', hour: '2-digit', minute: '2-digit' } as Intl.DateTimeFormatOptions,
  /** Ej: "2026-03-14" (ISO date only) */
  ISO_DATE: { year: 'numeric', month: '2-digit', day: '2-digit' } as Intl.DateTimeFormatOptions,
} as const

/**
 * Parsea un string ISO a Date, usando UTC para evitar problemas de timezone.
 * @param dateString - String en formato ISO (YYYY-MM-DD o YYYY-MM-DDTHH:MM:SS)
 * @returns Objeto Date
 */
export function parseISO(dateString: string): Date {
  return new Date(dateString)
}

/**
 * Formatea una fecha usando el locale especificado.
 * @param date - Fecha a formatear (Date o string ISO)
 * @param locale - Locale para formateo (default: 'es-ES')
 * @param options - Opciones de formateo
 * @returns String formateado
 */
export function formatDate(
  date: Date | string,
  locale = 'es-ES',
  options: Intl.DateTimeFormatOptions = DATE_FORMATS.LONG
): string {
  const dateObj = typeof date === 'string' ? parseISO(date) : date
  return dateObj.toLocaleDateString(locale, options)
}

/**
 * Formatea una fecha relativa (hoy, ayer, hace X días, etc.).
 * @param date - Fecha a formatear (Date o string ISO)
 * @param locale - Locale para formateo (default: 'es-ES')
 * @returns String con fecha relativa
 */
export function formatRelativeDate(date: Date | string, locale = 'es-ES'): string {
  const dateObj = typeof date === 'string' ? parseISO(date) : date
  const now = new Date()
  const diffMs = now.getTime() - dateObj.getTime()
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24))

  const rtf = new Intl.RelativeTimeFormat(locale, { numeric: 'auto' })

  if (diffDays === 0) {
    return rtf.format(0, 'day') // "hoy"
  }

  if (diffDays === 1) {
    return rtf.format(-1, 'day') // "ayer"
  }

  if (diffDays < 7) {
    return rtf.format(-diffDays, 'day') // "hace X días"
  }

  if (diffDays < 30) {
    const weeks = Math.floor(diffDays / 7)
    return rtf.format(-weeks, 'week')
  }

  if (diffDays < 365) {
    const months = Math.floor(diffDays / 30)
    return rtf.format(-months, 'month')
  }

  const years = Math.floor(diffDays / 365)
  return rtf.format(-years, 'year')
}

/**
 * Verifica si una fecha es hoy.
 * @param date - Fecha a verificar (Date o string ISO)
 * @returns true si la fecha es hoy
 */
export function isToday(date: Date | string): boolean {
  const dateObj = typeof date === 'string' ? parseISO(date) : date
  const now = new Date()
  return (
    dateObj.getUTCDate() === now.getUTCDate() &&
    dateObj.getUTCMonth() === now.getUTCMonth() &&
    dateObj.getUTCFullYear() === now.getUTCFullYear()
  )
}

/**
 * Verifica si una fecha es cumpleaños (mismo día y mes que hoy).
 * Solo debe usarse cuando birth_date_exact === true.
 * @param birthDateIso - Fecha de nacimiento en formato ISO
 * @returns true si hoy es el cumpleaños
 */
export function isBirthday(birthDateIso: string): boolean {
  const birth = new Date(birthDateIso)
  const now = new Date()
  return birth.getUTCDate() === now.getUTCDate() && birth.getUTCMonth() === now.getUTCMonth()
}

/**
 * Formatea una fecha como string ISO (YYYY-MM-DD).
 * @param date - Fecha a formatear
 * @returns String en formato YYYY-MM-DD
 */
export function toISODate(date: Date): string {
  const year = date.getUTCFullYear()
  const month = String(date.getUTCMonth() + 1).padStart(2, '0')
  const day = String(date.getUTCDate()).padStart(2, '0')
  return `${year}-${month}-${day}`
}

/**
 * Calcula la edad en años y meses desde una fecha de nacimiento.
 * Usa UTC para evitar problemas de timezone.
 * @param birthDateIso - Fecha de nacimiento en formato ISO
 * @returns Objeto con years y months
 */
export function calculateAge(birthDateIso: string): { years: number; months: number } {
  const birth = new Date(birthDateIso)
  const now = new Date()

  let years = now.getUTCFullYear() - birth.getUTCFullYear()
  let months = now.getUTCMonth() - birth.getUTCMonth()

  if (months < 0) {
    years--
    months += 12
  }

  // Si el día del mes aún no ha llegado, restar un mes
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
 * Calcula una fecha estimada de nacimiento a partir de años y meses.
 * @param years - Años de edad
 * @param months - Meses de edad
 * @returns String en formato YYYY-MM-DD
 */
export function estimateBirthDate(years: number, months: number): string {
  const now = new Date()
  let y = now.getUTCFullYear() - years
  let m = now.getUTCMonth() - months

  if (m < 0) {
    y--
    m += 12
  }

  // Usar el mismo día del mes, ajustado al rango válido
  const maxDay = new Date(Date.UTC(y, m + 1, 0)).getUTCDate()
  const d = Math.min(now.getUTCDate(), maxDay)

  return `${y}-${String(m + 1).padStart(2, '0')}-${String(d).padStart(2, '0')}`
}

/**
 * Verifica si una fecha es válida.
 * @param date - Fecha a verificar
 * @returns true si la fecha es válida
 */
export function isValidDate(date: Date | string): boolean {
  const dateObj = typeof date === 'string' ? parseISO(date) : date
  return !isNaN(dateObj.getTime())
}

/**
 * Obtiene el nombre del mes en español.
 * @param month - Número del mes (0-11)
 * @returns Nombre del mes en español
 */
export function getMonthName(month: number, locale = 'es-ES'): string {
  const date = new Date(2000, month, 1)
  return date.toLocaleDateString(locale, { month: 'long' })
}

/**
 * Formatea una fecha (solo fecha, sin hora) usando UTC para evitar problemas de timezone.
 * Ideal para fechas que vienen del backend en formato YYYY-MM-DD o YYYY-MM-DDT00:00:00Z.
 * 
 * @param dateStr - String de fecha en formato ISO (YYYY-MM-DD o YYYY-MM-DDTHH:MM:SSZ)
 * @param emptyValue - Valor a mostrar si la fecha es nula o vacía (default: '—')
 * @returns String formateado (ej: "23 mar 2026") o el valor vacío si no hay fecha
 */
export function formatDateOnly(
  dateStr: string | null | undefined,
  emptyValue = '—'
): string {
  if (!dateStr) return emptyValue
  
  // Extraer solo la parte de fecha (YYYY-MM-DD) del string ISO
  const datePart = dateStr.split('T')[0] ?? ''
  if (!datePart) return emptyValue

  const parts = datePart.split('-').map(Number)
  const year = parts[0]
  const month = parts[1]
  const day = parts[2]

  if (!year || month == null || !day) return emptyValue

  // Crear fecha usando UTC para evitar problemas de timezone
  // Los meses en JS son 0-indexados (0 = enero, 11 = diciembre)
  const date = new Date(Date.UTC(year, month - 1, day))
  
  // Formatear usando los componentes UTC directamente
  const dayNum = date.getUTCDate()
  const monthIndex = date.getUTCMonth()
  const yearNum = date.getUTCFullYear()
  
  // Nombres de meses abreviados en español
  const monthNames = ['ene', 'feb', 'mar', 'abr', 'may', 'jun', 'jul', 'ago', 'sep', 'oct', 'nov', 'dic']
  
  return `${dayNum} ${monthNames[monthIndex]} ${yearNum}`
}
