import type {
  VaccineApplication,
  CreateVaccineApplicationPayload,
  UpdateVaccineApplicationPayload,
} from '@/types/vaccineApplication'
import type { ApiResponse } from '@/types/shared'
import { get, post, put, del } from '@/services/http'

export const vaccineApplicationService = {
  // Obtener aplicaciones de un health_record
  getByHealthRecordId(healthRecordId: string): Promise<ApiResponse<VaccineApplication[]>> {
    return get(`/vaccine-applications/health-record/${healthRecordId}`)
  },

  // Obtener una aplicación por ID
  getById(id: string): Promise<ApiResponse<VaccineApplication>> {
    return get(`/vaccine-applications/${id}`)
  },

  create(payload: CreateVaccineApplicationPayload): Promise<ApiResponse<VaccineApplication>> {
    return post('/vaccine-applications', payload)
  },

  update(id: string, payload: UpdateVaccineApplicationPayload): Promise<ApiResponse<VaccineApplication>> {
    return put(`/vaccine-applications/${id}`, payload)
  },

  remove(id: string): Promise<{ message: string }> {
    return del(`/vaccine-applications/${id}`)
  }
}
