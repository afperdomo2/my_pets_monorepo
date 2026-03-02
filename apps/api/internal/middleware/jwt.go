package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/auth"
	"github.com/my-pets/api/internal/config"
)

// JWT returns a Gin middleware that validates the access_token HttpOnly cookie.
// On success it injects the following keys into the Gin context:
//
//	"userID"       uint
//	"email"        string
//	"isSystemUser" bool
//
// On failure it aborts with 401 Unauthorized.
func JWT(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenStr, err := c.Cookie("access_token")
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "authentication required"})
			return
		}

		claims, err := auth.ValidateClaims(cfg, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("isSystemUser", claims.IsSystemUser)
		c.Next()
	}
}
