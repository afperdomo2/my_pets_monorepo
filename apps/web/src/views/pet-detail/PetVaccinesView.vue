<script setup lang="ts">
import ConfirmDeleteModal from '@/components/health-tabs/ConfirmDeleteModal.vue'
import HealthRecordCreateModal from '@/components/health-tabs/HealthRecordCreateModal.vue'
import HealthRecordEditModal from '@/components/health-tabs/HealthRecordEditModal.vue'
import VaccineApplicationModal from '@/components/health-tabs/VaccineApplicationModal.vue'
import AppPagination from '@/components/ui/AppPagination.vue'
import PerPageSelector from '@/components/ui/PerPageSelector.vue'
import {
  useDeleteHealthRecord,
  useGetHealthRecordsByPetAndCategory,
  useUpdateHealthRecord,
} from '@/composables/useHealthRecords'
import { useGetPet } from '@/composables/usePets'
import { useCreateVaccineApplication } from '@/composables/useVaccineApplications'
import { HealthCatalogCategory } from '@/constants/healthRecord'
import { formatDateOnly } from '@/utils/date'
import {
  IconCalendar,
  IconEdit,
  IconPlus,
  IconRefresh,
  IconTrash,
  IconVaccine,
} from '@tabler/icons-vue'
import { computed, ref, type Ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const petId = computed(() => String(route.params.id))

const page = ref(1)
const perPage = ref(10)
const category = ref(HealthCatalogCategory.Vaccine)

const { data: pet } = useGetPet(petId.value)

const { data, isLoading, isError, refresh } = useGetHealthRecordsByPetAndCategory(
  petId as Ref<string>,
  category as Ref<string>,
  page,
  perPage,
)

const records = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.total ?? 0)
const totalPages = computed(() => data.value?.total_pages ?? 0)

function formatDosesText(record: (typeof records.value)[number]): string {
  const applied = record.applied_doses_count ?? 0
  if (record.total_doses) {
    return `${applied} de ${record.total_doses}`
  }
  return String(applied)
}

const showVaccineModal = ref(false)
const showConfirmModal = ref(false)
const showEditModal = ref(false)
const showApplicationModal = ref(false)
const recordToDelete = ref<(typeof records.value)[number] | null>(null)
const recordToEdit = ref<(typeof records.value)[number] | null>(null)
const healthRecordToApply = ref<(typeof records.value)[number] | null>(null)
const deletingId = ref<string | null>(null)

const deleteRecord = useDeleteHealthRecord()
const updateRecord = useUpdateHealthRecord()
const createApplication = useCreateVaccineApplication()

function openCreate() {
  showVaccineModal.value = true
}

function openEdit(record: (typeof records.value)[number]) {
  recordToEdit.value = record
  showEditModal.value = true
}

function openApplicationModal(record: (typeof records.value)[number]) {
  healthRecordToApply.value = record
  showApplicationModal.value = true
}

function openDeleteConfirm(record: (typeof records.value)[number]) {
  recordToDelete.value = record
  showConfirmModal.value = true
}

async function handleDeleteConfirm() {
  if (!recordToDelete.value) return
  deletingId.value = recordToDelete.value.id
  try {
    await deleteRecord.mutateAsync({ id: recordToDelete.value.id, petId: petId.value })
    showConfirmModal.value = false
    recordToDelete.value = null
  } finally {
    deletingId.value = null
  }
}
</script>

