<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
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
          <h1>Acceso restringido</h1>
          <p>El registro de nuevos usuarios es realizado por el administrador del sistema.</p>
        </div>

        <p class="info-text">
          Si necesitas acceder, comunícate con el administrador para que cree tu cuenta.
        </p>

        <button type="button" class="btn-primary" @click="router.push('/login')">
          Volver al inicio de sesión
        </button>

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

.info-text {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  padding: var(--space-3) var(--space-4);
  line-height: 1.6;
}
</style>
