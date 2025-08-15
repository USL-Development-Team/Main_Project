package repositories

import (
	"encoding/json"
	"fmt"
	"time"
	"usl-server/internal/config"
	"usl-server/internal/models"

	"github.com/supabase-community/supabase-go"
)

// UserRepository handles all user data access operations using Supabase Go client
// Exactly matches the patterns from JavaScript UserRepository
type UserRepository struct {
	client *supabase.Client
	config *config.Config
}

func NewUserRepository(client *supabase.Client, cfg *config.Config) *UserRepository {
	return &UserRepository{
		client: client,
		config: cfg,
	}
}

func (r *UserRepository) CreateUser(userData models.UserCreateRequest) (*models.User, error) {
	existingUser, err := r.FindUserByDiscordID(userData.DiscordID)
	if err == nil && existingUser != nil {
		return nil, fmt.Errorf("user with Discord ID %s already exists", userData.DiscordID)
	}

	insertData := models.PublicUsersInsert{
		Name:      userData.Name,
		DiscordId: userData.DiscordID,
		Active:    &userData.Active,
		Banned:    &userData.Banned,
	}

	data, _, err := r.client.From("users").Insert(insertData, false, "", "", "").Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	var result []models.PublicUsersSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created user: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no user returned after creation")
	}

	createdUser := r.convertToUser(result[0])
	return &createdUser, nil
}

func (r *UserRepository) FindUserByDiscordID(discordID string) (*models.User, error) {
	data, _, err := r.client.From("users").
		Select("*", "", false).
		Eq("discord_id", discordID).
		Single().
		Execute()

	if err != nil {
		return nil, err
	}

	var result []models.PublicUsersSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user data: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("user not found")
	}

	user := r.convertToUser(result[0])
	return &user, nil
}

func (r *UserRepository) UpdateUser(discordID string, userData models.UserUpdateRequest) (*models.User, error) {
	_, err := r.FindUserByDiscordID(discordID)
	if err != nil {
		return nil, fmt.Errorf("user with Discord ID %s not found", discordID)
	}

	// Restricted update: only name, active, and banned status (matches Google Sheets pattern)
	// Discord ID is immutable in this context
	updateData := models.PublicUsersUpdate{
		Name:   &userData.Name,
		Active: &userData.Active,
		Banned: &userData.Banned,
	}

	data, _, err := r.client.From("users").
		Update(updateData, "", "").
		Eq("discord_id", discordID).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	var result []models.PublicUsersSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated user: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no user returned after update")
	}

	updatedUser := r.convertToUser(result[0])
	return &updatedUser, nil
}

// DeleteUser marks user as inactive using Supabase client
func (r *UserRepository) DeleteUser(discordID string) (*models.User, error) {
	_, err := r.FindUserByDiscordID(discordID)
	if err != nil {
		return nil, fmt.Errorf("user with Discord ID %s not found", discordID)
	}

	inactive := false
	updateData := models.PublicUsersUpdate{
		Active: &inactive,
	}

	data, _, err := r.client.From("users").
		Update(updateData, "", "").
		Eq("discord_id", discordID).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}

	var result []models.PublicUsersSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse deleted user: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no user returned after deletion")
	}

	deletedUser := r.convertToUser(result[0])
	return &deletedUser, nil
}

// GetAllUsers gets all users with optional active filter using Supabase client
func (r *UserRepository) GetAllUsers(activeOnly bool) ([]*models.User, error) {
	query := r.client.From("users").Select("*", "", false).Order("name", nil)

	if activeOnly {
		query = query.Eq("active", "true")
	}

	data, _, err := query.Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get all users: %w", err)
	}

	var result []models.PublicUsersSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse users: %w", err)
	}

	users := make([]*models.User, len(result))
	for i, userSelect := range result {
		user := r.convertToUser(userSelect)
		users[i] = &user
	}

	return users, nil
}

