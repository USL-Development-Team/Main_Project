package repositories

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"
	"usl-server/internal/config"
	"usl-server/internal/models"

	"github.com/supabase-community/postgrest-go"
	"github.com/supabase-community/supabase-go"
)

// TrackerRepository handles all tracker data access operations using Supabase Go client
// Exactly matches the patterns from JavaScript TrackerRepository
type TrackerRepository struct {
	client *supabase.Client
	config *config.Config
}

// NewTrackerRepository creates a new tracker repository instance
func NewTrackerRepository(client *supabase.Client, cfg *config.Config) *TrackerRepository {
	return &TrackerRepository{
		client: client,
		config: cfg,
	}
}

// CreateTracker creates a new tracker, matching JavaScript createTracker()
func (r *TrackerRepository) CreateTracker(trackerData models.TrackerCreateRequest) (*models.UserTracker, error) {
	// Check for existing tracker by Discord ID
	existingTrackers, err := r.GetTrackersByDiscordID(trackerData.DiscordID, false)
	if err == nil && len(existingTrackers) > 0 {
		return nil, fmt.Errorf("tracker with Discord ID %s already exists", trackerData.DiscordID)
	}

	// Prepare insert data using generated types
	insertData := models.PublicUserTrackersInsert{
		DiscordId:                 trackerData.DiscordID,
		Url:                       trackerData.URL,
		OnesCurrentSeasonPeak:     r.intToInt32Ptr(trackerData.OnesCurrentSeasonPeak),
		OnesPreviousSeasonPeak:    r.intToInt32Ptr(trackerData.OnesPreviousSeasonPeak),
		OnesAllTimePeak:           r.intToInt32Ptr(trackerData.OnesAllTimePeak),
		OnesCurrentSeasonGames:    r.intToInt32Ptr(trackerData.OnesCurrentSeasonGames),
		OnesPreviousSeasonGames:   r.intToInt32Ptr(trackerData.OnesPreviousSeasonGames),
		TwosCurrentSeasonPeak:     r.intToInt32Ptr(trackerData.TwosCurrentSeasonPeak),
		TwosPreviousSeasonPeak:    r.intToInt32Ptr(trackerData.TwosPreviousSeasonPeak),
		TwosAllTimePeak:           r.intToInt32Ptr(trackerData.TwosAllTimePeak),
		TwosCurrentSeasonGames:    r.intToInt32Ptr(trackerData.TwosCurrentSeasonGames),
		TwosPreviousSeasonGames:   r.intToInt32Ptr(trackerData.TwosPreviousSeasonGames),
		ThreesCurrentSeasonPeak:   r.intToInt32Ptr(trackerData.ThreesCurrentSeasonPeak),
		ThreesPreviousSeasonPeak:  r.intToInt32Ptr(trackerData.ThreesPreviousSeasonPeak),
		ThreesAllTimePeak:         r.intToInt32Ptr(trackerData.ThreesAllTimePeak),
		ThreesCurrentSeasonGames:  r.intToInt32Ptr(trackerData.ThreesCurrentSeasonGames),
		ThreesPreviousSeasonGames: r.intToInt32Ptr(trackerData.ThreesPreviousSeasonGames),
		Valid:                     &trackerData.Valid,
		LastUpdated:               r.currentTimeStringPtr(),
	}

	// Insert using Supabase client
	data, _, err := r.client.From("user_trackers").Insert(insertData, false, "", "", "").Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to create tracker: %w", err)
	}

	// Parse the returned data
	var result []models.PublicUserTrackersSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse created tracker: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no tracker returned after creation")
	}

	// Convert to internal UserTracker model
	createdTracker := r.convertToUserTracker(result[0])
	return &createdTracker, nil
}

// GetTrackersByDiscordID finds trackers by Discord ID using Supabase client
func (r *TrackerRepository) GetTrackersByDiscordID(discordID string, validOnly bool) ([]*models.UserTracker, error) {
	query := r.client.From("user_trackers").
		Select("*", "", false).
		Eq("discord_id", discordID).
		Order("created_at", nil)

	if validOnly {
		query = query.Eq("valid", "true")
	}

	data, _, err := query.Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get trackers: %w", err)
	}

	var result []models.PublicUserTrackersSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tracker data: %w", err)
	}

	trackers := make([]*models.UserTracker, len(result))
	for i, trackerSelect := range result {
		tracker := r.convertToUserTracker(trackerSelect)
		trackers[i] = &tracker
	}

	return trackers, nil
}

