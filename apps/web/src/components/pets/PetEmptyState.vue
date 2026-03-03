<script setup lang="ts">
defineProps<{ searching?: boolean; query?: string }>()
defineEmits<{ add: [] }>()
</script>

<template>
  <div class="empty-state">
    <template v-if="searching">
      <div class="empty-icon">🔍</div>
      <p class="empty-title">Sin resultados</p>
      <p class="empty-sub">No encontramos mascotas que coincidan con "<strong>{{ query }}</strong>"</p>
    </template>
    <template v-else>
      <div class="empty-icon">🐾</div>
      <p class="empty-title">Todavía no hay mascotas</p>
      <p class="empty-sub">Agregá tu primera mascota para comenzar</p>
      <button class="btn-add" @click="$emit('add')">+ Agregar mascota</button>
    </template>
  </div>
</template>

<style scoped>
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: var(--space-3);
  padding: var(--space-16) var(--space-4);
  text-align: center;
  color: var(--color-text-tertiary);
}

.empty-icon {
  font-size: 3rem;
  line-height: 1;
  margin-bottom: var(--space-2);
  filter: grayscale(0.3);
}

.empty-title {
  font-family: var(--font-display);
  font-size: var(--text-xl);
  font-weight: 600;
  color: var(--color-text-secondary);
  margin: 0;
}

.empty-sub {
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
  margin: 0;
  max-width: 280px;
}

.btn-add {
  margin-top: var(--space-2);
  padding: var(--space-2) var(--space-6);
  background: var(--color-accent);
  color: var(--color-text-on-accent);
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast), transform var(--transition-fast);
}
.btn-add:hover {
  background: var(--color-accent-hover);
  transform: translateY(-1px);
}
</style>
