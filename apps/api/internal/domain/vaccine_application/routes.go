package vaccine_application

import "github.com/gin-gonic/gin"

// RegisterRoutes monta todos los endpoints de vaccine_application sobre el RouterGroup recibido.
//
// Rutas bajo /vaccine-applications:
//   - GET    /vaccine-applications/health-record/:id  → listar aplicaciones de un health_record
//   - GET    /vaccine-applications/:id                → obtener aplicación por ID
//   - POST   /vaccine-applications                    → crear aplicación
//   - PUT    /vaccine-applications/:id                → actualizar aplicación
//   - DELETE /vaccine-applications/:id                → eliminar aplicación
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	vaccineApps := rg.Group("/vaccine-applications")
	{
		// Listar aplicaciones de un health_record específico
		vaccineApps.GET("/health-record/:id", h.GetApplicationsByHealthRecord)

		// Obtener, crear, actualizar y eliminar aplicaciones
		vaccineApps.GET("/:id", h.GetApplicationByID)
		vaccineApps.POST("", h.CreateApplication)
		vaccineApps.PUT("/:id", h.UpdateApplication)
		vaccineApps.DELETE("/:id", h.DeleteApplication)
	}
}
