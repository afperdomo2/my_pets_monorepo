<script setup lang="ts">
import DatePicker from '@/components/ui/DatePicker.vue'
import { useGetHealthCatalogs } from '@/composables/useHealthCatalog'
import { useCreateHealthRecord } from '@/composables/useHealthRecords'
import { HealthCatalogCategory } from '@/constants/healthRecord'
import {
  IconBell,
  IconCalendar,
  IconCheck,
  IconPlus,
  IconSearch,
  IconX,
} from '@tabler/icons-vue'
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'

onMounted(() => {
  document.body.style.overflow = 'hidden'
})

onUnmounted(() => {
  document.body.style.overflow = ''
})

const props = defineProps<{
  petId: string
  petSpecies: string
  category: 'vaccine' | 'deworming'
}>()

const emit = defineEmits<{
  close: []
}>()

const categoryMap: Record<'vaccine' | 'deworming', string> = {
  vaccine: HealthCatalogCategory.Vaccine,
  deworming: HealthCatalogCategory.Deworming,
}

const categoryRef = ref(categoryMap[props.category])
const pageRef = ref(1)
const perPageRef = ref(100)
const speciesRef = computed(() => props.petSpecies)

const { data: catalogResponse } = useGetHealthCatalogs(categoryRef, pageRef, perPageRef, speciesRef)

const catalogSearch = ref('')
const selectedCatalogId = ref<string | null>(null)
const customName = ref('')
const showCustomInput = ref(false)

const filteredCatalog = computed(() => {
  const list = catalogResponse.value?.data || []
  const q = catalogSearch.value.toLowerCase().trim()
  if (q) {
    return list.filter((v) => v.name.toLowerCase().includes(q))
  }
  return list
})

function selectCustom() {
  showCustomInput.value = true
  selectedCatalogId.value = 'custom'
  if (catalogSearch.value.trim() && !customName.value) {
    customName.value = catalogSearch.value.trim()
  }
}

const selectedCatalog = computed(() => {
  if (selectedCatalogId.value === 'custom') return null
  return catalogResponse.value?.data?.find((v) => v.id === selectedCatalogId.value)
})

watch(catalogSearch, () => {
  if (selectedCatalogId.value === 'custom' && catalogSearch.value.trim() !== customName.value) {
    selectedCatalogId.value = null
  }
})

const applicationDate = ref('')
const nextDate = ref('')
const note = ref('')
const wantsNext = ref(false)

watch(applicationDate, (newDate) => {
  if (newDate && selectedCatalog.value?.frequency_months) {
    const d = new Date(newDate + 'T12:00:00')
    d.setMonth(d.getMonth() + selectedCatalog.value.frequency_months)
    nextDate.value = d.toISOString().split('T')[0] || ''
  }
})

watch(selectedCatalog, (val) => {
  if (!val) {
    wantsNext.value = false
    nextDate.value = ''
    return
  }
  wantsNext.value = !!val.frequency_months
  if (applicationDate.value && val.frequency_months) {
    const d = new Date(applicationDate.value + 'T12:00:00')
    d.setMonth(d.getMonth() + val.frequency_months)
    nextDate.value = d.toISOString().split('T')[0] || ''
  }
})

const canSave = computed(() => {
  const hasName = selectedCatalogId.value === 'custom'
    ? customName.value.trim().length > 0
    : !!selectedCatalogId.value
  const hasDate = !!applicationDate.value
  return hasName && hasDate
})

const createRecord = useCreateHealthRecord()

async function save() {
  if (!canSave.value) return

  const name = selectedCatalogId.value === 'custom' ? customName.value : (selectedCatalog.value?.name || '')
  const catalogId = selectedCatalogId.value === 'custom' ? undefined : (selectedCatalogId.value ?? undefined)

  try {
    await createRecord.mutateAsync({
      pet_id: props.petId,
      category: props.category === 'vaccine' ? HealthCatalogCategory.Vaccine : HealthCatalogCategory.Deworming,
      name,
      health_catalog_id: catalogId,
      application_date: applicationDate.value,
      next_dose_date: wantsNext.value && nextDate.value ? nextDate.value : undefined,
      notes: note.value || undefined,
    })

    emit('close')
  } catch(e) {
    console.error(e)
  }
}

