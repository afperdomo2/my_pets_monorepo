<script setup lang="ts">
import { ref, computed, type Ref } from 'vue'
import { useRoute } from 'vue-router'
import {
  IconPlus,
  IconStethoscope,
  IconCheck,
  IconChevronDown,
  IconChevronUp,
  IconFileText,
  IconTrash
} from '@tabler/icons-vue'
import { useGetExamsByPet, useCreateExam, useUpdateExam, useDeleteExam, useCompleteExam } from '@/composables/useExams'
import { ExamStatus, type ExamStatusType, type Exam } from '@/types/exam'
import DatePicker from '@/components/ui/DatePicker.vue'
import AppPagination from '@/components/ui/AppPagination.vue'

const route = useRoute()
const petId = computed(() => String(route.params.id))

const page = ref(1)
const perPage = ref(10)

const { data, isLoading, isError, refresh } = useGetExamsByPet(petId as Ref<string>, page, perPage)

const records = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.total ?? 0)
const totalPages = computed(() => data.value?.total_pages ?? 0)

const STATUS_CONFIG: Record<ExamStatusType, { label: string; className: string; icon: typeof IconFileText }> = {
  [ExamStatus.Scheduled]: { label: 'Programado', className: 'status--scheduled', icon: IconFileText },
  [ExamStatus.Completed]: { label: 'Completado', className: 'status--completed', icon: IconCheck },
}

function formatDate(dateStr: string | null): string {
  if (!dateStr) return '—'
  const date = new Date(dateStr)
  return date.toLocaleDateString('es-ES', { day: '2-digit', month: 'short', year: 'numeric' })
}

const expandedExamId = ref<string | null>(null)

function toggleExpand(id: string) {
  expandedExamId.value = expandedExamId.value === id ? null : id
}

function isExpanded(id: string): boolean {
  return expandedExamId.value === id
}

// Modal para crear/editar examen
const showExamModal = ref(false)
const editingExam = ref<Exam | null>(null)
const examName = ref('')
const examReason = ref('')
const examScheduledDate = ref('')
const examNotes = ref('')
const isCompleted = ref(false)
const examCompletedDate = ref('')

// Campos dinámicos del examen
const examFields = ref<{ name: string; value: string; unit: string }[]>([
  { name: '', value: '', unit: '' }
])

function addField() {
  examFields.value.push({ name: '', value: '', unit: '' })
}

function removeField(index: number) {
  examFields.value.splice(index, 1)
}

function openCreate() {
  editingExam.value = null
  examName.value = ''
  examReason.value = ''
  examScheduledDate.value = ''
  examNotes.value = ''
  isCompleted.value = false
  examCompletedDate.value = ''
  examFields.value = [{ name: '', value: '', unit: '' }]
  showExamModal.value = true
}

function openEdit(exam: Exam) {
  editingExam.value = exam
  examName.value = exam.name
  examReason.value = exam.reason ?? ''
  examScheduledDate.value = exam.scheduled_date ?? ''
  examNotes.value = exam.notes ?? ''
  isCompleted.value = exam.status === ExamStatus.Completed
  examCompletedDate.value = exam.completed_date ?? ''
  showExamModal.value = true
}

function getResultsFromFields(): Array<{ parameter_name: string; value: string; unit?: string }> {
  return examFields.value
    .filter(f => f.name.trim() && f.value.trim())
    .map(f => ({
      parameter_name: f.name.trim(),
      value: f.value.trim(),
      unit: f.unit.trim() || undefined,
    }))
}

const createExam = useCreateExam()
const updateExam = useUpdateExam()
const completeExam = useCompleteExam()

async function saveExam() {
  const results = getResultsFromFields()

  if (editingExam.value) {
    // Actualizar examen existente
    await updateExam.mutateAsync({
      id: editingExam.value.id,
      payload: {
        name: examName.value,
        reason: examReason.value || undefined,
        scheduled_date: examScheduledDate.value || undefined,
        notes: examNotes.value || undefined,
      },
    })

    // Si está completado y tiene resultados, usar el endpoint complete
    if (isCompleted.value && results.length > 0) {
      await completeExam.mutateAsync({
        id: editingExam.value.id,
        payload: {
          completed_date: examCompletedDate.value || new Date().toISOString().split('T')[0],
          results,
        },
      })
    }
  } else {
    // Crear nuevo examen
    const payload: {
      pet_id: string
      name: string
      reason?: string
      scheduled_date?: string
      notes?: string
      status?: 'scheduled' | 'completed'
      completed_date?: string
      results?: Array<{ parameter_name: string; value: string; unit?: string }>
    } = {
      pet_id: petId.value,
      name: examName.value,
      reason: examReason.value || undefined,
      scheduled_date: examScheduledDate.value || undefined,
      notes: examNotes.value || undefined,
    }

    if (isCompleted.value) {
      payload.status = 'completed'
      payload.completed_date = examCompletedDate.value || undefined
      payload.results = results
    }

    await createExam.mutateAsync(payload)
  }

  showExamModal.value = false
}

