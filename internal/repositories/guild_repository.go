package repositories

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"usl-server/internal/config"
	"usl-server/internal/models"

	"github.com/supabase-community/supabase-go"
)

const (
	// Database table names
	GuildsTable = "guilds"

	// Default values for guild creation
	DefaultGuildActiveStatus = true

	// Error messages
	ErrGuildAlreadyExists     = "guild with Discord ID %s already exists"
	ErrConfigValidationFailed = "config validation failed: %w"
	ErrGuildCreationFailed    = "failed to create guild: %w"
	ErrGuildParsingFailed     = "failed to parse created guild: %w"
	ErrNoGuildReturned        = "no guild returned after creation"
	ErrGuildConversionFailed  = "failed to convert guild: %w"
)

// GuildRepository handles all guild data access operations using Supabase Go client
type GuildRepository struct {
	client *supabase.Client
	config *config.Config
}

func NewGuildRepository(client *supabase.Client, cfg *config.Config) *GuildRepository {
	return &GuildRepository{
		client: client,
		config: cfg,
	}
}

// CreateGuild creates a new guild with default configuration
func (r *GuildRepository) CreateGuild(guildData models.GuildCreateRequest) (*models.Guild, error) {
	// Check if guild already exists
	existingGuild, err := r.FindGuildByDiscordID(guildData.DiscordGuildID)
	if err == nil && existingGuild != nil {
		return nil, fmt.Errorf(ErrGuildAlreadyExists, guildData.DiscordGuildID)
	}

	// Validate configuration
	if err := guildData.Config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	// Use a map to avoid ID conflicts - let Supabase generate the ID
	insertData := map[string]interface{}{
		"discord_guild_id": guildData.DiscordGuildID,
		"name":             guildData.Name,
		"slug":             guildData.Slug,
		"active":           DefaultGuildActiveStatus,
		"config":           guildData.Config,
	}

	data, _, err := r.client.From(GuildsTable).Insert(insertData, false, "", "", "").Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create guild: %w", err)
	}

	var result []models.PublicGuildsSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created guild: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no guild returned after creation")
	}

	createdGuild, err := r.convertToGuild(result[0])
	if err != nil {
		return nil, fmt.Errorf("failed to convert created guild: %w", err)
	}

	return createdGuild, nil
}

// FindGuildByDiscordID finds a guild by Discord guild ID
func (r *GuildRepository) FindGuildByDiscordID(discordGuildID string) (*models.Guild, error) {
	data, _, err := r.client.From("guilds").
		Select("*", "", false).
		Eq("discord_guild_id", discordGuildID).
		Single().
		Execute()

	if err != nil {
		return nil, err
	}

	var result models.PublicGuildsSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse guild data: %w", err)
	}

	guild, err := r.convertToGuild(result)
	if err != nil {
		return nil, fmt.Errorf("failed to convert guild: %w", err)
	}

	return guild, nil
}

// FindGuildByID finds a guild by internal ID
func (r *GuildRepository) FindGuildByID(guildID int64) (*models.Guild, error) {
	data, _, err := r.client.From("guilds").
		Select("*", "", false).
		Eq("id", strconv.FormatInt(guildID, 10)).
		Single().
		Execute()

	if err != nil {
		return nil, err
	}

	var result models.PublicGuildsSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse guild data: %w", err)
	}

	guild, err := r.convertToGuild(result)
	if err != nil {
		return nil, fmt.Errorf("failed to convert guild: %w", err)
	}

	return guild, nil
}

// FindGuildBySlug finds a guild by URL slug
// NOTE: Temporarily using mock guild approach until slug column is added to schema
func (r *GuildRepository) FindGuildBySlug(slug string) (*models.Guild, error) {
	// For USL, return a mock guild since the schema doesn't have slug column yet
	if slug == "usl" {
		return &models.Guild{
			ID:             1,
			DiscordGuildID: "1390537743385231451", // USL Discord Guild ID
			Name:           "USL",
			Slug:           "usl",
			Active:         true,
			Config:         models.GetDefaultGuildConfig(),
			Theme:          models.GetDefaultTheme(),
		}, nil
	}

	// For other slugs, return not found
	return nil, fmt.Errorf("guild with slug '%s' not found", slug)
}

// UpdateGuild updates an existing guild
func (r *GuildRepository) UpdateGuild(guildID int64, guildData models.GuildUpdateRequest) (*models.Guild, error) {
	// Verify guild exists
	_, err := r.FindGuildByID(guildID)
	if err != nil {
		return nil, fmt.Errorf("guild with ID %d not found", guildID)
	}

	// Validate configuration
	if err := guildData.Config.Validate(); err != nil {
		return nil, fmt.Errorf("config validation failed: %w", err)
	}

	updateData := map[string]interface{}{
		"name":   guildData.Name,
		"slug":   guildData.Slug,
		"active": guildData.Active,
		"config": guildData.Config,
	}

	data, _, err := r.client.From("guilds").
		Update(updateData, "", "").
		Eq("id", strconv.FormatInt(guildID, 10)).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("failed to update guild: %w", err)
	}

	var result []models.PublicGuildsSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated guild: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no guild returned after update")
	}

	updatedGuild, err := r.convertToGuild(result[0])
	if err != nil {
		return nil, fmt.Errorf("failed to convert updated guild: %w", err)
	}

	return updatedGuild, nil
}

