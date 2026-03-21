<script setup lang="ts">
import { IconTrash, IconX, IconAlertTriangle } from '@tabler/icons-vue'

defineProps<{
  modelValue: boolean
  recordName: string
  recordType?: string
  deleting?: boolean
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  confirm: []
}>()

function cancel() {
  emit('update:modelValue', false)
}

function getTypeLabel(type?: string): string {
  if (!type) return 'registro'
  const labels: Record<string, string> = {
    vaccine: 'vacuna',
    deworming: 'desparasitación',
    exam: 'examen',
    treatment: 'tratamiento',
  }
  return labels[type] ?? 'registro'
}
</script>

<template>
  <Teleport to="body">
    <Transition name="confirm-modal">
      <div v-if="modelValue" class="backdrop" @click.self="cancel">
        <div class="dialog" role="alertdialog" aria-modal="true" aria-labelledby="confirm-title">

          <!-- Close -->
          <button class="close-btn" :disabled="deleting" @click="cancel">
            <IconX :size="16" :stroke-width="2.5" />
          </button>

          <!-- Icon -->
          <div class="warning-icon">
            <IconAlertTriangle :size="28" :stroke-width="1.75" />
          </div>

          <!-- Content -->
          <div class="dialog-body">
            <h2 id="confirm-title" class="dialog-title">Eliminar registro</h2>
            <p class="dialog-desc">
              ¿Estás seguro de que querés eliminar
              <strong class="record-highlight">
                "{{ recordName }}"
              </strong>?
              <br />
              Este {{ getTypeLabel(recordType) }} se eliminará permanentemente y esta acción no se puede deshacer.
            </p>
          </div>

          <!-- Actions -->
          <div class="dialog-actions">
            <button class="btn-cancel" :disabled="deleting" @click="cancel">
              Cancelar
            </button>
            <button class="btn-confirm" :disabled="deleting" @click="emit('confirm')">
              <span v-if="deleting" class="btn-spinner" />
              <IconTrash v-else :size="15" :stroke-width="2.5" />
              <span>{{ deleting ? 'Eliminando…' : 'Sí, eliminar' }}</span>
            </button>
          </div>

        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
/* ── Backdrop ──────────────────────────────────────────────── */
.backdrop {
  position: fixed;
  inset: 0;
  background: rgba(26, 26, 24, 0.55);
  backdrop-filter: blur(6px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1100;
  padding: var(--space-4);
}

/* ── Dialog ────────────────────────────────────────────────── */
.dialog {
  position: relative;
  background: var(--color-surface);
  border-radius: var(--radius-xl);
  box-shadow:
    var(--shadow-xl),
    0 0 0 1px rgba(220, 38, 38, 0.08);
  width: 100%;
  max-width: 420px;
  padding: var(--space-8) var(--space-6) var(--space-6);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-4);
  text-align: center;
}

/* ── Close ─────────────────────────────────────────────────── */
.close-btn {
  position: absolute;
  top: var(--space-4);
  right: var(--space-4);
  width: 28px;
  height: 28px;
  display: flex;
  align-items: center;
  justify-content: center;
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-sm);
  background: var(--color-surface);
  color: var(--color-text-tertiary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast), border-color var(--transition-fast);
}
.close-btn:hover:not(:disabled) {
  background: var(--color-bg-alt);
  color: var(--color-text-primary);
  border-color: var(--color-border);
}
.close-btn:disabled {
  opacity: 0.4;
  cursor: not-allowed;
}

/* ── Warning icon ──────────────────────────────────────────── */
.warning-icon {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  background: var(--color-error-light);
  border: 2px solid var(--color-error-border);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-error);
  flex-shrink: 0;
}

/* ── Text ──────────────────────────────────────────────────── */
.dialog-body {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.dialog-title {
  font-family: var(--font-display);
  font-size: var(--text-xl);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.dialog-desc {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  line-height: 1.6;
  margin: 0;
}

.record-highlight {
  color: var(--color-text-primary);
  font-weight: 700;
}

/* ── Actions ───────────────────────────────────────────────── */
.dialog-actions {
  display: flex;
  gap: var(--space-3);
  width: 100%;
  margin-top: var(--space-2);
}

.btn-cancel {
  flex: 1;
  padding: var(--space-2) var(--space-4);
  background: transparent;
  color: var(--color-text-secondary);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 500;
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast);
}
.btn-cancel:hover:not(:disabled) {
  background: var(--color-bg-alt);
  border-color: var(--color-border);
}
.btn-cancel:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.btn-confirm {
  flex: 1;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-4);
  background: var(--color-error);
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast), transform var(--transition-fast);
}
.btn-confirm:hover:not(:disabled) {
  background: #b91c1c;
  transform: translateY(-1px);
}
.btn-confirm:disabled {
  opacity: 0.65;
  cursor: not-allowed;
  transform: none;
}

/* ── Spinner ───────────────────────────────────────────────── */
.btn-spinner {
  display: inline-block;
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.35);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
  flex-shrink: 0;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ── Transition ────────────────────────────────────────────── */
.confirm-modal-enter-active,
.confirm-modal-leave-active {
  transition: opacity var(--transition-fast);
}
.confirm-modal-enter-active .dialog,
.confirm-modal-leave-active .dialog {
  transition: transform var(--transition-fast), opacity var(--transition-fast);
}
.confirm-modal-enter-from,
.confirm-modal-leave-to {
  opacity: 0;
}
.confirm-modal-enter-from .dialog,
.confirm-modal-leave-to .dialog {
  transform: scale(0.95) translateY(8px);
  opacity: 0;
}
</style>