<template>
  <div class="tab-view">
    <div class="content-card">
      <!-- Header -->
      <div class="tab-header">
        <div class="tab-header-left">
          <h2 class="tab-title">
            <IconVaccine :size="20" :stroke-width="1.75" />
            Carnet de Vacunación
          </h2>
          <span class="panel-count">{{ total }} registros</span>
        </div>
        <div class="tab-header-right">
          <button class="btn-refresh" :disabled="isLoading" title="Refrescar" @click="refresh">
            <IconRefresh :size="16" :stroke-width="2" :class="{ spin: isLoading }" />
            Refrescar
          </button>
          <button class="btn-add" @click="openCreate">
            <IconPlus :size="16" :stroke-width="2.5" />
            Registrar
          </button>
        </div>
      </div>

      <!-- Tabla desktop -->
      <div class="table-wrapper">
        <table v-if="!isLoading && !isError && records.length > 0" class="vaccine-table">
          <thead>
            <tr>
              <th>Vacuna</th>
              <th class="th-center"># Dosis</th>
              <th class="th-center">Última aplicación</th>
              <th class="th-center">Próxima dosis</th>
              <th class="th-center">Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="record in records" :key="record.id" class="vaccine-row">
              <td>
                <div class="vaccine-cell">
                  <span class="vaccine-name">{{ record.name }}</span>
                  <span v-if="record.notes" class="vaccine-note" :title="record.notes">{{
                    record.notes
                  }}</span>
                </div>
              </td>
              <td class="td-center">
                <span class="doses-count">{{ formatDosesText(record) }}</span>
              </td>
              <td class="td-center">
                <div class="last-application-cell">
                  <div class="last-application-date">
                    <span>{{ formatDateOnly(record.last_dose_date) }}</span>
                  </div>
                </div>
              </td>
              <td class="td-center">
                <span class="date-cell">
                  {{ formatDateOnly(record.next_dose_date) }}
                </span>
              </td>
              <td class="td-center">
                <div class="action-buttons">
                  <button
                    class="btn-action btn-apply"
                    title="Aplicar dosis"
                    :disabled="createApplication.isPending.value"
                    @click="openApplicationModal(record)"
                  >
                    <IconVaccine :size="14" :stroke-width="2" />
                    Aplicar
                  </button>
                  <button
                    class="btn-action btn-edit"
                    title="Editar"
                    :disabled="updateRecord.isPending.value"
                    @click="openEdit(record)"
                  >
                    <IconEdit :size="14" :stroke-width="2" />
                    Editar
                  </button>
                  <button
                    class="btn-action btn-delete"
                    title="Eliminar"
                    :disabled="deletingId === record.id || deleteRecord.isPending.value"
                    @click="openDeleteConfirm(record)"
                  >
                    <IconTrash :size="14" :stroke-width="2" />
                    Eliminar
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>

        <!-- Loading -->
        <div v-if="isLoading" class="loading-state">
          <div class="spinner" />
          <p>Cargando registros...</p>
        </div>

        <!-- Error -->
        <div v-else-if="isError" class="error-state">
          <p>Error al cargar los registros</p>
          <button class="btn-retry" @click="refresh">Reintentar</button>
        </div>

        <!-- Empty -->
        <div v-else-if="records.length === 0" class="empty-state">
          <IconVaccine :size="40" :stroke-width="1.5" />
          <p>No hay vacunas registradas</p>
          <button class="btn-add-empty" @click="openCreate">
            <IconPlus :size="16" :stroke-width="2.5" />
            Aplicar primera vacuna
          </button>
        </div>
      </div>

      <!-- Vista card para móvil -->
      <div v-if="!isLoading && !isError && records.length > 0" class="record-cards">
        <div v-for="record in records" :key="`card-${record.id}`" class="record-card">
          <div class="record-card__top">
            <span class="record-card__vaccine">{{ record.name }}</span>
            <div class="record-card__actions">
              <button
                class="btn-apply-card"
                title="Aplicar dosis"
                :disabled="createApplication.isPending.value"
                @click="openApplicationModal(record)"
              >
                <IconVaccine :size="14" :stroke-width="2" />
                Aplicar
              </button>
              <button
                class="btn-edit-card"
                title="Editar"
                :disabled="updateRecord.isPending.value"
                @click="openEdit(record)"
              >
                <IconEdit :size="14" :stroke-width="2" />
                Editar
              </button>
              <button
                class="btn-delete-card"
                title="Eliminar"
                :disabled="deletingId === record.id || deleteRecord.isPending.value"
                @click="openDeleteConfirm(record)"
              >
                <IconTrash :size="14" :stroke-width="2" />
                Eliminar
              </button>
            </div>
          </div>

          <!-- # Dosis -->
          <div class="record-card__doses">
            <span class="record-card__doses-label"># Dosis</span>
            <span class="record-card__doses-value">{{ formatDosesText(record) }}</span>
          </div>

          <!-- Última aplicación -->
          <div class="record-card__last-application">
            <div class="record-card__last-application-date">
              <IconCalendar :size="14" :stroke-width="2" />
              <span>{{ formatDateOnly(record.last_dose_date) }}</span>
            </div>
          </div>

          <!-- Próxima dosis -->
          <div class="record-card__next-dose">
            <span class="record-card__next-dose-label">Próxima dosis</span>
            <span class="record-card__next-dose-value">{{
              formatDateOnly(record.next_dose_date)
            }}</span>
          </div>
        </div>
      </div>

      <!-- Footer con paginación -->
      <div v-if="!isLoading && !isError && records.length > 0" class="table-footer">
        <PerPageSelector v-model="perPage" :options="[10, 25, 50]" />
        <AppPagination
          :current-page="page"
          :total-pages="totalPages"
          :total-items="total"
          :per-page="perPage"
          @update:page="page = $event"
        />
      </div>
    </div>

    <!-- Modal registrar vacuna -->
    <HealthRecordCreateModal
      v-if="showVaccineModal"
      :pet-id="petId"
      :pet-species="pet?.species ?? 'dog'"
      category="vaccine"
      @close="showVaccineModal = false"
    />

    <!-- Modal confirmar eliminación -->
    <ConfirmDeleteModal
      v-model="showConfirmModal"
      :record-name="recordToDelete?.name ?? ''"
      record-type="vaccine"
      :deleting="deleteRecord.isPending.value"
      @confirm="handleDeleteConfirm"
    />

    <!-- Modal editar registro -->
    <HealthRecordEditModal
      v-if="showEditModal && recordToEdit"
      :record="recordToEdit"
      category="vaccine"
      @close="showEditModal = false; recordToEdit = null"
      @updated="refresh"
    />

    <!-- Modal aplicar dosis -->
    <VaccineApplicationModal
      v-if="showApplicationModal && healthRecordToApply"
      :health-record-id="healthRecordToApply.id"
      category="vaccine"
      :total-doses="healthRecordToApply.total_doses"
      @close="showApplicationModal = false; healthRecordToApply = null"
      @applied="refresh"
    />
  </div>
