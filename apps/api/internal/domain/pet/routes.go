package pet

import "github.com/gin-gonic/gin"

// RegisterRoutes mounts all pet endpoints under the given router group.
// Note: /life-stage must be registered before /:id to avoid route conflicts.
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	pets := rg.Group("/pets")
	{
		pets.GET("", h.GetPets)
		pets.GET("/life-stage", h.GetLifeStage)
		pets.GET("/:id", h.GetPet)
		pets.POST("", h.CreatePet)
		pets.PUT("/:id", h.UpdatePet)
		pets.DELETE("/:id", h.DeletePet)
	}
}