// UpdateConfig updates only the configuration for a guild
func (r *GuildRepository) UpdateConfig(guildID int64, config *models.GuildConfig) error {
	// Validate configuration
	if err := config.Validate(); err != nil {
		return fmt.Errorf("config validation failed: %w", err)
	}

	updateData := map[string]interface{}{
		"config": config,
	}

	_, _, err := r.client.From("guilds").
		Update(updateData, "", "").
		Eq("id", strconv.FormatInt(guildID, 10)).
		Execute()

	if err != nil {
		return fmt.Errorf("failed to update guild config: %w", err)
	}

	return nil
}

// GetConfig retrieves only the configuration for a guild
func (r *GuildRepository) GetConfig(guildID int64) (*models.GuildConfig, error) {
	data, _, err := r.client.From("guilds").
		Select("config", "", false).
		Eq("id", strconv.FormatInt(guildID, 10)).
		Single().
		Execute()

	if err != nil {
		return nil, err
	}

	var result struct {
		Config json.RawMessage `json:"config"`
	}
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config data: %w", err)
	}

	var config models.GuildConfig
	if err := json.Unmarshal(result.Config, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &config, nil
}

// GetAllGuilds retrieves all guilds with optional active filter
func (r *GuildRepository) GetAllGuilds(activeOnly bool) ([]*models.Guild, error) {
	query := r.client.From("guilds").Select("*", "", false).Order("name", nil)

	if activeOnly {
		query = query.Eq("active", "true")
	}

	data, _, err := query.Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get all guilds: %w", err)
	}

	var result []models.PublicGuildsSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse guilds: %w", err)
	}

	guilds := make([]*models.Guild, 0, len(result))
	for _, guildSelect := range result {
		guild, err := r.convertToGuild(guildSelect)
		if err != nil {
			fmt.Printf("Warning: failed to convert guild %s: %v\n", guildSelect.DiscordGuildId, err)
			continue
		}
		guilds = append(guilds, guild)
	}

	return guilds, nil
}

// DeactivateGuild marks a guild as inactive
func (r *GuildRepository) DeactivateGuild(guildID int64) (*models.Guild, error) {
	_, err := r.FindGuildByID(guildID)
	if err != nil {
		return nil, fmt.Errorf("guild with ID %d not found", guildID)
	}

	updateData := map[string]interface{}{
		"active": false,
	}

	data, _, err := r.client.From("guilds").
		Update(updateData, "", "").
		Eq("id", strconv.FormatInt(guildID, 10)).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("failed to deactivate guild: %w", err)
	}

	var result []models.PublicGuildsSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse deactivated guild: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no guild returned after deactivation")
	}

	deactivatedGuild, err := r.convertToGuild(result[0])
	if err != nil {
		return nil, fmt.Errorf("failed to convert deactivated guild: %w", err)
	}

	return deactivatedGuild, nil
}

// MigrateConfigToRelational is a placeholder for future migration to relational permissions
func (r *GuildRepository) MigrateConfigToRelational() error {
	// This function exists but is unused until needed
	// Extracts JSONB config to proper relational tables
	// You'll thank yourself in 12 months

	// Steps for future implementation:
	// 1. Use proper SQL client for transactions
	// 2. Extract all guild configs
	// 3. INSERT INTO guild_role_permissions SELECT ...
	// 4. UPDATE guilds SET config_migrated = true

	return fmt.Errorf("migration not implemented yet")
}

// Helper function to convert Supabase generated type to internal model
func (r *GuildRepository) convertToGuild(guildSelect models.PublicGuildsSelect) (*models.Guild, error) {
	// Parse timestamps
	createdAt, _ := time.Parse(time.RFC3339, guildSelect.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, guildSelect.UpdatedAt)

	// Parse configuration
	var config models.GuildConfig
	configBytes, err := json.Marshal(guildSelect.Config)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal config interface: %w", err)
	}
	if err := json.Unmarshal(configBytes, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return &models.Guild{
		ID:             guildSelect.Id,
		DiscordGuildID: guildSelect.DiscordGuildId,
		Name:           guildSelect.Name,
		Slug:           guildSelect.Slug,
		Active:         guildSelect.Active,
		Config:         config,
		CreatedAt:      createdAt,
		UpdatedAt:      updatedAt,
	}, nil
}
