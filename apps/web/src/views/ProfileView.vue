<script setup lang="ts">
import { ref, computed } from 'vue'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import {
  IconUser,
  IconMail,
  IconLock,
  IconCheck,
  IconAlertCircle,
  IconShieldCheck,
  IconBrandGoogle,
} from '@tabler/icons-vue'
import { useAuthStore } from '@/stores/auth'
import { authService } from '@/services/authService'
import { updateProfileSchema, changePasswordSchema } from '@/schemas/user'

const authStore = useAuthStore()

// ── Profile form ──────────────────────────────────────────────
const profileSubmitError = ref<string | null>(null)
const profileSuccess = ref(false)

const {
  defineField: defineProfileField,
  handleSubmit: handleProfileSubmit,
  errors: profileErrors,
} = useForm({
  validationSchema: toTypedSchema(updateProfileSchema),
  initialValues: {
    name: authStore.user?.name ?? '',
    email: authStore.user?.email ?? '',
  },
})

const [name, nameAttrs] = defineProfileField('name')
const [email, emailAttrs] = defineProfileField('email')

const handleUpdateProfile = handleProfileSubmit(async (values) => {
  profileSubmitError.value = null
  profileSuccess.value = false
  try {
    await authStore.updateProfile({ name: values.name, email: values.email })
    profileSuccess.value = true
    setTimeout(() => { profileSuccess.value = false }, 3000)
  } catch (e) {
    profileSubmitError.value = e instanceof Error ? e.message : 'Error al actualizar el perfil'
  }
})

// ── Password form ─────────────────────────────────────────────
const isLocalUser = computed(() => authStore.user?.auth_provider === 'local')

const passwordSubmitError = ref<string | null>(null)
const passwordSuccess = ref(false)

const {
  defineField: definePasswordField,
  handleSubmit: handlePasswordSubmit,
  errors: passwordErrors,
  resetForm: resetPasswordForm,
} = useForm({
  validationSchema: toTypedSchema(changePasswordSchema),
  initialValues: {
    current_password: '',
    new_password: '',
    confirm_password: '',
  },
})

const [currentPassword, currentPasswordAttrs] = definePasswordField('current_password')
const [newPassword, newPasswordAttrs] = definePasswordField('new_password')
const [confirmPassword, confirmPasswordAttrs] = definePasswordField('confirm_password')

const handleChangePassword = handlePasswordSubmit(async (values) => {
  passwordSubmitError.value = null
  passwordSuccess.value = false
  try {
    await authService.changePassword({
      current_password: values.current_password,
      new_password: values.new_password,
      confirm_password: values.confirm_password,
    })
    passwordSuccess.value = true
    resetPasswordForm()
    setTimeout(() => { passwordSuccess.value = false }, 4000)
  } catch (e) {
    passwordSubmitError.value = e instanceof Error ? e.message : 'Error al cambiar la contraseña'
  }
})

// ── User metadata ─────────────────────────────────────────────
const userInitials = computed(() => {
  const n = authStore.user?.name ?? ''
  return n.split(' ').slice(0, 2).map((w) => w[0]).join('').toUpperCase()
})

const joinedDate = computed(() => {
  const d = authStore.user?.created_at
  if (!d) return '—'
  return new Date(d).toLocaleDateString('es-ES', { year: 'numeric', month: 'long', day: 'numeric' })
})
</script>

