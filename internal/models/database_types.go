package models

import "encoding/json"

// User Guild Membership database types
type PublicUserGuildMembershipsSelect struct {
	Id             int64    `json:"id"`
	UserId         int64    `json:"user_id"`
	GuildId        int64    `json:"guild_id"`
	DiscordRoles   []string `json:"discord_roles"`
	UslPermissions []string `json:"usl_permissions"`
	JoinedAt       string   `json:"joined_at"`
	Active         bool     `json:"active"`
	CreatedAt      string   `json:"created_at"`
	UpdatedAt      string   `json:"updated_at"`
}

type PublicUserGuildMembershipsInsert struct {
	Id             *int64    `json:"id"`
	UserId         int64     `json:"user_id"`
	GuildId        int64     `json:"guild_id"`
	DiscordRoles   *[]string `json:"discord_roles"`
	UslPermissions *[]string `json:"usl_permissions"`
	JoinedAt       *string   `json:"joined_at"`
	Active         *bool     `json:"active"`
	CreatedAt      *string   `json:"created_at"`
	UpdatedAt      *string   `json:"updated_at"`
}

type PublicUserGuildMembershipsUpdate struct {
	Id             *int64    `json:"id"`
	UserId         *int64    `json:"user_id"`
	GuildId        *int64    `json:"guild_id"`
	DiscordRoles   *[]string `json:"discord_roles"`
	UslPermissions *[]string `json:"usl_permissions"`
	JoinedAt       *string   `json:"joined_at"`
	Active         *bool     `json:"active"`
	CreatedAt      *string   `json:"created_at"`
	UpdatedAt      *string   `json:"updated_at"`
}

// Player Effective MMR database types
type PublicPlayerEffectiveMmrSelect struct {
	Id             int64   `json:"id"`
	UserId         int64   `json:"user_id"`
	GuildId        int64   `json:"guild_id"`
	Mmr            int32   `json:"mmr"`
	TrueskillMu    float64 `json:"trueskill_mu"`
	TrueskillSigma float64 `json:"trueskill_sigma"`
	GamesPlayed    int32   `json:"games_played"`
	LastUpdated    string  `json:"last_updated"`
	CreatedAt      string  `json:"created_at"`
	UpdatedAt      string  `json:"updated_at"`
}

type PublicPlayerEffectiveMmrInsert struct {
	Id             *int64   `json:"id"`
	UserId         int64    `json:"user_id"`
	GuildId        int64    `json:"guild_id"`
	Mmr            *int32   `json:"mmr"`
	TrueskillMu    *float64 `json:"trueskill_mu"`
	TrueskillSigma *float64 `json:"trueskill_sigma"`
	GamesPlayed    *int32   `json:"games_played"`
	LastUpdated    *string  `json:"last_updated"`
	CreatedAt      *string  `json:"created_at"`
	UpdatedAt      *string  `json:"updated_at"`
}

type PublicPlayerEffectiveMmrUpdate struct {
	Id             *int64   `json:"id"`
	UserId         *int64   `json:"user_id"`
	GuildId        *int64   `json:"guild_id"`
	Mmr            *int32   `json:"mmr"`
	TrueskillMu    *float64 `json:"trueskill_mu"`
	TrueskillSigma *float64 `json:"trueskill_sigma"`
	GamesPlayed    *int32   `json:"games_played"`
	LastUpdated    *string  `json:"last_updated"`
	CreatedAt      *string  `json:"created_at"`
	UpdatedAt      *string  `json:"updated_at"`
}

// Player Historical MMR database types
type PublicPlayerHistoricalMmrSelect struct {
	Id                   int64    `json:"id"`
	UserId               int64    `json:"user_id"`
	GuildId              int64    `json:"guild_id"`
	MmrBefore            *int32   `json:"mmr_before"`
	MmrAfter             int32    `json:"mmr_after"`
	TrueskillMuBefore    *float64 `json:"trueskill_mu_before"`
	TrueskillMuAfter     float64  `json:"trueskill_mu_after"`
	TrueskillSigmaBefore *float64 `json:"trueskill_sigma_before"`
	TrueskillSigmaAfter  float64  `json:"trueskill_sigma_after"`
	ChangeReason         string   `json:"change_reason"`
	MatchId              *int64   `json:"match_id"`
	ChangedByUserId      *int64   `json:"changed_by_user_id"`
	CreatedAt            string   `json:"created_at"`
}

type PublicPlayerHistoricalMmrInsert struct {
	Id                   *int64   `json:"id"`
	UserId               int64    `json:"user_id"`
	GuildId              int64    `json:"guild_id"`
	MmrBefore            *int32   `json:"mmr_before"`
	MmrAfter             int32    `json:"mmr_after"`
	TrueskillMuBefore    *float64 `json:"trueskill_mu_before"`
	TrueskillMuAfter     float64  `json:"trueskill_mu_after"`
	TrueskillSigmaBefore *float64 `json:"trueskill_sigma_before"`
	TrueskillSigmaAfter  float64  `json:"trueskill_sigma_after"`
	ChangeReason         string   `json:"change_reason"`
	MatchId              *int64   `json:"match_id"`
	ChangedByUserId      *int64   `json:"changed_by_user_id"`
	CreatedAt            *string  `json:"created_at"`
}

type PublicPlayerHistoricalMmrUpdate struct {
	Id                   *int64   `json:"id"`
	UserId               *int64   `json:"user_id"`
	GuildId              *int64   `json:"guild_id"`
	MmrBefore            *int32   `json:"mmr_before"`
	MmrAfter             *int32   `json:"mmr_after"`
	TrueskillMuBefore    *float64 `json:"trueskill_mu_before"`
	TrueskillMuAfter     *float64 `json:"trueskill_mu_after"`
	TrueskillSigmaBefore *float64 `json:"trueskill_sigma_before"`
	TrueskillSigmaAfter  *float64 `json:"trueskill_sigma_after"`
	ChangeReason         *string  `json:"change_reason"`
	MatchId              *int64   `json:"match_id"`
	ChangedByUserId      *int64   `json:"changed_by_user_id"`
	CreatedAt            *string  `json:"created_at"`
}

// Guild database types
type PublicGuildsSelect struct {
	Id             int64           `json:"id"`
	DiscordGuildId string          `json:"discord_guild_id"`
	Name           string          `json:"name"`
	Active         bool            `json:"active"`
	Config         json.RawMessage `json:"config"`
	CreatedAt      string          `json:"created_at"`
	UpdatedAt      string          `json:"updated_at"`
}

type PublicGuildsInsert struct {
	Id             *int64           `json:"id"`
	DiscordGuildId string           `json:"discord_guild_id"`
	Name           string           `json:"name"`
	Active         *bool            `json:"active"`
	Config         *json.RawMessage `json:"config"`
	CreatedAt      *string          `json:"created_at"`
	UpdatedAt      *string          `json:"updated_at"`
}

type PublicGuildsUpdate struct {
	Id             *int64           `json:"id"`
	DiscordGuildId *string          `json:"discord_guild_id"`
	Name           *string          `json:"name"`
	Active         *bool            `json:"active"`
	Config         *json.RawMessage `json:"config"`
	CreatedAt      *string          `json:"created_at"`
	UpdatedAt      *string          `json:"updated_at"`
}

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
