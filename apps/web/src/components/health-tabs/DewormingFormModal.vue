<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import {
  IconX,
  IconArrowLeft,
  IconArrowRight,
  IconCheck,
  IconSearch,
  IconBell,
  IconPlus
} from '@tabler/icons-vue'
import DatePicker from '@/components/ui/DatePicker.vue'
import { useGetHealthCatalogs } from '@/composables/useHealthCatalog'
import { useCreateHealthRecord, useUpdateHealthRecord } from '@/composables/useHealthRecords'
import { HealthCatalogCategory } from '@/constants/healthRecord'
import type { HealthRecord } from '@/types'

onMounted(() => {
  document.body.style.overflow = 'hidden'
})

onUnmounted(() => {
  document.body.style.overflow = ''
})

const props = defineProps<{
  petId: string
  petSpecies: string
  editingRecord?: HealthRecord | null
}>()

const emit = defineEmits<{
  close: []
}>()

const isEditing = computed(() => !!props.editingRecord)

const currentStep = ref(1)
const totalSteps = 2

const steps = computed(() => [
  { number: 1, label: '¿Cuándo?' },
  { number: 2, label: '¿Cuándo próxima?' },
])

function nextStep() {
  if (currentStep.value < totalSteps) currentStep.value++
}

function prevStep() {
  if (currentStep.value > 1) currentStep.value--
}

const category = HealthCatalogCategory.Deworming
const categoryRef = ref(category)
const pageRef = ref(1)
const perPageRef = ref(100)
const speciesRef = computed(() => props.petSpecies)
const canFetchCatalog = computed(() => currentStep.value >= 1)

const { data: catalogResponse } = useGetHealthCatalogs(categoryRef, pageRef, perPageRef, speciesRef, canFetchCatalog)

const catalogSearch = ref('')
const selectedCatalogId = ref<string | null>(null)
const customName = ref('')

const filteredCatalog = computed(() => {
  const list = catalogResponse.value?.data || []
  const q = catalogSearch.value.toLowerCase().trim()
  if (q) {
    return list.filter((v) => v.name.toLowerCase().includes(q))
  }
  return list
})

function selectCustom() {
  selectedCatalogId.value = 'custom'
  if (catalogSearch.value.trim() && !customName.value) {
    customName.value = catalogSearch.value.trim()
  }
}

const selectedCatalog = computed(() => {
  if (selectedCatalogId.value === 'custom') return null
  return catalogResponse.value?.data?.find((v) => v.id === selectedCatalogId.value)
})

watch(catalogSearch, () => {
  if (selectedCatalogId.value === 'custom' && catalogSearch.value.trim() !== customName.value) {
    selectedCatalogId.value = null
  }
})

const applicationDate = ref('')
const nextDate = ref('')
const note = ref('')
const wantsNext = ref(false)

watch(applicationDate, (newDate) => {
  if (newDate && selectedCatalog.value?.frequency_months) {
    const d = new Date(newDate + 'T12:00:00')
    d.setMonth(d.getMonth() + selectedCatalog.value.frequency_months)
    nextDate.value = d.toISOString().split('T')[0] || ''
  }
})

watch(selectedCatalog, (val) => {
  if (!val) {
    wantsNext.value = false
    nextDate.value = ''
    return
  }
  wantsNext.value = !!val.frequency_months
  if (applicationDate.value && val.frequency_months) {
    const d = new Date(applicationDate.value + 'T12:00:00')
    d.setMonth(d.getMonth() + val.frequency_months)
    nextDate.value = d.toISOString().split('T')[0] || ''
  }
})

watch(() => props.editingRecord, (record) => {
  if (record) {
    selectedCatalogId.value = record.health_catalog_id || 'custom'
    customName.value = record.name
    applicationDate.value = record.application_date || ''
    nextDate.value = record.next_dose_date || ''
    note.value = record.notes || ''
    wantsNext.value = !!record.next_dose_date
  }
}, { immediate: true })

const canAdvance = computed(() => {
  if (currentStep.value === 1) {
    if (selectedCatalogId.value === 'custom') return customName.value.trim().length > 0
    return !!selectedCatalogId.value
  }
  if (currentStep.value === 2) {
    if (!applicationDate.value) return false
    return true
  }
  return false
})

