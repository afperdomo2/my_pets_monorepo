<script setup lang="ts">
import { onMounted, ref, computed } from 'vue'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import {
  IconPlus,
  IconSearch,
  IconAlertCircle,
  IconEdit,
  IconTrash,
  IconX,
  IconLock,
} from '@tabler/icons-vue'
import { useUserStore } from '@/stores/users'
import type { User } from '@/types/user'
import { createUserSchema, updateUserSchema } from '@/schemas/user'

const store = useUserStore()

// ── Modal state ────────────────────────────────────────────
type ModalMode = 'create' | 'edit'
const showModal = ref(false)
const modalMode = ref<ModalMode>('create')
const editingUser = ref<User | null>(null)

const submitError = ref<string | null>(null)
const submitting = ref(false)

// ── Search ─────────────────────────────────────────────────
const search = ref('')
const filteredUsers = computed(() => {
  const q = search.value.toLowerCase()
  if (!q) return store.users
  return store.users.filter(
    (u) => u.name.toLowerCase().includes(q) || u.email.toLowerCase().includes(q),
  )
})

onMounted(() => store.fetchUsers())

// ── Helpers ────────────────────────────────────────────────
function initials(name: string): string {
  return name
    .split(' ')
    .slice(0, 2)
    .map((w) => w[0])
    .join('')
    .toUpperCase()
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString('es-ES', {
    day: '2-digit',
    month: 'short',
    year: 'numeric',
  })
}

// ── Create form ─────────────────────────────────────────────
const {
  defineField: defineCreateField,
  handleSubmit: handleCreateSubmit,
  errors: createErrors,
  resetForm: resetCreateForm,
} = useForm({
  validationSchema: toTypedSchema(createUserSchema),
  initialValues: { name: '', email: '', password: '' },
})

const [createName, createNameAttrs] = defineCreateField('name')
const [createEmail, createEmailAttrs] = defineCreateField('email')
const [createPassword, createPasswordAttrs] = defineCreateField('password')

function openCreate() {
  modalMode.value = 'create'
  resetCreateForm({ values: { name: '', email: '', password: '' } })
  submitError.value = null
  showModal.value = true
}

const handleCreate = handleCreateSubmit(async (values) => {
  submitting.value = true
  submitError.value = null
  try {
    await store.createUser(values)
    showModal.value = false
  } catch (e) {
    submitError.value = e instanceof Error ? e.message : 'Error al crear usuario'
  } finally {
    submitting.value = false
  }
})

// ── Edit form ───────────────────────────────────────────────
const {
  defineField: defineEditField,
  handleSubmit: handleEditSubmit,
  errors: editErrors,
  resetForm: resetEditForm,
} = useForm({
  validationSchema: toTypedSchema(updateUserSchema),
  initialValues: { name: '', email: '' },
})

const [editName, editNameAttrs] = defineEditField('name')
const [editEmail, editEmailAttrs] = defineEditField('email')

function openEdit(user: User) {
  modalMode.value = 'edit'
  editingUser.value = user
  resetEditForm({ values: { name: user.name, email: user.email } })
  submitError.value = null
  showModal.value = true
}

const handleEdit = handleEditSubmit(async (values) => {
  if (!editingUser.value) return
  submitting.value = true
  submitError.value = null
  try {
    await store.updateUser(editingUser.value.id, values)
    showModal.value = false
  } catch (e) {
    submitError.value = e instanceof Error ? e.message : 'Error al actualizar usuario'
  } finally {
    submitting.value = false
  }
})

// ── Delete ─────────────────────────────────────────────────
const deletingId = ref<number | null>(null)

async function handleDelete(user: User) {
  if (!confirm(`¿Eliminar a ${user.name}? Esta acción no se puede deshacer.`)) return
  deletingId.value = user.id
  try {
    await store.deleteUser(user.id)
  } catch (e) {
    alert(e instanceof Error ? e.message : 'Error al eliminar usuario')
  } finally {
    deletingId.value = null
  }
}

