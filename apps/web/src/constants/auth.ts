/**
 * Constantes de autenticación y configuración de tokens.
 * Estos valores deben coincidir con la configuración del backend.
 */

/**
 * Proveedores de autenticación soportados.
 * Refleja el tipo AuthProvider del backend (apps/api/internal/models/user.go)
 * Uso: AuthProvider.Local, AuthProvider.Google
 */
export const AuthProvider = {
  Local: 'local',
  Google: 'google',
} as const

export type AuthProviderType = (typeof AuthProvider)[keyof typeof AuthProvider]

/**
 * Duración del access token en minutos.
 * Coincide con la configuración del backend.
 */
export const ACCESS_TOKEN_EXPIRY_MIN = 20

/**
 * Duración del refresh token en días.
 * Coincide con la configuración del backend.
 */
export const REFRESH_TOKEN_EXPIRY_DAYS = 20

/**
 * Nombres de las cookies JWT.
 */
export const ACCESS_TOKEN_COOKIE = 'access_token'
export const REFRESH_TOKEN_COOKIE = 'refresh_token'

/**
 * Rutas de autenticación.
 */
export const AUTH_ROUTES = {
  LOGIN: '/api/v1/auth/login',
  LOGOUT: '/api/v1/auth/logout',
  REFRESH: '/api/v1/auth/refresh',
  ME: '/api/v1/auth/me',
  GOOGLE: '/api/v1/auth/google',
  GOOGLE_CALLBACK: '/api/v1/auth/google/callback',
} as const
