<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { IconX } from '@tabler/icons-vue'
import DatePicker from '@/components/ui/DatePicker.vue'
import { useCreatePet, useUpdatePet } from '@/composables/usePets'
import { createPetSchema, updatePetSchema } from '@/schemas/pet'
import { estimatedBirthDate, toGrams } from '@/utils/pet'
import { PET_SPECIES } from '@/constants/species'
import { PET_SIZE_VALUES, PET_SIZE_LABELS, PET_SIZE_DESCRIPTIONS, PET_SIZE_ICONS } from '@/constants/petSize'
import type { Pet } from '@/types/pet'
import type { PetSize } from '@/constants/petSize'
import type { CreatePetFormValues } from '@/schemas/pet'

const props = defineProps<{
  modelValue: boolean
  mode: 'create' | 'edit'
  pet?: Pet
}>()

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
  success: []
}>()

const createPet = useCreatePet()
const updatePet = useUpdatePet()

const isPending = () =>
  props.mode === 'create' ? createPet.isPending.value : updatePet.isPending.value

// ── Date mode toggle ────────────────────────────────────────────────────────
const dateMode = ref<'exact' | 'age'>('age')

// Fields for "age" mode
const ageYears = ref(0)
const ageMonths = ref(0)

// ── Weight ──────────────────────────────────────────────────────────────────
const weightUnit = ref<'kg' | 'g'>('kg')
const weightInput = ref<number | ''>('')

// Computed weight in grams for sending to backend
const weightGrams = computed<number | null>(() => {
  if (weightInput.value === '' || weightInput.value <= 0) return null
  return toGrams(Number(weightInput.value), weightUnit.value)
})

// ── Form ─────────────────────────────────────────────────────────────────────
const { defineField, handleSubmit, errors, resetForm, setErrors, setFieldValue, values } = useForm<CreatePetFormValues>({
  validationSchema: computed(() =>
    toTypedSchema(props.mode === 'create' ? createPetSchema : updatePetSchema),
  ),
  initialValues: {
    name: '',
    species: '',
    breed: '',
    birth_date: estimatedBirthDate(0, 0),
    birth_date_exact: false,
    weight_grams: undefined,
    size: undefined,
  },
})

const [name, nameAttrs] = defineField('name')
const [species, speciesAttrs] = defineField('species')
const [breed, breedAttrs] = defineField('breed')
const [birthDate] = defineField('birth_date')

// Clear size when species changes away from 'dog'
watch(() => values.species, (newSpecies) => {
  if (newSpecies !== 'dog') {
    setFieldValue('size', undefined)
  }
})
// Sync dateMode toggle → birth_date_exact field
watch(dateMode, (mode) => {
  setFieldValue('birth_date_exact', mode === 'exact')
  if (mode === 'age') {
    setFieldValue('birth_date', estimatedBirthDate(ageYears.value, ageMonths.value))
  }
})

// Sync age inputs → birth_date field (in age mode)
watch([ageYears, ageMonths], ([y, m]) => {
  if (dateMode.value === 'age') {
    setFieldValue('birth_date', estimatedBirthDate(y, m))
  }
})

// ── Pre-fill form when editing ───────────────────────────────────────────────
watch(
  () => [props.modelValue, props.pet] as const,
  ([open, pet]) => {
    if (!open) return
    if (props.mode === 'edit' && pet) {
      dateMode.value = pet.birth_date_exact ? 'exact' : 'age'
      weightInput.value = ''
      weightUnit.value = 'kg'
      ageYears.value = 0
      ageMonths.value = 0
      resetForm({
        values: {
          name: pet.name,
          species: pet.species,
          breed: pet.breed ?? '',
          birth_date: pet.birth_date.slice(0, 10),
          birth_date_exact: pet.birth_date_exact,
          weight_grams: undefined,
          size: (pet.size as PetSize | undefined) ?? undefined,
        },
      })
    } else {
      dateMode.value = 'age'
      weightInput.value = ''
      weightUnit.value = 'kg'
      ageYears.value = 0
      ageMonths.value = 0
      resetForm({
        values: {
          name: '',
          species: '',
          breed: '',
          birth_date: estimatedBirthDate(0, 0),
          birth_date_exact: false,
          weight_grams: undefined,
          size: undefined,
        },
      })
    }
  },
  { immediate: true },
)