<template>
  <div class="profile-page">
    <!-- Header -->
    <div class="page-header">
      <div class="header-avatar">{{ userInitials }}</div>
      <div class="header-info">
        <h1 class="header-name">{{ authStore.user?.name }}</h1>
        <div class="header-meta">
          <span
            class="meta-badge"
            :class="authStore.user?.is_system_user ? 'meta-badge--admin' : 'meta-badge--user'"
          >
            {{ authStore.user?.is_system_user ? 'Administrador' : 'Usuario' }}
          </span>
          <span class="meta-dot" />
          <span class="meta-text">Miembro desde {{ joinedDate }}</span>
        </div>
      </div>
    </div>

    <!-- Cards grid -->
    <div class="cards-grid">

      <!-- Datos básicos -->
      <section class="card">
        <div class="card-header">
          <div class="card-icon">
            <IconUser :size="18" :stroke-width="1.75" />
          </div>
          <div>
            <h2 class="card-title">Datos personales</h2>
            <p class="card-subtitle">Actualiza tu nombre y dirección de correo</p>
          </div>
        </div>

        <form class="form" @submit.prevent="handleUpdateProfile">
          <div class="form-field">
            <label class="form-label" for="profile-name">Nombre completo</label>
            <div class="input-wrapper" :class="{ 'input-wrapper--error': profileErrors.name }">
              <IconUser :size="16" :stroke-width="1.75" class="input-icon" />
              <input
                id="profile-name"
                v-model="name"
                v-bind="nameAttrs"
                type="text"
                class="form-input"
                placeholder="Tu nombre"
                autocomplete="name"
              />
            </div>
            <p v-if="profileErrors.name" class="field-error">{{ profileErrors.name }}</p>
          </div>

          <div class="form-field">
            <label class="form-label" for="profile-email">Correo electrónico</label>
            <div class="input-wrapper" :class="{ 'input-wrapper--error': profileErrors.email }">
              <IconMail :size="16" :stroke-width="1.75" class="input-icon" />
              <input
                id="profile-email"
                v-model="email"
                v-bind="emailAttrs"
                type="email"
                class="form-input"
                placeholder="tu@correo.com"
                autocomplete="email"
              />
            </div>
            <p v-if="profileErrors.email" class="field-error">{{ profileErrors.email }}</p>
          </div>

          <Transition name="alert">
            <div v-if="profileSubmitError" class="alert alert--error">
              <IconAlertCircle :size="16" :stroke-width="1.75" />
              <span>{{ profileSubmitError }}</span>
            </div>
          </Transition>

          <Transition name="alert">
            <div v-if="profileSuccess" class="alert alert--success">
              <IconCheck :size="16" :stroke-width="2" />
              <span>Perfil actualizado correctamente</span>
            </div>
          </Transition>

          <div class="form-footer">
            <button type="submit" class="btn btn--primary" :disabled="authStore.loading">
              <span v-if="authStore.loading" class="btn-spinner" />
              <span>{{ authStore.loading ? 'Guardando...' : 'Guardar cambios' }}</span>
            </button>
          </div>
        </form>
      </section>

      <!-- Contraseña -->
      <section class="card">
        <div class="card-header">
          <div class="card-icon">
            <IconLock :size="18" :stroke-width="1.75" />
          </div>
          <div>
            <h2 class="card-title">Seguridad</h2>
            <p class="card-subtitle">Gestiona tu contraseña de acceso</p>
          </div>
        </div>

        <!-- Google user notice -->
        <div v-if="!isLocalUser" class="provider-notice">
          <IconBrandGoogle :size="20" :stroke-width="1.75" class="provider-icon" />
          <div>
            <p class="provider-title">Cuenta de Google</p>
            <p class="provider-desc">
              Tu cuenta está vinculada a Google. El cambio de contraseña no está disponible para este tipo de cuenta.
            </p>
          </div>
        </div>

        <!-- Password form (local users only) -->
        <form v-else class="form" @submit.prevent="handleChangePassword">
          <div class="form-field">
            <label class="form-label" for="current-password">Contraseña actual</label>
            <div class="input-wrapper" :class="{ 'input-wrapper--error': passwordErrors.current_password }">
              <IconLock :size="16" :stroke-width="1.75" class="input-icon" />
              <input
                id="current-password"
                v-model="currentPassword"
                v-bind="currentPasswordAttrs"
                type="password"
                class="form-input"
                placeholder="••••••••"
                autocomplete="current-password"
              />
            </div>
            <p v-if="passwordErrors.current_password" class="field-error">{{ passwordErrors.current_password }}</p>
          </div>

          <div class="form-field">
            <label class="form-label" for="new-password">Nueva contraseña</label>
            <div class="input-wrapper" :class="{ 'input-wrapper--error': passwordErrors.new_password }">
              <IconShieldCheck :size="16" :stroke-width="1.75" class="input-icon" />
              <input
                id="new-password"
                v-model="newPassword"
                v-bind="newPasswordAttrs"
                type="password"
                class="form-input"
                placeholder="Mínimo 8 caracteres"
                autocomplete="new-password"
              />
            </div>
            <p v-if="passwordErrors.new_password" class="field-error">{{ passwordErrors.new_password }}</p>
          </div>

          <div class="form-field">
            <label class="form-label" for="confirm-password">Confirmar contraseña</label>
            <div class="input-wrapper" :class="{ 'input-wrapper--error': passwordErrors.confirm_password }">
              <IconShieldCheck :size="16" :stroke-width="1.75" class="input-icon" />
              <input
                id="confirm-password"
                v-model="confirmPassword"
                v-bind="confirmPasswordAttrs"
                type="password"
                class="form-input"
                placeholder="Repite la nueva contraseña"
                autocomplete="new-password"
              />
            </div>
            <p v-if="passwordErrors.confirm_password" class="field-error">{{ passwordErrors.confirm_password }}</p>
          </div>

          <Transition name="alert">
            <div v-if="passwordSubmitError" class="alert alert--error">
              <IconAlertCircle :size="16" :stroke-width="1.75" />
              <span>{{ passwordSubmitError }}</span>
            </div>
          </Transition>

          <Transition name="alert">
            <div v-if="passwordSuccess" class="alert alert--success">
              <IconCheck :size="16" :stroke-width="2" />
              <span>Contraseña actualizada correctamente</span>
            </div>
          </Transition>

          <div class="form-footer">
            <button type="submit" class="btn btn--primary" :disabled="authStore.loading">
              <span v-if="authStore.loading" class="btn-spinner" />
              <span>{{ authStore.loading ? 'Cambiando...' : 'Cambiar contraseña' }}</span>
            </button>
          </div>
        </form>
      </section>

    </div>
  </div>
