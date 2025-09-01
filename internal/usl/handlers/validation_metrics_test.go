package handlers

import (
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

// TestValidationMetricsCollection tests that metrics are properly collected during validation
func TestValidationMetricsCollection(t *testing.T) {
	// Reset metrics for clean testing
	metricsMutex.Lock()
	validationMetrics = &ValidationMetrics{
		ErrorsByType:  make(map[string]int64),
		ErrorsByField: make(map[string]int64),
		LastReset:     time.Now(),
	}
	metricsMutex.Unlock()

	baseHandler := &BaseHandler{}

	t.Run("TestSuccessfulValidationMetrics", func(t *testing.T) {
		// Create valid form data
		formValues := url.Values{}
		formValues.Set("discord_id", "123456789012345678")
		formValues.Set("url", "https://rocketleague.tracker.network/profile/123")
		formValues.Set("ones_current_peak", "1500")

		req := httptest.NewRequest("POST", "/test", strings.NewReader(formValues.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Set("User-Agent", "Test-Agent/1.0")

		// Parse form
		err := req.ParseForm()
		if err != nil {
			t.Fatalf("Failed to parse form: %v", err)
		}

		tracker := baseHandler.buildTrackerFromForm(req)
		validation := baseHandler.validateTrackerWithMetrics(req, tracker)

		if !validation.IsValid {
			t.Errorf("Expected valid tracker, got errors: %+v", validation.Errors)
		}

		// Check metrics
		metrics := getValidationMetrics()
		if metrics.TotalValidations != 1 {
			t.Errorf("Expected 1 total validation, got %d", metrics.TotalValidations)
		}
		if metrics.SuccessfulValidations != 1 {
			t.Errorf("Expected 1 successful validation, got %d", metrics.SuccessfulValidations)
		}
		if metrics.FailedValidations != 0 {
			t.Errorf("Expected 0 failed validations, got %d", metrics.FailedValidations)
		}
		if metrics.SecurityIncidents != 0 {
			t.Errorf("Expected 0 security incidents, got %d", metrics.SecurityIncidents)
		}
	})

	t.Run("TestFailedValidationMetrics", func(t *testing.T) {
		// Reset for this test
		metricsMutex.Lock()
		validationMetrics = &ValidationMetrics{
			ErrorsByType:  make(map[string]int64),
			ErrorsByField: make(map[string]int64),
			LastReset:     time.Now(),
		}
		metricsMutex.Unlock()

		// Create invalid form data (empty Discord ID)
		formValues := url.Values{}
		formValues.Set("discord_id", "")
		formValues.Set("url", "https://rocketleague.tracker.network/profile/123")

		req := httptest.NewRequest("POST", "/test", strings.NewReader(formValues.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Parse form
		err := req.ParseForm()
		if err != nil {
			t.Fatalf("Failed to parse form: %v", err)
		}

		tracker := baseHandler.buildTrackerFromForm(req)
		validation := baseHandler.validateTrackerWithMetrics(req, tracker)

		if validation.IsValid {
			t.Error("Expected invalid tracker due to empty Discord ID")
		}

		// Check metrics
		metrics := getValidationMetrics()
		if metrics.TotalValidations != 1 {
			t.Errorf("Expected 1 total validation, got %d", metrics.TotalValidations)
		}
		if metrics.SuccessfulValidations != 0 {
			t.Errorf("Expected 0 successful validations, got %d", metrics.SuccessfulValidations)
		}
		if metrics.FailedValidations != 1 {
			t.Errorf("Expected 1 failed validation, got %d", metrics.FailedValidations)
		}

		// Check error categorization
		if metrics.ErrorsByType[ValidationCodeRequired] != 1 {
			t.Errorf("Expected 1 'required' error (Discord ID), got %d", metrics.ErrorsByType[ValidationCodeRequired])
		}
		if metrics.ErrorsByType[ValidationCodeNoData] != 1 {
			t.Errorf("Expected 1 'no_data' error (no playlist data), got %d", metrics.ErrorsByType[ValidationCodeNoData])
		}
		if metrics.ErrorsByField["discord_id"] != 1 {
			t.Errorf("Expected 1 discord_id error, got %d", metrics.ErrorsByField["discord_id"])
		}
	})

	t.Run("TestSecurityIncidentMetrics", func(t *testing.T) {
		// Reset for this test
		metricsMutex.Lock()
		validationMetrics = &ValidationMetrics{
			ErrorsByType:  make(map[string]int64),
			ErrorsByField: make(map[string]int64),
			LastReset:     time.Now(),
		}
		metricsMutex.Unlock()

		// Create malicious form data
		formValues := url.Values{}
		formValues.Set("discord_id", "'; DROP TABLE users; --")
		formValues.Set("url", "https://rocketleague.tracker.network/profile/123")

		req := httptest.NewRequest("POST", "/test", strings.NewReader(formValues.Encode()))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		// Parse form
		err := req.ParseForm()
		if err != nil {
			t.Fatalf("Failed to parse form: %v", err)
		}

		tracker := baseHandler.buildTrackerFromForm(req)
		validation := baseHandler.validateTrackerWithMetrics(req, tracker)

		if validation.IsValid {
			t.Error("Expected invalid tracker due to malicious input")
		}

		// Check metrics
		metrics := getValidationMetrics()
		if metrics.TotalValidations != 1 {
			t.Errorf("Expected 1 total validation, got %d", metrics.TotalValidations)
		}
		if metrics.SecurityIncidents != 1 {
			t.Errorf("Expected 1 security incident, got %d", metrics.SecurityIncidents)
		}
		if metrics.FailedValidations != 1 {
			t.Errorf("Expected 1 failed validation, got %d", metrics.FailedValidations)
		}
	})
}

// TestValidationMetricsAPI tests the metrics endpoint
func TestValidationMetricsAPI(t *testing.T) {
	// Reset metrics for clean testing
	metricsMutex.Lock()
	validationMetrics = &ValidationMetrics{
		TotalValidations:      10,
		SuccessfulValidations: 8,
		FailedValidations:     2,
		SecurityIncidents:     1,
		ErrorsByType:          map[string]int64{"required": 3, "invalid_format": 2},
		ErrorsByField:         map[string]int64{"discord_id": 4, "url": 1},
		LastReset:             time.Now(),
	}
	metricsMutex.Unlock()

	adminHandler := &AdminHandler{BaseHandler: &BaseHandler{}}

	req := httptest.NewRequest("GET", "/api/validation/metrics", nil)
	w := httptest.NewRecorder()

	adminHandler.ValidationMetricsAPI(w, req)

	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("Expected Content-Type application/json, got %s", contentType)
	}

	// Basic check that we got JSON response
	body := w.Body.String()
	if !strings.Contains(body, "total_validations") {
		t.Error("Response should contain 'total_validations' field")
	}
	if !strings.Contains(body, "success_rate") {
		t.Error("Response should contain 'success_rate' field")
	}
	if !strings.Contains(body, "top_error_types") {
		t.Error("Response should contain 'top_error_types' field")
	}
}

// TestSecurityIncidentDetection tests the security incident detection logic
func TestSecurityIncidentDetection(t *testing.T) {
	tests := []struct {
		name           string
		formData       map[string]string
		expectIncident bool
		expectedReason string
	}{
		{
			name: "SQL injection detected",
			formData: map[string]string{
				"discord_id": "'; DROP TABLE users; --",
			},
			expectIncident: true,
			expectedReason: "SQL injection attempt detected",
		},
		{
			name: "XSS attempt detected",
			formData: map[string]string{
				"url": "<script>alert('xss')</script>",
			},
			expectIncident: true,
			expectedReason: "XSS attempt detected",
		},
		{
			name: "Buffer overflow attempt detected",
			formData: map[string]string{
				"discord_id": strings.Repeat("1", 1001), // Over 1000 chars
			},
			expectIncident: true,
			expectedReason: "Buffer overflow attempt detected",
		},
		{
			name: "Normal data should not trigger incident",
			formData: map[string]string{
				"discord_id": "123456789012345678",
				"url":        "https://rocketleague.tracker.network/profile/123",
			},
			expectIncident: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with form data
			formValues := url.Values{}
			for key, value := range tt.formData {
				formValues.Set(key, value)
			}

			req := httptest.NewRequest("POST", "/test", strings.NewReader(formValues.Encode()))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
			if err := req.ParseForm(); err != nil {
				t.Fatalf("Failed to parse form: %v", err)
			}

			// Mock validation result (doesn't matter for this test)
			validation := &ValidationResult{IsValid: false}

			reason := detectSecurityIncident(req, validation)

			if tt.expectIncident {
				if reason == nil {
					t.Errorf("Expected security incident to be detected, but got nil")
				} else if *reason != tt.expectedReason {
					t.Errorf("Expected reason '%s', got '%s'", tt.expectedReason, *reason)
				}
			} else {
				if reason != nil {
					t.Errorf("Expected no security incident, but got: %s", *reason)
				}
			}
		})
	}
}
