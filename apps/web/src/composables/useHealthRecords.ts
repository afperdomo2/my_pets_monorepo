import { computed, nextTick, ref, watch, type Ref } from "vue";
import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
import { healthRecordService } from "@/services/healthRecordService";
import type { HealthRecord, CreateHealthRecordPayload, UpdateHealthRecordPayload, UpdateStatusPayload } from "@/types/healthRecord";

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
 * Obtiene los próximos registros pendientes con paginación acumulativa ("ver más").
 * Cada llamada a loadMore() carga la siguiente página y concatena los resultados.
 * Al invalidar el cache (ej: tras crear un registro), se reinicia a la página 1.
 *
 * @param perPage - Registros por página (default: 5)
 * @param category - Filtrar por categoría (opcional): 'vaccine', 'deworming', 'exam'
 */
export function useGetUpcomingRecordsPaged(perPage = 5, category?: string) {
  const queryClient = useQueryClient();

  // Página actual — siempre arranca en 1 al montar
  const currentPage = ref(1);
  const totalRef = ref(0);
  const totalPagesRef = ref(0);
  const isLoadingMore = ref(false);

  const queryKey = computed(() => ["health-records", "upcoming", { page: currentPage.value, perPage, category }]);

  // Inicializar desde cache si ya existe data para página 1 (evita flash de empty state al volver)
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

  // Acumular resultados cuando llega nueva data
  watch(query.data, (newData) => {
    if (!newData) return;
    totalRef.value = newData.total;
    totalPagesRef.value = newData.total_pages;
    if (currentPage.value === 1) {
      // Reset: reemplazar todo (ocurre tras invalidación del cache o al montar)
      accumulatedRecords.value = newData.data;
    } else {
      // Append: agregar nuevos registros a los existentes
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

  // Al invalidar el cache desde afuera (ej: useCreateHealthRecord), resetear a página 1
  async function refresh() {
    currentPage.value = 1;
    accumulatedRecords.value = [];
    await nextTick();
    await queryClient.invalidateQueries({ queryKey: ["health-records", "upcoming"], refetchType: "active" });
  }

  // isLoading solo es true cuando no hay data en cache y está fetching (primera carga real)
  // Si hay datos acumulados, nunca mostrar el spinner de carga inicial
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
 * Obtiene las próximas vacunas pendientes con paginación acumulativa ("ver más").
 * @param perPage - Registros por página (default: 5)
 */
export function useGetUpcomingVaccinesPaged(perPage = 5) {
  return useGetUpcomingRecordsPaged(perPage, 'vaccine');
}

/**
 * @deprecated Usar useGetUpcomingRecordsPaged en su lugar.
 * Obtiene los próximos registros pendientes de aplicación.
 * @param limit - Cantidad de registros a retornar (default: 10, max: 50)
 * @param category - Filtrar por categoría (opcional): 'vaccine', 'deworming', 'exam'
 */
export function useGetUpcomingRecords(limit = 10, category?: string) {
  const queryClient = useQueryClient();
  const queryKey = computed(() => ["health-records", "upcoming", { limit, category }]);

  const query = useQuery({
    queryKey,
    queryFn: () => healthRecordService.getUpcoming(1, limit, category),
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
 * @deprecated Usar useGetUpcomingVaccinesPaged en su lugar.
 * Obtiene las próximas vacunas pendientes de aplicación.
 * @param limit - Cantidad de registros a retornar (default: 10, max: 50)
 */
export function useGetUpcomingVaccines(limit = 10) {
  return useGetUpcomingRecords(limit, 'vaccine');
}
