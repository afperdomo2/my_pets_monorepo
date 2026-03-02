export interface User {
  id: number
  name: string
  email: string
  auth_provider: 'local' | 'google'
  is_system_user: boolean
  created_at: string
  updated_at: string
}

export interface CreateUserPayload {
  name: string
  email: string
  password: string
}

export interface UpdateUserPayload {
  name: string
  email: string
}

export interface LoginPayload {
  email: string
  password: string
}

export interface AuthResponse {
  data: User
  expires_in: number
}
