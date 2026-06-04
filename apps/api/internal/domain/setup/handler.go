package setup

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/domain/user"
	"github.com/my-pets/api/internal/validation"
)

// Handler holds setup handler dependencies.
type Handler struct {
	userRepo user.Repository
}

// NewHandler creates a new setup Handler.
func NewHandler(userRepo user.Repository) *Handler {
	return &Handler{userRepo: userRepo}
}

// SetupPayload is the request body for creating the first system user.
type SetupPayload struct {
	Name     string `json:"name"     binding:"required"`
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

// Status handles GET /api/v1/setup/status
// Retorna si el sistema ha sido inicializado (existe al menos un usuario).
//
//	@Summary	Verificar estado del setup
//	@Description	Verifica si el sistema ha sido inicializado (si existen usuarios).
//	@Tags		setup
//	@Produce	json
//	@Success	200	{object}	map[string]bool	"needs_setup: true si no existen usuarios"
//	@Router		/api/v1/setup/status [get]
func (h *Handler) Status(c *gin.Context) {
	hasUsers, err := h.userRepo.HasUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check setup status"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"needs_setup": !hasUsers})
}

// Create handles POST /api/v1/setup
// Crea el primer usuario sistema. Retorna 409 si ya existen usuarios.
//
//	@Summary	Crear el primer usuario sistema
//	@Description	Crea el primer usuario sistema. Solo funciona si no existen usuarios en el sistema.
//	@Tags		setup
//	@Accept		json
//	@Produce	json
//	@Param		user	body		SetupPayload			true	" Datos del primer usuario"
//	@Success	201		{object}	map[string]interface{}	"data: User"
//	@Failure	400		{object}	map[string]string		"error de validación"
//	@Failure	409		{object}	map[string]string		"sistema ya inicializado"
//	@Failure	500		{object}	map[string]string		"mensaje de error"
//	@Router		/api/v1/setup [post]
func (h *Handler) Create(c *gin.Context) {
	ctx := c.Request.Context()

	// Guard: reject if users already exist.
	hasUsers, err := h.userRepo.HasUsers(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to check setup status"})
		return
	}
	if hasUsers {
		c.JSON(http.StatusConflict, gin.H{"error": "system already initialized"})
		return
	}

	var payload SetupPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	created, err := h.userRepo.Create(ctx, user.CreateUserPayload{
		Name:     payload.Name,
		Email:    payload.Email,
		Password: payload.Password,
	}, true /* isSystemUser */)
	if errors.Is(err, user.ErrEmailTaken) {
		c.JSON(http.StatusConflict, gin.H{"error": "email already in use"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"data": created})
}
