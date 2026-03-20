<script setup lang="ts">
import { ref, computed, type Ref } from 'vue'
import { useRoute } from 'vue-router'
import { IconPlus, IconPill, IconDroplet } from '@tabler/icons-vue'
import { useGetHealthRecordsByPetAndCategory } from '@/composables/useHealthRecords'
import { useGetPet } from '@/composables/usePets'
import { HealthCatalogCategory } from '@/constants/healthRecord'
import type { HealthRecord } from '@/types'
import HealthRecordFormModal from '@/components/health-tabs/HealthRecordFormModal.vue'

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
    <div class="content-card">
      <!-- Header -->
      <div class="tab-header">
        <h2 class="tab-title">
          <IconPill :size="20" :stroke-width="1.75" />
          Control de Parásitos
        </h2>
        <button class="btn-add" @click="openCreate">
          <IconPlus :size="16" :stroke-width="2.5" />
          Registrar aplicación
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
        <IconPill :size="40" :stroke-width="1.5" />
        <p>No hay desparasitaciones registradas</p>
        <button class="btn-add-empty" @click="openCreate">
          <IconPlus :size="16" :stroke-width="2.5" />
          Registrar primera desparasitación
        </button>
      </div>

      <!-- List -->
      <div v-else class="records-list">
        <div
          v-for="record in records"
          :key="record.id"
          class="record-card"
        >
          <!-- Icon: pill (internal) or droplet (external) -->
          <div class="record-icon" :class="isInternalDeworming(record.name) ? 'icon--internal' : 'icon--external'">
            <IconPill v-if="isInternalDeworming(record.name)" :size="18" :stroke-width="1.75" />
            <IconDroplet v-else :size="18" :stroke-width="1.75" />
          </div>

          <div class="record-content">
            <div class="record-name">
              {{ record.name }}
              <span class="record-type">
                {{ isInternalDeworming(record.name) ? '💊 Interno' : '💧 Externo' }}
              </span>
            </div>
            <div class="record-dates">
              <span class="record-date">
                <span class="record-date-label">Aplicado:</span>
                {{ formatDate(record.application_date) }}
              </span>
              <span class="record-date">
                <span class="record-date-label">Próxima:</span>
                {{ formatDate(record.next_dose_date) }}
              </span>
            </div>
          </div>

          <div class="record-actions">
            <button class="btn-edit-record" title="Editar" @click="openEdit(record)">
              ✏️ Editar
            </button>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal -->
    <HealthRecordFormModal
      v-if="showCreateModal"
      v-model="showCreateModal"
      :pet-id="petId"
      :pet-species="pet?.species ?? 'dog'"
      :category="HealthCatalogCategory.Deworming"
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

.content-card {
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

/* Quick action */
.quick-action {
  padding: 0 var(--space-5) var(--space-4);
}

.quick-action-card {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-4);
  background: linear-gradient(135deg, #f0f9ff 0%, #e0f2fe 100%);
  border: 1px solid #bae6fd;
  border-radius: var(--radius-lg);
}

.quick-action-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #dbeafe;
  color: #2563eb;
  border-radius: var(--radius-md);
  flex-shrink: 0;
}

.quick-action-content {
  flex: 1;
  min-width: 0;
}

.quick-action-title {
  display: block;
  font-size: var(--text-sm);
  font-weight: 600;
  color: #1e40af;
}

.quick-action-desc {
  font-size: var(--text-xs);
  color: #64748b;
}

.btn-quick-action {
  padding: var(--space-2) var(--space-3);
  background: #2563eb;
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-xs);
  font-weight: 600;
  cursor: pointer;
  white-space: nowrap;
  transition: background var(--transition-fast);
}

.btn-quick-action:hover {
  background: #1d4ed8;
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
  padding: var(--space-3) var(--space-4);
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  transition: box-shadow var(--transition-fast);
}

.record-card:hover {
  box-shadow: var(--shadow-sm);
}

.record-icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-md);
  flex-shrink: 0;
}

.record-icon.icon--internal {
  background: #fef3e2;
  color: #c4714a;
}

.record-icon.icon--external {
  background: #f0f9ff;
  color: #0284c7;
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
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.record-type {
  font-size: var(--text-xs);
  font-weight: 500;
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
  padding: 3px var(--space-2);
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

.record-actions {
  display: flex;
  gap: var(--space-1);
}

.btn-mark-applied {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #22c55e;
  color: #fff;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  transition: background var(--transition-fast);
}

.btn-mark-applied:hover {
  background: #16a34a;
}

.btn-edit-record {
  padding: var(--space-1) var(--space-2);
  background: transparent;
  border: none;
  cursor: pointer;
  opacity: 0.5;
  font-size: 14px;
  transition: opacity var(--transition-fast);
}

.btn-edit-record:hover {
  opacity: 1;
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

  .record-actions {
    order: 3;
    margin-left: auto;
  }
}
</style>
