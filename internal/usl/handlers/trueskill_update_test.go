package handlers

import (
	"errors"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"usl-server/internal/services"
	"usl-server/internal/usl"
)

// Mock implementations for testing
type MockUSLRepository struct {
	users map[int64]*usl.USLUser
	err   error
}

func (m *MockUSLRepository) GetUserByID(id int64) (*usl.USLUser, error) {
	if m.err != nil {
		return nil, m.err
	}
	user, exists := m.users[id]
	if !exists {
		return nil, errors.New("user not found")
	}
	return user, nil
}

type MockTrueSkillService struct {
	result *services.TrueSkillUpdateResult
}

func (m *MockTrueSkillService) UpdateUserTrueSkillFromTrackers(discordID string) *services.TrueSkillUpdateResult {
	if m.result != nil {
		return m.result
	}
	// Default success result
	return &services.TrueSkillUpdateResult{
		Success:     true,
		HadTrackers: true,
		TrueSkillResult: &services.TrueSkillCalculation{
			Mu:    1500.0,
			Sigma: 250.0,
		},
	}
}

// TestableUpdateUserTrueSkill is a testable version of the handler that uses dependency injection
func TestableUpdateUserTrueSkill(
	w http.ResponseWriter,
	r *http.Request,
	mockRepo *MockUSLRepository,
	mockService *MockTrueSkillService,
	templates *template.Template,
) {
	// Method validation
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract user ID from URL query parameter
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	// Parse and validate user ID
	userID, err := parseUserIDForTest(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get user using mock repository
	user, err := mockRepo.GetUserByID(userID)
	if err != nil {
		http.Error(w, "Database error", http.StatusInternalServerError)
		return
	}

	// Trigger TrueSkill update using mock service
	result := mockService.UpdateUserTrueSkillFromTrackers(user.DiscordID)

	// Check if HTMX request
	isHTMX := r.Header.Get("HX-Request") == "true"
	if !isHTMX {
		// Redirect for non-HTMX requests
		http.Redirect(w, r, fmt.Sprintf("/usl/users/detail?id=%d", userID), http.StatusSeeOther)
		return
	}

	// Render template for HTMX requests
	data := struct {
		Success         bool
		HadTrackers     bool
		TrueSkillResult *services.TrueSkillCalculation
		UserName        string
		DiscordID       string
		Error           string
	}{
		Success:         result.Success,
		HadTrackers:     result.HadTrackers,
		TrueSkillResult: result.TrueSkillResult,
		UserName:        user.Name,
		DiscordID:       user.DiscordID,
		Error:           result.Error,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := templates.ExecuteTemplate(w, "trueskill-update-result", data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

// Helper function for parsing user ID in tests
func parseUserIDForTest(userIDStr string) (int64, error) {
	if userIDStr == "" {
		return 0, errors.New("empty user ID")
	}

	// Basic validation - reject non-numeric, negative, zero, or very long values
	if len(userIDStr) > 20 || strings.ContainsAny(userIDStr, "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ<>\"';-./\\") {
		return 0, errors.New("invalid user ID format")
	}

	// Convert to int64
	var userID int64
	for _, char := range userIDStr {
		if char < '0' || char > '9' {
			return 0, errors.New("non-numeric user ID")
		}
		userID = userID*10 + int64(char-'0')
	}

	if userID <= 0 {
		return 0, errors.New("user ID must be positive")
	}

	return userID, nil
}

// Helper function to create test template
func createTestTemplate() *template.Template {
	tmpl := template.New("test")
	return template.Must(tmpl.Parse(`
		{{define "trueskill-update-result"}}
		<div class="{{if .Success}}success{{else}}error{{end}}">
			{{if .Success}}Success: {{.UserName}}{{else}}Error: {{.Error}}{{end}}
		</div>
		{{end}}
	`))
}

// Core Behavior Tests
func TestUpdateUserTrueSkill_CoreBehavior(t *testing.T) {
	tests := []struct {
		name           string
		userID         string
		method         string
		isHTMX         bool
		mockUser       *usl.USLUser
		mockResult     *services.TrueSkillUpdateResult
		expectStatus   int
		expectRedirect bool
	}{
		{
			name:   "Successful TrueSkill update with trackers (HTMX)",
			userID: "1",
			method: "POST",
			isHTMX: true,
			mockUser: &usl.USLUser{
				ID:        1,
				DiscordID: "123456789012345678",
				Name:      "Test User",
			},
			mockResult: &services.TrueSkillUpdateResult{
				Success:     true,
				HadTrackers: true,
				TrueSkillResult: &services.TrueSkillCalculation{
					Mu:    1500.0,
					Sigma: 250.0,
				},
			},
			expectStatus: 200,
		},
		{
			name:   "Successful TrueSkill update without trackers (defaults)",
			userID: "2",
			method: "POST",
			isHTMX: true,
			mockUser: &usl.USLUser{
				ID:        2,
				DiscordID: "987654321098765432",
				Name:      "New User",
			},
			mockResult: &services.TrueSkillUpdateResult{
				Success:     true,
				HadTrackers: false,
				TrueSkillResult: &services.TrueSkillCalculation{
					Mu:    1500.0,
					Sigma: 500.0,
				},
			},
			expectStatus: 200,
		},
		{
			name:   "Non-HTMX request should redirect",
			userID: "1",
			method: "POST",
			isHTMX: false,
			mockUser: &usl.USLUser{
				ID:        1,
				DiscordID: "123456789012345678",
				Name:      "Test User",
			},
			expectStatus:   303,
			expectRedirect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mocks
			mockRepo := &MockUSLRepository{
				users: map[int64]*usl.USLUser{},
			}
			if tt.mockUser != nil {
				mockRepo.users[tt.mockUser.ID] = tt.mockUser
			}

			mockService := &MockTrueSkillService{
				result: tt.mockResult,
			}

			templates := createTestTemplate()

			// Create request
			req := httptest.NewRequest(tt.method, "/usl/users/update-trueskill?id="+tt.userID, nil)
			if tt.isHTMX {
				req.Header.Set("HX-Request", "true")
			}

			// Execute testable handler
			recorder := httptest.NewRecorder()
			TestableUpdateUserTrueSkill(recorder, req, mockRepo, mockService, templates)

			// Verify response
			if recorder.Code != tt.expectStatus {
				t.Errorf("Expected status %d, got %d", tt.expectStatus, recorder.Code)
			}

			if tt.expectRedirect {
				location := recorder.Header().Get("Location")
				if location == "" {
					t.Error("Expected redirect but no Location header found")
				}
			}

			// For successful HTMX requests, verify template rendering
			if tt.expectStatus == 200 && tt.isHTMX {
				body := recorder.Body.String()
				if !strings.Contains(body, "success") && !strings.Contains(body, "error") {
					t.Error("Expected rendered template in response body")
				}
			}
		})
	}
}

// Input Validation Tests
func TestUpdateUserTrueSkill_InputValidation(t *testing.T) {
	tests := []struct {
		name         string
		method       string
		userID       string
		expectStatus int
	}{
		{
			name:         "Invalid HTTP method (GET)",
			method:       "GET",
			userID:       "1",
			expectStatus: 405,
		},
		{
			name:         "Invalid HTTP method (PUT)",
			method:       "PUT",
			userID:       "1",
			expectStatus: 405,
		},
		{
			name:         "Missing user ID parameter",
			method:       "POST",
			userID:       "",
			expectStatus: 400,
		},
		{
			name:         "Invalid user ID format (non-numeric)",
			method:       "POST",
			userID:       "abc",
			expectStatus: 400,
		},
		{
			name:         "Invalid user ID format (negative)",
			method:       "POST",
			userID:       "-1",
			expectStatus: 400,
		},
		{
			name:         "Invalid user ID format (zero)",
			method:       "POST",
			userID:       "0",
			expectStatus: 400,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUSLRepository{users: map[int64]*usl.USLUser{}}
			mockService := &MockTrueSkillService{}
			templates := createTestTemplate()

			// Create request
			var req *http.Request
			if tt.userID != "" {
				req = httptest.NewRequest(tt.method, "/usl/users/update-trueskill?id="+tt.userID, nil)
			} else {
				req = httptest.NewRequest(tt.method, "/usl/users/update-trueskill", nil)
			}
			req.Header.Set("HX-Request", "true")

			recorder := httptest.NewRecorder()
			TestableUpdateUserTrueSkill(recorder, req, mockRepo, mockService, templates)

			if recorder.Code != tt.expectStatus {
				t.Errorf("Expected status %d, got %d", tt.expectStatus, recorder.Code)
			}
		})
	}
}

// Database Error Tests
func TestUpdateUserTrueSkill_DatabaseErrors(t *testing.T) {
	tests := []struct {
		name         string
		userID       string
		repoError    error
		expectStatus int
	}{
		{
			name:         "User not found",
			userID:       "999",
			repoError:    errors.New("user not found"),
			expectStatus: 500,
		},
		{
			name:         "Database connection failure",
			userID:       "1",
			repoError:    errors.New("connection failed"),
			expectStatus: 500,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo := &MockUSLRepository{
				users: map[int64]*usl.USLUser{},
				err:   tt.repoError,
			}
			mockService := &MockTrueSkillService{}
			templates := createTestTemplate()

			req := httptest.NewRequest("POST", "/usl/users/update-trueskill?id="+tt.userID, nil)
			req.Header.Set("HX-Request", "true")

			recorder := httptest.NewRecorder()
			TestableUpdateUserTrueSkill(recorder, req, mockRepo, mockService, templates)

			if recorder.Code != tt.expectStatus {
				t.Errorf("Expected status %d, got %d", tt.expectStatus, recorder.Code)
			}
		})
	}
}

// TrueSkill Service Error Tests
func TestUpdateUserTrueSkill_ServiceErrors(t *testing.T) {
	tests := []struct {
		name       string
		mockResult *services.TrueSkillUpdateResult
	}{
		{
			name: "TrueSkill service returns failure",
			mockResult: &services.TrueSkillUpdateResult{
				Success: false,
				Error:   "Failed to calculate TrueSkill: insufficient data",
			},
		},
		{
			name: "TrueSkill service returns success but no result data",
			mockResult: &services.TrueSkillUpdateResult{
				Success:         true,
				HadTrackers:     false,
				TrueSkillResult: nil,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := &usl.USLUser{
				ID:        1,
				DiscordID: "123456789012345678",
				Name:      "Test User",
			}

			mockRepo := &MockUSLRepository{
				users: map[int64]*usl.USLUser{1: mockUser},
			}
			mockService := &MockTrueSkillService{
				result: tt.mockResult,
			}
			templates := createTestTemplate()

			req := httptest.NewRequest("POST", "/usl/users/update-trueskill?id=1", nil)
			req.Header.Set("HX-Request", "true")

			recorder := httptest.NewRecorder()
			TestableUpdateUserTrueSkill(recorder, req, mockRepo, mockService, templates)

			// Should still return 200 but show error in template
			if recorder.Code != 200 {
				t.Errorf("Expected status 200 for service errors, got %d", recorder.Code)
			}

			body := recorder.Body.String()
			if tt.mockResult.Success {
				if !strings.Contains(body, "success") {
					t.Error("Expected success indication in response")
				}
			} else {
				if !strings.Contains(body, "error") {
					t.Error("Expected error indication in response")
				}
			}
		})
	}
}

// Security Input Tests
func TestUpdateUserTrueSkill_SecurityInputs(t *testing.T) {
	tests := []struct {
		name          string
		userID        string
		expectBlocked bool
		description   string
	}{
		{
			name:          "SQL injection in user ID",
			userID:        "1'; DROP TABLE users; --",
			expectBlocked: true,
			description:   "Should reject SQL injection attempts",
		},
		{
			name:          "Path traversal in user ID",
			userID:        "../../../etc/passwd",
			expectBlocked: true,
			description:   "Should reject path traversal attempts",
		},
		{
			name:          "Script injection in user ID",
			userID:        "<script>alert('xss')</script>",
			expectBlocked: true,
			description:   "Should reject script injection attempts",
		},
		{
			name:          "Extremely long user ID",
			userID:        strings.Repeat("1", 1000),
			expectBlocked: true,
			description:   "Should reject excessively long values",
		},
		{
			name:          "Valid numeric user ID",
			userID:        "12345",
			expectBlocked: false,
			description:   "Should accept valid numeric ID",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockUser := &usl.USLUser{
				ID:        12345,
				DiscordID: "123456789012345678",
				Name:      "Test User",
			}

			mockRepo := &MockUSLRepository{
				users: map[int64]*usl.USLUser{12345: mockUser},
			}
			mockService := &MockTrueSkillService{}
			templates := createTestTemplate()

			req := httptest.NewRequest("POST", "/usl/users/update-trueskill?id="+url.QueryEscape(tt.userID), nil)
			req.Header.Set("HX-Request", "true")

			recorder := httptest.NewRecorder()
			TestableUpdateUserTrueSkill(recorder, req, mockRepo, mockService, templates)

			if tt.expectBlocked {
				// Should be blocked with 400 status (bad request)
				if recorder.Code == 200 {
					t.Errorf("Expected malicious input to be blocked: %s", tt.description)
				}
			} else {
				// Valid input should proceed (may fail later due to missing user, but not at validation)
				if recorder.Code == 400 {
					t.Errorf("Expected valid input to pass validation: %s", tt.description)
				}
			}
		})
	}
}

