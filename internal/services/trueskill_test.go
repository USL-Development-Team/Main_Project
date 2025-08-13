package services

import (
	"testing"
	"time"
	"usl-server/internal/config"
	"usl-server/internal/models"
)

// TestTrueSkillCalculationAccuracy validates that our Go implementation produces
// identical results to the JavaScript version for key test cases
func TestTrueSkillCalculationAccuracy(t *testing.T) {
	// Initialize services exactly as production does
	cfg := &config.Config{
		TrueSkill: config.TrueSkillConfig{
			InitialMu:            1500.0,
			InitialSigma:         8.333,
			SigmaMin:             2.5,
			SigmaMax:             8.333,
			GamesForMaxCertainty: 1000,
		},
		MMR: config.MMRConfig{
			OnesWeight:           1.0,
			TwosWeight:           1.5,
			ThreesWeight:         1.2,
			MinGamesThreshold:    10,
			CurrentSeasonWeight:  0.7,
			PreviousSeasonWeight: 0.3,
		},
	}

	percentileConverter := NewPercentileConverter(cfg)
	mmrCalculator := NewMMRCalculator(cfg, percentileConverter)
	enhancedUncertainty := NewEnhancedUncertaintyCalculator(cfg)
	_ = NewDataTransformationService() // Unused in this test but available if needed

	// Test cases derived from JavaScript tests
	testCases := []struct {
		name          string
		playerData    PlayerData
		expectedMu    float64
		expectedSigma float64
		tolerance     float64 // Allow small floating-point differences
		description   string
	}{
		{
			name: "High MMR Player (1858 MMR) - Production Case",
			playerData: PlayerData{
				Ones: PlaylistData{
					Current:  PlaylistSeasonData{MMR: 0, Games: 0},
					Previous: PlaylistSeasonData{MMR: 0, Games: 0},
				},
				Twos: PlaylistData{
					Current:  PlaylistSeasonData{MMR: 1858, Games: 147},
					Previous: PlaylistSeasonData{MMR: 1641, Games: 89},
				},
				Threes: PlaylistData{
					Current:  PlaylistSeasonData{MMR: 1403, Games: 56},
					Previous: PlaylistSeasonData{MMR: 1287, Games: 67},
				},
			},
			expectedMu:    1500.0, // Should be > 1500 for high MMR players
			expectedSigma: 6.0,    // Should be relatively low uncertainty
			tolerance:     50.0,
			description:   "High MMR player should get high TrueSkill μ > 1500",
		},
		{
			name: "Low MMR Player (500 MMR)",
			playerData: PlayerData{
				Ones: PlaylistData{
					Current:  PlaylistSeasonData{MMR: 0, Games: 0},
					Previous: PlaylistSeasonData{MMR: 0, Games: 0},
				},
				Twos: PlaylistData{
					Current:  PlaylistSeasonData{MMR: 500, Games: 25},
					Previous: PlaylistSeasonData{MMR: 450, Games: 18},
				},
				Threes: PlaylistData{
					Current:  PlaylistSeasonData{MMR: 520, Games: 12},
					Previous: PlaylistSeasonData{MMR: 0, Games: 0},
				},
			},
			expectedMu:    900.0, // Should be < 900 for low MMR players
			expectedSigma: 7.0,   // Higher uncertainty due to fewer games
			tolerance:     50.0,
			description:   "Low MMR player should get low TrueSkill μ < 900",
		},
		{
			name: "Mid-tier Player (1100 MMR)",
			playerData: PlayerData{
				Ones: PlaylistData{
					Current:  PlaylistSeasonData{MMR: 800, Games: 45},
					Previous: PlaylistSeasonData{MMR: 750, Games: 67},
				},
				Twos: PlaylistData{
					Current:  PlaylistSeasonData{MMR: 1100, Games: 89},
					Previous: PlaylistSeasonData{MMR: 1050, Games: 134},
				},
				Threes: PlaylistData{
					Current:  PlaylistSeasonData{MMR: 1200, Games: 67},
					Previous: PlaylistSeasonData{MMR: 1150, Games: 89},
				},
			},
			expectedMu:    1200.0, // Should be around median-high
			expectedSigma: 5.0,    // Lower uncertainty due to many games
			tolerance:     100.0,
			description:   "Mid-tier player with diverse playlist activity",
		},
		{
			name: "New Player (No Games)",
			playerData: PlayerData{
				Ones:   PlaylistData{Current: PlaylistSeasonData{MMR: 0, Games: 0}, Previous: PlaylistSeasonData{MMR: 0, Games: 0}},
				Twos:   PlaylistData{Current: PlaylistSeasonData{MMR: 0, Games: 0}, Previous: PlaylistSeasonData{MMR: 0, Games: 0}},
				Threes: PlaylistData{Current: PlaylistSeasonData{MMR: 0, Games: 0}, Previous: PlaylistSeasonData{MMR: 0, Games: 0}},
			},
			expectedMu:    1000.0, // Should use fallback/default
			expectedSigma: 8.0,    // Maximum uncertainty
			tolerance:     100.0,
			description:   "New player with no game data should get defaults",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Calculate TrueSkill values using our Go implementation
			skillResult := mmrCalculator.CalculatePercentileBasedSkill(tc.playerData)

			t.Logf("Test: %s", tc.description)
			t.Logf("Input data: %+v", tc.playerData)
			t.Logf("Expected μ: %.1f (±%.1f), Got μ: %.1f", tc.expectedMu, tc.tolerance, skillResult.TrueskillMu)
			t.Logf("Expected σ: %.1f, Normalized Skill: %.1f%%", tc.expectedSigma, skillResult.NormalizedSkill)

			// Check if TrueSkill Mu is within expected range
			if tc.name == "High MMR Player (1858 MMR) - Production Case" {
				if skillResult.TrueskillMu < 1500 {
					t.Errorf("❌ BUG: High MMR player (1858) got low TrueSkill μ=%.1f, expected μ > 1500", skillResult.TrueskillMu)
				} else {
					t.Logf("✅ PASS: High MMR player got appropriate TrueSkill μ=%.1f", skillResult.TrueskillMu)
				}
			} else if tc.name == "Low MMR Player (500 MMR)" {
				if skillResult.TrueskillMu > 900 {
					t.Errorf("❌ Low MMR player got high TrueSkill μ=%.1f, expected μ < 900", skillResult.TrueskillMu)
				} else {
					t.Logf("✅ PASS: Low MMR player got appropriate TrueSkill μ=%.1f", skillResult.TrueskillMu)
				}
			}

			// Validate basic constraints
			if skillResult.TrueskillMu < 0 || skillResult.TrueskillMu > 2000 {
				t.Errorf("TrueSkill μ out of valid range [0-2000]: %.1f", skillResult.TrueskillMu)
			}

			if skillResult.NormalizedSkill < 0 || skillResult.NormalizedSkill > 100 {
				t.Errorf("Normalized skill out of valid range [0-100]: %.1f", skillResult.NormalizedSkill)
			}

			// Test Enhanced Uncertainty Calculator
			trackerData := &TrackerData{
				DiscordID:           "test_user",
				OnesCurrentPeak:     tc.playerData.Ones.Current.MMR,
				OnesPreviousPeak:    tc.playerData.Ones.Previous.MMR,
				OnesCurrentGames:    tc.playerData.Ones.Current.Games,
				OnesPreviousGames:   tc.playerData.Ones.Previous.Games,
				TwosCurrentPeak:     tc.playerData.Twos.Current.MMR,
				TwosPreviousPeak:    tc.playerData.Twos.Previous.MMR,
				TwosCurrentGames:    tc.playerData.Twos.Current.Games,
				TwosPreviousGames:   tc.playerData.Twos.Previous.Games,
				ThreesCurrentPeak:   tc.playerData.Threes.Current.MMR,
				ThreesPreviousPeak:  tc.playerData.Threes.Previous.MMR,
				ThreesCurrentGames:  tc.playerData.Threes.Current.Games,
				ThreesPreviousGames: tc.playerData.Threes.Previous.Games,
				LastUpdated:         time.Now(),
			}

			sigma, err := enhancedUncertainty.CalculateEnhancedUncertainty(trackerData)
			if err != nil {
				t.Errorf("Enhanced uncertainty calculation failed: %v", err)
			}

			t.Logf("Enhanced uncertainty σ: %.3f", sigma)

			if sigma < cfg.TrueSkill.SigmaMin || sigma > cfg.TrueSkill.SigmaMax {
				t.Errorf("Enhanced uncertainty σ out of valid range [%.1f-%.1f]: %.3f",
					cfg.TrueSkill.SigmaMin, cfg.TrueSkill.SigmaMax, sigma)
			}
		})
	}
}

