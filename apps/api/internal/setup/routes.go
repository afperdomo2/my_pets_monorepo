package setup

import "github.com/gin-gonic/gin"

// RegisterRoutes mounts the setup endpoints (all public — no JWT required).
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	rg.GET("/setup/status", h.Status)
	rg.POST("/setup", h.Create)
}
