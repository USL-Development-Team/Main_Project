package config

import (
	"os"
	"testing"
)

func TestGetEnvironmentConfig_ProductionRequiresAppBaseURL(t *testing.T) {
	originalEnv := os.Getenv("ENVIRONMENT")
	originalAppBaseURL := os.Getenv("APP_BASE_URL")
	defer func() {
		os.Setenv("ENVIRONMENT", originalEnv)
		if originalAppBaseURL != "" {
			os.Setenv("APP_BASE_URL", originalAppBaseURL)
		} else {
			os.Unsetenv("APP_BASE_URL")
		}
	}()

	t.Run("production panics without APP_BASE_URL", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "production")
		os.Unsetenv("APP_BASE_URL")

		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when APP_BASE_URL is missing in production")
			} else if r != "APP_BASE_URL environment variable is required in production" {
				t.Errorf("Expected specific panic message, got: %v", r)
			}
		}()

		GetEnvironmentConfig()
	})

	t.Run("production works with APP_BASE_URL", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "production")
		os.Setenv("APP_BASE_URL", "https://example.com")

		config := GetEnvironmentConfig()

		if config.Environment != Production {
			t.Errorf("Expected Production environment, got %v", config.Environment)
		}
		if config.AppBaseURL != "https://example.com" {
			t.Errorf("Expected AppBaseURL to be https://example.com, got %v", config.AppBaseURL)
		}
		if len(config.AllowedOrigins) != 1 || config.AllowedOrigins[0] != "https://example.com" {
			t.Errorf("Expected AllowedOrigins to contain AppBaseURL, got %v", config.AllowedOrigins)
		}
	})
}

func TestGetEnvironmentConfig_StagingRequiresAppBaseURL(t *testing.T) {
	originalEnv := os.Getenv("ENVIRONMENT")
	originalAppBaseURL := os.Getenv("APP_BASE_URL")
	defer func() {
		os.Setenv("ENVIRONMENT", originalEnv)
		if originalAppBaseURL != "" {
			os.Setenv("APP_BASE_URL", originalAppBaseURL)
		} else {
			os.Unsetenv("APP_BASE_URL")
		}
	}()

	t.Run("staging panics without APP_BASE_URL", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "staging")
		os.Unsetenv("APP_BASE_URL")

		defer func() {
			if r := recover(); r == nil {
				t.Error("Expected panic when APP_BASE_URL is missing in staging")
			} else if r != "APP_BASE_URL environment variable is required in staging" {
				t.Errorf("Expected specific panic message, got: %v", r)
			}
		}()

		GetEnvironmentConfig()
	})

	t.Run("staging works with APP_BASE_URL", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "staging")
		os.Setenv("APP_BASE_URL", "https://staging.example.com")

		config := GetEnvironmentConfig()

		if config.Environment != Staging {
			t.Errorf("Expected Staging environment, got %v", config.Environment)
		}
		if config.AppBaseURL != "https://staging.example.com" {
			t.Errorf("Expected AppBaseURL to be https://staging.example.com, got %v", config.AppBaseURL)
		}
	})
}

func TestGetEnvironmentConfig_DevelopmentUsesDefault(t *testing.T) {
	originalEnv := os.Getenv("ENVIRONMENT")
	originalAppBaseURL := os.Getenv("APP_BASE_URL")
	defer func() {
		os.Setenv("ENVIRONMENT", originalEnv)
		if originalAppBaseURL != "" {
			os.Setenv("APP_BASE_URL", originalAppBaseURL)
		} else {
			os.Unsetenv("APP_BASE_URL")
		}
	}()

	t.Run("development works without APP_BASE_URL", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "development")
		os.Unsetenv("APP_BASE_URL")

		config := GetEnvironmentConfig()

		if config.Environment != Development {
			t.Errorf("Expected Development environment, got %v", config.Environment)
		}
		if config.AppBaseURL != "http://localhost:8080" {
			t.Errorf("Expected default localhost URL, got %v", config.AppBaseURL)
		}
	})

	t.Run("development respects custom APP_BASE_URL", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "development")
		os.Setenv("APP_BASE_URL", "https://custom-dev.ngrok.io")

		config := GetEnvironmentConfig()

		if config.AppBaseURL != "https://custom-dev.ngrok.io" {
			t.Errorf("Expected custom URL, got %v", config.AppBaseURL)
		}
	})
}
