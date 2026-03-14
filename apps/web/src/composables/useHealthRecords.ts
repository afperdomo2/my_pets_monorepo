import { computed, nextTick, type Ref } from "vue";
import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
import { healthRecordService } from "@/services/healthRecordService";
import type { CreateHealthRecordPayload, UpdateHealthRecordPayload, UpdateStatusPayload } from "@/types/healthRecord";

// ═══════════════════════════════════════════════════════════════════════════════════
// Cache time constants — ajusta estos valores según los SLAs del módulo
// ═══════════════════════════════════════════════════════════════════════════════════
/** Stale time para listados paginados (en ms). */
const HEALTH_RECORD_LIST_STALE_TIME = 2 * 60_000; // 2 minutos

export function useGetHealthRecordsByPet(petId: Ref<string>, page: Ref<number>, perPage: Ref<number>) {
  const queryClient = useQueryClient();
  const queryKey = computed(() => ["health-records", "pet", petId.value, { page: page.value, perPage: perPage.value }]);

  const query = useQuery({
    queryKey,
    queryFn: () => healthRecordService.getByPetId(petId.value, page.value, perPage.value),
    staleTime: HEALTH_RECORD_LIST_STALE_TIME,
    enabled: computed(() => !!petId.value),
  });

  let refreshing = false;
  async function refresh() {
    if (refreshing) return;
    refreshing = true;
    try {
      page.value = 1;
      await nextTick();
      await queryClient.invalidateQueries({ queryKey: ["health-records", "pet", petId.value], refetchType: "active" });
    } finally {
      refreshing = false;
    }
  }

  return { ...query, refresh };
}

export function useGetHealthRecordsByPetAndCategory(petId: Ref<string>, category: Ref<string>, page: Ref<number>, perPage: Ref<number>) {
  const queryClient = useQueryClient();
  const queryKey = computed(() => ["health-records", "pet", petId.value, "category", category.value, { page: page.value, perPage: perPage.value }]);

  const query = useQuery({
    queryKey,
    queryFn: () => healthRecordService.getByPetIdAndCategory(petId.value, category.value, page.value, perPage.value),
    staleTime: HEALTH_RECORD_LIST_STALE_TIME,
    enabled: computed(() => !!petId.value && !!category.value),
  });

  let refreshing = false;
  async function refresh() {
    if (refreshing) return;
    refreshing = true;
    try {
      page.value = 1;
      await nextTick();
      await queryClient.invalidateQueries({ queryKey: ["health-records", "pet", petId.value, "category", category.value], refetchType: "active" });
    } finally {
      refreshing = false;
    }
  }

  return { ...query, refresh };
}

export function useCreateHealthRecord() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateHealthRecordPayload) => healthRecordService.create(payload).then((r) => r.data),
    onSuccess: (newItem) => {
      // Invalidar registros de la mascota para que las tablas se refresquen
      queryClient.invalidateQueries({ queryKey: ["health-records", "pet", newItem.pet_id] });
    },
  });
}

export function useUpdateHealthRecord() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: UpdateHealthRecordPayload }) =>
      healthRecordService.update(id, payload).then((r) => r.data),
    onSuccess: (updatedItem) => {
      queryClient.invalidateQueries({ queryKey: ["health-records", "pet", updatedItem.pet_id] });
    },
  });
}

export function useUpdateHealthRecordStatus() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: UpdateStatusPayload }) =>
      healthRecordService.updateStatus(id, payload).then((r) => r.data),
    onSuccess: (updatedItem) => {
      queryClient.invalidateQueries({ queryKey: ["health-records", "pet", updatedItem.pet_id] });
    },
  });
}

export function useDeleteHealthRecord() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, petId }: { id: string; petId: string }) => healthRecordService.remove(id),
    onSuccess: (_, variables) => {
      queryClient.invalidateQueries({ queryKey: ["health-records", "pet", variables.petId] });
    },
  });
}
