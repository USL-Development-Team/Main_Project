package repositories

import (
	"os"
	"testing"
	"usl-server/internal/config"
	"usl-server/internal/models"

	"github.com/supabase-community/supabase-go"
)

func TestGuildRepository_Integration(t *testing.T) {
	// Skip if not running integration tests
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Setup local Supabase client
	supabaseURL := "http://127.0.0.1:54321"
	supabaseKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZS1kZW1vIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImV4cCI6MTk4MzgxMjk5Nn0.EGIM96RAZx35lJzdJsyH-qQwv8Hdp7fsn3W0YpN81IU"

	client, err := supabase.NewClient(supabaseURL, supabaseKey, &supabase.ClientOptions{})
	if err != nil {
		t.Fatalf("Failed to create Supabase client: %v", err)
	}

	cfg := &config.Config{}
	repo := NewGuildRepository(client, cfg)

	t.Run("CreateGuild", func(t *testing.T) {
		// Create test guild with custom config
		testConfig := models.GuildConfig{
			Discord: models.DiscordConfig{
				AnnouncementChannelID: stringPtr("123456789012345678"),
				LeaderboardChannelID:  stringPtr("987654321098765432"),
				BotCommandPrefix:      "!test",
			},
			Permissions: models.PermissionConfig{
				AdminRoleIDs:     []string{"111111111111111111", "222222222222222222"},
				ModeratorRoleIDs: []string{"333333333333333333"},
			},
		}

		guildData := models.GuildCreateRequest{
			DiscordGuildID: "555555555555555555",
			Name:           "Test Guild",
			Config:         testConfig,
		}

		guild, err := repo.CreateGuild(guildData)
		if err != nil {
			t.Fatalf("Failed to create guild: %v", err)
		}

		// Verify guild was created correctly
		if guild.DiscordGuildID != guildData.DiscordGuildID {
			t.Errorf("Expected DiscordGuildID %s, got %s", guildData.DiscordGuildID, guild.DiscordGuildID)
		}

		if guild.Name != guildData.Name {
			t.Errorf("Expected Name %s, got %s", guildData.Name, guild.Name)
		}

		if !guild.Active {
			t.Error("Expected guild to be active")
		}

		// Verify configuration
		if guild.Config.Discord.BotCommandPrefix != "!test" {
			t.Errorf("Expected prefix '!test', got '%s'", guild.Config.Discord.BotCommandPrefix)
		}

		if len(guild.Config.Permissions.AdminRoleIDs) != 2 {
			t.Errorf("Expected 2 admin roles, got %d", len(guild.Config.Permissions.AdminRoleIDs))
		}

		t.Logf("✅ Created guild: ID=%d, DiscordID=%s", guild.ID, guild.DiscordGuildID)

		// Store guild ID for other tests
		testGuildID := guild.ID

		t.Run("FindGuildByDiscordID", func(t *testing.T) {
			foundGuild, err := repo.FindGuildByDiscordID(guildData.DiscordGuildID)
			if err != nil {
				t.Fatalf("Failed to find guild by Discord ID: %v", err)
			}

			if foundGuild.ID != testGuildID {
				t.Errorf("Expected guild ID %d, got %d", testGuildID, foundGuild.ID)
			}

			if foundGuild.Config.Discord.BotCommandPrefix != "!test" {
				t.Errorf("Config not loaded correctly: expected '!test', got '%s'", foundGuild.Config.Discord.BotCommandPrefix)
			}

			t.Logf("✅ Found guild by Discord ID: %s", foundGuild.DiscordGuildID)
		})

		t.Run("FindGuildByID", func(t *testing.T) {
			foundGuild, err := repo.FindGuildByID(testGuildID)
			if err != nil {
				t.Fatalf("Failed to find guild by ID: %v", err)
			}

			if foundGuild.DiscordGuildID != guildData.DiscordGuildID {
				t.Errorf("Expected Discord ID %s, got %s", guildData.DiscordGuildID, foundGuild.DiscordGuildID)
			}

			t.Logf("✅ Found guild by ID: %d", foundGuild.ID)
		})

		t.Run("GetConfig", func(t *testing.T) {
			config, err := repo.GetConfig(testGuildID)
			if err != nil {
				t.Fatalf("Failed to get guild config: %v", err)
			}

			if config.Discord.BotCommandPrefix != "!test" {
				t.Errorf("Expected prefix '!test', got '%s'", config.Discord.BotCommandPrefix)
			}

			if len(config.Permissions.AdminRoleIDs) != 2 {
				t.Errorf("Expected 2 admin roles, got %d", len(config.Permissions.AdminRoleIDs))
			}

			t.Logf("✅ Retrieved guild config successfully")
		})

		t.Run("UpdateConfig", func(t *testing.T) {
			newConfig := models.GuildConfig{
				Discord: models.DiscordConfig{
					AnnouncementChannelID: stringPtr("999888777666555444"),
					LeaderboardChannelID:  nil,   // Test null value
					BotCommandPrefix:      "!up", // Keep under 5 chars
				},
				Permissions: models.PermissionConfig{
					AdminRoleIDs:     []string{"111111111111111111"},                       // Reduced to 1
					ModeratorRoleIDs: []string{"333333333333333333", "444444444444444444"}, // Added 1
				},
			}

			err := repo.UpdateConfig(testGuildID, &newConfig)
			if err != nil {
				t.Fatalf("Failed to update config: %v", err)
			}

			// Verify update
			updatedConfig, err := repo.GetConfig(testGuildID)
			if err != nil {
				t.Fatalf("Failed to get updated config: %v", err)
			}

			if updatedConfig.Discord.BotCommandPrefix != "!up" {
				t.Errorf("Config not updated: expected '!up', got '%s'", updatedConfig.Discord.BotCommandPrefix)
			}

			if updatedConfig.Discord.LeaderboardChannelID != nil {
				t.Errorf("Expected LeaderboardChannelID to be nil, got %v", updatedConfig.Discord.LeaderboardChannelID)
			}

			if len(updatedConfig.Permissions.AdminRoleIDs) != 1 {
				t.Errorf("Expected 1 admin role, got %d", len(updatedConfig.Permissions.AdminRoleIDs))
			}

			if len(updatedConfig.Permissions.ModeratorRoleIDs) != 2 {
				t.Errorf("Expected 2 moderator roles, got %d", len(updatedConfig.Permissions.ModeratorRoleIDs))
			}

			t.Logf("✅ Updated guild config successfully")
		})

		t.Run("GetAllGuilds", func(t *testing.T) {
			allGuilds, err := repo.GetAllGuilds(false)
			if err != nil {
				t.Fatalf("Failed to get all guilds: %v", err)
			}

			if len(allGuilds) == 0 {
				t.Error("Expected at least 1 guild")
			}

			found := false
			for _, g := range allGuilds {
				if g.ID == testGuildID {
					found = true
					break
				}
			}

			if !found {
				t.Error("Test guild not found in all guilds list")
			}

			t.Logf("✅ Retrieved %d guilds", len(allGuilds))
		})

		t.Run("DeactivateGuild", func(t *testing.T) {
			deactivatedGuild, err := repo.DeactivateGuild(testGuildID)
			if err != nil {
				t.Fatalf("Failed to deactivate guild: %v", err)
			}

			if deactivatedGuild.Active {
				t.Error("Expected guild to be inactive")
			}

			// Verify with GetAllGuilds(activeOnly=true)
			activeGuilds, err := repo.GetAllGuilds(true)
			if err != nil {
				t.Fatalf("Failed to get active guilds: %v", err)
			}

			for _, g := range activeGuilds {
				if g.ID == testGuildID {
					t.Error("Deactivated guild found in active guilds list")
				}
			}

			t.Logf("✅ Deactivated guild successfully")
		})
	})

	t.Run("ValidationErrors", func(t *testing.T) {
		// Test invalid Discord snowflake
		invalidGuildData := models.GuildCreateRequest{
			DiscordGuildID: "invalid-id",
			Name:           "Invalid Guild",
			Config:         models.GetDefaultGuildConfig(),
		}

		_, err := repo.CreateGuild(invalidGuildData)
		if err == nil {
			t.Error("Expected validation error for invalid Discord ID")
		}

		t.Logf("✅ Validation correctly rejected invalid Discord ID: %v", err)

		// Test duplicate Discord ID
		validGuildData := models.GuildCreateRequest{
			DiscordGuildID: "666666666666666666",
			Name:           "First Guild",
			Config:         models.GetDefaultGuildConfig(),
		}

		_, err = repo.CreateGuild(validGuildData)
		if err != nil {
			t.Fatalf("Failed to create first guild: %v", err)
		}

		// Try to create another with same Discord ID
		duplicateGuildData := models.GuildCreateRequest{
			DiscordGuildID: "666666666666666666", // Same ID
			Name:           "Duplicate Guild",
			Config:         models.GetDefaultGuildConfig(),
		}

		_, err = repo.CreateGuild(duplicateGuildData)
		if err == nil {
			t.Error("Expected error for duplicate Discord ID")
		}

		t.Logf("✅ Correctly prevented duplicate Discord ID: %v", err)
	})
}

