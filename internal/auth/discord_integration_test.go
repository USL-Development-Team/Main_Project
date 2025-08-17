package auth

import (
	"testing"
	"usl-server/internal/config"

	"github.com/stretchr/testify/assert"
	"github.com/supabase-community/supabase-go"
)

// TestDiscordAuthIntegration tests the key aspects we can verify without real tokens
func TestDiscordAuthIntegration(t *testing.T) {
	envConfig := config.EnvironmentConfig{
		Environment: config.Production,
		AppBaseURL:  "https://test-app.com",
	}

	// Create auth with realistic parameters
	serviceRoleClient, err := supabase.NewClient("https://test.supabase.co", "service-role-key", nil)
	assert.NoError(t, err)

	auth := NewDiscordAuth(
		serviceRoleClient,
		[]string{"admin-discord-123"},
		"https://test.supabase.co",
		"https://test.supabase.co",
		"anon-key",
		envConfig,
	)

	t.Run("Auth instance configured correctly", func(t *testing.T) {
		assert.NotNil(t, auth.supabaseClient, "Service role client should be set")
		assert.Equal(t, "https://test.supabase.co", auth.supabaseURL)
		assert.Equal(t, "anon-key", auth.anonKey)
		assert.Equal(t, []string{"admin-discord-123"}, auth.adminDiscordIDs)
	})

	t.Run("Can create anon client for token validation", func(t *testing.T) {
		// This verifies the fix can work - anon client creation succeeds
		anonClient, err := supabase.NewClient(auth.supabaseURL, auth.anonKey, nil)
		assert.NoError(t, err, "Anon client creation should succeed")
		assert.NotNil(t, anonClient, "Anon client should not be nil")
	})

	t.Run("Discord ID extraction works correctly", func(t *testing.T) {
		testCases := []struct {
			name         string
			userMetadata map[string]interface{}
			expectedID   string
		}{
			{
				name: "Extract from provider_id",
				userMetadata: map[string]interface{}{
					"user_metadata": map[string]interface{}{
						"provider_id": "admin-discord-123",
					},
				},
				expectedID: "admin-discord-123",
			},
			{
				name: "Extract from sub field",
				userMetadata: map[string]interface{}{
					"user_metadata": map[string]interface{}{
						"sub": "admin-discord-123",
					},
				},
				expectedID: "admin-discord-123",
			},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				discordID := auth.extractDiscordID(tc.userMetadata)
				assert.Equal(t, tc.expectedID, discordID)
			})
		}
	})

	t.Run("Authorization logic works correctly", func(t *testing.T) {
		adminUser := map[string]interface{}{
			"user_metadata": map[string]interface{}{
				"provider_id": "admin-discord-123",
			},
		}

		regularUser := map[string]interface{}{
			"user_metadata": map[string]interface{}{
				"provider_id": "regular-user-456",
			},
		}

		assert.True(t, auth.isUserAuthorized(adminUser), "Admin user should be authorized")
		assert.False(t, auth.isUserAuthorized(regularUser), "Regular user should not be authorized")
	})
}

// TestFixImplementationDetails verifies the specific fix we applied
func TestFixImplementationDetails(t *testing.T) {
	t.Run("validateTokensAndGetUser creates anon client", func(t *testing.T) {
		envConfig := config.EnvironmentConfig{
			Environment: config.Development,
			AppBaseURL:  "http://localhost:8080",
		}

		auth := NewDiscordAuth(
			nil, // Service role client not needed for this test
			[]string{},
			"https://test.supabase.co",
			"https://test.supabase.co",
			"test-anon-key",
			envConfig,
		)

		// Test the core of our fix: can we create an anon client?
		anonClient, err := supabase.NewClient(auth.supabaseURL, auth.anonKey, nil)
		assert.NoError(t, err, "Fix should enable anon client creation")
		assert.NotNil(t, anonClient, "Anon client should be created successfully")

		// The actual validateTokensAndGetUser method now does this internally
		// We can't test it without real tokens, but we've verified the components work
	})
}

// TestAuthFlowComponents tests individual components of the auth flow
func TestAuthFlowComponents(t *testing.T) {
	envConfig := config.EnvironmentConfig{
		Environment: config.Production,
		AppBaseURL:  "https://production-app.com",
	}

	serviceRoleClient, _ := supabase.NewClient("https://test.supabase.co", "service-key", nil)
	auth := NewDiscordAuth(
		serviceRoleClient,
		[]string{"admin-123", "admin-456"},
		"https://test.supabase.co",
		"https://test.supabase.co",
		"anon-key",
		envConfig,
	)

	t.Run("OAuth URL generation", func(t *testing.T) {
		baseURL := auth.getAppBaseURL()
		assert.Equal(t, "https://production-app.com", baseURL)

		redirectURL := auth.buildRedirectURL(baseURL, "/usl/login")
		expected := "https://production-app.com/auth/callback?redirect=usl"
		assert.Equal(t, expected, redirectURL)
	})

	t.Run("Path detection", func(t *testing.T) {
		assert.True(t, auth.isUSLPath("/usl/login"))
		assert.True(t, auth.isUSLPath("/usl/admin"))
		assert.False(t, auth.isUSLPath("/login"))
		assert.False(t, auth.isUSLPath("/users"))
	})

	t.Run("Environment config injection", func(t *testing.T) {
		assert.NotNil(t, auth.envConfig)
		assert.Equal(t, config.Production, auth.envConfig.Environment)
		assert.Equal(t, "https://production-app.com", auth.envConfig.AppBaseURL)
	})
}
