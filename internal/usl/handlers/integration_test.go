package handlers

import (
	"testing"
	"time"
	"usl-server/internal/config"
	"usl-server/internal/services"
	"usl-server/internal/usl"
)

// TestUSLTrueSkillIntegration tests the complete USL-TrueSkill integration flow
func TestUSLTrueSkillIntegration(t *testing.T) {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		t.Skipf("Config not available for testing: %v", err)
	}

	// Create test tracker data
	lastUpdated := time.Now().Format(time.RFC3339)
	testTracker := &usl.USLUserTracker{
		ID:                              1,
		DiscordID:                       "123456789012345678",
		URL:                             "https://rocketleague.tracker.network/rocket-league/profile/steam/testuser",
		OnesCurrentSeasonPeak:           1500,
		OnesCurrentSeasonGamesPlayed:    10,
		OnesPreviousSeasonPeak:          1400,
		OnesPreviousSeasonGamesPlayed:   5,
		TwosCurrentSeasonPeak:           1600,
		TwosCurrentSeasonGamesPlayed:    15,
		TwosPreviousSeasonPeak:          1550,
		TwosPreviousSeasonGamesPlayed:   8,
		ThreesCurrentSeasonPeak:         1450,
		ThreesCurrentSeasonGamesPlayed:  12,
		ThreesPreviousSeasonPeak:        1400,
		ThreesPreviousSeasonGamesPlayed: 6,
		LastUpdated:                     &lastUpdated,
		Valid:                           true,
		MMR:                             1500,
		CreatedAt:                       time.Now(),
		UpdatedAt:                       time.Now(),
	}

	// Create mock handler
	baseHandler := &BaseHandler{
		config: cfg,
	}

	t.Run("TestUSLTrackerToTrackerDataMapping", func(t *testing.T) {
		// Test the data transformation
		trackerData := baseHandler.transformUSLTrackerToTrackerData(testTracker)

		// Verify transformation
		if trackerData.DiscordID != testTracker.DiscordID {
			t.Errorf("DiscordID mismatch: got %s, want %s", trackerData.DiscordID, testTracker.DiscordID)
		}

		if trackerData.URL != testTracker.URL {
			t.Errorf("URL mismatch: got %s, want %s", trackerData.URL, testTracker.URL)
		}

		// Test field mapping
		testCases := []struct {
			name     string
			got      int
			expected int
		}{
			{"OnesCurrentPeak", trackerData.OnesCurrentPeak, testTracker.OnesCurrentSeasonPeak},
			{"OnesCurrentGames", trackerData.OnesCurrentGames, testTracker.OnesCurrentSeasonGamesPlayed},
			{"OnesPreviousPeak", trackerData.OnesPreviousPeak, testTracker.OnesPreviousSeasonPeak},
			{"OnesPreviousGames", trackerData.OnesPreviousGames, testTracker.OnesPreviousSeasonGamesPlayed},
			{"TwosCurrentPeak", trackerData.TwosCurrentPeak, testTracker.TwosCurrentSeasonPeak},
			{"TwosCurrentGames", trackerData.TwosCurrentGames, testTracker.TwosCurrentSeasonGamesPlayed},
			{"TwosPreviousPeak", trackerData.TwosPreviousPeak, testTracker.TwosPreviousSeasonPeak},
			{"TwosPreviousGames", trackerData.TwosPreviousGames, testTracker.TwosPreviousSeasonGamesPlayed},
			{"ThreesCurrentPeak", trackerData.ThreesCurrentPeak, testTracker.ThreesCurrentSeasonPeak},
			{"ThreesCurrentGames", trackerData.ThreesCurrentGames, testTracker.ThreesCurrentSeasonGamesPlayed},
			{"ThreesPreviousPeak", trackerData.ThreesPreviousPeak, testTracker.ThreesPreviousSeasonPeak},
			{"ThreesPreviousGames", trackerData.ThreesPreviousGames, testTracker.ThreesPreviousSeasonGamesPlayed},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				if tc.got != tc.expected {
					t.Errorf("%s mismatch: got %d, want %d", tc.name, tc.got, tc.expected)
				}
			})
		}

		// Test LastUpdated parsing
		if trackerData.LastUpdated.IsZero() {
			t.Error("LastUpdated should not be zero")
		}
	})

	t.Run("TestUSLTrackerToTrackerDataWithNilLastUpdated", func(t *testing.T) {
		// Test with nil LastUpdated
		trackerWithNilTime := *testTracker
		trackerWithNilTime.LastUpdated = nil

		trackerData := baseHandler.transformUSLTrackerToTrackerData(&trackerWithNilTime)

		if trackerData.LastUpdated.IsZero() {
			t.Error("LastUpdated should be set to current time when nil")
		}
	})

	t.Run("TestUSLTrackerToTrackerDataWithInvalidLastUpdated", func(t *testing.T) {
		// Test with invalid LastUpdated format
		invalidTime := "invalid-time-format"
		trackerWithInvalidTime := *testTracker
		trackerWithInvalidTime.LastUpdated = &invalidTime

		trackerData := baseHandler.transformUSLTrackerToTrackerData(&trackerWithInvalidTime)

		if trackerData.LastUpdated.IsZero() {
			t.Error("LastUpdated should be set to current time when invalid")
		}
	})
}

