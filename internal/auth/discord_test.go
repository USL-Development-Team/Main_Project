package auth

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/supabase-community/supabase-go"
)

func TestDiscordAuth_LoginForm_RedirectURLs(t *testing.T) {
	tests := []struct {
		name           string
		appBaseURL     string
		requestPath    string
		expectedInHTML string
		description    string
	}{
		{
			name:           "Production USL path uses production URL",
			appBaseURL:     "https://rl-league-management.onrender.com",
			requestPath:    "/usl/login",
			expectedInHTML: "https://rl-league-management.onrender.com/auth/callback?redirect=usl",
			description:    "USL path should generate production redirect URL",
		},
		{
			name:           "Production main path uses production URL",
			appBaseURL:     "https://rl-league-management.onrender.com",
			requestPath:    "/login",
			expectedInHTML: "https://rl-league-management.onrender.com/auth/callback?redirect=main",
			description:    "Main app path should generate production redirect URL",
		},
		{
			name:           "Development USL path uses localhost",
			appBaseURL:     "http://localhost:8080",
			requestPath:    "/usl/login",
			expectedInHTML: "http://localhost:8080/auth/callback?redirect=usl",
			description:    "Development should use localhost URL",
		},
		{
			name:           "Development main path uses localhost",
			appBaseURL:     "http://localhost:8080",
			requestPath:    "/login",
			expectedInHTML: "http://localhost:8080/auth/callback?redirect=main",
			description:    "Development main app should use localhost URL",
		},
		{
			name:           "Staging environment uses staging URL",
			appBaseURL:     "https://staging.rl-league-management.onrender.com",
			requestPath:    "/usl/login",
			expectedInHTML: "https://staging.rl-league-management.onrender.com/auth/callback?redirect=usl",
			description:    "Staging should use staging-specific URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up environment
			oldAppBaseURL := os.Getenv("APP_BASE_URL")
			defer func() {
				if oldAppBaseURL != "" {
					os.Setenv("APP_BASE_URL", oldAppBaseURL)
				} else {
					os.Unsetenv("APP_BASE_URL")
				}
			}()

			os.Setenv("APP_BASE_URL", tt.appBaseURL)

			// Create auth instance
			supabaseClient, _ := supabase.NewClient("https://test.supabase.co", "test-key", nil)
			auth := NewDiscordAuth(
				supabaseClient,
				[]string{"test-admin-id"},
				"https://test.supabase.co",
				"https://test.supabase.co",
				"test-anon-key",
			)

			// Create request
			req := httptest.NewRequest("GET", tt.requestPath, nil)
			w := httptest.NewRecorder()

			// Call the method
			auth.LoginForm(w, req)

			// Check response
			body := w.Body.String()
			if !strings.Contains(body, tt.expectedInHTML) {
				t.Errorf("Expected HTML to contain %q, but it didn't.\nFull response: %s",
					tt.expectedInHTML, body)
			}

			// Verify it doesn't contain hardcoded localhost (unless that's what we expect)
			if tt.appBaseURL != "http://localhost:8080" {
				if strings.Contains(body, "http://127.0.0.1:8080") ||
					strings.Contains(body, "http://localhost:8080") {
					t.Errorf("Found hardcoded localhost in response when app base URL is %s.\nResponse: %s",
						tt.appBaseURL, body)
				}
			}
		})
	}
}

func TestDiscordAuth_RedirectURL_NoHardcodedValues(t *testing.T) {
	// This test ensures we never accidentally hardcode localhost again
	environments := []struct {
		name       string
		appBaseURL string
	}{
		{"production", "https://rl-league-management.onrender.com"},
		{"staging", "https://staging.rl-league-management.onrender.com"},
		{"development", "http://localhost:3000"}, // Different port to catch hardcoding
	}

	paths := []string{"/login", "/usl/login"}

	for _, env := range environments {
		for _, path := range paths {
			t.Run(env.name+"_"+path, func(t *testing.T) {
				// Set environment
				oldAppBaseURL := os.Getenv("APP_BASE_URL")
				defer func() {
					if oldAppBaseURL != "" {
						os.Setenv("APP_BASE_URL", oldAppBaseURL)
					} else {
						os.Unsetenv("APP_BASE_URL")
					}
				}()

				os.Setenv("APP_BASE_URL", env.appBaseURL)

				// Create auth instance
				supabaseClient, _ := supabase.NewClient("https://test.supabase.co", "test-key", nil)
				auth := NewDiscordAuth(
					supabaseClient,
					[]string{"test-admin-id"},
					"https://test.supabase.co",
					"https://test.supabase.co",
					"test-anon-key",
				)

				// Create request
				req := httptest.NewRequest("GET", path, nil)
				w := httptest.NewRecorder()

				// Call the method
				auth.LoginForm(w, req)

				// Check that response contains the expected base URL
				body := w.Body.String()
				if !strings.Contains(body, env.appBaseURL) {
					t.Errorf("Expected response to contain %q, but it didn't", env.appBaseURL)
				}

				// Check that it doesn't contain any hardcoded localhost values
				// (unless that's specifically what we set)
				hardcodedValues := []string{
					"http://127.0.0.1:8080",
					"http://localhost:8080", // The old hardcoded value
				}

				for _, hardcoded := range hardcodedValues {
					if hardcoded != env.appBaseURL && strings.Contains(body, hardcoded) {
						t.Errorf("Found hardcoded value %q in response when APP_BASE_URL is %q",
							hardcoded, env.appBaseURL)
					}
				}
			})
		}
	}
}

