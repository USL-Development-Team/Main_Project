package models

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// Guild represents a Discord server with USL integration
type Guild struct {
	ID             int64       `json:"id" db:"id"`
	DiscordGuildID string      `json:"discord_guild_id" db:"discord_guild_id" validate:"required,min=17,max=19"`
	Name           string      `json:"name" db:"name" validate:"required,min=1,max=100"`
	Active         bool        `json:"active" db:"active"`
	Config         GuildConfig `json:"config" db:"config"`
	CreatedAt      time.Time   `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at" db:"updated_at"`
}

// GuildCreateRequest represents data needed to create a new guild
type GuildCreateRequest struct {
	DiscordGuildID string      `json:"discord_guild_id" validate:"required,min=17,max=19"`
	Name           string      `json:"name" validate:"required,min=1,max=100"`
	Config         GuildConfig `json:"config"`
}

// GuildUpdateRequest represents data needed to update an existing guild
type GuildUpdateRequest struct {
	Name   string      `json:"name" validate:"required,min=1,max=100"`
	Active bool        `json:"active"`
	Config GuildConfig `json:"config"`
}

// Value implements driver.Valuer for database compatibility
func (g Guild) Value() (driver.Value, error) {
	return g.ID, nil
}

// Scan implements sql.Scanner for reading JSONB config from database
func (g *Guild) Scan(value any) error {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case []byte:
		return json.Unmarshal(v, &g.Config)
	case string:
		return json.Unmarshal([]byte(v), &g.Config)
	}

	return nil
}

// DisplayText returns the guild's display name for UI purposes
func (g *Guild) DisplayText() string {
	if g.Name == "" {
		return g.DiscordGuildID
	}
	return g.Name
}

// IsValidForUse checks if guild is active and properly configured
func (g *Guild) IsValidForUse() bool {
	return g.Active && g.Config.Validate() == nil
}

// GetDefaultConfig returns a new guild configuration with sensible defaults
func GetDefaultGuildConfig() GuildConfig {
	return GuildConfig{
		Discord: DiscordConfig{
			AnnouncementChannelID: nil,
			LeaderboardChannelID:  nil,
			BotCommandPrefix:      "!usl",
		},
		Permissions: PermissionConfig{
			AdminRoleIDs:     []string{},
			ModeratorRoleIDs: []string{},
		},
	}
}
