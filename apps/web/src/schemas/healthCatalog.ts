import { z } from 'zod'

const str = z.preprocess((v) => v ?? '', z.string())

// Categorías válidas para la guía de salud
const categoryEnum = z.enum(['vaccine', 'deworming', 'exam'], {
  required_error: 'La categoría es obligatoria',
  invalid_type_error: 'Categoría inválida',
})

export const createHealthCatalogSchema = z.object({
  name: str
    .pipe(z.string().min(1, 'El nombre es obligatorio').min(2, 'El nombre debe tener al menos 2 caracteres')),
  category: categoryEnum,
  description: str.pipe(z.string().max(1000, 'La descripción no puede superar 1000 caracteres')).optional(),
  frequency_months: z.number().min(1, 'La frecuencia debe ser mayor que 0').max(360, 'La frecuencia no puede exceder 360 meses').optional().nullable(),
  is_mandatory: z.boolean().optional().default(false),
})

export const updateHealthCatalogSchema = z.object({
  name: str
    .pipe(z.string().min(1, 'El nombre es obligatorio').min(2, 'El nombre debe tener al menos 2 caracteres')),
  category: categoryEnum,
  description: str.pipe(z.string().max(1000, 'La descripción no puede superar 1000 caracteres')).optional(),
  frequency_months: z.number().min(1, 'La frecuencia debe ser mayor que 0').max(360, 'La frecuencia no puede exceder 360 meses').optional().nullable(),
  is_mandatory: z.boolean().optional().default(false),
})

export type CreateHealthCatalogFormValues = z.infer<typeof createHealthCatalogSchema>
export type UpdateHealthCatalogFormValues = z.infer<typeof updateHealthCatalogSchema>
