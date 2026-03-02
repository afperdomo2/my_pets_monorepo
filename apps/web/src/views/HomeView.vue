<script setup lang="ts">
import { RouterLink } from 'vue-router'

const stats = [
  {
    label: 'Mascotas registradas',
    value: '12',
    change: '+2 este mes',
    positive: true,
    icon: `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><circle cx="11" cy="4" r="2"/><circle cx="18" cy="8" r="2"/><circle cx="4" cy="8" r="2"/><path d="M11 6c0 4-3 6-3 10a4 4 0 0 0 8 0c0-4-3-6-3-10"/><path d="M6.5 14c-1.5 1-2.5 2.5-2.5 4"/><path d="M17.5 14c1.5 1 2.5 2.5 2.5 4"/></svg>`,
    color: 'accent',
    to: '/pets',
  },
  {
    label: 'Vacunas al día',
    value: '9',
    change: '75% del total',
    positive: true,
    icon: `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="m18 2 4 4"/><path d="m17 7 3-3"/><path d="M19 9 8.7 19.3c-1 1-2.5 1-3.4 0l-1.6-1.6c-1-1-1-2.5 0-3.4L14 4"/><path d="m9 8 2 2"/><path d="m13 12 2 2"/></svg>`,
    color: 'green',
    to: '/vaccines',
  },
  {
    label: 'Vacunas pendientes',
    value: '3',
    change: 'Requieren atención',
    positive: false,
    icon: `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M10.29 3.86 1.82 18a2 2 0 0 0 1.71 3h16.94a2 2 0 0 0 1.71-3L13.71 3.86a2 2 0 0 0-3.42 0z"/><line x1="12" y1="9" x2="12" y2="13"/><line x1="12" y1="17" x2="12.01" y2="17"/></svg>`,
    color: 'orange',
    to: '/vaccines',
  },
  {
    label: 'Propietarios activos',
    value: '8',
    change: 'Registrados',
    positive: true,
    icon: `<svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round"><path d="M17 21v-2a4 4 0 0 0-4-4H5a4 4 0 0 0-4 4v2"/><circle cx="9" cy="7" r="4"/><path d="M23 21v-2a4 4 0 0 0-3-3.87"/><path d="M16 3.13a4 4 0 0 1 0 7.75"/></svg>`,
    color: 'blue',
    to: '/pets',
  },
]

const recentPets = [
  { name: 'Luna', species: 'Perro', breed: 'Golden Retriever', owner: 'Ana M.', status: 'Saludable', initials: 'LU' },
  { name: 'Simba', species: 'Gato', breed: 'Siamés', owner: 'Carlos R.', status: 'Vacuna pendiente', initials: 'SI' },
  { name: 'Manchas', species: 'Perro', breed: 'Dálmata', owner: 'Laura G.', status: 'Saludable', initials: 'MA' },
  { name: 'Kira', species: 'Gato', breed: 'Persa', owner: 'Pedro L.', status: 'En revisión', initials: 'KI' },
]

const upcomingVaccines = [
  { pet: 'Simba', vaccine: 'Antirrábica', due: 'En 3 días', urgent: true },
  { pet: 'Manchas', vaccine: 'Polivalente', due: 'En 8 días', urgent: false },
  { pet: 'Bolt', vaccine: 'Bivalente felina', due: 'En 15 días', urgent: false },
]
</script>

