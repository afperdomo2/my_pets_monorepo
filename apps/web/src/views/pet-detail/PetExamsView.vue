<script setup lang="ts">
import ConfirmDeleteModal from '@/components/health-tabs/ConfirmDeleteModal.vue'
import ExamEditModal from '@/components/health-tabs/ExamEditModal.vue'
import ExamResultsModal from '@/components/health-tabs/ExamResultsModal.vue'
import AppPagination from '@/components/ui/AppPagination.vue'
import PerPageSelector from '@/components/ui/PerPageSelector.vue'
import {
  useGetExamsByPet,
  useCreateExam,
  useDeleteExam,
  useUpdateExam,
} from '@/composables/useExams'
import { ExamStatus, type Exam } from '@/types/exam'
import { formatDateOnly } from '@/utils/date'
import {
  IconCalendar,
  IconCheck,
  IconEdit,
  IconEye,
  IconFileText,
  IconPlus,
  IconRefresh,
  IconStethoscope,
  IconTrash,
} from '@tabler/icons-vue'
import DatePicker from '@/components/ui/DatePicker.vue'
import { computed, ref, type Ref } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const petId = computed(() => String(route.params.id))

const page = ref(1)
const perPage = ref(10)

const { data, isLoading, isError, refresh } = useGetExamsByPet(petId as Ref<string>, page, perPage)

const records = computed(() => data.value?.data ?? [])
const total = computed(() => data.value?.total ?? 0)
const totalPages = computed(() => data.value?.total_pages ?? 0)

const STATUS_CONFIG: Record<string, { label: string; className: string }> = {
  [ExamStatus.Scheduled]: { label: 'Programado', className: 'status--scheduled' },
  [ExamStatus.Completed]: { label: 'Completado', className: 'status--completed' },
}

// ── Modals state ──
const showCreateModal = ref(false)
const showEditModal = ref(false)
const showResultsModal = ref(false)
const showConfirmModal = ref(false)

const examIdToEdit = ref<string | null>(null)
const examIdToView = ref<string | null>(null)
const recordToDelete = ref<Exam | null>(null)
const deletingId = ref<string | null>(null)

// ── Create form state ──
const examName = ref('')
const examReason = ref('')
const examScheduledDate = ref('')
const examNotes = ref('')
const isCompleted = ref(false)
const examCompletedDate = ref('')

const examFields = ref<{ name: string; value: string; unit: string }[]>([
  { name: '', value: '', unit: '' },
])

function addField() {
  examFields.value.push({ name: '', value: '', unit: '' })
}

function removeField(index: number) {
  examFields.value.splice(index, 1)
}

function resetCreateForm() {
  examName.value = ''
  examReason.value = ''
  examScheduledDate.value = ''
  examNotes.value = ''
  isCompleted.value = false
  examCompletedDate.value = ''
  examFields.value = [{ name: '', value: '', unit: '' }]
}

function openCreate() {
  resetCreateForm()
  showCreateModal.value = true
}

function openEdit(exam: Exam) {
  examIdToEdit.value = exam.id
  showEditModal.value = true
}

function openResults(exam: Exam) {
  examIdToView.value = exam.id
  showResultsModal.value = true
}

function openDeleteConfirm(exam: Exam) {
  recordToDelete.value = exam
  showConfirmModal.value = true
}

// ── Mutations ──
const createExam = useCreateExam()
const updateExam = useUpdateExam()
const deleteExam = useDeleteExam()

function getResultsFromFields(): Array<{ parameter_name: string; value: string; unit?: string }> {
  return examFields.value
    .filter(f => f.name.trim() && f.value.trim())
    .map(f => ({
      parameter_name: f.name.trim(),
      value: f.value.trim(),
      unit: f.unit.trim() || undefined,
    }))
}

async function saveExam() {
  const results = getResultsFromFields()

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
  showCreateModal.value = false
}

async function handleDeleteConfirm() {
  if (!recordToDelete.value) return
  deletingId.value = recordToDelete.value.id
  try {
    await deleteExam.mutateAsync({ id: recordToDelete.value.id, petId: petId.value })
    showConfirmModal.value = false
    recordToDelete.value = null
  } finally {
    deletingId.value = null
  }
}
</script>