// ── Submit ───────────────────────────────────────────────────────────────────
const handleSave = handleSubmit(async (formValues) => {
  try {
    if (props.mode === 'create') {
      const payload: CreatePetFormValues = { ...formValues }
      if (weightGrams.value !== null) {
        payload.weight_grams = weightGrams.value
      }
      await createPet.mutateAsync(payload)
    } else if (props.pet) {
      // For edit, only send the base fields (no weight_grams / life_stage / species)
      // Species cannot be changed after creation
      await updatePet.mutateAsync({
        id: props.pet.id,
        payload: {
          name: formValues.name,
          breed: formValues.breed,
          birth_date: formValues.birth_date,
          birth_date_exact: formValues.birth_date_exact,
          size: formValues.size,
        },
      })
    }
    close()
    emit('success')
  } catch (e) {
    setErrors({ name: e instanceof Error ? e.message : 'Error al guardar' })
  }
})

function close() {
  emit('update:modelValue', false)
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div v-if="modelValue" class="modal-backdrop" @click.self="close">
        <div class="modal" role="dialog" :aria-label="mode === 'create' ? 'Nueva mascota' : 'Editar mascota'">
          <!-- Header -->
          <div class="modal-header">
            <div class="modal-header__text">
              <h2 class="modal-title">
                {{ mode === 'create' ? 'Nueva mascota' : 'Editar mascota' }}
              </h2>
              <p class="modal-subtitle">
                {{ mode === 'create' ? 'Registrá los datos de tu nueva mascota' : 'Actualizá los datos de ' + (pet?.name ?? 'la mascota') }}
              </p>
            </div>
            <button class="modal-close" @click="close">
              <IconX :size="18" :stroke-width="2.5" />
            </button>
          </div>

          <!-- Form -->
          <form class="modal-form" @submit.prevent="handleSave">
            <!-- Row: Nombre + Especie -->
            <div class="form-row">
              <div class="field">
                <label class="field-label">Nombre <span class="required">*</span></label>
                <input
                  v-model="name"
                  v-bind="nameAttrs"
                  class="field-input"
                  :class="{ 'field-input--error': errors.name }"
                  placeholder="Ej. Luna"
                  autocomplete="off"
                />
                <p v-if="errors.name" class="field-error">{{ errors.name }}</p>
              </div>

              <div class="field">
                <label class="field-label">Especie <span class="required">*</span></label>
                  <select
                  v-model="species"
                  v-bind="speciesAttrs"
                  class="field-input field-select"
                  :class="{ 'field-input--error': errors.species }"
                  :disabled="mode === 'edit'"
                >
                  <option value="" disabled>Seleccionar…</option>
                  <option v-for="s in PET_SPECIES" :key="s.value" :value="s.value">
                    {{ s.icon }} {{ s.label }}
                  </option>
                </select>
                <p v-if="errors.species" class="field-error">{{ errors.species }}</p>
              </div>
            </div>

            <!-- Raza -->
            <div class="field">
              <label class="field-label">Raza</label>
              <input
                v-model="breed"
                v-bind="breedAttrs"
                class="field-input"
                placeholder="Ej. Labrador"
                autocomplete="off"
              />
            </div>

            <!-- Tamaño (solo para perros) -->
            <Transition name="size-fade">
              <div v-if="values.species === 'dog'" class="field">
                <label class="field-label">
                  Tamaño <span class="required">*</span>
                </label>
                <div class="size-grid">
                  <button
                    v-for="sizeVal in PET_SIZE_VALUES"
                    :key="sizeVal"
                    type="button"
                    class="size-card"
                    :class="{ 'size-card--active': values.size === sizeVal }"
                    @click="setFieldValue('size', sizeVal)"
                  >
                    <span class="size-card__icon">{{ PET_SIZE_ICONS[sizeVal] }}</span>
                    <span class="size-card__label">{{ PET_SIZE_LABELS[sizeVal] }}</span>
                    <span class="size-card__desc">{{ PET_SIZE_DESCRIPTIONS[sizeVal] }}</span>
                  </button>
                </div>
                <p v-if="errors.size" class="field-error">{{ errors.size }}</p>
              </div>
            </Transition>

            <!-- Fecha de nacimiento -->
            <div class="field">
              <label class="field-label">
                Fecha de nacimiento
                <span v-if="mode === 'edit' && !values.birth_date_exact" class="field-hint">(estimada)</span>
                <span class="required">*</span>
              </label>

              <!-- Toggle: Exacta / Por edad (solo en create) -->
              <div v-if="mode === 'create'" class="date-toggle">
                <button
                  type="button"
                  class="toggle-btn"
                  :class="{ 'toggle-btn--active': dateMode === 'age' }"
                  @click="dateMode = 'age'"
                >
                  Por edad estimada
                </button>
                <button
                  type="button"
                  class="toggle-btn"
                  :class="{ 'toggle-btn--active': dateMode === 'exact' }"
                  @click="dateMode = 'exact'"
                >
                  Fecha exacta
                </button>
              </div>

              <!-- Date picker: show in edit mode or create with exact date -->
              <DatePicker
                v-if="mode === 'edit' || (mode === 'create' && dateMode === 'exact')"
                v-model="birthDate"
                :max-date="new Date()"
                placeholder="Seleccionar fecha"
                :error="!!errors.birth_date"
              />

              <!-- Age inputs - only in create mode with age mode -->
              <template v-if="mode === 'create' && dateMode === 'age'">
                <div class="age-inputs">
                  <div class="age-field">
                    <input
                      v-model.number="ageYears"
                      class="field-input age-input"
                      type="number"
                      min="0"
                      max="30"
                      placeholder="0"
                    />
                    <span class="age-unit">años</span>
                  </div>
                  <div class="age-field">
                    <input
                      v-model.number="ageMonths"
                      class="field-input age-input"
                      type="number"
                      min="0"
                      max="11"
                      placeholder="0"
                    />
                    <span class="age-unit">meses</span>
                  </div>
                </div>
              </template>

              <p v-if="errors.birth_date" class="field-error">{{ errors.birth_date }}</p>
            </div>

            <!-- Peso (solo create) -->
            <template v-if="mode === 'create'">
              <div class="field">
                <label class="field-label">Peso <span class="field-hint">(opcional)</span></label>
                <div class="weight-row">
                   <input
                     v-model.number="weightInput"
                     class="field-input weight-input"
                     type="number"
                     min="0"
                     step="any"
                     placeholder="0"
                   />
                  <div class="unit-toggle">
                    <button
                      type="button"
                      class="unit-btn"
                      :class="{ 'unit-btn--active': weightUnit === 'kg' }"
                      @click="weightUnit = 'kg'"
                    >kg</button>
                    <button
                      type="button"
                      class="unit-btn"
                      :class="{ 'unit-btn--active': weightUnit === 'g' }"
                      @click="weightUnit = 'g'"
                    >g</button>
                  </div>
                </div>
              </div>
            </template>

            <!-- Actions -->
            <div class="modal-actions">
              <button type="button" class="btn-secondary" @click="close">Cancelar</button>
              <button type="submit" class="btn-primary" :disabled="isPending()">
                <span v-if="isPending()" class="btn-spinner" />
                <span v-else>{{ mode === 'create' ? 'Crear mascota' : 'Guardar cambios' }}</span>
              </button>
            </div>
          </form>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>

<style scoped>
/* ── Backdrop ─────────────────────────── */
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

@media (max-width: 480px) {
  .modal-backdrop {
    align-items: flex-end;
    padding: 0;
  }
}

/* ── Modal box ────────────────────────── */
.modal {
  background: var(--color-surface);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-xl);
  width: 100%;
  max-width: 520px;
  overflow: hidden;
}