const createRecord = useCreateHealthRecord()
const updateRecord = useUpdateHealthRecord()

async function save() {
  if (!canAdvance.value) return

  const name = selectedCatalogId.value === 'custom' ? customName.value : (selectedCatalog.value?.name || '')
  const catalogId = selectedCatalogId.value === 'custom' ? undefined : (selectedCatalogId.value ?? undefined)

  try {
    if (isEditing.value && props.editingRecord) {
      await updateRecord.mutateAsync({
        id: props.editingRecord.id,
        payload: {
          category: HealthCatalogCategory.Deworming,
          name,
          application_date: applicationDate.value || undefined,
          next_dose_date: nextDate.value || undefined,
          notes: note.value || undefined,
        },
      })
    } else {
      await createRecord.mutateAsync({
        pet_id: props.petId,
        category: HealthCatalogCategory.Deworming,
        name,
        health_catalog_id: catalogId,
        application_date: applicationDate.value || undefined,
        next_dose_date: wantsNext.value && nextDate.value ? nextDate.value : undefined,
        notes: note.value || undefined,
      })
    }

    emit('close')
  } catch(e) {
    console.error(e)
  }
}
</script>

<template>
  <Teleport to="body">
    <div class="modal-backdrop" @click.self="emit('close')">
      <div class="modal-container">
        <div class="modal-header">
          <h2>{{ isEditing ? 'Editar desparasitación' : 'Registrar desparasitación' }}</h2>
          <button class="btn-close" @click="emit('close')">
            <IconX :size="18" :stroke-width="2" />
          </button>
        </div>

        <div class="stepper-bar">
          <div
            v-for="step in steps"
            :key="step.number"
            class="step-item"
            :class="{
              'step-item--active': step.number === currentStep,
              'step-item--done': step.number < currentStep,
            }"
          >
            <div class="step-dot">
              <IconCheck v-if="step.number < currentStep" :size="12" :stroke-width="3" />
              <span v-else>{{ step.number }}</span>
            </div>
            <span class="step-label">{{ step.label }}</span>
          </div>
          <div class="step-connector">
            <div class="step-connector__fill" :style="{ width: `${((currentStep - 1) / (totalSteps - 1)) * 100}%` }" />
          </div>
        </div>

        <div class="modal-body">
          <!-- Paso 1: ¿Cuándo? -->
          <div v-if="currentStep === 1" class="step-content">
            <p class="step-instruction">
              Seleccioná el antiparasitario
            </p>

            <button
              class="catalog-option custom-option"
              :class="{ 'catalog-option--selected': selectedCatalogId === 'custom' }"
              @click="selectCustom"
            >
              <div class="catalog-option__info">
                <span class="catalog-option__name">Añadir manualmente</span>
                <span class="catalog-option__desc">Si no está en la lista</span>
              </div>
              <div v-if="selectedCatalogId === 'custom'" class="catalog-option__check">
                <IconCheck :size="16" :stroke-width="2.5" />
              </div>
              <div v-else class="catalog-option__check" style="color: var(--color-text-tertiary)">
                <IconPlus :size="16" :stroke-width="2" />
              </div>
            </button>

            <div class="catalog-search-box">
              <IconSearch class="catalog-search-icon" :size="16" :stroke-width="2" />
              <input
                v-model="catalogSearch"
                class="catalog-search-input"
                placeholder="Buscar antiparasitario..."
              />
            </div>

            <div class="catalog-list">
              <button
                v-for="item in filteredCatalog"
                :key="item.id"
                class="catalog-option"
                :class="{ 'catalog-option--selected': selectedCatalogId === item.id }"
                @click="selectedCatalogId = item.id"
              >
                <div class="catalog-option__info">
                  <span class="catalog-option__name">{{ item.name }}</span>
                  <span v-if="item.description" class="catalog-option__desc">{{ item.description }}</span>
                  <span v-if="item.frequency_months" class="catalog-option__freq">
                    Cada {{ item.frequency_months }} meses
                  </span>
                </div>
                <div v-if="selectedCatalogId === item.id" class="catalog-option__check">
                  <IconCheck :size="16" :stroke-width="2.5" />
                </div>
              </button>
            </div>

            <div v-if="selectedCatalogId === 'custom'" class="custom-input">
              <label class="field-label">Nombre *</label>
              <input
                v-model="customName"
                class="catalog-search-input"
                placeholder="Ej: Vacuna X..."
              />
            </div>
          </div>

          <!-- Paso 2: ¿Cuándo próxima? -->
          <div v-if="currentStep === 2" class="step-content">
            <p class="step-instruction">
              Indicá las fechas
            </p>

            <div v-if="selectedCatalog || customName" class="selected-info">
              <IconCheck :size="16" :stroke-width="2.5" />
              <span>{{ selectedCatalog?.name || customName }}</span>
            </div>

            <div class="date-field">
              <label class="field-label">Fecha de aplicación</label>
              <div class="date-picker-with-action">
                <DatePicker
                  v-model="applicationDate"
                  :max-date="new Date()"
                  placeholder="Seleccionar fecha"
                  unique-id="app-date"
                />
                <button
                  type="button"
                  class="btn-today"
                  @click="applicationDate = new Date().toISOString().split('T')[0] ?? ''"
                >
                  Hoy
                </button>
              </div>
            </div>

            <div class="note-field">
              <label class="field-label">Nota (opcional)</label>
              <textarea
                v-model="note"
                class="note-input"
                rows="2"
                placeholder="Ej: Lote #12345, veterinaria..."
              />
            </div>

            <div class="suggestion-card">
              <div class="suggestion-icon">
                <IconBell :size="18" :stroke-width="1.75" />
              </div>
              <div class="suggestion-body">
                <div class="suggestion-header">
                  <p class="suggestion-text">
                    <strong>Programar recordatorio</strong>
                    <template v-if="selectedCatalog?.frequency_months">
                      <br>Este antiparasitario se aplica cada <strong>{{ selectedCatalog.frequency_months }} meses</strong>.
                    </template>
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
                  <div style="margin-top: var(--space-2)">
                    <DatePicker
                      v-model="nextDate"
                      :min-date="applicationDate ? new Date(applicationDate) : new Date()"
                      placeholder="Fecha próxima aplicación"
                      unique-id="next-date"
                    />
                  </div>
                </template>
              </div>
            </div>
          </div>
        </div>

        <div class="modal-footer">
          <button v-if="currentStep > 1" class="btn-prev" @click="prevStep">
            <IconArrowLeft :size="16" :stroke-width="2" />
            Anterior
          </button>
          <div class="footer-spacer" />
          <button
            v-if="currentStep < totalSteps"
            class="btn-next"
            :disabled="!canAdvance"
            @click="nextStep"
          >
            Siguiente
            <IconArrowRight :size="16" :stroke-width="2" />
          </button>
          <button
            v-else
            class="btn-save"
            :disabled="!canAdvance || createRecord.isPending.value || updateRecord.isPending.value"
            @click="save"
          >
            <span v-if="createRecord.isPending.value || updateRecord.isPending.value" class="btn-spinner" />
            <template v-else>
              <IconCheck :size="16" :stroke-width="2.5" />
              {{ isEditing ? 'Guardar' : 'Registrar' }}
            </template>
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
/* Mismos estilos que HealthRecordFormModal */
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
  align-items: flex-end;
}

