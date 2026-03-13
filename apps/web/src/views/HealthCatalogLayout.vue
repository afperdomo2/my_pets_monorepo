<script setup lang="ts">
import { useRoute } from 'vue-router'

const route = useRoute()

const TABS = [
  { label: 'Vacunas', to: '/health-catalog/vaccine', name: 'vaccine' },
  { label: 'Desparasitación', to: '/health-catalog/deworming', name: 'deworming' },
  { label: 'Exámenes', to: '/health-catalog/exam', name: 'exam' },
]
</script>

<template>
  <div class="health-catalog-layout">
    <!-- Header General -->
    <div class="page-header">
      <div class="page-header__text">
        <h1 class="page-title">Guía de salud</h1>
        <p class="page-subtitle">Administra los esquemas preventivos de las mascotas</p>
      </div>
    </div>

    <!-- Tabs -->
    <div class="tabs-container">
      <nav class="tabs-nav">
        <router-link
          v-for="tab in TABS"
          :key="tab.name"
          :to="tab.to"
          class="tab-link"
          :class="{ 'tab-link--active': route.params.category === tab.name || route.path.includes(tab.to) }"
        >
          {{ tab.label }}
        </router-link>
      </nav>
    </div>

    <!-- Contenido dinámico (la tabla de HealthCatalogView) -->
    <div class="tab-content">
      <router-view />
    </div>
  </div>
</template>

<style scoped>
.health-catalog-layout {
  width: 100%;
  max-width: 100%;
  padding: var(--space-8) var(--space-10);
  display: flex;
  flex-direction: column;
  gap: var(--space-6);
  box-sizing: border-box;
}

@media (max-width: 1024px) {
  .health-catalog-layout {
    padding: var(--space-6) var(--space-6);
  }
}

@media (max-width: 768px) {
  .health-catalog-layout {
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

/* ── Tabs ──────────────────────────────────────────────── */
.tabs-container {
  border-bottom: 1px solid var(--color-border-light);
  margin-bottom: calc(var(--space-2) * -1); /* Pull content slightly closer */
}

.tabs-nav {
  display: flex;
  gap: var(--space-6);
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: none;
}

.tabs-nav::-webkit-scrollbar {
  display: none;
}

.tab-link {
  padding: var(--space-2) 0 var(--space-3);
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-text-tertiary);
  text-decoration: none;
  position: relative;
  transition: color var(--transition-fast);
  white-space: nowrap;
}

.tab-link:hover {
  color: var(--color-text-secondary);
}

.tab-link--active {
  color: var(--color-accent);
  font-weight: 600;
}

.tab-link--active::after {
  content: '';
  position: absolute;
  bottom: -1px;
  left: 0;
  right: 0;
  height: 2px;
  background: var(--color-accent);
  border-radius: 2px 2px 0 0;
}

/* ── Content ───────────────────────────────────────────── */
.tab-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  min-height: 0; /* Important for deeply nested scrolling */
}
</style>
