package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler returns a simple liveness response.
func Handler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "my-pets-api",
	})
}
