<script setup lang="ts">
import { ref, computed } from 'vue'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { IconX } from '@tabler/icons-vue'
import { useCreateUser, useUpdateUser } from '@/composables/useUsers'
import { useAuthStore } from '@/stores/auth'
import type { User, CreateUserPayload } from '@/types/user'
import { createUserSchema, updateUserSchema } from '@/schemas/user'

const props = defineProps<{
  mode: 'create' | 'edit'
  user?: User
}>()

const emit = defineEmits<{
  close: []
  saved: []
}>()

const authStore = useAuthStore()
const isSystemUser = computed(() => authStore.user?.is_system_user ?? false)

const submitError = ref<string | null>(null)
const createUser = useCreateUser()
const updateUser = useUpdateUser()

// ── Create form ─────────────────────────────────────────────
const {
  defineField: defineCreateField,
  handleSubmit: handleCreateSubmit,
  errors: createErrors,
} = useForm({
  validationSchema: toTypedSchema(createUserSchema),
  initialValues: { name: '', email: '', password: '', pet_limit: 5 },
})

const [createName, createNameAttrs] = defineCreateField('name')
const [createEmail, createEmailAttrs] = defineCreateField('email')
const [createPassword, createPasswordAttrs] = defineCreateField('password')
const [createPetLimit, createPetLimitAttrs] = defineCreateField('pet_limit')

const handleCreate = handleCreateSubmit(async (values) => {
  submitError.value = null
  try {
    const payload: CreateUserPayload = {
      name: values.name,
      email: values.email,
      password: values.password,
    }
    if (isSystemUser.value && values.pet_limit !== undefined) {
      payload.pet_limit = values.pet_limit
    }
    await createUser.mutateAsync(payload)
    emit('saved')
  } catch (e) {
    submitError.value = e instanceof Error ? e.message : 'Error al crear usuario'
  }
})

// ── Edit form ───────────────────────────────────────────────
const {
  defineField: defineEditField,
  handleSubmit: handleEditSubmit,
  errors: editErrors,
} = useForm({
  validationSchema: toTypedSchema(updateUserSchema),
  initialValues: { 
    name: props.user?.name ?? '', 
    email: props.user?.email ?? '',
    pet_limit: props.user?.pet_limit ?? 5
  },
})

const [editName, editNameAttrs] = defineEditField('name')
const [editEmail, editEmailAttrs] = defineEditField('email')
const [editPetLimit, editPetLimitAttrs] = defineEditField('pet_limit')

const handleEdit = handleEditSubmit(async (values) => {
  if (!props.user) return
  submitError.value = null
  try {
    const payload = { ...values }
    if (!isSystemUser.value) {
      delete payload.pet_limit
    }
    await updateUser.mutateAsync({ id: props.user.id, payload })
    emit('saved')
  } catch (e) {
    submitError.value = e instanceof Error ? e.message : 'Error al actualizar usuario'
  }
})
</script>

