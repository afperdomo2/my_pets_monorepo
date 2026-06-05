package pet

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestCalculateDogLifeStage(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		birth    time.Time
		size     SizeCategory
		expected string
	}{
		{"puppy (< 9 meses)", now.AddDate(0, 0, -120), SizeSmall, "puppy"},
		{"junior (9m – 4a)", now.AddDate(-2, 0, 0), SizeSmall, "junior"},
		{"adult small", now.AddDate(-6, 0, 0), SizeSmall, "adult"},
		{"senior small (≥ 9a)", now.AddDate(-10, 0, 0), SizeSmall, "senior"},
		{"senior medium (≥ 8a)", now.AddDate(-9, 0, 0), SizeMedium, "senior"},
		{"senior large (≥ 7a)", now.AddDate(-8, 0, 0), SizeLarge, "senior"},
		{"senior giant (≥ 5a)", now.AddDate(-6, 0, 0), SizeGiant, "senior"},
		{"adult giant (< 5a, but ≥ 4a)", now.AddDate(-4, -6, 0), SizeGiant, "adult"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, CalculateDogLifeStage(tt.birth, tt.size))
		})
	}
}

func TestCalculateCatLifeStage(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		birth    time.Time
		expected string
	}{
		{"kitten (< 1 año)", now.AddDate(0, -6, 0), "kitten"},
		{"young adult (1-6a)", now.AddDate(-3, 0, 0), "young_adult"},
		{"mature adult (7-10a)", now.AddDate(-8, 0, 0), "mature_adult"},
		{"senior (> 10a)", now.AddDate(-12, 0, 0), "senior"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, CalculateCatLifeStage(tt.birth))
		})
	}
}

func TestCalculateRabbitLifeStage(t *testing.T) {
	now := time.Now()
	zero := time.Time{}

	tests := []struct {
		name     string
		birth    time.Time
		expected string
	}{
		{"infant (< 3 meses)", now.AddDate(0, 0, -45), "infant"},
		{"juvenile (3-6m)", now.AddDate(0, -4, 0), "juvenile"},
		{"teenager (6-12m)", now.AddDate(0, -8, 0), "teenager"},
		{"adult (1-5a)", now.AddDate(-3, 0, 0), "adult"},
		{"senior (> 5a)", now.AddDate(-8, 0, 0), "senior"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.birth == zero {
				t.Skip("birth date calculation issue")
			}
			require.Equal(t, tt.expected, CalculateRabbitLifeStage(tt.birth))
		})
	}
}

func TestCalculateBirdLifeStage(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		birth    time.Time
		expected string
	}{
		{"hatchling (< 2 semanas)", now.AddDate(0, 0, -7), "hatchling"},
		{"nestling (2-4 semanas)", now.AddDate(0, 0, -21), "nestling"},
		{"fledgling (4-8 semanas)", now.AddDate(0, 0, -40), "fledgling"},
		{"juvenile (2m – 1a)", now.AddDate(0, -4, 0), "juvenile"},
		{"adult (1-10a)", now.AddDate(-5, 0, 0), "adulto"},
		{"senior (> 10a)", now.AddDate(-12, 0, 0), "senior"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expected, CalculateBirdLifeStage(tt.birth))
		})
	}
}
