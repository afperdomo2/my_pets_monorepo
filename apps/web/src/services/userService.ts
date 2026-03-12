import type { User, UserWithPetCount, CreateUserPayload, UpdateUserPayload } from '@/types/user'
import type { ApiResponse, PaginatedResponse } from '@/types/pet'

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

export const userService = {
  getAll(page = 1, perPage = 10): Promise<PaginatedResponse<UserWithPetCount>> {
    return request(`/users?page=${page}&per_page=${perPage}`)
  },

  getById(id: string): Promise<ApiResponse<UserWithPetCount>> {
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
