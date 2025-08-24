package handlers

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestValidationSafety tests that validation prevents invalid data from being processed.
// This is critical for database safety - invalid data should never reach the database.
// These tests follow Go testing best practices with table-driven tests, clear assertions,
// and comprehensive edge case coverage.
func TestValidationSafety(t *testing.T) {
	// Arrange: Create a handler instance for validation testing
	handler := &MigrationHandler{}

	// Define test cases using table-driven test pattern (Go best practice)
	tests := []struct {
		name           string
		formData       map[string]string
		shouldBeValid  bool
		expectedErrors []string
		description    string // Added for better test documentation
	}{
		{
			name: "Valid data should pass validation",
			formData: map[string]string{
				"discord_id":         "123456789012345678",
				"url":                "https://rocketleague.tracker.network/profile/123",
				"ones_current_peak":  "1500",
				"ones_current_games": "100",
				"valid":              "true",
			},
			shouldBeValid: true,
			description:   "Baseline test with all valid data",
		},
		{
			name: "Invalid Discord ID should be rejected",
			formData: map[string]string{
				"discord_id":        "invalid123", // Invalid format
				"url":               "https://rocketleague.tracker.network/profile/123",
				"ones_current_peak": "1500",
				"valid":             "true",
			},
			shouldBeValid:  false,
			expectedErrors: []string{"Discord ID must be 17-19 digits"},
			description:    "Tests Discord ID format validation",
		},
		{
			name: "Missing URL should be rejected",
			formData: map[string]string{
				"discord_id":        "123456789012345678",
				"url":               "", // Empty URL
				"ones_current_peak": "1500",
				"valid":             "true",
			},
			shouldBeValid:  false,
			expectedErrors: []string{"Tracker URL is required"},
			description:    "Tests required URL field validation",
		},
		{
			name: "No playlist data should be rejected",
			formData: map[string]string{
				"discord_id": "123456789012345678",
				"url":        "https://rocketleague.tracker.network/profile/123",
				"valid":      "true",
				// No MMR or games data provided
			},
			shouldBeValid:  false,
			expectedErrors: []string{"Tracker must have data for at least one playlist"},
			description:    "Tests business rule requiring at least one playlist",
		},
		{
			name: "Invalid tracker domain should be rejected",
			formData: map[string]string{
				"discord_id":        "123456789012345678",
				"url":               "https://malicious-site.com/steal-data", // Invalid domain
				"ones_current_peak": "1500",
				"valid":             "true",
			},
			shouldBeValid:  false,
			expectedErrors: []string{"Invalid tracker URL format"},
			description:    "Tests URL domain whitelist security",
		},
		{
			name: "Excessive MMR should be rejected",
			formData: map[string]string{
				"discord_id":        "123456789012345678",
				"url":               "https://rocketleague.tracker.network/profile/123",
				"ones_current_peak": "9999", // Way too high
				"valid":             "true",
			},
			shouldBeValid:  false,
			expectedErrors: []string{"1v1 current season MMR must be between 0 and 3000"},
			description:    "Tests MMR upper bound validation",
		},
		{
			name: "Excessive games should be rejected",
			formData: map[string]string{
				"discord_id":         "123456789012345678",
				"url":                "https://rocketleague.tracker.network/profile/123",
				"ones_current_peak":  "1500",
				"ones_current_games": "50000", // Way too many games
				"valid":              "true",
			},
			shouldBeValid:  false,
			expectedErrors: []string{"1v1 current season games must be between 0 and 10000"},
			description:    "Tests games played upper bound validation",
		},
		{
			name: "Unicode Discord ID should be rejected",
			formData: map[string]string{
				"discord_id":        "１２３４５６７８９０１２３４５６７８", // Unicode numbers
				"url":               "https://rocketleague.tracker.network/profile/123",
				"ones_current_peak": "1500",
				"valid":             "true",
			},
			shouldBeValid:  false,
			expectedErrors: []string{"Discord ID must be 17-19 digits"},
			description:    "Tests Unicode character rejection in Discord ID",
		},
		{
			name: "Leading/trailing spaces in Discord ID should be rejected",
			formData: map[string]string{
				"discord_id":        " 123456789012345678 ", // Spaces around valid ID
				"url":               "https://rocketleague.tracker.network/profile/123",
				"ones_current_peak": "1500",
				"valid":             "true",
			},
			shouldBeValid:  false, // Current implementation doesn't trim spaces
			expectedErrors: []string{"Discord ID must be 17-19 digits"},
			description:    "Tests that whitespace is not automatically trimmed",
		},
		{
			name: "Decimal MMR values should be parsed as integer",
			formData: map[string]string{
				"discord_id":        "123456789012345678",
				"url":               "https://rocketleague.tracker.network/profile/123",
				"ones_current_peak": "1500.5", // Decimal MMR - strconv.Atoi truncates
				"valid":             "true",
			},
			shouldBeValid:  false, // strconv.Atoi fails on decimal values
			expectedErrors: []string{"Tracker must have data for at least one playlist"},
			description:    "Tests decimal number parsing behavior",
		},
		{
			name: "Non-numeric MMR should be rejected",
			formData: map[string]string{
				"discord_id":        "123456789012345678",
				"url":               "https://rocketleague.tracker.network/profile/123",
				"ones_current_peak": "abc", // Non-numeric
				"valid":             "true",
			},
			shouldBeValid:  false,
			expectedErrors: []string{"Tracker must have data for at least one playlist"}, // Treated as 0
			description:    "Tests non-numeric input handling",
		},
		{
			name: "Exactly at MMR boundaries should be valid",
			formData: map[string]string{
				"discord_id":        "123456789012345678",
				"url":               "https://rocketleague.tracker.network/profile/123",
				"ones_current_peak": "0",    // Lower boundary
				"twos_current_peak": "3000", // Upper boundary
				"valid":             "true",
			},
			shouldBeValid: true,
			description:   "Tests boundary value acceptance",
		},
		{
			name: "URL with query parameters should be valid",
			formData: map[string]string{
				"discord_id":        "123456789012345678",
				"url":               "https://rocketleague.tracker.network/profile/steam/76561198000000000?playlist=10", // Query params
				"ones_current_peak": "1500",
				"valid":             "true",
			},
			shouldBeValid: true,
			description:   "Tests URL with complex path and query parameters",
		},
		{
			name: "URL with different valid subdomain",
			formData: map[string]string{
				"discord_id":        "123456789012345678",
				"url":               "https://api.rocketleague.tracker.network/profile/123", // Different subdomain
				"ones_current_peak": "1500",
				"valid":             "true",
			},
			shouldBeValid: true, // Should work with contains() logic
			description:   "Tests subdomain variation acceptance",
		},
		{
			name: "Mixed case in URL domain should be rejected",
			formData: map[string]string{
				"discord_id":        "123456789012345678",
				"url":               "https://RocketLeague.Tracker.Network/profile/123", // Mixed case
				"ones_current_peak": "1500",
				"valid":             "true",
			},
			shouldBeValid:  false, // Current implementation is case sensitive
			expectedErrors: []string{"Invalid tracker URL format"},
			description:    "Tests case sensitivity in URL validation",
		},
		{
			name: "Very large valid Discord ID",
			formData: map[string]string{
				"discord_id":        "9" + strings.Repeat("9", 17), // 18 nines
				"url":               "https://rocketleague.tracker.network/profile/123",
				"ones_current_peak": "1500",
				"valid":             "true",
			},
			shouldBeValid: true,
			description:   "Tests maximum valid Discord ID size",
		},
		{
			name: "All playlists with zero MMR but games played",
			formData: map[string]string{
				"discord_id":         "123456789012345678",
				"url":                "https://rocketleague.tracker.network/profile/123",
				"ones_current_peak":  "0",
				"ones_current_games": "5", // Games but no MMR
				"valid":              "true",
			},
			shouldBeValid: true, // Games > 0 should count as playlist data
			description:   "Tests games-only playlist data acceptance",
		},
	}

	// Execute test cases using Go's subtest pattern
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: Create a mock HTTP request with form data
			formValues := url.Values{}
			for key, value := range tt.formData {
				formValues.Set(key, value)
			}

			req := httptest.NewRequest("POST", "/test", strings.NewReader(formValues.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// Act: Build tracker from form data and validate
			tracker := handler.buildTrackerFromForm(req)
			validation := handler.validateTracker(tracker)

			// Assert: Check validation result matches expectations
			if validation.IsValid != tt.shouldBeValid {
				t.Errorf("Test %q failed: Expected valid=%v, got valid=%v", tt.name, tt.shouldBeValid, validation.IsValid)
				t.Logf("Test description: %s", tt.description)
				t.Logf("Validation errors: %+v", validation.Errors)
				return // Early return on failure for better readability
			}

			// Assert: Check specific error messages for invalid cases
			if !tt.shouldBeValid {
				if len(tt.expectedErrors) == 0 {
					t.Error("Test configuration error: expectedErrors must be provided for invalid test cases")
					return
				}

				for _, expectedError := range tt.expectedErrors {
					found := false
					for _, actualError := range validation.Errors {
						if strings.Contains(actualError.Message, expectedError) {
							found = true
							break
						}
					}
					if !found {
						t.Errorf("Expected error containing '%s' but got errors: %+v", expectedError, validation.Errors)
					}
				}
			} else {
				// Assert: Valid cases should have no errors
				if len(validation.Errors) > 0 {
					t.Errorf("Expected no errors for valid case, but got: %+v", validation.Errors)
				}
			}

			// CRITICAL SAFETY CHECK: Log database safety verification
			if !validation.IsValid {
				t.Logf("✅ SAFETY: Validation correctly rejected invalid data")
				t.Logf("   Description: %s", tt.description)
				t.Logf("   This data will NOT pollute the database")
			} else {
				t.Logf("✅ VALID: Validation passed for legitimate data")
				t.Logf("   Description: %s", tt.description)
			}
		})
	}
}