const deleteExam = useDeleteExam()

async function handleDelete(exam: Exam) {
  if (!confirm(`¿Eliminar el examen "${exam.name}"?`)) return
  await deleteExam.mutateAsync({ id: exam.id, petId: petId.value })
}
</script>

<template>
  <div class="tab-view">
    <div class="content-card">
      <!-- Header -->
      <div class="tab-header">
        <h2 class="tab-title">
          <IconStethoscope :size="20" :stroke-width="1.75" />
          Historial de Exámenes
        </h2>
        <button class="btn-add" @click="openCreate">
          <IconPlus :size="16" :stroke-width="2.5" />
          Registrar examen
        </button>
      </div>

      <!-- Loading -->
      <div v-if="isLoading" class="loading-state">
        <div class="spinner" />
        <p>Cargando registros...</p>
      </div>

      <!-- Error -->
      <div v-else-if="isError" class="error-state">
        <p>Error al cargar los registros</p>
        <button class="btn-retry" @click="refresh">Reintentar</button>
      </div>

      <!-- Empty -->
      <div v-else-if="records.length === 0" class="empty-state">
        <IconStethoscope :size="40" :stroke-width="1.5" />
        <p>No hay exámenes registrados</p>
        <button class="btn-add-empty" @click="openCreate">
          <IconPlus :size="16" :stroke-width="2.5" />
          Registrar primer examen
        </button>
      </div>

      <!-- Accordion list -->
      <div v-else class="accordion-list">
        <div
          v-for="exam in records"
          :key="exam.id"
          class="accordion-item"
          :class="{ 'accordion-item--expanded': isExpanded(exam.id) }"
        >
          <!-- Header -->
          <button class="accordion-header" @click="toggleExpand(exam.id)">
            <div class="accordion-header-left">
              <div class="accordion-icon">
                <IconFileText :size="18" :stroke-width="1.75" />
              </div>
              <div class="accordion-info">
                <span class="accordion-title">{{ exam.name }}</span>
                <span class="accordion-date">{{ formatDate(exam.scheduled_date || exam.completed_date) }}</span>
              </div>
            </div>
            <div class="accordion-header-right">
              <span
                class="status-badge"
                :class="STATUS_CONFIG[exam.status as ExamStatusType]?.className"
              >
                <component
                  :is="STATUS_CONFIG[exam.status as ExamStatusType]?.icon"
                  :size="12"
                  :stroke-width="2.5"
                />
                {{ STATUS_CONFIG[exam.status as ExamStatusType]?.label }}
              </span>
              <component
                :is="isExpanded(exam.id) ? IconChevronUp : IconChevronDown"
                :size="18"
                :stroke-width="2"
                class="accordion-chevron"
              />
            </div>
          </button>

          <!-- Content -->
          <div v-if="isExpanded(exam.id)" class="accordion-content">
            <!-- Reason -->
            <div v-if="exam.reason" class="reason-section">
              <h4 class="reason-title">Motivo</h4>
              <p class="reason-text">{{ exam.reason }}</p>
            </div>

            <!-- Notes -->
            <div v-if="exam.notes" class="notes-section">
              <h4 class="notes-title">Notas</h4>
              <p class="notes-text">{{ exam.notes }}</p>
            </div>

            <!-- Actions -->
            <div class="accordion-actions">
              <button class="btn-edit-exam" @click="openEdit(exam)">
                ✏️ Editar
              </button>
              <button class="btn-delete-exam" @click="handleDelete(exam)">
                <IconTrash :size="14" :stroke-width="2" />
                Eliminar
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Pagination -->
      <div v-if="totalPages > 1" class="pagination-bar">
        <AppPagination
          :current-page="page"
          :total-pages="totalPages"
          :total-items="total"
          :per-page="perPage"
          @update:page="page = $event"
        />
      </div>
    </div>

    <!-- Modal: Create/Edit Exam -->
    <Teleport to="body">
      <div v-if="showExamModal" class="modal-backdrop" @click.self="showExamModal = false">
        <div class="modal-container">
          <div class="modal-header">
            <h2>{{ editingExam ? 'Editar examen' : 'Registrar examen' }}</h2>
            <button class="btn-close" @click="showExamModal = false">✕</button>
          </div>

          <div class="modal-body">
            <div class="field">
              <label class="field-label">Nombre del examen *</label>
              <input
                v-model="examName"
                class="field-input"
                placeholder="Ej: Análisis de sangre, Perfil renal..."
              />
            </div>

            <div class="field">
              <label class="field-label">Motivo/Razón (opcional)</label>
              <textarea
                v-model="examReason"
                class="field-textarea"
                rows="2"
                placeholder="Motivo del examen..."
              />
            </div>

            <div class="field">
              <label class="field-label">Fecha programada</label>
              <div class="date-picker-row">
                <DatePicker
                  v-model="examScheduledDate"
                  placeholder="Seleccionar fecha"
                  unique-id="exam-scheduled-date"
                />
              </div>
            </div>

            <div class="field">
              <label class="field-toggle">
                <input v-model="isCompleted" type="checkbox" class="toggle-checkbox" />
                <span class="toggle-label">El examen ya está completado</span>
              </label>
            </div>

            <div v-if="isCompleted" class="completed-section">
              <div class="field">
                <label class="field-label">Fecha de realización *</label>
                <div class="date-picker-row">
                  <DatePicker
                    v-model="examCompletedDate"
                    :max-date="new Date()"
                    placeholder="Seleccionar fecha"
                    unique-id="exam-completed-date"
                  />
                  <button
                    type="button"
                    class="btn-today"
                    @click="examCompletedDate = new Date().toISOString().split('T')[0] ?? ''"
                  >
                    Hoy
                  </button>
                </div>
              </div>

              <div class="fields-section">
                <div class="fields-header">
                  <label class="field-label">Resultados (campos dinámicos)</label>
                  <button type="button" class="btn-add-field" @click="addField">
                    <IconPlus :size="14" :stroke-width="2.5" />
                    Añadir campo
                  </button>
                </div>

                <div class="dynamic-fields">
                  <div
                    v-for="(field, idx) in examFields"
                    :key="idx"
                    class="dynamic-field-row"
                  >
                    <input
                      v-model="field.name"
                      class="field-input field-name"
                      placeholder="Parámetro (ej: Glucosa)"
                    />
                    <input
                      v-model="field.value"
                      class="field-input field-value"
                      placeholder="Valor"
                    />
                    <input
                      v-model="field.unit"
                      class="field-input field-unit"
                      placeholder="Unidad"
                    />
                    <button
                      v-if="examFields.length > 1"
                      type="button"
                      class="btn-remove-field"
                      @click="removeField(idx)"
                    >
                      ✕
                    </button>
                  </div>
                </div>
              </div>
            </div>

            <div class="field">
              <label class="field-label">Notas adicionales</label>
              <textarea
                v-model="examNotes"
                class="field-textarea"
                rows="3"
                placeholder="Observaciones del veterinario..."
              />
            </div>
          </div>

          <div class="modal-footer">
            <button class="btn-cancel" @click="showExamModal = false">Cancelar</button>
            <button
              class="btn-save"
              :disabled="!examName || (isCompleted && !examCompletedDate) || createExam.isPending.value"
              @click="saveExam"
            >
              <span v-if="createExam.isPending.value" class="spinner-sm" />
              {{ editingExam ? 'Guardar cambios' : 'Registrar examen' }}
            </button>
          </div>
        </div>
      </div>
    </Teleport>
  </div>
