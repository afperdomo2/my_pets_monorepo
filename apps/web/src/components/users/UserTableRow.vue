<script setup lang="ts">
import { IconEdit, IconTrash, IconLock } from '@tabler/icons-vue'
import type { User } from '@/types/user'

defineProps<{
  user: User
  deletingId: string | null
}>()

const emit = defineEmits<{
  edit: [user: User]
  delete: [user: User]
}>()

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
</script>

<template>
  <tr class="user-row">
    <!-- Avatar + name + email -->
    <td class="td-user">
      <div class="avatar" :class="user.is_system_user ? 'avatar--admin' : 'avatar--regular'">
        {{ initials(user.name) }}
      </div>
      <div class="user-info">
        <span class="user-name">{{ user.name }}</span>
        <span class="user-email">{{ user.email }}</span>
      </div>
    </td>

    <!-- Role badge -->
    <td class="td-center">
      <span class="badge" :class="user.is_system_user ? 'badge--admin' : 'badge--user'">
        {{ user.is_system_user ? 'Administrador' : 'Usuario' }}
      </span>
    </td>

    <!-- Auth provider -->
    <td class="td-center">
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
    <td class="td-center td-date">{{ formatDate(user.created_at) }}</td>

    <!-- Actions -->
    <td class="td-center td-actions">
      <button class="action-btn action-btn--edit" title="Editar" @click="emit('edit', user)">
        <IconEdit :size="15" :stroke-width="2" />
      </button>
      <button
        class="action-btn action-btn--delete"
        title="Eliminar"
        :disabled="user.is_system_user || deletingId === user.id"
        @click="emit('delete', user)"
      >
        <IconTrash v-if="deletingId !== user.id" :size="15" :stroke-width="2" />
        <div v-else class="spinner spinner--sm" />
      </button>
    </td>
  </tr>
</template>

<style scoped>
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

.user-row td {
  padding: var(--space-3) var(--space-4);
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  vertical-align: middle;
}

.td-center {
  text-align: center;
}

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

.avatar--admin { background: #fef3c7; color: #d97706; }
.avatar--regular { background: var(--color-accent-light); color: var(--color-accent); }

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

.badge {
  display: inline-flex;
  align-items: center;
  padding: 2px var(--space-2);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: 600;
}

.badge--admin { background: #fef3c7; color: #b45309; }
.badge--user { background: #f0fdf4; color: #15803d; }

.provider {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  font-size: var(--text-xs);
  font-weight: 500;
  color: var(--color-text-tertiary);
}

.td-date {
  white-space: nowrap;
  font-size: var(--text-xs) !important;
  color: var(--color-text-tertiary) !important;
}

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

.action-btn--edit { color: var(--color-text-tertiary); }
.action-btn--edit:hover {
  background: #eff6ff;
  border-color: #bfdbfe;
  color: #2563eb;
}

.action-btn--delete { color: var(--color-text-tertiary); margin-left: var(--space-1); }
.action-btn--delete:hover:not(:disabled) {
  background: #fef2f2;
  border-color: #fecaca;
  color: var(--color-error);
}
.action-btn--delete:disabled { opacity: 0.3; cursor: not-allowed; }

.spinner {
  width: 28px;
  height: 28px;
  border: 3px solid var(--color-border-light);
  border-top-color: var(--color-accent);
  border-radius: var(--radius-full);
  animation: spin 0.7s linear infinite;
}

.spinner--sm { width: 16px; height: 16px; border-width: 2px; }

@keyframes spin { to { transform: rotate(360deg); } }
</style>
