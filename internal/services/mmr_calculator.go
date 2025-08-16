package services

import (
	"fmt"
	"math"
	"time"
	"usl-server/internal/config"
	"usl-server/internal/models"
)

// PlaylistData represents MMR and games for current and previous seasons
type PlaylistData struct {
	Current  PlaylistSeason `json:"current"`
	Previous PlaylistSeason `json:"previous"`
}

// PlaylistSeason represents MMR and games for a single season
type PlaylistSeason struct {
	MMR   int `json:"mmr"`
	Games int `json:"games"`
}

// PlayerMMRData represents structured player data for percentile calculations
// Exactly matches the JavaScript playerData structure
type PlayerMMRData struct {
	Ones   PlaylistData `json:"ones"`
	Twos   PlaylistData `json:"twos"`
	Threes PlaylistData `json:"threes"`
}

// Type aliases for TrueSkill service compatibility
type PlayerData = PlayerMMRData
type PlaylistSeasonData = PlaylistSeason
type PercentileSkillResult = SkillCalculationResult

// PlaylistBreakdown represents the calculation breakdown for each playlist
type PlaylistBreakdown struct {
	EffectiveMMR    float64  `json:"effectiveMMR"`
	NormalizedSkill *float64 `json:"normalizedSkill"` // Pointer to allow null
	Games           int      `json:"games"`
}

// SkillCalculationResult matches the JavaScript return structure exactly
type SkillCalculationResult struct {
	NormalizedSkill float64                      `json:"normalizedSkill"`
	TrueskillMu     float64                      `json:"trueskillMu"`
	TotalGames      int                          `json:"totalGames"`
	Breakdown       map[string]PlaylistBreakdown `json:"breakdown"`
	Weights         map[string]float64           `json:"weights"`
	AggregationInfo AggregationInfo              `json:"aggregationInfo"`
	Error           string                       `json:"error,omitempty"`
	Fallback        string                       `json:"fallback,omitempty"`
}

// AggregationInfo provides metadata about the calculation method
type AggregationInfo struct {
	Method    string `json:"method"`
	Converter string `json:"converter"`
	Timestamp string `json:"timestamp"`
}

// MMRCalculator handles all MMR calculation logic
// Exactly ports the JavaScript PercentileMMRCalculator
type MMRCalculator struct {
	config              *config.Config
	percentileConverter *PercentileConverter
}

// NewMMRCalculator creates a new MMR calculator with dependency injection
func NewMMRCalculator(cfg *config.Config, percentileConverter *PercentileConverter) *MMRCalculator {
	return &MMRCalculator{
		config:              cfg,
		percentileConverter: percentileConverter,
	}
}

// CalculatePercentileBasedSkill is the main calculation function
// Exact port of JavaScript calculatePercentileBasedSkill()
func (m *MMRCalculator) CalculatePercentileBasedSkill(playerData PlayerMMRData) SkillCalculationResult {

	// Validate input data
	if err := m.validatePlayerData(playerData); err != nil {
		return SkillCalculationResult{
			NormalizedSkill: 50.0,   // Default to median
			TrueskillMu:     1000.0, // Default μ for 50th percentile (0-2000 range)
			TotalGames:      0,
			Error:           err.Error(),
			Fallback:        "default_values",
		}
	}

	playlists := map[string]PlaylistData{
		"ones":   playerData.Ones,
		"twos":   playerData.Twos,
		"threes": playerData.Threes,
	}

	playlistEffectiveMMRs := make(map[string]float64)
	playlistNormalizedSkills := make(map[string]*float64)

	minGames := 10 // config.RELIABILITY_THRESHOLDS.MIN_GAMES_INCLUSION

	// Calculate effective MMR for each playlist (games-weighted pooling)
	for playlistName, data := range playlists {
		totalGames := data.Current.Games + data.Previous.Games
		var effectiveMMR float64

		if totalGames > 0 {
			effectiveMMR = float64(data.Current.MMR*data.Current.Games+data.Previous.MMR*data.Previous.Games) / float64(totalGames)
		} else {
			effectiveMMR = 0
		}

		playlistEffectiveMMRs[playlistName] = effectiveMMR

		// Convert to normalized skill if meets minimum requirements
		if totalGames >= minGames && effectiveMMR > 0 {
			playlistMapping := map[string]string{
				"ones":   "soloDuel",
				"twos":   "doubles",
				"threes": "standard",
			}

			normalizedSkill := m.percentileConverter.MMRToNormalizedSkill(effectiveMMR, playlistMapping[playlistName])
			playlistNormalizedSkills[playlistName] = &normalizedSkill
		} else {
			playlistNormalizedSkills[playlistName] = nil
		}
	}

	// Playlist weights from config (matches JavaScript weights)
	weights := map[string]float64{
		"ones":   m.config.MMR.OnesWeight,
		"twos":   m.config.MMR.TwosWeight,
		"threes": m.config.MMR.ThreesWeight,
	}

	// Aggregate skills across playlists
	aggregatedSkill := m.percentileConverter.AggregatePlaylistSkills(playlistNormalizedSkills, weights)

	// Convert to TrueSkill μ
	finalMu := m.percentileConverter.NormalizedSkillToTrueSkillMu(aggregatedSkill)

	// Calculate total games
	totalGames := (playerData.Ones.Current.Games + playerData.Ones.Previous.Games) +
		(playerData.Twos.Current.Games + playerData.Twos.Previous.Games) +
		(playerData.Threes.Current.Games + playerData.Threes.Previous.Games)

	breakdown := map[string]PlaylistBreakdown{
		"ones": {
			EffectiveMMR:    playlistEffectiveMMRs["ones"],
			NormalizedSkill: playlistNormalizedSkills["ones"],
			Games:           playerData.Ones.Current.Games + playerData.Ones.Previous.Games,
		},
		"twos": {
			EffectiveMMR:    playlistEffectiveMMRs["twos"],
			NormalizedSkill: playlistNormalizedSkills["twos"],
			Games:           playerData.Twos.Current.Games + playerData.Twos.Previous.Games,
		},
		"threes": {
			EffectiveMMR:    playlistEffectiveMMRs["threes"],
			NormalizedSkill: playlistNormalizedSkills["threes"],
			Games:           playerData.Threes.Current.Games + playerData.Threes.Previous.Games,
		},
	}

	return SkillCalculationResult{
		NormalizedSkill: math.Round(aggregatedSkill*100) / 100,
		TrueskillMu:     math.Round(finalMu*100) / 100,
		TotalGames:      totalGames,
		Breakdown:       breakdown,
		Weights:         weights,
		AggregationInfo: AggregationInfo{
			Method:    "percentile-based",
			Converter: "PercentileConverter",
			Timestamp: time.Now().UTC().Format(time.RFC3339),
		},
	}
}

