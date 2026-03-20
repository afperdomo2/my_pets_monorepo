<script setup lang="ts">
import { computed, ref } from 'vue'
import { useRoute, RouterLink } from 'vue-router'
import { IconArrowLeft, IconEdit, IconTrash, IconPaw, IconVaccine, IconPill, IconStethoscope } from '@tabler/icons-vue'
import { useGetPet, useDeletePet } from '@/composables/usePets'
import PetFormModal from '@/components/pets/PetFormModal.vue'
import ConfirmDeleteModal from '@/components/pets/ConfirmDeleteModal.vue'
import { calcAge, formatAge, formatBirthDate, formatWeight, isBirthdayToday, lifeStageLabel } from '@/utils/pet'

const route = useRoute()
const id = String(route.params.id)
const showEditModal = ref(false)
const showConfirmDelete = ref(false)

const { data: pet, isLoading, isError } = useGetPet(id)
const deletePet = useDeletePet()

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

function speciesEmoji(s: string) { return SPECIES_EMOJI[s] ?? '🐾' }
function speciesLabel(s: string) { return SPECIES_LABEL[s] ?? s }

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString('es-ES', {
    day: '2-digit', month: 'long', year: 'numeric',
  })
}

const ageText = computed(() => pet.value ? formatAge(calcAge(pet.value.birth_date)) : '')
const isBirthday = computed(
  () => !!pet.value && pet.value.birth_date_exact && isBirthdayToday(pet.value.birth_date),
)

const LIFE_STAGE_COLORS: Record<string, { bg: string; text: string }> = {
  puppy:        { bg: '#fef9c3', text: '#854d0e' },
  junior:       { bg: '#dcfce7', text: '#166534' },
  adult:        { bg: '#dbeafe', text: '#1e40af' },
  senior:       { bg: '#ede9fe', text: '#5b21b6' },
  geriatric:    { bg: '#fee2e2', text: '#991b1b' },
  kitten:       { bg: '#fef9c3', text: '#854d0e' },
  young_adult:  { bg: '#dcfce7', text: '#166534' },
  mature_adult: { bg: '#dbeafe', text: '#1e40af' },
  end_of_life:  { bg: '#fee2e2', text: '#991b1b' },
  infant:       { bg: '#fef9c3', text: '#854d0e' },
  juvenile:     { bg: '#dcfce7', text: '#166534' },
  teenager:     { bg: '#d1fae5', text: '#065f46' },
}

function lifeStageStyle(stage: string) {
  return LIFE_STAGE_COLORS[stage] ?? { bg: 'var(--color-bg-alt)', text: 'var(--color-text-secondary)' }
}

function openDeleteConfirm() {
  showConfirmDelete.value = true
}

async function handleDelete() {
  await deletePet.mutateAsync(id)
  showConfirmDelete.value = false
  window.location.href = '/pets'
}

const TABS = [
  { label: 'Vacunas', name: 'pet-detail-vaccines', icon: IconVaccine },
  { label: 'Desparasitación', name: 'pet-detail-deworming', icon: IconPill },
  { label: 'Exámenes', name: 'pet-detail-exams', icon: IconStethoscope },
]
</script>

