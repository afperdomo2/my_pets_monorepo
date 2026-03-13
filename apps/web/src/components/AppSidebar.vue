<script setup lang="ts">
import { computed, watch } from 'vue'
import { RouterLink, useRoute, useRouter } from 'vue-router'
import {
  IconDroplet,
  IconHome,
  IconPaw,
  IconVaccine,
  IconChartBar,
  IconUsers,
  IconSettings,
  IconLogout,
  IconListCheck,
} from '@tabler/icons-vue'
import { useAuthStore } from '@/stores/auth'

const props = defineProps<{ open?: boolean }>()
const emit = defineEmits<{ close: [] }>()

const route = useRoute()
const router = useRouter()
const authStore = useAuthStore()

// Close drawer when route changes (mobile nav)
watch(
  () => route.path,
  () => {
    if (props.open) emit('close')
  },
)

async function handleLogout() {
  await authStore.logout()
  router.push('/login')
}

const userInitials = computed(() => {
  const name = authStore.user?.name ?? ''
  return name
    .split(' ')
    .slice(0, 2)
    .map((w) => w[0])
    .join('')
    .toUpperCase()
})

const userRole = computed(() =>
  authStore.user?.is_system_user ? 'Administrador' : 'Usuario',
)

const navItems = [
  { to: '/',         name: 'home',     label: 'Inicio',    icon: IconHome },
  { to: '/pets',     name: 'pets',     label: 'Mascotas',  icon: IconPaw },
  { to: '/vaccines', name: 'vaccines', label: 'Vacunas',   icon: IconVaccine },
  { to: '/reports',  name: 'reports',  label: 'Reportes',  icon: IconChartBar },
]

const adminNavItems = [
  { to: '/users', name: 'users', label: 'Usuarios', icon: IconUsers },
  { to: '/vaccines-catalog', name: 'vaccines-catalog', label: 'Catálogo de vacunas', icon: IconListCheck },
]

const bottomItems = [
  { to: '/settings', name: 'settings', label: 'Configuración', icon: IconSettings },
]

function isActive(itemName: string): boolean {
  return route.name === itemName
}
</script>

<template>
  <!-- Mobile overlay -->
  <Transition name="overlay">
    <div v-if="open" class="sidebar-overlay" aria-hidden="true" @click="emit('close')" />
  </Transition>

  <aside class="sidebar" :class="{ 'sidebar--open': open }">
    <!-- Brand -->
    <div class="sidebar-brand">
      <div class="brand-icon">
        <IconDroplet :size="22" :stroke-width="2" />
      </div>
      <div class="brand-text">
        <span class="brand-name">My Pets</span>
        <span class="brand-tagline">Bienestar animal</span>
      </div>
    </div>

    <!-- Divider -->
    <div class="sidebar-divider" />

    <!-- Main navigation -->
    <nav class="sidebar-nav">
      <p class="nav-section-label">Principal</p>
      <RouterLink
        v-for="item in navItems"
        :key="item.name"
        :to="item.to"
        class="nav-item"
        :class="{ 'nav-item--active': isActive(item.name) }"
      >
        <span class="nav-icon"><component :is="item.icon" :size="20" :stroke-width="1.75" /></span>
        <span class="nav-label">{{ item.label }}</span>
        <span v-if="isActive(item.name)" class="nav-active-dot" />
      </RouterLink>
    </nav>

    <!-- Admin section (system users only) -->
    <nav v-if="authStore.user?.is_system_user" class="sidebar-nav">
      <p class="nav-section-label">Administración</p>
      <RouterLink
        v-for="item in adminNavItems"
        :key="item.name"
        :to="item.to"
        class="nav-item"
        :class="{ 'nav-item--active': isActive(item.name) }"
      >
        <span class="nav-icon"><component :is="item.icon" :size="20" :stroke-width="1.75" /></span>
        <span class="nav-label">{{ item.label }}</span>
        <span v-if="isActive(item.name)" class="nav-active-dot" />
      </RouterLink>
    </nav>

    <!-- Spacer -->
    <div class="sidebar-spacer" />

    <!-- Bottom: settings + user -->
    <div class="sidebar-bottom">
      <div class="sidebar-divider" />
      <nav class="sidebar-nav sidebar-nav--bottom">
        <RouterLink
          v-for="item in bottomItems"
          :key="item.name"
          :to="item.to"
          class="nav-item"
          :class="{ 'nav-item--active': isActive(item.name) }"
        >
          <span class="nav-icon"><component :is="item.icon" :size="20" :stroke-width="1.75" /></span>
          <span class="nav-label">{{ item.label }}</span>
        </RouterLink>
      </nav>

      <!-- User profile -->
      <RouterLink to="/profile" class="sidebar-user" :class="{ 'sidebar-user--active': isActive('profile') }">
        <div class="user-avatar">{{ userInitials }}</div>
        <div class="user-info">
          <span class="user-name">{{ authStore.user?.name ?? '—' }}</span>
          <span class="user-role">{{ userRole }}</span>
        </div>
        <button class="user-logout" title="Cerrar sesión" @click.prevent="handleLogout">
          <IconLogout :size="16" :stroke-width="1.75" />
        </button>
      </RouterLink>
    </div>
  </aside>
