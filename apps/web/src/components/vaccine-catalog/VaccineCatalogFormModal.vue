<script setup lang="ts">
import { ref, computed } from 'vue'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { IconX } from '@tabler/icons-vue'
import { useCreateVaccineCatalog, useUpdateVaccineCatalog } from '@/composables/useVaccineCatalog'
import type { VaccineCatalog } from '@/types/vaccineCatalog'
import { createVaccineCatalogSchema, updateVaccineCatalogSchema } from '@/schemas/vaccineCatalog'

const props = defineProps<{
  mode: 'create' | 'edit'
  vaccine?: VaccineCatalog
}>()

const emit = defineEmits<{
  close: []
  saved: []
}>()

const submitError = ref<string | null>(null)
const createVaccine = useCreateVaccineCatalog()
const updateVaccine = useUpdateVaccineCatalog()

// Especies disponibles — podrían venir de constantes o backend en el futuro
const availableSpecies = ['perro', 'gato', 'conejo', 'roedor', 'ave', 'reptil']
const selectedSpecies = ref<string[]>(props.vaccine?.species ?? [])

// ── Create form ─────────────────────────────────────────────
const {
  defineField: defineCreateField,
  handleSubmit: handleCreateSubmit,
  errors: createErrors,
} = useForm({
  validationSchema: toTypedSchema(createVaccineCatalogSchema),
  initialValues: {
    name: '',
    species: [],
    frequency_months: 12,
    is_mandatory: false,
  },
})

const [createName, createNameAttrs] = defineCreateField('name')
const [createFrequencyMonths, createFrequencyMonthsAttrs] = defineCreateField('frequency_months')
const [createIsMandatory, createIsMandatoryAttrs] = defineCreateField('is_mandatory')

const handleCreate = handleCreateSubmit(async (values) => {
  submitError.value = null
  try {
    await createVaccine.mutateAsync({
      name: values.name,
      species: selectedSpecies.value,
      frequency_months: values.frequency_months,
      is_mandatory: values.is_mandatory ?? false,
    })
    emit('saved')
  } catch (e) {
    submitError.value = e instanceof Error ? e.message : 'Error al crear vacuna'
  }
})

// ── Edit form ───────────────────────────────────────────────
const {
  defineField: defineEditField,
  handleSubmit: handleEditSubmit,
  errors: editErrors,
} = useForm({
  validationSchema: toTypedSchema(updateVaccineCatalogSchema),
  initialValues: {
    name: props.vaccine?.name ?? '',
    species: props.vaccine?.species ?? [],
    frequency_months: props.vaccine?.frequency_months ?? 12,
    is_mandatory: props.vaccine?.is_mandatory ?? false,
  },
})

const [editName, editNameAttrs] = defineEditField('name')
const [editFrequencyMonths, editFrequencyMonthsAttrs] = defineEditField('frequency_months')
const [editIsMandatory, editIsMandatoryAttrs] = defineEditField('is_mandatory')

const handleEdit = handleEditSubmit(async (values) => {
  if (!props.vaccine) return
  submitError.value = null
  try {
    await updateVaccine.mutateAsync({
      id: props.vaccine.id,
      payload: {
        name: values.name,
        species: selectedSpecies.value,
        frequency_months: values.frequency_months,
        is_mandatory: values.is_mandatory ?? false,
      },
    })
    emit('saved')
  } catch (e) {
    submitError.value = e instanceof Error ? e.message : 'Error al actualizar vacuna'
  }
})

function toggleSpecies(specie: string) {
  const idx = selectedSpecies.value.indexOf(specie)
  if (idx > -1) {
    selectedSpecies.value.splice(idx, 1)
  } else {
    selectedSpecies.value.push(specie)
  }
}

const hasSpeciesError = computed(() => {
  const createError = createErrors.value['species']
  const editError = editErrors.value['species']
  return props.mode === 'create' ? createError : editError
})
</script>

