<script setup lang="ts">
import { ref, computed, watch } from 'vue'
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
import { useGetPets } from '@/composables/usePets'
import { useGetHealthCatalogs } from '@/composables/useHealthCatalog'
import { useCreateHealthRecord } from '@/composables/useHealthRecords'
import PetAvatar from '@/components/pets/PetAvatar.vue'
import { onMounted, onUnmounted } from 'vue'

onMounted(() => {
  document.body.style.overflow = 'hidden'
})

onUnmounted(() => {
  document.body.style.overflow = ''
})

const props = defineProps<{
  /** Mascota preseleccionada (si viene desde una tarjeta del timeline) */
  preselectedPet?: string
  /** Nombre de la vacuna preseleccionada (para buscar y seleccionar automáticamente) */
  preselectedVaccine?: string
}>()

const emit = defineEmits<{
  close: []
}>()

// ── Stepper ─────────────────────────────────────
const currentStep = ref(props.preselectedPet && props.preselectedVaccine ? 3 : props.preselectedPet ? 2 : 1)
const totalSteps = 3

const steps = [
  { number: 1, label: '¿Quién?' },
  { number: 2, label: '¿Qué?' },
  { number: 3, label: '¿Cuándo?' },
]

function nextStep() {
  if (currentStep.value < totalSteps) currentStep.value++
}

function prevStep() {
  if (currentStep.value > 1) {
    if (props.preselectedPet && currentStep.value === 2) return
    if (props.preselectedPet && props.preselectedVaccine && currentStep.value === 3) return
    currentStep.value--
  }
}

// ── Paso 1: Seleccionar mascota ─────────────────
const { data: allPets } = useGetPets()

const selectedPetId = ref<string | null>(props.preselectedPet ?? null)
const selectedPet = computed(() => allPets.value?.find((p) => p.id === selectedPetId.value))

watch(() => allPets.value, (pets) => {
  if (props.preselectedPet && pets && pets.length > 0) {
    selectedPetId.value = props.preselectedPet
  }
}, { immediate: true })

// ── Paso 2: Seleccionar vacuna ──────────────────
const categoryRef = ref('vaccine')
const pageRef = ref(1)
const perPageRef = ref(100)
const speciesRef = computed(() => selectedPet.value?.species)
const canFetchVaccines = computed(() => currentStep.value >= 2 && !!speciesRef.value)

const { data: catalogResponse } = useGetHealthCatalogs(categoryRef, pageRef, perPageRef, speciesRef, canFetchVaccines)

const vaccineSearch = ref('')
const selectedVaccineId = ref<string | null>(null)
const customVaccineName = ref('')

const filteredVaccines = computed(() => {
  const list = catalogResponse.value?.data || []
  const q = vaccineSearch.value.toLowerCase().trim()
  if (q) {
    return list.filter((v) => v.name.toLowerCase().includes(q))
  }
  return list
})

function selectCustomVaccine() {
  selectedVaccineId.value = 'custom'
  if (vaccineSearch.value.trim() && !customVaccineName.value) {
    customVaccineName.value = vaccineSearch.value.trim()
  }
}

const selectedVaccine = computed(() => {
  if (selectedVaccineId.value === 'custom') return null
  return catalogResponse.value?.data?.find((v) => v.id === selectedVaccineId.value)
})

watch(vaccineSearch, () => {
  if (selectedVaccineId.value === 'custom' && vaccineSearch.value.trim() !== customVaccineName.value) {
    selectedVaccineId.value = null
  }
})

// Seleccionar automáticamente la vacuna preseleccionada por nombre
watch([catalogResponse, () => props.preselectedVaccine], ([catalog, preselectedVaccine]) => {
  if (!preselectedVaccine || !catalog?.data) return
  
  // Buscar la vacuna por nombre
  const vaccine = catalog.data.find((v) => v.name.toLowerCase() === preselectedVaccine.toLowerCase())
  if (vaccine) {
    selectedVaccineId.value = vaccine.id
  }
}, { immediate: true })

// ── Paso 3: Fecha y nota ────────────────────────
const applicationDate = ref('')
const boosterDate = ref('')
const note = ref('')
const wantsReminder = ref(false)

watch(applicationDate, (newDate) => {
  if (newDate) {
     const d = new Date(newDate + 'T12:00:00')
     d.setFullYear(d.getFullYear() + 1) // Default to 1 year for manual entries or existing records without catalog details
     boosterDate.value = d.toISOString().split('T')[0] || ''
  }
})

