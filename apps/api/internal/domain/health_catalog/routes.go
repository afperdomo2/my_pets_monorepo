package health_catalog

import "github.com/gin-gonic/gin"

// RegisterRoutes monta las rutas de la guía de salud sobre el RouterGroup recibido.
func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	healthCatalog := rg.Group("/health-catalogs")
	{
		healthCatalog.GET("", h.GetHealthCatalogs)
		healthCatalog.GET("/species/:species", h.GetHealthCatalogBySpecies)
		healthCatalog.GET("/:id", h.GetHealthCatalogByID)
		healthCatalog.POST("", h.CreateHealthCatalog)
		healthCatalog.PUT("/:id", h.UpdateHealthCatalog)
		healthCatalog.DELETE("/:id", h.DeleteHealthCatalog)
	}
}