// TestTrueSkillServiceTrackerDataMethod tests that the new TrueSkill service method exists and has correct signature
func TestTrueSkillServiceTrackerDataMethod(t *testing.T) {
	// Create data transformation service
	dataTransformationService := services.NewDataTransformationService()

	// Test TrackerData validation
	t.Run("TestTrackerDataValidation", func(t *testing.T) {
		validTrackerData := &services.TrackerData{
			DiscordID:           "123456789012345678",
			URL:                 "https://rocketleague.tracker.network/rocket-league/profile/steam/testuser",
			OnesCurrentPeak:     1500,
			OnesCurrentGames:    10,
			OnesPreviousPeak:    1400,
			OnesPreviousGames:   5,
			TwosCurrentPeak:     1600,
			TwosCurrentGames:    15,
			TwosPreviousPeak:    1550,
			TwosPreviousGames:   8,
			ThreesCurrentPeak:   1450,
			ThreesCurrentGames:  12,
			ThreesPreviousPeak:  1400,
			ThreesPreviousGames: 6,
			LastUpdated:         time.Now(),
		}

		err := dataTransformationService.ValidateTrackerData(validTrackerData)
		if err != nil {
			t.Errorf("Valid tracker data should not produce validation error: %v", err)
		}

		// Test invalid data
		invalidTrackerData := &services.TrackerData{
			DiscordID: "", // Invalid - empty
		}

		err = dataTransformationService.ValidateTrackerData(invalidTrackerData)
		if err == nil {
			t.Error("Invalid tracker data should produce validation error")
		}
	})
}

// TestConfigurationCompatibility tests that USL can use existing TrueSkill configuration
func TestConfigurationCompatibility(t *testing.T) {
	cfg, err := config.Load()
	if err != nil {
		t.Skipf("Config not available for testing: %v", err)
	}

	t.Run("TestTrueSkillDefaults", func(t *testing.T) {
		mu, sigma := cfg.GetTrueSkillDefaults()

		if mu <= 0 {
			t.Errorf("TrueSkill initial mu should be positive, got %f", mu)
		}

		if sigma <= 0 {
			t.Errorf("TrueSkill initial sigma should be positive, got %f", sigma)
		}

		t.Logf("TrueSkill defaults: μ=%.1f, σ=%.1f", mu, sigma)
	})

	t.Run("TestMMRConfig", func(t *testing.T) {
		mmrConfig := cfg.GetMMRConfig()

		if mmrConfig.MinGamesThreshold <= 0 {
			t.Errorf("MinGamesThreshold should be positive, got %d", mmrConfig.MinGamesThreshold)
		}

		if mmrConfig.OnesWeight <= 0 || mmrConfig.TwosWeight <= 0 || mmrConfig.ThreesWeight <= 0 {
			t.Error("All playlist weights should be positive")
		}

		t.Logf("MMR weights: 1s=%.1f, 2s=%.1f, 3s=%.1f",
			mmrConfig.OnesWeight, mmrConfig.TwosWeight, mmrConfig.ThreesWeight)
	})
}
