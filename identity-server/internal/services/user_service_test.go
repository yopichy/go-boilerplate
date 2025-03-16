package services

import (
	"identity-server/internal/models"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type TestUser struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
}

func loadTestUser(t *testing.T) TestUser {
	viper.Reset()
	viper.SetConfigName("user")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../models/testdata")

	if err := viper.ReadInConfig(); err != nil {
		t.Fatalf("Failed to load test user data: %v", err)
	}

	return TestUser{
		ID:           viper.GetString("users.test_user.id"),
		Username:     viper.GetString("users.test_user.username"),
		Email:        viper.GetString("users.test_user.email"),
		PasswordHash: viper.GetString("users.test_user.password_hash"),
	}
}

// MockUserRepository is a mock for user repository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) FindByID(id string) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func TestUserService_GetUser(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	testUser := loadTestUser(t)

	t.Run("successful get user", func(t *testing.T) {
		expectedUser := &models.User{
			ID:       testUser.ID,
			Username: testUser.Username,
			Email:    testUser.Email,
		}

		mockRepo.On("FindByID", testUser.ID).Return(expectedUser, nil)

		user, err := service.GetUser(testUser.ID)

		assert.NoError(t, err)
		assert.Equal(t, expectedUser, user)
		mockRepo.AssertExpectations(t)
	})

	t.Run("user not found", func(t *testing.T) {
		mockRepo.On("FindByID", "999").Return(nil, models.ErrUserNotFound)

		user, err := service.GetUser("999")

		assert.Error(t, err)
		assert.Nil(t, user)
		assert.Equal(t, models.ErrUserNotFound, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestUserService_ValidateCredentials(t *testing.T) {
	mockRepo := new(MockUserRepository)
	service := NewUserService(mockRepo)
	testUser := loadTestUser(t)

	t.Run("valid credentials", func(t *testing.T) {
		user := &models.User{
			Username: testUser.Username,
			Password: testUser.PasswordHash,
		}

		mockRepo.On("FindByUsername", testUser.Username).Return(user, nil)

		isValid := service.ValidateCredentials(testUser.Username, "correctpassword")
		assert.True(t, isValid)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		mockRepo.On("FindByUsername", testUser.Username).Return(nil, models.ErrUserNotFound)

		isValid := service.ValidateCredentials(testUser.Username, "wrongpassword")
		assert.False(t, isValid)
		mockRepo.AssertExpectations(t)
	})
}
