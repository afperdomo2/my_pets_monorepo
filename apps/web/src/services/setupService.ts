import type { User, UserPayload } from '@/types/user'

const BASE_URL = '/api/v1'

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${BASE_URL}${url}`, {
    headers: { 'Content-Type': 'application/json' },
    credentials: 'include',
    ...options,
  })
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: 'Unknown error' }))
    throw new Error(err.error ?? `HTTP ${res.status}`)
  }
  return res.json()
}

export interface SetupStatusResponse {
  needs_setup: boolean
}

export interface SetupPayload extends UserPayload {
  password: string
}

export const setupService = {
  checkStatus(): Promise<SetupStatusResponse> {
    return request('/setup/status')
  },

  createFirstUser(payload: SetupPayload): Promise<{ data: User }> {
    return request('/setup', {
      method: 'POST',
      body: JSON.stringify(payload),
    })
  },
}
