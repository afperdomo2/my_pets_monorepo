import { computed, type Ref } from "vue";
import { useQuery, useMutation, useQueryClient } from "@tanstack/vue-query";
import { vaccineApplicationService } from "@/services/vaccineApplicationService";
import type { VaccineApplication, CreateVaccineApplicationPayload, UpdateVaccineApplicationPayload } from "@/types/vaccineApplication";

/**
 * Obtiene las aplicaciones de un health_record.
 */
export function useGetVaccineApplicationsByHealthRecord(healthRecordId: Ref<string>) {
  const queryKey = computed(() => ["vaccine-applications", "health-record", healthRecordId.value]);

  return useQuery({
    queryKey,
    queryFn: () => vaccineApplicationService.getByHealthRecordId(healthRecordId.value).then((r) => r.data),
    enabled: computed(() => !!healthRecordId.value),
  });
}

/**
 * Obtiene una aplicación por ID.
 */
export function useGetVaccineApplication(id: Ref<string>) {
  const queryKey = computed(() => ["vaccine-applications", id.value]);

  return useQuery({
    queryKey,
    queryFn: () => vaccineApplicationService.getById(id.value).then((r) => r.data),
    enabled: computed(() => !!id.value),
  });
}

/**
 * Crea una nueva aplicación de vacuna.
 */
export function useCreateVaccineApplication() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (payload: CreateVaccineApplicationPayload) =>
      vaccineApplicationService.create(payload).then((r) => r.data),
    onSuccess: (newItem) => {
      // Invalidar aplicaciones del health_record afectado
      queryClient.invalidateQueries({ queryKey: ["vaccine-applications", "health-record", newItem.health_record_id] });
      // Invalidar health_records para actualizar la vista
      queryClient.invalidateQueries({ queryKey: ["health-records"] });
    },
  });
}

/**
 * Actualiza una aplicación de vacuna.
 */
export function useUpdateVaccineApplication() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: ({ id, payload }: { id: string; payload: UpdateVaccineApplicationPayload }) =>
      vaccineApplicationService.update(id, payload).then((r) => r.data),
    onSuccess: (updatedItem) => {
      queryClient.invalidateQueries({ queryKey: ["vaccine-applications", "health-record", updatedItem.health_record_id] });
    },
  });
}

/**
 * Elimina una aplicación de vacuna.
 */
export function useDeleteVaccineApplication() {
  const queryClient = useQueryClient();

  return useMutation({
    mutationFn: (variables: { id: string; healthRecordId: string }) =>
      vaccineApplicationService.remove(variables.id),
    onSuccess: (_data, variables) => {
      queryClient.invalidateQueries({ queryKey: ["vaccine-applications", "health-record", variables.healthRecordId] });
    },
  });
}
