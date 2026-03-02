package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/handlers"
)

func Setup(r *gin.Engine, h *handlers.PetHandler) {
	api := r.Group("/api/v1")
	{
		pets := api.Group("/pets")
		{
			pets.GET("", h.GetPets)
			pets.GET("/:id", h.GetPet)
			pets.POST("", h.CreatePet)
			pets.PUT("/:id", h.UpdatePet)
			pets.DELETE("/:id", h.DeletePet)
		}
	}
}
