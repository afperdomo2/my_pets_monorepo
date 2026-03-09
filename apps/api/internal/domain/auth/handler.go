package auth

import (
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/config"
	"github.com/my-pets/api/internal/domain/user"
	"github.com/my-pets/api/internal/models"
	"github.com/my-pets/api/internal/validation"
	"golang.org/x/crypto/bcrypt"
)

// Handler holds auth handler dependencies.
type Handler struct {
	cfg      *config.Config
	userRepo user.Repository
}

// NewHandler creates a new auth Handler.
func NewHandler(cfg *config.Config, userRepo user.Repository) *Handler {
	return &Handler{cfg: cfg, userRepo: userRepo}
}

// LoginPayload is the request body for local login.
type LoginPayload struct {
	Email    string `json:"email"    binding:"required,email"`
	Password string `json:"password" binding:"required,min=1"`
}

// setAuthCookies writes the access and refresh tokens as HttpOnly cookies.
func setAuthCookies(c *gin.Context, access, refresh string, secure bool) {
	sameSite := http.SameSiteLaxMode

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    access,
		Path:     "/",
		MaxAge:   int(accessTokenTTL.Seconds()),
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refresh,
		Path:     "/api/v1/auth/refresh",
		MaxAge:   int(refreshTokenTTL.Seconds()),
		HttpOnly: true,
		Secure:   secure,
		SameSite: sameSite,
	})
}

// clearAuthCookies removes both auth cookies by setting MaxAge to -1.
// The Path of each cookie must match exactly what was used when setting it.
func clearAuthCookies(c *gin.Context) {
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "access_token",
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
	})
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    "",
		Path:     "/api/v1/auth/refresh",
		MaxAge:   -1,
		HttpOnly: true,
	})
}

// Login handles POST /api/v1/auth/login
//
//	@Summary	Iniciar sesión con email y contraseña
//	@Tags		auth
//	@Accept		json
//	@Produce	json
//	@Param		credentials	body		LoginPayload			true	"Credenciales de acceso"
//	@Success	200			{object}	map[string]interface{}	"data: User"
//	@Failure	400			{object}	map[string]string		"error de validación"
//	@Failure	401			{object}	map[string]string		"credenciales inválidas"
//	@Router		/api/v1/auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var payload LoginPayload
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": validation.Translate(err)})
		return
	}

	ctx := c.Request.Context()

	// Fetch the stored hash — returns ErrNotFound or ErrWrongProvider if needed.
	hash, err := h.userRepo.GetPasswordByEmail(ctx, payload.Email)
	if errors.Is(err, user.ErrNotFound) || errors.Is(err, user.ErrWrongProvider) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(payload.Password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	u, err := h.userRepo.GetByEmail(ctx, payload.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	h.issueTokensAndRespond(c, u)
}

// Logout handles POST /api/v1/auth/logout
//
//	@Summary	Cerrar sesión — limpia las cookies de autenticación
//	@Tags		auth
//	@Produce	json
//	@Security	CookieAuth
//	@Success	200	{object}	map[string]string	"message: logged out"
//	@Failure	401	{object}	map[string]string	"autenticación requerida"
//	@Router		/api/v1/auth/logout [post]
func (h *Handler) Logout(c *gin.Context) {
	clearAuthCookies(c)
	c.JSON(http.StatusOK, gin.H{"message": "logged out"})
}

// Refresh handles POST /api/v1/auth/refresh
//
//	@Summary	Renovar token de acceso usando la cookie de refresh
//	@Tags		auth
//	@Produce	json
//	@Success	200	{object}	map[string]interface{}	"data: User"
//	@Failure	401	{object}	map[string]string		"token de refresh inválido o expirado"
//	@Router		/api/v1/auth/refresh [post]
func (h *Handler) Refresh(c *gin.Context) {
	cookie, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "no refresh token"})
		return
	}

	claims, err := ValidateClaims(h.cfg, cookie)
	if err != nil {
		clearAuthCookies(c)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired refresh token"})
		return
	}

	u, err := h.userRepo.GetByID(c.Request.Context(), claims.UserID)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}

	h.issueTokensAndRespond(c, u)
}

// Me handles GET /api/v1/auth/me
//
//	@Summary	Obtener el usuario actualmente autenticado
//	@Tags		auth
//	@Produce	json
//	@Security	CookieAuth
//	@Success	200	{object}	map[string]interface{}	"data: User"
//	@Failure	401	{object}	map[string]string		"no autorizado"
//	@Router		/api/v1/auth/me [get]
func (h *Handler) Me(c *gin.Context) {
	userID, _ := c.Get("userID")
	u, err := h.userRepo.GetByID(c.Request.Context(), userID.(string))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": u})
}

// issueTokensAndRespond generates both tokens, sets cookies, and returns the user.
func (h *Handler) issueTokensAndRespond(c *gin.Context, u models.User) {
	access, err := GenerateAccessToken(h.cfg, u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}
	refresh, err := GenerateRefreshToken(h.cfg, u)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	secure := h.cfg.GinMode == "release"
	setAuthCookies(c, access, refresh, secure)

	c.JSON(http.StatusOK, gin.H{
		"data":       u,
		"expires_in": int(time.Now().Add(accessTokenTTL).Unix()),
	})
}
