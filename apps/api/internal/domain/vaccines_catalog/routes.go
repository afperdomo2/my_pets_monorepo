package vaccines_catalog

import "github.com/gin-gonic/gin"

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	vaccinesCatalog := rg.Group("/vaccines-catalog")
	{
		vaccinesCatalog.GET("", h.GetVaccinesCatalog)
		vaccinesCatalog.GET("/species/:species", h.GetVaccinesBySpecies)
		vaccinesCatalog.GET("/:id", h.GetVaccineCatalogByID)

		vaccinesCatalog.POST("", h.CreateVaccineCatalog)
		vaccinesCatalog.PUT("/:id", h.UpdateVaccineCatalog)
		vaccinesCatalog.DELETE("/:id", h.DeleteVaccineCatalog)
	}
}
