package handlers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"webapi/config"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	cfg config.OAuth2Config
}

func NewAuthHandler(cfg config.OAuth2Config) *AuthHandler {
	return &AuthHandler{cfg: cfg}
}

// Login godoc
// @Summary Initiate OAuth2 login flow
// @Description Redirects user to the OAuth2 authorization server
// @Tags auth
// @Accept json
// @Produce json
// @Success 307 {string} string "Redirect to authorization server"
// @Router /login [get]
func (h *AuthHandler) Login(c *gin.Context) {
	state := generateRandomState()
	c.SetCookie("oauth_state", state, 3600, "/", "", false, true)

	authURL := fmt.Sprintf("%s/oauth/authorize?response_type=code&client_id=%s&redirect_uri=%s&state=%s&scope=read",
		h.cfg.AuthServerURL,
		h.cfg.ClientID,
		url.QueryEscape(h.cfg.RedirectURL),
		state,
	)

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// Callback godoc
// @Summary OAuth2 callback handler
// @Description Handles the OAuth2 callback and exchanges code for token
// @Tags auth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code"
// @Param state query string true "State parameter for CSRF protection"
// @Success 307 {string} string "Redirect to protected resource"
// @Failure 400 {object} map[string]string "Invalid request"
// @Failure 500 {object} map[string]string "Server error"
// @Router /callback [get]
func (h *AuthHandler) Callback(c *gin.Context) {
	state := c.Query("state")
	storedState, _ := c.Cookie("oauth_state")
	if state != storedState {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state"})
		return
	}

	code := c.Query("code")
	token, err := h.exchangeCodeForToken(code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get token"})
		return
	}

	c.SetCookie("access_token", token.AccessToken, 3600, "/", "", false, true)
	c.Redirect(http.StatusTemporaryRedirect, "/api/weather")
}

func (h *AuthHandler) exchangeCodeForToken(code string) (*TokenResponse, error) {
	data := url.Values{}
	data.Set("grant_type", "authorization_code")
	data.Set("code", code)
	data.Set("client_id", h.cfg.ClientID)
	data.Set("client_secret", h.cfg.ClientSecret)
	data.Set("redirect_uri", h.cfg.RedirectURL)

	resp, err := http.PostForm(h.cfg.AuthServerURL+"/oauth/token", data)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var token TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return nil, err
	}

	return &token, nil
}

func generateRandomState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// TokenResponse represents the OAuth2 token response
type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}