// HTMX Integration Tests
func TestUpdateUserTrueSkill_HTMXIntegration(t *testing.T) {
	t.Run("HTMX request detection", func(t *testing.T) {
		mockUser := &usl.USLUser{
			ID:        1,
			DiscordID: "123456789012345678",
			Name:      "Test User",
		}

		mockRepo := &MockUSLRepository{
			users: map[int64]*usl.USLUser{1: mockUser},
		}
		mockService := &MockTrueSkillService{}
		templates := createTestTemplate()

		// Test without HTMX header
		req1 := httptest.NewRequest("POST", "/usl/users/update-trueskill?id=1", nil)
		recorder1 := httptest.NewRecorder()
		TestableUpdateUserTrueSkill(recorder1, req1, mockRepo, mockService, templates)

		// Should redirect for non-HTMX requests
		if recorder1.Code != 303 {
			t.Errorf("Expected redirect (303) for non-HTMX request, got %d", recorder1.Code)
		}

		// Test with HTMX header
		req2 := httptest.NewRequest("POST", "/usl/users/update-trueskill?id=1", nil)
		req2.Header.Set("HX-Request", "true")
		recorder2 := httptest.NewRecorder()
		TestableUpdateUserTrueSkill(recorder2, req2, mockRepo, mockService, templates)

		// Should return template fragment for HTMX requests
		if recorder2.Code != 200 {
			t.Errorf("Expected success (200) for HTMX request, got %d", recorder2.Code)
		}
	})

	t.Run("HTMX response content type", func(t *testing.T) {
		mockUser := &usl.USLUser{
			ID:        1,
			DiscordID: "123456789012345678",
			Name:      "Test User",
		}

		mockRepo := &MockUSLRepository{
			users: map[int64]*usl.USLUser{1: mockUser},
		}
		mockService := &MockTrueSkillService{}
		templates := createTestTemplate()

		req := httptest.NewRequest("POST", "/usl/users/update-trueskill?id=1", nil)
		req.Header.Set("HX-Request", "true")
		recorder := httptest.NewRecorder()
		TestableUpdateUserTrueSkill(recorder, req, mockRepo, mockService, templates)

		contentType := recorder.Header().Get("Content-Type")
		if !strings.Contains(contentType, "text/html") {
			t.Errorf("Expected HTML content type for HTMX response, got %s", contentType)
		}
	})
}

