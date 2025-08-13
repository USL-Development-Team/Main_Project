package models

import (
	"database/sql/driver"
	"strings"
	"time"
)

// UserTracker represents the tracker entity, matching the Google Sheets UserTracker schema exactly
type UserTracker struct {
	ID                        int       `json:"id" db:"id"`
	DiscordID                 string    `json:"discord_id" db:"discord_id" validate:"required,min=17,max=19"`
	URL                       string    `json:"url" db:"url" validate:"required,max=1000"`
	OnesCurrentSeasonPeak     int       `json:"ones_current_season_peak" db:"ones_current_season_peak"`
	OnesPreviousSeasonPeak    int       `json:"ones_previous_season_peak" db:"ones_previous_season_peak"`
	OnesAllTimePeak           int       `json:"ones_all_time_peak" db:"ones_all_time_peak"`
	OnesCurrentSeasonGames    int       `json:"ones_current_season_games" db:"ones_current_season_games"`
	OnesPreviousSeasonGames   int       `json:"ones_previous_season_games" db:"ones_previous_season_games"`
	TwosCurrentSeasonPeak     int       `json:"twos_current_season_peak" db:"twos_current_season_peak"`
	TwosPreviousSeasonPeak    int       `json:"twos_previous_season_peak" db:"twos_previous_season_peak"`
	TwosAllTimePeak           int       `json:"twos_all_time_peak" db:"twos_all_time_peak"`
	TwosCurrentSeasonGames    int       `json:"twos_current_season_games" db:"twos_current_season_games"`
	TwosPreviousSeasonGames   int       `json:"twos_previous_season_games" db:"twos_previous_season_games"`
	ThreesCurrentSeasonPeak   int       `json:"threes_current_season_peak" db:"threes_current_season_peak"`
	ThreesPreviousSeasonPeak  int       `json:"threes_previous_season_peak" db:"threes_previous_season_peak"`
	ThreesAllTimePeak         int       `json:"threes_all_time_peak" db:"threes_all_time_peak"`
	ThreesCurrentSeasonGames  int       `json:"threes_current_season_games" db:"threes_current_season_games"`
	ThreesPreviousSeasonGames int       `json:"threes_previous_season_games" db:"threes_previous_season_games"`
	LastUpdated               time.Time `json:"last_updated" db:"last_updated"`
	Valid                     bool      `json:"valid" db:"valid"`
	CalculatedMMR             int       `json:"calculated_mmr" db:"calculated_mmr"`
	CreatedAt                 time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                 time.Time `json:"updated_at" db:"updated_at"`

	// Computed fields (not stored in database)
	PlatformInfo string `json:"platform_info,omitempty"`
	DisplayText  string `json:"display_text,omitempty"`
}

// TrackerCreateRequest matches the form data from AddUserTrackerForm.html
type TrackerCreateRequest struct {
	DiscordID                 string `json:"discord_id" validate:"required,min=17,max=19"`
	URL                       string `json:"url" validate:"required,max=1000"`
	OnesCurrentSeasonPeak     int    `json:"ones_current_season_peak"`
	OnesPreviousSeasonPeak    int    `json:"ones_previous_season_peak"`
	OnesAllTimePeak           int    `json:"ones_all_time_peak"`
	OnesCurrentSeasonGames    int    `json:"ones_current_season_games"`
	OnesPreviousSeasonGames   int    `json:"ones_previous_season_games"`
	TwosCurrentSeasonPeak     int    `json:"twos_current_season_peak"`
	TwosPreviousSeasonPeak    int    `json:"twos_previous_season_peak"`
	TwosAllTimePeak           int    `json:"twos_all_time_peak"`
	TwosCurrentSeasonGames    int    `json:"twos_current_season_games"`
	TwosPreviousSeasonGames   int    `json:"twos_previous_season_games"`
	ThreesCurrentSeasonPeak   int    `json:"threes_current_season_peak"`
	ThreesPreviousSeasonPeak  int    `json:"threes_previous_season_peak"`
	ThreesAllTimePeak         int    `json:"threes_all_time_peak"`
	ThreesCurrentSeasonGames  int    `json:"threes_current_season_games"`
	ThreesPreviousSeasonGames int    `json:"threes_previous_season_games"`
	Valid                     bool   `json:"valid"`
}