function closeModal() {
  showModal.value = false
}
</script>

<template>
  <div class="users-view">
    <!-- Header -->
    <div class="page-header">
      <div class="page-header__text">
        <h1 class="page-title">Gestión de usuarios</h1>
        <p class="page-subtitle">Administra los accesos al sistema</p>
      </div>
      <button class="btn-create" @click="openCreate">
        <IconPlus :size="16" :stroke-width="2.5" />
        Nuevo usuario
      </button>
    </div>

    <!-- Search + stats -->
    <div class="toolbar">
      <div class="search-box">
        <IconSearch class="search-icon" :size="16" :stroke-width="2" />
        <input v-model="search" class="search-input" placeholder="Buscar por nombre o email…" />
      </div>
      <div class="stats-pill">
        <span class="stats-num">{{ store.users.length }}</span>
        <span class="stats-label">{{ store.users.length === 1 ? 'usuario' : 'usuarios' }}</span>
      </div>
    </div>

    <!-- Loading -->
    <div v-if="store.loading" class="empty-state">
      <div class="spinner" />
      <p>Cargando usuarios…</p>
    </div>

    <!-- Error -->
    <div v-else-if="store.error" class="empty-state empty-state--error">
      <IconAlertCircle :size="40" :stroke-width="1.5" />
      <p>{{ store.error }}</p>
      <button class="btn-secondary" @click="store.fetchUsers">Reintentar</button>
    </div>

    <!-- Empty search -->
    <div v-else-if="filteredUsers.length === 0 && search" class="empty-state">
      <IconSearch :size="40" :stroke-width="1.5" />
      <p>Sin resultados para "<strong>{{ search }}</strong>"</p>
    </div>

    <!-- Table -->
    <div v-else class="table-wrapper">
      <table class="users-table">
        <thead>
          <tr>
            <th>Usuario</th>
            <th>Rol</th>
            <th>Acceso</th>
            <th>Creado</th>
            <th class="th-actions">Acciones</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="user in filteredUsers" :key="user.id" class="user-row">
            <!-- Avatar + name + email -->
            <td class="td-user">
              <div
                class="avatar"
                :class="user.is_system_user ? 'avatar--admin' : 'avatar--regular'"
              >
                {{ initials(user.name) }}
              </div>
              <div class="user-info">
                <span class="user-name">{{ user.name }}</span>
                <span class="user-email">{{ user.email }}</span>
              </div>
            </td>

            <!-- Role badge -->
            <td>
              <span class="badge" :class="user.is_system_user ? 'badge--admin' : 'badge--user'">
                {{ user.is_system_user ? 'Administrador' : 'Usuario' }}
              </span>
            </td>

            <!-- Auth provider -->
            <td>
              <span class="provider" :class="`provider--${user.auth_provider}`">
                <svg
                  v-if="user.auth_provider === 'google'"
                  width="14"
                  height="14"
                  viewBox="0 0 24 24"
                  fill="currentColor"
                >
                  <path
                    d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z"
                    fill="#4285F4"
                  />
                  <path
                    d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z"
                    fill="#34A853"
                  />
                  <path
                    d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l2.85-2.22.81-.62z"
                    fill="#FBBC05"
                  />
                  <path
                    d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z"
                    fill="#EA4335"
                  />
                </svg>
                <IconLock v-else :size="14" :stroke-width="2" />
                {{ user.auth_provider === 'google' ? 'Google' : 'Contraseña' }}
              </span>
            </td>

            <!-- Date -->
            <td class="td-date">{{ formatDate(user.created_at) }}</td>

            <!-- Actions -->
            <td class="td-actions">
              <button class="action-btn action-btn--edit" title="Editar" @click="openEdit(user)">
                <IconEdit :size="15" :stroke-width="2" />
              </button>
              <button
                class="action-btn action-btn--delete"
                title="Eliminar"
                :disabled="user.is_system_user || deletingId === user.id"
                @click="handleDelete(user)"
              >
                <IconTrash v-if="deletingId !== user.id" :size="15" :stroke-width="2" />
                <div v-else class="spinner spinner--sm" />
              </button>
            </td>
          </tr>
        </tbody>
      </table>
    </div>

    <!-- Modal -->
    <Teleport to="body">
      <div v-if="showModal" class="modal-backdrop" @click.self="closeModal">
        <div class="modal">
          <div class="modal-header">
            <h2 class="modal-title">
              {{ modalMode === 'create' ? 'Nuevo usuario' : 'Editar usuario' }}
            </h2>
            <button class="modal-close" @click="closeModal">
              <IconX :size="18" :stroke-width="2.5" />
            </button>
          </div>

          <!-- Create form -->
          <form v-if="modalMode === 'create'" class="modal-form" @submit.prevent="handleCreate">
            <div class="field">
              <label class="field-label">Nombre completo</label>
              <input
                v-model="createName"
                v-bind="createNameAttrs"
                class="field-input"
                :class="{ 'field-input--error': createErrors['name'] }"
                placeholder="Ej. María García"
                autocomplete="name"
              />
              <p v-if="createErrors['name']" class="field-error">{{ createErrors['name'] }}</p>
            </div>

            <div class="field">
              <label class="field-label">Correo electrónico</label>
              <input
                v-model="createEmail"
                v-bind="createEmailAttrs"
                class="field-input"
                :class="{ 'field-input--error': createErrors['email'] }"
                type="email"
                placeholder="correo@ejemplo.com"
                autocomplete="email"
              />
              <p v-if="createErrors['email']" class="field-error">{{ createErrors['email'] }}</p>
            </div>

            <div class="field">
              <label class="field-label">Contraseña</label>
              <input
                v-model="createPassword"
                v-bind="createPasswordAttrs"
                class="field-input"
                :class="{ 'field-input--error': createErrors['password'] }"
                type="password"
                placeholder="Mínimo 8 caracteres"
                autocomplete="new-password"
              />
              <p v-if="createErrors['password']" class="field-error">
                {{ createErrors['password'] }}
              </p>
            </div>

            <p v-if="submitError" class="submit-error">{{ submitError }}</p>

            <div class="modal-actions">
              <button type="button" class="btn-secondary" @click="closeModal">Cancelar</button>
              <button type="submit" class="btn-primary" :disabled="submitting">
                <div v-if="submitting" class="spinner spinner--sm spinner--white" />
                <span v-else>Crear usuario</span>
              </button>
            </div>
          </form>

          <!-- Edit form -->
          <form v-else class="modal-form" @submit.prevent="handleEdit">
            <div class="field">
              <label class="field-label">Nombre completo</label>
              <input
                v-model="editName"
                v-bind="editNameAttrs"
                class="field-input"
                :class="{ 'field-input--error': editErrors['name'] }"
                placeholder="Ej. María García"
                autocomplete="name"
              />
              <p v-if="editErrors['name']" class="field-error">{{ editErrors['name'] }}</p>
            </div>

            <div class="field">
              <label class="field-label">Correo electrónico</label>
              <input
                v-model="editEmail"
                v-bind="editEmailAttrs"
                class="field-input"
                :class="{ 'field-input--error': editErrors['email'] }"
                type="email"
                placeholder="correo@ejemplo.com"
                autocomplete="email"
              />
              <p v-if="editErrors['email']" class="field-error">{{ editErrors['email'] }}</p>
            </div>

            <p v-if="submitError" class="submit-error">{{ submitError }}</p>

            <div class="modal-actions">
              <button type="button" class="btn-secondary" @click="closeModal">Cancelar</button>
              <button type="submit" class="btn-primary" :disabled="submitting">
                <div v-if="submitting" class="spinner spinner--sm spinner--white" />
                <span v-else>Guardar cambios</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