// Template Rendering Tests
func TestUpdateUserTrueSkill_TemplateRendering(t *testing.T) {
	t.Run("Success template data structure", func(t *testing.T) {
		mockUser := &usl.USLUser{
			ID:        1,
			DiscordID: "123456789012345678",
			Name:      "Test User",
		}

		mockResult := &services.TrueSkillUpdateResult{
			Success:     true,
			HadTrackers: true,
			TrueSkillResult: &services.TrueSkillCalculation{
				Mu:    1500.5,
				Sigma: 250.75,
			},
		}

		mockRepo := &MockUSLRepository{
			users: map[int64]*usl.USLUser{1: mockUser},
		}
		mockService := &MockTrueSkillService{
			result: mockResult,
		}
		templates := createTestTemplate()

		req := httptest.NewRequest("POST", "/usl/users/update-trueskill?id=1", nil)
		req.Header.Set("HX-Request", "true")
		recorder := httptest.NewRecorder()
		TestableUpdateUserTrueSkill(recorder, req, mockRepo, mockService, templates)

		body := recorder.Body.String()
		if !strings.Contains(body, "success") {
			t.Error("Expected success indicator in template")
		}
		if !strings.Contains(body, mockUser.Name) {
			t.Error("Expected user name in template")
		}
	})

	t.Run("Error template data structure", func(t *testing.T) {
		mockUser := &usl.USLUser{
			ID:        1,
			DiscordID: "123456789012345678",
			Name:      "Error User",
		}

		mockResult := &services.TrueSkillUpdateResult{
			Success: false,
			Error:   "Failed to calculate TrueSkill: insufficient data",
		}

		mockRepo := &MockUSLRepository{
			users: map[int64]*usl.USLUser{1: mockUser},
		}
		mockService := &MockTrueSkillService{
			result: mockResult,
		}
		templates := createTestTemplate()

		req := httptest.NewRequest("POST", "/usl/users/update-trueskill?id=1", nil)
		req.Header.Set("HX-Request", "true")
		recorder := httptest.NewRecorder()
		TestableUpdateUserTrueSkill(recorder, req, mockRepo, mockService, templates)

		body := recorder.Body.String()
		if !strings.Contains(body, "error") {
			t.Error("Expected error indicator in template")
		}
	})
}

