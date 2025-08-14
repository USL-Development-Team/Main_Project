package services

import (
	"usl-server/internal/models"
	"usl-server/internal/repositories"
)

// GuildConfigProvider interface for getting guild configurations
type GuildConfigProvider interface {
	GetConfig(guildID int64) (*models.GuildConfig, error)
}

// PermissionService handles role-based permissions for guilds
type PermissionService struct {
	guildRepo GuildConfigProvider
}

// NewPermissionService creates a new permission service
func NewPermissionService(guildRepo *repositories.GuildRepository) *PermissionService {
	return &PermissionService{
		guildRepo: guildRepo,
	}
}

func (s *PermissionService) CanAddUsers(userRoles []string, guildID int64) bool {
	config, err := s.guildRepo.GetConfig(guildID)
	if err != nil {
		return false // fail secure
	}

	return config.HasModeratorRole(userRoles)
}

func (s *PermissionService) CanAddTrackers(userRoles []string, guildID int64) bool {
	config, err := s.guildRepo.GetConfig(guildID)
	if err != nil {
		return false // fail secure
	}

	return config.HasModeratorRole(userRoles)
}

func (s *PermissionService) CanRunAdminCommands(userRoles []string, guildID int64) bool {
	config, err := s.guildRepo.GetConfig(guildID)
	if err != nil {
		return false // fail secure
	}

	return config.HasAdminRole(userRoles)
}

func (s *PermissionService) CanManageGuildConfig(userRoles []string, guildID int64) bool {
	config, err := s.guildRepo.GetConfig(guildID)
	if err != nil {
		return false // fail secure
	}

	return config.HasAdminRole(userRoles)
}

// GetUserPermissions returns a summary of what the user can do in the guild
func (s *PermissionService) GetUserPermissions(userRoles []string, guildID int64) UserPermissions {
	config, err := s.guildRepo.GetConfig(guildID)
	if err != nil {
		return UserPermissions{} // no permissions if config can't be loaded
	}

	isAdmin := config.HasAdminRole(userRoles)
	isModerator := config.HasModeratorRole(userRoles)

	return UserPermissions{
		CanAddUsers:         isModerator,
		CanAddTrackers:      isModerator,
		CanRunAdminCommands: isAdmin,
		CanManageConfig:     isAdmin,
		IsAdmin:             isAdmin,
		IsModerator:         isModerator,
	}
}

func (s *PermissionService) ValidateGuildAccess(userRoles []string, guildID int64) bool {
	permissions := s.GetUserPermissions(userRoles, guildID)
	return permissions.HasAnyPermissions()
}

// UserPermissions represents what a user can do in a specific guild
type UserPermissions struct {
	CanAddUsers         bool `json:"can_add_users"`
	CanAddTrackers      bool `json:"can_add_trackers"`
	CanRunAdminCommands bool `json:"can_run_admin_commands"`
	CanManageConfig     bool `json:"can_manage_config"`
	IsAdmin             bool `json:"is_admin"`
	IsModerator         bool `json:"is_moderator"`
}

func (up UserPermissions) HasAnyPermissions() bool {
	return up.CanAddUsers || up.CanAddTrackers || up.CanRunAdminCommands || up.CanManageConfig
}

func (up UserPermissions) PermissionLevel() string {
	if up.IsAdmin {
		return "Admin"
	}
	if up.IsModerator {
		return "Moderator"
	}
	if up.HasAnyPermissions() {
		return "Member"
	}
	return "None"
}
