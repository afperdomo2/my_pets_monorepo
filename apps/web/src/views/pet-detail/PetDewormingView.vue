<script setup lang="ts">
import { ref, computed, type Ref } from 'vue'
import { useRoute } from 'vue-router'
import {
  IconPlus,
  IconRefresh,
  IconPill,
  IconDroplet,
  IconTrash,
  IconEdit,
  IconVaccine,
} from '@tabler/icons-vue'
import { useGetHealthRecordsByPetAndCategory, useDeleteHealthRecord, useUpdateHealthRecord } from '@/composables/useHealthRecords'
import { useCreateVaccineApplication } from '@/composables/useVaccineApplications'
import { useGetPet } from '@/composables/usePets'
import { HealthCatalogCategory } from '@/constants/healthRecord'
import type { HealthRecord } from '@/types'
import AppPagination from '@/components/ui/AppPagination.vue'
import PerPageSelector from '@/components/ui/PerPageSelector.vue'
import HealthRecordCreateModal from '@/components/health-tabs/HealthRecordCreateModal.vue'
import ConfirmDeleteModal from '@/components/health-tabs/ConfirmDeleteModal.vue'
import HealthRecordEditModal from '@/components/health-tabs/HealthRecordEditModal.vue'
import VaccineApplicationModal from '@/components/health-tabs/VaccineApplicationModal.vue'

const route = useRoute()
const petId = computed(() => String(route.params.id))

const page = ref(1)
const perPage = ref(10)
const category = ref(HealthCatalogCategory.Deworming)

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

function formatDate(dateStr: string | null): string {
  if (!dateStr) return '—'
  const date = new Date(dateStr)
  return date.toLocaleDateString('es-ES', { day: '2-digit', month: 'short', year: 'numeric' })
}

function isInternalDeworming(name: string): boolean {
  const lower = name.toLowerCase()
  return lower.includes('pastilla') || lower.includes('oral') || lower.includes('giardia') || lower.includes('nematodos') || lower.includes('céstodos')
}

const showCreateModal = ref(false)
const showConfirmModal = ref(false)
const showEditModal = ref(false)
const showApplicationModal = ref(false)
const recordToDelete = ref<HealthRecord | null>(null)
const recordToEdit = ref<HealthRecord | null>(null)
const healthRecordToApply = ref<HealthRecord | null>(null)
const deletingId = ref<string | null>(null)

const deleteRecord = useDeleteHealthRecord()
const updateRecord = useUpdateHealthRecord()
const createApplication = useCreateVaccineApplication()

function openCreate() {
  showCreateModal.value = true
}

function openEdit(record: HealthRecord) {
  recordToEdit.value = record
  showEditModal.value = true
}

function openApplicationModal(record: HealthRecord) {
  healthRecordToApply.value = record
  showApplicationModal.value = true
}

