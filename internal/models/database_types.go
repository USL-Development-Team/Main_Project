package models

type PublicUsersSelect struct {
	Active               bool    `json:"active"`
	Banned               bool    `json:"banned"`
	CreatedAt            string  `json:"created_at"`
	DiscordId            string  `json:"discord_id"`
	Id                   int64   `json:"id"`
	Mmr                  int32   `json:"mmr"`
	Name                 string  `json:"name"`
	TrueskillLastUpdated string  `json:"trueskill_last_updated"`
	TrueskillMu          float64 `json:"trueskill_mu"`
	TrueskillSigma       float64 `json:"trueskill_sigma"`
	UpdatedAt            string  `json:"updated_at"`
}

type PublicUsersInsert struct {
	Active               *bool    `json:"active"`
	Banned               *bool    `json:"banned"`
	CreatedAt            *string  `json:"created_at"`
	DiscordId            string   `json:"discord_id"`
	Id                   *int64   `json:"id"`
	Mmr                  *int32   `json:"mmr"`
	Name                 string   `json:"name"`
	TrueskillLastUpdated *string  `json:"trueskill_last_updated"`
	TrueskillMu          *float64 `json:"trueskill_mu"`
	TrueskillSigma       *float64 `json:"trueskill_sigma"`
	UpdatedAt            *string  `json:"updated_at"`
}

type PublicUsersUpdate struct {
	Active               *bool    `json:"active"`
	Banned               *bool    `json:"banned"`
	CreatedAt            *string  `json:"created_at"`
	DiscordId            *string  `json:"discord_id"`
	Id                   *int64   `json:"id"`
	Mmr                  *int32   `json:"mmr"`
	Name                 *string  `json:"name"`
	TrueskillLastUpdated *string  `json:"trueskill_last_updated"`
	TrueskillMu          *float64 `json:"trueskill_mu"`
	TrueskillSigma       *float64 `json:"trueskill_sigma"`
	UpdatedAt            *string  `json:"updated_at"`
}

type PublicUserTrackersSelect struct {
	CalculatedMmr             int32  `json:"calculated_mmr"`
	CreatedAt                 string `json:"created_at"`
	DiscordId                 string `json:"discord_id"`
	Id                        int64  `json:"id"`
	LastUpdated               string `json:"last_updated"`
	OnesAllTimePeak           int32  `json:"ones_all_time_peak"`
	OnesCurrentSeasonGames    int32  `json:"ones_current_season_games"`
	OnesCurrentSeasonPeak     int32  `json:"ones_current_season_peak"`
	OnesPreviousSeasonGames   int32  `json:"ones_previous_season_games"`
	OnesPreviousSeasonPeak    int32  `json:"ones_previous_season_peak"`
	ThreesAllTimePeak         int32  `json:"threes_all_time_peak"`
	ThreesCurrentSeasonGames  int32  `json:"threes_current_season_games"`
	ThreesCurrentSeasonPeak   int32  `json:"threes_current_season_peak"`
	ThreesPreviousSeasonGames int32  `json:"threes_previous_season_games"`
	ThreesPreviousSeasonPeak  int32  `json:"threes_previous_season_peak"`
	TwosAllTimePeak           int32  `json:"twos_all_time_peak"`
	TwosCurrentSeasonGames    int32  `json:"twos_current_season_games"`
	TwosCurrentSeasonPeak     int32  `json:"twos_current_season_peak"`
	TwosPreviousSeasonGames   int32  `json:"twos_previous_season_games"`
	TwosPreviousSeasonPeak    int32  `json:"twos_previous_season_peak"`
	UpdatedAt                 string `json:"updated_at"`
	Url                       string `json:"url"`
	Valid                     bool   `json:"valid"`
}

type PublicUserTrackersInsert struct {
	CalculatedMmr             *int32  `json:"calculated_mmr"`
	CreatedAt                 *string `json:"created_at"`
	DiscordId                 string  `json:"discord_id"`
	Id                        *int64  `json:"id"`
	LastUpdated               *string `json:"last_updated"`
	OnesAllTimePeak           *int32  `json:"ones_all_time_peak"`
	OnesCurrentSeasonGames    *int32  `json:"ones_current_season_games"`
	OnesCurrentSeasonPeak     *int32  `json:"ones_current_season_peak"`
	OnesPreviousSeasonGames   *int32  `json:"ones_previous_season_games"`
	OnesPreviousSeasonPeak    *int32  `json:"ones_previous_season_peak"`
	ThreesAllTimePeak         *int32  `json:"threes_all_time_peak"`
	ThreesCurrentSeasonGames  *int32  `json:"threes_current_season_games"`
	ThreesCurrentSeasonPeak   *int32  `json:"threes_current_season_peak"`
	ThreesPreviousSeasonGames *int32  `json:"threes_previous_season_games"`
	ThreesPreviousSeasonPeak  *int32  `json:"threes_previous_season_peak"`
	TwosAllTimePeak           *int32  `json:"twos_all_time_peak"`
	TwosCurrentSeasonGames    *int32  `json:"twos_current_season_games"`
	TwosCurrentSeasonPeak     *int32  `json:"twos_current_season_peak"`
	TwosPreviousSeasonGames   *int32  `json:"twos_previous_season_games"`
	TwosPreviousSeasonPeak    *int32  `json:"twos_previous_season_peak"`
	UpdatedAt                 *string `json:"updated_at"`
	Url                       string  `json:"url"`
	Valid                     *bool   `json:"valid"`
}

type PublicUserTrackersUpdate struct {
	CalculatedMmr             *int32  `json:"calculated_mmr"`
	CreatedAt                 *string `json:"created_at"`
	DiscordId                 *string `json:"discord_id"`
	Id                        *int64  `json:"id"`
	LastUpdated               *string `json:"last_updated"`
	OnesAllTimePeak           *int32  `json:"ones_all_time_peak"`
	OnesCurrentSeasonGames    *int32  `json:"ones_current_season_games"`
	OnesCurrentSeasonPeak     *int32  `json:"ones_current_season_peak"`
	OnesPreviousSeasonGames   *int32  `json:"ones_previous_season_games"`
	OnesPreviousSeasonPeak    *int32  `json:"ones_previous_season_peak"`
	ThreesAllTimePeak         *int32  `json:"threes_all_time_peak"`
	ThreesCurrentSeasonGames  *int32  `json:"threes_current_season_games"`
	ThreesCurrentSeasonPeak   *int32  `json:"threes_current_season_peak"`
	ThreesPreviousSeasonGames *int32  `json:"threes_previous_season_games"`
	ThreesPreviousSeasonPeak  *int32  `json:"threes_previous_season_peak"`
	TwosAllTimePeak           *int32  `json:"twos_all_time_peak"`
	TwosCurrentSeasonGames    *int32  `json:"twos_current_season_games"`
	TwosCurrentSeasonPeak     *int32  `json:"twos_current_season_peak"`
	TwosPreviousSeasonGames   *int32  `json:"twos_previous_season_games"`
	TwosPreviousSeasonPeak    *int32  `json:"twos_previous_season_peak"`
	UpdatedAt                 *string `json:"updated_at"`
	Url                       *string `json:"url"`
	Valid                     *bool   `json:"valid"`
}
