package health_record

import "github.com/gin-gonic/gin"

// RegisterRoutes monta todos los endpoints de health_record sobre el RouterGroup recibido.
//
// Rutas bajo /health-records:
//   - GET    /health-records                                    → todos los registros del usuario
//   - POST   /health-records                                    → crear (pet_id en body)
//   - PUT    /health-records/:record_id                         → actualizar
//   - PATCH  /health-records/:record_id/status                  → cambiar status
//   - DELETE /health-records/:record_id                         → eliminar
//   - GET    /health-records/pets/:pet_id/category/:category    → por mascota y categoría (más específico, primero)
//   - GET    /health-records/pets/:pet_id                       → por mascota
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	healthRecords := rg.Group("/health-records")
	{
		healthRecords.GET("", h.GetAllHealthRecords)
		healthRecords.POST("", h.CreateHealthRecord)
		healthRecords.PUT("/:record_id", h.UpdateHealthRecord)
		healthRecords.PATCH("/:record_id/status", h.UpdateHealthRecordStatus)
		healthRecords.DELETE("/:record_id", h.DeleteHealthRecord)

		// Rutas de listado por mascota — la más específica va primero para evitar
		// que Gin interprete "pets/:pet_id/category/:category" como registro con record_id="pets".
		pets := healthRecords.Group("/pets")
		{
			pets.GET("/:pet_id/category/:category", h.GetHealthRecordsByPetAndCategory)
			pets.GET("/:pet_id", h.GetHealthRecordsByPet)
		}
	}
}
