<script setup lang="ts">
import { ref, watch } from 'vue'
import { IconPlus, IconSearch, IconAlertCircle, IconRefresh } from '@tabler/icons-vue'
import { useGetUsers, useDeleteUser } from '@/composables/useUsers'
import { useUIStore } from '@/stores/ui'
import type { User } from '@/types/user'
import UserFormModal from '@/components/users/UserFormModal.vue'
import UserTableRow from '@/components/users/UserTableRow.vue'
import AppPagination from '@/components/ui/AppPagination.vue'
import PerPageSelector from '@/components/ui/PerPageSelector.vue'

const uiStore = useUIStore()

// ── Pagination state ────────────────────────────────────────
const page = ref(1)
const perPage = ref(uiStore.usersPerPage)

// Al cambiar perPage, sincronizar con store y resetear a página 1
watch(perPage, (val) => {
  uiStore.usersPerPage = val as 10 | 25 | 50
  page.value = 1
})

const { data, isLoading, isError, error: fetchError, refresh, isFetching } = useGetUsers(page, perPage)

const deleteUser = useDeleteUser()

// ── Search (filtrado local sobre la página actual) ──────────
const search = ref('')

// ── Modal state ─────────────────────────────────────────────
type ModalMode = 'create' | 'edit'
const showModal = ref(false)
const modalMode = ref<ModalMode>('create')
const editingUser = ref<User | null>(null)

function openCreate() {
  modalMode.value = 'create'
  editingUser.value = null
  showModal.value = true
}

function openEdit(user: User) {
  modalMode.value = 'edit'
  editingUser.value = user
  showModal.value = true
}

function closeModal() {
  showModal.value = false
}

// ── Delete ───────────────────────────────────────────────────
const deletingId = ref<string | null>(null)

async function handleDelete(user: User) {
  if (!confirm(`¿Eliminar a ${user.name}? Esta acción no se puede deshacer.`)) return
  deletingId.value = user.id
  try {
    await deleteUser.mutateAsync(user.id)
  } catch (e) {
    alert(e instanceof Error ? e.message : 'Error al eliminar usuario')
  } finally {
    deletingId.value = null
  }
}

// ── Derived ──────────────────────────────────────────────────
const users = ref<User[]>([])
const total = ref(0)
const totalPages = ref(1)

watch(
  data,
  (res) => {
    if (!res) return
    const q = search.value.toLowerCase()
    users.value = q
      ? res.data.filter(
          (u) => u.name.toLowerCase().includes(q) || u.email.toLowerCase().includes(q),
        )
      : res.data
    total.value = res.total
    totalPages.value = res.total_pages
  },
  { immediate: true },
)

watch(search, () => {
  if (!data.value) return
  const q = search.value.toLowerCase()
  users.value = q
    ? data.value.data.filter(
        (u) => u.name.toLowerCase().includes(q) || u.email.toLowerCase().includes(q),
      )
    : data.value.data
})
</script>

<template>
  <div class="users-view">
    <!-- Header -->
    <div class="page-header">
      <div class="page-header__text">
        <h1 class="page-title">Gestión de usuarios</h1>
        <p class="page-subtitle">Administra los accesos al sistema</p>
      </div>
      <div class="header-actions">
        <button class="btn-refresh" title="Refrescar" :disabled="isFetching" @click="refresh">
          <IconRefresh :size="16" :stroke-width="2" :class="{ spinning: isFetching }" />
          <span>Refrescar</span>
        </button>
        <button class="btn-create" @click="openCreate">
          <IconPlus :size="16" :stroke-width="2.5" />
          Nuevo usuario
        </button>
      </div>
    </div>

    <!-- Toolbar -->
    <div class="toolbar">
      <div class="search-box">
        <IconSearch class="search-icon" :size="16" :stroke-width="2" />
        <input v-model="search" class="search-input" placeholder="Buscar por nombre o email…" />
      </div>
      <div class="stats-pill">
        <span class="stats-num">{{ total }}</span>
        <span class="stats-label">{{ total === 1 ? 'usuario' : 'usuarios' }}</span>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="isLoading" class="empty-state">
      <div class="spinner" />
      <p>Cargando usuarios…</p>
    </div>

    <!-- Error -->
    <div v-else-if="isError" class="empty-state empty-state--error">
      <IconAlertCircle :size="40" :stroke-width="1.5" />
      <p>{{ fetchError?.message }}</p>
      <button class="btn-secondary" @click="refresh">Reintentar</button>
    </div>

    <!-- Empty search -->
    <div v-else-if="users.length === 0 && search" class="empty-state">
      <IconSearch :size="40" :stroke-width="1.5" />
      <p>Sin resultados para "<strong>{{ search }}</strong>"</p>
    </div>

    <!-- Table -->
    <div v-else class="table-wrapper">
      <table class="users-table">
        <thead>
          <tr>
            <th>Usuario</th>
            <th class="th-center">Rol</th>
            <th class="th-center">Acceso</th>
            <th class="th-center">Creado</th>
            <th class="th-center">Acciones</th>
          </tr>
        </thead>
        <tbody>
          <UserTableRow
            v-for="user in users"
            :key="user.id"
            :user="user"
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
    <UserFormModal
      v-if="showModal"
      :mode="modalMode"
      :user="editingUser ?? undefined"
      @close="closeModal"
      @saved="closeModal"
    />
  </div>
</template>

<style scoped>
/* ── Layout ────────────────────────────────────────────── */
.users-view {
  width: 100%;
  max-width: 100%;
  padding: var(--space-8) var(--space-10);
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
  box-sizing: border-box;
  overflow-x: hidden;
}

@media (max-width: 1024px) {
  .users-view {
    padding: var(--space-6) var(--space-6);
  }
}

@media (max-width: 768px) {
  .users-view {
    padding: var(--space-5) var(--space-4);
    gap: var(--space-4);
  }
}

/* ── Header ────────────────────────────────────────────── */
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

  .btn-refresh span {
    display: none;
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
  gap: var(--space-3);
}

.search-box {
  position: relative;
  flex: 1;
  max-width: 360px;
  min-width: 0;
}

@media (max-width: 480px) {
  .toolbar {
    flex-wrap: wrap;
  }

  .search-box {
    max-width: 100%;
    width: 100%;
  }

  .stats-pill {
    margin-left: 0;
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

.users-table {
  width: 100%;
  min-width: 560px;
  border-collapse: collapse;
}

.users-table thead tr {
  background: var(--color-bg);
  border-bottom: 1px solid var(--color-border-light);
}

.users-table th {
  padding: var(--space-3) var(--space-4);
  text-align: left;
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.users-table th.th-center {
  text-align: center;
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
