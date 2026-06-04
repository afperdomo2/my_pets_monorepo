package health

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler returns a simple liveness response.
//
//	@Summary	Verificación de salud del servicio
//	@Description	Endpoint de verificación de salud (liveness) del servicio.
//	@Tags		health
//	@Produce	json
//	@Success	200	{object}	map[string]string	"ok"
//	@Router		/health [get]
func Handler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"service": "my-pets-api",
	})
}
