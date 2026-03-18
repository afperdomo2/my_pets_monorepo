package dashboard

import "github.com/gin-gonic/gin"

// RegisterRoutes monta todos los endpoints de dashboard sobre el RouterGroup recibido.
//
// Rutas bajo /dashboard:
//   - GET /dashboard/summary → resumen del dashboard para el usuario autenticado
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	dashboard := rg.Group("/dashboard")
	{
		dashboard.GET("/summary", h.GetSummary)
	}
}
