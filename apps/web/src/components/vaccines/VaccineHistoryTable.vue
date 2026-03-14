<script setup lang="ts">
import { ref, computed } from 'vue'
import { useGetAllHealthRecords } from '@/composables/useHealthRecords'
import type { HealthRecordStatusType } from '@/constants/healthRecord'
import { HealthRecordStatus } from '@/constants/healthRecord'
import PetAvatar from '@/components/pets/PetAvatar.vue'
import AppPagination from '@/components/ui/AppPagination.vue'
import PerPageSelector from '@/components/ui/PerPageSelector.vue'

const page = ref(1)
const perPage = ref(10)

const { data, isLoading, isError, refresh } = useGetAllHealthRecords(page, perPage)

const records = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.total ?? 0)
const totalPages = computed(() => data.value?.total_pages ?? 0)

const STATUS_CONFIG: Record<HealthRecordStatusType, { label: string; className: string }> = {
  [HealthRecordStatus.Applied]: { label: 'Aplicada', className: 'status--uptodate' },
  [HealthRecordStatus.Pending]: { label: 'Pendiente', className: 'status--upcoming' },
  [HealthRecordStatus.Overdue]: { label: 'Vencida', className: 'status--overdue' },
}

function formatDate(dateStr: string | null): string {
  if (!dateStr) return '—'
  const date = new Date(dateStr)
  return date.toLocaleDateString('es-ES', { day: '2-digit', month: 'short', year: 'numeric' })
}
</script>

<template>
  <div class="history-panel">
    <div class="panel-header">
      <h2>Historial de vacunación</h2>
      <span class="panel-count">{{ total }} registros</span>
    </div>

    <div class="table-wrapper">
      <table v-if="!isLoading && !isError" class="history-table">
        <thead>
          <tr>
            <th>Mascota</th>
            <th>Vacuna</th>
            <th class="th-center">Fecha aplicación</th>
            <th class="th-center">Próxima dosis</th>
            <th class="th-center">Estado</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="record in records" :key="record.id" class="history-row">
            <td>
              <div class="pet-cell">
                <PetAvatar :species="record.pet.species" :name="record.pet.name" size="sm" />
                <span class="pet-cell-name">{{ record.pet.name }}</span>
              </div>
            </td>
            <td>
              <span class="vaccine-cell">{{ record.name }}</span>
            </td>
            <td class="td-center">
              <span class="date-cell">{{ formatDate(record.application_date) }}</span>
            </td>
            <td class="td-center">
              <span class="date-cell">{{ formatDate(record.due_date) }}</span>
            </td>
            <td class="td-center">
              <span class="status-badge" :class="STATUS_CONFIG[record.status as HealthRecordStatusType]?.className">
                {{ STATUS_CONFIG[record.status as HealthRecordStatusType]?.label }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>

      <div v-else-if="isLoading" class="loading-state">
        <div class="spinner" />
        <p>Cargando registros...</p>
      </div>

      <div v-else-if="isError" class="error-state">
        <p>Error al cargar los registros</p>
        <button class="btn-retry" @click="refresh">Reintentar</button>
      </div>
    </div>

    <!-- Paginación -->
    <div class="table-footer">
      <PerPageSelector v-model="perPage" :options="[10, 25, 50]" />
      <AppPagination
        v-model:current-page="page"
        :total-pages="totalPages"
        :total-items="total"
        :per-page="perPage"
      />
    </div>
  </div>
</template>

<style scoped>
/* ── Panel ────────────────────────── */
.history-panel {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-5) var(--space-6);
  border-bottom: 1px solid var(--color-border-light);
}

.panel-header h2 {
  font-family: var(--font-display);
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-text-primary);
}

.panel-count {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  font-weight: 500;
}

/* ── Tabla ────────────────────────── */
.table-wrapper {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.history-table {
  width: 100%;
  min-width: 620px;
  border-collapse: collapse;
}

.history-table thead tr {
  background: var(--color-bg);
  border-bottom: 1px solid var(--color-border-light);
}

.history-table th {
  padding: var(--space-3) var(--space-4);
  text-align: left;
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.history-table th.th-center {
  text-align: center;
}

.history-row {
  transition: background var(--transition-fast);
}

.history-row:hover {
  background: var(--color-bg-alt);
}

.history-row td {
  padding: var(--space-3) var(--space-4);
  border-bottom: 1px solid var(--color-border-light);
}

.history-row:last-child td {
  border-bottom: none;
}

.td-center {
  text-align: center;
}

/* ── Celdas ──────────────────────── */
.pet-cell {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.pet-cell-name {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
}

.vaccine-cell {
  font-size: var(--text-sm);
  color: var(--color-text-primary);
}

.date-cell {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

/* ── Badges de estado ────────────── */
.status-badge {
  display: inline-flex;
  align-items: center;
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

/* ── Estados de carga/error ──────── */
.loading-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-10) var(--space-4);
  gap: var(--space-3);
}

.loading-state p,
.error-state p {
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
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

.btn-retry {
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

.btn-retry:hover {
  background: var(--color-accent-hover);
}

/* ── Footer con paginación ───────── */
.table-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
  flex-wrap: wrap;
  padding: var(--space-4) var(--space-6);
  border-top: 1px solid var(--color-border-light);
  background: var(--color-bg);
}

@media (max-width: 600px) {
  .table-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-3);
  }
}
</style>
