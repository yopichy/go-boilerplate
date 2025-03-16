package handlers

import (
	"encoding/json"
	"identity-server/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockOAuthService struct {
	mock.Mock
}

func (m *MockOAuthService) ValidateClient(clientID, clientSecret string) bool {
	args := m.Called(clientID, clientSecret)
	return args.Bool(0)
}

func (m *MockOAuthService) GenerateAccessToken(userID string) (*models.Token, error) {
	args := m.Called(userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Token), args.Error(1)
}

func TestOAuthHandler_TokenEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockOAuthService)
	handler := NewOAuthHandler(mockService)

	t.Run("successful token generation", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		// Mock form data
		c.Request = httptest.NewRequest("POST", "/oauth/token", nil)
		c.Request.ParseForm()
		c.Request.Form.Set("grant_type", "client_credentials")
		c.Request.Form.Set("client_id", "test-client")
		c.Request.Form.Set("client_secret", "test-secret")

		expectedToken := &models.Token{
			AccessToken: "test-token",
			TokenType:   "Bearer",
			ExpiresIn:   3600,
		}

		mockService.On("ValidateClient", "test-client", "test-secret").Return(true)
		mockService.On("GenerateAccessToken", "test-client").Return(expectedToken, nil)

		handler.TokenEndpoint(c)

		assert.Equal(t, http.StatusOK, w.Code)

		var response models.Token
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, expectedToken.AccessToken, response.AccessToken)
		mockService.AssertExpectations(t)
	})

	t.Run("invalid client credentials", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/oauth/token", nil)
		c.Request.ParseForm()
		c.Request.Form.Set("grant_type", "client_credentials")
		c.Request.Form.Set("client_id", "invalid-client")
		c.Request.Form.Set("client_secret", "invalid-secret")

		mockService.On("ValidateClient", "invalid-client", "invalid-secret").Return(false)

		handler.TokenEndpoint(c)

		assert.Equal(t, http.StatusUnauthorized, w.Code)
		mockService.AssertExpectations(t)
	})
}
