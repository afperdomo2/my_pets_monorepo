import type { User, UserWithPetCount, CreateUserPayload, UpdateUserPayload } from '@/types/user'
import type { ApiResponse, PaginatedResponse } from '@/types/shared'
import { get, post, put, del } from '@/services/http'

const PER_PAGE_DEFAULT = 10

export const userService = {
  getAll(page = 1, perPage = PER_PAGE_DEFAULT): Promise<PaginatedResponse<UserWithPetCount>> {
    return get(`/users?page=${page}&per_page=${perPage}`)
  },

  getById(id: string): Promise<ApiResponse<UserWithPetCount>> {
    return get(`/users/${id}`)
  },

  create(payload: CreateUserPayload): Promise<ApiResponse<User>> {
    return post('/users', payload)
  },

  update(id: string, payload: UpdateUserPayload): Promise<ApiResponse<User>> {
    return put(`/users/${id}`, payload)
  },

  remove(id: string): Promise<{ message: string }> {
    return del(`/users/${id}`)
  },
}
