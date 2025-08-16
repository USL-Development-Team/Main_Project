package repositories

import (
	"encoding/json"
	"fmt"
	"time"
	"usl-server/internal/config"
	"usl-server/internal/models"

	"github.com/supabase-community/postgrest-go"
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

func (r *UserRepository) FindUserByID(userID int64) (*models.User, error) {
	data, _, err := r.client.From("users").
		Select("*", "", false).
		Eq("id", fmt.Sprintf("%d", userID)).
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

// GetUsersPaginated gets users with pagination and filtering using Supabase client
func (r *UserRepository) GetUsersPaginated(params *models.PaginationParams, filters *models.UserFilters) ([]*models.User, *models.PaginationMetadata, error) {

	query := r.client.From("users").Select("*", "", false)

	// Apply filters
	if filters != nil {
		// Apply search filter (name or discord_id)
		if filters.Search != "" {
			searchPattern := "%" + filters.Search + "%"
			query = query.Or(fmt.Sprintf("name.ilike.%s,discord_id.ilike.%s", searchPattern, searchPattern), "")
		}

		// Apply status filter
		if filters.Status != "" {
			switch filters.Status {
			case "active":
				query = query.Eq("active", "true").Eq("banned", "false")
			case "inactive":
				query = query.Eq("active", "false")
			case "banned":
				query = query.Eq("banned", "true")
			}
		}

		// Apply date filters
		if filters.CreatedAfter != nil {
			query = query.Gte("created_at", filters.CreatedAfter.Format(time.RFC3339))
		}
		if filters.CreatedBefore != nil {
			query = query.Lt("created_at", filters.CreatedBefore.Format(time.RFC3339))
		}
	}

	total, err := r.getUserCount(filters)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get user count: %w", err)
	}

	// Apply sorting
	if params.Sort != "" {
		if params.Order == "asc" {
			query = query.Order(params.Sort, &postgrest.OrderOpts{Ascending: true})
		} else {
			query = query.Order(params.Sort, &postgrest.OrderOpts{Ascending: false})
		}
	}

	// Apply pagination
	offset := params.CalculateOffset()
	query = query.Range(offset, offset+params.Limit-1, "")

	// Execute query
	data, _, err := query.Execute()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get paginated users: %w", err)
	}

	var result []models.PublicUsersSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse paginated users: %w", err)
	}

	// Convert to internal model
	users := make([]*models.User, len(result))
	for i, userSelect := range result {
		user := r.convertToUser(userSelect)
		users[i] = &user
	}

	// Calculate pagination metadata
	pagination := models.CalculatePagination(params, total)

	return users, &pagination, nil
}

// getUserCount gets the total count of users with filters applied
func (r *UserRepository) getUserCount(filters *models.UserFilters) (int64, error) {
	query := r.client.From("users").Select("id", "count", false)

	// Apply the same filters as the main query
	if filters != nil {
		// Apply search filter (name or discord_id)
		if filters.Search != "" {
			searchPattern := "%" + filters.Search + "%"
			query = query.Or(fmt.Sprintf("name.ilike.%s,discord_id.ilike.%s", searchPattern, searchPattern), "")
		}

		// Apply status filter
		if filters.Status != "" {
			switch filters.Status {
			case "active":
				query = query.Eq("active", "true").Eq("banned", "false")
			case "inactive":
				query = query.Eq("active", "false")
			case "banned":
				query = query.Eq("banned", "true")
			}
		}

		// Apply date filters
		if filters.CreatedAfter != nil {
			query = query.Gte("created_at", filters.CreatedAfter.Format(time.RFC3339))
		}
		if filters.CreatedBefore != nil {
			query = query.Lt("created_at", filters.CreatedBefore.Format(time.RFC3339))
		}
	}

	data, _, err := query.Execute()
	if err != nil {
		return 0, fmt.Errorf("failed to get user count: %w", err)
	}

	var countResult []map[string]interface{}
	err = json.Unmarshal(data, &countResult)
	if err != nil {
		return 0, fmt.Errorf("failed to parse user count: %w", err)
	}

	if len(countResult) == 0 {
		return 0, nil
	}

	// Extract count from Supabase response
	count, ok := countResult[0]["count"].(float64)
	if !ok {
		return int64(len(countResult)), nil // Fallback to result length
	}

	return int64(count), nil
}

