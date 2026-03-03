<script setup lang="ts">
import { watch } from 'vue'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { IconX } from '@tabler/icons-vue'
import { useCreatePet, useUpdatePet } from '@/composables/usePets'
import { petSchema } from '@/schemas/pet'
import type { Pet } from '@/types/pet'

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

const { defineField, handleSubmit, errors, resetForm, setErrors } = useForm({
  validationSchema: toTypedSchema(petSchema),
  initialValues: { name: '', species: '', breed: '', age: undefined, owner: '' },
})

const [name, nameAttrs] = defineField('name')
const [species, speciesAttrs] = defineField('species')
const [breed, breedAttrs] = defineField('breed')
const [age, ageAttrs] = defineField('age')
const [owner, ownerAttrs] = defineField('owner')

// Pre-fill form when editing
watch(
  () => [props.modelValue, props.pet] as const,
  ([open, pet]) => {
    if (!open) return
    if (props.mode === 'edit' && pet) {
      resetForm({
        values: {
          name: pet.name,
          species: pet.species,
          breed: pet.breed ?? '',
          age: pet.age,
          owner: pet.owner ?? '',
        },
      })
    } else {
      resetForm({ values: { name: '', species: '', breed: '', age: undefined, owner: '' } })
    }
  },
  { immediate: true },
)

const handleSave = handleSubmit(async (values) => {
  try {
    if (props.mode === 'create') {
      await createPet.mutateAsync(values)
    } else if (props.pet) {
      await updatePet.mutateAsync({ id: props.pet.id, payload: values })
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
                >
                  <option value="" disabled>Seleccionar…</option>
                  <option value="dog">🐕 Perro</option>
                  <option value="cat">🐈 Gato</option>
                  <option value="bird">🦜 Ave</option>
                  <option value="rabbit">🐇 Conejo</option>
                  <option value="fish">🐠 Pez</option>
                  <option value="other">🐾 Otro</option>
                </select>
                <p v-if="errors.species" class="field-error">{{ errors.species }}</p>
              </div>
            </div>

            <!-- Row: Raza + Edad -->
            <div class="form-row">
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

              <div class="field">
                <label class="field-label">Edad (años)</label>
                <input
                  v-model.number="age"
                  v-bind="ageAttrs"
                  class="field-input"
                  :class="{ 'field-input--error': errors.age }"
                  type="number"
                  placeholder="0"
                  min="0"
                  max="100"
                />
                <p v-if="errors.age" class="field-error">{{ errors.age }}</p>
              </div>
            </div>

            <!-- Dueño -->
            <div class="field">
              <label class="field-label">Dueño</label>
              <input
                v-model="owner"
                v-bind="ownerAttrs"
                class="field-input"
                placeholder="Nombre del dueño"
                autocomplete="off"
              />
            </div>

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

/* ── Modal box ────────────────────────── */
.modal {
  background: var(--color-surface);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-xl);
  width: 100%;
  max-width: 520px;
  overflow: hidden;
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

.field {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
}

.field-label {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-secondary);
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
