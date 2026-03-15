<script setup lang="ts">
import { ref } from 'vue'
import { IconShieldCheck, IconAlertTriangle, IconClockExclamation, IconPlus } from '@tabler/icons-vue'
import VaccineStatCard from '@/components/vaccines/VaccineStatCard.vue'
import VaccineUpcomingCard from '@/components/vaccines/VaccineUpcomingCard.vue'
import VaccineHistoryTable from '@/components/vaccines/VaccineHistoryTable.vue'
import VaccineRecordModal from '@/components/vaccines/VaccineRecordModal.vue'
import { useGetUpcomingVaccines } from '@/composables/useHealthRecords'

// ── Estado del modal ────────────────────────────
const showModal = ref(false)
const preselectedPet = ref<string | undefined>(undefined)
const preselectedVaccine = ref<string | undefined>(undefined)

function openModal(petId?: string, vaccineName?: string) {
  preselectedPet.value = petId
  preselectedVaccine.value = vaccineName
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  preselectedPet.value = undefined
  preselectedVaccine.value = undefined
}

// ── Próximas vacunas (datos reales) ─────────────────────
const { records, total, isLoading, isError, refresh } = useGetUpcomingVaccines(10)
</script>

<template>
  <div class="vaccines-dashboard">
    <!-- Header -->
    <div class="page-header">
      <div class="page-header__text">
        <h1 class="page-title">Vacunas</h1>
        <p class="page-subtitle">Estado de inmunidad de tus mascotas</p>
      </div>
      <button class="btn-register-vaccine" @click="openModal()">
        <IconPlus :size="16" :stroke-width="2.5" />
        Registrar vacuna
      </button>
    </div>

    <!-- Quick Stats -->
    <section class="stats-section">
      <VaccineStatCard
        label="Al día"
        value="5"
        description="Todo en orden"
        color="green"
        :icon="IconShieldCheck"
      />
      <VaccineStatCard
        label="Por vencer"
        value="2"
        description="Próximos 15 días"
        color="amber"
        :icon="IconAlertTriangle"
      />
      <VaccineStatCard
        label="Vencidas"
        value="1"
        description="Requiere atención"
        color="red"
        :icon="IconClockExclamation"
      />
    </section>

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
          :due-date="record.due_date"
          @register="openModal(record.pet_id, record.name)"
        />
      </div>
    </section>

    <!-- Historial global -->
    <section class="history-section">
      <VaccineHistoryTable />
    </section>

    <!-- Modal registrar vacuna -->
    <VaccineRecordModal
      v-if="showModal"
      :preselected-pet="preselectedPet"
      :preselected-vaccine="preselectedVaccine"
      @close="closeModal"
    />
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
  transition: background var(--transition-fast), transform var(--transition-fast),
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

/* ── Stats section ────────────────── */
.stats-section {
  display: grid;
  /* auto-fit + minmax: 3 cols en desktop, 2 en tablet, 1 en móvil — sin media queries */
  grid-template-columns: repeat(auto-fit, minmax(min(260px, 100%), 1fr));
  gap: var(--space-4);
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
