import type { VaccineCatalog, CreateVaccineCatalogPayload, UpdateVaccineCatalogPayload } from '@/types/vaccineCatalog'
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

export const vaccineCatalogService = {
  getAll(page = 1, perPage = 10, species?: string): Promise<PaginatedResponse<VaccineCatalog>> {
    let url = `/vaccines-catalog?page=${page}&per_page=${perPage}`
    if (species) {
      url += `&species=${species}`
    }
    return request(url)
  },

  getById(id: string): Promise<ApiResponse<VaccineCatalog>> {
    return request(`/vaccines-catalog/${id}`)
  },

  getBySpecies(species: string, page = 1, perPage = 10): Promise<PaginatedResponse<VaccineCatalog>> {
    return request(`/vaccines-catalog/species/${species}?page=${page}&per_page=${perPage}`)
  },

  create(payload: CreateVaccineCatalogPayload): Promise<ApiResponse<VaccineCatalog>> {
    return request('/vaccines-catalog', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },

  update(id: string, payload: UpdateVaccineCatalogPayload): Promise<ApiResponse<VaccineCatalog>> {
    return request(`/vaccines-catalog/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
  },

  remove(id: string): Promise<{ message: string }> {
    return request(`/vaccines-catalog/${id}`, { method: 'DELETE' })
  },
}
