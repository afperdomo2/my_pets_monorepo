package user

import (
	"context"
	"errors"
	"math"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/models"
	"github.com/my-pets/api/internal/validation"
)

type PetCountFn func(ctx context.Context, ownerID string) (int64, error)

type UserWithPetCount struct {
	models.User
	PetCount int64 `json:"pet_count"`
}

// Handler holds the dependencies for user HTTP handlers.
type Handler struct {
	repo       Repository
	petCountFn PetCountFn
}

// NewHandler constructs a Handler with the given Repository.
func NewHandler(repo Repository, petCountFn PetCountFn) *Handler {
	return &Handler{repo: repo, petCountFn: petCountFn}
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
// Solo usuarios sistema pueden listar todos los usuarios.
//
//	@Summary	Listar usuarios (paginado, solo usuario sistema)
//	@Description	Lista todos los usuarios del sistema con paginación. Solo accesible por usuarios sistema.
//	@Tags		users
//	@Produce	json
//	@Security	CookieAuth
//	@Param		page		query		int						false	"Número de página (por defecto 1)"
//	@Param		per_page	query		int						false	"Elementos por página (por defecto 10)"
//	@Success	200	{object}	map[string]interface{}	"data: []User, total: int, page: int, per_page: int, total_pages: int"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	403	{object}	map[string]string		"prohibido"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
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

	usersWithCount := make([]UserWithPetCount, len(users))
	for i, u := range users {
		petCount, _ := h.petCountFn(c.Request.Context(), u.ID)
		usersWithCount[i] = UserWithPetCount{
			User:     u,
			PetCount: petCount,
		}
	}

	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	c.JSON(http.StatusOK, gin.H{
		"data":        usersWithCount,
		"total":       total,
		"page":        page,
		"per_page":    perPage,
		"total_pages": totalPages,
	})
}

// GetUser handles GET /api/v1/users/:id
// Solo usuarios sistema pueden obtener un usuario por ID.
//
//	@Summary	Obtener un usuario por ID (solo usuario sistema)
//	@Description	Obtiene un usuario por su ID. Solo accesible por usuarios sistema.
//	@Tags		users
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string					true	"ID del usuario"
//	@Success	200	{object}	map[string]interface{}	"data: User"
//	@Failure	400	{object}	map[string]string		"id inválido"
//	@Failure	401	{object}	map[string]string		"autenticación requerida"
//	@Failure	403	{object}	map[string]string		"prohibido"
//	@Failure	404	{object}	map[string]string		"usuario no encontrado"
//	@Failure	500	{object}	map[string]string		"mensaje de error"
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

	petCount, _ := h.petCountFn(c.Request.Context(), u.ID)
	userWithCount := UserWithPetCount{
		User:     u,
		PetCount: petCount,
	}
	c.JSON(http.StatusOK, gin.H{"data": userWithCount})
}

// CreateUser handles POST /api/v1/users
// Solo usuarios sistema pueden crear nuevos usuarios.
//
//	@Summary	Crear un nuevo usuario (solo usuario sistema)
//	@Description	Solo un usuario con is_system_user=true puede llamar a este endpoint.
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		user	body		CreateUserPayload		true	"Datos del usuario"
//	@Success	201		{object}	map[string]interface{}	"data: User"
//	@Failure	400		{object}	map[string]string		"error de validación"
//	@Failure	401		{object}	map[string]string		"autenticación requerida"
//	@Failure	403		{object}	map[string]string		"prohibido"
//	@Failure	409		{object}	map[string]string		"email ya en uso"
//	@Failure	500		{object}	map[string]string		"mensaje de error"
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
// Solo el propio usuario o un usuario sistema pueden actualizar un usuario.
//
//	@Summary	Actualizar un usuario existente
//	@Description	Actualiza los datos de un usuario existente. El propio usuario o un usuario sistema pueden hacerlo.
//	@Tags		users
//	@Accept		json
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id		path		string					true	"ID del usuario"
//	@Param		user	body		UpdateUserPayload		true	"Datos del usuario"
//	@Success	200		{object}	map[string]interface{}	"data: User"
//	@Failure	400		{object}	map[string]string		"id inválido o error de validación"
//	@Failure	401		{object}	map[string]string		"autenticación requerida"
//	@Failure	403		{object}	map[string]string		"prohibido"
//	@Failure	404		{object}	map[string]string		"usuario no encontrado"
//	@Failure	409		{object}	map[string]string		"email ya en uso"
//	@Failure	500		{object}	map[string]string		"mensaje de error"
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

	callerIsSystemUser, _ := c.Get("isSystemUser")
	if callerIsSystemUser != true {
		payload.PetLimit = nil
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
// Solo el propio usuario o un usuario sistema pueden eliminar un usuario.
// Usuarios sistema no pueden ser eliminados.
//
//	@Summary	Eliminar un usuario
//	@Description	Elimina un usuario del sistema. Los usuarios sistema no pueden ser eliminados.
//	@Tags		users
//	@Produce	json
//	@Security	CookieAuth
//	@Param		id	path		string				true	"ID del usuario"
//	@Success	200	{object}	map[string]string	"message: user deleted"
//	@Failure	400	{object}	map[string]string	"id inválido"
//	@Failure	401	{object}	map[string]string	"autenticación requerida"
//	@Failure	403	{object}	map[string]string	"prohibido"
//	@Failure	404	{object}	map[string]string	"usuario no encontrado"
//	@Failure	500	{object}	map[string]string	"mensaje de error"
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
