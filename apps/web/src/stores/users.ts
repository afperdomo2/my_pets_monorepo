import { ref } from 'vue'
import { defineStore } from 'pinia'
import type { User, CreateUserPayload, UpdateUserPayload } from '@/types/user'
import { userService } from '@/services/userService'

export const useUserStore = defineStore('users', () => {
  const users = ref<User[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchUsers() {
    loading.value = true
    error.value = null
    try {
      const res = await userService.getAll()
      users.value = res.data
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Error al cargar usuarios'
    } finally {
      loading.value = false
    }
  }

  async function createUser(payload: CreateUserPayload): Promise<void> {
    const res = await userService.create(payload)
    users.value.push(res.data)
  }

  async function updateUser(id: number, payload: UpdateUserPayload): Promise<void> {
    const res = await userService.update(id, payload)
    const idx = users.value.findIndex((u) => u.id === id)
    if (idx !== -1) users.value[idx] = res.data
  }

  async function deleteUser(id: number): Promise<void> {
    await userService.remove(id)
    users.value = users.value.filter((u) => u.id !== id)
  }

  return { users, loading, error, fetchUsers, createUser, updateUser, deleteUser }
})
