package user

import (
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/validation"
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
// UUIDs are passed directly as strings from the URL path.
func parseID(c *gin.Context) (string, bool) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return "", false
	}
	return id, true
}

// canModifyUser returns true if the caller (from JWT context) is the target user
// or is a system user.
func canModifyUser(c *gin.Context, targetID string) bool {
	callerID, _ := c.Get("userID")
	isSystemUser, _ := c.Get("isSystemUser")
	if uid, ok := callerID.(string); ok && uid == targetID {
		return true
	}
	return isSystemUser == true
}

// GetUsers handles GET /api/v1/users
// Only system users may list all users.
//
//	@Summary	List users (paginated, system user only)
//	@Tags		users
//	@Produce	json
//	@Security	CookieAuth
//	@Param		page		query		int						false	"Page number (default 1)"
//	@Param		per_page	query		int						false	"Items per page (default 10)"
//	@Success	200	{object}	map[string]interface{}	"data: []User, total: int, page: int, per_page: int, total_pages: int"
//	@Failure	401	{object}	map[string]string		"authentication required"
//	@Failure	403	{object}	map[string]string		"forbidden"
//	@Failure	500	{object}	map[string]string		"error message"
//	@Router		/api/v1/users [get]
func (h *Handler) GetUsers(c *gin.Context) {
	isSystemUser, _ := c.Get("isSystemUser")
	if isSystemUser != true {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	perPage, _ := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if page < 1 {
		page = 1
	}
	if perPage < 1 {
		perPage = 10
	}

	users, total, err := h.repo.GetPaginated(c.Request.Context(), page, perPage)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	c.JSON(http.StatusOK, gin.H{
		"data":        users,
		"total":       total,
		"page":        page,
		"per_page":    perPage,
		"total_pages": totalPages,
	})
}

// GetUser handles GET /api/v1/users/:id
// Only system users may retrieve a user by ID.
//
//	@Summary	Get a user by ID (system user only)
//	@Tags		users
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string					true	"User ID"
//	@Success	200	{object}	map[string]interface{}	"data: User"
//	@Failure	400	{object}	map[string]string		"invalid id"
//	@Failure	401	{object}	map[string]string		"authentication required"
//	@Failure	403	{object}	map[string]string		"forbidden"
//	@Failure	404	{object}	map[string]string		"user not found"
//	@Failure	500	{object}	map[string]string		"error message"
//	@Router		/api/v1/users/{id} [get]
func (h *Handler) GetUser(c *gin.Context) {
	isSystemUser, _ := c.Get("isSystemUser")
	if isSystemUser != true {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
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
// Only system users may create new users.
//
//	@Summary	Create a new user (system user only)
//	@Description	Only a user with is_system_user=true may call this endpoint.
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		user	body		CreateUserPayload		true	"User data"
//	@Success	201		{object}	map[string]interface{}	"data: User"
//	@Failure	400		{object}	map[string]string		"validation error"
//	@Failure	401		{object}	map[string]string		"authentication required"
//	@Failure	403		{object}	map[string]string		"forbidden"
//	@Failure	409		{object}	map[string]string		"email already in use"
//	@Failure	500		{object}	map[string]string		"error message"
//	@Router		/api/v1/users [post]
func (h *Handler) CreateUser(c *gin.Context) {
	callerIsSystemUser, _ := c.Get("isSystemUser")
	if callerIsSystemUser != true {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	var payload CreateUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	ctx := c.Request.Context()
	systemExists, err := h.repo.SystemUserExists(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}
	// First user ever -> mark as system user
	isSystemUser := !systemExists

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
// Only the user themselves or a system user may update a user.
//
//	@Summary	Update an existing user
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id		path		string					true	"User ID"
//	@Param		user	body		UpdateUserPayload		true	"User data"
//	@Success	200		{object}	map[string]interface{}	"data: User"
//	@Failure	400		{object}	map[string]string		"invalid id or validation error"
//	@Failure	401		{object}	map[string]string		"authentication required"
//	@Failure	403		{object}	map[string]string		"forbidden"
//	@Failure	404		{object}	map[string]string		"user not found"
//	@Failure	409		{object}	map[string]string		"email already in use"
//	@Failure	500		{object}	map[string]string		"error message"
//	@Router		/api/v1/users/{id} [put]
func (h *Handler) UpdateUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if !canModifyUser(c, id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}
	var payload UpdateUserPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
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
// Only the user themselves or a system user may delete a user.
// System users cannot be deleted.
//
//	@Summary	Delete a user
//	@Tags		users
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string				true	"User ID"
//	@Success	200	{object}	map[string]string	"message: user deleted"
//	@Failure	400	{object}	map[string]string	"invalid id"
//	@Failure	401	{object}	map[string]string	"authentication required"
//	@Failure	403	{object}	map[string]string	"forbidden"
//	@Failure	404	{object}	map[string]string	"user not found"
//	@Failure	500	{object}	map[string]string	"error message"
//	@Router		/api/v1/users/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	id, ok := parseID(c)
	if !ok {
		return
	}
	if !canModifyUser(c, id) {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return
	}

	// Prevent deletion of system users.
	target, err := h.repo.GetByID(c.Request.Context(), id)
	if errors.Is(err, ErrNotFound) {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}
	if target.IsSystemUser {
		c.JSON(http.StatusForbidden, gin.H{"error": "system users cannot be deleted"})
		return
	}

	err = h.repo.Delete(c.Request.Context(), id)
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
