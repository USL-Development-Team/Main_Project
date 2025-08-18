package test

import (
	"os"
	"strings"
	"testing"

	"usl-server/internal/auth"
	"usl-server/internal/config"
)

// TestSmokeConfiguration validates basic application startup with production-like config
func TestSmokeConfiguration(t *testing.T) {
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

	t.Run("production environment can start with valid config", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "production")
		os.Setenv("APP_BASE_URL", "https://rl-league-management.onrender.com")

		// This should not panic
		envConfig := config.GetEnvironmentConfig()

		if envConfig.AppBaseURL != "https://rl-league-management.onrender.com" {
			t.Errorf("Expected production URL, got %v", envConfig.AppBaseURL)
		}
	})

	t.Run("staging environment can start with valid config", func(t *testing.T) {
		os.Setenv("ENVIRONMENT", "staging")
		os.Setenv("APP_BASE_URL", "https://staging.rl-league-management.onrender.com")

		// This should not panic
		envConfig := config.GetEnvironmentConfig()

		if envConfig.AppBaseURL != "https://staging.rl-league-management.onrender.com" {
			t.Errorf("Expected staging URL, got %v", envConfig.AppBaseURL)
		}
	})
}

// TestSmokeOAuthURLGeneration validates OAuth URLs contain real domains, not placeholders
func TestSmokeOAuthURLGeneration(t *testing.T) {
	testCases := []struct {
		name        string
		environment string
		appBaseURL  string
		supabaseURL string
	}{
		{
			name:        "production OAuth URL",
			environment: "production",
			appBaseURL:  "https://rl-league-management.onrender.com",
			supabaseURL: "https://fhdsksvvswfqutvcqqjr.supabase.co",
		},
		{
			name:        "staging OAuth URL",
			environment: "staging",
			appBaseURL:  "https://staging.rl-league-management.onrender.com",
			supabaseURL: "https://fhdsksvvswfqutvcqqjr.supabase.co",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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

			os.Setenv("ENVIRONMENT", tc.environment)
			os.Setenv("APP_BASE_URL", tc.appBaseURL)

			envConfig := config.GetEnvironmentConfig()

			discordAuth := auth.NewDiscordAuth(auth.DiscordAuthConfig{
				SupabaseClient:  nil, // supabase client not needed for URL generation
				AdminDiscordIDs: []string{"test-admin"},
				SupabaseURL:     tc.supabaseURL,
				PublicURL:       tc.supabaseURL, // publicURL same as URL
				AnonKey:         "test-anon-key",
				EnvConfig:       &envConfig,
			})

			appBaseURL := discordAuth.GetAppBaseURL()

			// Ensure no placeholder domains
			if strings.Contains(appBaseURL, "your-domain.com") {
				t.Errorf("OAuth URL contains placeholder domain: %s", appBaseURL)
			}
			if strings.Contains(appBaseURL, "staging.your-domain.com") {
				t.Errorf("OAuth URL contains placeholder domain: %s", appBaseURL)
			}

			if appBaseURL != tc.appBaseURL {
				t.Errorf("Expected app base URL %s, got %s", tc.appBaseURL, appBaseURL)
			}
		})
	}
}

// TestSmokeApplicationStartup validates the entire application can start
func TestSmokeApplicationStartup(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping application startup test in short mode")
	}

	// This would be a more complex test that actually starts the server
	// For now, just validate configuration loading doesn't panic

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

	os.Setenv("ENVIRONMENT", "production")
	os.Setenv("APP_BASE_URL", "https://test-production.com")
	appConfig, err := config.Load()
	if err != nil {
		t.Fatalf("Failed to load configuration: %v", err)
	}

	envConfig := config.GetEnvironmentConfig()
	if appConfig == nil {
		t.Error("Expected configuration, got nil")
	}
	if envConfig.AppBaseURL != "https://test-production.com" {
		t.Errorf("Expected test production URL, got %v", envConfig.AppBaseURL)
	}
}