@media (min-width: 37.5em) {
  .modal-backdrop {
    align-items: center;
  }
}

.modal-container {
  container-type: inline-size;
  container-name: modal;
  background: var(--color-surface);
  width: min(560px, 100%);
  max-height: min(90vh, 700px);
  border-radius: var(--radius-xl) var(--radius-xl) 0 0;
  box-shadow: var(--shadow-xl);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

@media (min-width: 37.5em) {
  .modal-container {
    border-radius: var(--radius-xl);
    max-height: 90vh;
  }
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
  min-height: 44px;
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

.stepper-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-3) var(--space-5);
  position: relative;
  flex-shrink: 0;
}

.step-connector {
  position: absolute;
  top: 50%;
  left: calc(var(--space-5) + 14px);
  right: calc(var(--space-5) + 14px);
  height: 2px;
  background: var(--color-border-light);
  transform: translateY(-4px);
  z-index: 0;
}

.step-connector__fill {
  height: 100%;
  background: var(--color-accent);
  transition: width var(--transition-base);
}

.step-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-1);
  z-index: 1;
}

.step-dot {
  width: 28px;
  height: 28px;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-xs);
  font-weight: 700;
  border: 2px solid var(--color-border);
  background: var(--color-surface);
  color: var(--color-text-tertiary);
  transition: all var(--transition-base);
}