// BulkUpdateUsers performs bulk updates on users
func (r *UserRepository) BulkUpdateUsers(operation *models.BulkOperation) (*models.BulkOperationResponse, error) {
	startTime := time.Now()
	response := &models.BulkOperationResponse{
		Results: make([]models.BulkOperationResult, 0),
		Errors:  make([]string, 0),
	}

	switch operation.Operation {
	case "update":
		return r.bulkUpdateUsersUpdate(operation, response, startTime)
	case "delete":
		return r.bulkUpdateUsersDelete(operation, response, startTime)
	default:
		response.Errors = append(response.Errors, fmt.Sprintf("unsupported operation: %s", operation.Operation))
		response.ProcessingTime = time.Since(startTime).String()
		return response, nil
	}
}

// bulkUpdateUsersUpdate handles bulk user updates
func (r *UserRepository) bulkUpdateUsersUpdate(operation *models.BulkOperation, response *models.BulkOperationResponse, startTime time.Time) (*models.BulkOperationResponse, error) {
	// If user IDs are specified, update those users
	if len(operation.UserIDs) > 0 {
		response.TotalRequested = len(operation.UserIDs)

		for _, userID := range operation.UserIDs {
			if err := r.updateSingleUser(userID, operation.Updates); err != nil {
				response.Failed++
				response.Results = append(response.Results, models.BulkOperationResult{
					ID:     userID,
					Status: "failed",
					Error:  err.Error(),
				})
			} else {
				response.Successful++
				response.Results = append(response.Results, models.BulkOperationResult{
					ID:     userID,
					Status: "success",
				})
			}
		}
	} else {
		// Update users based on filters
		// TODO: Implement filter-based bulk updates
		response.Errors = append(response.Errors, "filter-based bulk updates not yet implemented")
	}

	response.ProcessingTime = time.Since(startTime).String()
	return response, nil
}

// bulkUpdateUsersDelete handles bulk user deletion (marking inactive)
func (r *UserRepository) bulkUpdateUsersDelete(operation *models.BulkOperation, response *models.BulkOperationResponse, startTime time.Time) (*models.BulkOperationResponse, error) {
	if len(operation.UserIDs) > 0 {
		response.TotalRequested = len(operation.UserIDs)

		for _, userID := range operation.UserIDs {
			if _, err := r.DeleteUser(userID); err != nil {
				response.Failed++
				response.Results = append(response.Results, models.BulkOperationResult{
					ID:     userID,
					Status: "failed",
					Error:  err.Error(),
				})
			} else {
				response.Successful++
				response.Results = append(response.Results, models.BulkOperationResult{
					ID:     userID,
					Status: "success",
				})
			}
		}
	}

	response.ProcessingTime = time.Since(startTime).String()
	return response, nil
}

// updateSingleUser updates a single user with the provided updates
func (r *UserRepository) updateSingleUser(discordID string, updates map[string]interface{}) error {
	// Validate and sanitize updates
	allowedFields := map[string]bool{
		"name":   true,
		"active": true,
		"banned": true,
	}

	sanitizedUpdates := make(map[string]interface{})
	for key, value := range updates {
		if allowedFields[key] {
			sanitizedUpdates[key] = value
		}
	}

	if len(sanitizedUpdates) == 0 {
		return fmt.Errorf("no valid fields to update")
	}

	sanitizedUpdates["updated_at"] = time.Now().Format(time.RFC3339)

	// Execute update
	_, _, err := r.client.From("users").
		Update(sanitizedUpdates, "", "").
		Eq("discord_id", discordID).
		Execute()

	if err != nil {
		return fmt.Errorf("failed to update user %s: %w", discordID, err)
	}

	return nil
}
