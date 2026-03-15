<script setup lang="ts">
import { RouterLink } from 'vue-router'
import { IconPaw, IconShieldCheck, IconAlertTriangle, IconPlus, IconStethoscope, IconPill, IconVaccine } from '@tabler/icons-vue'
import PetAvatar from '@/components/pets/PetAvatar.vue'
import { getSpeciesLabel } from '@/constants/species'
import type { Component } from 'vue'

const stats: {
  label: string
  value: string
  change: string
  positive: boolean
  icon: Component
  color: string
  to: string
}[] = [
  {
    label: 'Mascotas registradas',
    value: '12',
    change: '+2 este mes',
    positive: true,
    icon: IconPaw,
    color: 'accent',
    to: '/pets',
  },
  {
    label: 'Mascotas al día',
    value: '9',
    change: '75% con salud al día',
    positive: true,
    icon: IconShieldCheck,
    color: 'green',
    to: '/health-records',
  },
  {
    label: 'Tareas pendientes',
    value: '3',
    change: 'Requieren atención',
    positive: false,
    icon: IconAlertTriangle,
    color: 'orange',
    to: '/health-records',
  },
]

const upcomingHealthTasks = [
  { pet: 'Romeo', species: 'dog', task: 'Desparasitante interno', category: 'deworming', due: 'En 2 días', urgent: false },
  { pet: 'Negra', species: 'dog', task: 'Perfil Senior', category: 'exam', due: 'En 5 días', urgent: false },
  { pet: 'Simba', species: 'cat', task: 'Antirrábica', category: 'vaccine', due: 'En 3 días', urgent: true },
  { pet: 'Manchas', species: 'dog', task: 'Polivalente', category: 'vaccine', due: 'En 8 días', urgent: false },
]

const getCategoryIcon = (category: string) => {
  switch (category) {
    case 'vaccine': return IconVaccine
    case 'deworming': return IconPill
    case 'exam': return IconStethoscope
    default: return IconShieldCheck
  }
}

const getCategoryColor = (category: string) => {
  switch (category) {
    case 'vaccine': return 'var(--color-accent)'
    case 'deworming': return '#c4714a'
    case 'exam': return '#2980b9'
    default: return 'var(--color-accent)'
  }
}

const getCategoryLabel = (category: string) => {
  switch (category) {
    case 'vaccine': return 'Vacuna'
    case 'deworming': return 'Desparasitación'
    case 'exam': return 'Examen'
    default: return 'Salud'
  }
}
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
        <IconPlus :size="16" :stroke-width="2.5" />
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
        <div class="stat-icon"><component :is="stat.icon" :size="22" :stroke-width="1.75" /></div>
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
      <!-- Próximas tareas de salud -->
      <div class="panel">
        <div class="panel-header">
          <h2>Próximas tareas de salud</h2>
          <RouterLink to="/health-records" class="panel-link">Ver todas</RouterLink>
        </div>
        <div class="panel-body">
          <div v-for="task in upcomingHealthTasks" :key="task.pet + task.task" class="task-row">
            <div class="task-indicator" :class="{ 'task-indicator--urgent': task.urgent }" />
            <!-- Mascota -->
            <div class="task-pet-info">
              <PetAvatar :species="task.species" :name="task.pet" size="md" />
              <div class="task-pet-text">
                <span class="task-pet-name">{{ task.pet }}</span>
                <span class="task-pet-species">{{ getSpeciesLabel(task.species) }}</span>
              </div>
            </div>
            <!-- Tarea de salud -->
            <div class="task-icon" :style="{ color: getCategoryColor(task.category) }">
              <component :is="getCategoryIcon(task.category)" :size="16" :stroke-width="1.75" />
            </div>
            <div class="task-info">
              <span class="task-name">{{ task.task }}</span>
              <span class="task-category">{{ getCategoryLabel(task.category) }}</span>
            </div>
            <span class="task-due" :class="{ 'task-due--urgent': task.urgent }">{{ task.due }}</span>
          </div>
        </div>

        <!-- CTA -->
        <div class="panel-footer">
          <RouterLink to="/health-records" class="btn-secondary"> Gestionar salud </RouterLink>
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

@media (max-width: 768px) {
  .home-view {
    padding: var(--space-5) var(--space-4);
    gap: var(--space-5);
  }
}

/* ── Page header ──────────────────────────────────────────── */
.page-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--space-4);
}

@media (max-width: 480px) {
  .page-header {
    flex-direction: column;
    gap: var(--space-3);
  }

  .btn-new-pet {
    width: 100%;
    justify-content: center;
  }
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
  transition:
    background var(--transition-fast),
    box-shadow var(--transition-fast);
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
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-4);
}

@media (max-width: 900px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

@media (max-width: 480px) {
  .stats-grid {
    grid-template-columns: 1fr 1fr;
    gap: var(--space-3);
  }
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
  transition:
    box-shadow var(--transition-base),
    transform var(--transition-base),
    border-color var(--transition-base);
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

.stat-card--accent .stat-icon {
  background: var(--color-accent-light);
  color: var(--color-accent);
}
.stat-card--green .stat-icon {
  background: #e8f5ee;
  color: #2e7d52;
}
.stat-card--orange .stat-icon {
  background: #fef3e2;
  color: #c4714a;
}
.stat-card--blue .stat-icon {
  background: #ebf4fb;
  color: #2980b9;
}

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
  .content-grid {
    grid-template-columns: 1fr;
  }
}

@media (max-width: 768px) {
  .content-grid {
    gap: var(--space-4);
  }
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

/* ── Task rows ────────────────────────────────────────────── */
.task-row {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-6);
  transition: background var(--transition-fast);
}

.task-row:hover {
  background: var(--color-bg-alt);
}

.task-indicator {
  width: 8px;
  height: 8px;
  border-radius: var(--radius-full);
  background: var(--color-accent-muted);
  flex-shrink: 0;
}

.task-indicator--urgent {
  background: var(--color-secondary);
  box-shadow: 0 0 0 3px var(--color-secondary-light);
}

/* Mascota */
.task-pet-info {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  min-width: 110px;
  flex-shrink: 0;
}

.task-pet-text {
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
}

.task-pet-name {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.task-pet-species {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  text-transform: capitalize;
}

/* Icono de tarea */
.task-icon {
  width: 28px;
  height: 28px;
  border-radius: var(--radius-md);
  background: var(--color-bg-alt);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

/* Info de la tarea */
.task-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
  margin-right: var(--space-3);
}

.task-name {
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.task-category {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
}

.task-due {
  font-size: var(--text-xs);
  font-weight: 500;
  color: var(--color-text-secondary);
  flex-shrink: 0;
  white-space: nowrap;
}

.task-due--urgent {
  color: var(--color-secondary);
  font-weight: 600;
}

/* ── Responsive: más espacio en pantallas grandes ─────────── */
@media (min-width: 768px) {
  .task-row {
    gap: var(--space-4);
  }
  
  .task-pet-info {
    min-width: 130px;
    gap: var(--space-3);
  }
  
  .task-info {
    margin-right: var(--space-4);
  }
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
  transition:
    background var(--transition-fast),
    border-color var(--transition-fast),
    color var(--transition-fast);
}

.btn-secondary:hover {
  background: var(--color-accent-light);
  border-color: var(--color-accent-muted);
  color: var(--color-accent-dark);
}
</style>