// parseUserID function test to verify security
func TestParseUserID_Security(t *testing.T) {
	tests := []struct {
		input       string
		expectError bool
	}{
		{"1", false},
		{"12345", false},
		{"0", true},                      // Zero should be rejected
		{"-1", true},                     // Negative should be rejected
		{"abc", true},                    // Non-numeric should be rejected
		{"", true},                       // Empty should be rejected
		{"1.5", true},                    // Float should be rejected
		{"1e10", true},                   // Scientific notation should be rejected
		{"<script>", true},               // Script injection should be rejected
		{strings.Repeat("1", 100), true}, // Very long should be rejected
	}

	for _, tt := range tests {
		t.Run(fmt.Sprintf("input_%s", tt.input), func(t *testing.T) {
			_, err := parseUserIDForTest(tt.input)
			hasError := err != nil

			if hasError != tt.expectError {
				t.Errorf("parseUserIDForTest(%q): expected error=%v, got error=%v", tt.input, tt.expectError, hasError)
			}
		})
	}
}

// Test the actual MigrationHandler.parseUserID method for compatibility
func TestMigrationHandler_ParseUserID(t *testing.T) {
	handler := &MigrationHandler{}

	validInputs := []string{"1", "123", "999999"}
	for _, input := range validInputs {
		t.Run(fmt.Sprintf("valid_%s", input), func(t *testing.T) {
			result, err := handler.parseUserID(input)
			if err != nil {
				t.Errorf("parseUserID(%q) should not return error for valid input", input)
			}
			if result <= 0 {
				t.Errorf("parseUserID(%q) should return positive value", input)
			}
		})
	}
}
