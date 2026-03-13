<script setup lang="ts">
// Datos mock del historial de vacunas
const history = [
  {
    id: '1',
    petName: 'Romeo',
    petInitials: 'RO',
    petSpecies: 'dog',
    vaccine: 'Antirrábica',
    appliedDate: '15 ene 2026',
    nextDose: '15 ene 2027',
    status: 'uptodate' as const,
  },
  {
    id: '2',
    petName: 'Luna',
    petInitials: 'LU',
    petSpecies: 'cat',
    vaccine: 'Triple felina',
    appliedDate: '03 dic 2025',
    nextDose: '03 dic 2026',
    status: 'uptodate' as const,
  },
  {
    id: '3',
    petName: 'Simba',
    petInitials: 'SI',
    petSpecies: 'cat',
    vaccine: 'Antirrábica',
    appliedDate: '20 feb 2025',
    nextDose: '20 feb 2026',
    status: 'overdue' as const,
  },
  {
    id: '4',
    petName: 'Manchas',
    petInitials: 'MA',
    petSpecies: 'dog',
    vaccine: 'Polivalente canina',
    appliedDate: '10 sep 2025',
    nextDose: '10 mar 2026',
    status: 'upcoming' as const,
  },
  {
    id: '5',
    petName: 'Bolt',
    petInitials: 'BO',
    petSpecies: 'dog',
    vaccine: 'Sextuple',
    appliedDate: '01 ago 2025',
    nextDose: '01 ago 2026',
    status: 'uptodate' as const,
  },
  {
    id: '6',
    petName: 'Kira',
    petInitials: 'KI',
    petSpecies: 'cat',
    vaccine: 'Leucemia felina',
    appliedDate: '28 nov 2025',
    nextDose: '28 nov 2026',
    status: 'uptodate' as const,
  },
]

const STATUS_CONFIG = {
  uptodate: { label: 'Al día', className: 'status--uptodate' },
  upcoming: { label: 'Próxima', className: 'status--upcoming' },
  overdue:  { label: 'Vencida', className: 'status--overdue' },
}
</script>

<template>
  <div class="history-panel">
    <div class="panel-header">
      <h2>Historial de vacunación</h2>
      <span class="panel-count">{{ history.length }} registros</span>
    </div>

    <div class="table-wrapper">
      <table class="history-table">
        <thead>
          <tr>
            <th>Mascota</th>
            <th>Vacuna</th>
            <th class="th-center">Fecha aplicación</th>
            <th class="th-center">Próxima dosis</th>
            <th class="th-center">Estado</th>
          </tr>
        </thead>
        <tbody>
          <tr v-for="row in history" :key="row.id" class="history-row">
            <td>
              <div class="pet-cell">
                <div class="pet-mini-avatar" :class="`pet-mini-avatar--${row.petSpecies}`">
                  {{ row.petInitials }}
                </div>
                <span class="pet-cell-name">{{ row.petName }}</span>
              </div>
            </td>
            <td>
              <span class="vaccine-cell">{{ row.vaccine }}</span>
            </td>
            <td class="td-center">
              <span class="date-cell">{{ row.appliedDate }}</span>
            </td>
            <td class="td-center">
              <span class="date-cell">{{ row.nextDose }}</span>
            </td>
            <td class="td-center">
              <span class="status-badge" :class="STATUS_CONFIG[row.status].className">
                {{ STATUS_CONFIG[row.status].label }}
              </span>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>

<style scoped>
/* ── Panel ────────────────────────── */
.history-panel {
  background: var(--color-surface);
  border: 1px solid var(--color-border);
  border-radius: var(--radius-lg);
  overflow: hidden;
}

.panel-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: var(--space-5) var(--space-6);
  border-bottom: 1px solid var(--color-border-light);
}

.panel-header h2 {
  font-family: var(--font-display);
  font-size: var(--text-base);
  font-weight: 600;
  color: var(--color-text-primary);
}

.panel-count {
  font-size: var(--text-xs);
  color: var(--color-text-tertiary);
  font-weight: 500;
}

/* ── Tabla ────────────────────────── */
.table-wrapper {
  overflow-x: auto;
  -webkit-overflow-scrolling: touch;
}

.history-table {
  width: 100%;
  min-width: 620px;
  border-collapse: collapse;
}

.history-table thead tr {
  background: var(--color-bg);
  border-bottom: 1px solid var(--color-border-light);
}

.history-table th {
  padding: var(--space-3) var(--space-4);
  text-align: left;
  font-size: var(--text-xs);
  font-weight: 600;
  color: var(--color-text-tertiary);
  text-transform: uppercase;
  letter-spacing: 0.06em;
}

.history-table th.th-center {
  text-align: center;
}

.history-row {
  transition: background var(--transition-fast);
}

.history-row:hover {
  background: var(--color-bg-alt);
}

.history-row td {
  padding: var(--space-3) var(--space-4);
  border-bottom: 1px solid var(--color-border-light);
}

.history-row:last-child td {
  border-bottom: none;
}

.td-center {
  text-align: center;
}

/* ── Celdas ──────────────────────── */
.pet-cell {
  display: flex;
  align-items: center;
  gap: var(--space-2);
}

.pet-mini-avatar {
  width: 30px;
  height: 30px;
  border-radius: var(--radius-full);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.65rem;
  font-weight: 700;
  flex-shrink: 0;
}

.pet-mini-avatar--dog    { background: #fef3c7; color: #92400e; }
.pet-mini-avatar--cat    { background: #ede9fe; color: #5b21b6; }
.pet-mini-avatar--bird   { background: #cffafe; color: #0e7490; }
.pet-mini-avatar--rabbit { background: #fce7f3; color: #9d174d; }
.pet-mini-avatar--fish   { background: #dbeafe; color: #1e40af; }

.pet-cell-name {
  font-size: var(--text-sm);
  font-weight: 600;
  color: var(--color-text-primary);
}

.vaccine-cell {
  font-size: var(--text-sm);
  color: var(--color-text-primary);
}

.date-cell {
  font-size: var(--text-sm);
  color: var(--color-text-secondary);
}

/* ── Badges de estado ────────────── */
.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 3px var(--space-3);
  border-radius: var(--radius-full);
  font-size: var(--text-xs);
  font-weight: 600;
  white-space: nowrap;
}

.status--uptodate {
  background: #E8F5EE;
  color: #2E7D52;
}

.status--upcoming {
  background: #FEF3E2;
  color: #C4714A;
}

.status--overdue {
  background: #FEF2F2;
  color: #DC2626;
}
</style>
