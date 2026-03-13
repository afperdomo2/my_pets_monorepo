<script setup lang="ts">
import { ref } from 'vue'
import { IconShieldCheck, IconAlertTriangle, IconClockExclamation, IconPlus } from '@tabler/icons-vue'
import VaccineStatCard from '@/components/vaccines/VaccineStatCard.vue'
import VaccineTimelineCard from '@/components/vaccines/VaccineTimelineCard.vue'
import VaccineHistoryTable from '@/components/vaccines/VaccineHistoryTable.vue'
import VaccineRecordModal from '@/components/vaccines/VaccineRecordModal.vue'

// ── Estado del modal ────────────────────────────
const showModal = ref(false)
const preselectedPet = ref<string | undefined>(undefined)

function openModal(petId?: string) {
  preselectedPet.value = petId
  showModal.value = true
}

function closeModal() {
  showModal.value = false
  preselectedPet.value = undefined
}

// ── Datos mock del timeline ─────────────────────
const timelineItems = [
  {
    id: '1',
    petName: 'Romeo',
    petInitials: 'RO',
    petSpecies: 'dog' as const,
    lifeStage: 'Adulto',
    vaccineName: 'Antirrábica',
    dueDate: '20 feb 2026',
    dueLabel: 'Vencida hace 21 días',
    urgency: 'urgent' as const,
    petId: '1',
  },
  {
    id: '2',
    petName: 'Simba',
    petInitials: 'SI',
    petSpecies: 'cat' as const,
    lifeStage: 'Adulto',
    vaccineName: 'Triple felina',
    dueDate: '18 mar 2026',
    dueLabel: 'En 5 días',
    urgency: 'soon' as const,
    petId: '3',
  },
  {
    id: '3',
    petName: 'Manchas',
    petInitials: 'MA',
    petSpecies: 'dog' as const,
    lifeStage: 'Senior',
    vaccineName: 'Polivalente canina',
    dueDate: '10 abr 2026',
    dueLabel: 'En 28 días',
    urgency: 'soon' as const,
    petId: '4',
  },
  {
    id: '4',
    petName: 'Bolt',
    petInitials: 'BO',
    petSpecies: 'dog' as const,
    lifeStage: 'Cachorro',
    vaccineName: 'Sextuple',
    dueDate: '01 ago 2026',
    dueLabel: 'En 4 meses',
    urgency: 'future' as const,
    petId: '5',
  },
]
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
        <span class="section-count">{{ timelineItems.length }} pendientes</span>
      </div>
      <div class="timeline-list">
        <VaccineTimelineCard
          v-for="item in timelineItems"
          :key="item.id"
          :pet-name="item.petName"
          :pet-initials="item.petInitials"
          :pet-species="item.petSpecies"
          :life-stage="item.lifeStage"
          :vaccine-name="item.vaccineName"
          :due-date="item.dueDate"
          :due-label="item.dueLabel"
          :urgency="item.urgency"
          @register="openModal(item.petId)"
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
      @close="closeModal"
    />
  </div>
</template>

<style scoped>
/* ── Layout principal ─────────────── */
.vaccines-dashboard {
  width: 100%;
  padding: var(--space-8) var(--space-10);
  display: flex;
  flex-direction: column;
  gap: var(--space-8);
}

@media (max-width: 768px) {
  .vaccines-dashboard {
    padding: var(--space-5) var(--space-4);
    gap: var(--space-5);
  }
}

/* ── Header ───────────────────────── */
.page-header {
  display: flex;
  align-items: flex-end;
  justify-content: space-between;
  gap: var(--space-4);
}

@media (max-width: 480px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: var(--space-3);
  }

  .btn-register-vaccine {
    width: 100%;
    justify-content: center;
  }
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

.btn-register-vaccine {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-5);
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

/* ── Stats section ────────────────── */
.stats-section {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-4);
}

@media (max-width: 768px) {
  .stats-section {
    grid-template-columns: 1fr;
    gap: var(--space-3);
  }
}

/* ── Section headers ──────────────── */
.section-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-4);
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
</style>
