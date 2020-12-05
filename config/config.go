package config

import (
	"context"
	"fmt"
	"github.com/coreos/go-oidc"
	"golang.org/x/oauth2"
	"os"
)

// AppConfig contains API/business configurations
type AppConfig struct {
	OIDCProvider                *oidc.Provider
	GoogleOauth                 oauth2.Config
	WebClientID                 string
	IOSClientID                 string
	AndroidClientID             string
	RevokeEndpoint              string
	GoogleServiceAccountKeyPath string
	TestMode                    bool
}

func (c *AppConfig) GetPlatformClientID(platform string) string {
	switch platform {
	case "web":
		return c.WebClientID
	case "ios":
		return c.IOSClientID
	case "android":
		return c.AndroidClientID
	}
	return c.WebClientID
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

	fmt.Println(getEnvAsBool("TEST_MODE", false))

	return &Config{
		Server: ServerConfig{
			Address: getEnv("ADDRESS", "localhost:4000"),
		},
		AppConfig: AppConfig{
			TestMode:                    getEnvAsBool("TEST_MODE", false),
			OIDCProvider:                provider,
			RevokeEndpoint:              getEnv("GOOGLE_OIDC_REVOKE_URL", "https://oauth2.googleapis.com/revoke"),
			WebClientID:                 os.Getenv("GOOGLE_OIDC_WEB_CLIENT_ID"),
			IOSClientID:                 os.Getenv("GOOGLE_OIDC_IOS_CLIENT_ID"),
			AndroidClientID:             os.Getenv("GOOGLE_OIDC_ANDROID_CLIENT_ID"),
			GoogleServiceAccountKeyPath: os.Getenv("GOOGLE_SERVICE_ACCOUNT_KEY"),
			GoogleOauth: oauth2.Config{
				ClientID:     os.Getenv("GOOGLE_OIDC_WEB_CLIENT_ID"),
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
