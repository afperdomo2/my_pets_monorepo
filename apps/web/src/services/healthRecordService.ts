import type {
  HealthRecord,
  CreateHealthRecordPayload,
  UpdateHealthRecordPayload,
} from '@/types/healthRecord'
import type { PaginatedResponse, ApiResponse } from '@/types/shared'
import { get, post, put, del } from '@/services/http'

const PER_PAGE_DEFAULT = 10

export const healthRecordService = {
  // Obtener todos los registros del usuario
  getAll(page = 1, perPage = PER_PAGE_DEFAULT): Promise<PaginatedResponse<HealthRecord>> {
    return get(`/health-records?page=${page}&per_page=${perPage}`)
  },

  // Obtener un registro por ID
  getById(id: string): Promise<ApiResponse<HealthRecord>> {
    return get(`/health-records/${id}`)
  },

  // Obtener registros de una mascota
  getByPetId(petId: string, page = 1, perPage = PER_PAGE_DEFAULT): Promise<PaginatedResponse<HealthRecord>> {
    return get(`/health-records/pets/${petId}?page=${page}&per_page=${perPage}`)
  },

  // Obtener registros de una mascota filtrados por categoría
  getByPetIdAndCategory(petId: string, category: string, page = 1, perPage = PER_PAGE_DEFAULT): Promise<PaginatedResponse<HealthRecord>> {
    return get(`/health-records/pets/${petId}/category/${category}?page=${page}&per_page=${perPage}`)
  },

  create(payload: CreateHealthRecordPayload): Promise<ApiResponse<HealthRecord>> {
    return post('/health-records', payload)
  },

  update(id: string, payload: UpdateHealthRecordPayload): Promise<ApiResponse<HealthRecord>> {
    return put(`/health-records/${id}`, payload)
  },

  remove(id: string): Promise<{ message: string }> {
    return del(`/health-records/${id}`)
  },

  // Obtener próximos registros con próxima dosis programada
  getUpcoming(page = 1, perPage = 10, category?: string): Promise<PaginatedResponse<HealthRecord>> {
    const params = new URLSearchParams({ page: String(page), per_page: String(perPage) })
    if (category) {
      params.append('category', category)
    }
    return get(`/health-records/upcoming?${params.toString()}`)
  }
}
