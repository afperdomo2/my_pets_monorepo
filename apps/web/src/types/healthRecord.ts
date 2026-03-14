export interface HealthRecord {
  id: string
  pet_id: string
  user_id: string
  health_catalog_id: string | null
  category: string
  name: string
  status: string
  application_date: string | null
  due_date: string
  notes: string | null
  created_at: string
  updated_at: string
}

export interface CreateHealthRecordPayload {
  pet_id: string
  health_catalog_id?: string
  category?: string
  name?: string
  status?: string
  application_date?: string
  due_date: string
  notes?: string
}

export interface UpdateHealthRecordPayload {
  category: string
  name: string
  status?: string
  application_date?: string
  due_date: string
  notes?: string
}

export interface UpdateStatusPayload {
  status: string
  application_date?: string
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