<template>
  <Teleport to="body">
    <div class="modal-backdrop" @click.self="emit('close')">
      <div class="modal">
        <div class="modal-header">
          <h2 class="modal-title">
            {{ mode === 'create' ? 'Nueva vacuna' : 'Editar vacuna' }}
          </h2>
          <button class="modal-close" @click="emit('close')">
            <IconX :size="18" :stroke-width="2.5" />
          </button>
        </div>

        <!-- Create form -->
        <form v-if="mode === 'create'" class="modal-form" @submit.prevent="handleCreate">
          <div class="field">
            <label class="field-label">Nombre de la vacuna</label>
            <input
              v-model="createName"
              v-bind="createNameAttrs"
              class="field-input"
              :class="{ 'field-input--error': createErrors['name'] }"
              placeholder="Ej. Vacuna rabia"
            />
            <p v-if="createErrors['name']" class="field-error">{{ createErrors['name'] }}</p>
          </div>

          <div class="field">
            <label class="field-label">Especies aplicables</label>
            <div class="species-checkboxes">
              <label v-for="specie in availableSpecies" :key="specie" class="checkbox-label">
                <input
                  type="checkbox"
                  :checked="selectedSpecies.includes(specie)"
                  @change="toggleSpecies(specie)"
                  class="checkbox-input"
                />
                <span class="checkbox-text">{{ specie }}</span>
              </label>
            </div>
            <p v-if="hasSpeciesError" class="field-error">{{ hasSpeciesError }}</p>
          </div>

          <div class="field">
            <label class="field-label">Frecuencia (meses)</label>
            <input
              v-model.number="createFrequencyMonths"
              v-bind="createFrequencyMonthsAttrs"
              class="field-input"
              :class="{ 'field-input--error': createErrors['frequency_months'] }"
              type="number"
              min="1"
              max="360"
              placeholder="12"
            />
            <p v-if="createErrors['frequency_months']" class="field-error">
              {{ createErrors['frequency_months'] }}
            </p>
          </div>

          <div class="field field--checkbox">
            <label class="checkbox-label">
              <input
                v-model="createIsMandatory"
                v-bind="createIsMandatoryAttrs"
                type="checkbox"
                class="checkbox-input"
              />
              <span class="checkbox-text">Vacuna obligatoria</span>
            </label>
          </div>

          <p v-if="submitError" class="submit-error">{{ submitError }}</p>

          <div class="modal-actions">
            <button type="button" class="btn-secondary" @click="emit('close')">Cancelar</button>
            <button type="submit" class="btn-primary" :disabled="createVaccine.isPending.value">
              <div v-if="createVaccine.isPending.value" class="spinner spinner--sm spinner--white" />
              <span v-else>Crear vacuna</span>
            </button>
          </div>
        </form>

        <!-- Edit form -->
        <form v-else class="modal-form" @submit.prevent="handleEdit">
          <div class="field">
            <label class="field-label">Nombre de la vacuna</label>
            <input
              v-model="editName"
              v-bind="editNameAttrs"
              class="field-input"
              :class="{ 'field-input--error': editErrors['name'] }"
              placeholder="Ej. Vacuna rabia"
            />
            <p v-if="editErrors['name']" class="field-error">{{ editErrors['name'] }}</p>
          </div>

          <div class="field">
            <label class="field-label">Especies aplicables</label>
            <div class="species-checkboxes">
              <label v-for="specie in availableSpecies" :key="specie" class="checkbox-label">
                <input
                  type="checkbox"
                  :checked="selectedSpecies.includes(specie)"
                  @change="toggleSpecies(specie)"
                  class="checkbox-input"
                />
                <span class="checkbox-text">{{ specie }}</span>
              </label>
            </div>
            <p v-if="hasSpeciesError" class="field-error">{{ hasSpeciesError }}</p>
          </div>

          <div class="field">
            <label class="field-label">Frecuencia (meses)</label>
            <input
              v-model.number="editFrequencyMonths"
              v-bind="editFrequencyMonthsAttrs"
              class="field-input"
              :class="{ 'field-input--error': editErrors['frequency_months'] }"
              type="number"
              min="1"
              max="360"
              placeholder="12"
            />
            <p v-if="editErrors['frequency_months']" class="field-error">
              {{ editErrors['frequency_months'] }}
            </p>
          </div>

          <div class="field field--checkbox">
            <label class="checkbox-label">
              <input
                v-model="editIsMandatory"
                v-bind="editIsMandatoryAttrs"
                type="checkbox"
                class="checkbox-input"
              />
              <span class="checkbox-text">Vacuna obligatoria</span>
            </label>
          </div>

          <p v-if="submitError" class="submit-error">{{ submitError }}</p>

          <div class="modal-actions">
            <button type="button" class="btn-secondary" @click="emit('close')">Cancelar</button>
            <button type="submit" class="btn-primary" :disabled="updateVaccine.isPending.value">
              <div v-if="updateVaccine.isPending.value" class="spinner spinner--sm spinner--white" />
              <span v-else>Guardar cambios</span>
            </button>
          </div>
        </form>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(0, 0, 0, 0.45);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: var(--space-4);
  animation: fade-in 0.15s ease;
}