function setToday() {
  applicationDate.value = new Date().toISOString().split('T')[0] || ''
}

const modalTitle = computed(() => {
  return props.category === 'vaccine' ? 'Registrar vacuna' : 'Registrar desparasitación'
})

const productLabel = computed(() => {
  return props.category === 'vaccine' ? 'Seleccioná la vacuna' : 'Seleccioná el antiparasitario'
})

const searchPlaceholder = computed(() => {
  return props.category === 'vaccine' ? 'Buscar vacuna...' : 'Buscar antiparasitario...'
})

const frequencyText = computed(() => {
  return props.category === 'vaccine' ? 'Esta vacuna' : 'Este antiparasitario'
})
</script>

<template>
  <Teleport to="body">
    <div class="modal-backdrop" @click.self="emit('close')">
      <div class="modal-container">
        <!-- Header -->
        <div class="modal-header">
          <h2>{{ modalTitle }}</h2>
          <button class="btn-close" @click="emit('close')">
            <IconX :size="18" :stroke-width="2" />
          </button>
        </div>

        <!-- Body - Dos columnas -->
        <div class="modal-body">
          <!-- Columna izquierda: Selección de producto -->
          <div class="column column--left">
            <h3 class="column-title">{{ productLabel }}</h3>

            <!-- Botón Añadir manualmente -->
            <button
              class="btn-add-manually"
              :class="{ 'btn-add-manually--active': showCustomInput }"
              @click="selectCustom"
            >
              <div class="btn-add-manually__icon">
                <IconPlus :size="18" :stroke-width="2" />
              </div>
              <div class="btn-add-manually__text">
                <span class="btn-add-manually__title">Añadir manualmente</span>
                <span class="btn-add-manually__desc">Si no está en la lista</span>
              </div>
            </button>

            <!-- Input personalizado (aparece debajo del botón) -->
            <div v-if="showCustomInput" class="custom-input-wrapper">
              <label class="field-label">
                Nombre <span class="required">*</span>
              </label>
              <input
                v-model="customName"
                class="field-input"
                placeholder="Ej: Vacuna X, Antiparasitario Y..."
                autofocus
              />
            </div>

            <!-- Buscador -->
            <div class="search-box">
              <IconSearch class="search-icon" :size="16" :stroke-width="2" />
              <input
                v-model="catalogSearch"
                class="search-input"
                :placeholder="searchPlaceholder"
              />
            </div>

            <!-- Lista de productos -->
            <div class="product-list">
              <button
                v-for="item in filteredCatalog"
                :key="item.id"
                class="product-item"
                :class="{ 'product-item--selected': selectedCatalogId === item.id }"
                @click="selectedCatalogId = item.id; showCustomInput = false"
              >
                <div class="product-item__info">
                  <span class="product-item__name">{{ item.name }}</span>
                  <span v-if="item.description" class="product-item__desc">{{ item.description }}</span>
                  <span v-if="item.frequency_months" class="product-item__freq">
                    Cada {{ item.frequency_months }} meses
                  </span>
                </div>
                <div v-if="selectedCatalogId === item.id" class="product-item__check">
                  <IconCheck :size="16" :stroke-width="2.5" />
                </div>
              </button>
            </div>
          </div>

          <!-- Columna derecha: Fechas y notas -->
          <div class="column column--right">
            <h3 class="column-title">Fechas y detalles</h3>

            <!-- Producto seleccionado (info) -->
            <div v-if="selectedCatalog || (showCustomInput && customName)" class="selected-info">
              <IconCheck :size="16" :stroke-width="2.5" />
              <span>{{ selectedCatalog?.name || customName }}</span>
            </div>

            <!-- Fecha de aplicación (obligatoria) -->
            <div class="field-group">
              <label class="field-label">
                Fecha de aplicación <span class="required">*</span>
              </label>
              <div class="date-field">
                <div class="date-input-wrapper">
                  <DatePicker
                    v-model="applicationDate"
                    :max-date="new Date()"
                    placeholder="Seleccionar fecha"
                    unique-id="app-date-create"
                  />
                </div>
                <button
                  type="button"
                  class="btn-today"
                  @click="setToday"
                >
                  <IconCalendar :size="14" :stroke-width="2" />
                  Hoy
                </button>
              </div>
            </div>

            <!-- Programar refuerzo -->
            <div class="suggestion-card">
              <div class="suggestion-icon">
                <IconBell :size="18" :stroke-width="1.75" />
              </div>
              <div class="suggestion-body">
                <div class="suggestion-header">
                  <p class="suggestion-text">
                    <strong>Programar refuerzo</strong>
                    <template v-if="selectedCatalog?.frequency_months">
                      <br>{{ frequencyText }} se aplica cada <strong>{{ selectedCatalog.frequency_months }} meses</strong>.
                    </template>
                  </p>
                  <label class="suggestion-toggle">
                    <input v-model="wantsNext" type="checkbox" class="toggle-checkbox" />
                    <span class="toggle-track">
                      <span class="toggle-thumb" />
                    </span>
                    <span class="toggle-label">{{ wantsNext ? 'Sí' : 'No' }}</span>
                  </label>
                </div>

                <template v-if="wantsNext">
                  <div class="next-date-wrapper">
                    <label class="field-label field-label--small">
                      Fecha próximo refuerzo
                    </label>
                    <DatePicker
                      v-model="nextDate"
                      :min-date="applicationDate ? new Date(applicationDate) : new Date()"
                      placeholder="Seleccionar fecha"
                      unique-id="next-date-create"
                    />
                  </div>
                </template>
              </div>
            </div>

            <!-- Nota (opcional) -->
            <div class="field-group">
              <label class="field-label">
                Nota <span class="optional">(opcional)</span>
              </label>
              <textarea
                v-model="note"
                class="note-input"
                rows="2"
                placeholder="Ej: Lote #12345, veterinaria..."
              />
            </div>
          </div>
        </div>

        <!-- Footer -->
        <div class="modal-footer">
          <div class="footer-spacer" />
          <button
            class="btn-save"
            :disabled="!canSave || createRecord.isPending.value"
            @click="save"
          >
            <span v-if="createRecord.isPending.value" class="btn-spinner" />
            <template v-else>
              <IconCheck :size="16" :stroke-width="2.5" />
              Registrar
            </template>
          </button>
        </div>
      </div>
    </div>
  </Teleport>