watch(selectedVaccine, (val) => {
   if (!val) {
     wantsReminder.value = false
     boosterDate.value = ''
     return
   }
   // Defaults it to false but calculate the future date anyway
   wantsReminder.value = false
   if (applicationDate.value && val.frequency_months) {
     const d = new Date(applicationDate.value + 'T12:00:00')
     d.setMonth(d.getMonth() + val.frequency_months)
     boosterDate.value = d.toISOString().split('T')[0] || ''
   }
})

// ── Validación por paso ─────────────────────────
const canAdvance = computed(() => {
  if (currentStep.value === 1) return !!selectedPetId.value
  if (currentStep.value === 2) {
    if (selectedVaccineId.value === 'custom') return customVaccineName.value.trim().length > 0
    return !!selectedVaccineId.value
  }
  if (currentStep.value === 3) {
    if (!applicationDate.value) return false
    if (wantsReminder.value && !boosterDate.value) return false
    return true
  }
  return false
})

// ── Solicitud ───────────────────────────────────
const createHealthRecord = useCreateHealthRecord()

async function save() {
  if (!canAdvance.value || createHealthRecord.isPending.value) return
  
  try {
     const name = selectedVaccineId.value === 'custom' ? customVaccineName.value : (selectedVaccine.value?.name || '')
     const catalogId = selectedVaccineId.value === 'custom' ? undefined : (selectedVaccineId.value ?? undefined)

     const basePayload = {
       pet_id: selectedPetId.value!,
       category: 'vaccine',
       name,
       health_catalog_id: catalogId,
       notes: note.value || undefined,
     }

     // 1. Aplicada
     await createHealthRecord.mutateAsync({
       ...basePayload,
       status: 'applied',
       application_date: applicationDate.value,
       due_date: applicationDate.value,
     })

     // 2. Refuerzo
     if (wantsReminder.value && boosterDate.value) {
       await createHealthRecord.mutateAsync({
         ...basePayload,
         status: 'pending',
         due_date: boosterDate.value,
         application_date: undefined
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
          <!-- Header -->
          <div class="modal-header">
            <h2>Registrar vacuna</h2>
            <button class="btn-close" @click="emit('close')">
               <IconX :size="18" :stroke-width="2" />
            </button>
          </div>

          <!-- Stepper indicator -->
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

          <!-- Step content -->
          <div class="modal-body">
            <!-- Paso 1: ¿Quién? -->
            <div v-if="currentStep === 1" class="step-content">
              <p class="step-instruction">Seleccioná la mascota que recibió la vacuna</p>
              <div class="pets-grid">
                <button
                  v-for="pet in allPets"
                  :key="pet.id"
                  class="pet-option"
                  :class="{
                    'pet-option--selected': selectedPetId === pet.id,
                    [`pet-option--${pet.species}`]: true,
                  }"
                  @click="selectedPetId = pet.id"
                >
                  <PetAvatar :species="pet.species" :name="pet.name" size="lg" />
                  <span class="pet-option__name">{{ pet.name }}</span>
                </button>
              </div>
            </div>

            <!-- Paso 2: ¿Qué? -->
            <div v-if="currentStep === 2" class="step-content">
              <p class="step-instruction">
                Seleccioná la vacuna aplicada
                <span v-if="selectedPet" class="step-instruction__pet">a {{ selectedPet.name }}</span>
              </p>

              <!-- Option for Manual Vaccine -->
              <button
                class="vaccine-option custom-vaccine"
                :class="{ 'vaccine-option--selected': selectedVaccineId === 'custom' }"
                @click="selectCustomVaccine"
                style="margin-bottom: var(--space-2);"
              >
                <div class="vaccine-option__info">
                  <span class="vaccine-option__name">Añadir otra vacuna manualmente</span>
                  <span class="vaccine-option__freq">Si no encuentras la vacuna en la lista</span>
                </div>
                <div v-if="selectedVaccineId === 'custom'" class="vaccine-option__check">
                    <IconCheck :size="16" :stroke-width="2.5" />
                </div>
                <div v-else class="vaccine-option__check" style="color: var(--color-text-tertiary)">
                    <IconPlus :size="16" :stroke-width="2" />
                </div>
              </button>

              <div class="vaccine-search-box">
                <IconSearch class="vaccine-search-icon" :size="16" :stroke-width="2" />
                <input
                  v-model="vaccineSearch"
                  class="vaccine-search-input"
                  placeholder="Buscar vacuna…"
                />
              </div>
              <div class="vaccine-list">
                <button
                  v-for="vaccine in filteredVaccines"
                  :key="vaccine.id"
                  class="vaccine-option"
                  :class="{ 'vaccine-option--selected': selectedVaccineId === vaccine.id }"
                  @click="selectedVaccineId = vaccine.id"
                >
                  <div class="vaccine-option__info">
                    <div style="display: flex; align-items: center; gap: 6px;">
                      <span class="vaccine-option__name">{{ vaccine.name }}</span>
                      <span v-if="vaccine.is_mandatory" class="badge-mandatory">Obligatoria</span>
                    </div>
                    <span v-if="vaccine.description" class="vaccine-option__desc">{{ vaccine.description }}</span>
                    <span class="vaccine-option__freq">
                      {{ vaccine.frequency_months ? `Se recomienda cada ${vaccine.frequency_months} meses` : 'Sin frecuencia' }}
                    </span>
                  </div>
                  <div v-if="selectedVaccineId === vaccine.id" class="vaccine-option__check">
                    <IconCheck :size="16" :stroke-width="2.5" />
                  </div>
                </button>
              </div>
              <p v-if="filteredVaccines.length === 0" class="no-results">
                 No se encontraron vacunas en el catálogo. Seleccioná añadir manualmente.
              </p>

              <!-- Input para nombre de vacuna manual -->
              <div v-if="selectedVaccineId === 'custom'" style="margin-top: 16px;">
                <label class="field-label">Nombre de la vacuna <span style="color:red; margin-left: 4px;">*</span></label>
                <input
                  v-model="customVaccineName"
                  class="vaccine-search-input"
                  style="padding-left: 12px; border-color: var(--color-accent);"
                  placeholder="Ej: Antirrábica..."
                />
              </div>
            </div>

            <!-- Paso 3: ¿Cuándo? -->
            <div v-if="currentStep === 3" class="step-content">
              <p class="step-instruction">
                Indicá la fecha de aplicación
                <span v-if="selectedVaccine || customVaccineName" class="step-instruction__pet">
                  de {{ selectedVaccine?.name || customVaccineName }}
                </span>
                <span v-if="selectedPet" class="step-instruction__pet">
                  a {{ selectedPet.name }}
                </span>
              </p>

              <!-- Mostrar vacuna seleccionada con opción a cambiar -->
              <div v-if="selectedVaccine || customVaccineName" class="selected-vaccine-info">
                <div class="selected-vaccine-info__content">
                  <IconCheck class="selected-vaccine-info__icon" :size="18" :stroke-width="2.5" />
                  <div class="selected-vaccine-info__text">
                    <span class="selected-vaccine-info__label">Vacuna seleccionada:</span>
                    <span class="selected-vaccine-info__name">{{ selectedVaccine?.name || customVaccineName }}</span>
                  </div>
                </div>
                <button 
                  v-if="!props.preselectedVaccine"
                  class="btn-change-vaccine" 
                  type="button"
                  @click="currentStep = 2"
                >
                  Cambiar
                </button>
              </div>

              <!-- Calendario para Fecha de Aplicación -->
              <div class="date-field">
                <label class="field-label">Fecha de aplicación <span style="color:red; margin-left: 4px;">*</span></label>
                <div class="date-picker-with-action">
                  <DatePicker
                    v-model="applicationDate"
                    :max-date="new Date()"
                    placeholder="Seleccionar fecha"
                    unique-id="application-date"
                  />
                  <button
                    type="button"
                    class="btn-today"
                    title="Establecer fecha de hoy"
                    @click="applicationDate = new Date().toISOString().split('T')[0]"
                  >
                    Hoy
                  </button>
                </div>
              </div>

              <!-- Calendario para Refuerzo, condicionado al toggle -->
              <div class="date-field" style="margin-top: var(--space-2)">
                  <div class="suggestion-card">
                    <div class="suggestion-icon">
                      <IconBell :size="18" :stroke-width="1.75" />
                    </div>
                    <div class="suggestion-body">
                      <div class="suggestion-header">
                        <p class="suggestion-text">
                          <strong>¿Programar refuerzo?</strong>
                          <template v-if="selectedVaccine && selectedVaccine.frequency_months">
                            <br>Esta vacuna tiene un intervalo de <strong>{{ selectedVaccine.frequency_months }} meses</strong>.
                          </template>
                        </p>
                        <label class="suggestion-toggle">
                          <input v-model="wantsReminder" type="checkbox" class="toggle-checkbox" />
                          <span class="toggle-track">
                            <span class="toggle-thumb" />
                          </span>
                          <span class="toggle-label">{{ wantsReminder ? 'Sí' : 'No' }}</span>
                        </label>
                      </div>

                      <template v-if="wantsReminder">
                         <div style="margin-top: 8px">
                           <DatePicker
                             v-model="boosterDate"
                             placeholder="Fecha del refuerzo"
                             unique-id="booster-date"
                           />
                         </div>
                      </template>
                    </div>
                  </div>
              </div>

              <div class="note-field" style="margin-top: var(--space-1)">
                <label class="field-label">Nota (opcional)</label>
                <textarea
                  v-model="note"
                  class="note-input"
                  rows="2"
                  placeholder="Ej: Lote #12345, veterinaria central..."
                />
              </div>

            </div>
          </div>

          <!-- Footer -->
          <div class="modal-footer">
            <button v-if="currentStep > 1 && !(props.preselectedPet && currentStep === 2)" class="btn-prev" @click="prevStep">
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
              :disabled="!canAdvance || createHealthRecord.isPending.value"
              @click="save"
            >
              <span v-if="createHealthRecord.isPending.value" class="btn-spinner" />
              <template v-else>
                 <IconCheck :size="16" :stroke-width="2.5" />
                 Guardar
              </template>
            </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
/* ── Backdrop ───────────────────────── */
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

/* ── Container ──────────────────────── */
.modal-container {
  background: var(--color-surface);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-xl);
  width: 100%;
  max-width: 600px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* ── Header ─────────────────────────── */
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-5) var(--space-6);
  border-bottom: 1px solid var(--color-border-light);
}

.modal-header h2 {
  font-family: var(--font-display);
  font-size: var(--text-lg);
  font-weight: 600;
  color: var(--color-text-primary);
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
  color: var(--color-text-tertiary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.btn-close:hover {
  background: var(--color-bg-alt);
  color: var(--color-text-primary);
}

/* ── Stepper bar ────────────────────── */
.stepper-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-4) var(--space-6);
  position: relative;
}

.step-connector {
  position: absolute;
  top: 50%;
  left: calc(var(--space-6) + 14px);
  right: calc(var(--space-6) + 14px);
  height: 2px;
  background: var(--color-border-light);
  transform: translateY(-4px);
  z-index: 0;
}

.step-connector__fill {
  height: 100%;
  background: var(--color-accent);
  transition: width var(--transition-base);
  border-radius: 1px;
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
  transition: color var(--transition-fast);
}

.step-item--active .step-label {
  color: var(--color-accent-dark);
  font-weight: 600;
}

/* ── Body ───────────────────────────── */
.modal-body {
  flex: 1;
  overflow-y: auto;
  padding: var(--space-5) var(--space-6);
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

.step-instruction__pet {
  color: var(--color-accent-dark);
  font-weight: 600;
}

/* ── Paso 1: Grid mascotas ──────────── */
.pets-grid {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: var(--space-3);
}

@media (max-width: 480px) {
  .pets-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.pet-option {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-4) var(--space-3);
  background: var(--color-surface);
  border: 2px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  cursor: pointer;
  transition: border-color var(--transition-fast), box-shadow var(--transition-fast),
    transform var(--transition-fast);
}

.pet-option:hover {
  border-color: var(--color-border);
  box-shadow: var(--shadow-sm);
  transform: translateY(-1px);
}

.pet-option--selected {
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(61, 122, 95, 0.15);
  background: var(--color-accent-light);
}

.pet-option__name {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
}

/* ── Paso 2: Búsqueda vacunas ───────── */
.vaccine-search-box {
  position: relative;
}

.vaccine-search-icon {
  position: absolute;
  left: var(--space-3);
  top: 50%;
  transform: translateY(-50%);
  color: var(--color-text-tertiary);
  pointer-events: none;
}

.vaccine-search-input {
  width: 100%;
  padding: var(--space-2) var(--space-3) var(--space-2) 2.25rem;
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  background: var(--color-bg);
  transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
  box-sizing: border-box;
}

.vaccine-search-input:focus {
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(61, 122, 95, 0.12);
  outline: none;
}

.vaccine-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  max-height: 280px;
  overflow-y: auto;
}

.vaccine-option {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-3) var(--space-4);
  background: var(--color-surface);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast);
  text-align: left;
}

