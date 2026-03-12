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
  pet_limit: z.number().min(0).optional(),
})

export const updateUserSchema = z.object({
  name: str
    .pipe(z.string().min(1, 'El nombre es obligatorio')),
  email: str
    .pipe(z.string().min(1, 'El correo electrónico es obligatorio').email('El correo electrónico no es válido')),
  pet_limit: z.number().min(0).optional(),
})

export const updateProfileSchema = z.object({
  name: str
    .pipe(z.string().min(1, 'El nombre es obligatorio').min(2, 'El nombre debe tener al menos 2 caracteres').max(100, 'El nombre no puede superar los 100 caracteres')),
  email: str
    .pipe(z.string().min(1, 'El correo electrónico es obligatorio').email('El correo electrónico no es válido')),
})

export const changePasswordSchema = z
  .object({
    current_password: str
      .pipe(z.string().min(1, 'La contraseña actual es obligatoria')),
    new_password: str
      .pipe(z.string().min(1, 'La nueva contraseña es obligatoria').min(8, 'La contraseña debe tener al menos 8 caracteres').max(72, 'La contraseña no puede superar los 72 caracteres')),
    confirm_password: str
      .pipe(z.string().min(1, 'Debes confirmar la nueva contraseña')),
  })
  .refine((data) => data.new_password === data.confirm_password, {
    message: 'Las contraseñas no coinciden',
    path: ['confirm_password'],
  })

export type LoginFormValues = z.infer<typeof loginSchema>
export type SetupFormValues = z.infer<typeof setupSchema>
export type CreateUserFormValues = z.infer<typeof createUserSchema>
export type UpdateUserFormValues = z.infer<typeof updateUserSchema>
export type UpdateProfileFormValues = z.infer<typeof updateProfileSchema>
export type ChangePasswordFormValues = z.infer<typeof changePasswordSchema>