@media (max-width: 480px) {
  .modal {
    border-radius: var(--radius-xl) var(--radius-xl) 0 0;
    max-width: 100%;
    max-height: 92vh;
    overflow-y: auto;
  }
}

/* ── Header ───────────────────────────── */
.modal-header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: var(--space-4);
  padding: var(--space-6) var(--space-6) var(--space-5);
  border-bottom: 1px solid var(--color-border-light);
  background: var(--color-bg);
}

.modal-header__text {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.modal-title {
  font-family: var(--font-display);
  font-size: var(--text-xl);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.modal-subtitle {
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
  margin: 0;
}

.modal-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  background: var(--color-surface);
  color: var(--color-text-tertiary);
  cursor: pointer;
  flex-shrink: 0;
  transition: background var(--transition-fast), color var(--transition-fast), border-color var(--transition-fast);
}
.modal-close:hover {
  background: var(--color-bg-alt);
  color: var(--color-text-primary);
  border-color: var(--color-border);
}

/* ── Form ─────────────────────────────── */
.modal-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  padding: var(--space-6);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-4);
}

@media (max-width: 480px) {
  .form-row {
    grid-template-columns: 1fr;
  }
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

.field-hint {
  font-weight: 400;
  color: var(--color-text-tertiary);
  font-size: var(--text-xs);
}

.required {
  color: var(--color-error);
  margin-left: 1px;
}

.field-input {
  padding: var(--space-2) var(--space-3);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  background: var(--color-surface);
  transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
  width: 100%;
}
.field-input:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(61, 122, 95, 0.12);
}
.field-input--error {
  border-color: var(--color-error);
}
.field-input--error:focus {
  box-shadow: 0 0 0 3px rgba(220, 38, 38, 0.12);
}