.vaccine-option:hover {
  background: var(--color-bg-alt);
  border-color: var(--color-border);
}

.vaccine-option--selected {
  background: var(--color-accent-light);
  border-color: var(--color-accent);
}

.vaccine-option__info {
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.vaccine-option__name {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
}

.vaccine-option__freq {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
}

.vaccine-option__check {
  color: var(--color-accent);
}

.custom-vaccine {
  background: #f8fafc;
  border: 1.5px dashed var(--color-border);
}

.custom-vaccine:hover {
  border-color: var(--color-accent);
}

.badge-mandatory {
  font-size: 0.65rem;
  font-weight: 700;
  text-transform: uppercase;
  background: #FEF2F2;
  color: #DC2626;
  border: 1px solid #FECACA;
  padding: 2px 6px;
  border-radius: var(--radius-sm);
}

.vaccine-option__desc {
  font-size: var(--text-xs);
  color: var(--color-text-secondary);
  margin-top: 2px;
  line-height: 1.3;
}

.no-results {
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
  text-align: center;
  padding: var(--space-4);
}

/* ── Paso 3: Campos ─────────────────── */
.field-label {
  display: flex;
  align-items: center;
  gap: var(--space-1);
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-bottom: var(--space-2);
}

.date-field,
.note-field {
  display: flex;
  flex-direction: column;
}

.date-picker-with-action {
  display: flex;
  gap: var(--space-2);
  align-items: flex-start;
}

.date-picker-with-action :deep(.date-picker) {
  flex: 1;
}

.btn-today {
  padding: var(--space-2) var(--space-3);
  background: var(--color-accent-light);
  color: var(--color-accent-dark);
  border: 1.5px solid var(--color-accent);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast),
    color var(--transition-fast);
  white-space: nowrap;
  flex-shrink: 0;
}

