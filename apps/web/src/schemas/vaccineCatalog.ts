import { z } from 'zod'

// Convierte undefined/null a "" para que Zod siempre evalúe con sus propios mensajes
const str = z.preprocess((v) => v ?? '', z.string())

export const createVaccineCatalogSchema = z.object({
  name: str
    .pipe(z.string().min(1, 'El nombre de la vacuna es obligatorio').min(2, 'El nombre debe tener al menos 2 caracteres')),
  species: z.array(str.pipe(z.string().min(1, 'La especie no puede estar vacía'))).min(1, 'Selecciona al menos una especie'),
  frequency_months: z.number().min(1, 'La frecuencia debe ser mayor que 0').max(360, 'La frecuencia no puede exceder 360 meses'),
  is_mandatory: z.boolean().optional().default(false),
})

export const updateVaccineCatalogSchema = z.object({
  name: str
    .pipe(z.string().min(1, 'El nombre de la vacuna es obligatorio').min(2, 'El nombre debe tener al menos 2 caracteres')),
  species: z.array(str.pipe(z.string().min(1, 'La especie no puede estar vacía'))).min(1, 'Selecciona al menos una especie'),
  frequency_months: z.number().min(1, 'La frecuencia debe ser mayor que 0').max(360, 'La frecuencia no puede exceder 360 meses'),
  is_mandatory: z.boolean().optional().default(false),
})

export type CreateVaccineCatalogFormValues = z.infer<typeof createVaccineCatalogSchema>
export type UpdateVaccineCatalogFormValues = z.infer<typeof updateVaccineCatalogSchema>
