package pet

// LifeStage represents a pet's life stage category.
type LifeStage string

const (
	LifeStagePuppy     LifeStage = "puppy"     // dog: cachorro
	LifeStageKitten    LifeStage = "kitten"    // cat: gatito
	LifeStageJunior    LifeStage = "junior"    // dog: joven
	LifeStageAdult     LifeStage = "adult"     // dog/cat: adulto
	LifeStageSenior    LifeStage = "senior"    // dog/cat: senior
	LifeStageGeriatric LifeStage = "geriatric" // dog: geriátrico
)

// CalculateLifeStage returns the life stage for dogs and cats based on weight in grams.
// Returns an empty string for species other than "dog" and "cat".
//
// Thresholds are weight-based approximations. After scientific review,
// these can be replaced with age+weight combined logic.
//
// Dog thresholds (by adult body weight proxy):
//   - puppy:    ≤ 2 000 g  (small/medium puppy stage)
//   - junior:   2 001 – 9 999 g  (small adult, young medium)
//   - adult:    10 000 – 29 999 g
//   - senior:   30 000 – 44 999 g
//   - geriatric: ≥ 45 000 g
//
// Cat thresholds:
//   - kitten:  ≤ 1 000 g
//   - junior:  1 001 – 2 499 g
//   - adult:   2 500 – 5 999 g
//   - senior:  ≥ 6 000 g
func CalculateLifeStage(species string, weightGrams int) string {
	switch species {
	case "dog":
		switch {
		case weightGrams <= 2000:
			return string(LifeStagePuppy)
		case weightGrams <= 9999:
			return string(LifeStageJunior)
		case weightGrams <= 29999:
			return string(LifeStageAdult)
		case weightGrams <= 44999:
			return string(LifeStageSenior)
		default:
			return string(LifeStageGeriatric)
		}
	case "cat":
		switch {
		case weightGrams <= 1000:
			return string(LifeStageKitten)
		case weightGrams <= 2499:
			return string(LifeStageJunior)
		case weightGrams <= 5999:
			return string(LifeStageAdult)
		default:
			return string(LifeStageSenior)
		}
	default:
		return ""
	}
}
