<script setup lang="ts">
defineProps<{
  modelValue: number
  options: number[]
  label?: string
}>()

defineEmits<{
  'update:modelValue': [value: number]
}>()
</script>

<template>
  <div class="per-page">
    <label class="per-page__label" for="per-page-select">
      {{ label ?? 'Mostrar' }}
    </label>
    <select
      id="per-page-select"
      class="per-page__select"
      :value="modelValue"
      @change="$emit('update:modelValue', Number(($event.target as HTMLSelectElement).value))"
    >
      <option v-for="opt in options" :key="opt" :value="opt">{{ opt }}</option>
    </select>
    <span class="per-page__suffix">por página</span>
  </div>
</template>

<style scoped>
.per-page {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
}

.per-page__label,
.per-page__suffix {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
  white-space: nowrap;
}

.per-page__select {
  height: 2rem;
  padding: 0 var(--space-6) 0 var(--space-3);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-surface);
  color: var(--color-text-primary);
  font-size: var(--text-sm);
  font-family: var(--font-body);
  appearance: none;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='%236B6A66' d='M6 8L1 3h10z'/%3E%3C/svg%3E");
  background-repeat: no-repeat;
  background-position: right var(--space-2) center;
  cursor: pointer;
  transition: border-color var(--transition-fast);
}

.per-page__select:hover {
  border-color: var(--color-accent-muted);
}

.per-page__select:focus {
  border-color: var(--color-accent);
  outline: none;
}
</style>
