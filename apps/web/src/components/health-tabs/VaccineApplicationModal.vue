<script setup lang="ts">
import DatePicker from '@/components/ui/DatePicker.vue'
import VaccineApplicationsList from '@/components/health-tabs/VaccineApplicationsList.vue'
import { useGetHealthRecord } from '@/composables/useHealthRecords'
import { useCreateVaccineApplication } from '@/composables/useVaccineApplications'
import { IconBell, IconCheck, IconInfoCircle, IconX } from '@tabler/icons-vue'
import { computed, onMounted, onUnmounted, ref, toRef, watch } from 'vue'

onMounted(() => {
  document.body.style.overflow = 'hidden'
})

onUnmounted(() => {
  document.body.style.overflow = ''
})

const props = defineProps<{
  healthRecordId: string
  category: 'vaccine' | 'deworming'
  preselectedDate?: string
}>()

const emit = defineEmits<{
  close: []
  applied: []
}>()

const healthRecordIdRef = toRef(props, 'healthRecordId')
const { data: healthRecord, isLoading: isRecordLoading } = useGetHealthRecord(healthRecordIdRef)

const totalDoses = computed(() => healthRecord.value?.total_doses ?? null)
const appliedDosesCount = computed(() => healthRecord.value?.applied_doses_count ?? 0)

const currentDoseNumber = computed(() => appliedDosesCount.value + 1)

const doseMessage = computed(() => {
  if (totalDoses.value === null || totalDoses.value === undefined) return null
  if (totalDoses.value === 1) return null
  if (currentDoseNumber.value === totalDoses.value) return 'Esta es la última dosis'
  return `Esta es la aplicación ${currentDoseNumber.value} de ${totalDoses.value}`
})

const applicationDate = ref(props.preselectedDate || '')
const note = ref('')
const wantsNext = ref(false)
const nextDate = ref('')

// Cuando se cargan los datos y totalDoses > 1, activar el toggle por defecto
watch(totalDoses, (val) => {
  if (val !== null && val !== undefined && val > 1) {
    wantsNext.value = true
  }
})

const createApplication = useCreateVaccineApplication()

const canSave = computed(() => {
  return !!applicationDate.value
})

async function save() {
  if (!canSave.value) return

  try {
    await createApplication.mutateAsync({
      health_record_id: props.healthRecordId,
      application_date: applicationDate.value,
      notes: note.value.trim() || undefined,
      next_dose_date: wantsNext.value && nextDate.value ? nextDate.value : null,
    })
    emit('applied')
    emit('close')
  } catch (e) {
    console.error('Error al aplicar dosis:', e)
  }
}

function setToday() {
  applicationDate.value = new Date().toISOString().split('T')[0] || ''
}

const modalTitle = computed(() => {
  return props.category === 'vaccine' ? 'Aplicar vacuna' : 'Aplicar desparasitación'
})

const notePlaceholder = computed(() => {
  return props.category === 'vaccine'
    ? 'Ej: Lote #12345, veterinario Dr. Pérez...'
    : 'Ej: Producto aplicado, zona de aplicación...'
})
</script>

<template>
  <Teleport to="body">
    <div class="modal-backdrop" @click.self="emit('close')">
      <div class="modal-container">
        <!-- Header -->
        <div class="modal-header">
          <h2>{{ modalTitle }}</h2>
          <button class="btn-close" @click="emit('close')">
            <IconX :size="18" :stroke-width="2" />
          </button>
        </div>

        <!-- Body -->
        <div class="modal-body">
          <!-- Indicador de dosis actual -->
          <div v-if="doseMessage && !isRecordLoading" class="dose-banner">
            <IconInfoCircle :size="18" :stroke-width="2" />
            <span>{{ doseMessage }}</span>
          </div>

          <!-- Fecha de aplicación -->
          <div class="field-group">
            <label class="field-label">
              Fecha de aplicación <span class="required">*</span>
            </label>
            <div class="date-field">
              <div class="date-input-wrapper">
                <DatePicker
                  v-model="applicationDate"
                  :max-date="new Date()"
                  placeholder="Seleccionar fecha"
                  unique-id="app-date-apply"
                />
              </div>
              <button
                type="button"
                class="btn-today"
                @click="setToday"
              >
                Hoy
              </button>
            </div>
          </div>

          <!-- Nota -->
          <div class="field-group">
            <label class="field-label">
              Nota <span class="optional">(opcional)</span>
            </label>
            <textarea
              v-model="note"
              class="note-input"
              rows="4"
              :placeholder="notePlaceholder"
            />
          </div>

          <!-- Programar refuerzo (oculto si es dosis única) -->
          <div v-if="totalDoses !== 1" class="suggestion-card">
            <div class="suggestion-icon">
              <IconBell :size="18" :stroke-width="1.75" />
            </div>
            <div class="suggestion-body">
              <div class="suggestion-header">
                <p class="suggestion-text">
                  <strong>Programar refuerzo</strong>
                  <br>Programa la próxima dosis para mantener el historial al día.
                </p>
                <label class="suggestion-toggle">
                  <input v-model="wantsNext" type="checkbox" class="toggle-checkbox" />
                  <span class="toggle-track">
                    <span class="toggle-thumb" />
                  </span>
                  <span class="toggle-label">{{ wantsNext ? 'Sí' : 'No' }}</span>
                </label>
              </div>

              <template v-if="wantsNext">
                <div class="next-date-wrapper">
                  <label class="field-label field-label--small">
                    Fecha próximo refuerzo
                  </label>
                  <DatePicker
                    v-model="nextDate"
                    :min-date="applicationDate ? new Date(applicationDate) : new Date()"
                    placeholder="Seleccionar fecha"
                    unique-id="next-date-apply"
                  />
                </div>
              </template>
            </div>
          </div>

          <!-- Dosis aplicadas -->
          <VaccineApplicationsList
            v-if="healthRecord?.vaccine_applications && healthRecord.vaccine_applications.length > 0"
            :applications="healthRecord.vaccine_applications"
            :category="props.category"
          />
        </div>

        <!-- Footer -->
        <div class="modal-footer">
          <div class="footer-spacer" />
          <button
            class="btn-cancel"
            @click="emit('close')"
          >
            Cancelar
          </button>
          <button
            class="btn-save"
            :disabled="!canSave || createApplication.isPending.value"
            @click="save"
          >
            <span v-if="createApplication.isPending.value" class="btn-spinner" />
            <template v-else>
              <IconCheck :size="16" :stroke-width="2.5" />
              Aplicar
            </template>
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
/* ── Backdrop ──────────────────────────────────────────────── */
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

