import { z } from 'zod'

// Convierte undefined/null a "" para que Zod siempre evalúe con sus propios mensajes
const str = z.preprocess((v) => v ?? '', z.string())

export const petSchema = z.object({
  name: str.pipe(z.string().min(1, 'El nombre es obligatorio')),
  species: str.pipe(z.string().min(1, 'La especie es obligatoria')),
  breed: z.string().optional(),
  age: z
    .number({ invalid_type_error: 'La edad debe ser un número' })
    .int('La edad debe ser un número entero')
    .min(0, 'La edad no puede ser negativa')
    .max(100, 'La edad no puede superar 100 años')
    .optional(),
  owner: z.string().optional(),
})

export type PetFormValues = z.infer<typeof petSchema>
