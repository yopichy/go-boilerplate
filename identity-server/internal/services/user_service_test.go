package services

import (
	"identity-server/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

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

	t.Run("successful get user", func(t *testing.T) {
		expectedUser := &models.User{
			ID:       "1",
			Username: "testuser",
			Email:    "test@example.com",
		}

		mockRepo.On("FindByID", "1").Return(expectedUser, nil)

		user, err := service.GetUser("1")

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

	t.Run("valid credentials", func(t *testing.T) {
		user := &models.User{
			Username: "testuser",
			Password: "$2a$10$somehashedpassword", // Assume this is a properly hashed password
		}

		mockRepo.On("FindByUsername", "testuser").Return(user, nil)

		isValid := service.ValidateCredentials("testuser", "correctpassword")
		assert.True(t, isValid)
		mockRepo.AssertExpectations(t)
	})

	t.Run("invalid credentials", func(t *testing.T) {
		mockRepo.On("FindByUsername", "testuser").Return(nil, models.ErrUserNotFound)

		isValid := service.ValidateCredentials("testuser", "wrongpassword")
		assert.False(t, isValid)
		mockRepo.AssertExpectations(t)
	})
}
