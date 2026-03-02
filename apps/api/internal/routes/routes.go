package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/handlers"
)

func Setup(r *gin.Engine) {
	api := r.Group("/api/v1")
	{
		pets := api.Group("/pets")
		{
			pets.GET("", handlers.GetPets)
			pets.GET("/:id", handlers.GetPet)
			pets.POST("", handlers.CreatePet)
			pets.PUT("/:id", handlers.UpdatePet)
			pets.DELETE("/:id", handlers.DeletePet)
		}
	}
}
