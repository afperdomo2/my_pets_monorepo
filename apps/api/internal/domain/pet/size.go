package pet

// SizeCategory represents the physical size classification of a dog.
// Only applicable when species == "dog"; for other species the field is nil.
type SizeCategory string

const (
	SizeSmall  SizeCategory = "small"  // up to ~10 kg  (e.g. Chihuahua, Poodle Toy)
	SizeMedium SizeCategory = "medium" // ~10–25 kg     (e.g. Beagle, Border Collie)
	SizeLarge  SizeCategory = "large"  // ~25–45 kg     (e.g. Labrador, German Shepherd)
	SizeGiant  SizeCategory = "giant"  // 45 kg+        (e.g. Great Dane, Saint Bernard)
)

// validSizeCategories is the canonical set of accepted size values.
var validSizeCategories = map[SizeCategory]bool{
	SizeSmall:  true,
	SizeMedium: true,
	SizeLarge:  true,
	SizeGiant:  true,
}

// IsValidSizeCategory returns true if s is one of the accepted size values.
func IsValidSizeCategory(s string) bool {
	return validSizeCategories[SizeCategory(s)]
}
