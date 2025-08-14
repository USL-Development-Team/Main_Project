package models

import (
	"testing"
)

func TestGuildConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  GuildConfig
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config with defaults",
			config: GuildConfig{
				Discord: DiscordConfig{
					BotCommandPrefix: "!usl",
				},
				Permissions: PermissionConfig{
					AdminRoleIDs:     []string{"123456789012345678"},
					ModeratorRoleIDs: []string{"987654321098765432"},
				},
			},
			wantErr: false,
		},
		{
			name: "empty config gets defaults",
			config: GuildConfig{
				Discord:     DiscordConfig{},
				Permissions: PermissionConfig{},
			},
			wantErr: false,
		},
		{
			name: "invalid admin role ID",
			config: GuildConfig{
				Discord: DiscordConfig{
					BotCommandPrefix: "!usl",
				},
				Permissions: PermissionConfig{
					AdminRoleIDs: []string{"invalid-id"},
				},
			},
			wantErr: true,
			errMsg:  "invalid Discord role ID format: invalid-id",
		},
		{
			name: "invalid moderator role ID",
			config: GuildConfig{
				Discord: DiscordConfig{
					BotCommandPrefix: "!usl",
				},
				Permissions: PermissionConfig{
					ModeratorRoleIDs: []string{"123"},
				},
			},
			wantErr: true,
			errMsg:  "invalid Discord role ID format: 123",
		},
		{
			name: "invalid announcement channel ID",
			config: GuildConfig{
				Discord: DiscordConfig{
					AnnouncementChannelID: stringPtr("invalid"),
					BotCommandPrefix:      "!usl",
				},
				Permissions: PermissionConfig{},
			},
			wantErr: true,
			errMsg:  "invalid Discord announcement channel ID: invalid",
		},
		{
			name: "invalid leaderboard channel ID",
			config: GuildConfig{
				Discord: DiscordConfig{
					LeaderboardChannelID: stringPtr("xyz"),
					BotCommandPrefix:     "!usl",
				},
				Permissions: PermissionConfig{},
			},
			wantErr: true,
			errMsg:  "invalid Discord leaderboard channel ID: xyz",
		},
		{
			name: "command prefix too long",
			config: GuildConfig{
				Discord: DiscordConfig{
					BotCommandPrefix: "!toolong",
				},
				Permissions: PermissionConfig{},
			},
			wantErr: true,
			errMsg:  "bot command prefix too long (max 5 characters): !toolong",
		},
		{
			name: "valid channel IDs",
			config: GuildConfig{
				Discord: DiscordConfig{
					AnnouncementChannelID: stringPtr("123456789012345678"),
					LeaderboardChannelID:  stringPtr("987654321098765432"),
					BotCommandPrefix:      "!usl",
				},
				Permissions: PermissionConfig{
					AdminRoleIDs:     []string{"111111111111111111"},
					ModeratorRoleIDs: []string{"222222222222222222"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()

			if tt.wantErr {
				if err == nil {
					t.Errorf("GuildConfig.Validate() error = nil, wantErr %v", tt.wantErr)
					return
				}
				if err.Error() != tt.errMsg {
					t.Errorf("GuildConfig.Validate() error = %v, want %v", err, tt.errMsg)
				}
			} else {
				if err != nil {
					t.Errorf("GuildConfig.Validate() error = %v, wantErr %v", err, tt.wantErr)
				}

				// Check that defaults were set
				if tt.config.Discord.BotCommandPrefix == "" {
					t.Error("Expected BotCommandPrefix to be set to default")
				}
			}
		})
	}
}

