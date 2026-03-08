package user

import "github.com/gin-gonic/gin"

// RegisterRoutes mounts all user endpoints under the given router group.
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	users := rg.Group("/users")
	{
		users.GET("", h.GetUsers)
		users.GET("/:id", h.GetUser)
		users.POST("", h.CreateUser)
		users.PUT("/:id", h.UpdateUser)
		users.DELETE("/:id", h.DeleteUser)
	}
}
