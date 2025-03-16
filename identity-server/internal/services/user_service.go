package services

import (
	"errors"
	"identity-server/internal/models"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB
}

func NewUserService(db *gorm.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)
	user.IsActive = true
	user.IsEmailVerified = false
	user.FailedLoginAttempts = 0
	
	return s.db.Create(user).Error
}

func (s *UserService) AuthenticateUser(username, password string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}

	if !user.IsActive {
		return nil, errors.New("account is disabled")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		// Increment failed login attempts
		s.db.Model(&user).UpdateColumn("failed_login_attempts", gorm.Expr("failed_login_attempts + ?", 1))
		
		// If too many failed attempts, disable account
		if user.FailedLoginAttempts >= 5 {
			s.db.Model(&user).Update("is_active", false)
			return nil, errors.New("account locked due to too many failed attempts")
		}
		
		return nil, errors.New("invalid credentials")
	}

	// Reset failed attempts and update last login
	now := time.Now()
	s.db.Model(&user).Updates(map[string]interface{}{
		"failed_login_attempts": 0,
		"last_login_at": &now,
	})

	return &user, nil
}

func (s *UserService) VerifyEmail(userID uint) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).Update("is_email_verified", true).Error
}

func (s *UserService) UpdateUserRole(userID uint, role string) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).Update("roles", role).Error
}

func (s *UserService) DeactivateUser(userID uint) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).Update("is_active", false).Error
}

func (s *UserService) ReactivateUser(userID uint) error {
	return s.db.Model(&models.User{}).Where("id = ?", userID).Updates(map[string]interface{}{
		"is_active": true,
		"failed_login_attempts": 0,
	}).Error
}