<script setup lang="ts">
import { IconEdit, IconTrash, IconShieldCheck } from '@tabler/icons-vue'
import { getSpeciesLabel } from '@/constants/species'
import type { HealthCatalog } from '@/types/healthCatalog'

defineProps<{
  item: HealthCatalog
  deletingId: string | null
}>()

const emit = defineEmits<{
  edit: [item: HealthCatalog]
  delete: [item: HealthCatalog]
}>()

// Mapa de categorías removido ya que la columna no se muestra

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString('es-ES', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  })
}

function frequencyLabel(months: number | null): string {
  if (months === null) return 'Dosis única'
  if (months === 1) return 'Mensual'
  if (months === 12) return 'Anual'
  return `${months} meses`
}
</script>

<template>
  <tr class="catalog-row">
    <!-- Nombre + icono -->
    <td class="td-item">
      <div class="item-icon">
        <IconShieldCheck :size="18" :stroke-width="1.75" />
      </div>
      <div class="item-info">
        <span class="item-name">{{ item.name }}</span>
        <span v-if="item.description" class="item-desc">{{ item.description }}</span>
      </div>
    </td>



    <!-- Especies -->
    <td class="td-species">
      <div class="species-tags">
        <span v-for="specie in item.species" :key="specie" class="species-tag">
          {{ getSpeciesLabel(specie) }}
        </span>
      </div>
    </td>

    <!-- Frecuencia -->
    <td class="td-center">
      <span class="frequency-badge" :class="{ 'frequency-badge--single': item.frequency_months === null }">
        {{ frequencyLabel(item.frequency_months) }}
      </span>
    </td>

    <!-- Obligatoria -->
    <td class="td-center">
      <span class="mandatory-badge" :class="item.is_mandatory ? 'mandatory-badge--yes' : 'mandatory-badge--no'">
        {{ item.is_mandatory ? 'Sí' : 'No' }}
      </span>
    </td>

    <!-- Fecha creación -->
    <td class="td-center td-date">{{ formatDate(item.created_at) }}</td>

    <!-- Acciones -->
    <td class="td-center td-actions">
      <button class="action-btn action-btn--edit" title="Editar" @click="emit('edit', item)">
        <IconEdit :size="15" :stroke-width="2" />
      </button>
      <button
        class="action-btn action-btn--delete"
        title="Eliminar"
        :disabled="deletingId === item.id"
        @click="emit('delete', item)"
      >
        <IconTrash v-if="deletingId !== item.id" :size="15" :stroke-width="2" />
        <div v-else class="spinner spinner--sm" />
      </button>
    </td>
  </tr>
</template>

<style scoped>
.catalog-row {
  border-bottom: 1px solid var(--color-border-light);
  transition: background var(--transition-fast);
}

.catalog-row:last-child {
  border-bottom: none;
}

.catalog-row:hover {
  background: #f9fafb;
}

.catalog-row td {
  padding: var(--space-3) var(--space-4);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  vertical-align: middle;
}

.td-center {
  text-align: center;
}

.td-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  min-width: 200px;
}

.item-icon {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  background: var(--color-accent-light);
  color: var(--color-accent);
}

.item-info {
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
}

.item-name {
  font-weight: 600;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.item-desc {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  line-height: 1.3;
}

/* ── Category badge ─────────────────────────────────────── */
.category-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px var(--space-2);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: 600;
  white-space: nowrap;
}

.category-badge--vaccine {
  background: #eff6ff;
  color: #2563eb;
}

.category-badge--deworming {
  background: #fdf4ff;
  color: #9333ea;
}

.category-badge--exam {
  background: #fff7ed;
  color: #c2410c;
}

/* ── Species ────────────────────────────────────────────── */
.td-species {
  min-width: 180px;
  text-align: center;
}

.species-tags {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-1);
  justify-content: center;
}

.species-tag {
  display: inline-block;
  background: #f0fdf4;
  color: #15803d;
  padding: 2px var(--space-2);
  border-radius: var(--radius-sm);
  font-size: var(--text-xs);
  font-weight: 500;
  white-space: nowrap;
}

/* ── Frequency ──────────────────────────────────────────── */
.frequency-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px var(--space-2);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: 600;
  background: #eff6ff;
  color: #2563eb;
}

.frequency-badge--single {
  background: #f3f4f6;
  color: #6b7280;
}

/* ── Mandatory ──────────────────────────────────────────── */
.mandatory-badge {
  display: inline-flex;
  align-items: center;
  padding: 2px var(--space-2);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: 600;
}

.mandatory-badge--yes {
  background: #fef3c7;
  color: #b45309;
}

.mandatory-badge--no {
  background: #f3f4f6;
  color: #6b7280;
}

/* ── Date ───────────────────────────────────────────────── */
.td-date {
  white-space: nowrap;
  font-size: var(--text-xs) !important;
  color: var(--color-text-tertiary) !important;
}

/* ── Actions ────────────────────────────────────────────── */
.td-actions {
  white-space: nowrap;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border: 1px solid transparent;
  border-radius: var(--radius-md);
  cursor: pointer;
  background: transparent;
  transition: background var(--transition-fast), border-color var(--transition-fast), color var(--transition-fast);
}

.action-btn--edit {
  color: var(--color-text-tertiary);
}

.action-btn--edit:hover {
  background: #eff6ff;
  border-color: #bfdbfe;
  color: #2563eb;
}

.action-btn--delete {
  color: var(--color-text-tertiary);
  margin-left: var(--space-1);
}

.action-btn--delete:hover:not(:disabled) {
  background: #fef2f2;
  border-color: #fecaca;
  color: var(--color-error);
}

.action-btn--delete:disabled {
  opacity: 0.3;
  cursor: not-allowed;
}

/* ── Spinner ────────────────────────────────────────────── */
.spinner {
  width: 28px;
  height: 28px;
  border: 3px solid var(--color-border-light);
  border-top-color: var(--color-accent);
  border-radius: var(--radius-full);
  animation: spin 0.7s linear infinite;
}

.spinner--sm {
  width: 16px;
  height: 16px;
  border-width: 2px;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
