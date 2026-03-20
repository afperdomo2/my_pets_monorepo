<script setup lang="ts">
import PetAvatar from '@/components/pets/PetAvatar.vue'
import { useDashboardSummary } from '@/composables/useDashboard'
import { useGetUpcomingRecordsPaged } from '@/composables/useHealthRecords'
import { getSpeciesLabel } from '@/constants/species'
import { useAuthStore } from '@/stores/auth'
import {
  IconAlertTriangle,
  IconPaw,
  IconPill,
  IconPlus,
  IconShieldCheck,
  IconStethoscope,
  IconVaccine,
} from '@tabler/icons-vue'
import { computed } from 'vue'
import { RouterLink } from 'vue-router'

const authStore = useAuthStore()
const userName = authStore.user?.name ?? ''

// Dashboard summary
const { totalPets, healthyPets, pendingTasks, overdueTasks } = useDashboardSummary()

// Stats computados basados en datos del dashboard
const stats = computed(() => [
  {
    label: 'Mascotas registradas',
    value: String(totalPets.value),
    change: totalPets.value > 0 ? `+${totalPets.value} en total` : 'No hay mascotas registradas',
    positive: true,
    icon: IconPaw,
    color: 'accent',
    to: '/pets',
  },
  {
    label: 'Mascotas al día',
    value: String(healthyPets.value),
    change:
      totalPets.value > 0
        ? `${Math.round((healthyPets.value / totalPets.value) * 100)}% con salud al día`
        : '0% con salud al día',
    positive: true,
    icon: IconShieldCheck,
    color: 'green',
    to: '/health-records',
  },
  {
    label: 'Tareas pendientes',
    value: String(pendingTasks.value),
    change: 'Tareas que están programadas',
    positive: false,
    icon: IconAlertTriangle,
    color: 'orange',
    to: '/vaccines',
  },
  {
    label: 'Tareas vencidas',
    value: String(overdueTasks.value),
    change: 'Requieren atención',
    positive: false,
    icon: IconAlertTriangle,
    color: 'red',
    to: '/vaccines',
  },
])

// API - Próximas tareas de salud (sin filtrar por categoría)
const { records: upcomingRecords, isLoading, isError, refresh } = useGetUpcomingRecordsPaged(5)

const getCategoryIcon = (category: string) => {
  switch (category) {
    case 'vaccine':
      return IconVaccine
    case 'deworming':
      return IconPill
    case 'exam':
      return IconStethoscope
    default:
      return IconShieldCheck
  }
}

const getCategoryColor = (category: string) => {
  switch (category) {
    case 'vaccine':
      return 'var(--color-accent)'
    case 'deworming':
      return '#c4714a'
    case 'exam':
      return '#2980b9'
    default:
      return 'var(--color-accent)'
  }
}

const getCategoryLabel = (category: string) => {
  switch (category) {
    case 'vaccine':
      return 'Vacuna'
    case 'deworming':
      return 'Desparasitación'
    case 'exam':
      return 'Examen'
    default:
      return 'Salud'
  }
}

const getCategoryRoute = (petId: string, category: string) => {
  switch (category) {
    case 'vaccine':
      return { name: 'pet-detail-vaccines', params: { id: petId } }
    case 'deworming':
      return { name: 'pet-detail-deworming', params: { id: petId } }
    case 'exam':
      return { name: 'pet-detail-exams', params: { id: petId } }
    default:
      return { name: 'pet-detail', params: { id: petId } }
  }
}

function formatDueDate(dueDate: string): string {
  const due = new Date(dueDate)
  const now = new Date()
  due.setUTCHours(0, 0, 0, 0)
  now.setUTCHours(0, 0, 0, 0)
  const diffMs = due.getTime() - now.getTime()
  const days = Math.ceil(diffMs / (1000 * 60 * 60 * 24))

  if (days < 0) return `Vencida hace ${Math.abs(days)} días`
  if (days === 0) return 'Hoy'
  if (days === 1) return 'Mañana'
  return `En ${days} días`
}

function isUrgent(dueDate: string): boolean {
  const due = new Date(dueDate)
  const now = new Date()
  due.setUTCHours(0, 0, 0, 0)
  now.setUTCHours(0, 0, 0, 0)
  const diffMs = due.getTime() - now.getTime()
  const days = Math.ceil(diffMs / (1000 * 60 * 60 * 24))
  return days <= 7 && days >= 0
}

function isOverdue(dueDate: string): boolean {
  const due = new Date(dueDate)
  const now = new Date()
  due.setUTCHours(0, 0, 0, 0)
  now.setUTCHours(0, 0, 0, 0)
  const diffMs = due.getTime() - now.getTime()
  const days = Math.ceil(diffMs / (1000 * 60 * 60 * 24))
  return days < 0
}
</script>

