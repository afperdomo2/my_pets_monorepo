import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { userService } from '@/services/userService'
import type { CreateUserPayload, UpdateUserPayload } from '@/types/user'

const USERS_KEY = ['users'] as const
const USERS_STALE_TIME = 60_000 // 60 s — debe coincidir con defaultOptions en main.ts

export function useGetUsers() {
  const queryClient = useQueryClient()
  const query = useQuery({
    queryKey: USERS_KEY,
    queryFn: () => userService.getAll().then((r) => r.data),
  })

  // Respeta el caché: solo va a la API si los datos tienen más de USERS_STALE_TIME
  function refresh() {
    void queryClient.fetchQuery({
      queryKey: USERS_KEY,
      queryFn: () => userService.getAll().then((r) => r.data),
      staleTime: USERS_STALE_TIME,
    })
  }

  return { ...query, refresh }
}

export function useGetUser(id: string) {
  return useQuery({
    queryKey: [...USERS_KEY, id],
    queryFn: () => userService.getById(id).then((r) => r.data),
    staleTime: 5 * 60_000, // detail: 5 min
  })
}

export function useCreateUser() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (payload: CreateUserPayload) => userService.create(payload).then((r) => r.data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: USERS_KEY })
    },
  })
}

export function useUpdateUser() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: UpdateUserPayload }) =>
      userService.update(id, payload).then((r) => r.data),
    onSuccess: (updatedUser) => {
      queryClient.setQueryData([...USERS_KEY, updatedUser.id], updatedUser)
      queryClient.invalidateQueries({ queryKey: USERS_KEY })
    },
  })
}

export function useDeleteUser() {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => userService.remove(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: USERS_KEY })
    },
  })
}
