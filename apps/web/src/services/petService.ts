import type { Pet, CreatePetPayload, UpdatePetPayload, ApiResponse, PaginatedResponse } from '@/types/pet'

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

export const petService = {
  getAll(page = 1, perPage = 10): Promise<PaginatedResponse<Pet>> {
    return request(`/pets?page=${page}&per_page=${perPage}`)
  },

  getById(id: string): Promise<ApiResponse<Pet>> {
    return request(`/pets/${id}`)
  },

  create(payload: CreatePetPayload): Promise<ApiResponse<Pet>> {
    return request('/pets', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },

  update(id: string, payload: UpdatePetPayload): Promise<ApiResponse<Pet>> {
    return request(`/pets/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
  },

  remove(id: string): Promise<{ message: string }> {
    return request(`/pets/${id}`, { method: 'DELETE' })
  },
}