func TestGuildRepository_DefaultConfig(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test")
	}

	// Setup local Supabase client
	supabaseURL := "http://127.0.0.1:54321"
	supabaseKey := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZS1kZW1vIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImV4cCI6MTk4MzgxMjk5Nn0.EGIM96RAZx35lJzdJsyH-qQwv8Hdp7fsn3W0YpN81IU"

	client, err := supabase.NewClient(supabaseURL, supabaseKey, &supabase.ClientOptions{})
	if err != nil {
		t.Fatalf("Failed to create Supabase client: %v", err)
	}

	cfg := &config.Config{}
	repo := NewGuildRepository(client, cfg)

	// Test creating guild with default config
	guildData := models.GuildCreateRequest{
		DiscordGuildID: "777777777777777777",
		Name:           "Default Config Guild",
		Config:         models.GetDefaultGuildConfig(),
	}

	guild, err := repo.CreateGuild(guildData)
	if err != nil {
		t.Fatalf("Failed to create guild with default config: %v", err)
	}

	// Verify defaults
	if guild.Config.Discord.BotCommandPrefix != "!usl" {
		t.Errorf("Expected default prefix '!usl', got '%s'", guild.Config.Discord.BotCommandPrefix)
	}

	if guild.Config.Discord.AnnouncementChannelID != nil {
		t.Errorf("Expected AnnouncementChannelID to be nil, got %v", guild.Config.Discord.AnnouncementChannelID)
	}

	if len(guild.Config.Permissions.AdminRoleIDs) != 0 {
		t.Errorf("Expected empty AdminRoleIDs, got %v", guild.Config.Permissions.AdminRoleIDs)
	}

	t.Logf("✅ Default config applied correctly")
}

// Helper function for creating string pointers in tests
func stringPtr(s string) *string {
	return &s
}

func init() {
	// Set environment variables for testing if not already set
	if os.Getenv("SUPABASE_URL") == "" {
		os.Setenv("SUPABASE_URL", "http://127.0.0.1:54321")
	}
	if os.Getenv("SUPABASE_SERVICE_ROLE_KEY") == "" {
		os.Setenv("SUPABASE_SERVICE_ROLE_KEY", "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJzdXBhYmFzZS1kZW1vIiwicm9sZSI6InNlcnZpY2Vfcm9sZSIsImV4cCI6MTk4MzgxMjk5Nn0.EGIM96RAZx35lJzdJsyH-qQwv8Hdp7fsn3W0YpN81IU")
	}
}
