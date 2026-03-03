import { useQuery, useMutation, useQueryClient } from '@tanstack/vue-query'
import { userService } from '@/services/userService'
import type { CreateUserPayload, UpdateUserPayload } from '@/types/user'

const USERS_KEY = ['users'] as const

export function useGetUsers() {
  return useQuery({
    queryKey: USERS_KEY,
    queryFn: () => userService.getAll().then((r) => r.data),
  })
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
