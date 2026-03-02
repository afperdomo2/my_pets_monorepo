<script setup lang="ts">
import { RouterLink, useRoute, useRouter } from 'vue-router'

const route = useRoute()
const router = useRouter()

function handleLogout() {
  router.push('/login')
}

const navItems = [
  {
    to: '/',
    name: 'home',
    label: 'Inicio',
    icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M3 9.5L12 3l9 6.5V20a1 1 0 0 1-1 1H4a1 1 0 0 1-1-1V9.5z"/><path d="M9 21V12h6v9"/></svg>`,
  },
  {
    to: '/pets',
    name: 'pets',
    label: 'Mascotas',
    icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="4" r="2"/><circle cx="18" cy="8" r="2"/><circle cx="4" cy="8" r="2"/><path d="M11 6c0 4-3 6-3 10a4 4 0 0 0 8 0c0-4-3-6-3-10"/><path d="M6.5 14c-1.5 1-2.5 2.5-2.5 4"/><path d="M17.5 14c1.5 1 2.5 2.5 2.5 4"/></svg>`,
  },
  {
    to: '/vaccines',
    name: 'vaccines',
    label: 'Vacunas',
    icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="m18 2 4 4"/><path d="m17 7 3-3"/><path d="M19 9 8.7 19.3c-1 1-2.5 1-3.4 0l-1.6-1.6c-1-1-1-2.5 0-3.4L14 4"/><path d="m9 8 2 2"/><path d="m13 12 2 2"/><path d="m8 16 2 2"/><path d="m14 8-3 3"/><path d="M5 20 4 21"/></svg>`,
  },
  {
    to: '/reports',
    name: 'reports',
    label: 'Reportes',
    icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M3 3v18h18"/><path d="M18 17V9"/><path d="M13 17V5"/><path d="M8 17v-3"/></svg>`,
  },
]

const bottomItems = [
  {
    to: '/settings',
    name: 'settings',
    label: 'Configuración',
    icon: `<svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"/><circle cx="12" cy="12" r="3"/></svg>`,
  },
]

function isActive(itemName: string): boolean {
  if (itemName === 'home') return route.name === 'home'
  return (route.name as string)?.startsWith(itemName) ?? false
}
</script>

<template>
  <aside class="sidebar">
    <!-- Brand -->
    <div class="sidebar-brand">
      <div class="brand-icon">
        <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
          <path d="M11 6c0 4-3 6-3 10a4 4 0 0 0 8 0c0-4-3-6-3-10"/>
          <circle cx="11" cy="4" r="1.5"/>
        </svg>
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
        <span class="nav-icon" v-html="item.icon" />
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
          <span class="nav-icon" v-html="item.icon" />
          <span class="nav-label">{{ item.label }}</span>
        </RouterLink>
      </nav>

      <!-- User profile -->
      <div class="sidebar-user">
        <div class="user-avatar">JD</div>
        <div class="user-info">
          <span class="user-name">Juan Doe</span>
          <span class="user-role">Administrador</span>
        </div>
        <button class="user-logout" title="Cerrar sesión" @click="handleLogout">
          <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
            <path d="M9 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h4"/><polyline points="16 17 21 12 16 7"/><line x1="21" y1="12" x2="9" y2="12"/>
          </svg>
        </button>
      </div>
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
}

.sidebar-user:hover {
  background: var(--sidebar-item-hover-bg);
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
</style>