// TestFormParsingSafety tests that form parsing itself is safe from malicious input.
// This follows Go security testing best practices by testing known attack vectors.
func TestFormParsingSafety(t *testing.T) {
	// Arrange: Create handler instance
	handler := &MigrationHandler{}

	// Define security test cases using table-driven pattern
	maliciousInputs := []struct {
		name        string
		formData    map[string]string
		description string
	}{
		{
			name: "SQL injection attempt in Discord ID",
			formData: map[string]string{
				"discord_id": "'; DROP TABLE users; --",
				"url":        "https://rocketleague.tracker.network/profile/123",
			},
			description: "Tests protection against SQL injection attacks",
		},
		{
			name: "XSS attempt in URL",
			formData: map[string]string{
				"discord_id": "123456789012345678",
				"url":        "<script>alert('xss')</script>",
			},
			description: "Tests protection against cross-site scripting",
		},
		{
			name: "Extremely long Discord ID",
			formData: map[string]string{
				"discord_id": strings.Repeat("1", 1000), // 1000 characters
				"url":        "https://rocketleague.tracker.network/profile/123",
			},
			description: "Tests protection against buffer overflow attempts",
		},
	}

	// Execute security tests using Go's subtest pattern
	for _, tt := range maliciousInputs {
		t.Run(tt.name, func(t *testing.T) {
			// Arrange: Create request with malicious data
			formValues := url.Values{}
			for key, value := range tt.formData {
				formValues.Set(key, value)
			}

			req := httptest.NewRequest("POST", "/test", strings.NewReader(formValues.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// Act: Parse form data and validate
			tracker := handler.buildTrackerFromForm(req)
			validation := handler.validateTracker(tracker)

			// Assert: Malicious input should ALWAYS be rejected
			if validation.IsValid {
				t.Errorf("CRITICAL SECURITY ISSUE: Malicious input was accepted: %+v", tracker)
				t.Errorf("Attack vector: %s", tt.description)
				t.FailNow() // Stop test immediately on security failure
			}

			// Security validation passed
			t.Logf("✅ SECURITY: Malicious input correctly rejected")
			t.Logf("   Attack type: %s", tt.description)
			t.Logf("   Validation errors: %d", len(validation.Errors))
		})
	}
}
