<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import { useRoute } from 'vue-router'
import { IconPlus, IconSearch, IconAlertCircle, IconRefresh } from '@tabler/icons-vue'
import { useGetHealthCatalogs, useDeleteHealthCatalog } from '@/composables/useHealthCatalog'
import { useUIStore } from '@/stores/ui'
import { PET_SPECIES } from '@/constants/species'
import type { HealthCatalog, HealthCatalogCategory } from '@/types/healthCatalog'
import HealthCatalogFormModal from '@/components/health-catalog/HealthCatalogFormModal.vue'
import HealthCatalogTableRow from '@/components/health-catalog/HealthCatalogTableRow.vue'
import AppPagination from '@/components/ui/AppPagination.vue'
import PerPageSelector from '@/components/ui/PerPageSelector.vue'

const uiStore = useUIStore()
const route = useRoute()

// El category actual viene de la URL. Si no hay (no debería pasar por el redirect), fallback a vaccine
const category = computed(() => {
  const cat = route.params.category as string
  if (['vaccine', 'deworming', 'exam'].includes(cat)) {
    return cat as HealthCatalogCategory
  }
  return 'vaccine'
})

// ── Species filter ──────────────────────────────────────────────
const speciesFilter = ref<string | undefined>(undefined)

// Al cambiar species, resetear a página 1
watch(speciesFilter, () => {
  page.value = 1
})

// ── Pagination state ────────────────────────────────────────
const page = ref(1)
const perPage = ref(uiStore.usersPerPage)

// Al cambiar perPage, sincronizar con store y resetear a página 1
watch(perPage, (val) => {
  uiStore.usersPerPage = val as 10 | 25 | 50
  page.value = 1
})

const { data, isLoading, isError, error: fetchError, refresh, isFetching } = useGetHealthCatalogs(category, page, perPage, speciesFilter)

const deleteItem = useDeleteHealthCatalog()

// ── Search (filtrado local sobre la página actual) ──────────
const search = ref('')

// ── Modal state ─────────────────────────────────────────────
type ModalMode = 'create' | 'edit'
const showModal = ref(false)
const modalMode = ref<ModalMode>('create')
const editingItem = ref<HealthCatalog | null>(null)

function openCreate() {
  modalMode.value = 'create'
  editingItem.value = null
  showModal.value = true
}

function openEdit(item: HealthCatalog) {
  modalMode.value = 'edit'
  editingItem.value = item
  showModal.value = true
}

function closeModal() {
  showModal.value = false
}

// ── Delete ───────────────────────────────────────────────────
const deletingId = ref<string | null>(null)

async function handleDelete(item: HealthCatalog) {
  if (!confirm(`¿Eliminar "${item.name}"? Esta acción no se puede deshacer.`)) return
  deletingId.value = item.id
  try {
    await deleteItem.mutateAsync(item.id)
  } catch (e) {
    alert(e instanceof Error ? e.message : 'Error al eliminar registro')
  } finally {
    deletingId.value = null
  }
}

// ── Derived ──────────────────────────────────────────────────
const items = ref<HealthCatalog[]>([])
const total = ref(0)
const totalPages = ref(1)

watch(
  data,
  (res) => {
    if (!res) return
    const q = search.value.toLowerCase()
    items.value = q
      ? res.data.filter((v) => v.name.toLowerCase().includes(q) || v.species.some((s) => s.toLowerCase().includes(q)))
      : res.data
    total.value = res.total
    totalPages.value = res.total_pages
  },
  { immediate: true },
)

watch(search, () => {
  if (!data.value) return
  const q = search.value.toLowerCase()
  items.value = q
    ? data.value.data.filter((v) => v.name.toLowerCase().includes(q) || v.species.some((s) => s.toLowerCase().includes(q)))
    : data.value.data
})
</script>

<template>
  <div class="health-catalog-view">
    <!-- Toolbar -->
    <div class="toolbar">
      <div class="toolbar-actions">
        <button class="btn-refresh" title="Refrescar" :disabled="isFetching" @click="refresh">
          <IconRefresh :size="16" :stroke-width="2" :class="{ spinning: isFetching }" />
          <span>Refrescar</span>
        </button>
        <button class="btn-create" @click="openCreate">
          <IconPlus :size="16" :stroke-width="2.5" />
          {{ category === 'vaccine' ? 'Nueva vacuna' : category === 'deworming' ? 'Nueva desparasitación' : 'Nuevo examen' }}
        </button>
      </div>

      <div class="toolbar-filters">
        <div class="search-box">
          <IconSearch class="search-icon" :size="16" :stroke-width="2" />
          <input v-model="search" class="search-input" placeholder="Buscar por nombre…" />
        </div>
        <select v-model="speciesFilter" class="species-filter">
          <option :value="undefined">Todas las especies</option>
          <option v-for="s in PET_SPECIES" :key="s.value" :value="s.value">
            {{ s.icon }} {{ s.label }}
          </option>
        </select>
      </div>
      <div class="stats-pill">
        <span class="stats-num">{{ total }}</span>
        <span class="stats-label">{{ total === 1 ? 'registro' : 'registros' }}</span>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="isLoading" class="empty-state">
      <div class="spinner" />
      <p>Cargando registros…</p>
    </div>

    <!-- Error -->
    <div v-else-if="isError" class="empty-state empty-state--error">
      <IconAlertCircle :size="40" :stroke-width="1.5" />
      <p>{{ fetchError?.message }}</p>
      <button class="btn-secondary" @click="refresh">Reintentar</button>
    </div>

    <!-- Empty search -->
    <div v-else-if="items.length === 0 && search" class="empty-state">
      <IconSearch :size="40" :stroke-width="1.5" />
      <p>Sin resultados para "<strong>{{ search }}</strong>"</p>
    </div>

    <!-- Empty state (no items) -->
    <div v-else-if="items.length === 0" class="empty-state">
      <IconSearch :size="40" :stroke-width="1.5" />
      <p>No hay registros para esta categoría</p>
      <button class="btn-secondary" @click="openCreate">Crear primer registro</button>
    </div>

    <!-- Table -->
    <div v-else class="table-wrapper">
      <table class="health-catalog-table">
        <thead>
          <tr>
            <th>Nombre</th>
            <th class="th-center th-species">Especies</th>
            <th class="th-center">Frecuencia</th>
            <th class="th-center">Obligatoria</th>
            <th class="th-center">Creado</th>
            <th class="th-center">Acciones</th>
          </tr>
        </thead>
        <tbody>
          <HealthCatalogTableRow
            v-for="item in items"
            :key="item.id"
            :item="item"
            :deleting-id="deletingId"
            @edit="openEdit"
            @delete="handleDelete"
          />
        </tbody>
      </table>
    </div>

    <!-- Footer: per-page + pagination -->
    <div v-if="!isLoading && !isError && totalPages > 0" class="table-footer">
      <PerPageSelector
        v-model="perPage"
        :options="[10, 25, 50]"
      />
      <AppPagination
        :current-page="page"
        :total-pages="totalPages"
        :total-items="total"
        :per-page="perPage"
        @update:page="page = $event"
      />
    </div>

    <!-- Modal -->
    <HealthCatalogFormModal
      v-if="showModal"
      :mode="modalMode"
      :item="editingItem ?? undefined"
      :default-category="category"
      @close="closeModal"
      @saved="closeModal"
    />
  </div>
