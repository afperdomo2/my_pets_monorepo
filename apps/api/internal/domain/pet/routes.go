package pet

import "github.com/gin-gonic/gin"

// RegisterRoutes mounts all pet endpoints under the given router group.
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	pets := rg.Group("/pets")
	{
		pets.GET("", h.GetPets)
		pets.GET("/:id", h.GetPet)
		pets.POST("", h.CreatePet)
		pets.PUT("/:id", h.UpdatePet)
		pets.DELETE("/:id", h.DeletePet)
	}
}
