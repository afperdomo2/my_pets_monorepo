// Canonical size categories for dogs — must mirror backend constants in
// apps/api/internal/domain/pet/size.go

export const PET_SIZE_VALUES = ['small', 'medium', 'large', 'giant'] as const

export type PetSize = (typeof PET_SIZE_VALUES)[number]

/** Labels in Spanish for display in the UI. */
export const PET_SIZE_LABELS: Record<PetSize, string> = {
  small: 'Pequeño',
  medium: 'Mediano',
  large: 'Grande',
  giant: 'Gigante',
}

/** Weight range descriptions (approximate, in Spanish). */
export const PET_SIZE_DESCRIPTIONS: Record<PetSize, string> = {
  small: 'hasta ~10 kg',
  medium: '~10 – 25 kg',
  large: '~25 – 45 kg',
  giant: 'más de 45 kg',
}

/** Emoji / icon hint for each size card. */
export const PET_SIZE_ICONS: Record<PetSize, string> = {
  small: '🐾',
  medium: '🐕',
  large: '🦮',
  giant: '🐻',
}
