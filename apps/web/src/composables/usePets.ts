import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { petService } from '@/services/petService'
import type { PetPayload } from '@/types/pet'

const PETS_KEY = ['pets'] as const

export function useGetPets() {
  return useQuery({
    queryKey: PETS_KEY,
    queryFn: () => petService.getAll().then((r) => r.data),
  })
}

export function useGetPet(id: string) {
  return useQuery({
    queryKey: [...PETS_KEY, id],
    queryFn: () => petService.getById(id).then((r) => r.data),
    staleTime: 5 * 60_000, // detail: 5 min
  })
}

export function useCreatePet() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: PetPayload) => petService.create(payload).then((r) => r.data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PETS_KEY })
    },
  })
}

export function useUpdatePet() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: PetPayload }) =>
      petService.update(id, payload).then((r) => r.data),
    onSuccess: (updatedPet) => {
      // Actualiza el detail en cache inmediatamente
      queryClient.setQueryData([...PETS_KEY, updatedPet.id], updatedPet)
      // Invalida la lista para que refleje el cambio
      queryClient.invalidateQueries({ queryKey: PETS_KEY })
    },
  })
}

export function useDeletePet() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => petService.remove(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: PETS_KEY })
    },
  })
}
