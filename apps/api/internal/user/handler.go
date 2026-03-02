package user

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// Handler holds the dependencies for user HTTP handlers.
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

// GetUsers handles GET /api/v1/users
func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users, "total": len(users)})
}

// GetUser handles GET /api/v1/users/:id
func (h *Handler) GetUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	u, err := h.repo.GetByID(c.Request.Context(), id)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": u})
}

// CreateUser handles POST /api/v1/users
// If no system user exists yet, the first created user is marked as the system user.
func (h *Handler) CreateUser(c *gin.Context) {
	var payload UserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx := c.Request.Context()
	isSystemUser, err := h.repo.SystemUserExists(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}
	// First user ever → mark as system user
	isSystemUser = !isSystemUser

	created, err := h.repo.Create(ctx, payload, isSystemUser)
	if errors.Is(err, ErrEmailTaken) {
		c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": created})
}

// UpdateUser handles PUT /api/v1/users/:id
func (h *Handler) UpdateUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	var payload UserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.repo.Update(c.Request.Context(), id, payload)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if errors.Is(err, ErrEmailTaken) {
		c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": updated})
}

// DeleteUser handles DELETE /api/v1/users/:id
func (h *Handler) DeleteUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	err := h.repo.Delete(c.Request.Context(), id)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}
