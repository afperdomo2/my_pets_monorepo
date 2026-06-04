<script setup lang="ts">
import { IconPlus } from '@tabler/icons-vue'

import VaccineHistoryTable from '@/components/vaccines/VaccineHistoryTable.vue'
import VaccineUpcomingCard from '@/components/vaccines/VaccineUpcomingCard.vue'
import { useGetUpcomingVaccinesPaged } from '@/composables/useHealthRecords'

// ── Próximas vacunas con paginación acumulativa ─────────────────────
const {
  records,
  total,
  hasMore,
  loadMore,
  isLoading,
  isLoadingMore,
  isFetching,
  isError,
  refresh,
} = useGetUpcomingVaccinesPaged(5)
</script>

<template>
  <div class="vaccines-dashboard">
    <!-- Header -->
    <div class="page-header">
      <div class="page-header__text">
        <h1 class="page-title">Vacunas</h1>
        <p class="page-subtitle">Estado de inmunidad de tus mascotas</p>
      </div>
    </div>

    <!-- Timeline de próximas aplicaciones -->
    <section class="timeline-section">
      <div class="section-header">
        <h2 class="section-title">Próximas aplicaciones</h2>
        <span class="section-count">{{ total }} pendientes</span>
      </div>

      <div v-if="isLoading" class="loading-state">
        <div class="spinner" />
        <p>Cargando próximas vacunas...</p>
      </div>

      <div v-else-if="isError" class="error-state">
        <p>Error al cargar las próximas vacunas</p>
        <button class="btn-retry" @click="refresh">Reintentar</button>
      </div>

      <div v-else-if="records.length === 0" class="empty-state">
        <p>¡No hay vacunas pendientes! 🎉</p>
        <span>Todas las mascotas están al día con sus vacunas.</span>
      </div>

      <div v-else class="timeline-list">
        <VaccineUpcomingCard
          v-for="record in records"
          :key="record.id"
          :record-id="record.id"
          :pet-id="record.pet_id"
          :pet-name="record.pet.name"
          :pet-species="record.pet.species"
          :vaccine-name="record.name"
          :next-dose-date="record.next_dose_date || null"
        />

        <!-- Botón Ver más -->
        <button
          v-if="hasMore"
          class="btn-load-more"
          :disabled="isLoadingMore || isFetching"
          @click="loadMore"
        >
          <span v-if="isLoadingMore || isFetching" class="btn-load-more__spinner" />
          <template v-else>
            <IconPlus :size="18" :stroke-width="2" />
            Ver más
          </template>
        </button>
      </div>
    </section>

    <!-- Historial global -->
    <section class="history-section">
      <VaccineHistoryTable />
    </section>

  </div>
</template>

<style scoped>
/* ── Layout principal ─────────────── */
.vaccines-dashboard {
  width: 100%;
  /* clamp(): mínimo 1rem, ideal 5vw, máximo 2.5rem — se adapta sin media queries */
  padding: var(--space-8) clamp(var(--space-4), 5vw, var(--space-10));
  display: flex;
  flex-direction: column;
  gap: clamp(var(--space-5), 4vw, var(--space-8));
}

/* ── Header ───────────────────────── */
.page-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: var(--space-4);
  flex-wrap: wrap; /* se rompe solo cuando no cabe */
}

.page-title {
  font-size: clamp(var(--text-xl), 4vw, var(--text-2xl));
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

.btn-register-vaccine {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-5);
  min-height: 44px; /* área táctil mínima */
  background: var(--color-accent);
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition:
    background var(--transition-fast),
    transform var(--transition-fast),
    box-shadow var(--transition-fast);
  white-space: nowrap;
  flex-shrink: 0;
}

.btn-register-vaccine:hover {
  background: var(--color-accent-hover);
  transform: translateY(-1px);
  box-shadow: var(--shadow-md);
}

/* En móvil el botón ocupa todo el ancho */
@media (max-width: 30em) {
  .page-header {
    flex-direction: column;
    align-items: stretch;
  }

  .btn-register-vaccine {
    justify-content: center;
    width: 100%;
  }
}

/* ── Section headers ──────────────── */
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-2);
  margin-bottom: var(--space-4);
  flex-wrap: wrap;
}

.section-title {
  font-family: var(--font-display);
  font-size: var(--text-lg);
  font-weight: 600;
  color: var(--color-text-primary);
}

.section-count {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  font-weight: 500;
}

/* ── Timeline list ────────────────── */
.timeline-section {
  display: flex;
  flex-direction: column;
}

.timeline-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

/* ── Botón ver más ────────────────── */
.btn-load-more {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  align-self: center;
  padding: var(--space-2) var(--space-6);
  min-height: 44px;
  background: var(--color-surface);
  color: var(--color-accent);
  border: 1.5px solid var(--color-accent-muted);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition:
    background var(--transition-fast),
    color var(--transition-fast),
    border-color var(--transition-fast),
    box-shadow var(--transition-fast);
  white-space: nowrap;
}

.btn-load-more:hover:not(:disabled) {
  background: var(--color-accent);
  color: #fff;
  border-color: var(--color-accent);
  box-shadow: var(--shadow-sm);
}

.btn-load-more:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-load-more__spinner {
  display: inline-block;
  width: 16px;
  height: 16px;
  border: 2px solid var(--color-accent-muted);
  border-top-color: var(--color-accent);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
}

/* ── History section ──────────────── */
.history-section {
  display: flex;
  flex-direction: column;
}

/* ── Loading / Error / Empty states ──────────────── */
.loading-state,
.error-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--space-3);
  padding: var(--space-10) var(--space-4);
  text-align: center;
}

.loading-state p,
.error-state p,
.empty-state p {
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
  margin: 0;
}

.empty-state span {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
}

.spinner {
  width: 32px;
  height: 32px;
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
  min-height: 44px;
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
</style>
