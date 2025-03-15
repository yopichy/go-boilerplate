package models

import (
	"time"

	"gorm.io/gorm"
)

type OAuthClient struct {
	gorm.Model
	ClientID     string `json:"client_id" gorm:"unique"`
	ClientSecret string `json:"client_secret"`
	RedirectURI  string `json:"redirect_uri"`
}

type OAuthToken struct {
	gorm.Model
	AccessToken  string `json:"access_token" gorm:"unique"`
	TokenType    string `json:"token_type"`
	ExpiresAt    time.Time
	RefreshToken string `json:"refresh_token" gorm:"unique"`
	UserID       uint
	ClientID     string
	Scope        string `json:"scope"`
}

type OAuthCode struct {
	gorm.Model
	Code        string `gorm:"unique"`
	UserID      uint
	ClientID    string
	RedirectURI string
	ExpiresAt   time.Time
	Scope       string
}
