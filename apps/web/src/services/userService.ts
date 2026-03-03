import type { User, CreateUserPayload, UpdateUserPayload } from '@/types/user'

const BASE_URL = '/api/v1'

interface ApiResponse<T> {
  data: T
  total?: number
}

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

export const userService = {
  getAll(): Promise<ApiResponse<User[]>> {
    return request('/users')
  },

  getById(id: string): Promise<ApiResponse<User>> {
    return request(`/users/${id}`)
  },

  create(payload: CreateUserPayload): Promise<ApiResponse<User>> {
    return request('/users', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },

  update(id: string, payload: UpdateUserPayload): Promise<ApiResponse<User>> {
    return request(`/users/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
  },

  remove(id: string): Promise<{ message: string }> {
    return request(`/users/${id}`, { method: 'DELETE' })
  },
}
