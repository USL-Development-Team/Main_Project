package models

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

type PublicUsersSelect struct {
	Active    bool   `json:"active"`
	Banned    bool   `json:"banned"`
	CreatedAt string `json:"created_at"`
	DiscordId string `json:"discord_id"`
	Id        int64  `json:"id"`
	Name      string `json:"name"`
	UpdatedAt string `json:"updated_at"`
}

type PublicUsersInsert struct {
	Active    *bool   `json:"active"`
	Banned    *bool   `json:"banned"`
	CreatedAt *string `json:"created_at"`
	DiscordId string  `json:"discord_id"`
	Id        *int64  `json:"id"`
	Name      string  `json:"name"`
	UpdatedAt *string `json:"updated_at"`
}

type PublicUsersUpdate struct {
	Active    *bool   `json:"active"`
	Banned    *bool   `json:"banned"`
	CreatedAt *string `json:"created_at"`
	DiscordId *string `json:"discord_id"`
	Id        *int64  `json:"id"`
	Name      *string `json:"name"`
	UpdatedAt *string `json:"updated_at"`
}

type PublicGuildsSelect struct {
	Active         bool        `json:"active"`
	Config         interface{} `json:"config"`
	CreatedAt      string      `json:"created_at"`
	DiscordGuildId string      `json:"discord_guild_id"`
	Id             int64       `json:"id"`
	Name           string      `json:"name"`
	UpdatedAt      string      `json:"updated_at"`
}

type PublicGuildsInsert struct {
	Active         *bool       `json:"active"`
	Config         interface{} `json:"config"`
	CreatedAt      *string     `json:"created_at"`
	DiscordGuildId string      `json:"discord_guild_id"`
	Id             *int64      `json:"id"`
	Name           string      `json:"name"`
	UpdatedAt      *string     `json:"updated_at"`
}

type PublicGuildsUpdate struct {
	Active         *bool       `json:"active"`
	Config         interface{} `json:"config"`
	CreatedAt      *string     `json:"created_at"`
	DiscordGuildId *string     `json:"discord_guild_id"`
	Id             *int64      `json:"id"`
	Name           *string     `json:"name"`
	UpdatedAt      *string     `json:"updated_at"`
}

type PublicUserGuildMembershipsSelect struct {
	Active         bool      `json:"active"`
	CreatedAt      string    `json:"created_at"`
	DiscordRoles   []*string `json:"discord_roles"`
	GuildId        int64     `json:"guild_id"`
	Id             int64     `json:"id"`
	JoinedAt       string    `json:"joined_at"`
	UpdatedAt      string    `json:"updated_at"`
	UserId         int64     `json:"user_id"`
	UslPermissions []*string `json:"usl_permissions"`
}

type PublicUserGuildMembershipsInsert struct {
	Active         *bool     `json:"active"`
	CreatedAt      *string   `json:"created_at"`
	DiscordRoles   []*string `json:"discord_roles"`
	GuildId        int64     `json:"guild_id"`
	Id             *int64    `json:"id"`
	JoinedAt       *string   `json:"joined_at"`
	UpdatedAt      *string   `json:"updated_at"`
	UserId         int64     `json:"user_id"`
	UslPermissions []*string `json:"usl_permissions"`
}

type PublicUserGuildMembershipsUpdate struct {
	Active         *bool     `json:"active"`
	CreatedAt      *string   `json:"created_at"`
	DiscordRoles   []*string `json:"discord_roles"`
	GuildId        *int64    `json:"guild_id"`
	Id             *int64    `json:"id"`
	JoinedAt       *string   `json:"joined_at"`
	UpdatedAt      *string   `json:"updated_at"`
	UserId         *int64    `json:"user_id"`
	UslPermissions []*string `json:"usl_permissions"`
}

type PublicPlayerEffectiveMmrSelect struct {
	CreatedAt      string  `json:"created_at"`
	GamesPlayed    int32   `json:"games_played"`
	GuildId        int64   `json:"guild_id"`
	Id             int64   `json:"id"`
	LastUpdated    string  `json:"last_updated"`
	Mmr            int32   `json:"mmr"`
	TrueskillMu    float64 `json:"trueskill_mu"`
	TrueskillSigma float64 `json:"trueskill_sigma"`
	UpdatedAt      string  `json:"updated_at"`
	UserId         int64   `json:"user_id"`
}