function openDeleteConfirm(record: HealthRecord) {
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
            <IconPill :size="20" :stroke-width="1.75" />
            Control de Parásitos
          </h2>
          <span class="panel-count">{{ total }} registros</span>
        </div>
        <div class="tab-header-right">
          <button class="btn-refresh" :disabled="isLoading" title="Refrescar" @click="refresh">
            <IconRefresh :size="16" :stroke-width="2" :class="{ 'spin': isLoading }" />
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
        <table v-if="!isLoading && !isError && records.length > 0" class="deworming-table">
          <thead>
            <tr>
              <th>Parásito</th>
              <th class="th-center">Tipo</th>
              <th class="th-center">Fecha aplicación</th>
              <th class="th-center">Próxima dosis</th>
              <th class="th-center">Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="record in records" :key="record.id" class="deworming-row">
              <td>
                <div class="deworming-cell">
                  <span class="deworming-name">{{ record.name }}</span>
                  <span
                    v-if="record.notes"
                    class="deworming-note"
                    :title="record.notes"
                  >{{ record.notes }}</span>
                </div>
              </td>
              <td class="td-center">
                <span class="deworming-type" :class="isInternalDeworming(record.name) ? 'deworming-type--internal' : 'deworming-type--external'">
                  <IconPill v-if="isInternalDeworming(record.name)" :size="14" :stroke-width="2" />
                  <IconDroplet v-else :size="14" :stroke-width="2" />
                  {{ isInternalDeworming(record.name) ? 'Interno' : 'Externo' }}
                </span>
              </td>
              <td class="td-center">
                <span class="date-cell">{{ formatDate(record.application_date) }}</span>
              </td>
              <td class="td-center">
                <span class="date-cell">
                  {{ formatDate(record.next_dose_date) }}
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
                  </button>
                  <button
                    class="btn-action btn-edit"
                    title="Editar"
                    :disabled="updateRecord.isPending.value"
                    @click="openEdit(record)"
                  >
                    <IconEdit :size="14" :stroke-width="2" />
                  </button>
                  <button
                    class="btn-action btn-delete"
                    title="Eliminar"
                    :disabled="deletingId === record.id || deleteRecord.isPending.value"
                    @click="openDeleteConfirm(record)"
                  >
                    <IconTrash :size="14" :stroke-width="2" />
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
          <IconPill :size="40" :stroke-width="1.5" />
          <p>No hay desparasitaciones registradas</p>
          <button class="btn-add-empty" @click="openCreate">
            <IconPlus :size="16" :stroke-width="2.5" />
            Registrar primera desparasitación
          </button>
        </div>
      </div>

      <!-- Vista card para móvil -->
      <div v-if="!isLoading && !isError && records.length > 0" class="record-cards">
        <div v-for="record in records" :key="`card-${record.id}`" class="record-card">
          <div class="record-card__top">
            <span class="record-card__name">{{ record.name }}</span>
            <span class="record-card__type" :class="isInternalDeworming(record.name) ? 'record-card__type--internal' : 'record-card__type--external'">
              <IconPill v-if="isInternalDeworming(record.name)" :size="12" :stroke-width="2" />
              <IconDroplet v-else :size="12" :stroke-width="2" />
              {{ isInternalDeworming(record.name) ? 'Interno' : 'Externo' }}
            </span>
            <div class="record-card__actions">
              <button
                class="btn-apply-card"
                title="Aplicar dosis"
                :disabled="createApplication.isPending.value"
                @click="openApplicationModal(record)"
              >
                <IconVaccine :size="14" :stroke-width="2" />
              </button>
              <button
                class="btn-edit-card"
                title="Editar"
                :disabled="updateRecord.isPending.value"
                @click="openEdit(record)"
              >
                <IconEdit :size="14" :stroke-width="2" />
              </button>
              <button
                class="btn-delete-card"
                title="Eliminar"
                :disabled="deletingId === record.id || deleteRecord.isPending.value"
                @click="openDeleteConfirm(record)"
              >
                <IconTrash :size="14" :stroke-width="2" />
              </button>
            </div>
          </div>
          <div class="record-card__dates">
            <div class="record-card__date-item">
              <span class="record-card__date-label">Aplicada</span>
              <span class="record-card__date-value">{{ formatDate(record.application_date) }}</span>
            </div>
            <div class="record-card__date-item">
              <span class="record-card__date-label">Próxima</span>
              <span class="record-card__date-value">
                {{ formatDate(record.next_dose_date) }}
              </span>
            </div>
          </div>
          <p v-if="record.notes" class="record-card__note" :title="record.notes">
            {{ record.notes }}
          </p>
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

    <!-- Modal -->
    <HealthRecordCreateModal
      v-if="showCreateModal"
      :pet-id="petId"
      :pet-species="pet?.species ?? 'dog'"
      category="deworming"
      @close="showCreateModal = false"
    />

    <!-- Modal confirmar eliminación -->
    <ConfirmDeleteModal
      v-model="showConfirmModal"
      :record-name="recordToDelete?.name ?? ''"
      record-type="deworming"
      :deleting="deleteRecord.isPending.value"
      @confirm="handleDeleteConfirm"
    />

    <!-- Modal editar registro -->
    <HealthRecordEditModal
      v-if="showEditModal && recordToEdit"
      :record="recordToEdit"
      category="deworming"
      @close="showEditModal = false; recordToEdit = null"
      @updated="refresh"
    />

    <!-- Modal aplicar dosis -->
    <VaccineApplicationModal
      v-if="showApplicationModal && healthRecordToApply"
      :health-record-id="healthRecordToApply.id"
      category="deworming"
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
  container-name: deworming-card;
  background: transparent;
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-xl);
  overflow: hidden;
  box-shadow: none;
}

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