<template>
  <div class="pet-detail">
    <!-- Back nav -->
    <RouterLink to="/pets" class="back-link">
      <IconArrowLeft :size="16" :stroke-width="2.5" />
      Volver a mascotas
    </RouterLink>

    <!-- Loading -->
    <div v-if="isLoading" class="feedback-state">
      <div class="spinner" />
      <p>Cargando…</p>
    </div>

    <!-- Error / not found -->
    <div v-else-if="isError || !pet" class="feedback-state feedback-state--error">
      <IconPaw :size="40" :stroke-width="1.5" />
      <p>Mascota no encontrada.</p>
      <RouterLink to="/pets" class="btn-back">Volver al listado</RouterLink>
    </div>

    <!-- Detail card -->
    <template v-else>
      <div class="detail-layout">
        <!-- Left: main card -->
        <div class="detail-card">
          <!-- Strip -->
          <div class="card-strip" :class="`strip--${pet.species}`" />

          <div class="card-content">
            <!-- Avatar -->
            <div class="species-avatar" :class="`avatar--${pet.species}`">
              <span class="species-emoji">{{ speciesEmoji(pet.species) }}</span>
            </div>

            <!-- Name + badges -->
            <div class="pet-header">
              <div class="pet-name-row">
                <h1 class="pet-name">{{ pet.name }}</h1>
                <span v-if="isBirthday" class="birthday-icon" title="Cumpleaños hoy">🎂</span>
              </div>
              <div class="badges">
                <span class="badge badge--species">{{ speciesLabel(pet.species) }}</span>
                <span v-if="pet.breed" class="badge badge--breed">{{ pet.breed }}</span>
              </div>
            </div>

            <!-- Stats row -->
            <div class="stats-row">
              <!-- Age -->
              <div class="stat">
                <span class="stat-value">{{ ageText }}</span>
                <span class="stat-label">
                  {{ pet.birth_date_exact ? 'Edad' : 'Edad estimada' }}
                </span>
              </div>

              <!-- Exact birth date (only when birth_date_exact) -->
              <div v-if="pet.birth_date_exact" class="stat">
                <span class="stat-value">{{ formatBirthDate(pet.birth_date) }}</span>
                <span class="stat-label">Nacimiento</span>
              </div>

              <!-- Weight -->
              <div v-if="pet.weight_grams !== null" class="stat">
                <span class="stat-value">{{ formatWeight(pet.weight_grams!) }}</span>
                <span class="stat-label">Peso</span>
              </div>

              <!-- Registered date -->
              <div class="stat">
                <span class="stat-value">{{ formatDate(pet.created_at) }}</span>
                <span class="stat-label">Registrado</span>
              </div>
            </div>

            <!-- Life stage badge -->
            <div v-if="pet.life_stage" class="life-stage-row">
              <span class="life-stage-label">Etapa de vida:</span>
              <span
                class="life-stage-badge"
                :style="{
                  background: lifeStageStyle(pet.life_stage).bg,
                  color: lifeStageStyle(pet.life_stage).text,
                }"
              >
                {{ lifeStageLabel(pet.life_stage) }}
              </span>
            </div>

            <!-- Actions -->
            <div class="card-actions">
              <button class="btn-edit" @click="showEditModal = true">
                <IconEdit :size="15" :stroke-width="2" />
                Editar
              </button>
              <button class="btn-delete" :disabled="deletePet.isPending.value" @click="openDeleteConfirm">
                <span v-if="deletePet.isPending.value" class="spinner-sm" />
                <IconTrash v-else :size="15" :stroke-width="2" />
                Eliminar
              </button>
            </div>
          </div>
        </div>

        <!-- Right: health card with tabs -->
        <div class="health-card">
          <div class="card-strip" :class="`strip--${pet.species}`" />

          <div class="health-card-content">
            <!-- Health header -->
            <div class="health-header">
              <h2 class="health-title">Gestión de salud</h2>
              <p class="health-subtitle">Historial médico de {{ pet.name }}</p>
            </div>

            <!-- Tabs navigation -->
            <nav class="tabs-nav">
              <RouterLink
                v-for="tab in TABS"
                :key="tab.name"
                :to="{ name: tab.name, params: { id } }"
                class="tab-link"
                :class="{ 'tab-link--active': route.name === tab.name }"
              >
                <component :is="tab.icon" :size="18" :stroke-width="1.75" />
                <span class="tab-label">{{ tab.label }}</span>
              </RouterLink>
            </nav>

            <!-- Tab content -->
            <div class="tab-content">
              <RouterView />
            </div>
          </div>
        </div>
      </div>
    </template>

    <!-- Edit modal -->
    <PetFormModal
      v-if="pet"
      v-model="showEditModal"
      mode="edit"
      :pet="pet ?? null"
    />

    <!-- Delete confirmation modal -->
    <ConfirmDeleteModal
      v-model="showConfirmDelete"
      :pet="pet ?? null"
      :deleting="deletePet.isPending.value"
      @confirm="handleDelete"
    />
  </div>
</template>

<style scoped>
.pet-detail {
  width: 100%;
  padding: var(--space-8) var(--space-10);
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

@media (max-width: 768px) {
  .pet-detail {
    padding: var(--space-5) var(--space-4);
    gap: var(--space-4);
  }
}

/* ── Back link ───────────────────────── */
.back-link {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-tertiary);
  text-decoration: none;
  transition: color var(--transition-fast);
  width: fit-content;
}
.back-link:hover {
  color: var(--color-accent);
}

/* ── Detail layout ───────────────────── */
.detail-layout {
  display: flex;
  gap: var(--space-6);
  align-items: flex-start;
  width: 100%;
}

@media (max-width: 900px) {
  .detail-layout {
    flex-direction: column;
  }
}

/* ── Detail card ─────────────────────── */
.detail-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-xl);
  overflow: hidden;
  box-shadow: var(--shadow-md);
  width: 340px;
  flex-shrink: 0;
}

@media (max-width: 900px) {
  .detail-card {
    width: 100%;
  }
}

@media (max-width: 768px) {
  .detail-card {
    max-width: 100%;
  }

  .card-content {
    padding: var(--space-5);
  }
}

