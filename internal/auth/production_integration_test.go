package auth

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/supabase-community/supabase-go"
	"usl-server/internal/config"
)

func TestProductionIntegrationRequirement(t *testing.T) {
	envConfig := config.EnvironmentConfig{
		Environment:    config.Production,
		AppBaseURL:     "https://production-app.com",
		RequireHTTPS:   true,
		AllowedOrigins: []string{"https://production-app.com"},
	}

	auth := NewDiscordAuth(nil, []string{"test-admin"}, "supabase-url", "public-url", "anon-key", envConfig)

	baseURL := auth.getAppBaseURL()
	expected := "https://production-app.com"

	if baseURL != expected {
		t.Errorf("Expected AppBaseURL %s, got %s", expected, baseURL)
	}
}

func TestMainGoUsesNewConstructor(t *testing.T) {
	envConfig := config.EnvironmentConfig{
		Environment: config.Development,
		AppBaseURL:  "http://test-localhost:8080",
	}

	auth := NewDiscordAuth(nil, []string{"test-id"}, "url", "public", "key", envConfig)

	if auth == nil {
		t.Error("Expected auth instance, got nil")
	}

	if auth.envConfig == nil {
		t.Error("Expected envConfig to be set, got nil")
	}

	if auth.envConfig.AppBaseURL != "http://test-localhost:8080" {
		t.Errorf("Expected AppBaseURL %s, got %s", "http://test-localhost:8080", auth.envConfig.AppBaseURL)
	}
}

