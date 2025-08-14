package models

import (
	"database/sql/driver"
	"time"
)

// PlayerEffectiveMMR represents current MMR/TrueSkill data for a user in a specific guild
type PlayerEffectiveMMR struct {
	ID             int64     `json:"id" db:"id"`
	UserID         int64     `json:"user_id" db:"user_id"`
	GuildID        int64     `json:"guild_id" db:"guild_id"`
	MMR            int       `json:"mmr" db:"mmr"`
	TrueSkillMu    float64   `json:"trueskill_mu" db:"trueskill_mu"`
	TrueSkillSigma float64   `json:"trueskill_sigma" db:"trueskill_sigma"`
	GamesPlayed    int       `json:"games_played" db:"games_played"`
	LastUpdated    time.Time `json:"last_updated" db:"last_updated"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// PlayerEffectiveMMRCreateRequest represents data needed to create a new player's MMR record
type PlayerEffectiveMMRCreateRequest struct {
	UserID         int64   `json:"user_id" validate:"required"`
	GuildID        int64   `json:"guild_id" validate:"required"`
	MMR            int     `json:"mmr" validate:"min=0"`
	TrueSkillMu    float64 `json:"trueskill_mu" validate:"min=0,max=5000"`
	TrueSkillSigma float64 `json:"trueskill_sigma" validate:"min=0,max=20"`
	GamesPlayed    int     `json:"games_played" validate:"min=0"`
}

// PlayerEffectiveMMRUpdateRequest represents data needed to update an existing player's MMR record
type PlayerEffectiveMMRUpdateRequest struct {
	MMR            int     `json:"mmr" validate:"min=0"`
	TrueSkillMu    float64 `json:"trueskill_mu" validate:"min=0,max=5000"`
	TrueSkillSigma float64 `json:"trueskill_sigma" validate:"min=0,max=20"`
	GamesPlayed    int     `json:"games_played" validate:"min=0"`
}

// Value implements driver.Valuer for database compatibility
func (p PlayerEffectiveMMR) Value() (driver.Value, error) {
	return p.ID, nil
}

// IsValidForCompetition checks if player has enough games and valid stats for competitive play
func (p *PlayerEffectiveMMR) IsValidForCompetition(minGames int) bool {
	return p.GamesPlayed >= minGames && p.TrueSkillMu > 0 && p.TrueSkillSigma > 0
}

// GetSkillUncertainty returns a normalized uncertainty value (0-1, where 0 is certain)
func (p *PlayerEffectiveMMR) GetSkillUncertainty() float64 {
	// Higher sigma = more uncertainty
	// Normalize to 0-1 range where 8.333 (initial) = 1.0 and 0 = 0.0
	maxSigma := 8.333
	if p.TrueSkillSigma >= maxSigma {
		return 1.0
	}
	return p.TrueSkillSigma / maxSigma
}

// GetSkillEstimate returns the conservative skill estimate (mu - 3*sigma)
func (p *PlayerEffectiveMMR) GetSkillEstimate() float64 {
	return p.TrueSkillMu - (3.0 * p.TrueSkillSigma)
}

// HasRecentActivity checks if the player has been active recently
func (p *PlayerEffectiveMMR) HasRecentActivity(daysSince int) bool {
	cutoff := time.Now().AddDate(0, 0, -daysSince)
	return p.LastUpdated.After(cutoff)
}