/* ── Modal container ───────────────────────────────────────── */
.modal-container {
  background: var(--color-surface);
  width: min(520px, 100%);
  max-height: min(90vh, 700px);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-xl);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* ── Header ────────────────────────────────────────────────── */
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-4) var(--space-5);
  border-bottom: 1px solid var(--color-border-light);
  flex-shrink: 0;
}

.modal-header h2 {
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

/* ── Body ──────────────────────────────────────────────────── */
.modal-body {
  padding: var(--space-5);
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  overflow-y: auto;
  flex: 1;
}

/* ── Campos del formulario ────────────────────────────────── */
.field-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.field-label {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-bottom: var(--space-1);
  display: block;
}

.required {
  color: #dc2626;
  font-weight: 700;
}

.optional {
  color: var(--color-text-tertiary);
  font-weight: 400;
  font-size: var(--text-xs);
}

.date-field {
  display: flex;
  gap: var(--space-2);
  align-items: stretch;
}

.date-input-wrapper {
  flex: 1;
  min-width: 0;
}

.btn-today {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-3);
  min-width: 80px;
  background: var(--color-accent-light);
  color: var(--color-accent-dark);
  border: 1.5px solid var(--color-accent);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: all var(--transition-fast);
  white-space: nowrap;
}

.btn-today:hover {
  background: var(--color-accent);
  color: #fff;
}

.note-input {
  width: 100%;
  padding: var(--space-3);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  background: var(--color-surface);
  resize: vertical;
  min-height: 100px;
  font-family: var(--font-body);
  box-sizing: border-box;
  transition: border-color var(--transition-fast);
}

.note-input:focus {
  border-color: var(--color-accent);
  outline: none;
}

/* ── Suggestion card ──────────────────────────────────────── */
.suggestion-card {
  display: flex;
  gap: var(--space-3);
  padding: var(--space-4);
  background: #f0f9ff;
  border: 1px solid #bae6fd;
  border-radius: var(--radius-lg);
  margin-top: auto;
}

.suggestion-icon {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-md);
  background: #dbeafe;
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.suggestion-body {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  min-width: 0;
}

.suggestion-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--space-2);
  flex-wrap: wrap;
}

.suggestion-text {
  font-size: var(--text-sm);
  color: #1e40af;
  line-height: 1.5;
  margin: 0;
}

.suggestion-toggle {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  cursor: pointer;
  flex-shrink: 0;
  min-height: 32px;
}

.toggle-checkbox {
  display: none;
}

.toggle-track {
  width: 36px;
  height: 20px;
  border-radius: var(--radius-full);
  background: var(--color-border);
  position: relative;
  transition: background var(--transition-fast);
  flex-shrink: 0;
}

.toggle-checkbox:checked + .toggle-track {
  background: var(--color-accent);
}

.toggle-thumb {
  width: 16px;
  height: 16px;
  border-radius: var(--radius-full);
  background: #fff;
  position: absolute;
  top: 2px;
  left: 2px;
  transition: transform var(--transition-fast);
  box-shadow: var(--shadow-xs);
}

.toggle-checkbox:checked + .toggle-track .toggle-thumb {
  transform: translateX(16px);
}

.toggle-label {
  font-size: var(--text-xs);
  font-weight: 600;
  color: #1e40af;
}

.next-date-wrapper {
  margin-top: var(--space-3);
}

.field-label--small {
  font-size: var(--text-xs);
  margin-bottom: var(--space-1);
}

/* ── Footer ───────────────────────────────────────────────── */
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

.btn-cancel {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
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

.btn-cancel:hover {
  background: var(--color-bg-alt);
  border-color: var(--color-text-tertiary);
}

.btn-save {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-5);
  min-height: 44px;
  background: #2e7d52;
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast);
}

.btn-save:hover:not(:disabled) {
  background: #256644;
}

.btn-save:disabled {
  background: var(--color-border);
  color: var(--color-text-tertiary);
  cursor: not-allowed;
}

.btn-spinner {
  display: inline-block;
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255,255,255,0.35);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ── Dose banner ──────────────────────────────────────────── */
.dose-banner {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-4);
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  color: #166534;
}

.dose-banner svg {
  color: #22c55e;
  flex-shrink: 0;
}
</style>