.card-strip {
  height: 6px;
  width: 100%;
}
.strip--dog    { background: linear-gradient(90deg, #f59e0b, #fbbf24); }
.strip--cat    { background: linear-gradient(90deg, #8b5cf6, #a78bfa); }
.strip--bird   { background: linear-gradient(90deg, #06b6d4, #22d3ee); }
.strip--rabbit { background: linear-gradient(90deg, #ec4899, #f472b6); }
.strip--fish   { background: linear-gradient(90deg, #3b82f6, #60a5fa); }
.strip--other  { background: linear-gradient(90deg, var(--color-accent), var(--color-accent-muted)); }

.card-content {
  padding: var(--space-8);
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

/* ── Health card ──────────────────────── */
.health-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-xl);
  overflow: hidden;
  box-shadow: var(--shadow-md);
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
}

@media (max-width: 900px) {
  .health-card {
    width: 100%;
  }
}

.health-card-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  padding: var(--space-6);
}

/* ── Health header ─────────────────────── */
.health-header {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.health-title {
  font-family: var(--font-display);
  font-size: var(--text-xl);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
  line-height: 1.2;
}

.health-subtitle {
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
  margin: 0;
  line-height: 1.4;
}

/* ── Tabs navigation ─────────────────── */
.tabs-nav {
  display: flex;
  gap: var(--space-2);
  padding: var(--space-1);
  background: var(--color-bg-alt);
  border-radius: var(--radius-lg);
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.tab-link {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-5);
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-tertiary);
  text-decoration: none;
  white-space: nowrap;
  border-radius: var(--radius-md);
  transition: all var(--transition-fast);
  flex: 1;
}

.tab-link:hover:not(.tab-link--active) {
  color: var(--color-text-secondary);
  background: var(--color-surface);
}

.tab-link--active {
  color: var(--color-accent-dark);
  background: var(--color-surface);
  box-shadow: var(--shadow-sm);
}

.tab-label {
  display: inline;
}

@media (max-width: 480px) {
  .tab-label {
    display: none;
  }
  
  .tab-link {
    padding: var(--space-3) var(--space-4);
  }
}

.tab-content {
  flex: 1;
  min-width: 0;
}

/* ── Species avatar ──────────────────── */
.species-avatar {
  width: 80px;
  height: 80px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
}
.avatar--dog    { background: #fef3c7; }
.avatar--cat    { background: #ede9fe; }
.avatar--bird   { background: #cffafe; }
.avatar--rabbit { background: #fce7f3; }
.avatar--fish   { background: #dbeafe; }
.avatar--other  { background: var(--color-accent-light); }

.species-emoji { font-size: 2.75rem; line-height: 1; }

.pet-header {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.pet-name-row {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.pet-name {
  font-family: var(--font-display);
  font-size: var(--text-3xl);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
  line-height: 1.15;
}

.birthday-icon {
  font-size: 1.75rem;
  line-height: 1;
  flex-shrink: 0;
}

/* ── Life stage ──────────────────────── */
.life-stage-row {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.life-stage-label {
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
  font-weight: 500;
}

.life-stage-badge {
  display: inline-flex;
  align-items: center;
  padding: 3px var(--space-3);
  border-radius: var(--radius-full);
  font-size: var(--text-sm);
  font-weight: 600;
  letter-spacing: 0.02em;
}

.badges {
  display: flex;
  flex-wrap: wrap;
  gap: var(--space-2);
}

.badge {
  display: inline-flex;
  align-items: center;
  padding: 3px var(--space-3);
  border-radius: var(--radius-full);
  font-size: var(--text-sm);
  font-weight: 600;
}
.badge--species {
  background: var(--color-accent-light);
  color: var(--color-accent-dark);
}
.badge--breed {
  background: var(--color-bg-alt);
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border-light);
}

/* ── Stats row ───────────────────────── */
.stats-row {
  display: flex;
  gap: var(--space-6);
  flex-wrap: wrap;
  padding: var(--space-4) 0;
  border-top: 1px solid var(--color-border-light);
  border-bottom: 1px solid var(--color-border-light);
}

.stat {
  display: flex;
  flex-direction: column;
  gap: 2px;
}
.stat-value {
  font-family: var(--font-display);
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-text-primary);
}
.stat-label {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

/* ── Card actions ────────────────────── */
.card-actions {
  display: flex;
  gap: var(--space-3);
}

.btn-edit {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-5);
  background: var(--color-accent);
  color: var(--color-text-on-accent);
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast), transform var(--transition-fast);
}
.btn-edit:hover {
  background: var(--color-accent-hover);
  transform: translateY(-1px);
}

.btn-delete {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  background: transparent;
  color: var(--color-error);
  border: 1.5px solid var(--color-error-border);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast);
}
.btn-delete:hover:not(:disabled) {
  background: var(--color-error-light);
  border-color: var(--color-error);
}
.btn-delete:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* ── Feedback states ─────────────────── */
.feedback-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-16);
  color: var(--color-text-tertiary);
  text-align: center;
}
.feedback-state--error {
  color: var(--color-error);
}

.btn-back {
  padding: var(--space-2) var(--space-5);
  background: var(--color-accent);
  color: var(--color-text-on-accent);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  text-decoration: none;
  transition: background var(--transition-fast);
}
.btn-back:hover { background: var(--color-accent-hover); }

/* ── Spinner ─────────────────────────── */
.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--color-border-light);
  border-top-color: var(--color-accent);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

.spinner-sm {
  display: inline-block;
  width: 13px;
  height: 13px;
  border: 2px solid rgba(220, 38, 38, 0.3);
  border-top-color: var(--color-error);
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