</template>

<style scoped>
/* Estilos similares a los existentes, adaptados para la nueva estructura */
.tab-view {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.content-card {
  background: transparent;
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-xl);
  overflow: hidden;
  box-shadow: none;
}

.tab-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  flex-wrap: wrap;
  padding: var(--space-4) var(--space-5);
  border-bottom: 1px solid var(--color-border-light);
}

.tab-title {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  font-family: var(--font-display);
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.btn-add {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  background: var(--color-accent);
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast);
}

.btn-add:hover {
  background: var(--color-accent-hover);
}

/* States */
.loading-state,
.error-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: var(--space-10) var(--space-4);
  gap: var(--space-3);
  color: var(--color-text-tertiary);
  text-align: center;
}

.empty-state p {
  margin: 0;
  font-size: var(--text-sm);
}

.btn-add-empty {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  background: var(--color-accent-light);
  color: var(--color-accent-dark);
  border: 1px solid var(--color-accent);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
}

.btn-retry {
  padding: var(--space-2) var(--space-4);
  background: var(--color-accent);
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
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

/* Accordion */
.accordion-list {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.accordion-item {
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.accordion-header {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  padding: var(--space-4);
  background: transparent;
  border: none;
  cursor: pointer;
  text-align: left;
  transition: background var(--transition-fast);
}

.accordion-header:hover {
  background: var(--color-bg-alt);
}

.accordion-header-left {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  flex: 1;
  min-width: 0;
}

.accordion-icon {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f0f9ff;
  color: #0284c7;
  border-radius: var(--radius-md);
  flex-shrink: 0;
}

.accordion-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.accordion-title {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
}

.accordion-date {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
}

.accordion-header-right {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px var(--space-2);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: 600;
  flex-shrink: 0;
}

.status--scheduled {
  background: #fef3e2;
  color: #c4714a;
}

.status--completed {
  background: #e8f5ee;
  color: #2e7d52;
}

.accordion-chevron {
  color: var(--color-text-tertiary);
  transition: transform var(--transition-fast);
}

.accordion-item--expanded .accordion-chevron {
  transform: rotate(180deg);
}

.accordion-content {
  padding: 0 var(--space-4) var(--space-4);
  border-top: 1px solid var(--color-border-light);
}

.reason-section,
.notes-section {
  padding-top: var(--space-4);
}

.reason-title,
.notes-title {
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin: 0 0 var(--space-2) 0;
}

.reason-text,
.notes-text {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  margin: 0;
  white-space: pre-wrap;
}

.accordion-actions {
  display: flex;
  gap: var(--space-2);
  margin-top: var(--space-4);
  padding-top: var(--space-3);
  border-top: 1px solid var(--color-border-light);
}

.btn-edit-exam,
.btn-delete-exam {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  padding: var(--space-2) var(--space-3);
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: var(--text-xs);
  font-weight: 500;
  cursor: pointer;
  transition: background var(--transition-fast);
}

.btn-edit-exam:hover {
  background: var(--color-bg-alt);
}

.btn-delete-exam {
  color: #dc2626;
  border-color: #fecaca;
}

.btn-delete-exam:hover {
  background: #fef2f2;
}

.pagination-bar {
  display: flex;
  justify-content: center;
  padding: var(--space-4) var(--space-5);
  border-top: 1px solid var(--color-border-light);
}

/* Modal */
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
  width: min(500px, 100%);
  max-height: 90vh;
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
}

.modal-header h2 {
  font-family: var(--font-display);
  font-size: var(--text-lg);
  font-weight: 600;
  margin: 0;
}

.btn-close {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  cursor: pointer;
  font-size: var(--text-lg);
}

.modal-body {
  padding: var(--space-5);
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.field {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.field-label {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-secondary);
}

.field-input,
.field-textarea {
  padding: var(--space-2) var(--space-3);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  background: var(--color-bg);
  font-family: var(--font-body);
}

.field-input:focus,
.field-textarea:focus {
  outline: none;
  border-color: var(--color-accent);
}

.date-picker-row {
  display: flex;
  gap: var(--space-2);
}

.btn-today {
  padding: var(--space-2) var(--space-3);
  background: var(--color-accent-light);
  color: var(--color-accent-dark);
  border: 1px solid var(--color-accent);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
}

.field-toggle {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  cursor: pointer;
}

.toggle-checkbox {
  width: 16px;
  height: 16px;
  cursor: pointer;
}

.toggle-label {
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-text-secondary);
}

.completed-section {
  padding: var(--space-3);
  background: var(--color-bg);
  border-radius: var(--radius-md);
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.fields-section {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.fields-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.btn-add-field {
  display: inline-flex;
  align-items: center;
  gap: var(--space-1);
  padding: var(--space-1) var(--space-2);
  background: transparent;
  color: var(--color-accent);
  border: 1px dashed var(--color-accent);
  border-radius: var(--radius-sm);
  font-size: var(--text-xs);
  font-weight: 600;
  cursor: pointer;
}

.dynamic-fields {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.dynamic-field-row {
  display: flex;
  gap: var(--space-2);
  align-items: center;
}

.field-name {
  flex: 2;
}

.field-value {
  flex: 1;
  min-width: 80px;
}

.field-unit {
  flex: 1;
  min-width: 60px;
}

.btn-remove-field {
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  cursor: pointer;
  color: var(--color-text-tertiary);
  flex-shrink: 0;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  padding: var(--space-4) var(--space-5);
  border-top: 1px solid var(--color-border-light);
}

.btn-cancel {
  padding: var(--space-2) var(--space-4);
  background: transparent;
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
}

.btn-save {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  background: var(--color-accent);
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
}

.btn-save:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.spinner-sm {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255,255,255,0.3);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}
</style>
