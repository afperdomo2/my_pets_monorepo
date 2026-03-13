package vaccines_catalog

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/validation"
)

// Handler holds the dependencies for vaccine catalog HTTP handlers.
type Handler struct {
	repo Repository
}

// NewHandler constructs a Handler with the given Repository.
func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

// parseID extracts and validates the :id path parameter.
func parseID(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return "", false
	}
	return id, true
}

// requireSystemUser returns true if the caller is a system user.
// If not, it responds with 403 Forbidden and returns false.
func requireSystemUser(c *gin.Context) bool {
	isSystemUser, _ := c.Get("isSystemUser")
	if isSystemUser != true {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return false
	}
	return true
}

// GetVaccinesCatalog handles GET /api/v1/vaccines-catalog
// Accessible by any authenticated user.
//
//	@Summary	Listar vacunas del catálogo (paginado)
//	@Tags		vaccines-catalog
//	@Produce	json
//	@Security	CookieAuth
//	@Param		page		query		int						false	"Número de página (por defecto 1)"
//	@Param		per_page	query		int						false	"Elementos por página (por defecto 10)"
//	@Param		species	query		string					false	"Filtrar por especie (dog, cat, bird, rabbit, fish, other)"
//	@Success	200	{object}	map[string]interface{}	"data: []VaccineCatalog, total: int, page: int, per_page: int, total_pages: int"
//	@Failure	401	{object}	map[string]string	"autenticación requerida"
//	@Failure	500	{object}	map[string]string	"mensaje de error"
//	@Router		/api/v1/vaccines-catalog [get]
func (h *Handler) GetVaccinesCatalog(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	species := c.Query("species")

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	// Validar species si se proporciona
	var speciesFilter *string
	if species != "" {
		validSpecies := map[string]bool{
			"dog":    true,
			"cat":    true,
			"bird":   true,
			"rabbit": true,
			"fish":   true,
			"other":  true,
		}
		if !validSpecies[species] {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid species filter"})
			return
		}
		speciesFilter = &species
	}

	vaccines, total, err := h.repo.GetPaginated(c.Request.Context(), page, perPage, speciesFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch vaccines catalog"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	c.JSON(http.StatusOK, gin.H{
		"data":        vaccines,
		"total":       total,
		"page":        page,
		"per_page":    perPage,
		"total_pages": totalPages,
	})
}

// GetVaccineCatalogByID handles GET /api/v1/vaccines-catalog/:id
// Accessible by any authenticated user.
//
//	@Summary	Obtener una vaccine del catálogo por ID
//	@Tags		vaccines-catalog
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string					true	"ID de la vaccine"
//	@Success	200	{object}	map[string]interface{}	"data: VaccineCatalog"
//	@Failure	400	{object}	map[string]string		"id inválido"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	404	{object}	map[string]string		"vaccine no encontrada"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/vaccines-catalog/{id} [get]
func (h *Handler) GetVaccineCatalogByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	v, err := h.repo.GetByID(c.Request.Context(), id)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "vaccine catalog not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch vaccine catalog"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": v})
}

// GetVaccinesBySpecies handles GET /api/v1/vaccines-catalog/species/:species
// Accessible by any authenticated user.
//
//	@Summary	Listar vacunas del catálogo por especie
//	@Tags		vaccines-catalog
//	@Produce	json
//	@Security	CookieAuth
//	@Param		species	path		string					true	"Especie (dog, cat, bird, rabbit, fish, other)"
//	@Success	200	{object}	map[string]interface{}	"data: []VaccineCatalog"
//	@Failure	400	{object}	map[string]string		"especie inválida"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/vaccines-catalog/species/{species} [get]
func (h *Handler) GetVaccinesBySpecies(c *gin.Context) {
	species := c.Param("species")
	if species == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid species"})
		return
	}

	validSpecies := map[string]bool{
		"dog":    true,
		"cat":    true,
		"bird":   true,
		"rabbit": true,
		"fish":   true,
		"other":  true,
	}
	if !validSpecies[species] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid species"})
		return
	}

	vaccines, err := h.repo.GetBySpecies(c.Request.Context(), species)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch vaccines by species"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": vaccines})
}

// CreateVaccineCatalog handles POST /api/v1/vaccines-catalog
// Only accessible by system users.
//
//	@Summary	Crear una nueva vaccine en el catálogo (solo usuario sistema)
//	@Tags		vaccines-catalog
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		vaccine	body		CreateVaccineCatalogPayload	true	"Datos de la vaccine"
//	@Success	201	{object}	map[string]interface{}	"data: VaccineCatalog"
//	@Failure	400	{object}	map[string]string		"error de validación"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	403	{object}	map[string]string		"prohibido"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/vaccines-catalog [post]
func (h *Handler) CreateVaccineCatalog(c *gin.Context) {
	if !requireSystemUser(c) {
		return
	}

	var payload CreateVaccineCatalogPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	created, err := h.repo.Create(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create vaccine catalog"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": created})
}

// UpdateVaccineCatalog handles PUT /api/v1/vaccines-catalog/:id
// Only accessible by system users.
//
//	@Summary	Actualizar una vaccine del catálogo (solo usuario sistema)
//	@Tags		vaccines-catalog
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id		path		string						true	"ID de la vaccine"
//	@Param		vaccine	body		UpdateVaccineCatalogPayload	true	"Datos de la vaccine"
//	@Success	200	{object}	map[string]interface{}	"data: VaccineCatalog"
//	@Failure	400	{object}	map[string]string			"id inválido o error de validación"
//	@Failure	401	{object}	map[string]string			"autenticación requerida"
//	@Failure	403	{object}	map[string]string			"prohibido"
//	@Failure	404	{object}	map[string]string			"vaccine no encontrada"
//	@Failure	500	{object}	map[string]string			"mensaje de error"
//	@Router		/api/v1/vaccines-catalog/{id} [put]
func (h *Handler) UpdateVaccineCatalog(c *gin.Context) {
	if !requireSystemUser(c) {
		return
	}

	id, ok := parseID(c)
	if !ok {
		return
	}

	var payload UpdateVaccineCatalogPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	updated, err := h.repo.Update(c.Request.Context(), id, payload)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "vaccine catalog not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update vaccine catalog"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

// DeleteVaccineCatalog handles DELETE /api/v1/vaccines-catalog/:id
// Only accessible by system users.
//
//	@Summary	Eliminar una vaccine del catálogo (solo usuario sistema)
//	@Tags		vaccines-catalog
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string				true	"ID de la vaccine"
//	@Success	200	{object}	map[string]string	"message: vaccine catalog deleted"
//	@Failure	400	{object}	map[string]string	"id inválido"
//	@Failure	401	{object}	map[string]string	"autenticación requerida"
//	@Failure	403	{object}	map[string]string	"prohibido"
//	@Failure	404	{object}	map[string]string	"vaccine no encontrada"
//	@Failure	500	{object}	map[string]string	"mensaje de error"
//	@Router		/api/v1/vaccines-catalog/{id} [delete]
func (h *Handler) DeleteVaccineCatalog(c *gin.Context) {
	if !requireSystemUser(c) {
		return
	}

	id, ok := parseID(c)
	if !ok {
		return
	}

	err := h.repo.Delete(c.Request.Context(), id)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "vaccine catalog not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete vaccine catalog"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "vaccine catalog deleted"})
}
