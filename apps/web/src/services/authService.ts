import type { User, UpdateProfilePayload, ChangePasswordPayload } from '@/types/user'
import type { LoginCredentials, LoginResponse, MeResponse } from '@/types/auth'
import { get, post, put, del } from '@/services/http'

export const authService = {
  login(payload: LoginCredentials): Promise<LoginResponse> {
    return post('/auth/login', payload)
  },

  logout(): Promise<{ message: string }> {
    return post('/auth/logout')
  },

  refresh(): Promise<{ access_token: string; expires_in: number }> {
    return post('/auth/refresh')
  },

  me(): Promise<MeResponse> {
    return get('/auth/me')
  },

  /** Redirects the browser to the Google OAuth initiation endpoint on the API. */
  googleLogin(): void {
    // Use a relative URL so the Vite dev proxy forwards it to the API,
    // and in production the same domain serves the API.
    window.location.href = '/api/v1/auth/google'
  },

  updateProfile(payload: UpdateProfilePayload): Promise<{ data: User }> {
    return put('/auth/profile', payload)
  },

  changePassword(payload: ChangePasswordPayload): Promise<{ message: string }> {
    return put('/auth/password', payload)
  },

  googleLogout(): Promise<{ message: string }> {
    return del('/auth/google/logout')
  },
}