type PublicPlayerEffectiveMmrInsert struct {
	CreatedAt      *string  `json:"created_at"`
	GamesPlayed    *int32   `json:"games_played"`
	GuildId        int64    `json:"guild_id"`
	Id             *int64   `json:"id"`
	LastUpdated    *string  `json:"last_updated"`
	Mmr            *int32   `json:"mmr"`
	TrueskillMu    *float64 `json:"trueskill_mu"`
	TrueskillSigma *float64 `json:"trueskill_sigma"`
	UpdatedAt      *string  `json:"updated_at"`
	UserId         int64    `json:"user_id"`
}

type PublicPlayerEffectiveMmrUpdate struct {
	CreatedAt      *string  `json:"created_at"`
	GamesPlayed    *int32   `json:"games_played"`
	GuildId        *int64   `json:"guild_id"`
	Id             *int64   `json:"id"`
	LastUpdated    *string  `json:"last_updated"`
	Mmr            *int32   `json:"mmr"`
	TrueskillMu    *float64 `json:"trueskill_mu"`
	TrueskillSigma *float64 `json:"trueskill_sigma"`
	UpdatedAt      *string  `json:"updated_at"`
	UserId         *int64   `json:"user_id"`
}

type PublicPlayerHistoricalMmrSelect struct {
	ChangeReason         string   `json:"change_reason"`
	ChangedByUserId      *int64   `json:"changed_by_user_id"`
	CreatedAt            string   `json:"created_at"`
	GuildId              int64    `json:"guild_id"`
	Id                   int64    `json:"id"`
	MatchId              *int64   `json:"match_id"`
	MmrAfter             int32    `json:"mmr_after"`
	MmrBefore            *int32   `json:"mmr_before"`
	TrueskillMuAfter     float64  `json:"trueskill_mu_after"`
	TrueskillMuBefore    *float64 `json:"trueskill_mu_before"`
	TrueskillSigmaAfter  float64  `json:"trueskill_sigma_after"`
	TrueskillSigmaBefore *float64 `json:"trueskill_sigma_before"`
	UserId               int64    `json:"user_id"`
}

type PublicPlayerHistoricalMmrInsert struct {
	ChangeReason         string   `json:"change_reason"`
	ChangedByUserId      *int64   `json:"changed_by_user_id"`
	CreatedAt            *string  `json:"created_at"`
	GuildId              int64    `json:"guild_id"`
	Id                   *int64   `json:"id"`
	MatchId              *int64   `json:"match_id"`
	MmrAfter             int32    `json:"mmr_after"`
	MmrBefore            *int32   `json:"mmr_before"`
	TrueskillMuAfter     float64  `json:"trueskill_mu_after"`
	TrueskillMuBefore    *float64 `json:"trueskill_mu_before"`
	TrueskillSigmaAfter  float64  `json:"trueskill_sigma_after"`
	TrueskillSigmaBefore *float64 `json:"trueskill_sigma_before"`
	UserId               int64    `json:"user_id"`
}

type PublicPlayerHistoricalMmrUpdate struct {
	ChangeReason         *string  `json:"change_reason"`
	ChangedByUserId      *int64   `json:"changed_by_user_id"`
	CreatedAt            *string  `json:"created_at"`
	GuildId              *int64   `json:"guild_id"`
	Id                   *int64   `json:"id"`
	MatchId              *int64   `json:"match_id"`
	MmrAfter             *int32   `json:"mmr_after"`
	MmrBefore            *int32   `json:"mmr_before"`
	TrueskillMuAfter     *float64 `json:"trueskill_mu_after"`
	TrueskillMuBefore    *float64 `json:"trueskill_mu_before"`
	TrueskillSigmaAfter  *float64 `json:"trueskill_sigma_after"`
	TrueskillSigmaBefore *float64 `json:"trueskill_sigma_before"`
	UserId               *int64   `json:"user_id"`
}
