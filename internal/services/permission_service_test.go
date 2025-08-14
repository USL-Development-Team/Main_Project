package services

import (
	"testing"
	"usl-server/internal/models"
)

// MockGuildRepository for testing
type MockGuildRepository struct {
	configs map[int64]*models.GuildConfig
	errors  map[int64]error
}

func NewMockGuildRepository() *MockGuildRepository {
	return &MockGuildRepository{
		configs: make(map[int64]*models.GuildConfig),
		errors:  make(map[int64]error),
	}
}

func (m *MockGuildRepository) GetConfig(guildID int64) (*models.GuildConfig, error) {
	if err, exists := m.errors[guildID]; exists {
		return nil, err
	}

	if config, exists := m.configs[guildID]; exists {
		return config, nil
	}

	return nil, &ConfigNotFoundError{GuildID: guildID}
}

func (m *MockGuildRepository) SetConfig(guildID int64, config *models.GuildConfig) {
	m.configs[guildID] = config
}

func (m *MockGuildRepository) SetError(guildID int64, err error) {
	m.errors[guildID] = err
}

type ConfigNotFoundError struct {
	GuildID int64
}

func (e *ConfigNotFoundError) Error() string {
	return "config not found"
}

func TestPermissionService_CanAddUsers(t *testing.T) {
	mockRepo := NewMockGuildRepository()
	service := &PermissionService{guildRepo: mockRepo}

	// Setup test guild config
	guildID := int64(1)
	config := &models.GuildConfig{
		Permissions: models.PermissionConfig{
			AdminRoleIDs:     []string{"123456789012345678"},
			ModeratorRoleIDs: []string{"987654321098765432"},
		},
	}
	mockRepo.SetConfig(guildID, config)

	tests := []struct {
		name      string
		userRoles []string
		guildID   int64
		want      bool
	}{
		{
			name:      "admin can add users",
			userRoles: []string{"123456789012345678"},
			guildID:   guildID,
			want:      true,
		},
		{
			name:      "moderator can add users",
			userRoles: []string{"987654321098765432"},
			guildID:   guildID,
			want:      true,
		},
		{
			name:      "regular user cannot add users",
			userRoles: []string{"111111111111111111"},
			guildID:   guildID,
			want:      false,
		},
		{
			name:      "no roles cannot add users",
			userRoles: []string{},
			guildID:   guildID,
			want:      false,
		},
		{
			name:      "config error fails secure",
			userRoles: []string{"123456789012345678"},
			guildID:   int64(999), // Non-existent guild
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := service.CanAddUsers(tt.userRoles, tt.guildID); got != tt.want {
				t.Errorf("PermissionService.CanAddUsers() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermissionService_CanRunAdminCommands(t *testing.T) {
	mockRepo := NewMockGuildRepository()
	service := &PermissionService{guildRepo: mockRepo}

	// Setup test guild config
	guildID := int64(1)
	config := &models.GuildConfig{
		Permissions: models.PermissionConfig{
			AdminRoleIDs:     []string{"123456789012345678"},
			ModeratorRoleIDs: []string{"987654321098765432"},
		},
	}
	mockRepo.SetConfig(guildID, config)

	tests := []struct {
		name      string
		userRoles []string
		want      bool
	}{
		{
			name:      "admin can run admin commands",
			userRoles: []string{"123456789012345678"},
			want:      true,
		},
		{
			name:      "moderator cannot run admin commands",
			userRoles: []string{"987654321098765432"},
			want:      false,
		},
		{
			name:      "regular user cannot run admin commands",
			userRoles: []string{"111111111111111111"},
			want:      false,
		},
		{
			name:      "multiple roles with admin",
			userRoles: []string{"111111111111111111", "123456789012345678", "222222222222222222"},
			want:      true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := service.CanRunAdminCommands(tt.userRoles, guildID); got != tt.want {
				t.Errorf("PermissionService.CanRunAdminCommands() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPermissionService_GetUserPermissions(t *testing.T) {
	mockRepo := NewMockGuildRepository()
	service := &PermissionService{guildRepo: mockRepo}

	// Setup test guild config
	guildID := int64(1)
	config := &models.GuildConfig{
		Permissions: models.PermissionConfig{
			AdminRoleIDs:     []string{"123456789012345678"},
			ModeratorRoleIDs: []string{"987654321098765432"},
		},
	}
	mockRepo.SetConfig(guildID, config)

	tests := []struct {
		name      string
		userRoles []string
		want      UserPermissions
	}{
		{
			name:      "admin gets all permissions",
			userRoles: []string{"123456789012345678"},
			want: UserPermissions{
				CanAddUsers:         true,
				CanAddTrackers:      true,
				CanRunAdminCommands: true,
				CanManageConfig:     true,
				IsAdmin:             true,
				IsModerator:         true,
			},
		},
		{
			name:      "moderator gets limited permissions",
			userRoles: []string{"987654321098765432"},
			want: UserPermissions{
				CanAddUsers:         true,
				CanAddTrackers:      true,
				CanRunAdminCommands: false,
				CanManageConfig:     false,
				IsAdmin:             false,
				IsModerator:         true,
			},
		},
		{
			name:      "regular user gets no permissions",
			userRoles: []string{"111111111111111111"},
			want: UserPermissions{
				CanAddUsers:         false,
				CanAddTrackers:      false,
				CanRunAdminCommands: false,
				CanManageConfig:     false,
				IsAdmin:             false,
				IsModerator:         false,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := service.GetUserPermissions(tt.userRoles, guildID)

			if got.CanAddUsers != tt.want.CanAddUsers {
				t.Errorf("CanAddUsers = %v, want %v", got.CanAddUsers, tt.want.CanAddUsers)
			}
			if got.CanAddTrackers != tt.want.CanAddTrackers {
				t.Errorf("CanAddTrackers = %v, want %v", got.CanAddTrackers, tt.want.CanAddTrackers)
			}
			if got.CanRunAdminCommands != tt.want.CanRunAdminCommands {
				t.Errorf("CanRunAdminCommands = %v, want %v", got.CanRunAdminCommands, tt.want.CanRunAdminCommands)
			}
			if got.CanManageConfig != tt.want.CanManageConfig {
				t.Errorf("CanManageConfig = %v, want %v", got.CanManageConfig, tt.want.CanManageConfig)
			}
			if got.IsAdmin != tt.want.IsAdmin {
				t.Errorf("IsAdmin = %v, want %v", got.IsAdmin, tt.want.IsAdmin)
			}
			if got.IsModerator != tt.want.IsModerator {
				t.Errorf("IsModerator = %v, want %v", got.IsModerator, tt.want.IsModerator)
			}
		})
	}
}

func TestUserPermissions_HasAnyPermissions(t *testing.T) {
	tests := []struct {
		name        string
		permissions UserPermissions
		want        bool
	}{
		{
			name: "admin has permissions",
			permissions: UserPermissions{
				CanAddUsers:         true,
				CanRunAdminCommands: true,
			},
			want: true,
		},
		{
			name: "only add users permission",
			permissions: UserPermissions{
				CanAddUsers: true,
			},
			want: true,
		},
		{
			name: "only add trackers permission",
			permissions: UserPermissions{
				CanAddTrackers: true,
			},
			want: true,
		},
		{
			name: "only admin commands permission",
			permissions: UserPermissions{
				CanRunAdminCommands: true,
			},
			want: true,
		},
		{
			name: "only config management permission",
			permissions: UserPermissions{
				CanManageConfig: true,
			},
			want: true,
		},
		{
			name: "no permissions",
			permissions: UserPermissions{
				CanAddUsers:         false,
				CanAddTrackers:      false,
				CanRunAdminCommands: false,
				CanManageConfig:     false,
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.permissions.HasAnyPermissions(); got != tt.want {
				t.Errorf("UserPermissions.HasAnyPermissions() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestUserPermissions_PermissionLevel(t *testing.T) {
	tests := []struct {
		name        string
		permissions UserPermissions
		want        string
	}{
		{
			name: "admin level",
			permissions: UserPermissions{
				IsAdmin: true,
			},
			want: "Admin",
		},
		{
			name: "moderator level",
			permissions: UserPermissions{
				IsModerator: true,
				IsAdmin:     false,
			},
			want: "Moderator",
		},
		{
			name: "member level (has some permissions but not admin/mod)",
			permissions: UserPermissions{
				CanAddUsers: true,
				IsAdmin:     false,
				IsModerator: false,
			},
			want: "Member",
		},
		{
			name: "no permissions",
			permissions: UserPermissions{
				CanAddUsers:         false,
				CanAddTrackers:      false,
				CanRunAdminCommands: false,
				CanManageConfig:     false,
				IsAdmin:             false,
				IsModerator:         false,
			},
			want: "None",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.permissions.PermissionLevel(); got != tt.want {
				t.Errorf("UserPermissions.PermissionLevel() = %v, want %v", got, tt.want)
			}
		})
	}
}
