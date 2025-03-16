package handlers

import (
	"encoding/json"
	"identity-server/config"
	"identity-server/internal/models"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

func setupTestConfig(t *testing.T) *config.Config {
	cfg, err := config.Load("../../../config/testdata")
	if err != nil {
		t.Fatalf("Failed to load test config: %v", err)
	}
	return cfg
}

type MockDB struct {
	mock.Mock
}

func (m *MockDB) Where(query interface{}, args ...interface{}) *gorm.DB {
	m.Called(query, args)
	return &gorm.DB{}
}

func (m *MockDB) First(dest interface{}, conds ...interface{}) *gorm.DB {
	m.Called(dest, conds)
	return &gorm.DB{}
}

func (m *MockDB) Create(value interface{}) *gorm.DB {
	m.Called(value)
	return &gorm.DB{}
}

func TestOAuthHandler_TokenEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	cfg := setupTestConfig(t)
	mockDB := new(MockDB)

	handler := &OAuthHandler{
		db: mockDB,
	}

	t.Run("successful token generation", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/oauth/token", nil)
		c.Request.ParseForm()
		c.Request.Form.Set("grant_type", "client_credentials")
		c.Request.Form.Set("client_id", cfg.OAuth2.ClientID)
		c.Request.Form.Set("client_secret", cfg.OAuth2.ClientSecret)

		// Mock client validation
		mockDB.On("Where", "client_id = ? AND client_secret = ?", cfg.OAuth2.ClientID, cfg.OAuth2.ClientSecret).Return(mockDB)
		mockDB.On("First", &models.OAuthClient{}, []interface{}{}).Return(nil)
		mockDB.On("Create", mock.AnythingOfType("*models.OAuthToken")).Return(nil)

		handler.Token(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		err := json.NewDecoder(w.Body).Decode(&response)
		assert.NoError(t, err)
		assert.Contains(t, response, "access_token")
		assert.Equal(t, "Bearer", response["token_type"])
		mockDB.AssertExpectations(t)
	})

	t.Run("invalid client credentials", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)

		c.Request = httptest.NewRequest("POST", "/oauth/token", nil)
		c.Request.ParseForm()
		c.Request.Form.Set("grant_type", "client_credentials")
		c.Request.Form.Set("client_id", "invalid-client")
		c.Request.Form.Set("client_secret", "invalid-secret")

		// Mock failed client validation
		mockDB.On("Where", "client_id = ? AND client_secret = ?", "invalid-client", "invalid-secret").Return(mockDB)
		mockDB.On("First", &models.OAuthClient{}, []interface{}{}).Return(&gorm.DB{Error: gorm.ErrRecordNotFound})

		handler.Token(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		mockDB.AssertExpectations(t)
	})
}
