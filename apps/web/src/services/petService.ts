import type { Pet, CreatePetPayload, UpdatePetPayload } from '@/types/pet'
import type { ApiResponse, PaginatedResponse } from '@/types/shared'
import { get, post, put, del } from '@/services/http'

const PER_PAGE_DEFAULT = 10

export const petService = {
  getAll(page = 1, perPage = PER_PAGE_DEFAULT): Promise<PaginatedResponse<Pet>> {
    return get(`/pets?page=${page}&per_page=${perPage}`)
  },

  getById(id: string): Promise<ApiResponse<Pet>> {
    return get(`/pets/${id}`)
  },

  create(payload: CreatePetPayload): Promise<ApiResponse<Pet>> {
    return post('/pets', payload)
  },

  update(id: string, payload: UpdatePetPayload): Promise<ApiResponse<Pet>> {
    return put(`/pets/${id}`, payload)
  },

  remove(id: string): Promise<{ message: string }> {
    return del(`/pets/${id}`)
  },
}
