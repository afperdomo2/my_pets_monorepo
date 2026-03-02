package pet

import (
	"errors"
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
func parseID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return id, true
}

// GetPets handles GET /api/v1/pets
//
//	@Summary	List all pets
//	@Tags		pets
//	@Produce	json
//	@Security	CookieAuth
//	@Success	200	{object}	map[string]interface{}	"data: []Pet, total: int"
//	@Failure	401	{object}	map[string]string		"authentication required"
//	@Failure	500	{object}	map[string]string		"error message"
//	@Router		/api/v1/pets [get]
func (h *Handler) GetPets(c *gin.Context) {
	pets, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch pets"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": pets, "total": len(pets)})
}

// GetPet handles GET /api/v1/pets/:id
//
//	@Summary	Get a pet by ID
//	@Tags		pets
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		int						true	"Pet ID"
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
	p, err := h.repo.GetByID(c.Request.Context(), id)
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
//	@Summary	Create a new pet
//	@Tags		pets
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		pet	body		PetPayload				true	"Pet data"
//	@Success	201	{object}	map[string]interface{}	"data: Pet"
//	@Failure	400	{object}	map[string]string		"validation error"
//	@Failure	401	{object}	map[string]string		"authentication required"
//	@Failure	500	{object}	map[string]string		"error message"
//	@Router		/api/v1/pets [post]
func (h *Handler) CreatePet(c *gin.Context) {
	var payload PetPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}
	created, err := h.repo.Create(c.Request.Context(), payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create pet"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": created})
}

// UpdatePet handles PUT /api/v1/pets/:id
//
//	@Summary	Update an existing pet
//	@Tags		pets
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		int						true	"Pet ID"
//	@Param		pet	body		PetPayload				true	"Pet data"
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
	var payload PetPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}
	updated, err := h.repo.Update(c.Request.Context(), id, payload)
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
//	@Summary	Delete a pet
//	@Tags		pets
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		int					true	"Pet ID"
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
	err := h.repo.Delete(c.Request.Context(), id)
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