// UpdateTracker updates an existing tracker using Supabase client
func (r *TrackerRepository) UpdateTracker(trackerID int, trackerData models.TrackerUpdateRequest) (*models.UserTracker, error) {
	// Prepare update data
	updateData := models.PublicUserTrackersUpdate{
		Url:                       &trackerData.URL,
		OnesCurrentSeasonPeak:     r.intToInt32Ptr(trackerData.OnesCurrentSeasonPeak),
		OnesPreviousSeasonPeak:    r.intToInt32Ptr(trackerData.OnesPreviousSeasonPeak),
		OnesAllTimePeak:           r.intToInt32Ptr(trackerData.OnesAllTimePeak),
		OnesCurrentSeasonGames:    r.intToInt32Ptr(trackerData.OnesCurrentSeasonGames),
		OnesPreviousSeasonGames:   r.intToInt32Ptr(trackerData.OnesPreviousSeasonGames),
		TwosCurrentSeasonPeak:     r.intToInt32Ptr(trackerData.TwosCurrentSeasonPeak),
		TwosPreviousSeasonPeak:    r.intToInt32Ptr(trackerData.TwosPreviousSeasonPeak),
		TwosAllTimePeak:           r.intToInt32Ptr(trackerData.TwosAllTimePeak),
		TwosCurrentSeasonGames:    r.intToInt32Ptr(trackerData.TwosCurrentSeasonGames),
		TwosPreviousSeasonGames:   r.intToInt32Ptr(trackerData.TwosPreviousSeasonGames),
		ThreesCurrentSeasonPeak:   r.intToInt32Ptr(trackerData.ThreesCurrentSeasonPeak),
		ThreesPreviousSeasonPeak:  r.intToInt32Ptr(trackerData.ThreesPreviousSeasonPeak),
		ThreesAllTimePeak:         r.intToInt32Ptr(trackerData.ThreesAllTimePeak),
		ThreesCurrentSeasonGames:  r.intToInt32Ptr(trackerData.ThreesCurrentSeasonGames),
		ThreesPreviousSeasonGames: r.intToInt32Ptr(trackerData.ThreesPreviousSeasonGames),
		Valid:                     &trackerData.Valid,
		LastUpdated:               r.currentTimeStringPtr(),
	}

	// Update using Supabase client
	data, _, err := r.client.From("user_trackers").
		Update(updateData, "", "").
		Eq("id", strconv.Itoa(trackerID)).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("failed to update tracker: %w", err)
	}

	var result []models.PublicUserTrackersSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse updated tracker: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no tracker returned after update")
	}

	updatedTracker := r.convertToUserTracker(result[0])
	return &updatedTracker, nil
}

// DeleteTracker marks tracker as invalid using Supabase client
func (r *TrackerRepository) DeleteTracker(trackerID int) (*models.UserTracker, error) {
	// Mark as invalid
	invalid := false
	updateData := models.PublicUserTrackersUpdate{
		Valid: &invalid,
	}

	data, _, err := r.client.From("user_trackers").
		Update(updateData, "", "").
		Eq("id", strconv.Itoa(trackerID)).
		Execute()

	if err != nil {
		return nil, fmt.Errorf("failed to delete tracker: %w", err)
	}

	var result []models.PublicUserTrackersSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse deleted tracker: %w", err)
	}

	if len(result) == 0 {
		return nil, fmt.Errorf("no tracker returned after deletion")
	}

	deletedTracker := r.convertToUserTracker(result[0])
	return &deletedTracker, nil
}

// GetAllTrackers gets all trackers with optional filters using Supabase client
func (r *TrackerRepository) GetAllTrackers(validOnly bool) ([]*models.UserTracker, error) {
	query := r.client.From("user_trackers").Select("*", "", false).Order("created_at", nil)

	if validOnly {
		query = query.Eq("valid", "true")
	}

	data, _, err := query.Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get all trackers: %w", err)
	}

	var result []models.PublicUserTrackersSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse trackers: %w", err)
	}

	trackers := make([]*models.UserTracker, len(result))
	for i, trackerSelect := range result {
		tracker := r.convertToUserTracker(trackerSelect)
		trackers[i] = &tracker
	}

	return trackers, nil
}