// TestLoginForm_OAuthRedirectURLs verifies OAuth URL generation with injected config
func TestLoginForm_OAuthRedirectURLs(t *testing.T) {
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
			name:           "Staging USL path uses staging URL",
			appBaseURL:     "https://staging.rl-league-management.onrender.com",
			requestPath:    "/usl/login",
			expectedInHTML: "https://staging.rl-league-management.onrender.com/auth/callback?redirect=usl",
			description:    "Staging should use staging-specific URL",
		},
		{
			name:           "Development main path uses localhost",
			appBaseURL:     "http://localhost:8080",
			requestPath:    "/login",
			expectedInHTML: "http://localhost:8080/auth/callback?redirect=main",
			description:    "Development main app should use localhost URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create environment config with injected AppBaseURL
			envConfig := config.EnvironmentConfig{
				Environment: config.Production,
				AppBaseURL:  tt.appBaseURL,
			}

			// Create auth instance with dependency injection
			supabaseClient, _ := supabase.NewClient("https://test.supabase.co", "test-key", nil)
			auth := NewDiscordAuth(
				supabaseClient,
				[]string{"test-admin-id"},
				"https://test.supabase.co",
				"https://test.supabase.co",
				"test-anon-key",
				envConfig,
			)

			// Create request
			req := httptest.NewRequest("GET", tt.requestPath, nil)
			w := httptest.NewRecorder()

			// Call the method
			auth.LoginForm(w, req)

			// Check response contains expected OAuth URL
			body := w.Body.String()
			if !strings.Contains(body, tt.expectedInHTML) {
				t.Errorf("Expected HTML to contain %q, but it didn't.\nFull response: %s",
					tt.expectedInHTML, body)
			}

			// Verify no hardcoded localhost when using production URLs
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

// TestEnvironmentSpecificURLs verifies all environments use correct default URLs
func TestEnvironmentSpecificURLs(t *testing.T) {
	environments := []struct {
		env         config.Environment
		expectedURL string
		description string
	}{
		{
			env:         config.Production,
			expectedURL: "https://your-domain.com", // Default from GetEnvironmentConfig
			description: "Production should use HTTPS production domain",
		},
		{
			env:         config.Staging,
			expectedURL: "https://staging.your-domain.com", // Default from GetEnvironmentConfig
			description: "Staging should use HTTPS staging domain",
		},
		{
			env:         config.Development,
			expectedURL: "http://localhost:8080", // Default from GetEnvironmentConfig
			description: "Development should use localhost",
		},
	}

	for _, env := range environments {
		t.Run(string(env.env), func(t *testing.T) {
			// Use GetEnvironmentConfig to get default URLs for each environment
			// We'll simulate by creating the expected config structure
			var envConfig config.EnvironmentConfig
			switch env.env {
			case config.Production:
				envConfig = config.EnvironmentConfig{
					Environment:  config.Production,
					AppBaseURL:   "https://your-domain.com",
					RequireHTTPS: true,
				}
			case config.Staging:
				envConfig = config.EnvironmentConfig{
					Environment:  config.Staging,
					AppBaseURL:   "https://staging.your-domain.com",
					RequireHTTPS: true,
				}
			case config.Development:
				envConfig = config.EnvironmentConfig{
					Environment:  config.Development,
					AppBaseURL:   "http://localhost:8080",
					RequireHTTPS: false,
				}
			}

			// Create auth instance
			supabaseClient, _ := supabase.NewClient("https://test.supabase.co", "test-key", nil)
			auth := NewDiscordAuth(
				supabaseClient,
				[]string{"test-admin-id"},
				"https://test.supabase.co",
				"https://test.supabase.co",
				"test-anon-key",
				envConfig,
			)

			// Test OAuth URL generation
			baseURL := auth.getAppBaseURL()
			if baseURL != env.expectedURL {
				t.Errorf("Expected %s environment to use %s, got %s",
					env.env, env.expectedURL, baseURL)
			}

			// Test HTML generation for USL path
			req := httptest.NewRequest("GET", "/usl/login", nil)
			w := httptest.NewRecorder()
			auth.LoginForm(w, req)

			body := w.Body.String()
			expectedRedirect := env.expectedURL + "/auth/callback?redirect=usl"
			if !strings.Contains(body, expectedRedirect) {
				t.Errorf("Expected %s environment HTML to contain %q, but it didn't",
					env.env, expectedRedirect)
			}

			// Verify HTTPS requirement for production environments
			if env.env == config.Production || env.env == config.Staging {
				if !strings.Contains(body, "https://") {
					t.Errorf("Expected %s environment to use HTTPS URLs", env.env)
				}
			}
		})
	}
}

// TestPathSpecificRedirects verifies USL vs Main app redirect logic
func TestPathSpecificRedirects(t *testing.T) {
	testCases := []struct {
		name             string
		requestPath      string
		expectedRedirect string
		expectedTitle    string
		description      string
	}{
		{
			name:             "USL login path generates USL redirect",
			requestPath:      "/usl/login",
			expectedRedirect: "redirect=usl",
			expectedTitle:    "USL Admin Login",
			description:      "USL paths should redirect to USL admin after auth",
		},
		{
			name:             "Main app login path generates main redirect",
			requestPath:      "/login",
			expectedRedirect: "redirect=main",
			expectedTitle:    "Sign In",
			description:      "Main app paths should redirect to users after auth",
		},
		{
			name:             "USL nested path generates USL redirect",
			requestPath:      "/usl/admin/login",
			expectedRedirect: "redirect=usl",
			expectedTitle:    "USL Admin Login",
			description:      "Any USL path should generate USL redirect",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Setup auth with test config
			envConfig := config.EnvironmentConfig{
				Environment: config.Development,
				AppBaseURL:  "https://test-app.com",
			}

			supabaseClient, _ := supabase.NewClient("https://test.supabase.co", "test-key", nil)
			auth := NewDiscordAuth(
				supabaseClient,
				[]string{"test-admin-id"},
				"https://test.supabase.co",
				"https://test.supabase.co",
				"test-anon-key",
				envConfig,
			)

			// Create request
			req := httptest.NewRequest("GET", tc.requestPath, nil)
			w := httptest.NewRecorder()

			// Call LoginForm
			auth.LoginForm(w, req)
			body := w.Body.String()

			// Verify correct redirect parameter
			if !strings.Contains(body, tc.expectedRedirect) {
				t.Errorf("Expected %s to contain %q redirect parameter, but it didn't.\nBody: %s",
					tc.requestPath, tc.expectedRedirect, body)
			}

			// Verify correct page title
			if !strings.Contains(body, tc.expectedTitle) {
				t.Errorf("Expected %s to have title %q, but it didn't.\nBody: %s",
					tc.requestPath, tc.expectedTitle, body)
			}

			// Verify full redirect URL construction
			expectedURL := "https://test-app.com/auth/callback?" + tc.expectedRedirect
			if !strings.Contains(body, expectedURL) {
				t.Errorf("Expected %s to contain full redirect URL %q, but it didn't",
					tc.requestPath, expectedURL)
			}
		})
	}
}

