<script setup lang="ts">
const SPECIES_EMOJI: Record<string, string> = {
  dog: '🐕',
  cat: '🐈',
  bird: '🦜',
  rabbit: '🐇',
  fish: '🐠',
  other: '🐾',
}

const SPECIES_BG: Record<string, string> = {
  dog: 'bg--dog',
  cat: 'bg--cat',
  bird: 'bg--bird',
  rabbit: 'bg--rabbit',
  fish: 'bg--fish',
  other: 'bg--other',
}

interface Props {
  species: string
  size?: 'sm' | 'md' | 'lg'
  showInitials?: boolean
  name?: string
}

withDefaults(defineProps<Props>(), {
  size: 'md',
  showInitials: false,
  name: '',
})

function getEmoji(species: string): string {
  return SPECIES_EMOJI[species] ?? '🐾'
}

function getBgClass(species: string): string {
  return SPECIES_BG[species] ?? 'bg--other'
}

function getInitials(name: string): string {
  if (!name) return '??'
  return name.slice(0, 2).toUpperCase()
}
</script>

<template>
  <div class="pet-avatar" :class="[getBgClass(species), `size--${size}`]">
    <span v-if="size === 'sm'" class="avatar-content">
      {{ showInitials ? getInitials(name) : getEmoji(species) }}
    </span>
    <span v-else-if="size === 'md'" class="avatar-content">
      {{ showInitials ? getInitials(name) : getEmoji(species) }}
    </span>
    <span v-else class="avatar-content">
      {{ showInitials ? getInitials(name) : getEmoji(species) }}
    </span>
  </div>
</template>

<style scoped>
.pet-avatar {
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-full);
  flex-shrink: 0;
  font-weight: 700;
}

/* Background colors per species */
.bg--dog    { background: #fef3c7; color: #92400e; }
.bg--cat    { background: #ede9fe; color: #5b21b6; }
.bg--bird   { background: #cffafe; color: #0e7490; }
.bg--rabbit { background: #fce7f3; color: #9d174d; }
.bg--fish   { background: #dbeafe; color: #1e40af; }
.bg--other  { background: var(--color-accent-light); color: var(--color-accent-dark); }

/* Sizes */
.size--sm {
  width: 24px;
  height: 24px;
  font-size: 0.75rem;
}

.size--md {
  width: 30px;
  height: 30px;
  font-size: 0.875rem;
}

.size--lg {
  width: 48px;
  height: 48px;
  font-size: 1.5rem;
}

.avatar-content {
  display: flex;
  align-items: center;
  justify-content: center;
  line-height: 1;
}
</style>
