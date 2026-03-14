/**
 * Tipos relacionados con autenticación.
 * Estos tipos reflejan las requests/responses de la API de autenticación.
 */

import type { User } from '@/types/user'
import type { ApiResponse } from '@/types/shared'

/**
 * Credenciales para login.
 */
export interface LoginCredentials {
  email: string
  password: string
}

/**
 * Respuesta de la API tras login exitoso.
 * Incluye el usuario y el tiempo de expiración del token.
 */
export interface LoginResponse extends ApiResponse<User> {
  expires_in: number
}

/**
 * Respuesta de la API para el endpoint /auth/me.
 */
// eslint-disable-next-line @typescript-eslint/no-empty-object-type
export interface MeResponse extends ApiResponse<User> {}

/**
 * Respuesta de la API para refresh de token.
 */
export interface RefreshResponse {
  access_token: string
  expires_in: number
}

/**
 * Payload para actualizar perfil de usuario.
 */
export interface UpdateProfilePayload {
  name: string
  email: string
}

/**
 * Payload para cambiar contraseña.
 * Solo válido para usuarios con auth_provider = 'local'.
 */
export interface ChangePasswordPayload {
  current_password: string
  new_password: string
  confirm_password: string
}
