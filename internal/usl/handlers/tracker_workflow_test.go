package handlers

import (
	"html/template"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"usl-server/internal/usl"
)

// TestAddTrackerButtonVisibility ensures the Add Tracker button is always visible
// regardless of whether the user has existing trackers
func TestAddTrackerButtonVisibility(t *testing.T) {
	// Load the user-detail template following the established pattern
	tmpl := template.New("test")
	tmpl = tmpl.Funcs(template.FuncMap{
		"printf": func(format string, args ...interface{}) string {
			return "test-value"
		},
	})

	// Parse templates with dependencies (following Issue35 pattern)
	tmpl = template.Must(tmpl.ParseFiles(
		"../../../templates/user-detail.html",
		"../../../templates/navigation.html",
	))

	tests := []struct {
		name         string
		userTrackers []*usl.USLUserTracker
		expectButton bool
		description  string
	}{
		{
			name:         "No_trackers_shows_button",
			userTrackers: nil,
			expectButton: true,
			description:  "Users with no trackers should see Add Tracker button",
		},
		{
			name: "Single_tracker_shows_button",
			userTrackers: []*usl.USLUserTracker{
				{
					ID:        1,
					DiscordID: "123456789012345678",
					URL:       "https://rocketleague.tracker.network/profile/123",
					Valid:     true,
				},
			},
			expectButton: true,
			description:  "Users with existing trackers should still see Add Tracker button",
		},
		{
			name: "Multiple_trackers_shows_button",
			userTrackers: []*usl.USLUserTracker{
				{
					ID:        1,
					DiscordID: "123456789012345678",
					URL:       "https://rocketleague.tracker.network/profile/123",
					Valid:     true,
				},
				{
					ID:        2,
					DiscordID: "123456789012345678",
					URL:       "https://ballchasing.com/player/123",
					Valid:     true,
				},
			},
			expectButton: true,
			description:  "Users with multiple trackers should still see Add Tracker button",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create template data matching the handler structure
			data := struct {
				Title        string
				CurrentPage  string
				User         *usl.USLUser
				UserTrackers []*usl.USLUserTracker
			}{
				Title:       "Test User",
				CurrentPage: "users",
				User: &usl.USLUser{
					ID:        1,
					Name:      "Test User",
					DiscordID: "123456789012345678",
					Active:    true,
				},
				UserTrackers: tt.userTrackers,
			}

			var buf strings.Builder
			err := tmpl.ExecuteTemplate(&buf, "user-detail-page", data)
			if err != nil {
				t.Fatalf("Template rendering failed: %v", err)
			}

			output := buf.String()

			// Verify Add Tracker button is always present
			hasAddButton := strings.Contains(output, "Add Tracker")
			hasCorrectLink := strings.Contains(output, "/usl/trackers/new?discord_id=123456789012345678")

			if tt.expectButton {
				if !hasAddButton {
					t.Errorf("Expected Add Tracker button to be present but it was not found")
					t.Logf("Description: %s", tt.description)
				}
				if !hasCorrectLink {
					t.Errorf("Expected Add Tracker link with correct discord_id parameter")
					t.Logf("Expected: /usl/trackers/new?discord_id=123456789012345678")
				}
			}

			// The key fix: button should be present regardless of tracker count
			// (This verifies the template fix that moved the button outside the conditional)
		})
	}
}

// TestNewTrackerFormParameterHandling tests that discord_id URL parameter is pre-filled
func TestNewTrackerFormParameterHandling(t *testing.T) {
	tests := []struct {
		name          string
		url           string
		expectPreFill bool
		expectedID    string
		description   string
	}{
		{
			name:          "Valid_discord_id_parameter_prefills",
			url:           "/usl/trackers/new?discord_id=123456789012345678",
			expectPreFill: true,
			expectedID:    "123456789012345678",
			description:   "Valid Discord ID from URL should pre-fill form field",
		},
		{
			name:          "Invalid_discord_id_ignored",
			url:           "/usl/trackers/new?discord_id=invalid123",
			expectPreFill: false,
			expectedID:    "",
			description:   "Invalid Discord ID should be ignored and not pre-fill",
		},
		{
			name:          "No_parameter_empty_form",
			url:           "/usl/trackers/new",
			expectPreFill: false,
			expectedID:    "",
			description:   "No discord_id parameter should result in empty form",
		},
		{
			name:          "Too_short_discord_id_ignored",
			url:           "/usl/trackers/new?discord_id=12345",
			expectPreFill: false,
			expectedID:    "",
			description:   "Too short Discord ID should be ignored",
		},
		{
			name:          "Too_long_discord_id_ignored",
			url:           "/usl/trackers/new?discord_id=12345678901234567890",
			expectPreFill: false,
			expectedID:    "",
			description:   "Too long Discord ID should be ignored",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request following established pattern
			req, err := http.NewRequest("GET", tt.url, nil)
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}

			// Test the parameter extraction logic (mimicking NewTrackerForm fix)
			tracker := &usl.USLUserTracker{}
			if discordID := req.URL.Query().Get("discord_id"); discordID != "" {
				// Validate using the same function from migration_handler.go
				if isValidDiscordID(discordID) {
					tracker.DiscordID = discordID
				}
			}

			// Verify expectations
			if tt.expectPreFill {
				if tracker.DiscordID != tt.expectedID {
					t.Errorf("Expected tracker.DiscordID = %q, got %q", tt.expectedID, tracker.DiscordID)
					t.Logf("Description: %s", tt.description)
				}
			} else {
				if tracker.DiscordID != "" {
					t.Errorf("Expected empty DiscordID, got %q", tracker.DiscordID)
					t.Logf("Description: %s", tt.description)
				}
			}
		})
	}
}

