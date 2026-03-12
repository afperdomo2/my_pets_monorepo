import type { User, LoginPayload, AuthResponse, UpdateProfilePayload, ChangePasswordPayload } from '@/types/user'

const BASE_URL = '/api/v1'

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${url}`, {
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    ...options,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'Unknown error' }))
    throw new Error(err.error ?? `HTTP ${res.status}`)
  }
  return res.json()
}

export const authService = {
  login(payload: LoginPayload): Promise<AuthResponse> {
    return request('/auth/login', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },

  logout(): Promise<{ message: string }> {
    return request('/auth/logout', { method: 'POST' })
  },

  refresh(): Promise<AuthResponse> {
    return request('/auth/refresh', { method: 'POST' })
  },

  me(): Promise<{ data: User }> {
    return request('/auth/me')
  },

  /** Redirects the browser to the Google OAuth initiation endpoint on the API. */
  googleLogin(): void {
    // Use a relative URL so the Vite dev proxy forwards it to the API,
    // and in production the same domain serves the API.
    window.location.href = '/api/v1/auth/google'
  },

  updateProfile(payload: UpdateProfilePayload): Promise<{ data: User }> {
    return request('/auth/profile', {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
  },

  changePassword(payload: ChangePasswordPayload): Promise<{ message: string }> {
    return request('/auth/password', {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
  },
}