.field-select {
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 24 24' fill='none' stroke='%239E9D99' stroke-width='2.5' stroke-linecap='round' stroke-linejoin='round'%3E%3Cpolyline points='6 9 12 15 18 9'%3E%3C/polyline%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right 12px center;
  padding-right: 2.25rem;
  cursor: pointer;
}

.field-error {
  font-size: var(--text-xs);
  color: var(--color-error);
  margin: 0;
}

/* ── Date toggle ──────────────────────── */
.date-toggle {
  display: flex;
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  overflow: hidden;
}

.toggle-btn {
  flex: 1;
  padding: var(--space-2) var(--space-3);
  background: transparent;
  border: none;
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-text-tertiary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}
.toggle-btn + .toggle-btn {
  border-left: 1.5px solid var(--color-border-light);
}
.toggle-btn--active {
  background: var(--color-accent);
  color: var(--color-text-on-accent);
  font-weight: 600;
}

/* ── Age inputs ───────────────────────── */
.age-inputs {
  display: flex;
  gap: var(--space-3);
}

.age-field {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  flex: 1;
}

.age-input {
  width: 70px;
  flex-shrink: 0;
  text-align: center;
}

.age-unit {
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
  white-space: nowrap;
}

/* ── Weight ───────────────────────────── */
.weight-row {
  display: flex;
  gap: var(--space-2);
  align-items: center;
}

.weight-input {
  flex: 1;
}

.unit-toggle {
  display: flex;
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  overflow: hidden;
  flex-shrink: 0;
}

.unit-btn {
  padding: var(--space-2) var(--space-3);
  background: transparent;
  border: none;
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-text-tertiary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}
.unit-btn + .unit-btn {
  border-left: 1.5px solid var(--color-border-light);
}
.unit-btn--active {
  background: var(--color-accent);
  color: var(--color-text-on-accent);
  font-weight: 600;
}

/* ── Actions ──────────────────────────── */
.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  padding-top: var(--space-2);
  border-top: 1px solid var(--color-border-light);
  margin-top: var(--space-2);
}

.btn-primary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-6);
  background: var(--color-accent);
  color: var(--color-text-on-accent);
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  min-width: 130px;
  transition: background var(--transition-fast), transform var(--transition-fast);
}
.btn-primary:hover:not(:disabled) {
  background: var(--color-accent-hover);
  transform: translateY(-1px);
}
.btn-primary:disabled {
  opacity: 0.65;
  cursor: not-allowed;
  transform: none;
}

.btn-secondary {
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
.btn-secondary:hover {
  background: var(--color-bg-alt);
  border-color: var(--color-border);
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

/* ── Size selector ────────────────────── */
.size-grid {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: var(--space-2);
}

@media (max-width: 480px) {
  .size-grid {
    grid-template-columns: repeat(2, 1fr);
  }
}

.size-card {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: var(--space-1);
  padding: var(--space-3) var(--space-2);
  background: var(--color-surface);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition:
    border-color var(--transition-fast),
    background var(--transition-fast),
    box-shadow var(--transition-fast);
  text-align: center;
}
.size-card:hover {
  border-color: var(--color-accent);
  background: var(--color-bg-alt);
}
.size-card--active {
  border-color: var(--color-accent);
  background: rgba(61, 122, 95, 0.07);
  box-shadow: 0 0 0 3px rgba(61, 122, 95, 0.15);
}

.size-card__icon {
  font-size: 1.5rem;
  line-height: 1;
}

.size-card__label {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
}

.size-card__desc {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  line-height: 1.3;
}

/* ── Size transition ──────────────────── */
.size-fade-enter-active,
.size-fade-leave-active {
  transition: opacity var(--transition-fast), transform var(--transition-fast);
}
.size-fade-enter-from,
.size-fade-leave-to {
  opacity: 0;
  transform: translateY(-6px);
}

/* ── Transition ───────────────────────── */
.modal-enter-active,
.modal-leave-active {
  transition: opacity var(--transition-fast);
}
.modal-enter-active .modal,
.modal-leave-active .modal {
  transition: transform var(--transition-fast), opacity var(--transition-fast);
}
.modal-enter-from,
.modal-leave-to {
  opacity: 0;
}
.modal-enter-from .modal,
.modal-leave-to .modal {
  transform: translateY(16px) scale(0.97);
  opacity: 0;
}
</style>
