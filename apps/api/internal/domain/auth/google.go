package auth

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/my-pets/api/internal/config"
	"github.com/my-pets/api/internal/domain/user"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

const googleUserInfoURL = "https://www.googleapis.com/oauth2/v2/userinfo"

// oauthStateCookieName is used to prevent CSRF in the OAuth flow.
const oauthStateCookieName = "oauth_state"

// GoogleHandler handles the Google OAuth flow.
type GoogleHandler struct {
	cfg      *config.Config
	userRepo user.Repository
	oauth    *oauth2.Config
}

// NewGoogleHandler creates a GoogleHandler with a configured OAuth2 client.
func NewGoogleHandler(cfg *config.Config, userRepo user.Repository) *GoogleHandler {
	oauthCfg := &oauth2.Config{
		ClientID:     cfg.GoogleClientID,
		ClientSecret: cfg.GoogleClientSecret,
		RedirectURL:  cfg.AppURL + "/api/v1/auth/google/callback",
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
	return &GoogleHandler{cfg: cfg, userRepo: userRepo, oauth: oauthCfg}
}

// Initiate handles GET /api/v1/auth/google
// Redirects the user to Google's consent page.
//
//	@Summary	Iniciar autenticación con Google
//	@Description	Inicia el flujo de autenticación OAuth con Google. Redirige a la pantalla de consentimiento de Google.
//	@Tags		auth
//	@Success	302	"Redirect to Google"
//	@Router		/api/v1/auth/google [get]
func (h *GoogleHandler) Initiate(c *gin.Context) {
	state := generateState()
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     oauthStateCookieName,
		Value:    state,
		Path:     "/",
		MaxAge:   300, // 5 minutes
		HttpOnly: true,
		Secure:   h.cfg.GinMode == "release",
		SameSite: http.SameSiteLaxMode,
	})
	url := h.oauth.AuthCodeURL(state, oauth2.AccessTypeOnline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

// Callback handles GET /api/v1/auth/google/callback
// Exchanges the code for tokens, upserts the user, sets auth cookies, redirects to frontend.
//
//	@Summary	Callback de autenticación con Google
//	@Description	Procesa el callback de Google OAuth, intercambia el código por tokens y redirige al frontend.
//	@Tags		auth
//	@Success	302	"Redirect to frontend"
//	@Failure	400	{object}	map[string]string	"invalid state or code"
//	@Router		/api/v1/auth/google/callback [get]
func (h *GoogleHandler) Callback(c *gin.Context) {
	// Validate state to prevent CSRF.
	stateCookie, err := c.Cookie(oauthStateCookieName)
	if err != nil || stateCookie != c.Query("state") {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid oauth state"})
		return
	}
	// Clear state cookie.
	http.SetCookie(c.Writer, &http.Cookie{
		Name: oauthStateCookieName, Value: "", Path: "/", MaxAge: -1, HttpOnly: true,
	})

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing authorization code"})
		return
	}

	ctx := context.Background()
	oauthToken, err := h.oauth.Exchange(ctx, code)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to exchange code"})
		return
	}

	googleInfo, err := fetchGoogleUserInfo(ctx, h.oauth, oauthToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch Google profile"})
		return
	}

	// Determine if this will be the system user (first user ever).
	systemExists, err := h.userRepo.SystemUserExists(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}
	isSystemUser := !systemExists

	u, err := h.userRepo.UpsertByGoogleID(c.Request.Context(), user.GoogleUserInfo{
		GoogleID: googleInfo.ID,
		Email:    googleInfo.Email,
		Name:     googleInfo.Name,
	}, isSystemUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "login failed"})
		return
	}

	// Issue our own JWT cookies and redirect to the frontend.
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

	c.Redirect(http.StatusTemporaryRedirect, h.cfg.FrontendURL+"/")
}

// googleProfile is the minimal shape of Google's userinfo response.
type googleProfile struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func fetchGoogleUserInfo(ctx context.Context, cfg *oauth2.Config, token *oauth2.Token) (*googleProfile, error) {
	client := cfg.Client(ctx, token)
	resp, err := client.Get(googleUserInfoURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var profile googleProfile
	if err := json.NewDecoder(resp.Body).Decode(&profile); err != nil {
		return nil, err
	}
	return &profile, nil
}

// generateState creates a random hex state string for CSRF protection.
func generateState() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "fallback-state"
	}
	return fmt.Sprintf("%x", b)
}
