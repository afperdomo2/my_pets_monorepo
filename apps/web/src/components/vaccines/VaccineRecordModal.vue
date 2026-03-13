<script setup lang="ts">
import { ref, computed } from 'vue'
import {
  IconX,
  IconArrowLeft,
  IconArrowRight,
  IconCheck,
  IconSearch,
  IconCalendar,
  IconBell,
} from '@tabler/icons-vue'

const props = defineProps<{
  /** Mascota preseleccionada (si viene desde una tarjeta del timeline) */
  preselectedPet?: string
}>()

const emit = defineEmits<{
  close: []
}>()

// ── Stepper ─────────────────────────────────────
const currentStep = ref(props.preselectedPet ? 2 : 1)
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
  if (currentStep.value > 1) currentStep.value--
}

// ── Paso 1: Seleccionar mascota ─────────────────
const mockPets = [
  { id: '1', name: 'Romeo', species: 'dog', initials: 'RO', lifeStage: 'Adulto' },
  { id: '2', name: 'Luna', species: 'cat', initials: 'LU', lifeStage: 'Joven' },
  { id: '3', name: 'Simba', species: 'cat', initials: 'SI', lifeStage: 'Adulto' },
  { id: '4', name: 'Manchas', species: 'dog', initials: 'MA', lifeStage: 'Senior' },
  { id: '5', name: 'Bolt', species: 'dog', initials: 'BO', lifeStage: 'Cachorro' },
  { id: '6', name: 'Kira', species: 'cat', initials: 'KI', lifeStage: 'Adulto' },
]

const selectedPetId = ref<string | null>(props.preselectedPet ?? null)
const selectedPet = computed(() => mockPets.find((p) => p.id === selectedPetId.value))

// ── Paso 2: Seleccionar vacuna ──────────────────
const vaccineSearch = ref('')
const mockVaccines = [
  { id: '1', name: 'Antirrábica', frequency: 'Anual', species: ['dog', 'cat'] },
  { id: '2', name: 'Polivalente canina', frequency: 'Anual', species: ['dog'] },
  { id: '3', name: 'Triple felina', frequency: 'Anual', species: ['cat'] },
  { id: '4', name: 'Sextuple', frequency: 'Anual', species: ['dog'] },
  { id: '5', name: 'Leucemia felina', frequency: 'Anual', species: ['cat'] },
  { id: '6', name: 'Bordetelosis', frequency: 'Cada 6 meses', species: ['dog'] },
]

const filteredVaccines = computed(() => {
  let list = mockVaccines
  // Filtrar por especie de la mascota seleccionada
  if (selectedPet.value) {
    list = list.filter((v) => v.species.includes(selectedPet.value!.species))
  }
  // Filtrar por texto de búsqueda
  const q = vaccineSearch.value.toLowerCase().trim()
  if (q) {
    list = list.filter((v) => v.name.toLowerCase().includes(q))
  }
  return list
})

const selectedVaccineId = ref<string | null>(null)
const selectedVaccine = computed(() => mockVaccines.find((v) => v.id === selectedVaccineId.value))

// ── Paso 3: Fecha y nota ────────────────────────
const applicationDate = ref('')
const note = ref('')
const wantsReminder = ref(true)

