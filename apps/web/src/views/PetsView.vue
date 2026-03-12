<script setup lang="ts">
import { ref, computed } from 'vue'
import { IconPlus, IconRefresh, IconSearch, IconAlertCircle } from '@tabler/icons-vue'
import { useGetPets, useDeletePet } from '@/composables/usePets'
import { useAuthStore } from '@/stores/auth'
import PetCard from '@/components/pets/PetCard.vue'
import PetFormModal from '@/components/pets/PetFormModal.vue'
import PetEmptyState from '@/components/pets/PetEmptyState.vue'
import type { Pet } from '@/types/pet'

const { data, allPets, total, hasMore, isLoading, isError, error, refresh, isFetching, loadMore } =
  useGetPets()
const deletePet = useDeletePet()
const authStore = useAuthStore()

const petLimit = computed(() => authStore.user?.pet_limit ?? 5)
const petCount = computed(() => total.value ?? 0)
const atLimit = computed(() => petCount.value >= petLimit.value)

const showModal = ref(false)
const modalMode = ref<'create' | 'edit'>('create')
const editingPet = ref<Pet | undefined>(undefined)
const deletingId = ref<string | null>(null)
const search = ref('')

// Filtrar sobre los datos acumulados (no va a la BD)
const filteredPets = computed(() => {
  const list = data.value ?? []
  const q = search.value.toLowerCase().trim()
  if (!q) return list
  return list.filter(
    (p) =>
      p.name.toLowerCase().includes(q) ||
      p.species.toLowerCase().includes(q) ||
      p.breed?.toLowerCase().includes(q),
  )
})

function openCreate() {
  editingPet.value = undefined
  modalMode.value = 'create'
  showModal.value = true
}

function openEdit(pet: Pet) {
  editingPet.value = pet
  modalMode.value = 'edit'
  showModal.value = true
}

async function handleDelete(id: string) {
  if (!confirm('¿Eliminar esta mascota?')) return
  deletingId.value = id
  try {
    await deletePet.mutateAsync(id)
  } finally {
    deletingId.value = null
  }
}
</script>

<template>
  <div class="pets-view">
    <!-- Header -->
    <div class="page-header">
      <div class="page-header__text">
        <h1 class="page-title">Mascotas</h1>
        <p class="page-subtitle">Administrá el registro de tus mascotas</p>
      </div>
      <div class="header-actions">
        <button
          class="btn-refresh"
          title="Refrescar"
          :disabled="isFetching"
          @click="refresh"
        >
          <IconRefresh :size="16" :stroke-width="2" :class="{ spinning: isFetching }" />
          <span>Refrescar</span>
        </button>
        <button
          class="btn-create"
          :disabled="atLimit"
          :title="atLimit ? 'Has alcanzado el límite de mascotas' : 'Nueva mascota'"
          @click="openCreate"
        >
          <IconPlus :size="16" :stroke-width="2.5" />
          Nueva mascota
        </button>
      </div>
    </div>

    <!-- Toolbar: search + stats -->
    <div class="toolbar">
      <div class="search-box">
        <IconSearch class="search-icon" :size="16" :stroke-width="2" />
        <input
          v-model="search"
          class="search-input"
          placeholder="Buscar por nombre, especie, raza…"
        />
      </div>
      <div v-if="!isLoading && !isError" class="stats-pill">
        <span class="stats-num">{{ petCount }} / {{ petLimit }}</span>
        <span class="stats-label">{{ petCount === 1 ? 'mascota' : 'mascotas' }}</span>
      </div>
    </div>

    <!-- Loading inicial -->
    <div v-if="isLoading && allPets.length === 0" class="feedback-state">
      <div class="spinner" />
      <p>Cargando mascotas…</p>
    </div>

    <!-- Error -->
    <div v-else-if="isError" class="feedback-state feedback-state--error">
      <IconAlertCircle :size="40" :stroke-width="1.5" />
      <p>{{ error?.message ?? 'Error al cargar mascotas' }}</p>
      <button class="btn-retry" @click="refresh">Reintentar</button>
    </div>

    <!-- Empty (sin mascotas) -->
    <PetEmptyState
      v-else-if="!filteredPets.length && !search"
      @add="openCreate"
    />

    <!-- Empty search -->
    <PetEmptyState
      v-else-if="!filteredPets.length && search"
      :searching="true"
      :query="search"
    />

    <!-- Grid -->
    <div v-else class="pets-grid">
      <PetCard
        v-for="pet in filteredPets"
        :key="pet.id"
        :pet="pet"
        :deleting="deletingId === pet.id"
        @edit="openEdit"
        @delete="handleDelete"
      />
    </div>

    <!-- Load more -->
    <div v-if="hasMore && !search && !isError" class="load-more-bar">
      <button class="btn-load-more" :disabled="isFetching" @click="loadMore">
        <div v-if="isFetching" class="spinner spinner--sm" />
        <span v-else>Ver más mascotas</span>
      </button>
    </div>

    <!-- Modal -->
    <PetFormModal
      v-model="showModal"
      :mode="modalMode"
      :pet="editingPet"
    />
  </div>
