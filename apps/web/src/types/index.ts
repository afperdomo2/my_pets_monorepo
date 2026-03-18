/**
 * Índice de tipos y contratos de la aplicación.
 * Centraliza todas las exportaciones de tipos para facilitar imports.
 */

// Tipos compartidos
export * from './shared'

// Tipos por dominio
export * from './pet'
export * from './user'
export * from './healthRecord'
export * from './healthCatalog'
export * from './dashboard'
// Auth se exporta al final para evitar conflictos con user.ts
export type {
  LoginCredentials,
  LoginResponse,
  MeResponse,
  RefreshResponse,
  UpdateProfilePayload,
  ChangePasswordPayload,
} from './auth'
