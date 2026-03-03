<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { setupService } from '@/services/setupService'
import { resetSetupCache } from '@/router'
import { setupSchema } from '@/schemas/user'

const router = useRouter()

const loading = ref(false)
const error = ref<string | null>(null)

const { defineField, handleSubmit, errors } = useForm({
  validationSchema: toTypedSchema(setupSchema),
  initialValues: { name: '', email: '', password: '', confirmPassword: '' },
})

const [name, nameAttrs] = defineField('name')
const [email, emailAttrs] = defineField('email')
const [password, passwordAttrs] = defineField('password')
const [confirmPassword, confirmPasswordAttrs] = defineField('confirmPassword')

const handleSetup = handleSubmit(async (values) => {
  error.value = null
  loading.value = true
  try {
    await setupService.createFirstUser({
      name: values.name,
      email: values.email,
      password: values.password,
    })
    resetSetupCache()
    router.push({ name: 'login', query: { success: 'account_created' } })
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Error al crear el usuario'
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <div class="auth-page">
    <!-- Panel decorativo izquierdo -->
    <div class="auth-panel">
      <div class="panel-content">
        <div class="panel-brand">
          <div class="panel-icon">
            <svg width="28" height="28" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
              <path d="M11 6c0 4-3 6-3 10a4 4 0 0 0 8 0c0-4-3-6-3-10"/>
              <circle cx="11" cy="4" r="1.5"/>
            </svg>
          </div>
          <span class="panel-brand-name">My Pets</span>
        </div>

        <div class="panel-headline">
          <h2>Configura tu<br /><em>sistema.</em></h2>
          <p>Crea el primer usuario administrador para comenzar a usar la plataforma.</p>
        </div>

        <div class="panel-features">
          <div class="feature-item">
            <span class="feature-dot" />
            <span>Usuario administrador único</span>
          </div>
          <div class="feature-item">
            <span class="feature-dot" />
            <span>Control total del sistema</span>
          </div>
          <div class="feature-item">
            <span class="feature-dot" />
            <span>Solo se realiza una vez</span>
          </div>
        </div>

        <!-- Decorative shapes -->
        <div class="panel-shape panel-shape--1" />
        <div class="panel-shape panel-shape--2" />
        <div class="panel-shape panel-shape--3" />
      </div>
    </div>

    <!-- Formulario derecho -->
    <div class="auth-form-container">
      <div class="auth-form-card">
        <div class="form-header">
          <h1>Configuración inicial</h1>
          <p>Crea el usuario administrador del sistema</p>
        </div>

        <form class="form" @submit.prevent="handleSetup">
          <div class="form-group">
            <label for="name">Nombre</label>
            <div class="input-wrapper">
              <span class="input-icon">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="8" r="4"/>
                  <path d="M4 20c0-4 3.6-7 8-7s8 3 8 7"/>
                </svg>
              </span>
              <input
                id="name"
                v-model="name"
                v-bind="nameAttrs"
                type="text"
                placeholder="Tu nombre completo"
                autocomplete="name"
              />
            </div>
            <p v-if="errors.name" class="field-error">{{ errors.name }}</p>
          </div>

          <div class="form-group">
            <label for="email">Correo electrónico</label>
            <div class="input-wrapper">
              <span class="input-icon">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
                  <rect width="20" height="16" x="2" y="4" rx="2"/>
                  <path d="m22 7-8.97 5.7a1.94 1.94 0 0 1-2.06 0L2 7"/>
                </svg>
              </span>
              <input
                id="email"
                v-model="email"
                v-bind="emailAttrs"
                type="email"
                placeholder="admin@correo.com"
                autocomplete="email"
              />
            </div>
            <p v-if="errors.email" class="field-error">{{ errors.email }}</p>
          </div>

          <div class="form-group">
            <label for="password">Contraseña</label>
            <div class="input-wrapper">
              <span class="input-icon">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
                  <rect width="18" height="11" x="3" y="11" rx="2" ry="2"/>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                </svg>
              </span>
              <input
                id="password"
                v-model="password"
                v-bind="passwordAttrs"
                type="password"
                placeholder="Mínimo 8 caracteres"
                autocomplete="new-password"
              />
            </div>
            <p v-if="errors.password" class="field-error">{{ errors.password }}</p>
          </div>

          <div class="form-group">
            <label for="confirm-password">Confirmar contraseña</label>
            <div class="input-wrapper">
              <span class="input-icon">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
                  <rect width="18" height="11" x="3" y="11" rx="2" ry="2"/>
                  <path d="M7 11V7a5 5 0 0 1 10 0v4"/>
                </svg>
              </span>
              <input
                id="confirm-password"
                v-model="confirmPassword"
                v-bind="confirmPasswordAttrs"
                type="password"
                placeholder="Repite tu contraseña"
                autocomplete="new-password"
              />
            </div>
            <p v-if="errors.confirmPassword" class="field-error">{{ errors.confirmPassword }}</p>
          </div>

          <p v-if="error" class="form-error">{{ error }}</p>

          <button type="submit" class="btn-primary" :class="{ 'btn--loading': loading }">
            <span v-if="!loading">Crear administrador</span>
            <span v-else class="btn-loader">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
                <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
              </svg>
            </span>
          </button>
        </form>
      </div>
    </div>
  </div>
</template>

<style scoped>
.auth-page {
  display: flex;
  min-height: 100vh;
  background: var(--color-bg);
}

/* ── Panel decorativo ─────────────────────────────────────── */
.auth-panel {
  flex: 0 0 460px;
  background: var(--color-accent);
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

@media (max-width: 900px) {
  .auth-panel {
    display: none;
  }
}

.panel-content {
  position: relative;
  z-index: 2;
  padding: var(--space-12);
  color: var(--color-text-inverse);
  display: flex;
  flex-direction: column;
  gap: var(--space-10);
}

.panel-brand {
  display: flex;
  align-items: center;
  gap: var(--space-3);
}

.panel-icon {
  width: 44px;
  height: 44px;
  background: rgba(255, 255, 255, 0.15);
  border-radius: var(--radius-md);
  display: flex;
  align-items: center;
  justify-content: center;
  backdrop-filter: blur(8px);
}

.panel-brand-name {
  font-family: var(--font-display);
  font-size: var(--text-xl);
  font-weight: 600;
  color: #fff;
}

.panel-headline h2 {
  font-family: var(--font-display);
  font-size: var(--text-4xl);
  font-weight: 600;
  line-height: 1.2;
  color: #fff;
  margin-bottom: var(--space-4);
}

.panel-headline h2 em {
  font-style: italic;
  opacity: 0.85;
}

.panel-headline p {
  font-size: var(--text-base);
  color: rgba(255, 255, 255, 0.75);
  line-height: 1.6;
  max-width: 320px;
}

.panel-features {
  display: flex;
  flex-direction: column;
  gap: var(--space-3);
}

.feature-item {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  font-size: var(--text-sm);
  color: rgba(255, 255, 255, 0.85);
  font-weight: 500;
}

.feature-dot {
  width: 6px;
  height: 6px;
  border-radius: var(--radius-full);
  background: rgba(255, 255, 255, 0.6);
  flex-shrink: 0;
}

/* Formas decorativas */
.panel-shape {
  position: absolute;
  border-radius: var(--radius-full);
  background: rgba(255, 255, 255, 0.06);
  z-index: 1;
}

.panel-shape--1 {
  width: 320px;
  height: 320px;
  top: -80px;
  right: -80px;
}

.panel-shape--2 {
  width: 200px;
  height: 200px;
  bottom: 60px;
  right: 40px;
}

.panel-shape--3 {
  width: 120px;
  height: 120px;
  bottom: -30px;
  left: 60px;
  background: rgba(255, 255, 255, 0.04);
}

/* ── Formulario ───────────────────────────────────────────── */
.auth-form-container {
  flex: 1;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: var(--space-8);
}

.auth-form-card {
  width: 100%;
  max-width: 400px;
  display: flex;
  flex-direction: column;
  gap: var(--space-8);
}

.form-header {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.form-header h1 {
  font-family: var(--font-display);
  font-size: var(--text-3xl);
  color: var(--color-text-primary);
}

.form-header p {
  font-size: var(--text-base);
  color: var(--color-text-secondary);
}

/* ── Form fields ──────────────────────────────────────────── */
.form {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.form-group label {
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-text-primary);
}

.input-wrapper {
  position: relative;
}

.input-icon {
  position: absolute;
  left: var(--space-3);
  top: 50%;
  transform: translateY(-50%);
  color: var(--color-text-tertiary);
  display: flex;
  align-items: center;
  pointer-events: none;
}

.input-wrapper input {
  width: 100%;
  padding: var(--space-3) var(--space-4) var(--space-3) calc(var(--space-3) + 16px + var(--space-2));
  border: 1.5px solid var(--color-border);
  border-radius: var(--radius-md);
  background: var(--color-surface);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
}

.input-wrapper input::placeholder {
  color: var(--color-text-tertiary);
}

.input-wrapper input:focus {
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px var(--color-accent-light);
}

/* ── Botón ────────────────────────────────────────────────── */
.btn-primary {
  width: 100%;
  padding: var(--space-3) var(--space-6);
  background: var(--color-accent);
  color: var(--color-text-inverse);
  font-size: var(--text-sm);
  font-weight: 600;
  border-radius: var(--radius-md);
  transition: background var(--transition-fast), transform var(--transition-fast), box-shadow var(--transition-fast);
  letter-spacing: 0.01em;
  height: 46px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.btn-primary:hover {
  background: var(--color-accent-hover);
  box-shadow: var(--shadow-md);
}

.btn-primary:active {
  transform: translateY(1px);
}

.btn--loading {
  pointer-events: none;
  opacity: 0.8;
}

.btn-loader svg {
  animation: spin 0.8s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}

/* ── Error message ────────────────────────────────────────── */
.form-error {
  font-size: var(--text-sm);
  color: #dc2626;
  background: #fef2f2;
  border: 1px solid #fecaca;
  border-radius: var(--radius-md);
  padding: var(--space-2) var(--space-3);
}

.field-error {
  font-size: var(--text-xs);
  color: #dc2626;
  margin-top: var(--space-1);
}
</style>