// TrackerUpdateRequest matches the form data from UpdateUserTrackerForm.html
type TrackerUpdateRequest struct {
	DiscordID                 string `json:"discord_id" validate:"required,min=17,max=19"`
	URL                       string `json:"url" validate:"required,max=1000"`
	OnesCurrentSeasonPeak     int    `json:"ones_current_season_peak"`
	OnesPreviousSeasonPeak    int    `json:"ones_previous_season_peak"`
	OnesAllTimePeak           int    `json:"ones_all_time_peak"`
	OnesCurrentSeasonGames    int    `json:"ones_current_season_games"`
	OnesPreviousSeasonGames   int    `json:"ones_previous_season_games"`
	TwosCurrentSeasonPeak     int    `json:"twos_current_season_peak"`
	TwosPreviousSeasonPeak    int    `json:"twos_previous_season_peak"`
	TwosAllTimePeak           int    `json:"twos_all_time_peak"`
	TwosCurrentSeasonGames    int    `json:"twos_current_season_games"`
	TwosPreviousSeasonGames   int    `json:"twos_previous_season_games"`
	ThreesCurrentSeasonPeak   int    `json:"threes_current_season_peak"`
	ThreesPreviousSeasonPeak  int    `json:"threes_previous_season_peak"`
	ThreesAllTimePeak         int    `json:"threes_all_time_peak"`
	ThreesCurrentSeasonGames  int    `json:"threes_current_season_games"`
	ThreesPreviousSeasonGames int    `json:"threes_previous_season_games"`
	Valid                     bool   `json:"valid"`
}

// TrackerStats represents tracker statistics, matching JavaScript getTrackerStats() output
type TrackerStats struct {
	TotalTrackers          int     `json:"totalTrackers"`
	ValidTrackers          int     `json:"validTrackers"`
	AverageMMR             float64 `json:"averageMMR"`
	UniqueUsers            int     `json:"uniqueUsers"`
	AverageGamesPerTracker float64 `json:"averageGamesPerTracker"`
}

// Value implements driver.Valuer for database compatibility
func (ut UserTracker) Value() (driver.Value, error) {
	return ut.ID, nil
}

// ParsePlatformInfo extracts platform information from tracker URL
// Matches JavaScript parseTrackerUrl() function
func (ut *UserTracker) ParsePlatformInfo() string {
	if ut.URL == "" {
		return "unknown"
	}

	url := strings.ToLower(ut.URL)

	if strings.Contains(url, "/steam/") {
		return "steam"
	} else if strings.Contains(url, "/epic/") {
		return "epic"
	} else if strings.Contains(url, "/psn/") {
		return "psn"
	} else if strings.Contains(url, "/xbl/") {
		return "xbox"
	}

	return "unknown"
}

// GenerateDisplayText creates display text for UI purposes
// Matches JavaScript tracker.displayText generation
func (ut *UserTracker) GenerateDisplayText() string {
	platformInfo := ut.ParsePlatformInfo()
	return "/" + platformInfo + " " + string(rune(ut.CalculatedMMR))
}

// TotalGamesPlayed calculates total games across all playlists
func (ut *UserTracker) TotalGamesPlayed() int {
	return ut.OnesCurrentSeasonGames + ut.OnesPreviousSeasonGames +
		ut.TwosCurrentSeasonGames + ut.TwosPreviousSeasonGames +
		ut.ThreesCurrentSeasonGames + ut.ThreesPreviousSeasonGames
}

// HasGameData checks if the tracker has meaningful game data
func (ut *UserTracker) HasGameData() bool {
	return ut.TotalGamesPlayed() > 0 ||
		ut.OnesCurrentSeasonPeak > 0 || ut.OnesPreviousSeasonPeak > 0 ||
		ut.TwosCurrentSeasonPeak > 0 || ut.TwosPreviousSeasonPeak > 0 ||
		ut.ThreesCurrentSeasonPeak > 0 || ut.ThreesPreviousSeasonPeak > 0
}

// IsValidTracker checks if tracker should be considered for calculations
func (ut *UserTracker) IsValidTracker() bool {
	return ut.Valid && ut.HasGameData()
}