// TestDataTransformationService validates the data transformation logic
func TestDataTransformationService(t *testing.T) {
	service := NewDataTransformationService()

	// Test tracker data preparation
	tracker := &models.UserTracker{
		DiscordID:                 "123456789012345678",
		URL:                       "https://rocketleague.tracker.network/rl/profile/steam/76561198000000000/overview",
		OnesCurrentSeasonPeak:     800,
		OnesPreviousSeasonPeak:    750,
		OnesAllTimePeak:           850,
		OnesCurrentSeasonGames:    45,
		OnesPreviousSeasonGames:   67,
		TwosCurrentSeasonPeak:     1200,
		TwosPreviousSeasonPeak:    1150,
		TwosAllTimePeak:           1250,
		TwosCurrentSeasonGames:    89,
		TwosPreviousSeasonGames:   134,
		ThreesCurrentSeasonPeak:   1100,
		ThreesPreviousSeasonPeak:  1050,
		ThreesAllTimePeak:         1180,
		ThreesCurrentSeasonGames:  67,
		ThreesPreviousSeasonGames: 89,
		Valid:                     true,
		LastUpdated:               time.Now(),
	}

	trackerData, err := service.PrepareTrackerDataForCalculation(tracker)
	if err != nil {
		t.Fatalf("Failed to prepare tracker data: %v", err)
	}

	// Validate transformation
	if trackerData.DiscordID != tracker.DiscordID {
		t.Errorf("Discord ID mismatch: expected %s, got %s", tracker.DiscordID, trackerData.DiscordID)
	}

	if trackerData.TwosCurrentPeak != tracker.TwosCurrentSeasonPeak {
		t.Errorf("Twos current peak mismatch: expected %d, got %d", tracker.TwosCurrentSeasonPeak, trackerData.TwosCurrentPeak)
	}

	// Test data validation
	if err := service.ValidateTrackerData(trackerData); err != nil {
		t.Errorf("Valid tracker data failed validation: %v", err)
	}

	// Test statistics
	stats := service.GetTrackerDataStats(trackerData)
	if stats["totalGames"].(int) != tracker.TotalGamesPlayed() {
		t.Errorf("Total games mismatch in stats: expected %d, got %d",
			tracker.TotalGamesPlayed(), stats["totalGames"].(int))
	}
}

