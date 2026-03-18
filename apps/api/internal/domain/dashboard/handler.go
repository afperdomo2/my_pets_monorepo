package dashboard

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler holds the dependencies para los handlers HTTP de dashboard.
type Handler struct {
	repo Repository
}

// NewHandler construye un Handler con el repositorio dado.
func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

// ownerID extrae el ID del usuario autenticado del contexto Gin (seteado por el middleware JWT).
func ownerID(c *gin.Context) string {
	return c.GetString("userID")
}

// GetSummary maneja GET /api/v1/dashboard/summary
// Obtiene un resumen de los datos principales del dashboard del usuario autenticado.
//
//	@Summary	Obtener resumen del dashboard
//	@Tags		dashboard
//	@Produce	json
//	@Security	CookieAuth
//	@Success	200	{object}	map[string]interface{}	"data: DashboardSummary"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/dashboard/summary [get]
func (h *Handler) GetSummary(c *gin.Context) {
	summary, err := h.repo.GetSummary(c.Request.Context(), ownerID(c))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error al obtener resumen del dashboard"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": summary})
}
