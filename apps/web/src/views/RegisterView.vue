<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const form = ref({
  name: '',
  email: '',
  password: '',
  confirmPassword: '',
})

const loading = ref(false)

function handleRegister() {
  loading.value = true
  // Sin lógica real por ahora — redirige directo al home
  setTimeout(() => {
    loading.value = false
    router.push('/')
  }, 700)
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
          <h2>Empieza tu<br /><em>historia.</em></h2>
          <p>Crea tu cuenta y lleva un registro completo del bienestar de tus mascotas.</p>
        </div>

        <div class="panel-features">
          <div class="feature-item">
            <span class="feature-dot" />
            <span>Registro ilimitado de mascotas</span>
          </div>
          <div class="feature-item">
            <span class="feature-dot" />
            <span>Alertas de vacunación</span>
          </div>
          <div class="feature-item">
            <span class="feature-dot" />
            <span>Reportes de salud</span>
          </div>
        </div>

        <div class="panel-shape panel-shape--1" />
        <div class="panel-shape panel-shape--2" />
        <div class="panel-shape panel-shape--3" />
      </div>
    </div>

    <!-- Formulario derecho -->
    <div class="auth-form-container">
      <div class="auth-form-card">
        <div class="form-header">
          <h1>Crear cuenta</h1>
          <p>Completa los datos para comenzar</p>
        </div>

        <form class="form" @submit.prevent="handleRegister">
          <div class="form-group">
            <label for="name">Nombre completo</label>
            <div class="input-wrapper">
              <span class="input-icon">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
                  <circle cx="12" cy="8" r="4"/><path d="M4 20c0-4 3.6-7 8-7s8 3 8 7"/>
                </svg>
              </span>
              <input
                id="name"
                v-model="form.name"
                type="text"
                placeholder="Juan García"
                autocomplete="name"
              />
            </div>
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
                v-model="form.email"
                type="email"
                placeholder="tu@correo.com"
                autocomplete="email"
              />
            </div>
          </div>

          <div class="form-row">
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
                  v-model="form.password"
                  type="password"
                  placeholder="Mínimo 8 caracteres"
                  autocomplete="new-password"
                />
              </div>
            </div>

            <div class="form-group">
              <label for="confirm">Confirmar</label>
              <div class="input-wrapper">
                <span class="input-icon">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.75" stroke-linecap="round" stroke-linejoin="round">
                    <path d="M12 22s8-4 8-10V5l-8-3-8 3v7c0 6 8 10 8 10z"/>
                  </svg>
                </span>
                <input
                  id="confirm"
                  v-model="form.confirmPassword"
                  type="password"
                  placeholder="Repetir contraseña"
                  autocomplete="new-password"
                />
              </div>
            </div>
          </div>

          <button type="submit" class="btn-primary" :class="{ 'btn--loading': loading }">
            <span v-if="!loading">Crear cuenta</span>
            <span v-else class="btn-loader">
              <svg width="18" height="18" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round">
                <path d="M21 12a9 9 0 1 1-6.219-8.56"/>
              </svg>
            </span>
          </button>
        </form>

        <p class="form-footer">
          ¿Ya tienes cuenta?
          <RouterLink to="/login">Iniciar sesión</RouterLink>
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
  flex: 0 0 420px;
  background: var(--color-accent-dark);
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
}

@media (max-width: 900px) {
  .auth-panel { display: none; }
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
  max-width: 300px;
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

.panel-shape {
  position: absolute;
  border-radius: var(--radius-full);
  background: rgba(255, 255, 255, 0.05);
  z-index: 1;
}

.panel-shape--1 { width: 300px; height: 300px; top: -80px; right: -80px; }
.panel-shape--2 { width: 180px; height: 180px; bottom: 80px; right: 30px; }
.panel-shape--3 { width: 100px; height: 100px; bottom: -20px; left: 50px; background: rgba(255,255,255,0.03); }

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
  max-width: 440px;
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

.form {
  display: flex;
  flex-direction: column;
  gap: var(--space-5);
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: var(--space-4);
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
</style>