// SearchUsers searches by name or Discord ID using Supabase client
func (r *UserRepository) SearchUsers(searchTerm string, maxResults int) ([]*models.User, error) {
	if searchTerm == "" {
		return []*models.User{}, nil
	}

	if maxResults <= 0 {
		maxResults = 50
	}

	// Supabase supports ilike for case-insensitive partial matching
	searchPattern := "%" + searchTerm + "%"

	data, _, err := r.client.From("users").
		Select("*", "", false).
		Or(fmt.Sprintf("name.ilike.%s,discord_id.ilike.%s", searchPattern, searchPattern), "").
		Order("name", nil).
		Limit(maxResults, "").
		Execute()

	if err != nil {
		return nil, fmt.Errorf("failed to search users: %w", err)
	}

	var result []models.PublicUsersSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	users := make([]*models.User, len(result))
	for i, userSelect := range result {
		user := r.convertToUser(userSelect)
		users[i] = &user
	}

	return users, nil
}

// UpdateUserTrueSkill updates a single user's TrueSkill values using Supabase client
// TODO: Update this function to work with player_effective_mmr table in new schema
func (r *UserRepository) UpdateUserTrueSkill(discordID string, trueskillMu, trueskillSigma float64, lastUpdated *time.Time) error {
	// This function needs to be updated to work with the new schema
	// where MMR data is stored in player_effective_mmr table
	return fmt.Errorf("UpdateUserTrueSkill not yet implemented for new schema")
}

// BatchUpdateTrueSkill updates TrueSkill values for multiple users
func (r *UserRepository) BatchUpdateTrueSkill(updates []models.User) (int, error) {
	successCount := 0

	for _, update := range updates {
		err := r.UpdateUserTrueSkill(update.DiscordID, update.TrueSkillMu, update.TrueSkillSigma, &update.TrueSkillLastUpdated)
		if err != nil {
			fmt.Printf("Failed to update TrueSkill for user %s: %v\n", update.DiscordID, err)
			continue
		}
		successCount++
	}

	return successCount, nil
}

// GetUserStats calculates user statistics using Supabase client
func (r *UserRepository) GetUserStats() (*models.UserStats, error) {
	data, _, err := r.client.From("users").Select("active,banned,mmr,trueskill_mu", "", false).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get user data for stats: %w", err)
	}

	var allUsers []models.PublicUsersSelect
	err = json.Unmarshal(data, &allUsers)
	if err != nil {
		return nil, fmt.Errorf("failed to parse user stats data: %w", err)
	}

	stats := &models.UserStats{}

	for _, user := range allUsers {
		stats.TotalUsers++
		if user.Active {
			stats.ActiveUsers++
			// TODO: Calculate average MMR and TrueSkill from player_effective_mmr table
		}
		if user.Banned {
			stats.BannedUsers++
		}
	}

	// TODO: Query player_effective_mmr table to get MMR statistics
	stats.AverageMMR = 0
	stats.AverageTrueSkillMu = 0

	return stats, nil
}

// Helper function to convert Supabase generated type to internal model
func (r *UserRepository) convertToUser(userSelect models.PublicUsersSelect) models.User {
	// Parse timestamps
	createdAt, _ := time.Parse(time.RFC3339, userSelect.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, userSelect.UpdatedAt)

	return models.User{
		ID:        int(userSelect.Id),
		Name:      userSelect.Name,
		DiscordID: userSelect.DiscordId,
		Active:    userSelect.Active,
		Banned:    userSelect.Banned,
		// MMR fields are now in player_effective_mmr table
		MMR:                  0,           // TODO: Query from player_effective_mmr
		TrueSkillMu:          0,           // TODO: Query from player_effective_mmr
		TrueSkillSigma:       0,           // TODO: Query from player_effective_mmr
		TrueSkillLastUpdated: time.Time{}, // TODO: Query from player_effective_mmr
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
	}
}
