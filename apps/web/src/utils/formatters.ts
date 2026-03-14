/**
 * Funciones de formateo genéricas y utilidades varias.
 * Funciones puras, sin dependencias de Vue.
 */

/**
 * Capitaliza la primera letra de un string.
 * @param str - String a capitalizar
 * @returns String con primera letra mayúscula
 */
export function capitalize(str: string): string {
  if (!str) return str
  return str.charAt(0).toUpperCase() + str.slice(1)
}

/**
 * Capitaliza cada palabra de un string.
 * @param str - String a capitalizar
 * @returns String con cada palabra capitalizada
 */
export function capitalizeWords(str: string): string {
  if (!str) return str
  return str
    .split(' ')
    .filter(Boolean)
    .map((word) => capitalize(word))
    .join(' ')
}

/**
 * Trunca un string a una longitud máxima, agregando ellipsis si es necesario.
 * @param str - String a truncar
 * @param maxLength - Longitud máxima
 * @returns String truncado
 */
export function truncate(str: string, maxLength: number): string {
  if (!str || str.length <= maxLength) return str
  return `${str.slice(0, maxLength)}...`
}

/**
 * Formatea un número con separadores de miles.
 * @param num - Número a formatear
 * @param locale - Locale para formateo (default: 'es-ES')
 * @returns String formateado
 */
export function formatNumber(num: number, locale = 'es-ES'): string {
  return new Intl.NumberFormat(locale).format(num)
}

/**
 * Formatea un peso en gramos a string legible.
 * - < 1000 g → "X g" (ej: "850 g")
 * - >= 1000 g → "X kg" con hasta 2 decimales (ej: "3.5 kg")
 * @param grams - Peso en gramos
 * @returns String formateado
 */
export function formatWeight(grams: number): string {
  if (grams < 1000) return `${grams} g`
  const kg = grams / 1000
  // Eliminar ceros trailing: 3.50 → "3.5", 4.00 → "4"
  return `${parseFloat(kg.toFixed(2))} kg`
}

/**
 * Convierte un peso ingresado por el usuario a gramos.
 * @param value - Valor numérico
 * @param unit - Unidad de medida ('kg' o 'g')
 * @returns Valor en gramos
 */
export function toGrams(value: number, unit: 'kg' | 'g'): number {
  return unit === 'kg' ? Math.round(value * 1000) : Math.round(value)
}

/**
 * Formatea una edad (años y meses) como string legible.
 * @param years - Años de edad
 * @param months - Meses de edad
 * @returns String formateado (ej: "3 años, 2 meses" | "5 meses" | "1 año")
 */
export function formatAge(years: number, months: number): string {
  if (years === 0 && months === 0) return 'Recién nacido'
  if (years === 0) return `${months} ${months === 1 ? 'mes' : 'meses'}`
  if (months === 0) return `${years} ${years === 1 ? 'año' : 'años'}`
  return `${years} ${years === 1 ? 'año' : 'años'}, ${months} ${months === 1 ? 'mes' : 'meses'}`
}

/**
 * Obtiene las iniciales de un nombre completo.
 * @param name - Nombre completo
 * @returns Iniciales (ej: "Juan Pérez" → "JP")
 */
export function getInitials(name: string): string {
  if (!name) return ''
  return name
    .split(' ')
    .filter(Boolean)
    .map((part) => part[0])
    .join('')
    .toUpperCase()
    .slice(0, 2)
}

/**
 * Valida un email con regex básica.
 * @param email - Email a validar
 * @returns true si el email es válido
 */
export function isValidEmail(email: string): boolean {
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  return emailRegex.test(email)
}

/**
 * Genera un ID único simple (para uso en frontend).
 * @returns String ID único
 */
export function generateId(): string {
  return `${Date.now()}-${Math.random().toString(36).slice(2, 9)}`
}

/**
 * Escapa caracteres HTML para prevenir XSS.
 * @param str - String a escapar
 * @returns String escapado
 */
export function escapeHtml(str: string): string {
  const htmlEntities: Record<string, string> = {
    '&': '&amp;',
    '<': '&lt;',
    '>': '&gt;',
    '"': '&quot;',
    "'": '&#39;',
  }
  return str.replace(/[&<>"']/g, (char) => htmlEntities[char] ?? char)
}
