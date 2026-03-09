package pet

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/validation"
)

// Handler holds the dependencies for pet HTTP handlers.
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

// ownerID extracts the authenticated user's ID from the Gin context (set by JWT middleware).
func ownerID(c *gin.Context) string {
	return c.GetString("userID")
}

// GetPets handles GET /api/v1/pets
//
//	@Summary	Listar mascotas del usuario autenticado (paginado)
//	@Tags		pets
//	@Produce	json
//	@Security	CookieAuth
//	@Param		page		query		int						false	"Número de página (por defecto 1)"
//	@Param		per_page	query		int						false	"Elementos por página (por defecto 10)"
//	@Success	200	{object}	map[string]interface{}	"data: []Pet, total: int, page: int, per_page: int, total_pages: int"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/pets [get]
func (h *Handler) GetPets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	pets, total, err := h.repo.GetPaginated(c.Request.Context(), ownerID(c), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch pets"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	c.JSON(http.StatusOK, gin.H{
		"data":        pets,
		"total":       total,
		"page":        page,
		"per_page":    perPage,
		"total_pages": totalPages,
	})
}

// GetPet handles GET /api/v1/pets/:id
//
//	@Summary	Obtener una mascota por ID (debe pertenecer al usuario autenticado)
//	@Tags		pets
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string					true	"ID de la mascota"
//	@Success	200	{object}	map[string]interface{}	"data: Pet"
//	@Failure	400	{object}	map[string]string		"id inválido"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	404	{object}	map[string]string		"mascota no encontrada"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/pets/{id} [get]
func (h *Handler) GetPet(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	p, err := h.repo.GetByID(c.Request.Context(), id, ownerID(c))
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "pet not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch pet"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": p})
}

// CreatePet handles POST /api/v1/pets
//
//	@Summary	Crear una nueva mascota para el usuario autenticado
//	@Tags		pets
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		pet	body		CreatePetPayload		true	"Datos de la mascota"
//	@Success	201	{object}	map[string]interface{}	"data: Pet"
//	@Failure	400	{object}	map[string]string		"error de validación"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/pets [post]
func (h *Handler) CreatePet(c *gin.Context) {
	var payload CreatePetPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	// Calculate life stage automatically if applicable
	if payload.WeightGrams != nil && (payload.Species == "dog" || payload.Species == "cat") {
		stage := CalculateLifeStage(payload.Species, *payload.WeightGrams)
		payload.LifeStage = &stage
	}

	created, err := h.repo.Create(c.Request.Context(), ownerID(c), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create pet"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": created})
}

// UpdatePet handles PUT /api/v1/pets/:id
//
//	@Summary	Actualizar una mascota (debe pertenecer al usuario autenticado)
//	@Tags		pets
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string					true	"ID de la mascota"
//	@Param		pet	body		UpdatePetPayload		true	"Datos de la mascota"
//	@Success	200	{object}	map[string]interface{}	"data: Pet"
//	@Failure	400	{object}	map[string]string		"id inválido o error de validación"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	404	{object}	map[string]string		"mascota no encontrada"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/pets/{id} [put]
func (h *Handler) UpdatePet(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var payload UpdatePetPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}
	updated, err := h.repo.Update(c.Request.Context(), id, ownerID(c), payload)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "pet not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update pet"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

// DeletePet handles DELETE /api/v1/pets/:id
//
//	@Summary	Eliminar una mascota (debe pertenecer al usuario autenticado)
//	@Tags		pets
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string				true	"ID de la mascota"
//	@Success	200	{object}	map[string]string	"message: pet deleted"
//	@Failure	400	{object}	map[string]string	"id inválido"
//	@Failure	401	{object}	map[string]string	"autenticación requerida"
//	@Failure	404	{object}	map[string]string	"mascota no encontrada"
//	@Failure	500	{object}	map[string]string	"mensaje de error"
//	@Router		/api/v1/pets/{id} [delete]
func (h *Handler) DeletePet(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	err := h.repo.Delete(c.Request.Context(), id, ownerID(c))
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "pet not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete pet"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "pet deleted"})
}
