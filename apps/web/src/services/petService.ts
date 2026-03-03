import type { Pet, PetPayload, ApiResponse } from '@/types/pet'

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
  getAll(): Promise<ApiResponse<Pet[]>> {
    return request('/pets')
  },

  getById(id: string): Promise<ApiResponse<Pet>> {
    return request(`/pets/${id}`)
  },

  create(payload: PetPayload): Promise<ApiResponse<Pet>> {
    return request('/pets', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },

  update(id: string, payload: PetPayload): Promise<ApiResponse<Pet>> {
    return request(`/pets/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
  },

  remove(id: string): Promise<{ message: string }> {
    return request(`/pets/${id}`, { method: 'DELETE' })
  },
}
