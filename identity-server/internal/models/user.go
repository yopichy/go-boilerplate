package models

import (
	"time"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username          string    `json:"username" gorm:"unique"`
	Email             string    `json:"email" gorm:"unique"`
	Password          string    `json:"-"` // "-" to never send password in JSON
	Roles             string    `json:"roles" gorm:"default:'user'"`
	IsActive          bool      `json:"is_active" gorm:"default:true"`
	IsEmailVerified   bool      `json:"is_email_verified" gorm:"default:false"`
	FailedLoginAttempts int     `json:"failed_login_attempts" gorm:"default:0"`
	LastLoginAt       *time.Time `json:"last_login_at"`
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
