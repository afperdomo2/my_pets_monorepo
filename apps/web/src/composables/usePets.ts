import { ref, computed } from 'vue'
import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { petService } from '@/services/petService'
import type { Pet, PetPayload } from '@/types/pet'

const PER_PAGE = 10

export function useGetPets() {
  const page = ref(1)
  // Acumulación de todas las mascotas cargadas (load-more)
  const allPets = ref<Pet[]>([])
  const totalPages = ref(1)
  const total = ref(0)

  const query = useQuery({
    queryKey: computed(() => ['pets', { page: page.value, perPage: PER_PAGE }]),
    queryFn: () => petService.getAll(page.value, PER_PAGE),
    staleTime: 60_000,
  })

  // Acumular resultados cuando cambia la query
  const hasMore = computed(() => page.value < totalPages.value)

  // Observar cambios en los datos para acumular
  const data = computed(() => {
    const res = query.data.value
    if (res) {
      total.value = res.total
      totalPages.value = res.total_pages
      // Merge: si la página es 1 reemplazar, si no, acumular
      const incoming = res.data
      if (page.value === 1) {
        allPets.value = incoming
      } else {
        // Evitar duplicados
        const existingIds = new Set(allPets.value.map((p) => p.id))
        const newItems = incoming.filter((p) => !existingIds.has(p.id))
        allPets.value = [...allPets.value, ...newItems]
      }
    }
    return allPets.value
  })

  function loadMore() {
    if (page.value < totalPages.value) {
      page.value++
    }
  }

  let refreshing = false
  // Refrescar: reinicia a página 1 y limpia acumulación
  function refresh() {
    if (refreshing) return
    refreshing = true
    allPets.value = []
    page.value = 1
    void query.refetch().finally(() => {
      refreshing = false
    })
  }

  return {
    data,
    allPets,
    total,
    totalPages,
    hasMore,
    page,
    isLoading: query.isLoading,
    isFetching: query.isFetching,
    isError: query.isError,
    error: query.error,
    loadMore,
    refresh,
  }
}

export function useGetPet(id: string) {
  return useQuery({
    queryKey: ['pets', id],
    queryFn: () => petService.getById(id).then((r) => r.data),
    staleTime: 5 * 60_000,
  })
}

export function useCreatePet() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: PetPayload) => petService.create(payload).then((r) => r.data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['pets'] })
    },
  })
}

export function useUpdatePet() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: PetPayload }) =>
      petService.update(id, payload).then((r) => r.data),
    onSuccess: (updatedPet) => {
      queryClient.setQueryData(['pets', updatedPet.id], updatedPet)
      queryClient.invalidateQueries({ queryKey: ['pets'] })
    },
  })
}

export function useDeletePet() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => petService.remove(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['pets'] })
    },
  })
}