.btn-today:hover {
  background: var(--color-accent);
  color: #fff;
  border-color: var(--color-accent);
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
  transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
  font-family: var(--font-body);
}

.note-input:focus {
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(61, 122, 95, 0.12);
  outline: none;
}

/* ── Sugerencia ─────────────────────── */
.suggestion-card {
  display: flex;
  gap: var(--space-3);
  padding: var(--space-4);
  background: #F0F9FF;
  border: 1px solid #BAE6FD;
  border-radius: var(--radius-lg);
}

.suggestion-icon {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-md);
  background: #DBEAFE;
  color: #2563EB;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.suggestion-body {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
}

.suggestion-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--space-2);
}

.suggestion-text {
  font-size: var(--text-sm);
  color: #1E40AF;
  line-height: 1.5;
  margin: 0;
}

/* ── Toggle switch ──────────────────── */
.suggestion-toggle {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  cursor: pointer;
  flex-shrink: 0;
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
  color: #1E40AF;
}

/* ── Footer ─────────────────────────── */
.modal-footer {
  display: flex;
  align-items: center;
  padding: var(--space-4) var(--space-6);
  border-top: 1px solid var(--color-border-light);
  gap: var(--space-3);
}

.footer-spacer {
  flex: 1;
}

.btn-prev {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-3);
  background: var(--color-surface);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  color: var(--color-text-secondary);
  font-size: var(--text-sm);
  font-weight: 500;
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast);
}

