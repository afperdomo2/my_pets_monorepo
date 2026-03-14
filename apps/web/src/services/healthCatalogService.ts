import type { HealthCatalog, CreateHealthCatalogPayload, UpdateHealthCatalogPayload } from '@/types/healthCatalog'
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

export const healthCatalogService = {
  getAll(category: string, page = 1, perPage = 10, species?: string): Promise<PaginatedResponse<HealthCatalog>> {
    let url = `/health-catalogs/category/${category}?page=${page}&per_page=${perPage}`
    if (species) {
      url += `&species=${species}`
    }
    return request(url)
  },

  getById(id: string): Promise<ApiResponse<HealthCatalog>> {
    return request(`/health-catalogs/${id}`)
  },

  create(payload: CreateHealthCatalogPayload): Promise<ApiResponse<HealthCatalog>> {
    return request('/health-catalogs', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },

  update(id: string, payload: UpdateHealthCatalogPayload): Promise<ApiResponse<HealthCatalog>> {
    return request(`/health-catalogs/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload),
    })
  },

  remove(id: string): Promise<{ message: string }> {
    return request(`/health-catalogs/${id}`, { method: 'DELETE' })
  },
}