</template>

<style scoped>
/* ── Backdrop ──────────────────────────────────────────────── */
.modal-backdrop {
  position: fixed;
  inset: 0;
  background: rgba(26, 26, 24, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: var(--space-4);
}

/* ── Modal container ───────────────────────────────────────── */
.modal-container {
  background: var(--color-surface);
  width: min(900px, 100%);
  max-height: min(90vh, 750px);
  border-radius: var(--radius-xl);
  box-shadow: var(--shadow-xl);
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

/* ── Header ────────────────────────────────────────────────── */
.modal-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-4) var(--space-5);
  border-bottom: 1px solid var(--color-border-light);
  flex-shrink: 0;
}

.modal-header h2 {
  font-family: var(--font-display);
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-text-primary);
  margin: 0;
}

.btn-close {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: transparent;
  border: none;
  border-radius: var(--radius-sm);
  color: var(--color-text-tertiary);
  cursor: pointer;
  transition: background var(--transition-fast), color var(--transition-fast);
}

.btn-close:hover {
  background: var(--color-bg-alt);
  color: var(--color-text-primary);
}

/* ── Body (2 columnas) ─────────────────────────────────────── */
.modal-body {
  display: grid;
  grid-template-columns: 1.2fr 1fr;
  gap: 0;
  flex: 1;
  overflow: hidden;
}

@media (max-width: 768px) {
  .modal-body {
    grid-template-columns: 1fr;
    overflow-y: auto;
  }
}

/* ── Columnas ──────────────────────────────────────────────── */
.column {
  padding: var(--space-5);
  display: flex;
  flex-direction: column;
  gap: var(--space-4);
}

.column--left {
  border-right: 1px solid var(--color-border-light);
  overflow-y: auto;
}

