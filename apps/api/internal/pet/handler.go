package pet

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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
func (h *Handler) GetPets(c *gin.Context) {
	pets, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch pets"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": pets, "total": len(pets)})
}

// GetPet handles GET /api/v1/pets/:id
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
func (h *Handler) CreatePet(c *gin.Context) {
	var payload PetPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
func (h *Handler) UpdatePet(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var payload PetPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
