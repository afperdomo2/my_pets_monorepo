export interface Pet {
  id: number
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