<template>
  <div class="home-view">
    <!-- Page header -->
    <div class="page-header">
      <div class="page-header-text">
        <h1>Buenos días, Juan</h1>
        <p>Aquí está el resumen de hoy — <span class="date">Lunes, 2 de marzo de 2026</span></p>
      </div>
      <RouterLink to="/pets" class="btn-new-pet">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <line x1="12" y1="5" x2="12" y2="19"/><line x1="5" y1="12" x2="19" y2="12"/>
        </svg>
        Nueva mascota
      </RouterLink>
    </div>

    <!-- Stats grid -->
    <div class="stats-grid">
      <RouterLink
        v-for="stat in stats"
        :key="stat.label"
        :to="stat.to"
        class="stat-card"
        :class="`stat-card--${stat.color}`"
      >
        <div class="stat-icon" v-html="stat.icon" />
        <div class="stat-body">
          <span class="stat-value">{{ stat.value }}</span>
          <span class="stat-label">{{ stat.label }}</span>
        </div>
        <span class="stat-change" :class="{ 'stat-change--negative': !stat.positive }">
          {{ stat.change }}
        </span>
      </RouterLink>
    </div>

    <!-- Content grid -->
    <div class="content-grid">
      <!-- Mascotas recientes -->
      <div class="panel">
        <div class="panel-header">
          <h2>Mascotas recientes</h2>
          <RouterLink to="/pets" class="panel-link">Ver todas</RouterLink>
        </div>
        <div class="panel-body">
          <div
            v-for="pet in recentPets"
            :key="pet.name"
            class="pet-row"
          >
            <div class="pet-avatar">{{ pet.initials }}</div>
            <div class="pet-info">
              <span class="pet-name">{{ pet.name }}</span>
              <span class="pet-meta">{{ pet.species }} · {{ pet.breed }} · {{ pet.owner }}</span>
            </div>
            <span
              class="pet-status"
              :class="{
                'pet-status--ok': pet.status === 'Saludable',
                'pet-status--warn': pet.status === 'Vacuna pendiente',
                'pet-status--info': pet.status === 'En revisión',
              }"
            >{{ pet.status }}</span>
          </div>
        </div>
      </div>

      <!-- Vacunas próximas -->
      <div class="panel">
        <div class="panel-header">
          <h2>Vacunas próximas</h2>
          <RouterLink to="/vaccines" class="panel-link">Ver calendario</RouterLink>
        </div>
        <div class="panel-body">
          <div
            v-for="vax in upcomingVaccines"
            :key="vax.pet + vax.vaccine"
            class="vax-row"
          >
            <div class="vax-indicator" :class="{ 'vax-indicator--urgent': vax.urgent }" />
            <div class="vax-info">
              <span class="vax-name">{{ vax.vaccine }}</span>
              <span class="vax-pet">Para <strong>{{ vax.pet }}</strong></span>
            </div>
            <span class="vax-due" :class="{ 'vax-due--urgent': vax.urgent }">{{ vax.due }}</span>
          </div>
        </div>

        <!-- CTA vacunas -->
        <div class="panel-footer">
          <RouterLink to="/vaccines" class="btn-secondary">
            Gestionar vacunas
          </RouterLink>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.home-view {
  padding: var(--space-8) var(--space-10);
  display: flex;
  flex-direction: column;
  gap: var(--space-8);
  width: 100%;
}

/* ── Page header ──────────────────────────────────────────── */
.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--space-4);
}

.page-header-text h1 {
  font-family: var(--font-display);
  font-size: var(--text-3xl);
  color: var(--color-text-primary);
  margin-bottom: var(--space-1);
}

.page-header-text p {
  font-size: var(--text-base);
  color: var(--color-text-secondary);
}

.date {
  color: var(--color-text-tertiary);
  font-size: var(--text-sm);
}

.btn-new-pet {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-5);
  background: var(--color-accent);
  color: #fff;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  text-decoration: none;
  transition: background var(--transition-fast), box-shadow var(--transition-fast);
  white-space: nowrap;
  flex-shrink: 0;
}

.btn-new-pet:hover {
  background: var(--color-accent-hover);
  box-shadow: var(--shadow-md);
  color: #fff;
}

/* ── Stats ────────────────────────────────────────────────── */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--space-4);
}

@media (max-width: 1200px) {
  .stats-grid { grid-template-columns: repeat(2, 1fr); }
}

.stat-card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--space-5);
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  text-decoration: none;
  transition: box-shadow var(--transition-base), transform var(--transition-base), border-color var(--transition-base);
  cursor: pointer;
}

.stat-card:hover {
  box-shadow: var(--shadow-md);
  transform: translateY(-2px);
  border-color: var(--color-border-light);
}

