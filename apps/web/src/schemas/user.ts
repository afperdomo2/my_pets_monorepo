import { z } from 'zod'

// Convierte undefined/null a "" para que Zod siempre evalúe con sus propios mensajes
const str = z.preprocess((v) => v ?? '', z.string())

export const loginSchema = z.object({
  email: str
    .pipe(z.string().min(1, 'El correo electrónico es obligatorio').email('El correo electrónico no es válido')),
  password: str
    .pipe(z.string().min(1, 'La contraseña es obligatoria')),
})

export const setupSchema = z
  .object({
    name: str
      .pipe(z.string().min(1, 'El nombre es obligatorio').min(2, 'El nombre debe tener al menos 2 caracteres')),
    email: str
      .pipe(z.string().min(1, 'El correo electrónico es obligatorio').email('El correo electrónico no es válido')),
    password: str
      .pipe(z.string().min(1, 'La contraseña es obligatoria').min(8, 'La contraseña debe tener al menos 8 caracteres')),
    confirmPassword: str
      .pipe(z.string().min(1, 'Debes confirmar la contraseña')),
  })
  .refine((data) => data.password === data.confirmPassword, {
    message: 'Las contraseñas no coinciden',
    path: ['confirmPassword'],
  })

export const createUserSchema = z.object({
  name: str
    .pipe(z.string().min(1, 'El nombre es obligatorio')),
  email: str
    .pipe(z.string().min(1, 'El correo electrónico es obligatorio').email('El correo electrónico no es válido')),
  password: str
    .pipe(z.string().min(1, 'La contraseña es obligatoria').min(8, 'La contraseña debe tener al menos 8 caracteres')),
})

export const updateUserSchema = z.object({
  name: str
    .pipe(z.string().min(1, 'El nombre es obligatorio')),
  email: str
    .pipe(z.string().min(1, 'El correo electrónico es obligatorio').email('El correo electrónico no es válido')),
})

export type LoginFormValues = z.infer<typeof loginSchema>
export type SetupFormValues = z.infer<typeof setupSchema>
export type CreateUserFormValues = z.infer<typeof createUserSchema>
export type UpdateUserFormValues = z.infer<typeof updateUserSchema>
