package pet

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsValidSizeCategory_Valid(t *testing.T) {
	valid := []string{"small", "medium", "large", "giant"}
	for _, s := range valid {
		t.Run(s, func(t *testing.T) {
			require.True(t, IsValidSizeCategory(s))
		})
	}
}

func TestIsValidSizeCategory_Invalid(t *testing.T) {
	invalid := []string{"", "big", "tiny", "xlarge", "SMALL", "MEDIUM"}
	for _, s := range invalid {
		t.Run(s, func(t *testing.T) {
			require.False(t, IsValidSizeCategory(s))
		})
	}
}

func TestSizeCategoryConstants(t *testing.T) {
	require.Equal(t, SizeCategory("small"), SizeSmall)
	require.Equal(t, SizeCategory("medium"), SizeMedium)
	require.Equal(t, SizeCategory("large"), SizeLarge)
	require.Equal(t, SizeCategory("giant"), SizeGiant)
}
