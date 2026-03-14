<script setup lang="ts">
import { ref, computed, watch, nextTick } from 'vue'
import { IconChevronLeft, IconChevronRight, IconX } from '@tabler/icons-vue'

interface Props {
  modelValue: string | null
  maxDate?: Date
  placeholder?: string
  error?: boolean
  uniqueId?: string
}

const props = withDefaults(defineProps<Props>(), {
  maxDate: undefined,
  placeholder: 'Seleccionar fecha',
  uniqueId: () => `dp-${Math.random().toString(36).substr(2, 9)}`
})

const emit = defineEmits<{
  'update:modelValue': [value: string | null]
}>()

const isOpen = ref(false)
const currentMonth = ref(new Date().getMonth())
const currentYear = ref(new Date().getFullYear())
const inputRef = ref<HTMLButtonElement | null>(null)

const dropdownPosition = ref<{ top: number; left: number }>({ top: 0, left: 0 })
const showAbove = ref(false)

const weekDays = ['Do', 'Lu', 'Ma', 'Mi', 'Ju', 'Vi', 'Sa']
const months = [
  'Enero', 'Febrero', 'Marzo', 'Abril', 'Mayo', 'Junio',
  'Julio', 'Agosto', 'Septiembre', 'Octubre', 'Noviembre', 'Diciembre'
]

const CALENDAR_HEIGHT = 340
const GAP = 8

const selectedDate = computed(() => {
  if (!props.modelValue) return null
  const parts = props.modelValue.split('-').map(Number)
  const year = parts[0] ?? 0
  const month = parts[1] ?? 1
  const day = parts[2] ?? 1
  return new Date(year, month - 1, day)
})

const calendarDays = computed(() => {
  const year = currentYear.value
  const month = currentMonth.value
  const firstDay = new Date(year, month, 1).getDay()
  const daysInMonth = new Date(year, month + 1, 0).getDate()
  const daysInPrevMonth = new Date(year, month, 0).getDate()

  const days: { date: number | null; isCurrentMonth: boolean; dateObj: Date | null }[] = []

  for (let i = firstDay - 1; i >= 0; i--) {
    const date = new Date(year, month - 1, daysInPrevMonth - i)
    days.push({
      date: daysInPrevMonth - i,
      isCurrentMonth: false,
      dateObj: date
    })
  }

  for (let i = 1; i <= daysInMonth; i++) {
    const date = new Date(year, month, i)
    days.push({
      date: i,
      isCurrentMonth: true,
      dateObj: date
    })
  }

  const remaining = 42 - days.length
  for (let i = 1; i <= remaining; i++) {
    const date = new Date(year, month + 1, i)
    days.push({
      date: i,
      isCurrentMonth: false,
      dateObj: date
    })
  }

  return days
})

const formattedValue = computed(() => {
  if (!props.modelValue) return ''
  const [year, month, day] = props.modelValue.split('-')
  return `${day}/${month}/${year}`
})

function calculatePosition() {
  if (!inputRef.value) return

  const rect = inputRef.value.getBoundingClientRect()
  const viewportHeight = window.innerHeight
  const viewportWidth = window.innerWidth

  const spaceBelow = viewportHeight - rect.bottom
  const spaceAbove = rect.top

  const centerX = rect.left + rect.width / 2
  const dropdownWidth = Math.min(300, viewportWidth - 32)

  let left = centerX - dropdownWidth / 2
  if (left < 16) left = 16
  if (left + dropdownWidth > viewportWidth - 16) left = viewportWidth - dropdownWidth - 16

  if (spaceBelow >= CALENDAR_HEIGHT || spaceBelow >= spaceAbove) {
    showAbove.value = false
    dropdownPosition.value = {
      top: rect.bottom + GAP,
      left: left
    }
  } else if (spaceAbove >= CALENDAR_HEIGHT) {
    showAbove.value = true
    dropdownPosition.value = {
      top: rect.top - CALENDAR_HEIGHT - GAP,
      left: left
    }
  } else {
    const centerY = viewportHeight / 2 - CALENDAR_HEIGHT / 2
    showAbove.value = false
    dropdownPosition.value = {
      top: centerY,
      left: left
    }
  }
}

function isToday(dateObj: Date | null): boolean {
  if (!dateObj) return false
  const today = new Date()
  return dateObj.getDate() === today.getDate() &&
    dateObj.getMonth() === today.getMonth() &&
    dateObj.getFullYear() === today.getFullYear()
}

