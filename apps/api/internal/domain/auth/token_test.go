package auth

import (
	"testing"

	"github.com/my-pets/api/internal/config"
	"github.com/my-pets/api/internal/models"
	"github.com/stretchr/testify/require"
)

func testCfg() *config.Config {
	return &config.Config{JWTSecret: "test-secret-for-unit-test"}
}

func TestGenerateAccessToken_Success(t *testing.T) {
	u := models.User{ID: "user-1", Email: "test@test.com", IsSystemUser: false}
	token, err := GenerateAccessToken(testCfg(), u)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestGenerateRefreshToken_Success(t *testing.T) {
	u := models.User{ID: "user-1", Email: "test@test.com", IsSystemUser: false}
	token, err := GenerateRefreshToken(testCfg(), u)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestGenerateRefreshToken_SystemUser(t *testing.T) {
	u := models.User{ID: "admin-1", Email: "admin@test.com", IsSystemUser: true}
	token, err := GenerateRefreshToken(testCfg(), u)
	require.NoError(t, err)
	require.NotEmpty(t, token)
}

func TestValidateClaims_ValidToken(t *testing.T) {
	u := models.User{ID: "user-1", Email: "test@test.com", IsSystemUser: false}
	token, err := GenerateAccessToken(testCfg(), u)
	require.NoError(t, err)

	claims, err := ValidateClaims(testCfg(), token)
	require.NoError(t, err)
	require.Equal(t, "user-1", claims.UserID)
	require.Equal(t, "test@test.com", claims.Email)
	require.False(t, claims.IsSystemUser)
}

func TestValidateClaims_InvalidToken(t *testing.T) {
	_, err := ValidateClaims(testCfg(), "invalid-token-string")
	require.Error(t, err)
}

func TestValidateClaims_WrongSecret(t *testing.T) {
	u := models.User{ID: "user-1", Email: "test@test.com"}
	token, err := GenerateAccessToken(testCfg(), u)
	require.NoError(t, err)

	otherCfg := &config.Config{JWTSecret: "different-secret"}
	_, err = ValidateClaims(otherCfg, token)
	require.Error(t, err)
}

func TestValidateClaims_RefreshToken(t *testing.T) {
	u := models.User{ID: "user-1", Email: "test@test.com", IsSystemUser: true}
	token, err := GenerateRefreshToken(testCfg(), u)
	require.NoError(t, err)

	claims, err := ValidateClaims(testCfg(), token)
	require.NoError(t, err)
	require.Equal(t, "user-1", claims.UserID)
	require.True(t, claims.IsSystemUser)
}
