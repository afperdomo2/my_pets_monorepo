import type {
  HealthRecord,
  CreateHealthRecordPayload,
  UpdateHealthRecordPayload,
  UpdateStatusPayload,
} from '@/types/healthRecord'
import type { PaginatedResponse, ApiResponse } from '@/types/shared'
import { get, post, put, patch, del } from '@/services/http'

const PER_PAGE_DEFAULT = 10

export const healthRecordService = {
  // Obtener todos los registros del usuario
  getAll(page = 1, perPage = PER_PAGE_DEFAULT): Promise<PaginatedResponse<HealthRecord>> {
    return get(`/health-records?page=${page}&per_page=${perPage}`)
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

  updateStatus(id: string, payload: UpdateStatusPayload): Promise<ApiResponse<HealthRecord>> {
    return patch(`/health-records/${id}/status`, payload)
  },

  remove(id: string): Promise<{ message: string }> {
    return del(`/health-records/${id}`)
  }
}
