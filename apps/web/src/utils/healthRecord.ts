/**
 * Utilidades para registros de salud (health records).
 * Funciones puras, sin dependencias de Vue.
 */

import type { HealthRecordStatusType, HealthCatalogCategoryType } from '@/constants/healthRecord'
import {
  HEALTH_RECORD_STATUS_LABELS,
  HEALTH_CATALOG_CATEGORY_LABELS,
  HEALTH_RECORD_STATUS_COLORS,
  HEALTH_CATALOG_CATEGORY_ICONS,
} from '@/constants/healthRecord'

/**
 * Obtiene el label en español para un estado de registro de salud.
 * @param status - El estado del registro
 * @returns Label en español
 */
export function getHealthRecordStatusLabel(status: HealthRecordStatusType | string): string {
  return HEALTH_RECORD_STATUS_LABELS[status as HealthRecordStatusType] ?? status
}

/**
 * Obtiene el color UI para un estado de registro de salud.
 * @param status - El estado del registro
 * @returns Color para UI (tailwind/css)
 */
export function getHealthRecordStatusColor(status: HealthRecordStatusType | string): string {
  return HEALTH_RECORD_STATUS_COLORS[status as HealthRecordStatusType] ?? 'gray'
}

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
 * - Su estado no es 'applied'
 * - Su due_date es anterior a hoy
 * @param dueDate - Fecha de vencimiento (ISO string)
 * @param status - Estado actual del registro
 * @returns true si está vencido
 */
export function isOverdue(dueDate: string, status: string): boolean {
  if (status === 'applied') return false

  const due = new Date(dueDate)
  const now = new Date()

  // Comparar solo fechas (sin horas)
  due.setUTCHours(0, 0, 0, 0)
  now.setUTCHours(0, 0, 0, 0)

  return due < now
}

/**
 * Calcula los días restantes hasta la fecha de vencimiento.
 * @param dueDate - Fecha de vencimiento (ISO string)
 * @returns Número de días restantes (negativo si ya venció)
 */
export function daysUntilDue(dueDate: string): number {
  const due = new Date(dueDate)
  const now = new Date()

  due.setUTCHours(0, 0, 0, 0)
  now.setUTCHours(0, 0, 0, 0)

  const diffMs = due.getTime() - now.getTime()
  return Math.ceil(diffMs / (1000 * 60 * 60 * 24))
}

/**
 * Formatea un mensaje de estado para la fecha de vencimiento.
 * @param dueDate - Fecha de vencimiento (ISO string)
 * @param status - Estado actual del registro
 * @returns Mensaje descriptivo
 */
export function getDueDateStatusMessage(dueDate: string, status: HealthRecordStatusType | string): string {
  if (status === 'applied') {
    return 'Completado'
  }

  const days = daysUntilDue(dueDate)

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
    return `Vence en ${days} días`
  }

  return `Vence el ${new Date(dueDate).toLocaleDateString('es-ES', {
    day: 'numeric',
    month: 'short',
  })}`
}
