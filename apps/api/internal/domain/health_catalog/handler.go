package health_catalog

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/validation"
)

// Handler agrupa las dependencias para los handlers HTTP de la guía de salud.
type Handler struct {
	repo Repository
}

// NewHandler construye un Handler con el Repository dado.
func NewHandler(repo Repository) *Handler {
	return &Handler{repo: repo}
}

// parseID extrae y valida el parámetro de ruta :id.
func parseID(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return "", false
	}
	return id, true
}

// requireSystemUser retorna true si el caller es un usuario sistema.
// Si no lo es, responde 403 Forbidden y retorna false.
func requireSystemUser(c *gin.Context) bool {
	isSystemUser, _ := c.Get("isSystemUser")
	if isSystemUser != true {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return false
	}
	return true
}

// GetHealthCatalogsByCategory maneja GET /api/v1/health-catalogs/category/:category
// Accesible por cualquier usuario autenticado.
//
//	@Summary	Listar registros de la guía de salud por categoría (paginado)
//	@Tags		health-catalogs
//	@Produce	json
//	@Security	CookieAuth
//	@Param		category	path		string					true	"Categoría (vaccine, deworming, exam)"
//	@Param		page		query		int						false	"Número de página (por defecto 1)"
//	@Param		per_page	query		int						false	"Elementos por página (por defecto 10)"
//	@Param		species		query		string					false	"Filtrar por especie (dog, cat, bird, rabbit, fish, other)"
//	@Success	200	{object}	map[string]interface{}	"data: []HealthCatalog, total: int, page: int, per_page: int, total_pages: int"
//	@Failure	400	{object}	map[string]string	"categoría inválida"
//	@Failure	401	{object}	map[string]string	"autenticación requerida"
//	@Failure	500	{object}	map[string]string	"mensaje de error"
//	@Router		/api/v1/health-catalogs/category/{category} [get]
func (h *Handler) GetHealthCatalogsByCategory(c *gin.Context) {
	category := c.Param("category")
	
	validCategories := map[string]bool{
		"vaccine":   true,
		"deworming": true,
		"exam":      true,
	}
	if !validCategories[category] {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	species := c.Query("species")

	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	// Validar especie si se proporciona
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

	items, total, err := h.repo.GetPaginatedByCategory(c.Request.Context(), category, page, perPage, speciesFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch health catalog"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	c.JSON(http.StatusOK, gin.H{
		"data":        items,
		"total":       total,
		"page":        page,
		"per_page":    perPage,
		"total_pages": totalPages,
	})
}

// GetHealthCatalogByID maneja GET /api/v1/health-catalogs/:id
// Accesible por cualquier usuario autenticado.
//
//	@Summary	Obtener un registro de la guía de salud por ID
//	@Tags		health-catalogs
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string					true	"ID del registro"
//	@Success	200	{object}	map[string]interface{}	"data: HealthCatalog"
//	@Failure	400	{object}	map[string]string		"id inválido"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	404	{object}	map[string]string		"registro no encontrado"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/health-catalogs/{id} [get]
func (h *Handler) GetHealthCatalogByID(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	item, err := h.repo.GetByID(c.Request.Context(), id)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "health catalog not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch health catalog"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": item})
}

// GetHealthCatalogBySpecies maneja GET /api/v1/health-catalogs/species/:species
// Accesible por cualquier usuario autenticado.
//
//	@Summary	Listar registros de la guía de salud por especie
//	@Tags		health-catalogs
//	@Produce	json
//	@Security	CookieAuth
//	@Param		species	path		string					true	"Especie (dog, cat, bird, rabbit, fish, other)"
//	@Success	200	{object}	map[string]interface{}	"data: []HealthCatalog"
//	@Failure	400	{object}	map[string]string		"especie inválida"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/health-catalogs/species/{species} [get]
func (h *Handler) GetHealthCatalogBySpecies(c *gin.Context) {
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

	items, err := h.repo.GetBySpecies(c.Request.Context(), species)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch health catalog by species"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": items})
}

// CreateHealthCatalog maneja POST /api/v1/health-catalogs
// Solo accesible por usuarios sistema.
//
//	@Summary	Crear un nuevo registro en la guía de salud (solo usuario sistema)
//	@Tags		health-catalogs
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		item	body		CreateHealthCatalogPayload	true	"Datos del registro"
//	@Success	201	{object}	map[string]interface{}	"data: HealthCatalog"
//	@Failure	400	{object}	map[string]string		"error de validación"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	403	{object}	map[string]string		"prohibido"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/health-catalogs [post]
func (h *Handler) CreateHealthCatalog(c *gin.Context) {
	if !requireSystemUser(c) {
		return
	}

	var payload CreateHealthCatalogPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	created, err := h.repo.Create(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create health catalog"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": created})
}

// UpdateHealthCatalog maneja PUT /api/v1/health-catalogs/:id
// Solo accesible por usuarios sistema.
//
//	@Summary	Actualizar un registro de la guía de salud (solo usuario sistema)
//	@Tags		health-catalogs
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id		path		string						true	"ID del registro"
//	@Param		item	body		UpdateHealthCatalogPayload	true	"Datos del registro"
//	@Success	200	{object}	map[string]interface{}	"data: HealthCatalog"
//	@Failure	400	{object}	map[string]string			"id inválido o error de validación"
//	@Failure	401	{object}	map[string]string			"autenticación requerida"
//	@Failure	403	{object}	map[string]string			"prohibido"
//	@Failure	404	{object}	map[string]string			"registro no encontrado"
//	@Failure	500	{object}	map[string]string			"mensaje de error"
//	@Router		/api/v1/health-catalogs/{id} [put]
func (h *Handler) UpdateHealthCatalog(c *gin.Context) {
	if !requireSystemUser(c) {
		return
	}

	id, ok := parseID(c)
	if !ok {
		return
	}

	var payload UpdateHealthCatalogPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	updated, err := h.repo.Update(c.Request.Context(), id, payload)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "health catalog not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update health catalog"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

// DeleteHealthCatalog maneja DELETE /api/v1/health-catalogs/:id
// Solo accesible por usuarios sistema.
//
//	@Summary	Eliminar un registro de la guía de salud (solo usuario sistema)
//	@Tags		health-catalogs
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string				true	"ID del registro"
//	@Success	200	{object}	map[string]string	"message: health catalog deleted"
//	@Failure	400	{object}	map[string]string	"id inválido"
//	@Failure	401	{object}	map[string]string	"autenticación requerida"
//	@Failure	403	{object}	map[string]string	"prohibido"
//	@Failure	404	{object}	map[string]string	"registro no encontrado"
//	@Failure	500	{object}	map[string]string	"mensaje de error"
//	@Router		/api/v1/health-catalogs/{id} [delete]
func (h *Handler) DeleteHealthCatalog(c *gin.Context) {
	if !requireSystemUser(c) {
		return
	}

	id, ok := parseID(c)
	if !ok {
		return
	}

	err := h.repo.Delete(c.Request.Context(), id)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "health catalog not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete health catalog"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "health catalog deleted"})
}
