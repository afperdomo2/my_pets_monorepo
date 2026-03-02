<template>
  <div class="placeholder-view">
    <div class="placeholder-icon">
      <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <path d="M12.22 2h-.44a2 2 0 0 0-2 2v.18a2 2 0 0 1-1 1.73l-.43.25a2 2 0 0 1-2 0l-.15-.08a2 2 0 0 0-2.73.73l-.22.38a2 2 0 0 0 .73 2.73l.15.1a2 2 0 0 1 1 1.72v.51a2 2 0 0 1-1 1.74l-.15.09a2 2 0 0 0-.73 2.73l.22.38a2 2 0 0 0 2.73.73l.15-.08a2 2 0 0 1 2 0l.43.25a2 2 0 0 1 1 1.73V20a2 2 0 0 0 2 2h.44a2 2 0 0 0 2-2v-.18a2 2 0 0 1 1-1.73l.43-.25a2 2 0 0 1 2 0l.15.08a2 2 0 0 0 2.73-.73l.22-.39a2 2 0 0 0-.73-2.73l-.15-.08a2 2 0 0 1-1-1.74v-.5a2 2 0 0 1 1-1.74l.15-.09a2 2 0 0 0 .73-2.73l-.22-.38a2 2 0 0 0-2.73-.73l-.15.08a2 2 0 0 1-2 0l-.43-.25a2 2 0 0 1-1-1.73V4a2 2 0 0 0-2-2z"/>
        <circle cx="12" cy="12" r="3"/>
      </svg>
    </div>
    <div class="placeholder-text">
      <h1>Configuración</h1>
      <p>Las opciones de configuración de la cuenta y la aplicación estarán disponibles próximamente.</p>
    </div>
    <span class="placeholder-badge">Próximamente</span>
  </div>
</template>

<style scoped>
.placeholder-view {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: calc(100vh - 80px);
  gap: var(--space-5);
  text-align: center;
  padding: var(--space-8) var(--space-10);
}

.placeholder-icon {
  width: 72px;
  height: 72px;
  background: var(--color-accent-light);
  border-radius: var(--radius-xl);
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--color-accent);
}

.placeholder-text {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.placeholder-text h1 {
  font-family: var(--font-display);
  font-size: var(--text-2xl);
  color: var(--color-text-primary);
}

.placeholder-text p {
  font-size: var(--text-base);
  color: var(--color-text-secondary);
  max-width: 380px;
  line-height: 1.6;
}

.placeholder-badge {
  font-size: var(--text-xs);
  font-weight: 600;
  letter-spacing: 0.08em;
  text-transform: uppercase;
  color: var(--color-accent);
  background: var(--color-accent-light);
  padding: var(--space-1) var(--space-4);
  border-radius: var(--radius-full);
}
</style>