<template>
  <div class="tab-view">
    <div class="content-card">
      <!-- Header -->
      <div class="tab-header">
        <div class="tab-header-left">
          <h2 class="tab-title">
            <IconStethoscope :size="20" :stroke-width="1.75" />
            Historial de Exámenes
          </h2>
          <span class="panel-count">{{ total }} registros</span>
        </div>
        <div class="tab-header-right">
          <button class="btn-refresh" :disabled="isLoading" title="Refrescar" @click="refresh">
            <IconRefresh :size="16" :stroke-width="2" :class="{ spin: isLoading }" />
            Refrescar
          </button>
          <button class="btn-add" @click="openCreate">
            <IconPlus :size="16" :stroke-width="2.5" />
            Registrar
          </button>
        </div>
      </div>

      <!-- Tabla desktop -->
      <div class="table-wrapper">
        <table v-if="!isLoading && !isError && records.length > 0" class="exam-table">
          <thead>
            <tr>
              <th>Nombre</th>
              <th class="th-center">Estado</th>
              <th class="th-center">Fecha programada</th>
              <th class="th-center">Fecha realización</th>
              <th class="th-center">Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-for="exam in records" :key="exam.id" class="exam-row">
              <td>
                <div class="exam-cell">
                  <span class="exam-name">{{ exam.name }}</span>
                  <span v-if="exam.reason" class="exam-reason" :title="exam.reason">{{
                    exam.reason
                  }}</span>
                </div>
              </td>
              <td class="td-center">
                <span
                  class="status-badge"
                  :class="STATUS_CONFIG[exam.status]?.className ?? ''"
                >
                  <IconCheck
                    v-if="exam.status === ExamStatus.Completed"
                    :size="12"
                    :stroke-width="2.5"
                  />
                  <IconFileText
                    v-else
                    :size="12"
                    :stroke-width="2.5"
                  />
                  {{ STATUS_CONFIG[exam.status]?.label ?? exam.status }}
                </span>
              </td>
              <td class="td-center">
                <span class="date-cell">{{ formatDateOnly(exam.scheduled_date) }}</span>
              </td>
              <td class="td-center">
                <span class="date-cell">{{ formatDateOnly(exam.completed_date) }}</span>
              </td>
              <td class="td-center">
                <div class="action-buttons">
                  <button
                    class="btn-action btn-view"
                    title="Ver resultados"
                    @click="openResults(exam)"
                  >
                    <IconEye :size="14" :stroke-width="2" />
                    Ver resultados
                  </button>
                  <button
                    class="btn-action btn-edit"
                    title="Editar"
                    :disabled="updateExam.isPending.value"
                    @click="openEdit(exam)"
                  >
                    <IconEdit :size="14" :stroke-width="2" />
                    Editar
                  </button>
                  <button
                    class="btn-action btn-delete"
                    title="Eliminar"
                    :disabled="deletingId === exam.id || deleteExam.isPending.value"
                    @click="openDeleteConfirm(exam)"
                  >
                    <IconTrash :size="14" :stroke-width="2" />
                    Eliminar
                  </button>
                </div>
              </td>
            </tr>
          </tbody>
        </table>

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
      </div>

      <!-- Vista card para móvil -->
      <div v-if="!isLoading && !isError && records.length > 0" class="record-cards">
        <div v-for="exam in records" :key="`card-${exam.id}`" class="record-card">
          <div class="record-card__top">
            <div class="record-card__info">
              <span class="record-card__name">{{ exam.name }}</span>
              <span
                class="status-badge"
                :class="STATUS_CONFIG[exam.status]?.className ?? ''"
              >
                <IconCheck
                  v-if="exam.status === ExamStatus.Completed"
                  :size="10"
                  :stroke-width="2.5"
                />
                <IconFileText
                  v-else
                  :size="10"
                  :stroke-width="2.5"
                />
                {{ STATUS_CONFIG[exam.status]?.label ?? exam.status }}
              </span>
            </div>
            <div class="record-card__actions">
              <button
                class="btn-action btn-view"
                title="Ver resultados"
                @click="openResults(exam)"
              >
                <IconEye :size="14" :stroke-width="2" />
              </button>
              <button
                class="btn-action btn-edit"
                title="Editar"
                :disabled="updateExam.isPending.value"
                @click="openEdit(exam)"
              >
                <IconEdit :size="14" :stroke-width="2" />
              </button>
              <button
                class="btn-action btn-delete"
                title="Eliminar"
                :disabled="deletingId === exam.id || deleteExam.isPending.value"
                @click="openDeleteConfirm(exam)"
              >
                <IconTrash :size="14" :stroke-width="2" />
              </button>
            </div>
          </div>

          <div v-if="exam.reason" class="record-card__detail">
            <span class="record-card__label">Motivo</span>
            <span class="record-card__value">{{ exam.reason }}</span>
          </div>

          <div class="record-card__dates">
            <div class="record-card__date-item">
              <span class="record-card__label">Programado</span>
              <span class="record-card__value">{{ formatDateOnly(exam.scheduled_date) }}</span>
            </div>
            <div class="record-card__date-item">
              <span class="record-card__label">Realizado</span>
              <span class="record-card__value">{{ formatDateOnly(exam.completed_date) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Footer con paginación -->
      <div v-if="!isLoading && !isError && records.length > 0" class="table-footer">
        <PerPageSelector v-model="perPage" :options="[10, 25, 50]" />
        <AppPagination
          :current-page="page"
          :total-pages="totalPages"
          :total-items="total"
          :per-page="perPage"
          @update:page="page = $event"
        />
      </div>
    </div>

    <!-- Modal crear examen -->
    <Teleport to="body">
      <div v-if="showCreateModal" class="modal-backdrop" @click.self="showCreateModal = false">
        <div class="modal-container">
          <div class="modal-header">
            <h2>Registrar examen</h2>
            <button class="btn-close" @click="showCreateModal = false">✕</button>
          </div>

          <div class="modal-body">
            <div class="field">
              <label class="field-label">Nombre del examen <span class="required">*</span></label>
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
                <label class="field-label">Fecha de realización <span class="required">*</span></label>
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
                    <IconCalendar :size="16" :stroke-width="2" />
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
            <button class="btn-cancel" @click="showCreateModal = false">Cancelar</button>
            <button
              class="btn-save"
              :disabled="!examName || (isCompleted && !examCompletedDate) || createExam.isPending.value"
              @click="saveExam"
            >
              <span v-if="createExam.isPending.value" class="spinner-sm" />
              Registrar examen
            </button>
          </div>
        </div>
      </div>
    </Teleport>

    <!-- Modal editar examen -->
    <ExamEditModal
      v-if="showEditModal && examIdToEdit"
      :exam-id="examIdToEdit"
      @close="showEditModal = false; examIdToEdit = null"
      @updated="refresh"
    />

    <!-- Modal ver resultados -->
    <ExamResultsModal
      v-if="showResultsModal && examIdToView"
      :exam-id="examIdToView"
      @close="showResultsModal = false; examIdToView = null"
    />

    <!-- Modal confirmar eliminación -->
    <ConfirmDeleteModal
      v-model="showConfirmModal"
      :record-name="recordToDelete?.name ?? ''"
      record-type="exam"
      :deleting="deleteExam.isPending.value"
      @confirm="handleDeleteConfirm"
    />
  </div>
</template>

<style scoped>
.tab-view {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.content-card {
  container-type: inline-size;
  container-name: exam-card;
  background: transparent;
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-xl);
  overflow: hidden;
  box-shadow: none;
}

/* ── Header ───────────────────────── */
.tab-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  flex-wrap: wrap;
  padding: var(--space-4) var(--space-5);
  border-bottom: 1px solid var(--color-border-light);
}

.tab-header-left {
  display: flex;
  align-items: center;
  gap: var(--space-3);
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

.panel-count {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  font-weight: 500;
}

.tab-header-right {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.btn-refresh {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  background: var(--color-bg-alt);
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition:
    background var(--transition-fast),
    color var(--transition-fast);
  flex-shrink: 0;
}

.btn-refresh:hover:not(:disabled) {
  background: var(--color-accent-light);
  color: var(--color-accent-dark);
  border-color: var(--color-accent);
}

.btn-refresh:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-refresh .spin {
  animation: spin 0.8s linear infinite;
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

/* ── Tabla ────────────────────────── */
.table-wrapper {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.exam-table {
  width: 100%;
  min-width: 700px;
  border-collapse: collapse;
}

.exam-table thead tr {
  background: var(--color-bg);
  border-bottom: 1px solid var(--color-border-light);
}

.exam-table th {
  padding: var(--space-3) var(--space-4);
  text-align: left;
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.exam-table th.th-center {
  text-align: center;
}

.exam-row {
  transition: background var(--transition-fast);
}

.exam-row:hover {
  background: var(--color-bg-alt);
}

.exam-row td {
  padding: var(--space-3) var(--space-4);
  border-bottom: 1px solid var(--color-border-light);
}

.exam-row:last-child td {
  border-bottom: none;
}

.td-center {
  text-align: center;
}

/* ── Celdas ───────────────────────── */
.exam-cell {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.exam-name {
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  font-weight: 500;
}

.exam-reason {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  line-height: 1.4;
}

.date-cell {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  white-space: nowrap;
}

/* ── Badges de estado ─────────────── */
.status-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px var(--space-3);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: 600;
  white-space: nowrap;
}

.status--scheduled {
  background: #fef3e2;
  color: #c4714a;
}

.status--completed {
  background: #e8f5ee;
  color: #2e7d52;
}

/* ── Botones de acción ────────────── */
.action-buttons {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  justify-content: center;
}

.btn-action {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-1);
  padding: var(--space-2) var(--space-3);
  border-radius: var(--radius-md);
  font-size: var(--text-xs);
  font-weight: 600;
  cursor: pointer;
  border: 1px solid transparent;
  transition:
    background var(--transition-fast),
    color var(--transition-fast),
    border-color var(--transition-fast);
  white-space: nowrap;
}

/* Ver resultados - Bordered con fondo claro */
.btn-view {
  background: #e0f2fe;
  color: #0284c7;
  border-color: #bae6fd;
}

.btn-view:hover:not(:disabled) {
  background: #bae6fd;
  border-color: #0284c7;
}

/* Editar - Bordered con fondo claro */
.btn-edit {
  background: #fef3c7;
  color: #d97706;
  border-color: #fde68a;
}

.btn-edit:hover:not(:disabled) {
  background: #fde68a;
  border-color: #f59e0b;
}

/* Eliminar - Bordered con fondo claro */
.btn-delete {
  background: #fef2f2;
  color: #dc2626;
  border-color: #fecaca;
}

.btn-delete:hover:not(:disabled) {
  background: #fee2e2;
  border-color: #fca5a5;
}

.btn-delete:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

.btn-view:disabled,
.btn-edit:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

/* ── Estados de carga/error/vacío ─── */
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
  to {
    transform: rotate(360deg);
  }
}

/* ── Vista card para móvil (< 700px) ─ */
.record-cards {
  display: none;
}

@container exam-card (max-width: 699px) {
  .table-wrapper {
    display: none;
  }

  .record-cards {
    display: flex;
    flex-direction: column;
    gap: 0;
  }

  .record-card {
    display: flex;
    flex-direction: column;
    gap: var(--space-3);
    padding: var(--space-4);
    border-bottom: 1px solid var(--color-border-light);
  }

  .record-card:last-child {
    border-bottom: none;
  }

  .record-card__top {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: var(--space-2);
  }

  .record-card__info {
    display: flex;
    align-items: center;
    gap: var(--space-2);
    flex: 1;
    min-width: 0;
  }

  .record-card__name {
    font-size: var(--text-sm);
    font-weight: 600;
    color: var(--color-text-primary);
    flex: 1;
    min-width: 0;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .record-card__actions {
    display: flex;
    gap: var(--space-1);
    flex-shrink: 0;
  }

  .record-card__detail {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .record-card__dates {
    display: flex;
    gap: var(--space-4);
  }

  .record-card__date-item {
    display: flex;
    flex-direction: column;
    gap: 2px;
  }

  .record-card__label {
    font-size: 0.65rem;
    font-weight: 600;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--color-text-tertiary);
  }

  .record-card__value {
    font-size: var(--text-xs);
    color: var(--color-text-secondary);
  }
}

/* ── Footer con paginación ─────────── */
.table-footer {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
  flex-wrap: wrap;
  padding: var(--space-4) var(--space-5);
  border-top: 1px solid var(--color-border-light);
  background: var(--color-bg);
}

@container exam-card (max-width: 480px) {
  .tab-header {
    padding: var(--space-3) var(--space-4);
  }

  .table-footer {
    flex-direction: column;
    align-items: stretch;
    padding: var(--space-3) var(--space-4);
    gap: var(--space-3);
  }
}

/* ── Modal (create) ───────────────── */
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

.required {
  color: var(--color-error);
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

.field-textarea {
  resize: vertical;
}

.date-picker-row {
  display: flex;
  gap: var(--space-2);
}

.btn-today {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
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

.dynamic-fields .field-input {
  background: #fff;
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
