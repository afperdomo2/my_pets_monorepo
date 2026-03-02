<template>
  <div class="placeholder-view">
    <div class="placeholder-icon">
      <svg width="32" height="32" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round">
        <path d="m18 2 4 4"/><path d="m17 7 3-3"/>
        <path d="M19 9 8.7 19.3c-1 1-2.5 1-3.4 0l-1.6-1.6c-1-1-1-2.5 0-3.4L14 4"/>
        <path d="m9 8 2 2"/><path d="m13 12 2 2"/>
      </svg>
    </div>
    <div class="placeholder-text">
      <h1>Vacunas</h1>
      <p>El calendario de vacunación y alertas de salud estará disponible próximamente.</p>
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