@keyframes fade-in {
  from {
    opacity: 0;
  }
  to {
    opacity: 1;
  }
}

.modal {
  background: #fff;
  border-radius: var(--radius-lg);
  box-shadow: 0 20px 60px rgba(0, 0, 0, 0.15), 0 4px 16px rgba(0, 0, 0, 0.08);
  width: 100%;
  max-width: 500px;
  animation: slide-up 0.2s ease;
}

@keyframes slide-up {
  from {
    opacity: 0;
    transform: translateY(12px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 480px) {
  .modal-backdrop {
    align-items: flex-end;
    padding: 0;
  }
  .modal {
    border-radius: var(--radius-lg) var(--radius-lg) 0 0;
    max-width: 100%;
    animation: slide-up-mobile 0.25s ease;
  }
  @keyframes slide-up-mobile {
    from {
      transform: translateY(100%);
    }
    to {
      transform: translateY(0);
    }
  }
}

.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-5) var(--space-6) var(--space-4);
  border-bottom: 1px solid var(--color-border-light);
}

.modal-title {
  font-size: var(--text-lg);
  font-weight: 700;
  color: var(--color-text-primary);
  margin: 0;
}

.modal-close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 30px;
  height: 30px;
  border: none;
  background: transparent;
  border-radius: var(--radius-md);
  color: var(--color-text-tertiary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.modal-close:hover {
  background: #f3f4f6;
  color: var(--color-text-primary);
}

.modal-form {
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
  padding: var(--space-5) var(--space-6) var(--space-6);
}

.field {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.field--checkbox {
  gap: 0;
}

.field-label {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-secondary);
}

.field-input {
  padding: var(--space-2) var(--space-3);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  background: #fff;
  transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
}

.field-input:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(66, 184, 131, 0.15);
}

.field-input--error {
  border-color: var(--color-error);
}

.field-input--error:focus {
  box-shadow: 0 0 0 3px rgba(220, 38, 38, 0.12);
}

.field-error {
  font-size: var(--text-xs);
  color: var(--color-error);
  margin: 0;
}

.species-checkboxes {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(120px, 1fr));
  gap: var(--space-2);
  margin-top: var(--space-2);
}

.checkbox-label {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  cursor: pointer;
  user-select: none;
}

.checkbox-input {
  width: 18px;
  height: 18px;
  cursor: pointer;
  accent-color: var(--color-accent);
}

.checkbox-text {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.submit-error {
  font-size: var(--text-sm);
  color: var(--color-error);
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: var(--radius-md);
  padding: var(--space-2) var(--space-3);
  margin: 0;
}

.modal-actions {
  display: flex;
  justify-content: flex-end;
  gap: var(--space-3);
  padding-top: var(--space-2);
}

.btn-primary {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-5);
  background: var(--color-accent);
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast);
  min-width: 140px;
}

.btn-primary:hover:not(:disabled) {
  background: #369a6e;
}

.btn-primary:disabled {
  opacity: 0.7;
  cursor: not-allowed;
}

.btn-secondary {
  padding: var(--space-2) var(--space-4);
  background: #fff;
  color: var(--color-text-secondary);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 500;
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast);
}

.btn-secondary:hover {
  background: var(--color-bg);
  border-color: #d1d5db;
}

.spinner {
  width: 28px;
  height: 28px;
  border: 3px solid var(--color-border-light);
  border-top-color: var(--color-accent);
  border-radius: var(--radius-full);
  animation: spin 0.7s linear infinite;
}

.spinner--sm {
  width: 16px;
  height: 16px;
  border-width: 2px;
}

.spinner--white {
  border-color: rgba(255, 255, 255, 0.3);
  border-top-color: #fff;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}
</style>
