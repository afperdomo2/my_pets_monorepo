package user

import (
	"errors"
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
func parseID(c *gin.Context) (int, bool) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return 0, false
	}
	return id, true
}

// canModifyUser returns true if the caller (from JWT context) is the target user
// or is a system user.
func canModifyUser(c *gin.Context, targetID int) bool {
	callerID, _ := c.Get("userID")
	isSystemUser, _ := c.Get("isSystemUser")
	if uid, ok := callerID.(uint); ok && int(uid) == targetID {
		return true
	}
	return isSystemUser == true
}

// GetUsers handles GET /api/v1/users
// Only system users may list all users.
//
//	@Summary	List all users (system user only)
//	@Tags		users
//	@Produce	json
//	@Security	CookieAuth
//	@Success	200	{object}	map[string]interface{}	"data: []User, total: int"
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
	users, err := h.repo.GetAll(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": users, "total": len(users)})
}

// GetUser handles GET /api/v1/users/:id
// Only system users may retrieve a user by ID.
//
//	@Summary	Get a user by ID (system user only)
//	@Tags		users
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		int						true	"User ID"
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
// Only the user themselves or a system user may update a user.
//
//	@Summary	Update an existing user
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id		path		int						true	"User ID"
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
//	@Param		id	path		int					true	"User ID"
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

	// Prevent deletion of system users
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