function isSelected(dateObj: Date | null): boolean {
  if (!dateObj || !selectedDate.value) return false
  return dateObj.getDate() === selectedDate.value.getDate() &&
    dateObj.getMonth() === selectedDate.value.getMonth() &&
    dateObj.getFullYear() === selectedDate.value.getFullYear()
}

function isDisabled(dateObj: Date | null): boolean {
  if (!dateObj) return true
  if (!props.maxDate) return false
  const max = new Date(props.maxDate)
  max.setHours(23, 59, 59, 999)
  return dateObj > max
}

function selectDay(dateObj: Date | null) {
  if (!dateObj || isDisabled(dateObj)) return
  const formatted = `${dateObj.getFullYear()}-${String(dateObj.getMonth() + 1).padStart(2, '0')}-${String(dateObj.getDate()).padStart(2, '0')}`
  emit('update:modelValue', formatted)
  isOpen.value = false
}

function prevMonth() {
  if (currentMonth.value === 0) {
    currentMonth.value = 11
    currentYear.value--
  } else {
    currentMonth.value--
  }
}

function nextMonth() {
  if (currentMonth.value === 11) {
    currentMonth.value = 0
    currentYear.value++
  } else {
    currentMonth.value++
  }
}

function toggle() {
  isOpen.value = !isOpen.value
  if (isOpen.value && props.modelValue) {
    const parts = props.modelValue.split('-').map(Number)
    currentMonth.value = (parts[1] ?? 1) - 1
    currentYear.value = parts[0] ?? new Date().getFullYear()
  }
}

function close() {
  isOpen.value = false
}

watch(isOpen, async (open) => {
  if (open) {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    if ((window as any).__activeDatePicker && (window as any).__activeDatePicker !== props.uniqueId) {
       document.dispatchEvent(new CustomEvent('close-date-picker'))
    }
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    (window as any).__activeDatePicker = props.uniqueId

    await nextTick()
    calculatePosition()
    document.addEventListener('click', handleDocumentClick)
    document.addEventListener('scroll', calculatePosition, true)
    document.addEventListener('close-date-picker', handleCloseEvent)
  } else {
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    if ((window as any).__activeDatePicker === props.uniqueId) {
       // eslint-disable-next-line @typescript-eslint/no-explicit-any
       (window as any).__activeDatePicker = null
    }
    document.removeEventListener('click', handleDocumentClick)
    document.removeEventListener('scroll', calculatePosition, true)
    document.removeEventListener('close-date-picker', handleCloseEvent)
  }
})

function handleCloseEvent() {
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  if (isOpen.value && (window as any).__activeDatePicker !== props.uniqueId) {
    isOpen.value = false
  }
}

function handleDocumentClick(event: MouseEvent) {
  const target = event.target as HTMLElement
  if (!target.closest('.date-picker') && !target.closest('.date-picker-dropdown')) {
    isOpen.value = false
  }
}
</script>

<template>
  <div class="date-picker" :class="{ 'date-picker--error': error }">
    <button
      ref="inputRef"
      type="button"
      class="date-picker__input"
      @click="toggle"
    >
      <span v-if="modelValue" class="date-picker__value">{{ formattedValue }}</span>
      <span v-else class="date-picker__placeholder">{{ placeholder }}</span>
      <IconChevronRight
        class="date-picker__icon"
        :class="{ 'date-picker__icon--open': isOpen }"
        :size="18"
      />
    </button>

    <Teleport to="body">
      <Transition name="dropdown">
        <div
          v-if="isOpen"
          class="date-picker-dropdown"
          :style="{
            top: `${dropdownPosition.top}px`,
            left: `${dropdownPosition.left}px`,
          }"
        >
          <div class="date-picker-dropdown__header">
            <button
              type="button"
              class="date-picker-dropdown__nav"
              @click="prevMonth"
            >
              <IconChevronLeft :size="20" />
            </button>
            <span class="date-picker-dropdown__title">
              {{ months[currentMonth] }} {{ currentYear }}
            </span>
            <button
              type="button"
              class="date-picker-dropdown__nav"
              @click="nextMonth"
            >
              <IconChevronRight :size="20" />
            </button>
            <button
              type="button"
              class="date-picker-dropdown__close"
              @click="close"
            >
              <IconX :size="18" />
            </button>
          </div>

          <div class="date-picker-dropdown__weekdays">
            <span v-for="day in weekDays" :key="day" class="date-picker-dropdown__weekday">
              {{ day }}
            </span>
          </div>

          <div class="date-picker-dropdown__grid">
            <button
              v-for="(day, index) in calendarDays"
              :key="index"
              type="button"
              class="date-picker-dropdown__day"
              :class="{
                'date-picker-dropdown__day--other-month': !day.isCurrentMonth,
                'date-picker-dropdown__day--today': isToday(day.dateObj),
                'date-picker-dropdown__day--selected': isSelected(day.dateObj),
                'date-picker-dropdown__day--disabled': isDisabled(day.dateObj),
              }"
              :disabled="isDisabled(day.dateObj)"
              @click="selectDay(day.dateObj)"
            >
              {{ day.date }}
            </button>
          </div>
        </div>
      </Transition>
    </Teleport>

    <Teleport to="body">
      <div
        v-if="isOpen"
        class="date-picker-overlay"
        @click="close"
      />
    </Teleport>
  </div>