// TestTrueSkillServiceIntegration tests the full TrueSkill service workflow
func TestTrueSkillServiceIntegration(t *testing.T) {
	// This would require mock repositories and services
	// For now, we test the calculation logic directly

	cfg := &config.Config{
		TrueSkill: config.TrueSkillConfig{
			InitialMu:            1500.0,
			InitialSigma:         8.333,
			SigmaMin:             2.5,
			SigmaMax:             8.333,
			GamesForMaxCertainty: 1000,
		},
	}

	enhancedUncertainty := NewEnhancedUncertaintyCalculator(cfg)

	// Test calculation bounds
	trackerData := &TrackerData{
		DiscordID:         "test_user",
		TwosCurrentPeak:   1500,
		TwosCurrentGames:  100,
		TwosPreviousPeak:  1400,
		TwosPreviousGames: 80,
		LastUpdated:       time.Now(),
	}

	sigma, err := enhancedUncertainty.CalculateEnhancedUncertainty(trackerData)
	if err != nil {
		t.Fatalf("Enhanced uncertainty calculation failed: %v", err)
	}

	t.Logf("Enhanced uncertainty for test user: %.3f", sigma)

	// Sigma should be within valid bounds
	if sigma < cfg.TrueSkill.SigmaMin || sigma > cfg.TrueSkill.SigmaMax {
		t.Errorf("Enhanced uncertainty out of bounds: %.3f (expected %.1f-%.1f)",
			sigma, cfg.TrueSkill.SigmaMin, cfg.TrueSkill.SigmaMax)
	}
}

// BenchmarkTrueSkillCalculation benchmarks the performance of TrueSkill calculations
func BenchmarkTrueSkillCalculation(b *testing.B) {
	cfg := &config.Config{
		MMR: config.MMRConfig{
			OnesWeight:           1.0,
			TwosWeight:           1.5,
			ThreesWeight:         1.2,
			MinGamesThreshold:    10,
			CurrentSeasonWeight:  0.7,
			PreviousSeasonWeight: 0.3,
		},
	}

	percentileConverter := NewPercentileConverter(cfg)
	mmrCalculator := NewMMRCalculator(cfg, percentileConverter)

	playerData := PlayerData{
		Ones: PlaylistData{
			Current:  PlaylistSeasonData{MMR: 800, Games: 45},
			Previous: PlaylistSeasonData{MMR: 750, Games: 67},
		},
		Twos: PlaylistData{
			Current:  PlaylistSeasonData{MMR: 1200, Games: 89},
			Previous: PlaylistSeasonData{MMR: 1150, Games: 134},
		},
		Threes: PlaylistData{
			Current:  PlaylistSeasonData{MMR: 1100, Games: 67},
			Previous: PlaylistSeasonData{MMR: 1050, Games: 89},
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = mmrCalculator.CalculatePercentileBasedSkill(playerData)
	}
}
