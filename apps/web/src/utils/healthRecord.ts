/**
 * Utilidades para registros de salud (health records).
 * Funciones puras, sin dependencias de Vue.
 * 
 * Nota: Las funciones relacionadas con status fueron eliminadas
 * ya que el campo status fue removido de health_records.
 */

import type { HealthCatalogCategoryType } from '@/constants/healthRecord'
import {
  HEALTH_CATALOG_CATEGORY_LABELS,
  HEALTH_CATALOG_CATEGORY_ICONS,
} from '@/constants/healthRecord'

/**
 * Obtiene el label en español para una categoría de registro de salud.
 * @param category - La categoría del registro
 * @returns Label en español
 */
export function getHealthCatalogCategoryLabel(category: HealthCatalogCategoryType | string): string {
  return HEALTH_CATALOG_CATEGORY_LABELS[category as HealthCatalogCategoryType] ?? category
}

/**
 * Obtiene el icono/emoji para una categoría de registro de salud.
 * @param category - La categoría del registro
 * @returns Icono/emoji para UI
 */
export function getHealthCatalogCategoryIcon(category: HealthCatalogCategoryType | string): string {
  return HEALTH_CATALOG_CATEGORY_ICONS[category as HealthCatalogCategoryType] ?? '📋'
}

/**
 * Verifica si un registro de salud está vencido.
 * Un registro está vencido si:
 * - Su next_dose_date es anterior a hoy
 * - No tiene application_date
 * @param nextDoseDate - Fecha de próxima dosis (ISO string)
 * @param applicationDate - Fecha de aplicación (ISO string)
 * @returns true si está vencido
 */
export function isOverdue(nextDoseDate: string | null, applicationDate?: string | null): boolean {
  if (!nextDoseDate || applicationDate) return false

  const due = new Date(nextDoseDate)
  const now = new Date()

  // Comparar solo fechas (sin horas)
  due.setUTCHours(0, 0, 0, 0)
  now.setUTCHours(0, 0, 0, 0)

  return due < now
}

/**
 * Calcula los días restantes hasta la fecha de vencimiento.
 * @param nextDoseDate - Fecha de próxima dosis (ISO string)
 * @returns Número de días restantes (negativo si ya venció)
 */
export function daysUntilDue(nextDoseDate: string | null): number {
  if (!nextDoseDate) return 0
  
  const due = new Date(nextDoseDate)
  const now = new Date()

  due.setUTCHours(0, 0, 0, 0)
  now.setUTCHours(0, 0, 0, 0)

  const diffMs = due.getTime() - now.getTime()
  return Math.ceil(diffMs / (1000 * 60 * 60 * 24))
}

/**
 * Verifica si un registro es urgente (vence en menos de 7 días).
 */
export function isUrgent(nextDoseDate: string | null, applicationDate?: string | null): boolean {
  if (!nextDoseDate || applicationDate) return false
  
  const days = daysUntilDue(nextDoseDate)
  return days >= 0 && days <= 7
}

/**
 * Formatea un mensaje de estado para la fecha de próxima dosis.
 * @param nextDoseDate - Fecha de próxima dosis (ISO string)
 * @param applicationDate - Fecha de aplicación (ISO string)
 * @returns Mensaje descriptivo
 */
export function getNextDoseDateStatusMessage(nextDoseDate: string | null, applicationDate?: string | null): string {
  if (applicationDate) {
    return 'Completado'
  }

  if (!nextDoseDate) {
    return 'Sin programar'
  }

  const days = daysUntilDue(nextDoseDate)

  if (days < 0) {
    return `Vencido hace ${Math.abs(days)} ${Math.abs(days) === 1 ? 'día' : 'días'}`
  }

  if (days === 0) {
    return 'Vence hoy'
  }

  if (days === 1) {
    return 'Vence mañana'
  }

  if (days <= 7) {
    return `Vence en ${days} ${days === 1 ? 'día' : 'días'}`
  }

  return `Vence el ${new Date(nextDoseDate).toLocaleDateString('es-ES', {
    day: 'numeric',
    month: 'short',
  })}`
}