func TestGuildConfig_HasAdminRole(t *testing.T) {
	config := GuildConfig{
		Permissions: PermissionConfig{
			AdminRoleIDs:     []string{"123456789012345678", "111111111111111111"},
			ModeratorRoleIDs: []string{"987654321098765432"},
		},
	}

	tests := []struct {
		name      string
		userRoles []string
		want      bool
	}{
		{
			name:      "user has admin role",
			userRoles: []string{"123456789012345678", "other-role"},
			want:      true,
		},
		{
			name:      "user has different admin role",
			userRoles: []string{"111111111111111111"},
			want:      true,
		},
		{
			name:      "user has only moderator role",
			userRoles: []string{"987654321098765432"},
			want:      false,
		},
		{
			name:      "user has no matching roles",
			userRoles: []string{"999999999999999999"},
			want:      false,
		},
		{
			name:      "user has no roles",
			userRoles: []string{},
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := config.HasAdminRole(tt.userRoles); got != tt.want {
				t.Errorf("GuildConfig.HasAdminRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGuildConfig_HasModeratorRole(t *testing.T) {
	config := GuildConfig{
		Permissions: PermissionConfig{
			AdminRoleIDs:     []string{"123456789012345678"},
			ModeratorRoleIDs: []string{"987654321098765432", "222222222222222222"},
		},
	}

	tests := []struct {
		name      string
		userRoles []string
		want      bool
	}{
		{
			name:      "user has admin role (should have moderator access)",
			userRoles: []string{"123456789012345678"},
			want:      true,
		},
		{
			name:      "user has moderator role",
			userRoles: []string{"987654321098765432"},
			want:      true,
		},
		{
			name:      "user has different moderator role",
			userRoles: []string{"222222222222222222"},
			want:      true,
		},
		{
			name:      "user has both admin and moderator roles",
			userRoles: []string{"123456789012345678", "987654321098765432"},
			want:      true,
		},
		{
			name:      "user has no matching roles",
			userRoles: []string{"999999999999999999"},
			want:      false,
		},
		{
			name:      "user has no roles",
			userRoles: []string{},
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := config.HasModeratorRole(tt.userRoles); got != tt.want {
				t.Errorf("GuildConfig.HasModeratorRole() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsValidDiscordSnowflake(t *testing.T) {
	tests := []struct {
		name string
		id   string
		want bool
	}{
		{
			name: "valid 18 digit snowflake",
			id:   "123456789012345678",
			want: true,
		},
		{
			name: "valid 19 digit snowflake",
			id:   "1234567890123456789",
			want: true,
		},
		{
			name: "valid 17 digit snowflake",
			id:   "12345678901234567",
			want: true,
		},
		{
			name: "too short",
			id:   "1234567890123456",
			want: false,
		},
		{
			name: "too long",
			id:   "12345678901234567890",
			want: false,
		},
		{
			name: "contains letters",
			id:   "12345678901234567a",
			want: false,
		},
		{
			name: "empty string",
			id:   "",
			want: false,
		},
		{
			name: "only letters",
			id:   "abcdefghijklmnopqr",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isValidDiscordSnowflake(tt.id); got != tt.want {
				t.Errorf("isValidDiscordSnowflake() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDefaultGuildConfig(t *testing.T) {
	config := GetDefaultGuildConfig()

	// Check defaults are set correctly
	if config.Discord.BotCommandPrefix != "!usl" {
		t.Errorf("Expected BotCommandPrefix to be '!usl', got '%s'", config.Discord.BotCommandPrefix)
	}

	if config.Discord.AnnouncementChannelID != nil {
		t.Errorf("Expected AnnouncementChannelID to be nil, got %v", config.Discord.AnnouncementChannelID)
	}

	if config.Discord.LeaderboardChannelID != nil {
		t.Errorf("Expected LeaderboardChannelID to be nil, got %v", config.Discord.LeaderboardChannelID)
	}

	if len(config.Permissions.AdminRoleIDs) != 0 {
		t.Errorf("Expected AdminRoleIDs to be empty, got %v", config.Permissions.AdminRoleIDs)
	}

	if len(config.Permissions.ModeratorRoleIDs) != 0 {
		t.Errorf("Expected ModeratorRoleIDs to be empty, got %v", config.Permissions.ModeratorRoleIDs)
	}

	// Validate the default config
	if err := config.Validate(); err != nil {
		t.Errorf("Default config should be valid, got error: %v", err)
	}
}

// Helper function for creating string pointers in tests
func stringPtr(s string) *string {
	return &s
}
