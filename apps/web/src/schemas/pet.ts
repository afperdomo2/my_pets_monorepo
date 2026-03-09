import { z } from 'zod'

// Convierte undefined/null a "" para que Zod siempre evalúe con sus propios mensajes
const str = z.preprocess((v) => v ?? '', z.string())

const SPECIES = ['dog', 'cat', 'bird', 'rabbit', 'fish', 'other'] as const

// Schema compartido de campos comunes
const basePetSchema = z.object({
  name: str.pipe(z.string().min(1, 'El nombre es obligatorio').max(100)),
  species: str.pipe(z.string().min(1, 'La especie es obligatoria')).refine(
    (v) => (SPECIES as readonly string[]).includes(v),
    { message: 'Especie no válida' },
  ),
  breed: z.string().max(100).optional(),
  birth_date: str.pipe(z.string().min(1, 'La fecha de nacimiento es obligatoria')),
  birth_date_exact: z.boolean(),
})

// Schema para crear mascota — incluye peso y etapa de vida
export const createPetSchema = basePetSchema.extend({
  weight_grams: z
    .number({ invalid_type_error: 'El peso debe ser un número' })
    .int('El peso debe ser un número entero')
    .min(1, 'El peso debe ser mayor a 0')
    .optional(),
  life_stage: z.string().optional(),
})

// Schema para editar mascota — sin peso ni etapa de vida
export const updatePetSchema = basePetSchema

export type CreatePetFormValues = z.infer<typeof createPetSchema>
export type UpdatePetFormValues = z.infer<typeof updatePetSchema>