// ── Validación por paso ─────────────────────────
const canAdvance = computed(() => {
  if (currentStep.value === 1) return !!selectedPetId.value
  if (currentStep.value === 2) return !!selectedVaccineId.value
  if (currentStep.value === 3) return !!applicationDate.value
  return false
})
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
                  v-for="pet in mockPets"
                  :key="pet.id"
                  class="pet-option"
                  :class="{
                    'pet-option--selected': selectedPetId === pet.id,
                    [`pet-option--${pet.species}`]: true,
                  }"
                  @click="selectedPetId = pet.id"
                >
                  <div class="pet-option__avatar">{{ pet.initials }}</div>
                  <span class="pet-option__name">{{ pet.name }}</span>
                  <span class="pet-option__stage">{{ pet.lifeStage }}</span>
                </button>
              </div>
            </div>

            <!-- Paso 2: ¿Qué? -->
            <div v-if="currentStep === 2" class="step-content">
              <p class="step-instruction">
                Seleccioná la vacuna aplicada
                <span v-if="selectedPet" class="step-instruction__pet">a {{ selectedPet.name }}</span>
              </p>
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
                    <span class="vaccine-option__name">{{ vaccine.name }}</span>
                    <span class="vaccine-option__freq">{{ vaccine.frequency }}</span>
                  </div>
                  <div v-if="selectedVaccineId === vaccine.id" class="vaccine-option__check">
                    <IconCheck :size="16" :stroke-width="2.5" />
                  </div>
                </button>
              </div>
              <p v-if="filteredVaccines.length === 0" class="no-results">
                No se encontraron vacunas
              </p>
            </div>

            <!-- Paso 3: ¿Cuándo? -->
            <div v-if="currentStep === 3" class="step-content">
              <p class="step-instruction">
                Indicá la fecha de aplicación
                <span v-if="selectedVaccine" class="step-instruction__pet">
                  de {{ selectedVaccine.name }}
                </span>
              </p>

              <div class="date-field">
                <label class="field-label">
                  <IconCalendar :size="14" :stroke-width="2" />
                  Fecha de aplicación
                </label>
                <input v-model="applicationDate" type="date" class="date-input" />
              </div>

              <div class="note-field">
                <label class="field-label">Nota (opcional)</label>
                <textarea
                  v-model="note"
                  class="note-input"
                  rows="3"
                  placeholder="Ej: Aplicada por Dr. Pérez en clínica veterinaria…"
                />
              </div>

              <!-- Sugerencia inteligente -->
              <div v-if="selectedVaccine" class="suggestion-card">
                <div class="suggestion-icon">
                  <IconBell :size="18" :stroke-width="1.75" />
                </div>
                <div class="suggestion-body">
                  <p class="suggestion-text">
                    <strong>{{ selectedVaccine.name }}</strong> suele aplicarse
                    <strong>{{ selectedVaccine.frequency?.toLowerCase() }}</strong>.
                    ¿Querés que te avise cuando corresponda la próxima dosis?
                  </p>
                  <label class="suggestion-toggle">
                    <input v-model="wantsReminder" type="checkbox" class="toggle-checkbox" />
                    <span class="toggle-track">
                      <span class="toggle-thumb" />
                    </span>
                    <span class="toggle-label">{{ wantsReminder ? 'Sí, avisame' : 'No, gracias' }}</span>
                  </label>
                </div>
              </div>
            </div>
          </div>

          <!-- Footer -->
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
              :disabled="!canAdvance"
              @click="emit('close')"
            >
              <IconCheck :size="16" :stroke-width="2.5" />
              Guardar
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
  max-width: 520px;
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

.pet-option__avatar {
  width: 48px;
  height: 48px;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-sm);
  font-weight: 700;
  background: var(--color-bg-alt);
  color: var(--color-text-secondary);
  transition: background var(--transition-fast);
}

.pet-option--dog .pet-option__avatar    { background: #fef3c7; color: #92400e; }
.pet-option--cat .pet-option__avatar    { background: #ede9fe; color: #5b21b6; }
.pet-option--bird .pet-option__avatar   { background: #cffafe; color: #0e7490; }
.pet-option--rabbit .pet-option__avatar { background: #fce7f3; color: #9d174d; }

.pet-option__name {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
}

.pet-option__stage {
  font-size: 0.65rem;
  font-weight: 500;
  padding: 1px var(--space-2);
  border-radius: var(--radius-full);
  background: var(--color-bg-alt);
  color: var(--color-text-tertiary);
  border: 1px solid var(--color-border-light);
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
}

.vaccine-list {
  display: flex;
  flex-direction: column;
  gap: 2px;
  max-height: 240px;
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

.date-input {
  padding: var(--space-2) var(--space-3);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  background: var(--color-bg);
  transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
}

.date-input:focus {
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(61, 122, 95, 0.12);
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
  gap: var(--space-3);
}

.suggestion-text {
  font-size: var(--text-sm);
  color: #1E40AF;
  line-height: 1.5;
}

/* ── Toggle switch ──────────────────── */
.suggestion-toggle {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  cursor: pointer;
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

/* ── Modal transitions ──────────────── */
.modal-enter-active {
  transition: opacity 200ms ease;
}

.modal-enter-active .modal-container {
  transition: opacity 200ms ease, transform 200ms cubic-bezier(0.34, 1.56, 0.64, 1);
}

.modal-leave-active {
  transition: opacity 150ms ease;
}

.modal-leave-active .modal-container {
  transition: opacity 150ms ease, transform 150ms ease;
}

.modal-enter-from {
  opacity: 0;
}

.modal-enter-from .modal-container {
  opacity: 0;
  transform: scale(0.95) translateY(10px);
}

.modal-leave-to {
  opacity: 0;
}

.modal-leave-to .modal-container {
  opacity: 0;
  transform: scale(0.97) translateY(5px);
}
</style>