// SearchTrackers searches by Discord ID or URL using Supabase client
func (r *TrackerRepository) SearchTrackers(searchTerm string, maxResults int) ([]*models.UserTracker, error) {
	if searchTerm == "" {
		return []*models.UserTracker{}, nil
	}

	if maxResults <= 0 {
		maxResults = 50
	}

	// Supabase supports ilike for case-insensitive partial matching
	searchPattern := "%" + searchTerm + "%"

	data, _, err := r.client.From("user_trackers").
		Select("*", "", false).
		Or(fmt.Sprintf("discord_id.ilike.%s,url.ilike.%s", searchPattern, searchPattern), "").
		Order("created_at", nil).
		Limit(maxResults, "").
		Execute()

	if err != nil {
		return nil, fmt.Errorf("failed to search trackers: %w", err)
	}

	var result []models.PublicUserTrackersSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, fmt.Errorf("failed to parse search results: %w", err)
	}

	trackers := make([]*models.UserTracker, len(result))
	for i, trackerSelect := range result {
		tracker := r.convertToUserTracker(trackerSelect)
		trackers[i] = &tracker
	}

	return trackers, nil
}

// GetTrackerStats calculates tracker statistics using Supabase client
func (r *TrackerRepository) GetTrackerStats() (*models.TrackerStats, error) {
	data, _, err := r.client.From("user_trackers").Select("valid,calculated_mmr", "", false).Execute()
	if err != nil {
		return nil, fmt.Errorf("failed to get tracker data for stats: %w", err)
	}

	var allTrackers []models.PublicUserTrackersSelect
	err = json.Unmarshal(data, &allTrackers)
	if err != nil {
		return nil, fmt.Errorf("failed to parse tracker stats data: %w", err)
	}

	stats := &models.TrackerStats{}
	var totalMMR float64
	var totalGames int

	for _, tracker := range allTrackers {
		stats.TotalTrackers++
		if tracker.Valid {
			stats.ValidTrackers++
			totalMMR += float64(tracker.CalculatedMmr)
			// Calculate total games for this tracker
			games := int(tracker.OnesCurrentSeasonGames + tracker.OnesPreviousSeasonGames +
				tracker.TwosCurrentSeasonGames + tracker.TwosPreviousSeasonGames +
				tracker.ThreesCurrentSeasonGames + tracker.ThreesPreviousSeasonGames)
			totalGames += games
		}
	}

	if stats.ValidTrackers > 0 {
		stats.AverageMMR = totalMMR / float64(stats.ValidTrackers)
		stats.AverageGamesPerTracker = float64(totalGames) / float64(stats.ValidTrackers)
	}

	// Count unique users
	uniqueUsers := make(map[string]bool)
	for _, tracker := range allTrackers {
		uniqueUsers[tracker.DiscordId] = true
	}
	stats.UniqueUsers = len(uniqueUsers)

	return stats, nil
}

// GetTrackersPaginated gets trackers with pagination and filtering using Supabase client
func (r *TrackerRepository) GetTrackersPaginated(params *models.PaginationParams, filters *models.TrackerFilters) ([]*models.UserTracker, *models.PaginationMetadata, error) {

	query := r.client.From("user_trackers").Select("*", "", false)

	// Apply filters
	if filters != nil {
		// Apply valid filter
		if filters.Valid != nil {
			if *filters.Valid {
				query = query.Eq("valid", "true")
			} else {
				query = query.Eq("valid", "false")
			}
		}

		// Apply playlist filter (based on peak values)
		if filters.Playlist != "" {
			switch filters.Playlist {
			case "ones":
				query = query.Gt("ones_current_season_peak", "0")
			case "twos":
				query = query.Gt("twos_current_season_peak", "0")
			case "threes":
				query = query.Gt("threes_current_season_peak", "0")
			}
		}

		// Apply Discord ID filter
		if filters.DiscordID != "" {
			query = query.Eq("discord_id", filters.DiscordID)
		}

		// Apply date filters
		if filters.CreatedAfter != nil {
			query = query.Gte("created_at", filters.CreatedAfter.Format(time.RFC3339))
		}
		if filters.CreatedBefore != nil {
			query = query.Lt("created_at", filters.CreatedBefore.Format(time.RFC3339))
		}
	}

	total, err := r.getTrackerCount(filters)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get tracker count: %w", err)
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
		return nil, nil, fmt.Errorf("failed to get paginated trackers: %w", err)
	}

	var result []models.PublicUserTrackersSelect
	err = json.Unmarshal(data, &result)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to parse paginated trackers: %w", err)
	}

	// Convert to internal model
	trackers := make([]*models.UserTracker, len(result))
	for i, trackerSelect := range result {
		tracker := r.convertToUserTracker(trackerSelect)
		trackers[i] = &tracker
	}

	// Calculate pagination metadata
	pagination := models.CalculatePagination(params, total)

	return trackers, &pagination, nil
}

