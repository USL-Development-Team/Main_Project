package usl

import (
	"time"
)

// USLUser represents a user in the USL-specific migration table
// Mirrors the Google Sheets User.csv structure exactly
type USLUser struct {
	ID                   int64     `json:"id" db:"id"`
	Name                 string    `json:"name" db:"name"`
	DiscordID            string    `json:"discord_id" db:"discord_id"`
	Active               bool      `json:"active" db:"active"`
	Banned               bool      `json:"banned" db:"banned"`
	MMR                  int       `json:"mmr" db:"mmr"`
	TrueSkillMu          float64   `json:"trueskill_mu" db:"trueskill_mu"`
	TrueSkillSigma       float64   `json:"trueskill_sigma" db:"trueskill_sigma"`
	TrueSkillLastUpdated *string   `json:"trueskill_last_updated" db:"trueskill_last_updated"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// GetTrueSkillLastUpdatedFormatted returns a formatted date string for display
func (u *USLUser) GetTrueSkillLastUpdatedFormatted() string {
	if u.TrueSkillLastUpdated == nil || *u.TrueSkillLastUpdated == "" {
		return "Never"
	}
	return *u.TrueSkillLastUpdated
}

// GetLastUpdatedFormatted returns a formatted date string for display
func (t *USLUserTracker) GetLastUpdatedFormatted() string {
	if t.LastUpdated == nil || *t.LastUpdated == "" {
		return "Never"
	}
	return *t.LastUpdated
}

// USLUserTracker represents a user tracker in the USL-specific migration table
// Mirrors the Google Sheets UserTracker.csv structure exactly
type USLUserTracker struct {
	ID                              int64     `json:"id" db:"id"`
	DiscordID                       string    `json:"discord_id" db:"discord_id"`
	URL                             string    `json:"url" db:"url"`
	OnesCurrentSeasonPeak           int       `json:"ones_current_season_peak" db:"ones_current_season_peak"`
	OnesPreviousSeasonPeak          int       `json:"ones_previous_season_peak" db:"ones_previous_season_peak"`
	OnesAllTimePeak                 int       `json:"ones_all_time_peak" db:"ones_all_time_peak"`
	OnesCurrentSeasonGamesPlayed    int       `json:"ones_current_season_games_played" db:"ones_current_season_games_played"`
	OnesPreviousSeasonGamesPlayed   int       `json:"ones_previous_season_games_played" db:"ones_previous_season_games_played"`
	TwosCurrentSeasonPeak           int       `json:"twos_current_season_peak" db:"twos_current_season_peak"`
	TwosPreviousSeasonPeak          int       `json:"twos_previous_season_peak" db:"twos_previous_season_peak"`
	TwosAllTimePeak                 int       `json:"twos_all_time_peak" db:"twos_all_time_peak"`
	TwosCurrentSeasonGamesPlayed    int       `json:"twos_current_season_games_played" db:"twos_current_season_games_played"`
	TwosPreviousSeasonGamesPlayed   int       `json:"twos_previous_season_games_played" db:"twos_previous_season_games_played"`
	ThreesCurrentSeasonPeak         int       `json:"threes_current_season_peak" db:"threes_current_season_peak"`
	ThreesPreviousSeasonPeak        int       `json:"threes_previous_season_peak" db:"threes_previous_season_peak"`
	ThreesAllTimePeak               int       `json:"threes_all_time_peak" db:"threes_all_time_peak"`
	ThreesCurrentSeasonGamesPlayed  int       `json:"threes_current_season_games_played" db:"threes_current_season_games_played"`
	ThreesPreviousSeasonGamesPlayed int       `json:"threes_previous_season_games_played" db:"threes_previous_season_games_played"`
	LastUpdated                     *string   `json:"last_updated" db:"last_updated"`
	Valid                           bool      `json:"valid" db:"valid"`
	MMR                             int       `json:"mmr" db:"mmr"`
	CreatedAt                       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                       time.Time `json:"updated_at" db:"updated_at"`
}

type USLUserCSV struct {
	Name                 string `csv:"name"`
	DiscordID            string `csv:"discord id"`
	Active               string `csv:"active"`
	Banned               string `csv:"banned"`
	MMR                  string `csv:"mmr"`
	TrueSkillMu          string `csv:"trueskill_mu"`
	TrueSkillSigma       string `csv:"trueskill_sigma"`
	TrueSkillLastUpdated string `csv:"trueskill_last_updated"`
}

type USLUserTrackerCSV struct {
	DiscordID                       string `csv:"discord id"`
	URL                             string `csv:"url"`
	OnesCurrentSeasonPeak           string `csv:"ones current season peak"`
	OnesPreviousSeasonPeak          string `csv:"ones previous season peak"`
	OnesAllTimePeak                 string `csv:"ones all time peak"`
	OnesCurrentSeasonGamesPlayed    string `csv:"ones current season games played"`
	OnesPreviousSeasonGamesPlayed   string `csv:"ones previous season games played"`
	TwosCurrentSeasonPeak           string `csv:"twos current season peak"`
	TwosPreviousSeasonPeak          string `csv:"twos previous season peak"`
	TwosAllTimePeak                 string `csv:"twos all time peak"`
	TwosCurrentSeasonGamesPlayed    string `csv:"twos current season games played"`
	TwosPreviousSeasonGamesPlayed   string `csv:"twos previous season games played"`
	ThreesCurrentSeasonPeak         string `csv:"threes current season peak"`
	ThreesPreviousSeasonPeak        string `csv:"threes previous season peak"`
	ThreesAllTimePeak               string `csv:"threes all time peak"`
	ThreesCurrentSeasonGamesPlayed  string `csv:"threes current season games played"`
	ThreesPreviousSeasonGamesPlayed string `csv:"threes previous season games played"`
	LastUpdated                     string `csv:"last updated"`
	Valid                           string `csv:"valid"`
	MMR                             string `csv:"mmr"`
}

// IsValidForPlay checks if user can participate in matches
func (u *USLUser) IsValidForPlay() bool {
	return u.Active && !u.Banned
}

// HasTrueSkillData checks if the user has been processed by TrueSkill engine
func (u *USLUser) HasTrueSkillData() bool {
	return u.TrueSkillMu != 25.0 || u.TrueSkillSigma != 8.333333
}

// DisplayName returns the user's display name for UI purposes
func (u *USLUser) DisplayName() string {
	if u.Name == "" {
		return u.DiscordID
	}
	return u.Name
}

// IsValidTracker checks if tracker has enough data to be considered valid
func (t *USLUserTracker) IsValidTracker() bool {
	return t.Valid
}

// TotalGames calculates total games across all playlists
func (t *USLUserTracker) TotalGames() int {
	return t.OnesCurrentSeasonGamesPlayed + t.TwosCurrentSeasonGamesPlayed + t.ThreesCurrentSeasonGamesPlayed
}

// HasEnoughGames checks if tracker has sufficient games for MMR calculation
func (t *USLUserTracker) HasEnoughGames() bool {
	return t.TotalGames() >= 10 // Arbitrary threshold
}