// CalculatePercentileBasedSkillLegacy provides backward compatibility
// Exact port of JavaScript calculatePercentileBasedSkillLegacy()
func (m *MMRCalculator) CalculatePercentileBasedSkillLegacy(
	onesCurrentPeak, onesCurrentGames, onesPreviousPeak, onesPreviousGames,
	twosCurrentPeak, twosCurrentGames, twosPreviousPeak, twosPreviousGames,
	threesCurrentPeak, threesCurrentGames, threesPreviousPeak, threesPreviousGames int) SkillCalculationResult {

	playerData := PlayerMMRData{
		Ones: PlaylistData{
			Current:  PlaylistSeason{MMR: onesCurrentPeak, Games: onesCurrentGames},
			Previous: PlaylistSeason{MMR: onesPreviousPeak, Games: onesPreviousGames},
		},
		Twos: PlaylistData{
			Current:  PlaylistSeason{MMR: twosCurrentPeak, Games: twosCurrentGames},
			Previous: PlaylistSeason{MMR: twosPreviousPeak, Games: twosPreviousGames},
		},
		Threes: PlaylistData{
			Current:  PlaylistSeason{MMR: threesCurrentPeak, Games: threesCurrentGames},
			Previous: PlaylistSeason{MMR: threesPreviousPeak, Games: threesPreviousGames},
		},
	}

	return m.CalculatePercentileBasedSkill(playerData)
}

// CalculatePercentileSkillFromTracker calculates from UserTracker model
// Matches JavaScript calculatePercentileSkillFromRow()
func (m *MMRCalculator) CalculatePercentileSkillFromTracker(tracker models.UserTracker) SkillCalculationResult {
	playerData := PlayerMMRData{
		Ones: PlaylistData{
			Current:  PlaylistSeason{MMR: tracker.OnesCurrentSeasonPeak, Games: tracker.OnesCurrentSeasonGames},
			Previous: PlaylistSeason{MMR: tracker.OnesPreviousSeasonPeak, Games: tracker.OnesPreviousSeasonGames},
		},
		Twos: PlaylistData{
			Current:  PlaylistSeason{MMR: tracker.TwosCurrentSeasonPeak, Games: tracker.TwosCurrentSeasonGames},
			Previous: PlaylistSeason{MMR: tracker.TwosPreviousSeasonPeak, Games: tracker.TwosPreviousSeasonGames},
		},
		Threes: PlaylistData{
			Current:  PlaylistSeason{MMR: tracker.ThreesCurrentSeasonPeak, Games: tracker.ThreesCurrentSeasonGames},
			Previous: PlaylistSeason{MMR: tracker.ThreesPreviousSeasonPeak, Games: tracker.ThreesPreviousSeasonGames},
		},
	}

	return m.CalculatePercentileBasedSkill(playerData)
}

// validatePlayerData ensures the input data is valid
func (m *MMRCalculator) validatePlayerData(playerData PlayerMMRData) error {
	// Basic validation - could be expanded based on JavaScript validation logic
	if playerData.Ones.Current.Games < 0 || playerData.Ones.Previous.Games < 0 ||
		playerData.Twos.Current.Games < 0 || playerData.Twos.Previous.Games < 0 ||
		playerData.Threes.Current.Games < 0 || playerData.Threes.Previous.Games < 0 {
		return fmt.Errorf("games played cannot be negative")
	}

	if playerData.Ones.Current.MMR < 0 || playerData.Ones.Previous.MMR < 0 ||
		playerData.Twos.Current.MMR < 0 || playerData.Twos.Previous.MMR < 0 ||
		playerData.Threes.Current.MMR < 0 || playerData.Threes.Previous.MMR < 0 {
		return fmt.Errorf("MMR values cannot be negative")
	}

	return nil
}
