package models

import (
	"database/sql/driver"
	"time"
)

// User represents the user entity, matching the Google Sheets User schema exactly
type User struct {
	ID                   int       `json:"id" db:"id"`
	Name                 string    `json:"name" db:"name" validate:"required,min=1,max=50"`
	DiscordID            string    `json:"discord_id" db:"discord_id" validate:"required,min=17,max=19"`
	Active               bool      `json:"active" db:"active"`
	Banned               bool      `json:"banned" db:"banned"`
	MMR                  int       `json:"mmr" db:"mmr"`
	TrueSkillMu          float64   `json:"trueskill_mu" db:"trueskill_mu"`
	TrueSkillSigma       float64   `json:"trueskill_sigma" db:"trueskill_sigma"`
	TrueSkillLastUpdated time.Time `json:"trueskill_last_updated" db:"trueskill_last_updated"`
	CreatedAt            time.Time `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time `json:"updated_at" db:"updated_at"`
}

// UserCreateRequest matches the form data from AddUserForm.html
type UserCreateRequest struct {
	Name      string `json:"name" validate:"required,min=1,max=50"`
	DiscordID string `json:"discord_id" validate:"required,min=17,max=19"`
	Active    bool   `json:"active"`
	Banned    bool   `json:"banned"`
	MMR       int    `json:"mmr"`
}

// UserUpdateRequest matches the form data from UpdateUserForm.html
type UserUpdateRequest struct {
	Name      string `json:"name" validate:"required,min=1,max=50"`
	DiscordID string `json:"discord_id" validate:"required,min=17,max=19"`
	Active    bool   `json:"active"`
	Banned    bool   `json:"banned"`
	MMR       int    `json:"mmr"`
}

// UserStats represents user statistics, matching JavaScript getUserStats() output
type UserStats struct {
	TotalUsers         int     `json:"totalUsers"`
	ActiveUsers        int     `json:"activeUsers"`
	BannedUsers        int     `json:"bannedUsers"`
	AverageMMR         float64 `json:"averageMMR"`
	AverageTrueSkillMu float64 `json:"averageTrueSkillMu"`
}

// Value implements driver.Valuer for database compatibility
func (u User) Value() (driver.Value, error) {
	return u.ID, nil
}

// DisplayText returns the user's display name for UI purposes
func (u *User) DisplayText() string {
	if u.Name == "" {
		return u.DiscordID
	}
	return u.Name
}

// IsValidForPlay checks if user can participate in matches
func (u *User) IsValidForPlay() bool {
	return u.Active && !u.Banned
}

// HasTrueSkillData checks if the user has been processed by TrueSkill engine
func (u *User) HasTrueSkillData() bool {
	return u.TrueSkillMu != 0 || u.TrueSkillSigma != 0
}