.tab-header-right {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.panel-count {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  font-weight: 500;
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
  transition: background var(--transition-fast), color var(--transition-fast);
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

/* States */
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
  to { transform: rotate(360deg); }
}

/* ── Tabla ────────────────────────── */
.table-wrapper {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.deworming-table {
  width: 100%;
  min-width: 700px;
  border-collapse: collapse;
}

.deworming-table thead tr {
  background: var(--color-bg);
  border-bottom: 1px solid var(--color-border-light);
}

.deworming-table th {
  padding: var(--space-3) var(--space-4);
  text-align: left;
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.deworming-table th.th-center {
  text-align: center;
}

.deworming-row {
  transition: background var(--transition-fast);
}

.deworming-row:hover {
  background: var(--color-bg-alt);
}

.deworming-row td {
  padding: var(--space-3) var(--space-4);
  border-bottom: 1px solid var(--color-border-light);
}

.deworming-row:last-child td {
  border-bottom: none;
}

.td-center {
  text-align: center;
}

/* ── Celdas ───────────────────────── */
.deworming-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.deworming-name {
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  font-weight: 500;
}

.deworming-note {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.4;
}

.deworming-type {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 4px;
  padding: 4px var(--space-2);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: 600;
  white-space: nowrap;
}

.deworming-type--internal {
  background: #fef3e2;
  color: #c4714a;
}

.deworming-type--external {
  background: #f0f9ff;
  color: #0284c7;
}

.date-cell {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  white-space: nowrap;
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
  gap: 4px;
  padding: 4px var(--space-3);
  border-radius: var(--radius-md);
  font-size: var(--text-xs);
  font-weight: 600;
  cursor: pointer;
  border: 1px solid transparent;
  transition: background var(--transition-fast), color var(--transition-fast);
  white-space: nowrap;
}

.btn-edit {
  background: var(--color-bg-alt);
  color: var(--color-text-secondary);
  border-color: var(--color-border);
}

.btn-edit:hover {
  background: var(--color-accent-light);
  color: var(--color-accent-dark);
  border-color: var(--color-accent);
}

.btn-delete {
  background: var(--color-bg-alt);
  color: var(--color-text-secondary);
  border-color: var(--color-border);
}

.btn-delete:hover {
  background: #fef2f2;
  color: #dc2626;
  border-color: #fecaca;
}

.btn-delete-card {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 6px;
  background: transparent;
  color: var(--color-text-tertiary);
  border: 1px solid transparent;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.btn-delete-card:hover:not(:disabled) {
  background: #fef2f2;
  color: #dc2626;
}

.btn-delete-card:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.btn-edit-card {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 6px;
  background: transparent;
  color: var(--color-text-tertiary);
  border: 1px solid transparent;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.btn-edit-card:hover:not(:disabled) {
  background: var(--color-accent-light);
  color: var(--color-accent-dark);
}

.btn-edit-card:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.btn-apply-card {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 6px;
  background: transparent;
  color: var(--color-text-tertiary);
  border: 1px solid transparent;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.btn-apply-card:hover:not(:disabled) {
  background: #e0f2fe;
  color: #0284c7;
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

.btn-apply {
  background: var(--color-bg-alt);
  color: var(--color-text-secondary);
  border-color: var(--color-border);
}

.btn-apply:hover {
  background: #e0f2fe;
  color: #0284c7;
  border-color: #bae6fd;
}

/* ── Vista card para móvil (< 700px) ─ */
.record-cards {
  display: none;
}

@container deworming-card (max-width: 699px) {
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

  .record-card__name {
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text-primary);
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .record-card__type {
    display: inline-flex;
    align-items: center;
    gap: 3px;
    padding: 3px var(--space-2);
    border-radius: var(--radius-full);
    font-size: 0.6rem;
    font-weight: 600;
    white-space: nowrap;
    flex-shrink: 0;
  }

  .record-card__type--internal {
    background: #fef3e2;
    color: #c4714a;
  }

  .record-card__type--external {
    background: #f0f9ff;
    color: #0284c7;
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

@container deworming-card (max-width: 480px) {
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
