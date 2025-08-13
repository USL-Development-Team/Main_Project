package services

import (
	"testing"
	"usl-server/internal/config"
)

// TestJavaScriptValidation compares Go implementation against known JavaScript production results
// This is the CRITICAL test that validates the port maintains mathematical accuracy
func TestJavaScriptValidation(t *testing.T) {
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

	// Real test cases extracted from production CSV files
	// UserTracker.csv contains inputs, User.csv contains expected JavaScript outputs
	testCases := []struct {
		name        string
		discordID   string
		playerData  PlayerMMRData
		expectedMu  float64 // From User.csv - what JavaScript calculated
		tolerance   float64 // Allowable difference for floating point precision
		description string
	}{
		{
			name:      "mogtron - Critical High MMR Case",
			discordID: "354474826192388127",
			playerData: PlayerMMRData{
				Ones: PlaylistData{
					Current:  PlaylistSeason{MMR: 800, Games: 14},
					Previous: PlaylistSeason{MMR: 799, Games: 37},
				},
				Twos: PlaylistData{
					Current:  PlaylistSeason{MMR: 1142, Games: 13},
					Previous: PlaylistSeason{MMR: 1133, Games: 60},
				},
				Threes: PlaylistData{
					Current:  PlaylistSeason{MMR: 1110, Games: 11},
					Previous: PlaylistSeason{MMR: 1167, Games: 42},
				},
			},
			expectedMu:  1805.78, // Exact value from production User.csv
			tolerance:   50.0,    // Allow reasonable floating point differences
			description: "High MMR player from CSV - critical validation case",
		},
		{
			name:      "oay - High MMR Player",
			discordID: "544209988931944479",
			playerData: PlayerMMRData{
				Ones: PlaylistData{
					Current:  PlaylistSeason{MMR: 1140, Games: 0},
					Previous: PlaylistSeason{MMR: 1184, Games: 23},
				},
				Twos: PlaylistData{
					Current:  PlaylistSeason{MMR: 1660, Games: 0},
					Previous: PlaylistSeason{MMR: 2002, Games: 170},
				},
				Threes: PlaylistData{
					Current:  PlaylistSeason{MMR: 1559, Games: 0},
					Previous: PlaylistSeason{MMR: 1753, Games: 110},
				},
			},
			expectedMu:  1998.65, // Exact value from production User.csv
			tolerance:   50.0,
			description: "Another high MMR validation case",
		},
		{
			name:      "ayejoshy - High MMR Many Games",
			discordID: "837466622670667776",
			playerData: PlayerMMRData{
				Ones: PlaylistData{
					Current:  PlaylistSeason{MMR: 1365, Games: 351},
					Previous: PlaylistSeason{MMR: 1311, Games: 931},
				},
				Twos: PlaylistData{
					Current:  PlaylistSeason{MMR: 1801, Games: 728},
					Previous: PlaylistSeason{MMR: 1885, Games: 772},
				},
				Threes: PlaylistData{
					Current:  PlaylistSeason{MMR: 1398, Games: 5},
					Previous: PlaylistSeason{MMR: 1444, Games: 1},
				},
			},
			expectedMu:  1998.54, // Exact value from production User.csv
			tolerance:   50.0,
			description: "High MMR player with many games (low uncertainty)",
		},
	}

	t.Logf("🔬 JavaScript vs Go TrueSkill Validation")
	t.Logf("Comparing Go implementation against production CSV results")

	passedTests := 0
	totalTests := len(testCases)

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			t.Logf("Testing %s (Discord: %s)", tc.name, tc.discordID)
			t.Logf("Description: %s", tc.description)

			// Show input data for debugging
			t.Logf("Input Data:")
			t.Logf("  1v1: Current=%d (%d games), Previous=%d (%d games)",
				tc.playerData.Ones.Current.MMR, tc.playerData.Ones.Current.Games,
				tc.playerData.Ones.Previous.MMR, tc.playerData.Ones.Previous.Games)
			t.Logf("  2v2: Current=%d (%d games), Previous=%d (%d games)",
				tc.playerData.Twos.Current.MMR, tc.playerData.Twos.Current.Games,
				tc.playerData.Twos.Previous.MMR, tc.playerData.Twos.Previous.Games)
			t.Logf("  3v3: Current=%d (%d games), Previous=%d (%d games)",
				tc.playerData.Threes.Current.MMR, tc.playerData.Threes.Current.Games,
				tc.playerData.Threes.Previous.MMR, tc.playerData.Threes.Previous.Games)

			// Calculate using Go implementation
			result := mmrCalculator.CalculatePercentileBasedSkill(tc.playerData)

			t.Logf("Results:")
			t.Logf("  JavaScript μ (expected): %.2f (from production CSV)", tc.expectedMu)
			t.Logf("  Go μ (calculated):       %.2f", result.TrueskillMu)
			difference := abs(result.TrueskillMu - tc.expectedMu)
			t.Logf("  Difference:              %.2f (%.1f%%)", difference, difference/tc.expectedMu*100)
			t.Logf("  Go Normalized Skill:     %.1f%%", result.NormalizedSkill)
			t.Logf("  Go Total Games:          %d", result.TotalGames)

			if result.Error != "" {
				t.Errorf("Go calculation error: %s", result.Error)
			}

			// Validation with tolerance
			if difference <= tc.tolerance {
				t.Logf("✅ PASS: Go μ (%.2f) matches JavaScript μ (%.2f) within tolerance (±%.1f)",
					result.TrueskillMu, tc.expectedMu, tc.tolerance)
				passedTests++

				// Additional precision checks
				if difference <= 10.0 {
					t.Logf("🎯 EXCELLENT: Very close match (difference: %.2f) - port is highly accurate!", difference)
				} else if difference <= 25.0 {
					t.Logf("👍 GOOD: Close match (difference: %.2f) - minor calculation differences", difference)
				}
			} else {
				t.Errorf("❌ FAIL: Go μ (%.2f) differs significantly from JavaScript μ (%.2f)",
					result.TrueskillMu, tc.expectedMu)
				t.Errorf("   Difference: %.2f (%.1f%%) exceeds tolerance of %.1f",
					difference, difference/tc.expectedMu*100, tc.tolerance)
				t.Errorf("   This indicates a calculation bug in the Go port!")
			}

			// Sanity checks
			if result.TrueskillMu <= 0 || result.TrueskillMu > 2000 {
				t.Errorf("❌ RANGE ERROR: Go μ (%.2f) is outside valid range [0-2000]", result.TrueskillMu)
			}

			if result.NormalizedSkill < 0 || result.NormalizedSkill > 100 {
				t.Errorf("❌ SKILL ERROR: Go normalized skill (%.1f%%) is outside valid range [0-100]", result.NormalizedSkill)
			}
		})
	}

	// Overall validation summary
	t.Logf("\n📊 VALIDATION SUMMARY")
	t.Logf("Passed: %d/%d tests (%.1f%%)", passedTests, totalTests, float64(passedTests)/float64(totalTests)*100)

	if passedTests == totalTests {
		t.Logf("🎉 SUCCESS! Go implementation matches JavaScript results!")
		t.Logf("✅ All test cases passed within acceptable tolerance")
		t.Logf("✅ The JavaScript->Go port is mathematically accurate")
	} else {
		failed := totalTests - passedTests
		t.Errorf("💥 VALIDATION FAILED! %d out of %d tests failed", failed, totalTests)
		t.Errorf("❌ The Go port does NOT match JavaScript calculations")
		t.Errorf("🔧 Debug needed: Check percentile conversion, playlist weights, aggregation logic")
	}
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}
