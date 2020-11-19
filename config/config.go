package config

import (
	"context"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"os"
)

// AppConfig contains API/business configurations
type AppConfig struct {
	OIDCProvider *oidc.Provider
	GoogleOauth  oauth2.Config
}

// ServerConfig contains server configurations (HTTP, etc)
type ServerConfig struct {
	Address string
}

// DatabaseConfig contains database configurations
type DatabaseConfig struct {
	URI                  string
	MigrationsDir        string
	MigrationsLogVerbose bool
}

type Config struct {
	Server    ServerConfig
	AppConfig AppConfig
	Database  DatabaseConfig
}

// NewConfig returns a Config object populated with values from environment variables or defaults
func NewConfig() *Config {
	provider, err := oidc.NewProvider(context.TODO(), "https://accounts.google.com")
	if err != nil {
		panic(err)
	}

	return &Config{
		Server: ServerConfig{
			Address: getEnv("ADDRESS", "localhost:4000"),
		},
		AppConfig: AppConfig{
			OIDCProvider: provider,
			GoogleOauth: oauth2.Config{
				ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
				ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
				Endpoint:     provider.Endpoint(),
				RedirectURL:  os.Getenv("GOOGLE_OAUTH_REDIRECT_URL"),
				Scopes: []string{
					"openid",
					"profile",
					"email",
				},
			},
		},
		Database: DatabaseConfig{
			URI:                  os.Getenv("DB_URI"),
			MigrationsDir:        getEnv("DB_MIGRATIONS_DIR", "file://migrations"),
			MigrationsLogVerbose: getEnvAsBool("DB_MIGRATIONS_VERBOSE", false),
		},
	}
}
