<script setup lang="ts">
import { useGetExamById, useUpdateExam } from '@/composables/useExams'
import { computed, onMounted, onUnmounted, ref } from 'vue'
import { IconCheck, IconX } from '@tabler/icons-vue'
import DatePicker from '@/components/ui/DatePicker.vue'

const props = defineProps<{
  examId: string
}>()

const emit = defineEmits<{
  close: []
  updated: []
}>()

onMounted(() => {
  document.body.style.overflow = 'hidden'
})

onUnmounted(() => {
  document.body.style.overflow = ''
})

const examIdRef = computed(() => props.examId)
const { refetch } = useGetExamById(examIdRef)

const loading = ref(true)
const name = ref('')
const reason = ref('')
const scheduledDate = ref('')
const notes = ref('')

onMounted(async () => {
  const { data: fresh } = await refetch()
  if (fresh) {
    name.value = fresh.name
    reason.value = fresh.reason ?? ''
    scheduledDate.value = fresh.scheduled_date?.split('T')[0] ?? ''
    notes.value = fresh.notes ?? ''
  }
  loading.value = false
})

const updateExam = useUpdateExam()

const canSave = computed(() => {
  return name.value.trim().length > 0 && !updateExam.isPending.value
})

async function save() {
  if (!canSave.value) return
  try {
    await updateExam.mutateAsync({
      id: props.examId,
      payload: {
        name: name.value.trim(),
        reason: reason.value.trim() || undefined,
        scheduled_date: scheduledDate.value || undefined,
        notes: notes.value.trim() || undefined,
      },
    })
    emit('updated')
    emit('close')
  } catch (e) {
    console.error('Error al actualizar examen:', e)
  }
}
</script>

<template>
  <Teleport to="body">
    <div class="modal-backdrop" @click.self="emit('close')">
      <div class="modal-container">
        <div class="modal-header">
          <h2>Editar examen</h2>
          <button class="btn-close" @click="emit('close')">
            <IconX :size="18" :stroke-width="2" />
          </button>
        </div>

        <div v-if="loading" class="modal-loading">
          <div class="spinner" />
          <p>Cargando datos del examen...</p>
        </div>

        <template v-else>
          <div class="modal-body">
            <div class="field-group">
              <label class="field-label">
                Nombre del examen <span class="required">*</span>
              </label>
              <input
                v-model="name"
                class="field-input"
                placeholder="Ej: Análisis de sangre, Perfil renal..."
                autofocus
              />
            </div>

            <div class="field-group">
              <label class="field-label">
                Motivo/Razón <span class="optional">(opcional)</span>
              </label>
              <textarea
                v-model="reason"
                class="field-textarea"
                rows="2"
                placeholder="Motivo del examen..."
              />
            </div>

            <div class="field-group">
              <label class="field-label">
                Fecha programada <span class="optional">(opcional)</span>
              </label>
              <DatePicker
                v-model="scheduledDate"
                placeholder="Seleccionar fecha"
                unique-id="exam-edit-scheduled-date"
              />
            </div>

            <div class="field-group">
              <label class="field-label">
                Notas adicionales <span class="optional">(opcional)</span>
              </label>
              <textarea
                v-model="notes"
                class="field-textarea"
                rows="3"
                placeholder="Observaciones del veterinario..."
              />
            </div>
          </div>

          <div class="modal-footer">
            <div class="footer-spacer" />
            <button class="btn-cancel" @click="emit('close')">
              Cancelar
            </button>
            <button
              class="btn-save"
              :disabled="!canSave"
              @click="save"
            >
              <span v-if="updateExam.isPending.value" class="btn-spinner" />
              <template v-else>
                <IconCheck :size="16" :stroke-width="2.5" />
                Guardar cambios
              </template>
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
  width: min(520px, 100%);
  max-height: min(90vh, 600px);
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
  gap: var(--space-4);
  overflow-y: auto;
  flex: 1;
}

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

.field-input,
.field-textarea {
  width: 100%;
  padding: var(--space-2) var(--space-3);
  min-height: 44px;
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  background: var(--color-surface);
  transition: border-color var(--transition-fast);
  box-sizing: border-box;
}

.field-input:focus,
.field-textarea:focus {
  border-color: var(--color-accent);
  outline: none;
}

.field-textarea {
  resize: vertical;
  min-height: 60px;
  font-family: var(--font-body);
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
</style>
