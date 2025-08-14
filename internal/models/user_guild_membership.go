package models

import (
	"database/sql/driver"
	"time"
)

// UserGuildMembership represents a user's membership in a Discord guild
type UserGuildMembership struct {
	ID             int64     `json:"id" db:"id"`
	UserID         int64     `json:"user_id" db:"user_id"`
	GuildID        int64     `json:"guild_id" db:"guild_id"`
	DiscordRoles   []string  `json:"discord_roles" db:"discord_roles"`
	USLPermissions []string  `json:"usl_permissions" db:"usl_permissions"`
	JoinedAt       time.Time `json:"joined_at" db:"joined_at"`
	Active         bool      `json:"active" db:"active"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// USL Permission constants
const (
	PermissionAdmin       = "admin"
	PermissionModerator   = "moderator"
	PermissionStatsViewer = "stats_viewer"
	PermissionUserManager = "user_manager"
)

// UserGuildMembershipCreateRequest represents data needed to create a new membership
type UserGuildMembershipCreateRequest struct {
	UserID         int64    `json:"user_id" validate:"required"`
	GuildID        int64    `json:"guild_id" validate:"required"`
	DiscordRoles   []string `json:"discord_roles"`
	USLPermissions []string `json:"usl_permissions"`
	Active         bool     `json:"active"`
}

// UserGuildMembershipUpdateRequest represents data needed to update an existing membership
type UserGuildMembershipUpdateRequest struct {
	DiscordRoles   []string `json:"discord_roles"`
	USLPermissions []string `json:"usl_permissions"`
	Active         bool     `json:"active"`
}

// Value implements driver.Valuer for database compatibility
func (m UserGuildMembership) Value() (driver.Value, error) {
	return m.ID, nil
}

func (m *UserGuildMembership) HasPermission(permission string) bool {
	for _, p := range m.USLPermissions {
		if p == permission {
			return true
		}
	}
	return false
}

func (m *UserGuildMembership) HasDiscordRole(roleID string) bool {
	for _, r := range m.DiscordRoles {
		if r == roleID {
			return true
		}
	}
	return false
}

func (m *UserGuildMembership) IsAdmin() bool {
	return m.HasPermission(PermissionAdmin)
}

func (m *UserGuildMembership) IsModerator() bool {
	return m.HasPermission(PermissionAdmin) || m.HasPermission(PermissionModerator)
}

func (m *UserGuildMembership) CanAddUsers() bool {
	return m.IsModerator() || m.HasPermission(PermissionUserManager)
}

func (m *UserGuildMembership) CanViewStats() bool {
	return m.IsAdmin() || m.IsModerator() || m.HasPermission(PermissionStatsViewer)
}

func (m *UserGuildMembership) GetPermissionLevel() string {
	if m.IsAdmin() {
		return "Admin"
	}
	if m.HasPermission(PermissionModerator) {
		return "Moderator"
	}
	if m.HasPermission(PermissionUserManager) {
		return "User Manager"
	}
	if m.HasPermission(PermissionStatsViewer) {
		return "Stats Viewer"
	}
	return "Member"
}

func (m *UserGuildMembership) IsValidMembership() bool {
	return m.Active && m.UserID > 0 && m.GuildID > 0
}

func (m *UserGuildMembership) GetMembershipDuration() time.Duration {
	return time.Since(m.JoinedAt)
}

// AddPermission adds a USL permission if not already present
func (m *UserGuildMembership) AddPermission(permission string) {
	if !m.HasPermission(permission) {
		m.USLPermissions = append(m.USLPermissions, permission)
	}
}

// RemovePermission removes a USL permission if present
func (m *UserGuildMembership) RemovePermission(permission string) {
	m.USLPermissions = removeStringFromSlice(m.USLPermissions, permission)
}

// AddDiscordRole adds a Discord role if not already present
func (m *UserGuildMembership) AddDiscordRole(roleID string) {
	if !m.HasDiscordRole(roleID) {
		m.DiscordRoles = append(m.DiscordRoles, roleID)
	}
}

// RemoveDiscordRole removes a Discord role if present
func (m *UserGuildMembership) RemoveDiscordRole(roleID string) {
	m.DiscordRoles = removeStringFromSlice(m.DiscordRoles, roleID)
}

// removeStringFromSlice removes the first occurrence of a string from a slice
func removeStringFromSlice(slice []string, item string) []string {
	for i, s := range slice {
		if s == item {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}
