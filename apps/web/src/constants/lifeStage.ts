/**
 * Constantes de etapas de vida (life stages) para mascotas.
 * Estos valores deben coincidir exactamente con los definidos en el backend:
 * apps/api/internal/domain/pet/lifestage.go
 */

// ==================== PERRO ====================

/**
 * Etapas de vida para perros según edad y tamaño.
 * Objeto const para usar como enum: LifeStageForDog.Puppy
 */
export const LifeStageForDog = {
  Puppy: 'puppy',
  Junior: 'junior',
  Adult: 'adult',
  Senior: 'senior',
} as const

export type DogLifeStage = (typeof LifeStageForDog)[keyof typeof LifeStageForDog]

/**
 * Labels en español para etapas de vida de perros.
 */
export const DOG_LIFE_STAGE_LABELS: Record<DogLifeStage, string> = {
  [LifeStageForDog.Puppy]: 'Cachorro',
  [LifeStageForDog.Junior]: 'Joven adulto',
  [LifeStageForDog.Adult]: 'Adulto',
  [LifeStageForDog.Senior]: 'Senior',
}

/**
 * Descripciones de cada etapa de vida para perros.
 */
export const DOG_LIFE_STAGE_DESCRIPTIONS: Record<DogLifeStage, string> = {
  [LifeStageForDog.Puppy]: 'Cachorro (0-9 meses)',
  [LifeStageForDog.Junior]: 'Joven adulto (9 meses - 4 años)',
  [LifeStageForDog.Adult]: 'Adulto maduro (4 años - umbral senior)',
  [LifeStageForDog.Senior]: 'Senior (umbral según tamaño de raza)',
}

// ==================== GATO ====================

/**
 * Etapas de vida para gatos según edad.
 */
export const LifeStageForCat = {
  Kitten: 'kitten',
  YoungAdult: 'young_adult',
  MatureAdult: 'mature_adult',
  Senior: 'senior',
} as const

export type CatLifeStage = (typeof LifeStageForCat)[keyof typeof LifeStageForCat]

/**
 * Labels en español para etapas de vida de gatos.
 */
export const CAT_LIFE_STAGE_LABELS: Record<CatLifeStage, string> = {
  [LifeStageForCat.Kitten]: 'Gatito',
  [LifeStageForCat.YoungAdult]: 'Joven adulto',
  [LifeStageForCat.MatureAdult]: 'Adulto maduro',
  [LifeStageForCat.Senior]: 'Senior',
}

/**
 * Descripciones de cada etapa de vida para gatos.
 */
export const CAT_LIFE_STAGE_DESCRIPTIONS: Record<CatLifeStage, string> = {
  [LifeStageForCat.Kitten]: 'Gatito (nacimiento - 1 año)',
  [LifeStageForCat.YoungAdult]: 'Joven adulto (1 - 6 años)',
  [LifeStageForCat.MatureAdult]: 'Adulto maduro (7 - 10 años)',
  [LifeStageForCat.Senior]: 'Senior (> 10 años)',
}

// ==================== CONEJO ====================

/**
 * Etapas de vida para conejos según edad.
 */
export const LifeStageForRabbit = {
  Infant: 'infant',
  Juvenile: 'juvenile',
  Teenager: 'teenager',
  Adult: 'adult',
  Senior: 'senior',
} as const

export type RabbitLifeStage = (typeof LifeStageForRabbit)[keyof typeof LifeStageForRabbit]

/**
 * Labels en español para etapas de vida de conejos.
 */
export const RABBIT_LIFE_STAGE_LABELS: Record<RabbitLifeStage, string> = {
  [LifeStageForRabbit.Infant]: 'Infancia',
  [LifeStageForRabbit.Juvenile]: 'Juvenil',
  [LifeStageForRabbit.Teenager]: 'Adolescente',
  [LifeStageForRabbit.Adult]: 'Adulto',
  [LifeStageForRabbit.Senior]: 'Senior',
}

/**
 * Descripciones de cada etapa de vida para conejos.
 */
export const RABBIT_LIFE_STAGE_DESCRIPTIONS: Record<RabbitLifeStage, string> = {
  [LifeStageForRabbit.Infant]: 'Infancia (nacimiento - 3 meses)',
  [LifeStageForRabbit.Juvenile]: 'Juvenil (3 - 6 meses)',
  [LifeStageForRabbit.Teenager]: 'Adolescente (6 - 12 meses)',
  [LifeStageForRabbit.Adult]: 'Adulto (1 - 5 años)',
  [LifeStageForRabbit.Senior]: 'Senior (> 5 años)',
}

// ==================== AVE ====================

/**
 * Etapas de vida para aves según edad.
 */
export const LifeStageForBird = {
  Hatchling: 'hatchling',
  Nestling: 'nestling',
  Fledgling: 'fledgling',
  Juvenile: 'juvenile',
  Adult: 'adult',
  Senior: 'senior',
} as const

export type BirdLifeStage = (typeof LifeStageForBird)[keyof typeof LifeStageForBird]

/**
 * Labels en español para etapas de vida de aves.
 */
export const BIRD_LIFE_STAGE_LABELS: Record<BirdLifeStage, string> = {
  [LifeStageForBird.Hatchling]: 'Eclosionado',
  [LifeStageForBird.Nestling]: 'Nidícola',
  [LifeStageForBird.Fledgling]: 'Volantón',
  [LifeStageForBird.Juvenile]: 'Juvenil',
  [LifeStageForBird.Adult]: 'Adulto',
  [LifeStageForBird.Senior]: 'Senior',
}

/**
 * Descripciones de cada etapa de vida para aves.
 */
export const BIRD_LIFE_STAGE_DESCRIPTIONS: Record<BirdLifeStage, string> = {
  [LifeStageForBird.Hatchling]: 'Eclosionado (0 - 2 semanas)',
  [LifeStageForBird.Nestling]: 'Nidícola (2 - 4 semanas)',
  [LifeStageForBird.Fledgling]: 'Volantón (4 - 8 semanas)',
  [LifeStageForBird.Juvenile]: 'Juvenil (2 meses - 1 año)',
  [LifeStageForBird.Adult]: 'Adulto (1 - 10+ años)',
  [LifeStageForBird.Senior]: 'Senior (> 10 - 15 años)',
}

// ==================== UNION Y HELPERS ====================

/**
 * Unión de todas las etapas de vida posibles.
 */
export type LifeStage = DogLifeStage | CatLifeStage | RabbitLifeStage | BirdLifeStage

/**
 * Mapa completo de labels en español para cualquier life stage.
 */
export const LIFE_STAGE_LABELS: Record<LifeStage, string> = {
  ...DOG_LIFE_STAGE_LABELS,
  ...CAT_LIFE_STAGE_LABELS,
  ...RABBIT_LIFE_STAGE_LABELS,
  ...BIRD_LIFE_STAGE_LABELS,
}

/**
 * Obtiene el label en español para una etapa de vida.
 * @param stage - El valor de la etapa de vida
 * @returns El label en español, o el mismo valor si no se encuentra
 */
export function getLifeStageLabel(stage: string): string {
  return LIFE_STAGE_LABELS[stage as LifeStage] ?? stage
}

/**
 * Valida si un string es una etapa de vida válida.
 */
export function isValidLifeStage(stage: string): stage is LifeStage {
  return Object.values({ ...LifeStageForDog, ...LifeStageForCat, ...LifeStageForRabbit, ...LifeStageForBird }).includes(stage as LifeStage)
}
