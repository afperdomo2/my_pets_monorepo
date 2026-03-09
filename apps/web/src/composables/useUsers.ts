import { type Ref, computed, nextTick } from "vue";
import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
import { userService } from "@/services/userService";
import type { CreateUserPayload, UpdateUserPayload } from "@/types/user";

const USERS_STALE_TIME = 30 * 60_000;
const USER_STALE_TIME = 30 * 60_000;

export function useGetUsers(page: Ref<number>, perPage: Ref<number>) {
  const queryClient = useQueryClient();
  const queryKey = computed(() => ["users", { page: page.value, perPage: perPage.value }]);

  const query = useQuery({
    queryKey,
    queryFn: () => userService.getAll(page.value, perPage.value),
    staleTime: USERS_STALE_TIME,
  });

  let refreshing = false;
  async function refresh() {
    if (refreshing) return;
    refreshing = true;
    try {
      // 1. Resetear page sin invalidar aún: la query de página 1 no está stale,
      //    así que useQuery NO lanzará fetch reactivo al cambiar el queryKey.
      page.value = 1;

      // 2. Esperar a que Vue flush el watcher interno de TanStack que llama
      //    observer.setOptions() con el nuevo queryKey. Después de este tick,
      //    el observer ya apunta a la query de página 1 sin haber hecho fetch.
      await nextTick();

      // 3. Invalidar TODO el caché de users y refetchear solo el observer activo
      //    (página 1). refetchType:'active' garantiza exactamente un fetch.
      await queryClient.invalidateQueries({ queryKey: ["users"], refetchType: "active" });
    } finally {
      refreshing = false;
    }
  }

  return { ...query, refresh };
}

export function useGetUser(id: string) {
  return useQuery({
    queryKey: ["users", id],
    queryFn: () => userService.getById(id).then((r) => r.data),
    staleTime: USER_STALE_TIME,
  });
}

export function useCreateUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateUserPayload) => userService.create(payload).then((r) => r.data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  });
}

export function useUpdateUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: UpdateUserPayload }) =>
      userService.update(id, payload).then((r) => r.data),
    onSuccess: (updatedUser) => {
      queryClient.setQueryData(["users", updatedUser.id], updatedUser);
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  });
}

export function useDeleteUser() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => userService.remove(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["users"] });
    },
  });
}
