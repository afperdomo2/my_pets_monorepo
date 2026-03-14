import type { 
  HealthRecord, 
  CreateHealthRecordPayload, 
  UpdateHealthRecordPayload, 
  UpdateStatusPayload,
  ApiResponse,
  PaginatedResponse
} from '@/types/healthRecord'

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

export const healthRecordService = {
  // Obtener registros de una mascota
  getByPetId(petId: string, page = 1, perPage = 10): Promise<PaginatedResponse<HealthRecord>> {
    return request(`/health-records/pets/${petId}?page=${page}&per_page=${perPage}`)
  },

  // Obtener registros de una mascota filtrados por categoría
  getByPetIdAndCategory(petId: string, category: string, page = 1, perPage = 10): Promise<PaginatedResponse<HealthRecord>> {
    return request(`/health-records/pets/${petId}/category/${category}?page=${page}&per_page=${perPage}`)
  },

  create(payload: CreateHealthRecordPayload): Promise<ApiResponse<HealthRecord>> {
    return request('/health-records', {
      method: 'POST',
      body: JSON.stringify(payload)
    })
  },

  update(id: string, payload: UpdateHealthRecordPayload): Promise<ApiResponse<HealthRecord>> {
    return request(`/health-records/${id}`, {
      method: 'PUT',
      body: JSON.stringify(payload)
    })
  },

  updateStatus(id: string, payload: UpdateStatusPayload): Promise<ApiResponse<HealthRecord>> {
    return request(`/health-records/${id}/status`, {
      method: 'PATCH',
      body: JSON.stringify(payload)
    })
  },

  remove(id: string): Promise<{ message: string }> {
    return request(`/health-records/${id}`, { method: 'DELETE' })
  }
}