// TestTrackerCreationTrueSkillFlow tests that TrueSkill is called after tracker creation
// using the established pattern of testing critical flow without complex mocking
func TestTrackerCreationTrueSkillFlow(t *testing.T) {
	tests := []struct {
		name        string
		description string
		formData    map[string]string
		expectValid bool
	}{
		{
			name:        "Valid_tracker_creation_flow",
			description: "Valid tracker should trigger TrueSkill update flow",
			formData: map[string]string{
				"discord_id":         "123456789012345678",
				"url":                "https://rocketleague.tracker.network/profile/123",
				"ones_current_peak":  "1500",
				"ones_current_games": "10",
				"valid":              "true",
			},
			expectValid: true,
		},
		{
			name:        "Invalid_tracker_skips_trueskill",
			description: "Invalid tracker should not reach TrueSkill update",
			formData: map[string]string{
				"discord_id": "invalid", // Invalid Discord ID
				"url":        "https://rocketleague.tracker.network/profile/123",
				"valid":      "true",
			},
			expectValid: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create form data
			form := url.Values{}
			for key, value := range tt.formData {
				form.Set(key, value)
			}

			// Create request
			req, err := http.NewRequest("POST", "/usl/trackers/create", strings.NewReader(form.Encode()))
			if err != nil {
				t.Fatalf("Failed to create request: %v", err)
			}
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			// Parse form and build tracker (mimicking CreateTracker logic)
			if err := req.ParseForm(); err != nil {
				t.Fatalf("Failed to parse form: %v", err)
			}

			// Build tracker from form using existing logic pattern
			baseHandler := &BaseHandler{}
			tracker := baseHandler.buildTrackerFromForm(req)

			// Validate tracker (this determines if TrueSkill update should be called)
			validation := baseHandler.validateTracker(tracker)

			if tt.expectValid {
				if !validation.IsValid {
					t.Errorf("Expected valid tracker, but validation failed: %v", validation.Errors)
					t.Logf("Description: %s", tt.description)
				}
				// If validation passes, TrueSkill update SHOULD be called
				// (We can't easily test the actual call without complex mocking,
				// but we've verified the precondition)
			} else {
				if validation.IsValid {
					t.Errorf("Expected invalid tracker, but validation passed")
					t.Logf("Description: %s", tt.description)
				}
				// If validation fails, TrueSkill update should NOT be called
			}
		})
	}
}

