<script setup lang="ts">
import { computed } from 'vue'

const props = defineProps<{
  currentPage: number
  totalPages: number
  totalItems: number
  perPage: number
}>()

const emit = defineEmits<{
  'update:page': [page: number]
}>()

function go(page: number) {
  if (page < 1 || page > props.totalPages || page === props.currentPage) return
  emit('update:page', page)
}

// Genera la secuencia de páginas con ellipsis: [1] ... [4] [5] [6] ... [10]
const pages = computed<(number | '...')[]>(() => {
  const total = props.totalPages
  const current = props.currentPage
  if (total <= 7) {
    return Array.from({ length: total }, (_, i) => i + 1)
  }
  const result: (number | '...')[] = [1]
  if (current > 3) result.push('...')
  const start = Math.max(2, current - 1)
  const end = Math.min(total - 1, current + 1)
  for (let i = start; i <= end; i++) result.push(i)
  if (current < total - 2) result.push('...')
  result.push(total)
  return result
})

const from = computed(() => Math.min((props.currentPage - 1) * props.perPage + 1, props.totalItems))
const to = computed(() => Math.min(props.currentPage * props.perPage, props.totalItems))
</script>

<template>
  <div class="pagination" aria-label="Paginación">
    <span class="pagination__info">
      {{ from }}–{{ to }} de {{ totalItems }}
    </span>

    <div class="pagination__controls">
      <!-- Primera -->
      <button
        class="pagination__btn"
        :disabled="currentPage === 1"
        aria-label="Primera página"
        @click="go(1)"
      >
        «
      </button>
      <!-- Anterior -->
      <button
        class="pagination__btn"
        :disabled="currentPage === 1"
        aria-label="Página anterior"
        @click="go(currentPage - 1)"
      >
        ‹
      </button>

      <!-- Páginas -->
      <template v-for="(p, i) in pages" :key="i">
        <span v-if="p === '...'" class="pagination__ellipsis">…</span>
        <button
          v-else
          class="pagination__btn pagination__btn--page"
          :class="{ 'pagination__btn--active': p === currentPage }"
          :aria-current="p === currentPage ? 'page' : undefined"
          @click="go(p)"
        >
          {{ p }}
        </button>
      </template>

      <!-- Siguiente -->
      <button
        class="pagination__btn"
        :disabled="currentPage === totalPages"
        aria-label="Página siguiente"
        @click="go(currentPage + 1)"
      >
        ›
      </button>
      <!-- Última -->
      <button
        class="pagination__btn"
        :disabled="currentPage === totalPages"
        aria-label="Última página"
        @click="go(totalPages)"
      >
        »
      </button>
    </div>
  </div>
</template>

<style scoped>
.pagination {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-4);
  flex-wrap: wrap;
}

.pagination__info {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

.pagination__controls {
  display: flex;
  align-items: center;
  gap: var(--space-1);
}

.pagination__btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 2rem;
  height: 2rem;
  padding: 0 var(--space-2);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-sm);
  background: var(--color-surface);
  color: var(--color-text-primary);
  font-size: var(--text-sm);
  font-family: var(--font-body);
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast), color var(--transition-fast);
}

.pagination__btn:hover:not(:disabled) {
  background: var(--color-bg-alt);
  border-color: var(--color-accent-muted);
}

.pagination__btn:disabled {
  opacity: 0.38;
  cursor: not-allowed;
}

.pagination__btn--active {
  background: var(--color-accent);
  border-color: var(--color-accent);
  color: var(--color-text-on-accent);
  font-weight: 500;
}

.pagination__btn--active:hover:not(:disabled) {
  background: var(--color-accent-hover);
  border-color: var(--color-accent-hover);
}

.pagination__ellipsis {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  min-width: 2rem;
  height: 2rem;
  font-size: var(--text-sm);
  color: var(--color-text-tertiary);
  user-select: none;
}

@media (max-width: 480px) {
  .pagination {
    justify-content: center;
  }
  .pagination__info {
    width: 100%;
    text-align: center;
  }
}
</style>
