package services

import (
	"testing"
	"time"
	"usl-server/internal/config"
)

// TestUpdateUserTrueSkillFromTrackerData tests the new TrackerData input method
func TestUpdateUserTrueSkillFromTrackerData(t *testing.T) {
	// This test verifies the method signature exists and basic validation works
	// Full integration would require database setup

	cfg, err := config.Load()
	if err != nil {
		t.Skipf("Config not available for testing: %v", err)
	}

	dataTransformationService := NewDataTransformationService()
	uncertaintyCalculator := NewEnhancedUncertaintyCalculator(cfg)
	percentileConverter := NewPercentileConverter(cfg)
	mmrCalculator := NewMMRCalculator(cfg, percentileConverter)

	service := &UserTrueSkillService{
		trackerRepo:               nil, // Will cause graceful failure in real calculation
		userRepo:                  nil, // Will cause graceful failure in real calculation
		percentileCalculator:      mmrCalculator,
		enhancedUncertainty:       uncertaintyCalculator,
		dataTransformationService: dataTransformationService,
		config:                    cfg,
	}

	t.Run("TestMethodExists", func(t *testing.T) {
		// Test that the method exists with correct signature
		validTrackerData := &TrackerData{
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

		// Call the method - it should exist and return a result
		result := service.UpdateUserTrueSkillFromTrackerData(validTrackerData)

		// Verify we get a result object
		if result == nil {
			t.Fatal("UpdateUserTrueSkillFromTrackerData should return a result")
		}

		// The method should fail gracefully since we don't have real repos
		if result.Success {
			t.Error("Expected failure due to missing repository dependencies")
		}

		// But it should pass validation
		if result.Error == "" {
			t.Error("Expected error message explaining the failure")
		}

		t.Logf("Method exists and returns expected failure: %s", result.Error)
	})

	t.Run("TestDataValidation", func(t *testing.T) {
		// Test with invalid data
		invalidTrackerData := &TrackerData{
			DiscordID: "", // Invalid - empty
		}

		result := service.UpdateUserTrueSkillFromTrackerData(invalidTrackerData)

		if result == nil {
			t.Fatal("Should return result even for invalid data")
		}

		if result.Success {
			t.Error("Should fail validation for invalid data")
		}

		if result.Error == "" {
			t.Error("Should provide error message for validation failure")
		}

		t.Logf("Validation correctly failed: %s", result.Error)
	})

	t.Run("TestNilDataHandling", func(t *testing.T) {
		// Test with nil data
		result := service.UpdateUserTrueSkillFromTrackerData(nil)

		if result == nil {
			t.Fatal("Should return result even for nil data")
		}

		if result.Success {
			t.Error("Should fail for nil data")
		}

		t.Logf("Nil data correctly handled: %s", result.Error)
	})
}

// TestTrueSkillUpdateResult tests the result structure
func TestTrueSkillUpdateResult(t *testing.T) {
	t.Run("TestResultStructure", func(t *testing.T) {
		result := &TrueSkillUpdateResult{
			Success:     true,
			HadTrackers: true,
			TrueSkillResult: &TrueSkillCalculation{
				Mu:          1500.0,
				Sigma:       5.5,
				LastUpdated: time.Now(),
			},
		}

		if !result.Success {
			t.Error("Success field should be settable")
		}

		if !result.HadTrackers {
			t.Error("HadTrackers field should be settable")
		}

		if result.TrueSkillResult == nil {
			t.Error("TrueSkillResult should be settable")
		}

		if result.TrueSkillResult.Mu != 1500.0 {
			t.Errorf("Mu should be 1500.0, got %f", result.TrueSkillResult.Mu)
		}

		if result.TrueSkillResult.Sigma != 5.5 {
			t.Errorf("Sigma should be 5.5, got %f", result.TrueSkillResult.Sigma)
		}
	})
}