</template>

<style scoped>
.date-picker {
  position: relative;
  width: 100%;
}

.date-picker__input {
  display: flex;
  align-items: center;
  justify-content: space-between;
  width: 100%;
  padding: var(--space-2) var(--space-3) var(--space-2) var(--space-8);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  background: var(--color-surface);
  cursor: pointer;
  transition: border-color var(--transition-fast), box-shadow var(--transition-fast);
}

.date-picker__input:focus {
  outline: none;
  border-color: var(--color-accent);
  box-shadow: 0 0 0 3px rgba(61, 122, 95, 0.12);
}

.date-picker--error .date-picker__input {
  border-color: var(--color-error);
}

.date-picker--error .date-picker__input:focus {
  box-shadow: 0 0 0 3px rgba(220, 38, 38, 0.12);
}

.date-picker__value {
  color: var(--color-text-primary);
}

.date-picker__placeholder {
  color: var(--color-text-tertiary);
}

.date-picker__icon {
  color: var(--color-text-tertiary);
  transition: transform var(--transition-fast);
  flex-shrink: 0;
}

.date-picker__icon--open {
  transform: rotate(90deg);
}

.date-picker-overlay {
  position: fixed;
  inset: 0;
  z-index: 999;
  background: transparent;
}

.date-picker-dropdown {
  position: fixed;
  z-index: 1000;
  background: var(--color-surface);
  border: 1px solid var(--color-border-light);
  border-radius: var(--radius-lg);
  box-shadow: var(--shadow-xl);
  padding: var(--space-3);
  width: 300px;
  max-width: calc(100vw - 32px);
}

.date-picker-dropdown__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: var(--space-3);
}

.date-picker-dropdown__nav {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border: none;
  border-radius: var(--radius-md);
  background: transparent;
  color: var(--color-text-secondary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.date-picker-dropdown__nav:hover {
  background: var(--color-bg-alt);
  color: var(--color-text-primary);
}

.date-picker-dropdown__title {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
  flex: 1;
  text-align: center;
}

.date-picker-dropdown__close {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border: none;
  border-radius: var(--radius-md);
  background: transparent;
  color: var(--color-text-tertiary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.date-picker-dropdown__close:hover {
  background: var(--color-bg-alt);
  color: var(--color-text-primary);
}

.date-picker-dropdown__weekdays {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: var(--space-1);
  margin-bottom: var(--space-2);
}

.date-picker-dropdown__weekday {
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-text-tertiary);
  height: 28px;
}

.date-picker-dropdown__grid {
  display: grid;
  grid-template-columns: repeat(7, 1fr);
  gap: var(--space-1);
}

.date-picker-dropdown__day {
  display: flex;
  align-items: center;
  justify-content: center;
  height: 36px;
  border: none;
  border-radius: var(--radius-md);
  background: transparent;
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.date-picker-dropdown__day:hover:not(:disabled) {
  background: var(--color-bg-alt);
}

.date-picker-dropdown__day--other-month {
  color: var(--color-text-tertiary);
}

.date-picker-dropdown__day--today {
  font-weight: 600;
  color: var(--color-accent);
}

.date-picker-dropdown__day--selected {
  background: var(--color-accent);
  color: var(--color-text-on-accent);
}

.date-picker-dropdown__day--selected:hover {
  background: var(--color-accent-hover);
}

.date-picker-dropdown__day--disabled {
  color: var(--color-text-tertiary);
  opacity: 0.4;
  cursor: not-allowed;
}

.dropdown-enter-active,
.dropdown-leave-active {
  transition: opacity 0.15s ease, transform 0.15s ease;
}

.dropdown-enter-from,
.dropdown-leave-to {
  opacity: 0;
  transform: scale(0.95);
}
</style>
