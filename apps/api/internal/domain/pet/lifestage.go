package pet

import "time"

// LifeStage represents a pet's life stage category.
type LifeStage string

const (
	LifeStagePuppy  LifeStage = "puppy"  // dog: cachorro (< 9 meses)
	LifeStageJunior LifeStage = "junior" // dog: joven adulto (9 meses – 4 años)
	LifeStageAdult  LifeStage = "adult"  // dog: adulto maduro (4 años – umbral senior)
	LifeStageSenior LifeStage = "senior" // dog: senior (umbral según tamaño de raza)

	LifeStageKitten    LifeStage = "kitten" // cat: gatito  (≤ 1 000 g)
	LifeStageAdultCat  LifeStage = "adult"  // cat: adulto  (2 500 – 5 999 g)
	LifeStageSeniorCat LifeStage = "senior" // cat: senior  (≥ 6 000 g)
	LifeStageJuniorCat LifeStage = "junior" // cat: joven   (1 001 – 2 499 g)
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

// CalculateCatLifeStage returns the life stage for a cat based on weight in grams.
func CalculateCatLifeStage(weightGrams int) string {
	switch {
	case weightGrams <= 1000:
		return string(LifeStageKitten)
	case weightGrams <= 2499:
		return string(LifeStageJuniorCat)
	case weightGrams <= 5999:
		return string(LifeStageAdultCat)
	default:
		return string(LifeStageSeniorCat)
	}
}

// CalculateLifeStage is kept for backward compatibility.
// For dogs it now requires birthDate and size; for cats it uses weightGrams.
// Deprecated: use CalculateDogLifeStage or CalculateCatLifeStage directly.
func CalculateLifeStage(species string, weightGrams int) string {
	if species == "cat" {
		return CalculateCatLifeStage(weightGrams)
	}
	return ""
}
