<script setup lang="ts">
import { useUpdateHealthRecord } from '@/composables/useHealthRecords'
import type { HealthRecord } from '@/types/healthRecord'
import { IconCheck, IconX } from '@tabler/icons-vue'
import { computed, onMounted, onUnmounted, ref } from 'vue'

onMounted(() => {
  document.body.style.overflow = 'hidden'
})

onUnmounted(() => {
  document.body.style.overflow = ''
})

const props = defineProps<{
  record: HealthRecord
  category: 'vaccine' | 'deworming'
}>()

const emit = defineEmits<{
  close: []
  updated: [record: HealthRecord]
}>()

const name = ref(props.record.name)
const note = ref(props.record.notes || '')
const totalDoses = ref<number | undefined>(props.record.total_doses ?? undefined)

const updateRecord = useUpdateHealthRecord()

const canSave = computed(() => {
  return name.value.trim().length > 0
})

async function save() {
  if (!canSave.value) return

  try {
    const updated = await updateRecord.mutateAsync({
      id: props.record.id,
      payload: {
        category: props.record.category,
        name: name.value.trim(),
        notes: note.value.trim() || undefined,
        total_doses: totalDoses.value,
      },
    })
    emit('updated', updated)
    emit('close')
  } catch (e) {
    console.error('Error al actualizar registro:', e)
  }
}

const modalTitle = computed(() => {
  return props.category === 'vaccine' ? 'Editar vacuna' : 'Editar desparasitación'
})

const nameLabel = computed(() => {
  return props.category === 'vaccine' ? 'Nombre de la vacuna' : 'Nombre del antiparasitario'
})

const namePlaceholder = computed(() => {
  return props.category === 'vaccine' ? 'Ej: Vacuna contra la rabia' : 'Ej: Antiparasitario externo'
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
          <!-- Nombre -->
          <div class="field-group">
            <label class="field-label">
              {{ nameLabel }} <span class="required">*</span>
            </label>
            <input
              v-model="name"
              class="field-input"
              :placeholder="namePlaceholder"
              autofocus
            />
          </div>

          <!-- Total de dosis (opcional) -->
          <div class="field-group">
            <label class="field-label">
              Total de dosis <span class="optional">(opcional)</span>
            </label>
            <input
              v-model.number="totalDoses"
              type="number"
              min="1"
              class="field-input field-input--white"
              placeholder="Ej: 3"
            />
            <span class="field-help">Cantidad total de dosis que indicará el veterinario</span>
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
              placeholder="Ej: Lote #12345, veterinatoria donde se aplicó..."
            />
          </div>
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
            :disabled="!canSave || updateRecord.isPending.value"
            @click="save"
          >
            <span v-if="updateRecord.isPending.value" class="btn-spinner" />
            <template v-else>
              <IconCheck :size="16" :stroke-width="2.5" />
              Guardar cambios
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
  max-height: min(90vh, 600px);
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

.field-input {
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

.field-input:focus {
  border-color: var(--color-accent);
  outline: none;
}

.field-input--white {
  background: #fff;
}

.field-help {
  display: block;
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  margin-top: 4px;
  line-height: 1.4;
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
</style>
