<script setup lang="ts">
import PetAvatar from '@/components/pets/PetAvatar.vue'
import { getSpeciesLabel } from '@/constants/species'
import { formatDate } from '@/utils/date'
import { IconAlertTriangle } from '@tabler/icons-vue'
import { computed } from 'vue'

interface Props {
  petName: string
  petSpecies: string
  petId: string
  vaccineName: string
  nextDoseDate: string | null
  recordId: string
  urgency?: 'urgent' | 'soon' | 'future'
}

const props = withDefaults(defineProps<Props>(), {
  urgency: 'soon',
})

// Calcular días restantes
const daysUntilDue = computed(() => {
  if (!props.nextDoseDate) return 0
  const due = new Date(props.nextDoseDate)
  const now = new Date()
  due.setUTCHours(0, 0, 0, 0)
  now.setUTCHours(0, 0, 0, 0)
  const diffMs = due.getTime() - now.getTime()
  return Math.ceil(diffMs / (1000 * 60 * 60 * 24))
})

// Label de urgencia
const urgencyLabel = computed(() => {
  const days = daysUntilDue.value
  if (days < 0) {
    return 'Vencida'
  }
  return 'Programado'
})

// Badge class según estado
const badgeClass = computed(() => {
  const days = daysUntilDue.value
  if (days < 0) return 'status-badge--overdue'
  return 'status-badge--pending'
})

// Clase CSS según urgencia
const urgencyClass = computed(() => {
  const days = daysUntilDue.value
  if (days < 0) return 'upcoming-card--overdue'
  if (days <= 7) return 'upcoming-card--urgent'
  if (days <= 15) return 'upcoming-card--soon'
  return 'upcoming-card--future'
})

</script>

<template>
  <div class="upcoming-card-wrapper">
  <div class="upcoming-card" :class="urgencyClass">
    <div class="upcoming-card__content">
      <!-- Mascota -->
      <div class="pet-info">
        <PetAvatar :species="petSpecies" :name="petName" size="md" />
        <div class="pet-info__text">
          <span class="pet-info__name">{{ petName }}</span>
          <span class="pet-info__species">{{ getSpeciesLabel(petSpecies) }}</span>
        </div>
      </div>

      <!-- Vacuna -->
      <div class="vaccine-info">
        <span class="vaccine-info__name">{{ vaccineName }}</span>
      </div>

      <!-- Fecha de vencimiento y estado -->
      <div class="due-info">
        <span class="due-info__date">{{ nextDoseDate ? formatDate(nextDoseDate) : 'Sin programar' }}</span>
        <span class="status-badge" :class="badgeClass">
          <IconAlertTriangle v-if="daysUntilDue < 0" :size="12" :stroke-width="2.5" />
          <IconCheck v-else :size="12" :stroke-width="2.5" />
          {{ urgencyLabel }}
        </span>
      </div>
    </div>

  </div>
  </div>
</template>

<style scoped>
/* ── Wrapper: el container real ─────────────────────────── */
.upcoming-card-wrapper {
  container-type: inline-size;
  container-name: upcoming-card;
}

/* ── Card container ─────────────────────────────────────── */
.upcoming-card {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
  padding: var(--space-4);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  transition:
    box-shadow var(--transition-fast),
    transform var(--transition-fast);
}

.upcoming-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
}

/* Estados de urgencia */
.upcoming-card--overdue {
  border-left: 4px solid var(--color-error);
}

.upcoming-card--urgent {
  border-left: 4px solid var(--color-warning, #f59e0b);
}

.upcoming-card--soon {
  border-left: 4px solid var(--color-info, #3b82f6);
}

.upcoming-card--future {
  border-left: 4px solid var(--color-success, #10b981);
}

/* ── Content: columna por defecto (móvil) ─────────────── */
.upcoming-card__content {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

/* ── En containers anchos (≥520px): layout en fila ──────── */
@container upcoming-card (min-width: 520px) {
  .upcoming-card {
    flex-direction: row;
    align-items: center;
    gap: var(--space-4);
  }

  .upcoming-card__content {
    flex-direction: row;
    align-items: center;
    flex: 1;
    gap: var(--space-4);
    min-width: 0;
  }

  .pet-info {
    min-width: 140px;
    flex-shrink: 0;
  }

  .vaccine-info {
    flex: 1;
    min-width: 0;
  }

  .due-info {
    min-width: 180px;
    flex-shrink: 0;
  }
}

/* ── Pet info ───────────────────────────────────────────── */
.pet-info {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  min-width: 0;
}

.pet-info__text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.pet-info__name {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.pet-info__species {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  text-transform: capitalize;
}

/* ── Vaccine info ───────────────────────────────────────── */
.vaccine-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
  min-width: 0;
}

.vaccine-info__name {
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

/* ── Badges de estado ───────────────── */
.status-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
  padding: 3px var(--space-2);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: 600;
  white-space: nowrap;
  width: fit-content;
}

.status-badge--pending {
  background: #fef3c7;
  color: #92400e;
}

.status-badge--overdue {
  background: #fef2f2;
  color: #dc2626;
}

/* ── Due date info ───────────────── */
.due-info {
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.due-info__date {
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-text-primary);
}

</style>
