package auth

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/my-pets/api/internal/config"
	"github.com/my-pets/api/internal/user"
)

const (
	accessTokenTTL  = 20 * time.Minute
	refreshTokenTTL = 20 * 24 * time.Hour // 20 days
)

// Claims are the JWT payload fields embedded in both access and refresh tokens.
type Claims struct {
	UserID       string `json:"uid"`
	Email        string `json:"email"`
	IsSystemUser bool   `json:"is_system_user"`
	jwt.RegisteredClaims
}

// GenerateAccessToken creates a signed JWT access token valid for 20 minutes.
func GenerateAccessToken(cfg *config.Config, u user.User) (string, error) {
	claims := Claims{
		UserID:       u.ID,
		Email:        u.Email,
		IsSystemUser: u.IsSystemUser,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(u.ID),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(accessTokenTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

// GenerateRefreshToken creates a signed JWT refresh token valid for 20 days.
// It only contains the subject (user ID) to minimize exposure.
func GenerateRefreshToken(cfg *config.Config, u user.User) (string, error) {
	claims := Claims{
		UserID:       u.ID,
		Email:        u.Email,
		IsSystemUser: u.IsSystemUser,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   fmt.Sprint(u.ID),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(refreshTokenTTL)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.JWTSecret))
}

// ValidateClaims parses and validates a JWT string, returning its claims.
func ValidateClaims(cfg *config.Config, tokenStr string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.JWTSecret), nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
