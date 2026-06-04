<script setup lang="ts">
import VaccineApplicationsList from '@/components/health-tabs/VaccineApplicationsList.vue'
import { useGetVaccineApplicationsByHealthRecord } from '@/composables/useVaccineApplications'
import { IconX } from '@tabler/icons-vue'
import { onMounted, onUnmounted, toRef } from 'vue'

const props = defineProps<{
  healthRecordId: string
  category: 'vaccine' | 'deworming'
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

const recordIdRef = toRef(props, 'healthRecordId')
const { data: applications, isLoading } = useGetVaccineApplicationsByHealthRecord(recordIdRef)
</script>

<template>
  <Teleport to="body">
    <div class="modal-backdrop" @click.self="emit('close')">
      <div class="modal-container">
        <!-- Header -->
        <div class="modal-header">
          <h2>{{ category === 'vaccine' ? 'Dosis aplicadas' : 'Dosis aplicadas' }}</h2>
          <button class="btn-close" @click="emit('close')">
            <IconX :size="18" :stroke-width="2" />
          </button>
        </div>

        <!-- Body -->
        <div class="modal-body">
          <div v-if="isLoading" class="loading-state">
            <div class="spinner" />
            <p>Cargando aplicaciones...</p>
          </div>

          <div v-else-if="applications && applications.length > 0" class="applications-wrapper">
            <VaccineApplicationsList
              :applications="applications"
              :category="props.category"
            />
          </div>

          <div v-else class="empty-state">
            <p>No hay dosis aplicadas registradas</p>
          </div>
        </div>

        <!-- Footer -->
        <div class="modal-footer">
          <div class="footer-spacer" />
          <button class="btn-close-modal" @click="emit('close')">
            Cerrar
          </button>
        </div>
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

.modal-body {
  padding: var(--space-5);
  overflow-y: auto;
  flex: 1;
}

.applications-wrapper {
  margin-top: 0;
}

.applications-wrapper :deep(.vaccine-applications-list) {
  margin-top: 0;
  padding-top: 0;
  border-top: none;
}

.loading-state,
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

.loading-state p,
.empty-state p {
  margin: 0;
  font-size: var(--text-sm);
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
  background: var(--color-bg-alt);
  color: var(--color-text-secondary);
  border: 1.5px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast);
}

.btn-close-modal:hover {
  background: var(--color-accent-light);
  border-color: var(--color-accent);
  color: var(--color-accent-dark);
}
</style>
