export interface User {
  id: string
  name: string
  email: string
  auth_provider: 'local' | 'google'
  is_system_user: boolean
  pet_limit: number
  created_at: string
  updated_at: string
}

export interface UserWithPetCount extends User {
  pet_count: number
}

export interface CreateUserPayload {
  name: string
  email: string
  password: string
  pet_limit?: number
}

export interface UpdateUserPayload {
  name: string
  email: string
  pet_limit?: number
}

/**
 * @deprecated Usar desde @/types/auth
 */
export interface UpdateProfilePayload {
  name: string
  email: string
}

/**
 * @deprecated Usar desde @/types/auth
 */
export interface ChangePasswordPayload {
  current_password: string
  new_password: string
  confirm_password: string
}

/**
 * @deprecated Usar LoginCredentials desde @/types/auth
 */
export interface LoginPayload {
  email: string
  password: string
}

/**
 * @deprecated Usar LoginResponse desde @/types/auth
 */
export interface AuthResponse {
  data: User
  expires_in: number
}
