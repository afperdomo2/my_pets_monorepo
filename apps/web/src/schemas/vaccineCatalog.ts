import { z } from 'zod'

const str = z.preprocess((v) => v ?? '', z.string())

export const createVaccineCatalogSchema = z.object({
  name: str
    .pipe(z.string().min(1, 'El nombre de la vacunas es obligatorio').min(2, 'El nombre debe tener al menos 2 caracteres')),
  frequency_months: z.number().min(1, 'La frecuencia debe ser mayor que 0').max(360, 'La frecuencia no puede exceder 360 meses').optional().nullable(),
  is_mandatory: z.boolean().optional().default(false),
})

export const updateVaccineCatalogSchema = z.object({
  name: str
    .pipe(z.string().min(1, 'El nombre de la vacunas es obligatorio').min(2, 'El nombre debe tener al menos 2 caracteres')),
  frequency_months: z.number().min(1, 'La frecuencia debe ser mayor que 0').max(360, 'La frecuencia no puede exceder 360 meses').optional().nullable(),
  is_mandatory: z.boolean().optional().default(false),
})

export type CreateVaccineCatalogFormValues = z.infer<typeof createVaccineCatalogSchema>
export type UpdateVaccineCatalogFormValues = z.infer<typeof updateVaccineCatalogSchema>
