package auth

import (
	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/config"
	"github.com/my-pets/api/internal/domain/user"
)

// RegisterRoutes mounts all auth endpoints.
// Public routes are registered directly on the given group (no JWT middleware).
// Protected auth routes (e.g. /logout, /me) must be wrapped by the caller with
// the JWT middleware after calling this function.
func RegisterRoutes(rg *gin.RouterGroup, cfg *config.Config, userRepo user.Repository) *Handler {
	h := NewHandler(cfg, userRepo)
	gh := NewGoogleHandler(cfg, userRepo)

	auth := rg.Group("/auth")
	{
		// Public — no JWT required
		auth.POST("/login", h.Login)
		auth.POST("/refresh", h.Refresh)
		auth.GET("/google", gh.Initiate)
		auth.GET("/google/callback", gh.Callback)
	}

	return h
}

// RegisterProtectedRoutes mounts auth endpoints that require a valid JWT.
// Call this after applying the JWT middleware to the router group.
func RegisterProtectedRoutes(rg *gin.RouterGroup, h *Handler) {
	auth := rg.Group("/auth")
	{
		auth.POST("/logout", h.Logout)
		auth.GET("/me", h.Me)
		auth.PUT("/profile", h.UpdateProfile)
		auth.PUT("/password", h.ChangePassword)
	}
}