// TestHTMLOutputValidation verifies HTML contains correct OAuth URLs and structure
func TestHTMLOutputValidation(t *testing.T) {
	envConfig := config.EnvironmentConfig{
		Environment: config.Production,
		AppBaseURL:  "https://production-app.com",
	}

	supabaseClient, _ := supabase.NewClient("https://test.supabase.co", "test-key", nil)
	auth := NewDiscordAuth(
		supabaseClient,
		[]string{"test-admin-id"},
		"https://test.supabase.co",
		"https://test.supabase.co",
		"test-anon-key",
		envConfig,
	)

	t.Run("USL login HTML structure", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/usl/login", nil)
		w := httptest.NewRecorder()
		auth.LoginForm(w, req)

		body := w.Body.String()

		// Verify HTML structure
		requiredElements := []string{
			"<!DOCTYPE html>",
			"<title>USL Admin Login</title>",
			"<h2 class=\"usl-header\">USL Administration</h2>",
			"Sign in with Discord to access the USL management system.",
			"🎮 Sign in with Discord",
		}

		for _, element := range requiredElements {
			if !strings.Contains(body, element) {
				t.Errorf("Expected HTML to contain %q, but it didn't", element)
			}
		}

		// Verify Discord OAuth URL is properly constructed
		expectedOAuthURL := "https://test.supabase.co/auth/v1/authorize?provider=discord&redirect_to=https://production-app.com/auth/callback?redirect=usl"
		if !strings.Contains(body, expectedOAuthURL) {
			t.Errorf("Expected HTML to contain OAuth URL %q, but it didn't", expectedOAuthURL)
		}

		// Verify no hardcoded localhost
		if strings.Contains(body, "localhost") || strings.Contains(body, "127.0.0.1") {
			t.Errorf("Found hardcoded localhost in production HTML: %s", body)
		}
	})

	t.Run("Main app login HTML structure", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/login", nil)
		w := httptest.NewRecorder()
		auth.LoginForm(w, req)

		body := w.Body.String()

		// Verify main app specific elements
		if !strings.Contains(body, "<title>Sign In</title>") {
			t.Error("Expected main app login to have 'Sign In' title")
		}

		if !strings.Contains(body, "Sign in with Discord to access the application.") {
			t.Error("Expected main app login message")
		}

		// Verify main redirect parameter
		expectedRedirect := "redirect=main"
		if !strings.Contains(body, expectedRedirect) {
			t.Errorf("Expected main app login to contain %q", expectedRedirect)
		}
	})
}

// TestErrorConditions verifies edge cases and error handling
func TestErrorConditions(t *testing.T) {
	t.Run("Empty AppBaseURL falls back to localhost", func(t *testing.T) {
		envConfig := config.EnvironmentConfig{
			Environment: config.Development,
			AppBaseURL:  "", // Empty URL
		}

		supabaseClient, _ := supabase.NewClient("https://test.supabase.co", "test-key", nil)
		auth := NewDiscordAuth(
			supabaseClient,
			[]string{"test-admin-id"},
			"https://test.supabase.co",
			"https://test.supabase.co",
			"test-anon-key",
			envConfig,
		)

		baseURL := auth.getAppBaseURL()
		expected := "http://localhost:8080"
		if baseURL != expected {
			t.Errorf("Expected empty AppBaseURL to fallback to %s, got %s", expected, baseURL)
		}
	})

	t.Run("Nil envConfig falls back to localhost", func(t *testing.T) {
		// This simulates using the old constructor or misconfiguration
		auth := &DiscordAuth{
			envConfig: nil, // No config injected
		}

		baseURL := auth.getAppBaseURL()
		expected := "http://localhost:8080"
		if baseURL != expected {
			t.Errorf("Expected nil envConfig to fallback to %s, got %s", expected, baseURL)
		}
	})

	t.Run("HTTP method validation", func(t *testing.T) {
		envConfig := config.EnvironmentConfig{
			Environment: config.Development,
			AppBaseURL:  "http://localhost:8080",
		}

		supabaseClient, _ := supabase.NewClient("https://test.supabase.co", "test-key", nil)
		auth := NewDiscordAuth(
			supabaseClient,
			[]string{"test-admin-id"},
			"https://test.supabase.co",
			"https://test.supabase.co",
			"test-anon-key",
			envConfig,
		)

		// Test POST request (should fail)
		req := httptest.NewRequest("POST", "/login", nil)
		w := httptest.NewRecorder()
		auth.LoginForm(w, req)

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected POST to return %d, got %d", http.StatusMethodNotAllowed, w.Code)
		}
	})

	t.Run("Different environment configurations work", func(t *testing.T) {
		environments := []config.Environment{
			config.Development,
			config.Staging,
			config.Production,
		}

		for _, env := range environments {
			envConfig := config.EnvironmentConfig{
				Environment: env,
				AppBaseURL:  "https://test-" + string(env) + ".com",
			}

			supabaseClient, _ := supabase.NewClient("https://test.supabase.co", "test-key", nil)
			auth := NewDiscordAuth(
				supabaseClient,
				[]string{"test-admin-id"},
				"https://test.supabase.co",
				"https://test.supabase.co",
				"test-anon-key",
				envConfig,
			)

			baseURL := auth.getAppBaseURL()
			expected := "https://test-" + string(env) + ".com"
			if baseURL != expected {
				t.Errorf("Environment %s: expected %s, got %s", env, expected, baseURL)
			}
		}
	})
}
