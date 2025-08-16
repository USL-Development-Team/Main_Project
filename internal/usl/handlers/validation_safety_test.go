package handlers

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

// TestValidationSafety tests that validation prevents invalid data from being processed
// This is critical for database safety - invalid data should never reach the database
func TestValidationSafety(t *testing.T) {
	// Create a handler instance for validation testing
	handler := &MigrationHandler{}

	tests := []struct {
		name           string
		formData       map[string]string
		shouldBeValid  bool
		expectedErrors []string
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
			expectedErrors: []string{"Profile URL is required"},
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
			expectedErrors: []string{"Must provide MMR data for at least one playlist"},
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
			expectedErrors: []string{"Must be a valid Rocket League tracker URL"},
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
			expectedErrors: []string{"Current peak must be between 0 and 3000"},
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
			expectedErrors: []string{"Current games must be between 0 and 10000"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock HTTP request with form data
			formValues := url.Values{}
			for key, value := range tt.formData {
				formValues.Set(key, value)
			}

			req := httptest.NewRequest("POST", "/test", strings.NewReader(formValues.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// Build tracker from form data (this is what the real handler does)
			tracker := handler.buildTrackerFromForm(req)

			// Validate the tracker (this is the critical safety check)
			validation := handler.validateTracker(tracker)

			// Check if validation result matches expectations
			if validation.IsValid != tt.shouldBeValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.shouldBeValid, validation.IsValid)
				t.Logf("Validation errors: %+v", validation.Errors)
			}

			// Check specific error messages for invalid cases
			if !tt.shouldBeValid {
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
			}

			// CRITICAL: If validation failed, this data should NEVER reach the database
			if !validation.IsValid {
				t.Logf("✅ Validation correctly rejected invalid data: %+v", tracker)
				t.Logf("✅ This data will NOT pollute the database")
			} else {
				t.Logf("✅ Validation passed for valid data: %+v", tracker)
			}
		})
	}
}

// TestFormParsingSafety tests that form parsing itself is safe
func TestFormParsingSafety(t *testing.T) {
	handler := &MigrationHandler{}

	maliciousInputs := []struct {
		name     string
		formData map[string]string
	}{
		{
			name: "SQL injection attempt in Discord ID",
			formData: map[string]string{
				"discord_id": "'; DROP TABLE users; --",
				"url":        "https://rocketleague.tracker.network/profile/123",
			},
		},
		{
			name: "XSS attempt in URL",
			formData: map[string]string{
				"discord_id": "123456789012345678",
				"url":        "<script>alert('xss')</script>",
			},
		},
		{
			name: "Extremely long Discord ID",
			formData: map[string]string{
				"discord_id": strings.Repeat("1", 1000), // 1000 characters
				"url":        "https://rocketleague.tracker.network/profile/123",
			},
		},
	}

	for _, tt := range maliciousInputs {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with malicious data
			formValues := url.Values{}
			for key, value := range tt.formData {
				formValues.Set(key, value)
			}

			req := httptest.NewRequest("POST", "/test", strings.NewReader(formValues.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// Parse form data
			tracker := handler.buildTrackerFromForm(req)

			// Validate - this should catch malicious input
			validation := handler.validateTracker(tracker)

			// Malicious input should always be rejected
			if validation.IsValid {
				t.Errorf("SECURITY ISSUE: Malicious input was accepted: %+v", tracker)
			} else {
				t.Logf("✅ Security: Malicious input correctly rejected")
				t.Logf("   Input: %+v", tt.formData)
				t.Logf("   Errors: %+v", validation.Errors)
			}
		})
	}
}