</template>

<style scoped>
/* ── Page layout ─────────────────────────────────────────────── */
.profile-page {
  padding: var(--space-8) var(--space-8) var(--space-12);
  max-width: 900px;
}

/* ── Header ──────────────────────────────────────────────────── */
.page-header {
  display: flex;
  align-items: center;
  gap: var(--space-5);
  margin-bottom: var(--space-8);
}

.header-avatar {
  width: 64px;
  height: 64px;
  border-radius: var(--radius-full);
  background: var(--color-accent);
  color: var(--color-text-inverse);
  display: flex;
  align-items: center;
  justify-content: center;
  font-family: var(--font-display);
  font-size: var(--text-xl);
  font-weight: 600;
  letter-spacing: 0.02em;
  flex-shrink: 0;
}

.header-info {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.header-name {
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  font-weight: 600;
  color: var(--color-text-primary);
  line-height: 1.2;
}

.header-meta {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.meta-badge {
  font-size: var(--text-xs);
  font-weight: 600;
  padding: 2px var(--space-2);
  border-radius: var(--radius-full);
  letter-spacing: 0.03em;
}

.meta-badge--admin {
  background: var(--color-accent-light);
  color: var(--color-accent-dark);
}

.meta-badge--user {
  background: var(--color-bg-alt);
  color: var(--color-text-secondary);
}

.meta-dot {
  width: 3px;
  height: 3px;
  border-radius: var(--radius-full);
  background: var(--color-text-tertiary);
}

.meta-text {
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
}

/* ── Cards grid ──────────────────────────────────────────────── */
.cards-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(380px, 1fr));
  gap: var(--space-6);
  align-items: start;
}

/* ── Card ────────────────────────────────────────────────────── */
.card {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  padding: var(--space-6);
  box-shadow: var(--shadow-sm);
}

.card-header {
  display: flex;
  align-items: flex-start;
  gap: var(--space-4);
  margin-bottom: var(--space-6);
  padding-bottom: var(--space-5);
  border-bottom: 1px solid var(--color-border-light);
}

