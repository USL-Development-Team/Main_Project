package config

import (
	"fmt"
	"os"
)

// Environment represents the deployment environment
type Environment string

const (
	Development Environment = "development"
	Staging     Environment = "staging"
	Production  Environment = "production"
)

// GetEnvironment determines the current environment
func GetEnvironment() Environment {
	env := os.Getenv("ENVIRONMENT")
	switch env {
	case "staging":
		return Staging
	case "production":
		return Production
	default:
		return Development
	}
}

// EnvironmentConfig provides environment-specific configuration
type EnvironmentConfig struct {
	Environment         Environment
	AppBaseURL          string
	RequireHTTPS        bool
	AllowedOrigins      []string
	LogLevel            string
	SessionCookieSecure bool
	EnableDebugLogging  bool
}

// GetEnvironmentConfig returns configuration based on environment
func GetEnvironmentConfig() EnvironmentConfig {
	env := GetEnvironment()

	switch env {
	case Production:
		appBaseURL := getEnv("APP_BASE_URL", "")
		if appBaseURL == "" {
			panic("APP_BASE_URL environment variable is required in production")
		}
		return EnvironmentConfig{
			Environment:         Production,
			AppBaseURL:          appBaseURL,
			RequireHTTPS:        true,
			AllowedOrigins:      []string{appBaseURL},
			LogLevel:            "info",
			SessionCookieSecure: true,
			EnableDebugLogging:  false,
		}
	case Staging:
		appBaseURL := getEnv("APP_BASE_URL", "")
		if appBaseURL == "" {
			panic("APP_BASE_URL environment variable is required in staging")
		}
		return EnvironmentConfig{
			Environment:         Staging,
			AppBaseURL:          appBaseURL,
			RequireHTTPS:        true,
			AllowedOrigins:      []string{appBaseURL},
			LogLevel:            "debug",
			SessionCookieSecure: true,
			EnableDebugLogging:  true,
		}
	default: // Development
		return EnvironmentConfig{
			Environment:         Development,
			AppBaseURL:          getEnv("APP_BASE_URL", "http://localhost:8080"),
			RequireHTTPS:        false,
			AllowedOrigins:      []string{"http://localhost:8080", "http://127.0.0.1:8080"},
			LogLevel:            "debug",
			SessionCookieSecure: false,
			EnableDebugLogging:  true,
		}
	}
}

// GetSupabaseConfig returns environment-appropriate Supabase configuration
func GetSupabaseConfig() SupabaseConfig {
	env := GetEnvironment()

	// Base configuration from environment variables
	baseConfig := SupabaseConfig{
		URL:            getEnv("SUPABASE_URL", ""),
		AnonKey:        getEnv("SUPABASE_ANON_KEY", ""),
		ServiceRoleKey: getEnv("SUPABASE_SERVICE_ROLE_KEY", ""),
	}

	// Environment-specific logic
	switch env {
	case Development:
		// In development, PublicURL might be different (ngrok)
		baseConfig.PublicURL = getEnv("SUPABASE_PUBLIC_URL", baseConfig.URL)
	default:
		// In staging/production, PublicURL is same as URL
		baseConfig.PublicURL = baseConfig.URL
	}

	// Validate required fields
	if baseConfig.URL == "" {
		panic(fmt.Sprintf("SUPABASE_URL is required for %s environment", env))
	}

	return baseConfig
}
