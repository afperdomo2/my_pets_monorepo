<script setup lang="ts">
import type { Component } from 'vue'

defineProps<{
  label: string
  value: string | number
  description: string
  color: 'green' | 'amber' | 'red'
  icon: Component
}>()
</script>

<template>
  <div class="stat-card" :class="`stat-card--${color}`">
    <div class="stat-card__icon">
      <component :is="icon" :size="24" :stroke-width="1.75" />
    </div>
    <div class="stat-card__body">
      <span class="stat-card__value">{{ value }}</span>
      <span class="stat-card__label">{{ label }}</span>
    </div>
    <span class="stat-card__badge" :class="`stat-card__badge--${color}`">
      {{ description }}
    </span>
  </div>
</template>

<style scoped>
/* container query: el componente responde a su propio tamaño, no al viewport */
.stat-card {
  container-type: inline-size;
  container-name: stat-card;

  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--space-5);
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  transition: box-shadow var(--transition-base), transform var(--transition-base),
    border-color var(--transition-base);
  cursor: default;
}

.stat-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
  border-color: var(--color-border-light);
}

/* ── Icono ─────────────────────────── */
.stat-card__icon {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.stat-card--green .stat-card__icon {
  background: #E8F5EE;
  color: #2E7D52;
}
.stat-card--amber .stat-card__icon {
  background: #FEF3E2;
  color: #C4714A;
}
.stat-card--red .stat-card__icon {
  background: #FEF2F2;
  color: #DC2626;
}

/* ── Body ──────────────────────────── */
.stat-card__body {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.stat-card__value {
  font-family: var(--font-display);
  font-size: var(--text-3xl);
  font-weight: 600;
  color: var(--color-text-primary);
  line-height: 1;
}

.stat-card__label {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  font-weight: 400;
}

/* ── Badge ─────────────────────────── */
.stat-card__badge {
  font-size: var(--text-xs);
  font-weight: 500;
  padding: 2px var(--space-2);
  border-radius: var(--radius-full);
  align-self: flex-start;
}

.stat-card__badge--green {
  background: #E8F5EE;
  color: #2E7D52;
}
.stat-card__badge--amber {
  background: #FEF3E2;
  color: #C4714A;
}
.stat-card__badge--red {
  background: #FEF2F2;
  color: #DC2626;
}

/* ── Container query: layout horizontal cuando la card es ancha ── */
/* Ej: cuando el grid pone 1 columna y la card ocupa todo el ancho */
@container stat-card (min-width: 320px) {
  .stat-card {
    flex-direction: row;
    align-items: center;
    gap: var(--space-4);
    padding: var(--space-4) var(--space-5);
  }

  .stat-card__body {
    flex: 1;
  }

  .stat-card__badge {
    align-self: center;
    white-space: nowrap;
  }
}
</style>
