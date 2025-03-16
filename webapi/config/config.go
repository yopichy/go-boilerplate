package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Auth     AuthConfig
	OAuth2   OAuth2Config
	Logging  LoggingConfig
}

type ServerConfig struct {
	Port int
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	Database string
	Username string
	Password string
}

type AuthConfig struct {
	JWTSecret string
}

type OAuth2Config struct {
	AuthServerURL string
	ClientID      string
	ClientSecret  string
	RedirectURL   string
}

type LoggingConfig struct {
	FilePath string
	Level    string
	Format   string
}

func Load(paths ...string) (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	if len(paths) > 0 {
		for _, path := range paths {
			viper.AddConfigPath(path)
		}
	} else {
		viper.AddConfigPath(".")
	}

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
