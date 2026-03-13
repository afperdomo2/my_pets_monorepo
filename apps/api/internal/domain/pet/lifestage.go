package pet

import "time"

// LifeStage represents a pet's life stage category.
type LifeStage string

const (
	// Dog life stages
	LifeStagePuppy  LifeStage = "puppy"  // dog: cachorro (< 9 meses)
	LifeStageJunior LifeStage = "junior" // dog: joven adulto (9 meses – 4 años)
	LifeStageAdult  LifeStage = "adult"  // dog: adulto maduro (4 años – umbral senior)
	LifeStageSenior LifeStage = "senior" // dog: senior (umbral según tamaño de raza)

	// Cat life stages (based on age)
	LifeStageKitten      LifeStage = "kitten"       // cat: gatito (nacimiento – 1 año)
	LifeStageYoungAdult  LifeStage = "young_adult"  // cat: joven adulto (1 – 6 años)
	LifeStageMatureAdult LifeStage = "mature_adult" // cat: adulto maduro (7 – 10 años)
	LifeStageSeniorCat   LifeStage = "senior"       // cat: senior (> 10 años)
	LifeStageEndOfLife   LifeStage = "end_of_life"  // cat/any: fin de vida (cualquier edad, evaluación clínica)

	// Rabbit life stages (based on age)
	LifeStageRabbitInfant   LifeStage = "infant"   // rabbit: infancia (nacimiento – 3 meses)
	LifeStageRabbitJuvenile LifeStage = "juvenile" // rabbit: juvenil/adolescente (3 – 6 meses)
	LifeStageRabbitTeenager LifeStage = "teenager" // rabbit: teenager (6 – 12 meses)
	LifeStageRabbitAdult    LifeStage = "adult"    // rabbit: adulto (1 – 5 años)
	LifeStageRabbitSenior   LifeStage = "senior"   // rabbit: senior (> 5 años)

	// Bird life stages (based on age)
	LifeStageBirdHatchling LifeStage = "hatchling" // bird: 0–2 semanas
	LifeStageBirdNestling  LifeStage = "nestling"  // bird: 2–4 semanas
	LifeStageBirdFledgling LifeStage = "fledgling" // bird: 4–8 semanas
	LifeStageBirdJuvenile  LifeStage = "juvenile"  // bird: 2 meses – 1 año
	LifeStageBirdAdult     LifeStage = "adulto"    // bird: 1 – 10+ años
	LifeStageBirdSenior    LifeStage = "senior"    // bird: > 10–15 años
)

// seniorThreshold maps dog size to the age (in years) at which they become senior.
// Based on breed-size research.
var seniorThreshold = map[SizeCategory]float64{
	SizeSmall:  9.0, // Pequeño: senior a los 9 años
	SizeMedium: 8.0, // Mediano: senior a los 8 años
	SizeLarge:  7.0, // Grande:  senior a los 7 años
	SizeGiant:  5.0, // Gigante: senior a los 5 años
}

// CalculateDogLifeStage returns the life stage for a dog based on age and size category.
//
// Stages:
//   - puppy:  < 9 meses (~0.75 años)
//   - junior: 9 meses – 4 años
//   - senior: ≥ umbral por tamaño de raza
//   - adult:  resto (4 años – umbral senior)
func CalculateDogLifeStage(birthDate time.Time, size SizeCategory) string {
	age := time.Since(birthDate).Hours() / 24 / 365.25

	if age < 0.75 { // ~9 meses
		return string(LifeStagePuppy)
	}

	if age < 4.0 {
		return string(LifeStageJunior)
	}

	if threshold, ok := seniorThreshold[size]; ok && age >= threshold {
		return string(LifeStageSenior)
	}

	return string(LifeStageAdult)
}

// CalculateCatLifeStage returns the life stage for a cat based on age.
//
// Stages (fuente: datos clínicos felinos):
//   - kitten:       nacimiento – 1 año
//   - young_adult:  1 – 6 años
//   - mature_adult: 7 – 10 años
//   - senior:       > 10 años
func CalculateCatLifeStage(birthDate time.Time) string {
	age := time.Since(birthDate).Hours() / 24 / 365.25

	switch {
	case age < 1.0:
		return string(LifeStageKitten)
	case age < 7.0:
		return string(LifeStageYoungAdult)
	case age <= 10.0:
		return string(LifeStageMatureAdult)
	default:
		return string(LifeStageSeniorCat)
	}
}

// CalculateRabbitLifeStage returns the life stage for a rabbit based on age.
//
// Stages (fuente: datos nutricionales/clínicos lagomorfos):
//   - infant:   nacimiento – 3 meses
//   - juvenile: 3 – 6 meses
//   - teenager: 6 – 12 meses
//   - adult:    1 – 5 años
//   - senior:   > 5 años
func CalculateRabbitLifeStage(birthDate time.Time) string {
	age := time.Since(birthDate).Hours() / 24 / 365.25

	switch {
	case age < 0.25: // ~3 meses
		return string(LifeStageRabbitInfant)
	case age < 0.5: // ~6 meses
		return string(LifeStageRabbitJuvenile)
	case age < 1.0: // ~12 meses
		return string(LifeStageRabbitTeenager)
	case age <= 5.0:
		return string(LifeStageRabbitAdult)
	default:
		return string(LifeStageRabbitSenior)
	}
}

// CalculateBirdLifeStage returns the life stage for a bird based on age.
//
// Stages (based on general avian development):
//   - hatchling:  0 – 2 semanas
//   - nestling:   2 – 4 semanas
//   - fledgling:  4 – 8 semanas
//   - juvenile:   2 meses – 1 año
//   - adulto:     1 – 10+ años
//   - senior:     > 10 – 15 años
func CalculateBirdLifeStage(birthDate time.Time) string {
	age := time.Since(birthDate).Hours() / 24 / 365.25

	switch {
	case age < 0.04: // ~2 semanas
		return string(LifeStageBirdHatchling)
	case age < 0.08: // ~4 semanas
		return string(LifeStageBirdNestling)
	case age < 0.15: // ~8 semanas
		return string(LifeStageBirdFledgling)
	case age < 1.0: // ~1 año
		return string(LifeStageBirdJuvenile)
	case age <= 10.0:
		return string(LifeStageBirdAdult)
	default:
		return string(LifeStageBirdSenior)
	}
}

// CalculateLifeStage is kept for backward compatibility.
// Deprecated: use CalculateDogLifeStage, CalculateCatLifeStage or CalculateRabbitLifeStage directly.
func CalculateLifeStage(species string, weightGrams int) string {
	return ""
}