</template>

<style scoped>
.tab-view {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.content-card {
  container-type: inline-size;
  container-name: vaccine-card;
  background: transparent;
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-xl);
  overflow: hidden;
  box-shadow: none;
}

/* ── Header ───────────────────────── */
.tab-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  flex-wrap: wrap;
  padding: var(--space-4) var(--space-5);
  border-bottom: 1px solid var(--color-border-light);
}

.tab-header-left {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.tab-title {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-family: var(--font-display);
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.panel-count {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  font-weight: 500;
}

.tab-header-right {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.btn-refresh {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  background: var(--color-bg-alt);
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition:
    background var(--transition-fast),
    color var(--transition-fast);
  flex-shrink: 0;
}

.btn-refresh:hover:not(:disabled) {
  background: var(--color-accent-light);
  color: var(--color-accent-dark);
  border-color: var(--color-accent);
}

.btn-refresh:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-refresh .spin {
  animation: spin 0.8s linear infinite;
}

.btn-add {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  background: var(--color-accent);
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast);
}

.btn-add:hover {
  background: var(--color-accent-hover);
}

/* ── Tabla ────────────────────────── */
.table-wrapper {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.vaccine-table {
  width: 100%;
  min-width: 700px;
  border-collapse: collapse;
}

.vaccine-table thead tr {
  background: var(--color-bg);
  border-bottom: 1px solid var(--color-border-light);
}

.vaccine-table th {
  padding: var(--space-3) var(--space-4);
  text-align: left;
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.vaccine-table th.th-center {
  text-align: center;
}

.vaccine-row {
  transition: background var(--transition-fast);
}

.vaccine-row:hover {
  background: var(--color-bg-alt);
}

.vaccine-row td {
  padding: var(--space-3) var(--space-4);
  border-bottom: 1px solid var(--color-border-light);
}

.vaccine-row:last-child td {
  border-bottom: none;
}

.td-center {
  text-align: center;
}

/* ── Celdas ───────────────────────── */
.vaccine-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.vaccine-name {
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  font-weight: 500;
}

.vaccine-note {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.4;
}

.date-cell {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  white-space: nowrap;
}

/* ── Celda # Dosis ────────────── */
.doses-count {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
}

/* ── Celda Última aplicación ────────────── */
.last-application-cell {
  display: flex;
  flex-direction: column;
  gap: 4px;
  align-items: center;
}

.last-application-date {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

/* ── Estilos para móvil ────────────── */
.record-card__doses {
  margin-top: var(--space-3);
  padding-top: var(--space-3);
  border-top: 1px solid var(--color-border-light);
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.record-card__doses-label {
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.record-card__doses-value {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
}

.record-card__last-application {
  margin-top: var(--space-3);
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.record-card__last-application-date {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.record-card__next-dose {
  margin-top: var(--space-3);
  padding-top: var(--space-3);
  border-top: 1px solid var(--color-border-light);
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.record-card__next-dose-label {
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.record-card__next-dose-value {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

/* ── Badges de estado ─────────────── */
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px var(--space-3);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: 600;
  white-space: nowrap;
}

.status--uptodate {
  background: #e8f5ee;
  color: #2e7d52;
}

.status--upcoming {
  background: #fef3e2;
  color: #c4714a;
}

.status--overdue {
  background: #fef2f2;
  color: #dc2626;
}

/* ── Botones de acción ────────────── */
.action-buttons {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  justify-content: center;
}

.btn-action {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-1);
  padding: var(--space-2) var(--space-3);
  border-radius: var(--radius-md);
  font-size: var(--text-xs);
  font-weight: 600;
  cursor: pointer;
  border: 1px solid transparent;
  transition:
    background var(--transition-fast),
    color var(--transition-fast),
    border-color var(--transition-fast);
  white-space: nowrap;
}

/* Aplicar - Bordered con fondo claro */
.btn-apply {
  background: #e0f2fe;
  color: #0284c7;
  border-color: #bae6fd;
}

.btn-apply:hover:not(:disabled) {
  background: #bae6fd;
  border-color: #0284c7;
}

/* Editar - Bordered con fondo claro */
.btn-edit {
  background: #fef3c7;
  color: #d97706;
  border-color: #fde68a;
}

.btn-edit:hover:not(:disabled) {
  background: #fde68a;
  border-color: #f59e0b;
}

/* Eliminar - Bordered con fondo claro */
.btn-delete {
  background: #fef2f2;
  color: #dc2626;
  border-color: #fecaca;
}

.btn-delete:hover:not(:disabled) {
  background: #fee2e2;
  border-color: #fca5a5;
}

/* Eliminar - Bordered con fondo claro */
.btn-delete-card {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-1);
  padding: var(--space-2) var(--space-3);
  background: #fef2f2;
  color: #dc2626;
  border: 1px solid #fecaca;
  border-radius: var(--radius-md);
  font-size: var(--text-xs);
  font-weight: 600;
  cursor: pointer;
  transition:
    background var(--transition-fast),
    color var(--transition-fast),
    border-color var(--transition-fast);
  white-space: nowrap;
}

.btn-delete-card:hover:not(:disabled) {
  background: #fee2e2;
  color: #dc2626;
  border-color: #fca5a5;
}

.btn-delete-card:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

/* Editar - Bordered con fondo claro */
.btn-edit-card {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-1);
  padding: var(--space-2) var(--space-3);
  background: #fef3c7;
  color: #d97706;
  border: 1px solid #fde68a;
  border-radius: var(--radius-md);
  font-size: var(--text-xs);
  font-weight: 600;
  cursor: pointer;
  transition:
    background var(--transition-fast),
    color var(--transition-fast),
    border-color var(--transition-fast);
  white-space: nowrap;
}

.btn-edit-card:hover:not(:disabled) {
  background: #fde68a;
  color: #d97706;
  border-color: #f59e0b;
}

.btn-edit-card:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

/* Aplicar - Bordered con fondo claro */
.btn-apply-card {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-1);
  padding: var(--space-2) var(--space-3);
  background: #e0f2fe;
  color: #0284c7;
  border: 1px solid #bae6fd;
  border-radius: var(--radius-md);
  font-size: var(--text-xs);
  font-weight: 600;
  cursor: pointer;
  transition:
    background var(--transition-fast),
    color var(--transition-fast),
    border-color var(--transition-fast);
  white-space: nowrap;
}

.btn-apply-card:hover:not(:disabled) {
  background: #bae6fd;
  color: #0284c7;
  border-color: #0284c7;
}

.btn-apply-card:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.record-card__actions {
  display: flex;
  align-items: center;
  gap: var(--space-1);
}

/* ── Vista card para móvil (< 700px) ─ */
.record-cards {
  display: none;
}

@container vaccine-card (max-width: 699px) {
  .table-wrapper {
    display: none;
  }

  .record-cards {
    display: flex;
    flex-direction: column;
    gap: 0;
  }

  .record-card {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
    padding: var(--space-4);
    border-bottom: 1px solid var(--color-border-light);
  }

  .record-card:last-child {
    border-bottom: none;
  }

  .record-card__top {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-2);
  }

  .record-card__vaccine {
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text-primary);
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .record-card__dates {
    display: flex;
    gap: var(--space-4);
  }

  .record-card__date-item {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .record-card__date-label {
    font-size: 0.65rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--color-text-tertiary);
  }

  .record-card__date-value {
    font-size: var(--text-xs);
    color: var(--color-text-secondary);
  }

  .record-card__note {
    font-size: var(--text-xs);
    color: var(--color-text-tertiary);
    margin: 0;
    line-height: 1.4;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .record-card__actions {
    display: flex;
    gap: var(--space-2);
  }
}

/* ── Estados de carga/error/vacío ─── */
.loading-state,
.error-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-10) var(--space-4);
  gap: var(--space-3);
  color: var(--color-text-tertiary);
  text-align: center;
}

.empty-state p {
  margin: 0;
  font-size: var(--text-sm);
}

.btn-add-empty {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  background: var(--color-accent-light);
  color: var(--color-accent-dark);
  border: 1px solid var(--color-accent);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
}

.btn-retry {
  padding: var(--space-2) var(--space-4);
  background: var(--color-accent);
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
}

.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* ── Footer con paginación ─────────── */
.table-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
  flex-wrap: wrap;
  padding: var(--space-4) var(--space-5);
  border-top: 1px solid var(--color-border-light);
  background: var(--color-bg);
}

@container vaccine-card (max-width: 480px) {
  .tab-header {
    padding: var(--space-3) var(--space-4);
  }

  .table-footer {
    flex-direction: column;
    align-items: stretch;
    padding: var(--space-3) var(--space-4);
    gap: var(--space-3);
  }
}
</style>
