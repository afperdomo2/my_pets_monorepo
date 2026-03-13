export const PET_SPECIES = [
  { value: 'dog', label: 'Perro', icon: '🐕' },
  { value: 'cat', label: 'Gato', icon: '🐈' },
  { value: 'bird', label: 'Ave', icon: '🦜' },
  { value: 'rabbit', label: 'Conejo', icon: '🐇' },
  { value: 'fish', label: 'Pez', icon: '🐠' },
  { value: 'other', label: 'Otro', icon: '🐾' },
] as const

export type PetSpecies = typeof PET_SPECIES[number]['value']

export function getSpeciesLabel(value: string): string {
  return PET_SPECIES.find(s => s.value === value)?.label ?? value
}

export function getSpeciesValue(label: string): string {
  return PET_SPECIES.find(s => s.label === label)?.value ?? label
}

export function getSpeciesIcon(value: string): string {
  return PET_SPECIES.find(s => s.value === value)?.icon ?? '🐾'
}
