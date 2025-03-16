package config

import (
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name    string
		setup   func()
		want    *Config
		wantErr bool
	}{
		{
			name: "successful config load",
			setup: func() {
				viper.Reset()
				viper.SetConfigName("config")
				viper.SetConfigType("yaml")
				viper.AddConfigPath("./testdata")
			},
			want: &Config{
				Server: ServerConfig{
					Port: 8001,
				},
				Database: DatabaseConfig{
					Driver:   "postgres",
					Host:     "localhost",
					Port:     5432,
					Database: "testdb",
					Username: "testuser",
					Password: "testpass",
				},
				Auth: AuthConfig{
					JWTSecret: "test-secret-key",
				},
				OAuth2: OAuth2Config{
					AuthServerURL: "http://localhost:8000",
					ClientID:      "test-client",
					ClientSecret:  "test-secret",
					RedirectURL:   "http://localhost:8001/callback",
				},
				Logging: LoggingConfig{
					FilePath: "./testdata/test.log",
					Level:    "info",
					Format:   "json",
				},
			},
			wantErr: false,
		},
		{
			name: "config file not found",
			setup: func() {
				viper.Reset()
				viper.SetConfigName("nonexistent")
				viper.SetConfigType("yaml")
				viper.AddConfigPath("./nonexistent")
			},
			want:    nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setup()

			got, err := Load()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, got)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, got)
			assert.Equal(t, tt.want, got)
		})
	}
}
