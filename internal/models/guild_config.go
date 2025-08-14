package models

import (
	"fmt"
	"regexp"
)

const (
	// Default bot command prefix for new guilds
	DefaultBotCommandPrefix = "!usl"

	// Maximum length for bot command prefix
	MaxBotCommandPrefixLength = 5

	// Discord snowflake ID validation pattern (17-19 digits)
	DiscordSnowflakePattern = `^\d{17,19}$`
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

// Discord snowflake ID validation regex
var discordSnowflakeRegex = regexp.MustCompile(DiscordSnowflakePattern)

// Validate ensures the guild configuration is valid
func (gc *GuildConfig) Validate() error {
	gc.setDefaults()

	if err := gc.validateRoleIDs(); err != nil {
		return err
	}

	if err := gc.validateChannelIDs(); err != nil {
		return err
	}

	return gc.validateBotCommandPrefix()
}

// setDefaults applies sensible default values to the configuration
func (gc *GuildConfig) setDefaults() {
	if gc.Discord.BotCommandPrefix == "" {
		gc.Discord.BotCommandPrefix = DefaultBotCommandPrefix
	}
}

// validateRoleIDs validates all role IDs in the configuration
func (gc *GuildConfig) validateRoleIDs() error {
	allRoleIDs := append(gc.Permissions.AdminRoleIDs, gc.Permissions.ModeratorRoleIDs...)
	for _, roleID := range allRoleIDs {
		if !isValidDiscordSnowflake(roleID) {
			return fmt.Errorf("invalid Discord role ID format: %s", roleID)
		}
	}
	return nil
}

// validateChannelIDs validates channel IDs if they are provided
func (gc *GuildConfig) validateChannelIDs() error {
	if err := validateOptionalChannelID(gc.Discord.AnnouncementChannelID, "announcement"); err != nil {
		return err
	}

	return validateOptionalChannelID(gc.Discord.LeaderboardChannelID, "leaderboard")
}

// validateBotCommandPrefix validates the bot command prefix
func (gc *GuildConfig) validateBotCommandPrefix() error {
	if len(gc.Discord.BotCommandPrefix) > MaxBotCommandPrefixLength {
		return fmt.Errorf("bot command prefix too long (max %d characters): %s",
			MaxBotCommandPrefixLength, gc.Discord.BotCommandPrefix)
	}
	return nil
}

// validateOptionalChannelID validates a channel ID if it's provided
func validateOptionalChannelID(channelID *string, channelType string) error {
	if channelID != nil && !isValidDiscordSnowflake(*channelID) {
		return fmt.Errorf("invalid Discord %s channel ID: %s", channelType, *channelID)
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

// hasAnyRole efficiently checks if any user role matches any allowed role
func hasAnyRole(userRoles, allowedRoles []string) bool {
	// Create a set of allowed roles for O(1) lookup
	allowedSet := make(map[string]bool, len(allowedRoles))
	for _, role := range allowedRoles {
		allowedSet[role] = true
	}

	// Check if any user role exists in the allowed set
	for _, userRole := range userRoles {
		if allowedSet[userRole] {
			return true
		}
	}

	return false
}
