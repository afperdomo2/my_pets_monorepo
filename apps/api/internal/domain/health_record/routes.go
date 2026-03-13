package health_record

import "github.com/gin-gonic/gin"

// RegisterRoutes monta todos los endpoints de health_record bajo el grupo de rutas dado.
// Se reutiliza el parámetro :id del grupo /pets/:id para evitar conflictos con el router de Gin.
// Dentro de los handlers, c.Param("id") contiene el pet_id.
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	records := rg.Group("/pets/:id/health-records")
	{
		records.GET("", h.GetHealthRecords)
		records.GET("/:record_id", h.GetHealthRecord)
		records.POST("", h.CreateHealthRecord)
		records.PUT("/:record_id", h.UpdateHealthRecord)
		records.PATCH("/:record_id/status", h.UpdateHealthRecordStatus)
		records.DELETE("/:record_id", h.DeleteHealthRecord)
	}
}