</template>

<style scoped>
/* ── Sidebar shell ────────────────────────────────────────── */
.sidebar {
  width: var(--sidebar-width);
  min-height: 100vh;
  background: var(--sidebar-bg);
  border-right: 1px solid var(--sidebar-border);
  display: flex;
  flex-direction: column;
  padding: var(--space-5) 0;
  position: fixed;
  top: 0;
  left: 0;
  bottom: 0;
  z-index: 100;
}

/* ── Brand ────────────────────────────────────────────────── */
.sidebar-brand {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: 0 var(--space-5) var(--space-4);
}

.brand-icon {
  width: 40px;
  height: 40px;
  background: var(--color-accent-light);
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-accent);
  flex-shrink: 0;
}

.brand-text {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.brand-name {
  font-family: var(--font-display);
  font-size: var(--text-lg);
  font-weight: 600;
  color: var(--color-text-primary);
  line-height: 1.2;
}

.brand-tagline {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  letter-spacing: 0.02em;
  font-weight: 400;
}

/* ── Divider ──────────────────────────────────────────────── */
.sidebar-divider {
  height: 1px;
  background: var(--color-border-light);
  margin: var(--space-2) var(--space-4);
}

/* ── Nav ──────────────────────────────────────────────────── */
.sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: 2px;
  padding: var(--space-3) var(--space-3);
}

.nav-section-label {
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-text-tertiary);
  letter-spacing: 0.08em;
  text-transform: uppercase;
  padding: var(--space-2) var(--space-2) var(--space-1);
}

.nav-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-2) var(--space-3);
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
  font-weight: 500;
  text-decoration: none;
  transition: background var(--transition-fast), color var(--transition-fast);
  position: relative;
}

.nav-item:hover {
  background: var(--sidebar-item-hover-bg);
  color: var(--color-text-primary);
}

.nav-item--active {
  background: var(--sidebar-item-active-bg);
  color: var(--sidebar-item-active-text);
  font-weight: 600;
}

.nav-item--active:hover {
  background: var(--sidebar-item-active-bg);
  color: var(--sidebar-item-active-text);
}

.nav-icon {
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  color: inherit;
}

.nav-icon :deep(svg) {
  display: block;
}

.nav-label {
  flex: 1;
}

.nav-active-dot {
  width: 6px;
  height: 6px;
  border-radius: var(--radius-full);
  background: var(--color-accent);
  flex-shrink: 0;
}

/* ── Spacer ───────────────────────────────────────────────── */
.sidebar-spacer {
  flex: 1;
}

/* ── Bottom ───────────────────────────────────────────────── */
.sidebar-bottom {
  display: flex;
  flex-direction: column;
  gap: 0;
}

.sidebar-nav--bottom {
  padding-bottom: var(--space-2);
}

/* ── User card ────────────────────────────────────────────── */
.sidebar-user {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-4);
  margin: 0 var(--space-3);
  border-radius: var(--radius-md);
  transition: background var(--transition-fast);
  text-decoration: none;
  cursor: pointer;
}

.sidebar-user:hover {
  background: var(--sidebar-item-hover-bg);
}

.sidebar-user--active {
  background: var(--sidebar-item-active-bg);
}

.sidebar-user--active .user-name {
  color: var(--sidebar-item-active-text);
}

.sidebar-user--active .user-role {
  color: var(--color-accent-muted);
}

.user-avatar {
  width: 34px;
  height: 34px;
  border-radius: var(--radius-full);
  background: var(--color-accent);
  color: var(--color-text-inverse);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-xs);
  font-weight: 600;
  letter-spacing: 0.02em;
  flex-shrink: 0;
}

.user-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
}

.user-name {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-role {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
}

.user-logout {
  background: none;
  border: none;
  color: var(--color-text-tertiary);
  padding: var(--space-1);
  border-radius: var(--radius-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: color var(--transition-fast), background var(--transition-fast);
  cursor: pointer;
  flex-shrink: 0;
}

.user-logout:hover {
  color: var(--color-error);
  background: #FEF2F2;
}

/* ── Mobile: drawer behavior ──────────────────────────────── */
.sidebar-overlay {
  display: none;
}

@media (max-width: 768px) {
  .sidebar {
    transform: translateX(-100%);
    transition: transform var(--transition-base);
    box-shadow: none;
    top: 0;
  }

  .sidebar--open {
    transform: translateX(0);
    box-shadow: var(--shadow-xl);
  }

  .sidebar-overlay {
    display: block;
    position: fixed;
    inset: 0;
    background: rgba(26, 26, 24, 0.45);
    z-index: 99;
    backdrop-filter: blur(2px);
  }
}

/* ── Overlay transition ───────────────────────────────────── */
.overlay-enter-active,
.overlay-leave-active {
  transition: opacity var(--transition-base);
}
.overlay-enter-from,
.overlay-leave-to {
  opacity: 0;
}
</style>