// getTrackerCount gets the total count of trackers with filters applied
func (r *TrackerRepository) getTrackerCount(filters *models.TrackerFilters) (int64, error) {
	query := r.client.From("user_trackers").Select("id", "count", false)

	// Apply the same filters as the main query
	if filters != nil {
		// Apply valid filter
		if filters.Valid != nil {
			if *filters.Valid {
				query = query.Eq("valid", "true")
			} else {
				query = query.Eq("valid", "false")
			}
		}

		// Apply playlist filter (based on peak values)
		if filters.Playlist != "" {
			switch filters.Playlist {
			case "ones":
				query = query.Gt("ones_current_season_peak", "0")
			case "twos":
				query = query.Gt("twos_current_season_peak", "0")
			case "threes":
				query = query.Gt("threes_current_season_peak", "0")
			}
		}

		// Apply Discord ID filter
		if filters.DiscordID != "" {
			query = query.Eq("discord_id", filters.DiscordID)
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
		return 0, fmt.Errorf("failed to get tracker count: %w", err)
	}

	var countResult []map[string]interface{}
	err = json.Unmarshal(data, &countResult)
	if err != nil {
		return 0, fmt.Errorf("failed to parse tracker count: %w", err)
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

// BulkUpdateTrackers performs bulk updates on trackers
func (r *TrackerRepository) BulkUpdateTrackers(operation *models.BulkOperation) (*models.BulkOperationResponse, error) {
	startTime := time.Now()
	response := &models.BulkOperationResponse{
		Results: make([]models.BulkOperationResult, 0),
		Errors:  make([]string, 0),
	}

	switch operation.Operation {
	case "update":
		return r.bulkUpdateTrackersUpdate(operation, response, startTime)
	case "delete":
		return r.bulkUpdateTrackersDelete(operation, response, startTime)
	default:
		response.Errors = append(response.Errors, fmt.Sprintf("unsupported operation: %s", operation.Operation))
		response.ProcessingTime = time.Since(startTime).String()
		return response, nil
	}
}

// bulkUpdateTrackersUpdate handles bulk tracker updates
func (r *TrackerRepository) bulkUpdateTrackersUpdate(operation *models.BulkOperation, response *models.BulkOperationResponse, startTime time.Time) (*models.BulkOperationResponse, error) {
	// If tracker IDs are specified, update those trackers
	if len(operation.UserIDs) > 0 { // Using UserIDs field to store tracker IDs
		response.TotalRequested = len(operation.UserIDs)

		for _, trackerIDStr := range operation.UserIDs {
			if trackerID, err := strconv.Atoi(trackerIDStr); err != nil {
				response.Failed++
				response.Results = append(response.Results, models.BulkOperationResult{
					ID:     trackerIDStr,
					Status: "failed",
					Error:  "invalid tracker ID",
				})
			} else if err := r.updateSingleTracker(trackerID, operation.Updates); err != nil {
				response.Failed++
				response.Results = append(response.Results, models.BulkOperationResult{
					ID:     trackerIDStr,
					Status: "failed",
					Error:  err.Error(),
				})
			} else {
				response.Successful++
				response.Results = append(response.Results, models.BulkOperationResult{
					ID:     trackerIDStr,
					Status: "success",
				})
			}
		}
	} else {
		// Update trackers based on filters
		// TODO: Implement filter-based bulk updates
		response.Errors = append(response.Errors, "filter-based bulk updates not yet implemented")
	}

	response.ProcessingTime = time.Since(startTime).String()
	return response, nil
}

// bulkUpdateTrackersDelete handles bulk tracker deletion (marking invalid)
func (r *TrackerRepository) bulkUpdateTrackersDelete(operation *models.BulkOperation, response *models.BulkOperationResponse, startTime time.Time) (*models.BulkOperationResponse, error) {
	if len(operation.UserIDs) > 0 { // Using UserIDs field to store tracker IDs
		response.TotalRequested = len(operation.UserIDs)

		for _, trackerIDStr := range operation.UserIDs {
			if trackerID, err := strconv.Atoi(trackerIDStr); err != nil {
				response.Failed++
				response.Results = append(response.Results, models.BulkOperationResult{
					ID:     trackerIDStr,
					Status: "failed",
					Error:  "invalid tracker ID",
				})
			} else if _, err := r.DeleteTracker(trackerID); err != nil {
				response.Failed++
				response.Results = append(response.Results, models.BulkOperationResult{
					ID:     trackerIDStr,
					Status: "failed",
					Error:  err.Error(),
				})
			} else {
				response.Successful++
				response.Results = append(response.Results, models.BulkOperationResult{
					ID:     trackerIDStr,
					Status: "success",
				})
			}
		}
	}

	response.ProcessingTime = time.Since(startTime).String()
	return response, nil
}

// updateSingleTracker updates a single tracker with the provided updates
func (r *TrackerRepository) updateSingleTracker(trackerID int, updates map[string]interface{}) error {
	// Validate and sanitize updates
	allowedFields := map[string]bool{
		"url":                          true,
		"ones_current_season_peak":     true,
		"ones_previous_season_peak":    true,
		"ones_all_time_peak":           true,
		"ones_current_season_games":    true,
		"ones_previous_season_games":   true,
		"twos_current_season_peak":     true,
		"twos_previous_season_peak":    true,
		"twos_all_time_peak":           true,
		"twos_current_season_games":    true,
		"twos_previous_season_games":   true,
		"threes_current_season_peak":   true,
		"threes_previous_season_peak":  true,
		"threes_all_time_peak":         true,
		"threes_current_season_games":  true,
		"threes_previous_season_games": true,
		"valid":                        true,
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

	// Add updated_at and last_updated timestamps
	sanitizedUpdates["updated_at"] = time.Now().Format(time.RFC3339)
	sanitizedUpdates["last_updated"] = time.Now().Format(time.RFC3339)

	// Execute update
	_, _, err := r.client.From("user_trackers").
		Update(sanitizedUpdates, "", "").
		Eq("id", strconv.Itoa(trackerID)).
		Execute()

	if err != nil {
		return fmt.Errorf("failed to update tracker %d: %w", trackerID, err)
	}

	return nil
}

// Helper function to convert Supabase generated type to internal model
func (r *TrackerRepository) convertToUserTracker(trackerSelect models.PublicUserTrackersSelect) models.UserTracker {
	// Parse timestamps
	createdAt, _ := time.Parse(time.RFC3339, trackerSelect.CreatedAt)
	updatedAt, _ := time.Parse(time.RFC3339, trackerSelect.UpdatedAt)
	lastUpdated, _ := time.Parse(time.RFC3339, trackerSelect.LastUpdated)

	return models.UserTracker{
		ID:                        int(trackerSelect.Id),
		DiscordID:                 trackerSelect.DiscordId,
		URL:                       trackerSelect.Url,
		OnesCurrentSeasonPeak:     int(trackerSelect.OnesCurrentSeasonPeak),
		OnesPreviousSeasonPeak:    int(trackerSelect.OnesPreviousSeasonPeak),
		OnesAllTimePeak:           int(trackerSelect.OnesAllTimePeak),
		OnesCurrentSeasonGames:    int(trackerSelect.OnesCurrentSeasonGames),
		OnesPreviousSeasonGames:   int(trackerSelect.OnesPreviousSeasonGames),
		TwosCurrentSeasonPeak:     int(trackerSelect.TwosCurrentSeasonPeak),
		TwosPreviousSeasonPeak:    int(trackerSelect.TwosPreviousSeasonPeak),
		TwosAllTimePeak:           int(trackerSelect.TwosAllTimePeak),
		TwosCurrentSeasonGames:    int(trackerSelect.TwosCurrentSeasonGames),
		TwosPreviousSeasonGames:   int(trackerSelect.TwosPreviousSeasonGames),
		ThreesCurrentSeasonPeak:   int(trackerSelect.ThreesCurrentSeasonPeak),
		ThreesPreviousSeasonPeak:  int(trackerSelect.ThreesPreviousSeasonPeak),
		ThreesAllTimePeak:         int(trackerSelect.ThreesAllTimePeak),
		ThreesCurrentSeasonGames:  int(trackerSelect.ThreesCurrentSeasonGames),
		ThreesPreviousSeasonGames: int(trackerSelect.ThreesPreviousSeasonGames),
		CalculatedMMR:             int(trackerSelect.CalculatedMmr),
		Valid:                     trackerSelect.Valid,
		LastUpdated:               lastUpdated,
		CreatedAt:                 createdAt,
		UpdatedAt:                 updatedAt,
	}
}

// Helper functions for improved code readability

// intToInt32Ptr converts an int to a pointer to int32
func (r *TrackerRepository) intToInt32Ptr(value int) *int32 {
	converted := int32(value)
	return &converted
}

// currentTimeStringPtr returns a pointer to the current time as RFC3339 string
func (r *TrackerRepository) currentTimeStringPtr() *string {
	now := time.Now().Format(time.RFC3339)
	return &now
}