.column--right {
  background: var(--color-bg);
  overflow-y: auto;
}

@media (max-width: 768px) {
  .column--left {
    border-right: none;
    border-bottom: 1px solid var(--color-border-light);
  }
}

.column-title {
  font-family: var(--font-display);
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-secondary);
  margin: 0 0 var(--space-1) 0;
  text-transform: uppercase;
  letter-spacing: 0.05em;
}

/* ── Botón Añadir manualmente ─────────────────────────────── */
.btn-add-manually {
  display: flex;
  align-items: center;
  gap: var(--space-3);
  padding: var(--space-3) var(--space-4);
  background: #eff6ff;
  border: 1.5px dashed #93c5fd;
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: all var(--transition-fast);
  text-align: left;
  width: 100%;
}

.btn-add-manually:hover {
  background: #dbeafe;
  border-color: #60a5fa;
}

.btn-add-manually--active {
  background: #dbeafe;
  border-color: #3b82f6;
}

.btn-add-manually__icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #3b82f6;
  color: #fff;
  border-radius: var(--radius-md);
  flex-shrink: 0;
}

.btn-add-manually__text {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.btn-add-manually__title {
  font-size: var(--text-sm);
  font-weight: 600;
  color: #1e40af;
}

.btn-add-manually__desc {
  font-size: var(--text-xs);
  color: #64748b;
}

/* ── Input personalizado ───────────────────────────────────── */
.custom-input-wrapper {
  animation: slideDown 0.2s ease-out;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-8px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.field-label {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-secondary);
  margin-bottom: var(--space-2);
  display: block;
}

.required {
  color: #dc2626;
  font-weight: 700;
}

.optional {
  color: var(--color-text-tertiary);
  font-weight: 400;
  font-size: var(--text-xs);
}

.field-input {
  width: 100%;
  padding: var(--space-2) var(--space-3);
  min-height: 44px;
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  background: var(--color-bg);
  transition: border-color var(--transition-fast);
  box-sizing: border-box;
}

.field-input:focus {
  border-color: var(--color-accent);
  outline: none;
}

/* ── Buscador ─────────────────────────────────────────────── */
.search-box {
  position: relative;
}

.search-icon {
  position: absolute;
  left: var(--space-3);
  top: 50%;
  transform: translateY(-50%);
  color: var(--color-text-tertiary);
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: var(--space-2) var(--space-3) var(--space-2) 2.25rem;
  min-height: 44px;
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  background: var(--color-bg);
  transition: border-color var(--transition-fast);
  box-sizing: border-box;
}

.search-input:focus {
  border-color: var(--color-accent);
  outline: none;
}

/* ── Lista de productos ───────────────────────────────────── */
.product-list {
  display: flex;
  flex-direction: column;
  gap: 4px;
  max-height: 280px;
  overflow-y: auto;
}

.product-item {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: var(--space-2);
  padding: var(--space-3) var(--space-4);
  min-height: 44px;
  background: var(--color-surface);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  cursor: pointer;
  transition: background var(--transition-fast), border-color var(--transition-fast);
  text-align: left;
}

.product-item:hover {
  background: var(--color-bg-alt);
  border-color: var(--color-border);
}

.product-item--selected {
  background: var(--color-accent-light);
  border-color: var(--color-accent);
}

.product-item__info {
  display: flex;
  flex-direction: column;
  gap: 1px;
  min-width: 0;
}

.product-item__name {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
}

.product-item__desc {
  font-size: var(--text-xs);
  color: var(--color-text-secondary);
}

.product-item__freq {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
}

.product-item__check {
  color: var(--color-accent);
  flex-shrink: 0;
}

/* ── Info seleccionado ────────────────────────────────────── */
.selected-info {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-3);
  background: var(--color-surface-secondary);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
}

