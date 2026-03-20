import { computed, nextTick, ref, watch, type Ref } from "vue";
import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
import { healthRecordService } from "@/services/healthRecordService";
import type { HealthRecord, CreateHealthRecordPayload, UpdateHealthRecordPayload } from "@/types/healthRecord";

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
      // Invalidar el listado global
      queryClient.invalidateQueries({ queryKey: ["health-records", "all"] });
      // Invalidar próximos registros
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
 * Obtiene los próximos registros con próxima dosis programada.
 */
export function useGetUpcomingRecordsPaged(perPage = 5, category?: string) {
  const queryClient = useQueryClient();

  const currentPage = ref(1);
  const totalRef = ref(0);
  const totalPagesRef = ref(0);
  const isLoadingMore = ref(false);

  const queryKey = computed(() => ["health-records", "upcoming", { page: currentPage.value, perPage, category }]);

  const cachedPage1 = queryClient.getQueryData<{ data: HealthRecord[]; total: number; total_pages: number }>(
    ["health-records", "upcoming", { page: 1, perPage, category }]
  );
  const accumulatedRecords = ref<HealthRecord[]>(cachedPage1?.data ?? []);
  if (cachedPage1) {
    totalRef.value = cachedPage1.total;
    totalPagesRef.value = cachedPage1.total_pages;
  }

  const query = useQuery({
    queryKey,
    queryFn: () => healthRecordService.getUpcoming(currentPage.value, perPage, category),
    staleTime: HEALTH_RECORD_UPCOMING_STALE_TIME,
  });

  watch(query.data, (newData) => {
    if (!newData) return;
    totalRef.value = newData.total;
    totalPagesRef.value = newData.total_pages;
    if (currentPage.value === 1) {
      accumulatedRecords.value = newData.data;
    } else {
      accumulatedRecords.value = [...accumulatedRecords.value, ...newData.data];
    }
    isLoadingMore.value = false;
  });

  const hasMore = computed(() => currentPage.value < totalPagesRef.value);

  async function loadMore() {
    if (!hasMore.value || isLoadingMore.value || query.isFetching.value) return;
    isLoadingMore.value = true;
    currentPage.value++;
  }

  async function refresh() {
    currentPage.value = 1;
    accumulatedRecords.value = [];
    await nextTick();
    await queryClient.invalidateQueries({ queryKey: ["health-records", "upcoming"], refetchType: "active" });
  }

  const isLoading = computed(() => query.isLoading.value && accumulatedRecords.value.length === 0);

  return {
    records: accumulatedRecords,
    total: totalRef,
    totalPages: totalPagesRef,
    currentPage,
    hasMore,
    loadMore,
    refresh,
    isLoading,
    isLoadingMore,
    isFetching: computed(() => query.isFetching.value),
    isError: computed(() => query.isError.value),
  };
}

/**
 * Obtiene las próximas vacunas pendientes con paginación acumulativa.
 */
export function useGetUpcomingVaccinesPaged(perPage = 5) {
  return useGetUpcomingRecordsPaged(perPage, 'vaccine');
}
