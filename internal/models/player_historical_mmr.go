package models

import (
	"database/sql/driver"
	"time"
)

// PlayerHistoricalMMR represents a change in a player's MMR/TrueSkill values
type PlayerHistoricalMMR struct {
	ID                   int64     `json:"id" db:"id"`
	UserID               int64     `json:"user_id" db:"user_id"`
	GuildID              int64     `json:"guild_id" db:"guild_id"`
	MMRBefore            *int      `json:"mmr_before" db:"mmr_before"`
	MMRAfter             int       `json:"mmr_after" db:"mmr_after"`
	TrueSkillMuBefore    *float64  `json:"trueskill_mu_before" db:"trueskill_mu_before"`
	TrueSkillMuAfter     float64   `json:"trueskill_mu_after" db:"trueskill_mu_after"`
	TrueSkillSigmaBefore *float64  `json:"trueskill_sigma_before" db:"trueskill_sigma_before"`
	TrueSkillSigmaAfter  float64   `json:"trueskill_sigma_after" db:"trueskill_sigma_after"`
	ChangeReason         string    `json:"change_reason" db:"change_reason"`
	MatchID              *int64    `json:"match_id" db:"match_id"`
	ChangedByUserID      *int64    `json:"changed_by_user_id" db:"changed_by_user_id"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
}

// ChangeReason constants for different types of MMR changes
const (
	ChangeReasonMatchResult      = "match_result"
	ChangeReasonManualAdjustment = "manual_adjustment"
	ChangeReasonSeasonReset      = "season_reset"
	ChangeReasonInitialSetup     = "initial_setup"
	ChangeReasonRecalculation    = "recalculation"
)

// PlayerHistoricalMMRCreateRequest represents data needed to create a new historical MMR record
type PlayerHistoricalMMRCreateRequest struct {
	UserID               int64    `json:"user_id" validate:"required"`
	GuildID              int64    `json:"guild_id" validate:"required"`
	MMRBefore            *int     `json:"mmr_before"`
	MMRAfter             int      `json:"mmr_after" validate:"min=0"`
	TrueSkillMuBefore    *float64 `json:"trueskill_mu_before"`
	TrueSkillMuAfter     float64  `json:"trueskill_mu_after" validate:"min=0,max=5000"`
	TrueSkillSigmaBefore *float64 `json:"trueskill_sigma_before"`
	TrueSkillSigmaAfter  float64  `json:"trueskill_sigma_after" validate:"min=0,max=20"`
	ChangeReason         string   `json:"change_reason" validate:"required,oneof=match_result manual_adjustment season_reset initial_setup recalculation"`
	MatchID              *int64   `json:"match_id"`
	ChangedByUserID      *int64   `json:"changed_by_user_id"`
}

// Value implements driver.Valuer for database compatibility
func (p PlayerHistoricalMMR) Value() (driver.Value, error) {
	return p.ID, nil
}

// GetMMRChange returns the difference in MMR (positive = gain, negative = loss)
func (p *PlayerHistoricalMMR) GetMMRChange() int {
	if p.MMRBefore == nil {
		return p.MMRAfter // Initial setup case
	}
	return p.MMRAfter - *p.MMRBefore
}

// GetTrueSkillMuChange returns the difference in TrueSkill mu
func (p *PlayerHistoricalMMR) GetTrueSkillMuChange() float64 {
	if p.TrueSkillMuBefore == nil {
		return p.TrueSkillMuAfter // Initial setup case
	}
	return p.TrueSkillMuAfter - *p.TrueSkillMuBefore
}

// GetTrueSkillSigmaChange returns the difference in TrueSkill sigma
func (p *PlayerHistoricalMMR) GetTrueSkillSigmaChange() float64 {
	if p.TrueSkillSigmaBefore == nil {
		return p.TrueSkillSigmaAfter // Initial setup case
	}
	return p.TrueSkillSigmaAfter - *p.TrueSkillSigmaBefore
}

// IsInitialRecord checks if this is the player's first MMR record
func (p *PlayerHistoricalMMR) IsInitialRecord() bool {
	return p.MMRBefore == nil && p.TrueSkillMuBefore == nil && p.TrueSkillSigmaBefore == nil
}

// IsMatchResult checks if this change was due to a match result
func (p *PlayerHistoricalMMR) IsMatchResult() bool {
	return p.ChangeReason == ChangeReasonMatchResult
}

// IsManualAdjustment checks if this change was made manually by an admin
func (p *PlayerHistoricalMMR) IsManualAdjustment() bool {
	return p.ChangeReason == ChangeReasonManualAdjustment
}

// GetChangeDescription returns a human-readable description of the change
func (p *PlayerHistoricalMMR) GetChangeDescription() string {
	mmrChange := p.GetMMRChange()

	switch p.ChangeReason {
	case ChangeReasonMatchResult:
		if mmrChange > 0 {
			return "Won match"
		} else if mmrChange < 0 {
			return "Lost match"
		}
		return "Match played"
	case ChangeReasonManualAdjustment:
		return "Manual adjustment"
	case ChangeReasonSeasonReset:
		return "Season reset"
	case ChangeReasonInitialSetup:
		return "Initial setup"
	case ChangeReasonRecalculation:
		return "Recalculation"
	default:
		return "Unknown change"
	}
}
