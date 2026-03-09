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
//	@Summary	List pets for the authenticated user (paginated)
//	@Tags		pets
//	@Produce	json
//	@Security	CookieAuth
//	@Param		page		query		int						false	"Page number (default 1)"
//	@Param		per_page	query		int						false	"Items per page (default 10)"
//	@Success	200	{object}	map[string]interface{}	"data: []Pet, total: int, page: int, per_page: int, total_pages: int"
//	@Failure	401	{object}	map[string]string		"authentication required"
//	@Failure	500	{object}	map[string]string		"error message"
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
//	@Summary	Get a pet by ID (must belong to the authenticated user)
//	@Tags		pets
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string					true	"Pet ID"
//	@Success	200	{object}	map[string]interface{}	"data: Pet"
//	@Failure	400	{object}	map[string]string		"invalid id"
//	@Failure	401	{object}	map[string]string		"authentication required"
//	@Failure	404	{object}	map[string]string		"pet not found"
//	@Failure	500	{object}	map[string]string		"error message"
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
//	@Summary	Create a new pet for the authenticated user
//	@Tags		pets
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		pet	body		CreatePetPayload		true	"Pet data"
//	@Success	201	{object}	map[string]interface{}	"data: Pet"
//	@Failure	400	{object}	map[string]string		"validation error"
//	@Failure	401	{object}	map[string]string		"authentication required"
//	@Failure	500	{object}	map[string]string		"error message"
//	@Router		/api/v1/pets [post]
func (h *Handler) CreatePet(c *gin.Context) {
	var payload CreatePetPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
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
//	@Summary	Update a pet (must belong to the authenticated user)
//	@Tags		pets
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string					true	"Pet ID"
//	@Param		pet	body		UpdatePetPayload		true	"Pet data"
//	@Success	200	{object}	map[string]interface{}	"data: Pet"
//	@Failure	400	{object}	map[string]string		"invalid id or validation error"
//	@Failure	401	{object}	map[string]string		"authentication required"
//	@Failure	404	{object}	map[string]string		"pet not found"
//	@Failure	500	{object}	map[string]string		"error message"
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
//	@Summary	Delete a pet (must belong to the authenticated user)
//	@Tags		pets
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string				true	"Pet ID"
//	@Success	200	{object}	map[string]string	"message: pet deleted"
//	@Failure	400	{object}	map[string]string	"invalid id"
//	@Failure	401	{object}	map[string]string	"authentication required"
//	@Failure	404	{object}	map[string]string	"pet not found"
//	@Failure	500	{object}	map[string]string	"error message"
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

// GetLifeStage handles GET /api/v1/pets/life-stage
//
//	@Summary	Calculate life stage for a dog or cat based on weight
//	@Tags		pets
//	@Produce	json
//	@Security	CookieAuth
//	@Param		species			query		string	true	"Pet species (dog or cat)"
//	@Param		weight_grams	query		int		true	"Weight in grams"
//	@Success	200	{object}	map[string]string	"life_stage: string"
//	@Failure	400	{object}	map[string]string	"missing or invalid parameters"
//	@Failure	401	{object}	map[string]string	"authentication required"
//	@Router		/api/v1/pets/life-stage [get]
func (h *Handler) GetLifeStage(c *gin.Context) {
	species := c.Query("species")
	if species == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "species is required"})
		return
	}

	weightStr := c.Query("weight_grams")
	if weightStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "weight_grams is required"})
		return
	}

	weightGrams, err := strconv.Atoi(weightStr)
	if err != nil || weightGrams < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "weight_grams must be a positive integer"})
		return
	}

	stage := CalculateLifeStage(species, weightGrams)
	c.JSON(http.StatusOK, gin.H{"life_stage": stage})
}