<template>
  <div class="home-view">
    <!-- Page header -->
    <div class="page-header">
      <div class="page-header-text">
        <h1>Buenos días, {{ userName }}</h1>
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
        <span
          v-if="stat.change"
          class="stat-change"
          :class="{ 'stat-change--negative': !stat.positive }"
        >
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
        </div>

        <!-- Loading state -->
        <div v-if="isLoading" class="panel-loading">
          <div class="spinner" />
          <p>Cargando tareas...</p>
        </div>

        <!-- Error state -->
        <div v-else-if="isError" class="panel-error">
          <p>Error al cargar las tareas</p>
          <button class="btn-retry" @click="refresh">Reintentar</button>
        </div>

        <!-- Empty state -->
        <div v-else-if="upcomingRecords.length === 0" class="panel-empty">
          <p>¡No hay tareas pendientes! 🎉</p>
          <span>Todas las mascotas están al día.</span>
        </div>

        <!-- Records list -->
        <div v-else class="panel-body">
          <RouterLink
            v-for="record in upcomingRecords"
            :key="record.id"
            :to="getCategoryRoute(record.pet_id, record.category)"
            class="task-row"
            :class="{
              'task-row--urgent': record.next_dose_date && isUrgent(record.next_dose_date),
              'task-row--overdue': record.next_dose_date && isOverdue(record.next_dose_date),
            }"
          >
            <!-- Mascota -->
            <div class="task-pet-info">
              <PetAvatar :species="record.pet.species" :name="record.pet.name" size="md" />
              <div class="task-pet-text">
                <span class="task-pet-name">{{ record.pet.name }}</span>
                <span class="task-pet-species">{{ getSpeciesLabel(record.pet.species) }}</span>
              </div>
            </div>
            <!-- Tarea de salud -->
            <div class="task-details">
              <div class="task-icon" :style="{ color: getCategoryColor(record.category) }">
                <component :is="getCategoryIcon(record.category)" :size="16" :stroke-width="1.75" />
              </div>
              <div class="task-info">
                <span class="task-name">{{ record.name }}</span>
                <span class="task-category">{{ getCategoryLabel(record.category) }}</span>
              </div>
            </div>
            <!-- Fecha -->
            <span class="task-due" :class="{ 'task-due--urgent': record.next_dose_date && isUrgent(record.next_dose_date) }">
              {{ record.next_dose_date ? formatDueDate(record.next_dose_date) : 'Sin programar' }}
            </span>
          </RouterLink>
        </div>

        <!-- CTA -->
        <div class="panel-footer">
          <RouterLink to="/vaccines" class="btn-secondary"> Ver vacunación </RouterLink>
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
  grid-template-columns: repeat(4, 1fr);
  gap: var(--space-4);
}

@media (max-width: 1200px) {
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
.stat-card--red .stat-icon {
  background: #fef2f2;
  color: #dc2626;
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

/* ── Loading / Error / Empty states ───────────────────────── */
.panel-loading,
.panel-error,
.panel-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--space-3);
  padding: var(--space-8) var(--space-4);
  text-align: center;
}

.panel-loading p,
.panel-error p,
.panel-empty p {
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
  margin: 0;
}

.panel-empty span {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
}

.spinner {
  width: 28px;
  height: 28px;
  border: 3px solid var(--color-border);
  border-top-color: var(--color-accent);
  border-radius: 50%;
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.btn-retry {
  padding: var(--space-2) var(--space-4);
  min-height: 40px;
  background: var(--color-accent);
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast);
}

.btn-retry:hover {
  background: var(--color-accent-hover);
}

/* ── Task rows: grid de 4 columnas (mascota | tarea x2 | fecha) ─ */
.task-row {
  display: grid;
  grid-template-columns: 1fr 2fr 1fr;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-6);
  text-decoration: none;
  border-left: 3px solid transparent;
  transition:
    background var(--transition-fast),
    border-color var(--transition-fast);
}

.task-row:hover {
  background: var(--color-bg-alt);
}

.task-row--urgent {
  border-left-color: var(--color-warning, #f59e0b);
}

.task-row--overdue {
  border-left-color: var(--color-error);
}

.task-row--overdue .task-due {
  color: var(--color-error);
  font-weight: 600;
}

/* Mascota */
.task-pet-info {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  min-width: 0;
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

/* Detalles de tarea (icono + info) */
.task-details {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  min-width: 0;
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
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
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
  text-align: right;
  white-space: nowrap;
}

.task-due--urgent {
  color: var(--color-secondary);
  font-weight: 600;
}

/* ── Responsive ───────────────────────────────────────────── */
@media (max-width: 600px) {
  .task-row {
    grid-template-columns: 1fr;
    gap: var(--space-2);
    padding: var(--space-3) var(--space-4);
    border-bottom: 1px solid var(--color-border-light);
  }

  .task-row:last-child {
    border-bottom: none;
  }

  .task-pet-info {
    gap: var(--space-2);
  }

  .task-details {
    order: 3;
    margin-top: var(--space-1);
    gap: var(--space-1);
  }

  .task-due {
    text-align: left;
    order: 4;
    margin-top: var(--space-1);
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
