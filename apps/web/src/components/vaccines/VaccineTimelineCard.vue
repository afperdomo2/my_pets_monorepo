<script setup lang="ts">
defineProps<{
  petName: string
  petInitials: string
  petSpecies: string
  lifeStage: string
  vaccineName: string
  dueDate: string
  dueLabel: string
  urgency: 'urgent' | 'soon' | 'future'
}>()

defineEmits<{
  register: []
}>()

const URGENCY_LABELS: Record<string, string> = {
  urgent: 'Vencida',
  soon: 'Próxima',
  future: 'Programada',
}
</script>

<template>
  <div class="timeline-card-wrapper">
  <div class="timeline-card" :class="`timeline-card--${urgency}`">
    <div class="timeline-card__strip" />
    <div class="timeline-card__content">
      <!-- Mascota -->
      <div class="timeline-card__pet">
        <div class="pet-avatar" :class="`pet-avatar--${petSpecies}`">
          {{ petInitials }}
        </div>
        <div class="pet-details">
          <div class="pet-name-row">
            <span class="pet-name">{{ petName }}</span>
            <span class="life-stage-badge">{{ lifeStage }}</span>
          </div>
          <span class="vaccine-name">{{ vaccineName }}</span>
        </div>
      </div>

      <!-- Fecha & estado -->
      <div class="timeline-card__meta">
        <div class="due-info">
          <span class="due-date">{{ dueDate }}</span>
          <span class="due-label" :class="`due-label--${urgency}`">
            {{ dueLabel }}
          </span>
        </div>
        <span class="urgency-tag" :class="`urgency-tag--${urgency}`">
          {{ URGENCY_LABELS[urgency] }}
        </span>
      </div>

      <!-- Acción -->
      <button class="btn-register" @click="$emit('register')">
        Registrar aplicación
      </button>
    </div>
  </div>
  </div>
</template>

<style scoped>
/* ── Wrapper: container real ─────────────────────────────── */
.timeline-card-wrapper {
  container-type: inline-size;
  container-name: timeline-card;
}

.timeline-card {
  display: flex;
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  overflow: hidden;
  transition: box-shadow var(--transition-base), transform var(--transition-base);
  cursor: pointer;
}

.timeline-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-1px);
}

/* ── Borde lateral de urgencia ────── */
.timeline-card__strip {
  width: 5px;
  flex-shrink: 0;
}

.timeline-card--urgent .timeline-card__strip {
  background: linear-gradient(180deg, #DC2626, #EF4444);
}
.timeline-card--soon .timeline-card__strip {
  background: linear-gradient(180deg, #D97706, #F59E0B);
}
.timeline-card--future .timeline-card__strip {
  background: linear-gradient(180deg, var(--color-accent), var(--color-accent-muted));
}

/* ── Contenido ───────────────────── */
.timeline-card__content {
  flex: 1;
  padding: var(--space-4) var(--space-5);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

@container timeline-card (min-width: 40em) {
  .timeline-card__content {
    flex-direction: row;
    align-items: center;
    gap: var(--space-4);
  }
}

/* ── Mascota ─────────────────────── */
.timeline-card__pet {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  flex: 1;
  min-width: 0;
}

.pet-avatar {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-sm);
  font-weight: 700;
  letter-spacing: 0.02em;
  flex-shrink: 0;
}

.pet-avatar--dog    { background: #fef3c7; color: #92400e; }
.pet-avatar--cat    { background: #ede9fe; color: #5b21b6; }
.pet-avatar--bird   { background: #cffafe; color: #0e7490; }
.pet-avatar--rabbit { background: #fce7f3; color: #9d174d; }
.pet-avatar--fish   { background: #dbeafe; color: #1e40af; }
.pet-avatar--other  { background: var(--color-accent-light); color: var(--color-accent-dark); }

.pet-details {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.pet-name-row {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.pet-name {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
}

.life-stage-badge {
  font-size: 0.65rem;
  font-weight: 600;
  padding: 1px var(--space-2);
  border-radius: var(--radius-full);
  background: var(--color-bg-alt);
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border-light);
  white-space: nowrap;
}

.vaccine-name {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
}

/* ── Meta ────────────────────────── */
.timeline-card__meta {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

@container timeline-card (min-width: 40em) {
  .timeline-card__meta {
    flex-direction: column;
    align-items: flex-end;
    gap: var(--space-1);
  }
}

.due-info {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

@container timeline-card (min-width: 40em) {
  .due-info {
    text-align: right;
  }
}

.due-date {
  font-size: var(--text-xs);
  color: var(--color-text-secondary);
  font-weight: 500;
}

.due-label {
  font-size: var(--text-xs);
  font-weight: 600;
}

.due-label--urgent { color: #DC2626; }
.due-label--soon   { color: #D97706; }
.due-label--future { color: var(--color-accent); }

.urgency-tag {
  font-size: 0.65rem;
  font-weight: 600;
  padding: 2px var(--space-2);
  border-radius: var(--radius-full);
  text-transform: uppercase;
  letter-spacing: 0.06em;
  white-space: nowrap;
}

.urgency-tag--urgent { background: #FEF2F2; color: #DC2626; }
.urgency-tag--soon   { background: #FEF3E2; color: #D97706; }
.urgency-tag--future { background: #E8F5EE; color: #2E7D52; }

/* ── Botón ───────────────────────── */
.btn-register {
  padding: var(--space-2) var(--space-4);
  min-height: 44px;
  background: var(--color-surface);
  color: var(--color-accent);
  border: 1.5px solid var(--color-accent-muted);
  border-radius: var(--radius-md);
  font-size: var(--text-xs);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast),
    border-color var(--transition-fast);
  white-space: nowrap;
  align-self: flex-start;
}

@container timeline-card (min-width: 40em) {
  .btn-register {
    align-self: center;
    flex-shrink: 0;
  }
}

.btn-register:hover {
  background: var(--color-accent);
  color: #fff;
  border-color: var(--color-accent);
}
</style>
