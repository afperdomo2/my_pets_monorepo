<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(defineProps<{
  current: number
  limit: number
  showLabel?: 'top' | 'bottom' | 'none'
  labelText?: string
  colorGreen?: string
  colorYellow?: string
  colorOrange?: string
  colorFull?: string
  warningThreshold?: number
  dangerThreshold?: number
}>(), {
  showLabel: 'none',
  colorGreen: '#22c55e',
  colorYellow: '#eab308',
  colorOrange: '#f97316',
  colorFull: '#dc2626',
  warningThreshold: 50,
  dangerThreshold: 80,
})

const percentage = computed(() => {
  if (props.limit <= 0) return 0
  return Math.min((props.current / props.limit) * 100, 100)
})

const barColor = computed(() => {
  if (percentage.value >= 100) return props.colorFull
  if (percentage.value >= props.dangerThreshold) return props.colorOrange
  if (percentage.value >= props.warningThreshold) return props.colorYellow
  return props.colorGreen
})

const labelDisplay = computed(() => props.labelText || `${props.current} / ${props.limit}`)
</script>

<template>
  <div class="resource-progress">
    <span v-if="showLabel === 'top'" class="label label--top">
      {{ labelDisplay }}
    </span>
    <div class="progress-bar">
      <div
        class="progress-fill"
        :style="{ width: `${percentage}%`, backgroundColor: barColor }"
      />
    </div>
    <span v-if="showLabel === 'bottom'" class="label label--bottom">
      {{ labelDisplay }}
    </span>
  </div>
</template>

<style scoped>
.resource-progress {
  display: flex;
  flex-direction: column;
  gap: var(--space-1);
  min-width: 80px;
}

.progress-bar {
  height: 6px;
  background: var(--color-border-light);
  border-radius: var(--radius-full);
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: var(--radius-full);
  transition: width 0.3s ease, background-color 0.3s ease;
}

.label {
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-text-secondary);
}

.label--top {
  text-align: left;
}

.label--bottom {
  text-align: right;
}
</style>