func TestDiscordAuth_EnvironmentIntegration(t *testing.T) {
	// Test that the auth system properly integrates with environment configuration
	testCases := []struct {
		name        string
		environment string
		appBaseURL  string
		expectHTTPS bool
	}{
		{
			name:        "Production environment uses HTTPS",
			environment: "production",
			appBaseURL:  "https://rl-league-management.onrender.com",
			expectHTTPS: true,
		},
		{
			name:        "Staging environment uses HTTPS",
			environment: "staging",
			appBaseURL:  "https://staging.rl-league-management.onrender.com",
			expectHTTPS: true,
		},
		{
			name:        "Development environment allows HTTP",
			environment: "development",
			appBaseURL:  "http://localhost:8080",
			expectHTTPS: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set up environment
			oldEnv := os.Getenv("ENVIRONMENT")
			oldAppBaseURL := os.Getenv("APP_BASE_URL")
			defer func() {
				if oldEnv != "" {
					os.Setenv("ENVIRONMENT", oldEnv)
				} else {
					os.Unsetenv("ENVIRONMENT")
				}
				if oldAppBaseURL != "" {
					os.Setenv("APP_BASE_URL", oldAppBaseURL)
				} else {
					os.Unsetenv("APP_BASE_URL")
				}
			}()

			os.Setenv("ENVIRONMENT", tc.environment)
			os.Setenv("APP_BASE_URL", tc.appBaseURL)

			// Create auth instance
			supabaseClient, _ := supabase.NewClient("https://test.supabase.co", "test-key", nil)
			auth := NewDiscordAuth(
				supabaseClient,
				[]string{"test-admin-id"},
				"https://test.supabase.co",
				"https://test.supabase.co",
				"test-anon-key",
			)

			// Test request
			req := httptest.NewRequest("GET", "/usl/login", nil)
			w := httptest.NewRecorder()

			auth.LoginForm(w, req)

			// Verify the response uses the correct protocol
			body := w.Body.String()
			if tc.expectHTTPS {
				if !strings.Contains(body, "https://") {
					t.Errorf("Expected HTTPS URL in %s environment, but didn't find it", tc.environment)
				}
				if strings.Contains(body, "http://") && !strings.Contains(body, "https://") {
					t.Errorf("Found HTTP instead of HTTPS in %s environment", tc.environment)
				}
			} else {
				// Development can use HTTP
				if !strings.Contains(body, tc.appBaseURL) {
					t.Errorf("Expected to find %s in response", tc.appBaseURL)
				}
			}
		})
	}
}

// Test helper to verify OAuth URL construction
func TestOAuthURLConstruction(t *testing.T) {
	testCases := []struct {
		name           string
		supabaseURL    string
		appBaseURL     string
		path           string
		expectedFormat string
	}{
		{
			name:           "USL OAuth URL format",
			supabaseURL:    "https://test.supabase.co",
			appBaseURL:     "https://myapp.com",
			path:           "/usl/login",
			expectedFormat: "https://test.supabase.co/auth/v1/authorize?provider=discord&redirect_to=https://myapp.com/auth/callback?redirect=usl",
		},
		{
			name:           "Main app OAuth URL format",
			supabaseURL:    "https://test.supabase.co",
			appBaseURL:     "https://myapp.com",
			path:           "/login",
			expectedFormat: "https://test.supabase.co/auth/v1/authorize?provider=discord&redirect_to=https://myapp.com/auth/callback?redirect=main",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Set environment
			oldAppBaseURL := os.Getenv("APP_BASE_URL")
			defer func() {
				if oldAppBaseURL != "" {
					os.Setenv("APP_BASE_URL", oldAppBaseURL)
				} else {
					os.Unsetenv("APP_BASE_URL")
				}
			}()

			os.Setenv("APP_BASE_URL", tc.appBaseURL)

			// Create auth instance
			supabaseClient, _ := supabase.NewClient(tc.supabaseURL, "test-key", nil)
			auth := NewDiscordAuth(
				supabaseClient,
				[]string{"test-admin-id"},
				tc.supabaseURL,
				tc.supabaseURL,
				"test-anon-key",
			)

			// Create request
			req := httptest.NewRequest("GET", tc.path, nil)
			w := httptest.NewRecorder()

			// Call the method
			auth.LoginForm(w, req)

			// Check that the OAuth URL is properly formatted
			body := w.Body.String()
			if !strings.Contains(body, tc.expectedFormat) {
				t.Errorf("Expected OAuth URL format not found.\nExpected: %s\nBody: %s",
					tc.expectedFormat, body)
			}
		})
	}
}
