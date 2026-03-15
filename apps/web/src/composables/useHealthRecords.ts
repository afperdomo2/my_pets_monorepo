import { computed, nextTick, type Ref } from "vue";
import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
import { healthRecordService } from "@/services/healthRecordService";
import type { CreateHealthRecordPayload, UpdateHealthRecordPayload, UpdateStatusPayload } from "@/types/healthRecord";

// ═══════════════════════════════════════════════════════════════════════════════════
// Cache time constants — ajusta estos valores según los SLAs del módulo
// ═══════════════════════════════════════════════════════════════════════════════════
/** Stale time para listados paginados (en ms). */
const HEALTH_RECORD_LIST_STALE_TIME = 2 * 60_000; // 2 minutos
/** Stale time para próximos registros (en ms). */
const HEALTH_RECORD_UPCOMING_STALE_TIME = 5 * 60_000; // 5 minutos

export function useGetAllHealthRecords(page: Ref<number>, perPage: Ref<number>) {
  const queryClient = useQueryClient();
  const queryKey = computed(() => ["health-records", "all", { page: page.value, perPage: perPage.value }]);

  const query = useQuery({
    queryKey,
    queryFn: () => healthRecordService.getAll(page.value, perPage.value),
    staleTime: HEALTH_RECORD_LIST_STALE_TIME,
  });

  let refreshing = false;
  async function refresh() {
    if (refreshing) return;
    refreshing = true;
    try {
      page.value = 1;
      await nextTick();
      await queryClient.invalidateQueries({ queryKey: ["health-records", "all"], refetchType: "active" });
    } finally {
      refreshing = false;
    }
  }

  return { ...query, refresh };
}

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
      // Invalidar registros de la mascota afectada
      queryClient.invalidateQueries({ queryKey: ["health-records", "pet", newItem.pet_id] });
      // Invalidar el listado global (VaccineHistoryTable)
      queryClient.invalidateQueries({ queryKey: ["health-records", "all"] });
      // Invalidar próximas vacunas/registros (VaccinesView upcoming)
      queryClient.invalidateQueries({ queryKey: ["health-records", "upcoming"] });
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
    mutationFn: (variables: { id: string; petId: string }) => healthRecordService.remove(variables.id),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["health-records", "pet", variables.petId] });
    },
  });
}

/**
 * Obtiene los próximos registros pendientes de aplicación.
 * @param limit - Cantidad de registros a retornar (default: 10, max: 50)
 * @param category - Filtrar por categoría (opcional): 'vaccine', 'deworming', 'exam'
 */
export function useGetUpcomingRecords(limit = 10, category?: string) {
  const queryClient = useQueryClient();
  const queryKey = computed(() => ["health-records", "upcoming", { limit, category }]);

  const query = useQuery({
    queryKey,
    queryFn: () => healthRecordService.getUpcoming(limit, category),
    staleTime: HEALTH_RECORD_UPCOMING_STALE_TIME,
  });

  const records = computed(() => query.data.value?.data ?? []);
  const total = computed(() => query.data.value?.total ?? 0);

  let refreshing = false;
  async function refresh() {
    if (refreshing) return;
    refreshing = true;
    try {
      await queryClient.invalidateQueries({ queryKey: ["health-records", "upcoming"], refetchType: "active" });
    } finally {
      refreshing = false;
    }
  }

  return {
    ...query,
    records,
    total,
    refresh,
  };
}

/**
 * Obtiene las próximas vacunas pendientes de aplicación.
 * @param limit - Cantidad de registros a retornar (default: 10, max: 50)
 */
export function useGetUpcomingVaccines(limit = 10) {
  return useGetUpcomingRecords(limit, 'vaccine');
}
