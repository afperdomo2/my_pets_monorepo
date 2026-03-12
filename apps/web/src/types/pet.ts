export interface Pet {
  id: string
  user_id: string
  name: string
  species: string
  breed: string
  birth_date: string        // ISO "YYYY-MM-DDT..." from Go time.Time
  birth_date_exact: boolean
  weight_grams: number | null
  life_stage: string | null // null for non-dog/cat species
  created_at: string
  updated_at: string
}

export interface CreatePetPayload {
  name: string
  species: string
  breed?: string
  birth_date: string        // "YYYY-MM-DD"
  birth_date_exact: boolean
  weight_grams?: number     // grams, only at creation
  life_stage?: string       // only at creation, only dog/cat
}

export interface UpdatePetPayload {
  name: string
  species: string
  breed?: string
  birth_date: string
  birth_date_exact: boolean
  // weight_grams and life_stage intentionally excluded
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