</template>

<style scoped>
.pets-view {
  width: 100%;
  padding: var(--space-8) var(--space-10);
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

@media (max-width: 768px) {
  .pets-view {
    padding: var(--space-5) var(--space-4);
    gap: var(--space-4);
  }
}

/* ── Header ─────────────────────────── */
.page-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: var(--space-4);
}

@media (max-width: 480px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-3);
  }

  .header-actions {
    width: 100%;
  }

  .btn-create {
    flex: 1;
    justify-content: center;
  }
}

.page-title {
  font-size: var(--text-2xl);
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0 0 var(--space-1);
  line-height: 1.2;
}

.page-subtitle {
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
  margin: 0;
}

.header-actions {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  flex-shrink: 0;
}

.btn-refresh {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-3);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  background: var(--color-surface);
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast), color var(--transition-fast);
  font-size: var(--text-sm);
  font-weight: 500;
}

.btn-refresh:hover:not(:disabled) {
  background: var(--color-bg-alt);
  border-color: var(--color-border);
  color: var(--color-text-primary);
}

.btn-refresh:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-create {
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
  white-space: nowrap;
}
.btn-create:hover:not(:disabled) {
  background: var(--color-accent-hover);
  transform: translateY(-1px);
}
.btn-create:disabled {
  background: var(--color-border);
  color: var(--color-text-tertiary);
  cursor: not-allowed;
}

/* ── Toolbar ─────────────────────────── */
.toolbar {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

@media (max-width: 480px) {
  .toolbar {
    flex-wrap: wrap;
  }

  .search-box {
    max-width: 100%;
    flex: 1 1 100%;
  }
}

.search-box {
  position: relative;
  flex: 1;
  max-width: 400px;
}

.search-icon {
  position: absolute;
  left: var(--space-3);
  top: 50%;
  transform: translateY(-50%);
  color: var(--color-text-tertiary);
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: var(--space-2) var(--space-3) var(--space-2) 2.25rem;
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  background: var(--color-surface);
  transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
  box-sizing: border-box;
}
.search-input:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(61, 122, 95, 0.12);
}

.stats-pill {
  margin-left: auto;
  display: flex;
  align-items: baseline;
  gap: var(--space-1);
  background: var(--color-accent-light);
  border-radius: var(--radius-full);
  padding: var(--space-1) var(--space-3);
}
.stats-num {
  font-size: var(--text-base);
  font-weight: 700;
  color: var(--color-accent);
}
.stats-label {
  font-size: var(--text-xs);
  color: var(--color-accent);
  font-weight: 500;
}

/* ── Grid ────────────────────────────── */
.pets-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(260px, 1fr));
  gap: var(--space-4);
}

/* ── Load more ───────────────────────── */
.load-more-bar {
  display: flex;
  justify-content: center;
  padding: var(--space-2) 0 var(--space-4);
}

.btn-load-more {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-6);
  border: 1.5px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-surface);
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
  font-weight: 500;
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast), color var(--transition-fast);
}

.btn-load-more:hover:not(:disabled) {
  background: var(--color-accent-light);
  border-color: var(--color-accent-muted);
  color: var(--color-accent-dark);
}

.btn-load-more:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

/* ── Feedback states ─────────────────── */
.feedback-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--space-3);
  padding: var(--space-16) var(--space-4);
  color: var(--color-text-tertiary);
  text-align: center;
}
.feedback-state--error {
  color: var(--color-error);
}

.btn-retry {
  padding: var(--space-2) var(--space-5);
  background: transparent;
  color: var(--color-error);
  border: 1.5px solid var(--color-error-border);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast);
}
.btn-retry:hover {
  background: var(--color-error-light);
}

/* ── Spinner ─────────────────────────── */
.spinner {
  width: 32px;
  height: 32px;
  border: 3px solid var(--color-border-light);
  border-top-color: var(--color-accent);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

.spinner--sm {
  width: 16px;
  height: 16px;
  border-width: 2px;
}

.spinning {
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>