.selected-info svg {
  color: var(--color-success, #10b981);
  flex-shrink: 0;
}

/* ── Campos de fecha ──────────────────────────────────────── */
.field-group {
  display: flex;
  flex-direction: column;
  gap: var(--space-2);
}

.date-field {
  display: flex;
  gap: var(--space-2);
  align-items: stretch;
}

.date-input-wrapper {
  flex: 1;
  min-width: 0;
}

.btn-today {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: var(--space-1);
  padding: var(--space-2) var(--space-3);
  min-width: 100px;
  background: var(--color-accent-light);
  color: var(--color-accent-dark);
  border: 1.5px solid var(--color-accent);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: all var(--transition-fast);
  white-space: nowrap;
}

.btn-today:hover {
  background: var(--color-accent);
  color: #fff;
}

/* ── Nota ─────────────────────────────────────────────────── */
.note-input {
  padding: var(--space-3);
  border: 1.5px solid var(--color-border-light);
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  color: var(--color-text-primary);
  background: var(--color-bg);
  resize: vertical;
  min-height: 60px;
  font-family: var(--font-body);
  width: 100%;
  box-sizing: border-box;
}

.note-input:focus {
  border-color: var(--color-accent);
  outline: none;
}

/* ── Suggestion card ──────────────────────────────────────── */
.suggestion-card {
  display: flex;
  gap: var(--space-3);
  padding: var(--space-4);
  background: #f0f9ff;
  border: 1px solid #bae6fd;
  border-radius: var(--radius-lg);
  margin-top: auto;
}

.suggestion-icon {
  width: 36px;
  height: 36px;
  border-radius: var(--radius-md);
  background: #dbeafe;
  color: #2563eb;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.suggestion-body {
  display: flex;
  flex-direction: column;
  flex-grow: 1;
  min-width: 0;
}

.suggestion-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  gap: var(--space-2);
  flex-wrap: wrap;
}

.suggestion-text {
  font-size: var(--text-sm);
  color: #1e40af;
  line-height: 1.5;
  margin: 0;
}

.suggestion-toggle {
  display: flex;
  align-items: center;
  gap: var(--space-2);
  cursor: pointer;
  flex-shrink: 0;
  min-height: 32px;
}

.toggle-checkbox {
  display: none;
}

.toggle-track {
  width: 36px;
  height: 20px;
  border-radius: var(--radius-full);
  background: var(--color-border);
  position: relative;
  transition: background var(--transition-fast);
  flex-shrink: 0;
}

.toggle-checkbox:checked + .toggle-track {
  background: var(--color-accent);
}

.toggle-thumb {
  width: 16px;
  height: 16px;
  border-radius: var(--radius-full);
  background: #fff;
  position: absolute;
  top: 2px;
  left: 2px;
  transition: transform var(--transition-fast);
  box-shadow: var(--shadow-xs);
}

.toggle-checkbox:checked + .toggle-track .toggle-thumb {
  transform: translateX(16px);
}

.toggle-label {
  font-size: var(--text-xs);
  font-weight: 600;
  color: #1e40af;
}

.next-date-wrapper {
  margin-top: var(--space-3);
}

.field-label--small {
  font-size: var(--text-xs);
  margin-bottom: var(--space-1);
}

/* ── Footer ───────────────────────────────────────────────── */
.modal-footer {
  display: flex;
  align-items: center;
  padding: var(--space-3) var(--space-5);
  border-top: 1px solid var(--color-border-light);
  gap: var(--space-3);
  flex-shrink: 0;
}

.footer-spacer {
  flex: 1;
}

.btn-save {
  display: inline-flex;
  align-items: center;
  gap: var(--space-2);
  padding: var(--space-2) var(--space-5);
  min-height: 44px;
  background: #2e7d52;
  color: #fff;
  border: none;
  border-radius: var(--radius-md);
  font-size: var(--text-sm);
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition-fast);
}

.btn-save:hover:not(:disabled) {
  background: #256644;
}

.btn-save:disabled {
  background: var(--color-border);
  color: var(--color-text-tertiary);
  cursor: not-allowed;
}

.btn-spinner {
  display: inline-block;
  width: 14px;
  height: 14px;
  border: 2px solid rgba(255,255,255,0.35);
  border-top-color: #fff;
  border-radius: 50%;
  animation: spin 0.6s linear infinite;
}

@keyframes spin {
  to { transform: rotate(360deg); }
}
</style>
