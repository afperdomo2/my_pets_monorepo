import { type Ref, computed, nextTick } from "vue";
import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
import { healthCatalogService } from "@/services/healthCatalogService";
import type { CreateHealthCatalogPayload, UpdateHealthCatalogPayload } from "@/types/healthCatalog";

// ═══════════════════════════════════════════════════════════════════════════════════
// Cache time constants — ajusta estos valores según los SLAs del módulo
// ═══════════════════════════════════════════════════════════════════════════════════
/** Stale time para listados paginados (en ms). Reduce para refrescar más a menudo. */
const HEALTH_CATALOG_LIST_STALE_TIME = 2 * 60_000; // 2 minutos para listado

/** Stale time para detalles individuales (en ms). Más largo porque cambios son menos frecuentes. */
const HEALTH_CATALOG_ITEM_STALE_TIME = 5 * 60_000; // 5 minutos para detalle individual

export function useGetHealthCatalogs(category: Ref<string>, page: Ref<number>, perPage: Ref<number>, species: Ref<string | undefined>) {
  const queryClient = useQueryClient();
  const queryKey = computed(() => ["health-catalog", { category: category.value, page: page.value, perPage: perPage.value, species: species.value }]);

  const query = useQuery({
    queryKey,
    queryFn: () => healthCatalogService.getAll(category.value, page.value, perPage.value, species.value),
    staleTime: HEALTH_CATALOG_LIST_STALE_TIME,
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

      // 3. Invalidar TODO el caché de health-catalog y refetchear solo el observer activo
      //    (página 1). refetchType:'active' garantiza exactamente un fetch.
      await queryClient.invalidateQueries({ queryKey: ["health-catalog"], refetchType: "active" });
    } finally {
      refreshing = false;
    }
  }

  return { ...query, refresh };
}

export function useGetHealthCatalog(id: string) {
  return useQuery({
    queryKey: ["health-catalog", id],
    queryFn: () => healthCatalogService.getById(id).then((r) => r.data),
    staleTime: HEALTH_CATALOG_ITEM_STALE_TIME,
  });
}

export function useCreateHealthCatalog() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateHealthCatalogPayload) => healthCatalogService.create(payload).then((r) => r.data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["health-catalog"] });
    },
  });
}

export function useUpdateHealthCatalog() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: UpdateHealthCatalogPayload }) =>
      healthCatalogService.update(id, payload).then((r) => r.data),
    onSuccess: (updatedItem) => {
      queryClient.setQueryData(["health-catalog", updatedItem.id], updatedItem);
      queryClient.invalidateQueries({ queryKey: ["health-catalog"] });
    },
  });
}

export function useDeleteHealthCatalog() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => healthCatalogService.remove(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["health-catalog"] });
    },
  });
}