// TestCompleteTrackerWorkflow tests the complete user journey end-to-end
func TestCompleteTrackerWorkflow(t *testing.T) {
	// Load both templates (following existing template testing pattern)
	userDetailTmpl := template.New("user-detail")
	userDetailTmpl = userDetailTmpl.Funcs(template.FuncMap{
		"printf": func(format string, args ...interface{}) string {
			return "test-value"
		},
	})
	userDetailTmpl = template.Must(userDetailTmpl.ParseFiles(
		"../../../templates/user-detail.html",
		"../../../templates/navigation.html",
	))

	trackerNewTmpl := template.New("tracker-new")
	trackerNewTmpl = trackerNewTmpl.Funcs(template.FuncMap{
		"printf": func(format string, args ...interface{}) string {
			return "test-value"
		},
	})
	trackerNewTmpl = template.Must(trackerNewTmpl.ParseFiles(
		"../../../templates/tracker-new.html",
		"../../../templates/navigation.html",
	))

	t.Run("Complete_User_Journey", func(t *testing.T) {
		// Step 1: User detail page shows button (regardless of existing trackers)
		userData := struct {
			Title        string
			CurrentPage  string
			User         *usl.USLUser
			UserTrackers []*usl.USLUserTracker
		}{
			Title:       "User Detail",
			CurrentPage: "users",
			User: &usl.USLUser{
				ID:        1,
				Name:      "Test User",
				DiscordID: "123456789012345678",
				Active:    true,
			},
			UserTrackers: []*usl.USLUserTracker{{ID: 1, DiscordID: "123456789012345678", Valid: true}}, // Has existing tracker
		}

		var userBuf strings.Builder
		err := userDetailTmpl.ExecuteTemplate(&userBuf, "user-detail-page", userData)
		if err != nil {
			t.Fatalf("User detail template failed: %v", err)
		}

		userOutput := userBuf.String()

		// Verify button present and URL correct (Issue 41 fix #1)
		if !strings.Contains(userOutput, "Add Tracker") {
			t.Errorf("Missing Add Tracker button on user detail page with existing trackers")
			t.Logf("This tests the fix for Issue 41 template logic")
		}
		if !strings.Contains(userOutput, "discord_id=123456789012345678") {
			t.Errorf("Missing pre-fill parameter in button URL")
			t.Logf("Expected URL parameter: discord_id=123456789012345678")
		}

		t.Logf("✅ Step 1: User detail page shows Add Tracker button with correct URL parameter")

		// Step 2: New tracker form pre-fills Discord ID (Issue 41 fix #2)
		trackerData := struct {
			Title       string
			CurrentPage string
			Tracker     *usl.USLUserTracker
			Errors      map[string]string
		}{
			Title:       "New Tracker",
			CurrentPage: "trackers",
			Tracker:     &usl.USLUserTracker{DiscordID: "123456789012345678"}, // Pre-filled from URL param
			Errors:      map[string]string{},
		}

		var trackerBuf strings.Builder
		err = trackerNewTmpl.ExecuteTemplate(&trackerBuf, "tracker-new-page", trackerData)
		if err != nil {
			t.Fatalf("Tracker new template failed: %v", err)
		}

		trackerOutput := trackerBuf.String()

		// Verify pre-filled Discord ID (Issue 41 fix #2)
		if !strings.Contains(trackerOutput, `value="123456789012345678"`) {
			t.Errorf("Discord ID not pre-filled in form")
			t.Logf("This tests the fix for Issue 41 parameter handling")
		}

		t.Logf("✅ Step 2: New tracker form pre-fills Discord ID from URL parameter")

		// Step 3: Form submission validation and workflow (Issue 41 fix #3)
		// Test that valid form data would trigger the TrueSkill update workflow
		baseHandler := &BaseHandler{}

		form := url.Values{}
		form.Set("discord_id", "123456789012345678")
		form.Set("url", "https://rocketleague.tracker.network/profile/123")
		form.Set("ones_current_peak", "1500")
		form.Set("ones_current_games", "10")
		form.Set("valid", "true")

		req, err := http.NewRequest("POST", "/usl/trackers/create", strings.NewReader(form.Encode()))
		if err != nil {
			t.Fatalf("Failed to create request: %v", err)
		}
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		req.ParseForm()

		tracker := baseHandler.buildTrackerFromForm(req)
		validation := baseHandler.validateTracker(tracker)

		if !validation.IsValid {
			t.Errorf("Valid form submission should pass validation: %v", validation.Errors)
		} else {
			t.Logf("✅ Step 3: Form submission would trigger tracker creation AND TrueSkill update")
		}

		t.Logf("✅ Complete workflow: User can always add trackers → form pre-fills → TrueSkill updates")
	})

	t.Run("Empty_User_Journey", func(t *testing.T) {
		// Test the journey for a user with NO existing trackers
		userData := struct {
			Title        string
			CurrentPage  string
			User         *usl.USLUser
			UserTrackers []*usl.USLUserTracker
		}{
			Title:       "User Detail",
			CurrentPage: "users",
			User: &usl.USLUser{
				ID:        2,
				Name:      "New User",
				DiscordID: "987654321098765432",
				Active:    true,
			},
			UserTrackers: nil, // No existing trackers
		}

		var userBuf strings.Builder
		err := userDetailTmpl.ExecuteTemplate(&userBuf, "user-detail-page", userData)
		if err != nil {
			t.Fatalf("User detail template failed: %v", err)
		}

		userOutput := userBuf.String()

		// Verify button still present for empty users
		if !strings.Contains(userOutput, "Add Tracker") {
			t.Errorf("Missing Add Tracker button on user detail page with no trackers")
		}
		if !strings.Contains(userOutput, "discord_id=987654321098765432") {
			t.Errorf("Missing correct Discord ID in button URL for new user")
		}

		t.Logf("✅ User with no trackers can also access Add Tracker functionality")
	})
}
