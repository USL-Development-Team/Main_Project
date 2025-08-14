package models

import (
	"fmt"
	"regexp"
)

// GuildConfig represents the JSONB configuration for a Discord guild
type GuildConfig struct {
	Discord     DiscordConfig    `json:"discord"`
	Permissions PermissionConfig `json:"permissions"`
}

// DiscordConfig contains Discord-specific integration settings
type DiscordConfig struct {
	AnnouncementChannelID *string `json:"announcement_channel_id"`
	LeaderboardChannelID  *string `json:"leaderboard_channel_id"`
	BotCommandPrefix      string  `json:"bot_command_prefix"`
}

// PermissionConfig defines role-based permissions for the guild
type PermissionConfig struct {
	AdminRoleIDs     []string `json:"admin_role_ids"`
	ModeratorRoleIDs []string `json:"moderator_role_ids"`
}

// Discord snowflake ID validation (18-19 digits)
var discordSnowflakeRegex = regexp.MustCompile(`^\d{17,19}$`)

// Validate ensures the guild configuration is valid
func (gc *GuildConfig) Validate() error {
	// Set sensible defaults
	if gc.Discord.BotCommandPrefix == "" {
		gc.Discord.BotCommandPrefix = "!usl"
	}

	// Validate Discord snowflake IDs format
	allRoleIDs := append(gc.Permissions.AdminRoleIDs, gc.Permissions.ModeratorRoleIDs...)
	for _, roleID := range allRoleIDs {
		if !isValidDiscordSnowflake(roleID) {
			return fmt.Errorf("invalid Discord role ID format: %s", roleID)
		}
	}

	// Validate channel IDs if provided
	if gc.Discord.AnnouncementChannelID != nil && !isValidDiscordSnowflake(*gc.Discord.AnnouncementChannelID) {
		return fmt.Errorf("invalid Discord announcement channel ID: %s", *gc.Discord.AnnouncementChannelID)
	}

	if gc.Discord.LeaderboardChannelID != nil && !isValidDiscordSnowflake(*gc.Discord.LeaderboardChannelID) {
		return fmt.Errorf("invalid Discord leaderboard channel ID: %s", *gc.Discord.LeaderboardChannelID)
	}

	// Validate command prefix
	if len(gc.Discord.BotCommandPrefix) > 5 {
		return fmt.Errorf("bot command prefix too long (max 5 characters): %s", gc.Discord.BotCommandPrefix)
	}

	return nil
}

// isValidDiscordSnowflake checks if a string is a valid Discord snowflake ID
func isValidDiscordSnowflake(id string) bool {
	return discordSnowflakeRegex.MatchString(id)
}

// HasAdminRole checks if any of the provided role IDs match admin roles
func (gc *GuildConfig) HasAdminRole(userRoleIDs []string) bool {
	return hasAnyRole(userRoleIDs, gc.Permissions.AdminRoleIDs)
}

// HasModeratorRole checks if any of the provided role IDs match moderator or admin roles
func (gc *GuildConfig) HasModeratorRole(userRoleIDs []string) bool {
	allModRoles := append(gc.Permissions.AdminRoleIDs, gc.Permissions.ModeratorRoleIDs...)
	return hasAnyRole(userRoleIDs, allModRoles)
}

// hasAnyRole checks if any user role matches any allowed role
func hasAnyRole(userRoles, allowedRoles []string) bool {
	for _, userRole := range userRoles {
		for _, allowedRole := range allowedRoles {
			if userRole == allowedRole {
				return true
			}
		}
	}
	return false
}