</template>

<style scoped>
/* ── Layout ────────────────────────────────────────────── */
.health-catalog-view {
  width: 100%;
  max-width: 100%;
  padding: var(--space-6) 0;
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
  box-sizing: border-box;
  overflow-x: hidden;
}

@media (max-width: 1024px) {
  .health-catalog-view {
    padding: var(--space-4) 0;
  }
}

@media (max-width: 768px) {
  .health-catalog-view {
    padding: var(--space-3) 0;
    gap: var(--space-4);
  }
}

/* ── Toolbar ───────────────────────────────────────────── */
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
  background: #fff;
  padding: var(--space-4);
  border-radius: var(--radius-lg);
  border: 1px solid var(--color-border-light);
}

.toolbar-actions {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  border-right: 1px solid var(--color-border-light);
  padding-right: var(--space-4);
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
  transition: background var(--transition-fast), border-color var(--transition-fast),
    color var(--transition-fast);
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

.spinning {
  animation: spin 0.7s linear infinite;
}

.btn-create {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-5);
  background: var(--color-accent);
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast), transform var(--transition-fast);
  white-space: nowrap;
  flex-shrink: 0;
}

.btn-create:hover {
  background: #369a6e;
  transform: translateY(-1px);
}

/* ── Toolbar ───────────────────────────────────────────── */
.toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
}

.toolbar-filters {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  flex: 1;
  max-width: 500px;
}

.search-box {
  position: relative;
  flex: 1;
  min-width: 0;
}

.species-filter {
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  background: #fff;
  cursor: pointer;
  min-width: 140px;
}

.species-filter:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(66, 184, 131, 0.15);
}

@media (max-width: 768px) {
  .toolbar {
    flex-wrap: wrap;
    flex-direction: column;
    align-items: stretch;
  }

  .toolbar-actions {
    border-right: none;
    padding-right: 0;
    border-bottom: 1px solid var(--color-border-light);
    padding-bottom: var(--space-3);
    justify-content: space-between;
  }

  .toolbar-filters {
    max-width: 100%;
  }

  .stats-pill {
    margin-left: 0;
    align-self: flex-start;
  }
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
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  background: var(--color-bg);
  transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
  box-sizing: border-box;
}

.search-input:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(66, 184, 131, 0.15);
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

/* ── Table ─────────────────────────────────────────────── */
.table-wrapper {
  background: #fff;
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  overflow: hidden;
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  width: 100%;
}

.health-catalog-table {
  width: 100%;
  min-width: 760px;
  border-collapse: collapse;
}

.health-catalog-table thead tr {
  background: var(--color-bg);
  border-bottom: 1px solid var(--color-border-light);
}

.health-catalog-table th {
  padding: var(--space-3) var(--space-4);
  text-align: left;
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.health-catalog-table th.th-center {
  text-align: center;
}

.health-catalog-table th.th-species {
  min-width: 180px;
}

/* ── Empty / loading states ────────────────────────────── */
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--space-3);
  padding: var(--space-12) var(--space-4);
  color: var(--color-text-tertiary);
  text-align: center;
}

.empty-state--error {
  color: var(--color-error);
}

/* ── Spinner ───────────────────────────────────────────── */
.spinner {
  width: 28px;
  height: 28px;
  border: 3px solid var(--color-border-light);
  border-top-color: var(--color-accent);
  border-radius: var(--radius-full);
  animation: spin 0.7s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* ── Buttons ───────────────────────────────────────────── */
.btn-secondary {
  padding: var(--space-2) var(--space-4);
  background: #fff;
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 500;
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast);
}

.btn-secondary:hover {
  background: var(--color-bg);
  border-color: #d1d5db;
}

/* ── Table footer ──────────────────────────────────────── */
.table-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
  flex-wrap: wrap;
}

@media (max-width: 600px) {
  .table-footer {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-3);
  }
}
</style>
