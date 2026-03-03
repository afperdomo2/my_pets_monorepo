import { defineStore } from 'pinia'
import { ref } from 'vue'
import { petService } from '@/services/petService'
import type { Pet, PetPayload } from '@/types/pet'

export const usePetStore = defineStore('pets', () => {
  const pets = ref<Pet[]>([])
  const loading = ref(false)
  const error = ref<string | null>(null)

  async function fetchPets() {
    loading.value = true
    error.value = null
    try {
      const res = await petService.getAll()
      pets.value = res.data
    } catch (e) {
      error.value = (e as Error).message
    } finally {
      loading.value = false
    }
  }

  async function createPet(payload: PetPayload) {
    const res = await petService.create(payload)
    pets.value.push(res.data)
    return res.data
  }

  async function updatePet(id: string, payload: PetPayload) {
    const res = await petService.update(id, payload)
    const idx = pets.value.findIndex(p => p.id === id)
    if (idx !== -1) pets.value[idx] = res.data
    return res.data
  }

  async function deletePet(id: string) {
    await petService.remove(id)
    pets.value = pets.value.filter(p => p.id !== id)
  }

  return { pets, loading, error, fetchPets, createPet, updatePet, deletePet }
})
