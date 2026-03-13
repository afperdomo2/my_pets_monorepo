import { type Ref, computed, nextTick } from "vue";
import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
import { vaccineCatalogService } from "@/services/vaccineCatalogService";
import type { CreateVaccineCatalogPayload, UpdateVaccineCatalogPayload } from "@/types/vaccineCatalog";

// ── Cache constants ────────────────────────────────────────────
// Controla cuánto tiempo los datos se consideran frescos (no "stale")
// Aumentar estos valores para reducir fetches innecesarios
const VACCINES_STALE_TIME = 2 * 60_000; // 2 minutos para listado
const VACCINE_STALE_TIME = 5 * 60_000; // 5 minutos para detalle individual

export function useGetVaccinesCatalog(page: Ref<number>, perPage: Ref<number>, species: Ref<string | undefined>) {
  const queryClient = useQueryClient();
  const queryKey = computed(() => ["vaccines-catalog", { page: page.value, perPage: perPage.value, species: species.value }]);

  const query = useQuery({
    queryKey,
    queryFn: () => vaccineCatalogService.getAll(page.value, perPage.value, species.value),
    staleTime: VACCINES_STALE_TIME,
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

      // 3. Invalidar TODO el caché de vaccines-catalog y refetchear solo el observer activo
      //    (página 1). refetchType:'active' garantiza exactamente un fetch.
      await queryClient.invalidateQueries({ queryKey: ["vaccines-catalog"], refetchType: "active" });
    } finally {
      refreshing = false;
    }
  }

  return { ...query, refresh };
}

export function useGetVaccineCatalog(id: string) {
  return useQuery({
    queryKey: ["vaccines-catalog", id],
    queryFn: () => vaccineCatalogService.getById(id).then((r) => r.data),
    staleTime: VACCINE_STALE_TIME,
  });
}

export function useCreateVaccineCatalog() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (payload: CreateVaccineCatalogPayload) => vaccineCatalogService.create(payload).then((r) => r.data),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["vaccines-catalog"] });
    },
  });
}

export function useUpdateVaccineCatalog() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: UpdateVaccineCatalogPayload }) =>
      vaccineCatalogService.update(id, payload).then((r) => r.data),
    onSuccess: (updatedVaccine) => {
      queryClient.setQueryData(["vaccines-catalog", updatedVaccine.id], updatedVaccine);
      queryClient.invalidateQueries({ queryKey: ["vaccines-catalog"] });
    },
  });
}

export function useDeleteVaccineCatalog() {
  const queryClient = useQueryClient();
  return useMutation({
    mutationFn: (id: string) => vaccineCatalogService.remove(id),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["vaccines-catalog"] });
    },
  });
}