.step-item--active .step-dot {
  border-color: var(--color-accent);
  background: var(--color-accent);
  color: #fff;
}

.step-item--done .step-dot {
  border-color: var(--color-accent);
  background: var(--color-accent-light);
  color: var(--color-accent);
}

.step-label {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  font-weight: 500;
}

.step-item--active .step-label {
  color: var(--color-accent-dark);
  font-weight: 600;
}

.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-4) var(--space-5);
}

.step-content {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.step-instruction {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  line-height: 1.5;
  margin: 0;
}

.catalog-search-box {
  position: relative;
}

.catalog-search-icon {
  position: absolute;
  left: var(--space-3);
  top: 50%;
  transform: translateY(-50%);
  color: var(--color-text-tertiary);
  pointer-events: none;
}

.catalog-search-input {
  width: 100%;
  padding: var(--space-2) var(--space-3) var(--space-2) 2.25rem;
  min-height: 44px;
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  background: var(--color-bg);
  transition: border-color var(--transition-fast);
  box-sizing: border-box;
}

.catalog-search-input:focus {
  border-color: var(--color-accent);
  outline: none;
}

.catalog-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  max-height: 200px;
  overflow-y: auto;
}

.catalog-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-4);
  min-height: 44px;
  background: var(--color-surface);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast);
  text-align: left;
}

.catalog-option:hover {
  background: var(--color-bg-alt);
  border-color: var(--color-border);
}

.catalog-option--selected {
  background: var(--color-accent-light);
  border-color: var(--color-accent);
}

.catalog-option__info {
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
}

.catalog-option__name {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
}

.catalog-option__desc {
  font-size: var(--text-xs);
  color: var(--color-text-secondary);
}

.catalog-option__freq {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
}

.catalog-option__check {
  color: var(--color-accent);
  flex-shrink: 0;
}

.custom-option {
  background: #f8fafc;
  border: 1.5px dashed var(--color-border);
}

.custom-option:hover {
  border-color: var(--color-accent);
}

.custom-input {
  margin-top: var(--space-2);
}

.field-label {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-bottom: var(--space-2);
  display: block;
}

.date-picker-with-action {
  display: flex;
  gap: var(--space-2);
  align-items: flex-start;
  flex-wrap: wrap;
}

.btn-today {
  padding: var(--space-2) var(--space-3);
  min-height: 44px;
  background: var(--color-accent-light);
  color: var(--color-accent-dark);
  border: 1.5px solid var(--color-accent);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
}

.btn-today:hover {
  background: var(--color-accent);
  color: #fff;
}

.note-input {
  padding: var(--space-3);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  background: var(--color-bg);
  resize: vertical;
  min-height: 60px;
  font-family: var(--font-body);
  width: 100%;
  box-sizing: border-box;
}

.note-input:focus {
  border-color: var(--color-accent);
  outline: none;
}

.selected-info {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3);
  background: var(--color-surface-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
}

.selected-info svg {
  color: var(--color-success, #10b981);
}

.suggestion-card {
  display: flex;
  gap: var(--space-3);
  padding: var(--space-4);
  background: #f0f9ff;
  border: 1px solid #bae6fd;
  border-radius: var(--radius-lg);
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
  min-height: 44px;
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

.btn-prev {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-3);
  min-height: 44px;
  background: var(--color-surface);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
  font-weight: 500;
  cursor: pointer;
}

.btn-prev:hover {
  background: var(--color-bg-alt);
}

.btn-next {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-5);
  min-height: 44px;
  background: var(--color-accent);
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
}

.btn-next:hover:not(:disabled) {
  background: var(--color-accent-hover);
}

.btn-next:disabled {
  background: var(--color-border);
  color: var(--color-text-tertiary);
  cursor: not-allowed;
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
