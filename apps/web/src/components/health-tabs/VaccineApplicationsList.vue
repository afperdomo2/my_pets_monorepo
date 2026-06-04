<script setup lang="ts">
import type { VaccineApplication } from '@/types/vaccineApplication'
import { IconVaccine, IconInfoCircle } from '@tabler/icons-vue'
import { formatDateOnly } from '@/utils/date'

defineProps<{
  applications: VaccineApplication[]
  category: 'vaccine' | 'deworming'
}>()

defineEmits<{
  delete: [id: string]
}>()
</script>

<template>
  <div class="vaccine-applications-list">
    <h3 class="section-title">
      <IconVaccine :size="16" :stroke-width="2" />
      {{ category === 'vaccine' ? 'Dosis aplicadas' : 'Dosis aplicadas' }}
    </h3>

    <div v-if="applications.length === 0" class="empty-applications">
      <IconInfoCircle :size="16" :stroke-width="2" />
      <span>Sin aplicaciones registradas</span>
    </div>

    <div v-else class="applications-scroll">
      <div
        v-for="(app, index) in applications"
        :key="app.id"
        class="application-row"
      >
        <span class="dose-badge">#{{ index + 1 }}</span>
        <div class="dose-info">
          <span class="dose-date">{{ formatDateOnly(app.application_date) }}</span>
          <span v-if="app.notes" class="dose-notes">{{ app.notes }}</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.vaccine-applications-list {
  margin-top: var(--space-4);
  padding-top: var(--space-4);
  border-top: 1px solid var(--color-border-light);
}

.section-title {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-secondary);
  margin: 0 0 var(--space-3) 0;
}

.empty-applications {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-4);
  background: var(--color-bg-alt);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
}

.empty-applications svg {
  flex-shrink: 0;
  color: var(--color-text-tertiary);
}

.applications-scroll {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.application-row {
  display: flex;
  align-items: flex-start;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-4);
  background: var(--color-bg);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  transition: background var(--transition-fast);
}

.application-row:hover {
  background: var(--color-bg-alt);
}

.dose-badge {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 32px;
  height: 24px;
  padding: 0 var(--space-2);
  background: var(--color-accent-light);
  color: var(--color-accent-dark);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: 700;
  flex-shrink: 0;
  margin-top: 1px;
}

.dose-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.dose-date {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
}

.dose-notes {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  line-height: 1.4;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>
