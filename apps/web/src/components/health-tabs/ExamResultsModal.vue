<script setup lang="ts">
import { useGetExamById } from '@/composables/useExams'
import { ExamStatus } from '@/types/exam'
import { formatDateOnly } from '@/utils/date'
import {
  IconClipboardList,
  IconFileText,
  IconX,
} from '@tabler/icons-vue'
import { computed, onMounted, onUnmounted, ref } from 'vue'

const props = defineProps<{
  examId: string
}>()

const emit = defineEmits<{
  close: []
}>()

onMounted(() => {
  document.body.style.overflow = 'hidden'
})

onUnmounted(() => {
  document.body.style.overflow = ''
})

const examIdRef = computed(() => props.examId)
const { data: examWithResults, refetch } = useGetExamById(examIdRef)

const loading = ref(true)

onMounted(async () => {
  await refetch()
  loading.value = false
})

const statusConfig: Record<string, { label: string; class: string }> = {
  [ExamStatus.Scheduled]: { label: 'Programado', class: 'status-scheduled' },
  [ExamStatus.Completed]: { label: 'Completado', class: 'status-completed' },
}

const hasResults = computed(() => {
  return examWithResults.value?.results && examWithResults.value.results.length > 0
})
</script>

<template>
  <Teleport to="body">
    <div class="modal-backdrop" @click.self="emit('close')">
      <div class="modal-container">
        <div class="modal-header">
          <div class="modal-header-left">
            <IconClipboardList :size="20" :stroke-width="1.75" />
            <h2>Resultados del examen</h2>
          </div>
          <button class="btn-close" @click="emit('close')">
            <IconX :size="18" :stroke-width="2" />
          </button>
        </div>

        <div v-if="loading" class="modal-loading">
          <div class="spinner" />
          <p>Cargando resultados...</p>
        </div>

        <template v-else-if="examWithResults">
          <div class="modal-body">
            <div class="exam-info-card">
              <div class="exam-info-top">
                <h3 class="exam-name">{{ examWithResults.name }}</h3>
                <span
                  class="status-badge"
                  :class="statusConfig[examWithResults.status]?.class ?? ''"
                >
                  {{ statusConfig[examWithResults.status]?.label ?? examWithResults.status }}
                </span>
              </div>

              <div class="exam-info-dates">
                <div class="date-item">
                  <span class="date-label">Programado</span>
                  <span class="date-value">{{ formatDateOnly(examWithResults.scheduled_date) }}</span>
                </div>
                <div class="date-item">
                  <span class="date-label">Realizado</span>
                  <span class="date-value">{{ formatDateOnly(examWithResults.completed_date) }}</span>
                </div>
              </div>

              <div v-if="examWithResults.reason" class="exam-info-section">
                <span class="section-label">Motivo</span>
                <p class="section-text">{{ examWithResults.reason }}</p>
              </div>

              <div v-if="examWithResults.notes" class="exam-info-section">
                <span class="section-label">Notas</span>
                <p class="section-text">{{ examWithResults.notes }}</p>
              </div>
            </div>

            <div class="results-section">
              <h3 class="results-title">
                <IconFileText :size="16" :stroke-width="2" />
                Resultados
              </h3>

              <div
                v-if="!hasResults"
                class="no-results"
              >
                <IconFileText :size="32" :stroke-width="1.5" />
                <p v-if="examWithResults.status === ExamStatus.Scheduled">
                  El examen aún no se ha realizado. No hay resultados disponibles.
                </p>
                <p v-else>
                  Este examen no tiene resultados registrados.
                </p>
              </div>

              <table v-else class="results-table">
                <thead>
                  <tr>
                    <th>Parámetro</th>
                    <th class="th-center">Valor</th>
                    <th class="th-center">Unidad</th>
                  </tr>
                </thead>
                <tbody>
                  <tr
                    v-for="result in examWithResults.results"
                    :key="result.id"
                    class="result-row"
                  >
                    <td>
                      <span class="param-name">{{ result.parameter_name }}</span>
                    </td>
                    <td class="td-center">
                      <span class="param-value">{{ result.value }}</span>
                    </td>
                    <td class="td-center">
                      <span class="param-unit">{{ result.unit ?? '—' }}</span>
                    </td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>

          <div class="modal-footer">
            <div class="footer-spacer" />
            <button class="btn-close-modal" @click="emit('close')">
              Cerrar
            </button>
          </div>
        </template>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(26, 26, 24, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: var(--space-4);
}

.modal-container {
  background: var(--color-surface);
  width: min(560px, 100%);
  max-height: min(90vh, 700px);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-xl);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-4) var(--space-5);
  border-bottom: 1px solid var(--color-border-light);
  flex-shrink: 0;
}

.modal-header-left {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.modal-header-left h2 {
  font-family: var(--font-display);
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.btn-close {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  color: var(--color-text-tertiary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.btn-close:hover {
  background: var(--color-bg-alt);
  color: var(--color-text-primary);
}

.modal-loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-10) var(--space-4);
  gap: var(--space-3);
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
  to { transform: rotate(360deg); }
}

.modal-body {
  padding: var(--space-5);
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
  overflow-y: auto;
  flex: 1;
}

.exam-info-card {
  background: var(--color-bg);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  padding: var(--space-4);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.exam-info-top {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-2);
}

.exam-name {
  font-family: var(--font-display);
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
  flex: 1;
  min-width: 0;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px var(--space-3);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: 600;
  white-space: nowrap;
  flex-shrink: 0;
}

.status-scheduled {
  background: #fef3e2;
  color: #c4714a;
}

.status-completed {
  background: #e8f5ee;
  color: #2e7d52;
}

.exam-info-dates {
  display: flex;
  gap: var(--space-6);
}

.date-item {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.date-label {
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.date-value {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.exam-info-section {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.section-label {
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

.section-text {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  margin: 0;
  white-space: pre-wrap;
  line-height: 1.5;
}

.results-section {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.results-title {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-family: var(--font-display);
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.no-results {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-6) var(--space-4);
  color: var(--color-text-tertiary);
  text-align: center;
}

.no-results p {
  margin: 0;
  font-size: var(--text-sm);
  line-height: 1.5;
}

.results-table {
  width: 100%;
  border-collapse: collapse;
}

.results-table thead tr {
  background: var(--color-bg);
  border-bottom: 1px solid var(--color-border-light);
}

.results-table th {
  padding: var(--space-3) var(--space-4);
  text-align: left;
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.results-table th.th-center {
  text-align: center;
}

.result-row {
  transition: background var(--transition-fast);
}

.result-row:hover {
  background: var(--color-bg-alt);
}

.result-row td {
  padding: var(--space-3) var(--space-4);
  border-bottom: 1px solid var(--color-border-light);
}

.result-row:last-child td {
  border-bottom: none;
}

.td-center {
  text-align: center;
}

.param-name {
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-text-primary);
}

.param-value {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
}

.param-unit {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.modal-footer {
  display: flex;
  align-items: center;
  padding: var(--space-3) var(--space-5);
  border-top: 1px solid var(--color-border-light);
  gap: var(--space-3);
  flex-shrink: 0;
}

.footer-spacer {
  flex: 1;
}

.btn-close-modal {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-5);
  min-height: 44px;
  background: transparent;
  color: var(--color-text-secondary);
  border: 1.5px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast);
}

.btn-close-modal:hover {
  background: var(--color-bg-alt);
  border-color: var(--color-text-tertiary);
}
</style>
