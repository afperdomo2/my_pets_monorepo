import type {
  Exam,
  ExamWithResults,
  CreateExamPayload,
  UpdateExamPayload,
  ScheduleExamPayload,
  CompleteExamPayload,
} from '@/types/exam'
import type { PaginatedResponse, ApiResponse } from '@/types/shared'
import { get, post, put, patch, del } from '@/services/http'

const PER_PAGE_DEFAULT = 10

export const examService = {
  // Obtener todos los exámenes del usuario
  getAll(page = 1, perPage = PER_PAGE_DEFAULT): Promise<PaginatedResponse<Exam>> {
    return get(`/exams?page=${page}&per_page=${perPage}`)
  },

  // Obtener exámenes de una mascota
  getByPetId(petId: string, page = 1, perPage = PER_PAGE_DEFAULT): Promise<PaginatedResponse<Exam>> {
    return get(`/exams/pets/${petId}?page=${page}&per_page=${perPage}`)
  },

  // Obtener examen por ID con resultados
  getById(id: string): Promise<ApiResponse<ExamWithResults>> {
    return get(`/exams/${id}`)
  },

  create(payload: CreateExamPayload): Promise<ApiResponse<Exam>> {
    return post('/exams', payload)
  },

  update(id: string, payload: UpdateExamPayload): Promise<ApiResponse<Exam>> {
    return put(`/exams/${id}`, payload)
  },

  // Programar examen
  schedule(id: string, payload: ScheduleExamPayload): Promise<ApiResponse<Exam>> {
    return patch(`/exams/${id}/schedule`, payload)
  },

  // Completar examen
  complete(id: string, payload: CompleteExamPayload): Promise<ApiResponse<Exam>> {
    return patch(`/exams/${id}/complete`, payload)
  },

  remove(id: string): Promise<{ message: string }> {
    return del(`/exams/${id}`)
  }
}
