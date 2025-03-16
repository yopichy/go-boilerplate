package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"identity-server/internal/models"
	"identity-server/internal/services"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OAuthHandler struct {
	db          *gorm.DB
	userService *services.UserService
}

func NewOAuthHandler(db *gorm.DB) *OAuthHandler {
	return &OAuthHandler{
		db:          db,
		userService: services.NewUserService(db),
	}
}

// Authorize godoc
// @Summary OAuth2 authorization endpoint
// @Description Authorization endpoint for OAuth2 flows
// @Tags oauth2
// @Accept json
// @Produce json
// @Param client_id query string true "Client ID"
// @Param response_type query string true "Response Type (code or token)"
// @Param redirect_uri query string true "Redirect URI"
// @Param scope query string false "Requested scopes"
// @Param state query string false "State parameter for CSRF protection"
// @Success 302 {string} string "Redirect to client with code or token"
// @Failure 400 {object} map[string]string "Invalid request"
// @Router /oauth/authorize [get]
func (h *OAuthHandler) Authorize(c *gin.Context) {
	clientID := c.Query("client_id")
	responseType := c.Query("response_type")
	redirectURI := c.Query("redirect_uri")
	state := c.Query("state")
	scope := c.Query("scope")

	// Validate client
	var client models.OAuthClient
	if err := h.db.Where("client_id = ?", clientID).First(&client).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_client"})
		return
	}

	// For implicit flow (token response_type)
	if responseType == "token" {
		accessToken := generateRandomString(32)

		token := models.OAuthToken{
			AccessToken: accessToken,
			TokenType:   "Bearer",
			ExpiresAt:   time.Now().Add(time.Hour),
			ClientID:    clientID,
			Scope:       scope,
		}
		h.db.Create(&token)

		// Redirect with token
		redirectURL := fmt.Sprintf("%s#access_token=%s&token_type=Bearer&expires_in=3600&state=%s",
			redirectURI,
			accessToken,
			state,
		)
		c.Redirect(http.StatusTemporaryRedirect, redirectURL)
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported_response_type"})
}

// Token godoc
// @Summary OAuth2 token endpoint
// @Description Token endpoint for OAuth2 flows
// @Tags oauth2
// @Accept x-www-form-urlencoded
// @Produce json
// @Param grant_type formData string true "Grant Type (authorization_code, password, client_credentials)"
// @Param code formData string false "Authorization Code (required for authorization_code grant)"
// @Param client_id formData string true "Client ID"
// @Param client_secret formData string true "Client Secret"
// @Param redirect_uri formData string false "Redirect URI (required for authorization_code grant)"
// @Param username formData string false "Username (required for password grant)"
// @Param password formData string false "Password (required for password grant)"
// @Success 200 {object} models.OAuthToken
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 401 {object} map[string]string "Invalid client credentials"
// @Router /oauth/token [post]
func (h *OAuthHandler) Token(c *gin.Context) {
	grantType := c.PostForm("grant_type")
	clientID := c.PostForm("client_id")
	clientSecret := c.PostForm("client_secret")

	// Verify client
	var client models.OAuthClient
	if err := h.db.Where("client_id = ? AND client_secret = ?", clientID, clientSecret).First(&client).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_client"})
		return
	}

	switch grantType {
	case "password":
		username := c.PostForm("username")
		password := c.PostForm("password")

		// Use UserService for authentication
		user, err := h.userService.AuthenticateUser(username, password)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant", "description": err.Error()})
			return
		}

		// Generate tokens
		accessToken := generateRandomString(32)
		refreshToken := generateRandomString(32)

		token := models.OAuthToken{
			AccessToken:  accessToken,
			TokenType:    "Bearer",
			ExpiresAt:    time.Now().Add(time.Hour),
			RefreshToken: refreshToken,
			UserID:       user.ID,
			ClientID:     clientID,
			Scope:        user.Roles, // Use user roles as scope
		}
		h.db.Create(&token)

		c.JSON(http.StatusOK, gin.H{
			"access_token":  accessToken,
			"token_type":    "Bearer",
			"expires_in":    3600,
			"refresh_token": refreshToken,
			"scope":         user.Roles,
		})

	case "authorization_code":
		code := c.PostForm("code")
		redirectURI := c.PostForm("redirect_uri")

		// Verify authorization code
		var authCode models.OAuthCode
		if err := h.db.Where("code = ? AND client_id = ?", code, clientID).
			Where("expires_at > ?", time.Now()).
			First(&authCode).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant"})
			return
		}

		// Verify redirect URI matches the one used in authorization request
		if authCode.RedirectURI != redirectURI {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid_grant"})
			return
		}

		// Delete used code
		h.db.Delete(&authCode)

		// Generate tokens
		accessToken := generateRandomString(32)
		refreshToken := generateRandomString(32)

		token := models.OAuthToken{
			AccessToken:  accessToken,
			TokenType:    "Bearer",
			ExpiresAt:    time.Now().Add(time.Hour),
			RefreshToken: refreshToken,
			UserID:       authCode.UserID,
			ClientID:     clientID,
			Scope:        authCode.Scope,
		}
		h.db.Create(&token)

		c.JSON(http.StatusOK, gin.H{
			"access_token":  accessToken,
			"token_type":    "Bearer",
			"expires_in":    3600,
			"refresh_token": refreshToken,
			"scope":         token.Scope,
		})

	case "client_credentials":
		// Scope is optional for client credentials
		scope := c.PostForm("scope")
		if scope == "" {
			scope = "client"
		}

		// Generate tokens (no refresh token for client credentials)
		accessToken := generateRandomString(32)

		token := models.OAuthToken{
			AccessToken: accessToken,
			TokenType:   "Bearer",
			ExpiresAt:   time.Now().Add(time.Hour),
			ClientID:    clientID,
			Scope:       scope,
		}
		h.db.Create(&token)

		c.JSON(http.StatusOK, gin.H{
			"access_token": accessToken,
			"token_type":   "Bearer",
			"expires_in":   3600,
			"scope":        scope,
		})

	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported_grant_type"})
	}
}

// RegisterClient godoc
// @Summary Register OAuth client
// @Description Register a new OAuth client application
// @Tags oauth2
// @Accept json
// @Produce json
// @Param client body models.OAuthClient true "OAuth client details"
// @Success 201 {object} models.OAuthClient
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Internal server error"
// @Router /oauth/clients [post]
func (h *OAuthHandler) RegisterClient(c *gin.Context) {
	var client models.OAuthClient
	if err := c.ShouldBindJSON(&client); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if client.ClientID == "" || client.ClientSecret == "" || client.RedirectURI == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "client_id, client_secret and redirect_uri are required"})
		return
	}

	if err := h.db.Create(&client).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not register client"})
		return
	}

	c.JSON(http.StatusCreated, client)
}

func generateRandomString(n int) string {
	b := make([]byte, n)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)[:n]
}