<template>
  <Teleport to="body">
    <div class="modal-backdrop" @click.self="emit('close')">
      <div class="modal">
        <div class="modal-header">
          <h2 class="modal-title">
            {{ mode === 'create' ? 'Nuevo usuario' : 'Editar usuario' }}
          </h2>
          <button class="modal-close" @click="emit('close')">
            <IconX :size="18" :stroke-width="2.5" />
          </button>
        </div>

        <!-- Create form -->
        <form v-if="mode === 'create'" class="modal-form" @submit.prevent="handleCreate">
          <div class="field">
            <label class="field-label">Nombre completo</label>
            <input
              v-model="createName"
              v-bind="createNameAttrs"
              class="field-input"
              :class="{ 'field-input--error': createErrors['name'] }"
              placeholder="Ej. María García"
              autocomplete="name"
            />
            <p v-if="createErrors['name']" class="field-error">{{ createErrors['name'] }}</p>
          </div>

          <div class="field">
            <label class="field-label">Correo electrónico</label>
            <input
              v-model="createEmail"
              v-bind="createEmailAttrs"
              class="field-input"
              :class="{ 'field-input--error': createErrors['email'] }"
              type="email"
              placeholder="correo@ejemplo.com"
              autocomplete="email"
            />
            <p v-if="createErrors['email']" class="field-error">{{ createErrors['email'] }}</p>
          </div>

          <div class="field">
            <label class="field-label">Contraseña</label>
            <input
              v-model="createPassword"
              v-bind="createPasswordAttrs"
              class="field-input"
              :class="{ 'field-input--error': createErrors['password'] }"
              type="password"
              placeholder="Mínimo 8 caracteres"
              autocomplete="new-password"
            />
            <p v-if="createErrors['password']" class="field-error">
              {{ createErrors['password'] }}
            </p>
          </div>

          <div v-if="isSystemUser" class="field">
            <label class="field-label">Límite de mascotas</label>
            <input
              v-model="createPetLimit"
              v-bind="createPetLimitAttrs"
              class="field-input"
              :class="{ 'field-input--error': createErrors['pet_limit'] }"
              type="number"
              min="0"
              placeholder="5"
            />
            <p v-if="createErrors['pet_limit']" class="field-error">
              {{ createErrors['pet_limit'] }}
            </p>
          </div>

          <p v-if="submitError" class="submit-error">{{ submitError }}</p>

          <div class="modal-actions">
            <button type="button" class="btn-secondary" @click="emit('close')">Cancelar</button>
            <button type="submit" class="btn-primary" :disabled="createUser.isPending.value">
              <div v-if="createUser.isPending.value" class="spinner spinner--sm spinner--white" />
              <span v-else>Crear usuario</span>
            </button>
          </div>
        </form>

        <!-- Edit form -->
        <form v-else class="modal-form" @submit.prevent="handleEdit">
          <div class="field">
            <label class="field-label">Nombre completo</label>
            <input
              v-model="editName"
              v-bind="editNameAttrs"
              class="field-input"
              :class="{ 'field-input--error': editErrors['name'] }"
              placeholder="Ej. María García"
              autocomplete="name"
            />
            <p v-if="editErrors['name']" class="field-error">{{ editErrors['name'] }}</p>
          </div>

          <div class="field">
            <label class="field-label">Correo electrónico</label>
            <input
              v-model="editEmail"
              v-bind="editEmailAttrs"
              class="field-input"
              :class="{ 'field-input--error': editErrors['email'] }"
              type="email"
              placeholder="correo@ejemplo.com"
              autocomplete="email"
            />
            <p v-if="editErrors['email']" class="field-error">{{ editErrors['email'] }}</p>
          </div>

          <div v-if="isSystemUser" class="field">
            <label class="field-label">Límite de mascotas</label>
            <input
              v-model="editPetLimit"
              v-bind="editPetLimitAttrs"
              class="field-input"
              :class="{ 'field-input--error': editErrors['pet_limit'] }"
              type="number"
              min="0"
              placeholder="5"
            />
            <p v-if="editErrors['pet_limit']" class="field-error">
              {{ editErrors['pet_limit'] }}
            </p>
          </div>

          <p v-if="submitError" class="submit-error">{{ submitError }}</p>

          <div class="modal-actions">
            <button type="button" class="btn-secondary" @click="emit('close')">Cancelar</button>
            <button type="submit" class="btn-primary" :disabled="updateUser.isPending.value">
              <div v-if="updateUser.isPending.value" class="spinner spinner--sm spinner--white" />
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
  from { opacity: 0; }
  to { opacity: 1; }
}

.modal {
  background: #fff;
  border-radius: var(--radius-lg);
  box-shadow: 0 20px 60px rgba(0,0,0,0.15), 0 4px 16px rgba(0,0,0,0.08);
  width: 100%;
  max-width: 440px;
  animation: slide-up 0.2s ease;
}

@keyframes slide-up {
  from { opacity: 0; transform: translateY(12px); }
  to { opacity: 1; transform: translateY(0); }
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
    from { transform: translateY(100%); }
    to { transform: translateY(0); }
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

.field-input--error { border-color: var(--color-error); }
.field-input--error:focus { box-shadow: 0 0 0 3px rgba(220, 38, 38, 0.12); }

.field-error {
  font-size: var(--text-xs);
  color: var(--color-error);
  margin: 0;
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
  min-width: 120px;
}

.btn-primary:hover:not(:disabled) { background: #369a6e; }
.btn-primary:disabled { opacity: 0.7; cursor: not-allowed; }

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

.spinner--sm { width: 16px; height: 16px; border-width: 2px; }
.spinner--white { border-color: rgba(255,255,255,0.3); border-top-color: #fff; }

@keyframes spin { to { transform: rotate(360deg); } }
</style>
