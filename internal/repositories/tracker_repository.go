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
		OnesCurrentSeasonPeak:     func() *int32 { v := int32(trackerData.OnesCurrentSeasonPeak); return &v }(),
		OnesPreviousSeasonPeak:    func() *int32 { v := int32(trackerData.OnesPreviousSeasonPeak); return &v }(),
		OnesAllTimePeak:           func() *int32 { v := int32(trackerData.OnesAllTimePeak); return &v }(),
		OnesCurrentSeasonGames:    func() *int32 { v := int32(trackerData.OnesCurrentSeasonGames); return &v }(),
		OnesPreviousSeasonGames:   func() *int32 { v := int32(trackerData.OnesPreviousSeasonGames); return &v }(),
		TwosCurrentSeasonPeak:     func() *int32 { v := int32(trackerData.TwosCurrentSeasonPeak); return &v }(),
		TwosPreviousSeasonPeak:    func() *int32 { v := int32(trackerData.TwosPreviousSeasonPeak); return &v }(),
		TwosAllTimePeak:           func() *int32 { v := int32(trackerData.TwosAllTimePeak); return &v }(),
		TwosCurrentSeasonGames:    func() *int32 { v := int32(trackerData.TwosCurrentSeasonGames); return &v }(),
		TwosPreviousSeasonGames:   func() *int32 { v := int32(trackerData.TwosPreviousSeasonGames); return &v }(),
		ThreesCurrentSeasonPeak:   func() *int32 { v := int32(trackerData.ThreesCurrentSeasonPeak); return &v }(),
		ThreesPreviousSeasonPeak:  func() *int32 { v := int32(trackerData.ThreesPreviousSeasonPeak); return &v }(),
		ThreesAllTimePeak:         func() *int32 { v := int32(trackerData.ThreesAllTimePeak); return &v }(),
		ThreesCurrentSeasonGames:  func() *int32 { v := int32(trackerData.ThreesCurrentSeasonGames); return &v }(),
		ThreesPreviousSeasonGames: func() *int32 { v := int32(trackerData.ThreesPreviousSeasonGames); return &v }(),
		Valid:                     &trackerData.Valid,
		LastUpdated:               func() *string { now := time.Now().Format(time.RFC3339); return &now }(),
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
		OnesCurrentSeasonPeak:     func() *int32 { v := int32(trackerData.OnesCurrentSeasonPeak); return &v }(),
		OnesPreviousSeasonPeak:    func() *int32 { v := int32(trackerData.OnesPreviousSeasonPeak); return &v }(),
		OnesAllTimePeak:           func() *int32 { v := int32(trackerData.OnesAllTimePeak); return &v }(),
		OnesCurrentSeasonGames:    func() *int32 { v := int32(trackerData.OnesCurrentSeasonGames); return &v }(),
		OnesPreviousSeasonGames:   func() *int32 { v := int32(trackerData.OnesPreviousSeasonGames); return &v }(),
		TwosCurrentSeasonPeak:     func() *int32 { v := int32(trackerData.TwosCurrentSeasonPeak); return &v }(),
		TwosPreviousSeasonPeak:    func() *int32 { v := int32(trackerData.TwosPreviousSeasonPeak); return &v }(),
		TwosAllTimePeak:           func() *int32 { v := int32(trackerData.TwosAllTimePeak); return &v }(),
		TwosCurrentSeasonGames:    func() *int32 { v := int32(trackerData.TwosCurrentSeasonGames); return &v }(),
		TwosPreviousSeasonGames:   func() *int32 { v := int32(trackerData.TwosPreviousSeasonGames); return &v }(),
		ThreesCurrentSeasonPeak:   func() *int32 { v := int32(trackerData.ThreesCurrentSeasonPeak); return &v }(),
		ThreesPreviousSeasonPeak:  func() *int32 { v := int32(trackerData.ThreesPreviousSeasonPeak); return &v }(),
		ThreesAllTimePeak:         func() *int32 { v := int32(trackerData.ThreesAllTimePeak); return &v }(),
		ThreesCurrentSeasonGames:  func() *int32 { v := int32(trackerData.ThreesCurrentSeasonGames); return &v }(),
		ThreesPreviousSeasonGames: func() *int32 { v := int32(trackerData.ThreesPreviousSeasonGames); return &v }(),
		Valid:                     &trackerData.Valid,
		LastUpdated:               func() *string { now := time.Now().Format(time.RFC3339); return &now }(),
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
