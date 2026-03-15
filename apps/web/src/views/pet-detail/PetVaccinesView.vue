<script setup lang="ts">
import { ref, computed, type Ref } from 'vue'
import { useRoute } from 'vue-router'
import { IconPlus, IconVaccine, IconCheck, IconClock, IconAlertTriangle } from '@tabler/icons-vue'
import { useGetHealthRecordsByPetAndCategory } from '@/composables/useHealthRecords'
import { useGetPet } from '@/composables/usePets'
import { HealthRecordStatus, HealthCatalogCategory } from '@/constants/healthRecord'
import type { HealthRecord } from '@/types'
import AppPagination from '@/components/ui/AppPagination.vue'
import HealthRecordFormModal from '@/components/health-tabs/HealthRecordFormModal.vue'

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

const STATUS_CONFIG = {
  [HealthRecordStatus.Applied]: { label: 'Aplicada', className: 'status--uptodate', icon: IconCheck },
  [HealthRecordStatus.Pending]: { label: 'Pendiente', className: 'status--upcoming', icon: IconClock },
  [HealthRecordStatus.Overdue]: { label: 'Vencida', className: 'status--overdue', icon: IconAlertTriangle },
}

function formatDate(dateStr: string | null): string {
  if (!dateStr) return '—'
  const date = new Date(dateStr)
  return date.toLocaleDateString('es-ES', { day: '2-digit', month: 'short', year: 'numeric' })
}

const showCreateModal = ref(false)
const editingRecord = ref<HealthRecord | null>(null)

function openCreate() {
  editingRecord.value = null
  showCreateModal.value = true
}

function openEdit(record: HealthRecord) {
  editingRecord.value = record
  showCreateModal.value = true
}
</script>

<template>
  <div class="tab-view">
    <!-- Header -->
      <div class="tab-header">
        <h2 class="tab-title">
          <IconVaccine :size="20" :stroke-width="1.75" />
          Carnet de Vacunación
        </h2>
        <button class="btn-add" @click="openCreate">
          <IconPlus :size="16" :stroke-width="2.5" />
          Registrar vacuna
        </button>
      </div>

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
          Registrar primera vacuna
        </button>
      </div>

      <!-- List -->
      <div v-else class="records-list">
        <div
          v-for="record in records"
          :key="record.id"
          class="record-card"
        >
          <div class="record-icon">
            <IconVaccine :size="18" :stroke-width="1.75" />
          </div>
          <div class="record-content">
            <div class="record-name">{{ record.name }}</div>
            <div class="record-dates">
              <span class="record-date">
                <span class="record-date-label">Aplicada:</span>
                {{ formatDate(record.application_date) }}
              </span>
              <span class="record-date">
                <span class="record-date-label">Próxima:</span>
                {{ formatDate(record.due_date) }}
              </span>
            </div>
          </div>
          <div class="record-status">
            <span
              class="status-badge"
              :class="STATUS_CONFIG[record.status as keyof typeof STATUS_CONFIG]?.className"
            >
              <component
                :is="STATUS_CONFIG[record.status as keyof typeof STATUS_CONFIG]?.icon"
                :size="12"
                :stroke-width="2.5"
              />
              {{ STATUS_CONFIG[record.status as keyof typeof STATUS_CONFIG]?.label }}
            </span>
          </div>
          <button class="btn-edit-record" title="Editar" @click="openEdit(record)">
            ✏️
          </button>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="pagination-bar">
        <AppPagination
          :current-page="page"
          :total-pages="totalPages"
          :total-items="total"
          :per-page="perPage"
          @update:page="page = $event"
        />
      </div>

    <!-- Modal -->
    <HealthRecordFormModal
      v-if="showCreateModal"
      v-model="showCreateModal"
      :pet-id="petId"
      :pet-species="pet?.species ?? 'dog'"
      :category="HealthCatalogCategory.Vaccine"
      :editing-record="editingRecord"
      @close="showCreateModal = false"
    />
  </div>
</template>

<style scoped>
.tab-view {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
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

/* Records list */
.records-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.record-card {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-4);
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  transition: box-shadow var(--transition-fast);
}

.record-card:hover {
  box-shadow: var(--shadow-sm);
}

.record-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--color-accent-light);
  color: var(--color-accent);
  border-radius: var(--radius-md);
  flex-shrink: 0;
}

.record-content {
  flex: 1;
  min-width: 0;
}

.record-name {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
  margin-bottom: 2px;
}

.record-dates {
  display: flex;
  gap: var(--space-4);
  flex-wrap: wrap;
}

.record-date {
  font-size: var(--text-xs);
  color: var(--color-text-secondary);
}

.record-date-label {
  color: var(--color-text-tertiary);
  margin-right: 2px;
}

.record-status {
  flex-shrink: 0;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px var(--space-3);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: 600;
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

.btn-edit-record {
  padding: var(--space-2);
  background: transparent;
  border: none;
  cursor: pointer;
  opacity: 0.5;
  transition: opacity var(--transition-fast);
}

.btn-edit-record:hover {
  opacity: 1;
}

.pagination-bar {
  display: flex;
  justify-content: center;
  padding: var(--space-4) var(--space-5);
  border-top: 1px solid var(--color-border-light);
}

@media (max-width: 600px) {
  .record-card {
    flex-wrap: wrap;
  }

  .record-content {
    order: 1;
    width: 100%;
  }

  .record-status {
    order: 2;
  }

  .btn-edit-record {
    order: 3;
    margin-left: auto;
  }
}
</style>
