package pet

import (
	"errors"
	"math"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/domain/user"
	"github.com/my-pets/api/internal/validation"
)

// Handler holds the dependencies for pet HTTP handlers.
type Handler struct {
	repo     Repository
	userRepo user.Repository
}

// NewHandler constructs a Handler with the given Repository.
func NewHandler(repo Repository, userRepo user.Repository) *Handler {
	return &Handler{repo: repo, userRepo: userRepo}
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
//	@Description	Obtiene todas las mascotas del usuario autenticado con paginación.
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
//	@Description	Obtiene una mascota específica por su ID. Debe pertenecer al usuario autenticado.
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
//	@Description	Registra una nueva mascota para el usuario autenticado. Calcula la etapa de vida automáticamente según especie y edad.
//	@Tags		pets
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		pet	body		CreatePetPayload		true	"Datos de la mascota"
//	@Success	201	{object}	map[string]interface{}	"data: Pet"
//	@Failure	400	{object}	map[string]string		"error de validación"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	403	{object}	map[string]string		"límite de mascotas alcanzado"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/pets [post]
func (h *Handler) CreatePet(c *gin.Context) {
	var payload CreatePetPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	// Conditional size validation: required for dogs, forbidden for other species.
	if payload.Species == "dog" {
		if payload.Size == nil || *payload.Size == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "el tamaño es obligatorio para perros"})
			return
		}
		if !IsValidSizeCategory(*payload.Size) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tamaño no válido"})
			return
		}
	} else {
		payload.Size = nil
	}

	uid := ownerID(c)

	currentUser, err := h.userRepo.GetByID(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	count, err := h.repo.CountByOwner(c.Request.Context(), uid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to count pets"})
		return
	}

	if count >= int64(currentUser.PetLimit) {
		c.JSON(http.StatusForbidden, gin.H{"error": "pet limit reached"})
		return
	}

	// Calculate life stage automatically
	if payload.Species == "dog" {
		// Dogs: life stage based on age + size (birth date already validated above)
		birthDate, err := time.Parse("2006-01-02", payload.BirthDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "formato de fecha inválido"})
			return
		}
		stage := CalculateDogLifeStage(birthDate, SizeCategory(*payload.Size))
		payload.LifeStage = &stage
	} else if payload.Species == "cat" || payload.Species == "rabbit" {
		// Cats and rabbits: life stage based on age (birth date already present in payload)
		birthDate, err := time.Parse("2006-01-02", payload.BirthDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "formato de fecha inválido"})
			return
		}
		var stage string
		if payload.Species == "cat" {
			stage = CalculateCatLifeStage(birthDate)
		} else {
			stage = CalculateRabbitLifeStage(birthDate)
		}
		payload.LifeStage = &stage
	} else if payload.Species == "bird" {
		// Birds: life stage based on age
		birthDate, err := time.Parse("2006-01-02", payload.BirthDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "formato de fecha inválido"})
			return
		}
		stage := CalculateBirdLifeStage(birthDate)
		payload.LifeStage = &stage
	}

	created, err := h.repo.Create(c.Request.Context(), uid, payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create pet"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": created})
}

// UpdatePet handles PUT /api/v1/pets/:id
//
//	@Summary	Actualizar una mascota (debe pertenecer al usuario autenticado)
//	@Description	Actualiza los datos de una mascota existente. Recalcula la etapa de vida si cambia la fecha de nacimiento.
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

	// Get existing pet to get species (which cannot be changed)
	existingPet, err := h.repo.GetByID(c.Request.Context(), id, ownerID(c))
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "pet not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch pet"})
		return
	}

	// Conditional size validation: required for dogs, forbidden for other species.
	if existingPet.Species == "dog" {
		if payload.Size == nil || *payload.Size == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "el tamaño es obligatorio para perros"})
			return
		}
		if !IsValidSizeCategory(*payload.Size) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "tamaño no válido"})
			return
		}
	} else {
		payload.Size = nil
	}

	// Recalculate life stage when birth date or size may have changed.
	if existingPet.Species == "dog" || existingPet.Species == "cat" || existingPet.Species == "rabbit" || existingPet.Species == "bird" {
		birthDate, err := time.Parse("2006-01-02", payload.BirthDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "formato de fecha inválido"})
			return
		}
		var stage string
		switch existingPet.Species {
		case "dog":
			stage = CalculateDogLifeStage(birthDate, SizeCategory(*payload.Size))
		case "cat":
			stage = CalculateCatLifeStage(birthDate)
		case "rabbit":
			stage = CalculateRabbitLifeStage(birthDate)
		case "bird":
			stage = CalculateBirdLifeStage(birthDate)
		}
		payload.LifeStage = &stage
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
//	@Description	Elimina una mascota y todos sus datos asociados (registros de salud, exámenes, etc.).
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
