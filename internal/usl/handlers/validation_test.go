package handlers

import (
	"testing"
	"usl-server/internal/usl"
)

func TestValidateTracker(t *testing.T) {
	// Create a handler instance (we don't need real dependencies for validation tests)
	handler := &MigrationHandler{}

	tests := []struct {
		name        string
		tracker     *usl.USLUserTracker
		expectValid bool
		expectError string
	}{
		{
			name: "Valid tracker",
			tracker: &usl.USLUserTracker{
				DiscordID:               "123456789012345678", // 18 digits
				URL:                     "https://rocketleague.tracker.network/profile/123",
				OnesCurrentSeasonPeak:   1500,
				TwosCurrentSeasonPeak:   1600,
				ThreesCurrentSeasonPeak: 1700,
				Valid:                   true,
			},
			expectValid: true,
		},
		{
			name: "Empty Discord ID",
			tracker: &usl.USLUserTracker{
				DiscordID: "",
				URL:       "https://rocketleague.tracker.network/profile/123",
			},
			expectValid: false,
			expectError: "Discord ID is required",
		},
		{
			name: "Invalid Discord ID - too short",
			tracker: &usl.USLUserTracker{
				DiscordID: "12345", // Only 5 digits
				URL:       "https://rocketleague.tracker.network/profile/123",
			},
			expectValid: false,
			expectError: "Discord ID must be 17-19 digits",
		},
		{
			name: "Invalid Discord ID - too long",
			tracker: &usl.USLUserTracker{
				DiscordID: "12345678901234567890", // 20 digits
				URL:       "https://rocketleague.tracker.network/profile/123",
			},
			expectValid: false,
			expectError: "Discord ID must be 17-19 digits",
		},
		{
			name: "Invalid Discord ID - non-numeric",
			tracker: &usl.USLUserTracker{
				DiscordID: "1234567890123456ab", // Contains letters
				URL:       "https://rocketleague.tracker.network/profile/123",
			},
			expectValid: false,
			expectError: "Discord ID must be 17-19 digits",
		},
		{
			name: "Invalid URL - wrong domain",
			tracker: &usl.USLUserTracker{
				DiscordID: "123456789012345678",
				URL:       "https://example.com/profile/123",
			},
			expectValid: false,
			expectError: "Invalid tracker URL format",
		},
		{
			name: "Valid URL - ballchasing",
			tracker: &usl.USLUserTracker{
				DiscordID:             "123456789012345678",
				URL:                   "https://ballchasing.com/player/123",
				OnesCurrentSeasonPeak: 1000, // Need at least one playlist data
			},
			expectValid: true,
		},
		{
			name: "Valid URL - rltracker.pro",
			tracker: &usl.USLUserTracker{
				DiscordID:             "123456789012345678",
				URL:                   "https://rltracker.pro/player/123",
				TwosCurrentSeasonPeak: 1000, // Need at least one playlist data
			},
			expectValid: true,
		},
		{
			name: "Invalid MMR - too high",
			tracker: &usl.USLUserTracker{
				DiscordID:             "123456789012345678",
				URL:                   "https://rocketleague.tracker.network/profile/123",
				OnesCurrentSeasonPeak: 5000, // Way too high
			},
			expectValid: false,
			expectError: "1v1 current season MMR must be between 0 and 3000",
		},
		{
			name: "Invalid MMR - negative (treated as no data)",
			tracker: &usl.USLUserTracker{
				DiscordID:             "123456789012345678",
				URL:                   "https://rocketleague.tracker.network/profile/123",
				OnesCurrentSeasonPeak: -100, // Negative - treated as no data
			},
			expectValid: false,
			expectError: "Tracker must have data for at least one playlist",
		},
		{
			name: "Invalid games - too high",
			tracker: &usl.USLUserTracker{
				DiscordID:                    "123456789012345678",
				URL:                          "https://rocketleague.tracker.network/profile/123",
				OnesCurrentSeasonPeak:        1000,  // Need at least one playlist data
				OnesCurrentSeasonGamesPlayed: 15000, // Too many games
			},
			expectValid: false,
			expectError: "1v1 current season games must be between 0 and 10000",
		},
		{
			name: "Invalid case - no playlist data",
			tracker: &usl.USLUserTracker{
				DiscordID:               "123456789012345678",
				URL:                     "https://rocketleague.tracker.network/profile/123",
				OnesCurrentSeasonPeak:   0,
				TwosCurrentSeasonPeak:   0,
				ThreesCurrentSeasonPeak: 0,
			},
			expectValid: false,
			expectError: "Tracker must have data for at least one playlist",
		},
		{
			name: "Invalid case - URL required",
			tracker: &usl.USLUserTracker{
				DiscordID:             "123456789012345678",
				URL:                   "", // Empty URL
				OnesCurrentSeasonPeak: 1000,
			},
			expectValid: false,
			expectError: "Tracker URL is required",
		},
		{
			name: "Valid edge case - max values",
			tracker: &usl.USLUserTracker{
				DiscordID:                      "123456789012345678",
				URL:                            "https://rocketleague.tracker.network/profile/123",
				OnesCurrentSeasonPeak:          3000,
				OnesCurrentSeasonGamesPlayed:   10000,
				TwosCurrentSeasonPeak:          3000,
				TwosCurrentSeasonGamesPlayed:   10000,
				ThreesCurrentSeasonPeak:        3000,
				ThreesCurrentSeasonGamesPlayed: 10000,
			},
			expectValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.validateTracker(tt.tracker)

			if result.IsValid != tt.expectValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectValid, result.IsValid)
			}

			if !tt.expectValid && tt.expectError != "" {
				// Check if any error contains the expected message
				found := false
				for _, err := range result.Errors {
					if err.Message == tt.expectError {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected error message '%s', but got errors: %+v", tt.expectError, result.Errors)
				}
			}

			if tt.expectValid && len(result.Errors) > 0 {
				t.Errorf("Expected no errors for valid tracker, but got: %+v", result.Errors)
			}
		})
	}
}

// Test individual validation functions
func TestIsValidDiscordID(t *testing.T) {
	tests := []struct {
		discordID string
		expected  bool
	}{
		{"123456789012345678", true},    // 18 digits (valid)
		{"12345678901234567", true},     // 17 digits (valid)
		{"1234567890123456789", true},   // 19 digits (valid)
		{"1234567890123456", false},     // 16 digits (too short)
		{"12345678901234567890", false}, // 20 digits (too long)
		{"1234567890123456ab", false},   // Contains letters
		{"", false},                     // Empty
		{"abcdefghijklmnopqr", false},   // All letters
	}

	for _, tt := range tests {
		t.Run(tt.discordID, func(t *testing.T) {
			result := isValidDiscordID(tt.discordID)
			if result != tt.expected {
				t.Errorf("isValidDiscordID(%q) = %v; expected %v", tt.discordID, result, tt.expected)
			}
		})
	}
}

func TestIsValidTrackerURL(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"https://rocketleague.tracker.network/profile/123", true},
		{"https://www.rocketleague.tracker.network/profile/123", true},
		{"http://rocketleague.tracker.network/profile/123", true},
		{"https://ballchasing.com/player/123", true},
		{"https://rltracker.pro/player/123", true},
		{"https://example.com/profile/123", false},
		{"not-a-url", false},
		{"", false},
		{"https://google.com", false},
		{"ftp://rocketleague.tracker.network/profile/123", true}, // Still contains valid host
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := isValidTrackerURL(tt.url)
			if result != tt.expected {
				t.Errorf("isValidTrackerURL(%q) = %v; expected %v", tt.url, result, tt.expected)
			}
		})
	}
}
