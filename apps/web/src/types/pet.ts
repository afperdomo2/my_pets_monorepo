export interface Pet {
  id: string
  name: string
  species: string
  breed: string
  age: number
  owner: string
  created_at: string
  updated_at: string
}

export interface PetPayload {
  name: string
  species: string
  breed?: string
  age?: number
  owner?: string
}

export interface ApiResponse<T> {
  data: T
  total?: number
}

export interface PaginatedResponse<T> {
  data: T[]
  total: number
  page: number
  per_page: number
  total_pages: number
}