.card-icon {
  width: 38px;
  height: 38px;
  border-radius: var(--radius-md);
  background: var(--color-accent-light);
  color: var(--color-accent);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.card-title {
  font-family: var(--font-display);
  font-size: var(--text-lg);
  font-weight: 600;
  color: var(--color-text-primary);
  line-height: 1.3;
}

.card-subtitle {
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
  margin-top: 2px;
}

/* ── Form ────────────────────────────────────────────────────── */
.form {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

.form-field {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.form-label {
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-text-secondary);
}

.input-wrapper {
  position: relative;
  display: flex;
  align-items: center;
}

.input-icon {
  position: absolute;
  left: var(--space-3);
  color: var(--color-text-tertiary);
  pointer-events: none;
  flex-shrink: 0;
}

.form-input {
  width: 100%;
  padding: var(--space-3) var(--space-3) var(--space-3) calc(var(--space-3) + 16px + var(--space-2));
  font-family: var(--font-body);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  background: var(--color-bg);
  border: 1.5px solid var(--color-border);
  border-radius: var(--radius-md);
  outline: none;
  transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
}

.form-input::placeholder {
  color: var(--color-text-tertiary);
}

.form-input:focus {
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px var(--color-accent-light);
}

.input-wrapper--error .form-input {
  border-color: var(--color-error);
}

.input-wrapper--error .form-input:focus {
  box-shadow: 0 0 0 3px var(--color-error-light);
}

/* ── Field error ─────────────────────────────────────────────── */
.field-error {
  font-size: var(--text-xs);
  color: var(--color-error);
  margin-top: 2px;
}

/* ── Alerts ──────────────────────────────────────────────────── */
.alert {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-4);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 500;
}

.alert--error {
  background: var(--color-error-light);
  color: var(--color-error);
  border: 1px solid var(--color-error-border);
}

.alert--success {
  background: var(--color-accent-light);
  color: var(--color-accent-dark);
  border: 1px solid rgba(61, 122, 95, 0.2);
}

/* ── Form footer ─────────────────────────────────────────────── */
.form-footer {
  display: flex;
  justify-content: flex-end;
  padding-top: var(--space-2);
}

/* ── Button ──────────────────────────────────────────────────── */
.btn {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-5);
  border: none;
  border-radius: var(--radius-md);
  font-family: var(--font-body);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast), opacity var(--transition-fast), transform var(--transition-fast);
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn--primary {
  background: var(--color-accent);
  color: var(--color-text-inverse);
}

.btn--primary:not(:disabled):hover {
  background: var(--color-accent-hover);
  transform: translateY(-1px);
}

.btn--primary:not(:disabled):active {
  transform: translateY(0);
}

/* ── Spinner ─────────────────────────────────────────────────── */
.btn-spinner {
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255, 255, 255, 0.35);
  border-top-color: #fff;
  border-radius: var(--radius-full);
  animation: spin 0.7s linear infinite;
  flex-shrink: 0;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ── Google provider notice ──────────────────────────────────── */
.provider-notice {
  display: flex;
  gap: var(--space-4);
  padding: var(--space-4);
  background: var(--color-bg-alt);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-md);
}

.provider-icon {
  color: var(--color-text-tertiary);
  flex-shrink: 0;
  margin-top: 2px;
}

.provider-title {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-bottom: var(--space-1);
}

.provider-desc {
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
  line-height: 1.5;
}

/* ── Alert transition ────────────────────────────────────────── */
.alert-enter-active,
.alert-leave-active {
  transition: opacity var(--transition-base), transform var(--transition-base);
}

.alert-enter-from,
.alert-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}

/* ── Responsive ──────────────────────────────────────────────── */
@media (max-width: 768px) {
  .profile-page {
    padding: var(--space-6) var(--space-4) var(--space-10);
  }

  .cards-grid {
    grid-template-columns: 1fr;
  }

  .page-header {
    gap: var(--space-4);
  }

  .header-avatar {
    width: 52px;
    height: 52px;
    font-size: var(--text-lg);
  }

  .header-name {
    font-size: var(--text-xl);
  }
}
</style>
