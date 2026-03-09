<script setup lang="ts">
import { computed } from 'vue'
import { IconEdit, IconTrash } from '@tabler/icons-vue'
import type { Pet } from '@/types/pet'
import { calcAge, formatAge, isBirthdayToday } from '@/utils/pet'

defineEmits<{ edit: [pet: Pet]; delete: [id: string] }>()

const SPECIES_EMOJI: Record<string, string> = {
  dog: '🐕',
  cat: '🐈',
  bird: '🦜',
  rabbit: '🐇',
  fish: '🐠',
  other: '🐾',
}

const SPECIES_LABEL: Record<string, string> = {
  dog: 'Perro',
  cat: 'Gato',
  bird: 'Ave',
  rabbit: 'Conejo',
  fish: 'Pez',
  other: 'Otro',
}

function speciesEmoji(s: string) {
  return SPECIES_EMOJI[s] ?? '🐾'
}

function speciesLabel(s: string) {
  return SPECIES_LABEL[s] ?? s
}

const props = defineProps<{ pet: Pet; deleting?: boolean }>()

const ageText = computed(() => formatAge(calcAge(props.pet.birth_date)))
const isBirthday = computed(
  () => props.pet.birth_date_exact && isBirthdayToday(props.pet.birth_date),
)
</script>

<template>
  <article class="pet-card">
    <!-- Top accent strip per species -->
    <div class="card-strip" :class="`strip--${pet.species}`" />

    <!-- Card body -->
    <div class="card-body">
      <!-- Species avatar -->
      <div class="species-avatar" :class="`avatar--${pet.species}`">
        <span class="species-emoji">{{ speciesEmoji(pet.species) }}</span>
      </div>

        <div class="card-info">
          <div class="pet-name-row">
            <RouterLink :to="`/pets/${pet.id}`" class="pet-name">{{ pet.name }}</RouterLink>
            <span v-if="isBirthday" class="birthday-icon" title="Cumpleaños hoy">🎂</span>
          </div>

          <div class="badges">
            <span class="badge badge--species">{{ speciesLabel(pet.species) }}</span>
            <span v-if="pet.breed" class="badge badge--breed">{{ pet.breed }}</span>
          </div>

          <div class="card-meta">
            <span class="meta-item">
              <span class="meta-dot" />
              {{ ageText }}
            </span>
          </div>
        </div>
    </div>

    <!-- Actions -->
    <div class="card-actions">
      <button class="action-btn action-btn--edit" title="Editar" @click="$emit('edit', pet)">
        <IconEdit :size="15" :stroke-width="2" />
      </button>
      <button
        class="action-btn action-btn--delete"
        title="Eliminar"
        :disabled="deleting"
        @click="$emit('delete', pet.id)"
      >
        <span v-if="deleting" class="spinner-sm" />
        <IconTrash v-else :size="15" :stroke-width="2" />
      </button>
    </div>

    <!-- Invisible full-card link for keyboard / click-anywhere -->
    <RouterLink :to="`/pets/${pet.id}`" class="card-link" aria-hidden="true" tabindex="-1" />
  </article>
</template>

<style scoped>
.pet-card {
  position: relative;
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  overflow: hidden;
  display: flex;
  flex-direction: column;
  transition:
    transform var(--transition-fast),
    box-shadow var(--transition-fast),
    border-color var(--transition-fast);
}

.pet-card:hover {
  transform: translateY(-3px);
  box-shadow: var(--shadow-lg);
  border-color: var(--color-border);
}

/* ── Accent strip ─────────────────────── */
.card-strip {
  height: 4px;
  width: 100%;
  flex-shrink: 0;
}
.strip--dog    { background: linear-gradient(90deg, #f59e0b, #fbbf24); }
.strip--cat    { background: linear-gradient(90deg, #8b5cf6, #a78bfa); }
.strip--bird   { background: linear-gradient(90deg, #06b6d4, #22d3ee); }
.strip--rabbit { background: linear-gradient(90deg, #ec4899, #f472b6); }
.strip--fish   { background: linear-gradient(90deg, #3b82f6, #60a5fa); }
.strip--other  { background: linear-gradient(90deg, var(--color-accent), var(--color-accent-muted)); }

/* ── Body ─────────────────────────────── */
.card-body {
  display: flex;
  align-items: flex-start;
  gap: var(--space-4);
  padding: var(--space-5) var(--space-5) var(--space-3);
  flex: 1;
}

/* ── Species avatar ───────────────────── */
.species-avatar {
  width: 52px;
  height: 52px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.avatar--dog    { background: #fef3c7; }
.avatar--cat    { background: #ede9fe; }
.avatar--bird   { background: #cffafe; }
.avatar--rabbit { background: #fce7f3; }
.avatar--fish   { background: #dbeafe; }
.avatar--other  { background: var(--color-accent-light); }

.species-emoji { font-size: 1.75rem; line-height: 1; }

/* ── Info ─────────────────────────────── */
.card-info {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
  min-width: 0;
  flex: 1;
}

.pet-name {
  font-family: var(--font-display);
  font-size: var(--text-lg);
  font-weight: 600;
  color: var(--color-text-primary);
  text-decoration: none;
  line-height: 1.2;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  transition: color var(--transition-fast);
  position: relative;
  z-index: 1;
}
.pet-name:hover { color: var(--color-accent); }

.pet-name-row {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  min-width: 0;
}

.birthday-icon {
  font-size: 1rem;
  line-height: 1;
  flex-shrink: 0;
  position: relative;
  z-index: 1;
}

.badges {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-1);
}

.badge {
  display: inline-flex;
  align-items: center;
  padding: 2px var(--space-2);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: 600;
  letter-spacing: 0.02em;
}
.badge--species {
  background: var(--color-accent-light);
  color: var(--color-accent-dark);
}
.badge--breed {
  background: var(--color-bg-alt);
  color: var(--color-text-secondary);
}

.card-meta {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-2);
  align-items: center;
}
.meta-item {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
}
.meta-dot {
  width: 3px;
  height: 3px;
  border-radius: 50%;
  background: var(--color-text-tertiary);
  flex-shrink: 0;
}

/* ── Actions ──────────────────────────── */
.card-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-1);
  padding: var(--space-2) var(--space-3) var(--space-3);
  position: relative;
  z-index: 1;
}

.action-btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border: 1px solid transparent;
  border-radius: var(--radius-sm);
  background: transparent;
  cursor: pointer;
  transition:
    background var(--transition-fast),
    border-color var(--transition-fast),
    color var(--transition-fast);
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
}
.action-btn--delete:hover:not(:disabled) {
  background: var(--color-error-light);
  border-color: var(--color-error-border);
  color: var(--color-error);
}
.action-btn--delete:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

/* ── Full-card invisible link ─────────── */
.card-link {
  position: absolute;
  inset: 0;
  z-index: 0;
}

/* ── Spinner ──────────────────────────── */
.spinner-sm {
  display: inline-block;
  width: 13px;
  height: 13px;
  border: 2px solid var(--color-border);
  border-top-color: var(--color-error);
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