.btn-prev:hover {
  background: var(--color-bg-alt);
  border-color: var(--color-border);
}

.btn-next {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-5);
  background: var(--color-accent);
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast), transform var(--transition-fast);
}

.btn-next:hover:not(:disabled) {
  background: var(--color-accent-hover);
  transform: translateY(-1px);
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
  background: #2E7D52;
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast), transform var(--transition-fast);
}

.btn-save:hover:not(:disabled) {
  background: #256644;
  transform: translateY(-1px);
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

/* ── Selected vaccine info (Step 3) ───────────────────────────────────── */
.selected-vaccine-info {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-3);
  padding: var(--space-3);
  background: var(--color-surface-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  margin-bottom: var(--space-4);
}

.selected-vaccine-info__content {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.selected-vaccine-info__icon {
  color: var(--color-success);
  flex-shrink: 0;
}

.selected-vaccine-info__text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.selected-vaccine-info__label {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
}

.selected-vaccine-info__name {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
}

.btn-change-vaccine {
  padding: var(--space-1) var(--space-3);
  background: transparent;
  color: var(--color-accent);
  border: 1px solid var(--color-accent);
  border-radius: var(--radius-sm);
  font-size: var(--text-xs);
  font-weight: 500;
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
  white-space: nowrap;
}

.btn-change-vaccine:hover {
  background: var(--color-accent);
  color: #fff;
}
</style>
