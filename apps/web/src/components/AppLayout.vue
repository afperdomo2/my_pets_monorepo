<script setup lang="ts">
import { ref } from 'vue'
import { RouterView } from 'vue-router'
import { IconDroplet, IconMenu2 } from '@tabler/icons-vue'
import AppSidebar from './AppSidebar.vue'

const sidebarOpen = ref(false)

function openSidebar() {
  sidebarOpen.value = true
}

function closeSidebar() {
  sidebarOpen.value = false
}
</script>

<template>
  <div class="app-layout">
    <!-- Mobile topbar -->
    <header class="mobile-topbar">
      <button class="hamburger" :aria-expanded="sidebarOpen" aria-label="Abrir menú" @click="openSidebar">
        <IconMenu2 :size="22" :stroke-width="2" />
      </button>
      <div class="topbar-brand">
        <IconDroplet :size="18" :stroke-width="2" class="topbar-brand-icon" />
        <span class="topbar-brand-name">My Pets</span>
      </div>
    </header>

    <!-- Sidebar -->
    <AppSidebar :open="sidebarOpen" @close="closeSidebar" />

    <!-- Main content -->
    <main class="app-main">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>
.app-layout {
  display: flex;
  min-height: 100vh;
  background: var(--color-bg);
}

.app-main {
  flex: 1;
  margin-left: var(--sidebar-width);
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  min-width: 0;
  overflow-x: hidden;
}

/* ── Mobile topbar ─────────────────────── */
.mobile-topbar {
  display: none;
}

@media (max-width: 768px) {
  .mobile-topbar {
    display: flex;
    align-items: center;
    gap: var(--space-3);
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    height: 56px;
    background: var(--color-surface);
    border-bottom: 1px solid var(--color-border);
    padding: 0 var(--space-4);
    z-index: 90;
    box-shadow: var(--shadow-xs);
  }

  .app-main {
    margin-left: 0;
    padding-top: 56px;
  }
}

/* ── Hamburger button ──────────────────── */
.hamburger {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 38px;
  height: 38px;
  border: none;
  border-radius: var(--radius-md);
  background: transparent;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
  flex-shrink: 0;
}

.hamburger:hover {
  background: var(--color-bg-alt);
  color: var(--color-text-primary);
}

/* ── Topbar brand ──────────────────────── */
.topbar-brand {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.topbar-brand-icon {
  color: var(--color-accent);
}

.topbar-brand-name {
  font-family: var(--font-display);
  font-size: var(--text-lg);
  font-weight: 600;
  color: var(--color-text-primary);
  line-height: 1;
}
</style>