.stat-icon {
  width: 44px;
  height: 44px;
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
}

.stat-card--accent .stat-icon { background: var(--color-accent-light); color: var(--color-accent); }
.stat-card--green  .stat-icon { background: #E8F5EE; color: #2E7D52; }
.stat-card--orange .stat-icon { background: #FEF3E2; color: #C4714A; }
.stat-card--blue   .stat-icon { background: #EBF4FB; color: #2980B9; }

.stat-body {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.stat-value {
  font-family: var(--font-display);
  font-size: var(--text-3xl);
  font-weight: 600;
  color: var(--color-text-primary);
  line-height: 1;
}

.stat-label {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  font-weight: 400;
}

.stat-change {
  font-size: var(--text-xs);
  font-weight: 500;
  color: var(--color-accent);
  background: var(--color-accent-light);
  padding: 2px var(--space-2);
  border-radius: var(--radius-full);
  align-self: flex-start;
}

.stat-change--negative {
  color: var(--color-secondary);
  background: var(--color-secondary-light);
}

/* ── Content grid ─────────────────────────────────────────── */
.content-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-6);
}

@media (max-width: 900px) {
  .content-grid { grid-template-columns: 1fr; }
}

/* ── Panels ───────────────────────────────────────────────── */
.panel {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-5) var(--space-6);
  border-bottom: 1px solid var(--color-border-light);
}

.panel-header h2 {
  font-family: var(--font-display);
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-text-primary);
}

.panel-link {
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-accent);
  text-decoration: none;
}

.panel-link:hover {
  color: var(--color-accent-dark);
}

.panel-body {
  padding: var(--space-2) 0;
  flex: 1;
}

.panel-footer {
  padding: var(--space-4) var(--space-6);
  border-top: 1px solid var(--color-border-light);
}

/* ── Pet rows ─────────────────────────────────────────────── */
.pet-row {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-6);
  transition: background var(--transition-fast);
}

.pet-row:hover {
  background: var(--color-bg-alt);
}

.pet-avatar {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-full);
  background: var(--color-accent-light);
  color: var(--color-accent-dark);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-xs);
  font-weight: 700;
  letter-spacing: 0.02em;
  flex-shrink: 0;
}

.pet-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.pet-name {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
}

.pet-meta {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.pet-status {
  font-size: var(--text-xs);
  font-weight: 500;
  padding: 3px var(--space-2);
  border-radius: var(--radius-full);
  white-space: nowrap;
  flex-shrink: 0;
}

.pet-status--ok   { background: #E8F5EE; color: #2E7D52; }
.pet-status--warn { background: #FEF3E2; color: #C4714A; }
.pet-status--info { background: #EBF4FB; color: #2980B9; }

/* ── Vaccine rows ─────────────────────────────────────────── */
.vax-row {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-4) var(--space-6);
  transition: background var(--transition-fast);
}

.vax-row:hover {
  background: var(--color-bg-alt);
}

.vax-indicator {
  width: 8px;
  height: 8px;
  border-radius: var(--radius-full);
  background: var(--color-accent-muted);
  flex-shrink: 0;
}

.vax-indicator--urgent {
  background: var(--color-secondary);
  box-shadow: 0 0 0 3px var(--color-secondary-light);
}

.vax-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.vax-name {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
}

.vax-pet {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
}

.vax-due {
  font-size: var(--text-xs);
  font-weight: 500;
  color: var(--color-text-secondary);
}

.vax-due--urgent {
  color: var(--color-secondary);
  font-weight: 600;
}

/* ── Secondary button ─────────────────────────────────────── */
.btn-secondary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 100%;
  padding: var(--space-2) var(--space-4);
  background: var(--color-surface);
  border: 1.5px solid var(--color-border);
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
  font-weight: 500;
  border-radius: var(--radius-md);
  text-decoration: none;
  transition: background var(--transition-fast), border-color var(--transition-fast), color var(--transition-fast);
}

.btn-secondary:hover {
  background: var(--color-accent-light);
  border-color: var(--color-accent-muted);
  color: var(--color-accent-dark);
}
</style>