/* ── Layout ────────────────────────────────────────────── */
.users-view {
  width: 100%;
  padding: var(--space-8) var(--space-10);
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
}

/* ── Header ────────────────────────────────────────────── */
.page-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: var(--space-4);
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
}

.users-table {
  width: 100%;
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

.th-actions {
  text-align: right;
}

.user-row {
  border-bottom: 1px solid var(--color-border-light);
  transition: background var(--transition-fast);
}

.user-row:last-child {
  border-bottom: none;
}

.user-row:hover {
  background: #f9fafb;
}

.users-table td {
  padding: var(--space-3) var(--space-4);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  vertical-align: middle;
}

/* ── Avatar + user cell ────────────────────────────────── */
.td-user {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.avatar {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-xs);
  font-weight: 700;
  letter-spacing: 0.02em;
  flex-shrink: 0;
}

.avatar--admin {
  background: #fef3c7;
  color: #d97706;
}

.avatar--regular {
  background: var(--color-accent-light);
  color: var(--color-accent);
}

.user-info {
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
}

.user-name {
  font-weight: 600;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-email {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

/* ── Badge ─────────────────────────────────────────────── */
.badge {
  display: inline-flex;
  align-items: center;
  padding: 2px var(--space-2);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: 600;
}

.badge--admin {
  background: #fef3c7;
  color: #b45309;
}

.badge--user {
  background: #f0fdf4;
  color: #15803d;
}

/* ── Provider ──────────────────────────────────────────── */
.provider {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  font-size: var(--text-xs);
  font-weight: 500;
  color: var(--color-text-tertiary);
}

/* ── Date ──────────────────────────────────────────────── */
.td-date {
  white-space: nowrap;
  font-size: var(--text-xs) !important;
  color: var(--color-text-tertiary) !important;
}

/* ── Action buttons ────────────────────────────────────── */
.td-actions {
  text-align: right;
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
  transition: background var(--transition-fast), border-color var(--transition-fast),
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

.spinner--sm {
  width: 16px;
  height: 16px;
  border-width: 2px;
}

.spinner--white {
  border-color: rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

/* ── Buttons ───────────────────────────────────────────── */
.btn-primary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-5);
  background: var(--color-accent);
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast);
  min-width: 120px;
}

.btn-primary:hover:not(:disabled) {
  background: #369a6e;
}

.btn-primary:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

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

/* ── Modal ─────────────────────────────────────────────── */
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: var(--space-4);
  animation: fade-in 0.15s ease;
}

@keyframes fade-in {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.modal {
  background: #fff;
  border-radius: var(--radius-lg);
  box-shadow:
    0 20px 60px rgba(0, 0, 0, 0.15),
    0 4px 16px rgba(0, 0, 0, 0.08);
  width: 100%;
  max-width: 440px;
  animation: slide-up 0.2s ease;
}

@keyframes slide-up {
  from {
    opacity: 0;
    transform: translateY(12px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-5) var(--space-6) var(--space-4);
  border-bottom: 1px solid var(--color-border-light);
}

.modal-title {
  font-size: var(--text-lg);
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0;
}

.modal-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border: none;
  background: transparent;
  border-radius: var(--radius-md);
  color: var(--color-text-tertiary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.modal-close:hover {
  background: #f3f4f6;
  color: var(--color-text-primary);
}

.modal-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  padding: var(--space-5) var(--space-6) var(--space-6);
}

/* ── Form fields ───────────────────────────────────────── */
.field {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.field-label {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-secondary);
}

.field-input {
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  background: #fff;
  transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
}

.field-input:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(66, 184, 131, 0.15);
}

.field-input--error {
  border-color: var(--color-error);
}

.field-input--error:focus {
  box-shadow: 0 0 0 3px rgba(220, 38, 38, 0.12);
}

.field-error {
  font-size: var(--text-xs);
  color: var(--color-error);
  margin: 0;
}

.submit-error {
  font-size: var(--text-sm);
  color: var(--color-error);
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: var(--radius-md);
  padding: var(--space-2) var(--space-3);
  margin: 0;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  padding-top: var(--space-2);
}
</style>
