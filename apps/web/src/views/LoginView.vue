<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useForm } from 'vee-validate'
import { toTypedSchema } from '@vee-validate/zod'
import { useAuthStore } from '@/stores/auth'
import { authService } from '@/services/authService'
import { loginSchema } from '@/schemas/user'

const router = useRouter()
const route = useRoute()
const authStore = useAuthStore()

const error = ref<string | null>(null)

const successMessage = computed(() =>
  route.query['success'] === 'account_created'
    ? 'Usuario creado exitosamente. Inicia sesión para continuar.'
    : null,
)

onMounted(async () => {
  // Try to restore session from cookies if not already authenticated
  if (!authStore.isAuthenticated && !authStore.loading) {
    await authStore.initSession()
  }
})

const { defineField, handleSubmit, errors } = useForm({
  validationSchema: toTypedSchema(loginSchema),
  initialValues: { email: '', password: '' },
})

const [email, emailAttrs] = defineField('email')
const [password, passwordAttrs] = defineField('password')

const handleLogin = handleSubmit(async (values) => {
  error.value = null
  try {
    await authStore.login(values.email, values.password)
    router.push('/')
  } catch (e) {
    error.value = e instanceof Error ? e.message : 'Error al iniciar sesión'
  }
})

function handleGoogleLogin() {
  authService.googleLogin()
}
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
          <h2>Cuida lo que más<br /><em>importa.</em></h2>
          <p>Gestiona la salud y el bienestar de tus mascotas en un solo lugar.</p>
        </div>

        <div class="panel-features">
          <div class="feature-item">
            <span class="feature-dot" />
            <span>Historial médico completo</span>
          </div>
          <div class="feature-item">
            <span class="feature-dot" />
            <span>Calendario de vacunas</span>
          </div>
          <div class="feature-item">
            <span class="feature-dot" />
            <span>Seguimiento de bienestar</span>
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
          <h1>Bienvenido</h1>
          <p>Inicia sesión para continuar</p>
        </div>

        <form class="form" @submit.prevent="handleLogin">
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
                placeholder="tu@correo.com"
                autocomplete="email"
              />
            </div>
            <p v-if="errors.email" class="field-error">{{ errors.email }}</p>
          </div>

          <div class="form-group">
            <label for="password">
              Contraseña
              <a href="#" class="forgot-link">¿Olvidaste tu contraseña?</a>
            </label>
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
                placeholder="••••••••"
                autocomplete="current-password"
              />
            </div>
            <p v-if="errors.password" class="field-error">{{ errors.password }}</p>
          </div>

          <p v-if="successMessage" class="form-success">{{ successMessage }}</p>
          <p v-if="error" class="form-error">{{ error }}</p>

          <button type="submit" class="btn-primary" :class="{ 'btn--loading': authStore.loading }">
            <span v-if="!authStore.loading">Ingresar</span>
            <span v-else class="btn-loader">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
                <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
              </svg>
            </span>
          </button>
        </form>

        <div class="divider"><span>o</span></div>

        <button type="button" class="btn-google" @click="handleGoogleLogin">
          <svg width="18" height="18" viewBox="0 0 24 24" fill="none">
            <path d="M22.56 12.25c0-.78-.07-1.53-.2-2.25H12v4.26h5.92c-.26 1.37-1.04 2.53-2.21 3.31v2.77h3.57c2.08-1.92 3.28-4.74 3.28-8.09z" fill="#4285F4"/>
            <path d="M12 23c2.97 0 5.46-.98 7.28-2.66l-3.57-2.77c-.98.66-2.23 1.06-3.71 1.06-2.86 0-5.29-1.93-6.16-4.53H2.18v2.84C3.99 20.53 7.7 23 12 23z" fill="#34A853"/>
            <path d="M5.84 14.09c-.22-.66-.35-1.36-.35-2.09s.13-1.43.35-2.09V7.07H2.18C1.43 8.55 1 10.22 1 12s.43 3.45 1.18 4.93l3.66-2.84z" fill="#FBBC05"/>
            <path d="M12 5.38c1.62 0 3.06.56 4.21 1.64l3.15-3.15C17.45 2.09 14.97 1 12 1 7.7 1 3.99 3.47 2.18 7.07l3.66 2.84c.87-2.6 3.3-4.53 6.16-4.53z" fill="#EA4335"/>
          </svg>
          Continuar con Google
        </button>

        <p class="form-footer">
          ¿No tienes cuenta?
          <RouterLink to="/register">Crear cuenta</RouterLink>
        </p>
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
  display: flex;
  align-items: center;
  justify-content: space-between;
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-text-primary);
}

.forgot-link {
  font-size: var(--text-xs);
  font-weight: 400;
  color: var(--color-accent);
  transition: color var(--transition-fast);
}

.forgot-link:hover {
  color: var(--color-accent-dark);
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

/* ── Footer del formulario ────────────────────────────────── */
.form-footer {
  text-align: center;
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.form-footer a {
  color: var(--color-accent);
  font-weight: 600;
  margin-left: var(--space-1);
}

.form-footer a:hover {
  color: var(--color-accent-dark);
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

/* ── Success message ──────────────────────────────────────── */
.form-success {
  font-size: var(--text-sm);
  color: #16a34a;
  background: #f0fdf4;
  border: 1px solid #bbf7d0;
  border-radius: var(--radius-md);
  padding: var(--space-2) var(--space-3);
}

/* ── Divider ──────────────────────────────────────────────── */
.divider {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  color: var(--color-text-tertiary);
  font-size: var(--text-xs);
}

.divider::before,
.divider::after {
  content: '';
  flex: 1;
  height: 1px;
  background: var(--color-border);
}

/* ── Google button ────────────────────────────────────────── */
.btn-google {
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-6);
  background: var(--color-surface);
  border: 1.5px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 500;
  color: var(--color-text-primary);
  cursor: pointer;
  height: 46px;
  transition: background var(--transition-fast), box-shadow var(--transition-fast);
}

.btn-google:hover {
  background: var(--color-bg);
  box-shadow: var(--shadow-sm);
}
</style>
