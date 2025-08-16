package services

import (
	"fmt"
	"strconv"
	"time"
	"usl-server/internal/models"
)

// DataTransformationService handles data transformations between different formats and structures.
// Centralizes data mapping logic to ensure consistency across the application.
//
// Service Responsibilities:
// - Transform raw sheet data to structured objects
// - Map between different data representations
// - Validate data structure integrity
// - Provide consistent data formatting
//
// Exact port of JavaScript DataTransformationService
type DataTransformationService struct{}

// TrackerData represents structured tracker data for calculations
type TrackerData struct {
	DiscordID           string    `json:"discordId"`
	URL                 string    `json:"url"`
	OnesCurrentPeak     int       `json:"onesCurrentPeak"`
	OnesPreviousPeak    int       `json:"onesPreviousPeak"`
	OnesAllTimePeak     int       `json:"onesAllTimePeak"`
	OnesCurrentGames    int       `json:"onesCurrentGames"`
	OnesPreviousGames   int       `json:"onesPreviousGames"`
	TwosCurrentPeak     int       `json:"twosCurrentPeak"`
	TwosPreviousPeak    int       `json:"twosPreviousPeak"`
	TwosAllTimePeak     int       `json:"twosAllTimePeak"`
	TwosCurrentGames    int       `json:"twosCurrentGames"`
	TwosPreviousGames   int       `json:"twosPreviousGames"`
	ThreesCurrentPeak   int       `json:"threesCurrentPeak"`
	ThreesPreviousPeak  int       `json:"threesPreviousPeak"`
	ThreesAllTimePeak   int       `json:"threesAllTimePeak"`
	ThreesCurrentGames  int       `json:"threesCurrentGames"`
	ThreesPreviousGames int       `json:"threesPreviousGames"`
	LastUpdated         time.Time `json:"lastUpdated"`
}

// NewDataTransformationService creates a new data transformation service
func NewDataTransformationService() *DataTransformationService {
	return &DataTransformationService{}
}

// TransformRowDataToTracker transforms raw row data to structured tracker data object
// Exact port of JavaScript transformRowDataToTracker() function
func (s *DataTransformationService) TransformRowDataToTracker(rowData []interface{}) (*TrackerData, error) {
	if len(rowData) < 17 {
		return nil, fmt.Errorf("insufficient columns in row data, expected at least 17, got %d", len(rowData))
	}

	// Helper function to safely parse interface{} to string
	safeString := func(val interface{}) string {
		if val == nil {
			return ""
		}
		return fmt.Sprintf("%v", val)
	}

	var lastUpdated time.Time
	if rowData[17] != nil {
		if dateStr := safeString(rowData[17]); dateStr != "" {
			if parsed, err := time.Parse(time.RFC3339, dateStr); err == nil {
				lastUpdated = parsed
			} else {
				lastUpdated = time.Now()
			}
		} else {
			lastUpdated = time.Now()
		}
	} else {
		lastUpdated = time.Now()
	}

	return &TrackerData{
		DiscordID:           safeString(rowData[0]),
		URL:                 safeString(rowData[1]),
		OnesCurrentPeak:     s.safeParseNumber(rowData[2], 0),
		OnesPreviousPeak:    s.safeParseNumber(rowData[3], 0),
		OnesAllTimePeak:     s.safeParseNumber(rowData[4], 0),
		OnesCurrentGames:    s.safeParseNumber(rowData[5], 0),
		OnesPreviousGames:   s.safeParseNumber(rowData[6], 0),
		TwosCurrentPeak:     s.safeParseNumber(rowData[7], 0),
		TwosPreviousPeak:    s.safeParseNumber(rowData[8], 0),
		TwosAllTimePeak:     s.safeParseNumber(rowData[9], 0),
		TwosCurrentGames:    s.safeParseNumber(rowData[10], 0),
		TwosPreviousGames:   s.safeParseNumber(rowData[11], 0),
		ThreesCurrentPeak:   s.safeParseNumber(rowData[12], 0),
		ThreesPreviousPeak:  s.safeParseNumber(rowData[13], 0),
		ThreesAllTimePeak:   s.safeParseNumber(rowData[14], 0),
		ThreesCurrentGames:  s.safeParseNumber(rowData[15], 0),
		ThreesPreviousGames: s.safeParseNumber(rowData[16], 0),
		LastUpdated:         lastUpdated,
	}, nil
}

// PrepareTrackerDataForCalculation prepares tracker object for calculations
// Exact port of JavaScript prepareTrackerDataForCalculation() function
func (s *DataTransformationService) PrepareTrackerDataForCalculation(trackerObject *models.Tracker) (*TrackerData, error) {
	if trackerObject == nil {
		return nil, fmt.Errorf("tracker object is required")
	}

	return &TrackerData{
		DiscordID:           trackerObject.DiscordID,
		URL:                 trackerObject.URL,
		OnesCurrentPeak:     trackerObject.OnesCurrentSeasonPeak,
		OnesPreviousPeak:    trackerObject.OnesPreviousSeasonPeak,
		OnesAllTimePeak:     trackerObject.OnesAllTimePeak,
		OnesCurrentGames:    trackerObject.OnesCurrentSeasonGames,
		OnesPreviousGames:   trackerObject.OnesPreviousSeasonGames,
		TwosCurrentPeak:     trackerObject.TwosCurrentSeasonPeak,
		TwosPreviousPeak:    trackerObject.TwosPreviousSeasonPeak,
		TwosAllTimePeak:     trackerObject.TwosAllTimePeak,
		TwosCurrentGames:    trackerObject.TwosCurrentSeasonGames,
		TwosPreviousGames:   trackerObject.TwosPreviousSeasonGames,
		ThreesCurrentPeak:   trackerObject.ThreesCurrentSeasonPeak,
		ThreesPreviousPeak:  trackerObject.ThreesPreviousSeasonPeak,
		ThreesAllTimePeak:   trackerObject.ThreesAllTimePeak,
		ThreesCurrentGames:  trackerObject.ThreesCurrentSeasonGames,
		ThreesPreviousGames: trackerObject.ThreesPreviousSeasonGames,
		LastUpdated:         trackerObject.LastUpdated,
	}, nil
}

// TransformRowDataToUser transforms raw row data to structured user object
// Exact port of JavaScript transformRowDataToUser() function
func (s *DataTransformationService) TransformRowDataToUser(rowData []interface{}) (*models.User, error) {
	if len(rowData) < 6 {
		return nil, fmt.Errorf("insufficient columns in user row data, expected at least 6, got %d", len(rowData))
	}

	// Helper function to safely parse interface{} to string
	safeString := func(val interface{}) string {
		if val == nil {
			return ""
		}
		return fmt.Sprintf("%v", val)
	}

	// Helper function to safely parse interface{} to bool
	safeBool := func(val interface{}) bool {
		if val == nil {
			return false
		}
		switch v := val.(type) {
		case bool:
			return v
		case string:
			return v == "true" || v == "TRUE" || v == "1"
		case float64:
			return v != 0
		case int:
			return v != 0
		default:
			return false
		}
	}

	// Helper function to safely parse interface{} to float64
	safeFloat := func(val interface{}, defaultVal float64) float64 {
		if val == nil {
			return defaultVal
		}
		switch v := val.(type) {
		case float64:
			return v
		case int:
			return float64(v)
		case string:
			if f, err := strconv.ParseFloat(v, 64); err == nil {
				return f
			}
		}
		return defaultVal
	}

	var createdAt, updatedAt, trueskillLastUpdated time.Time
	now := time.Now()

	if len(rowData) > 8 && rowData[8] != nil {
		if dateStr := safeString(rowData[8]); dateStr != "" {
			if parsed, err := time.Parse(time.RFC3339, dateStr); err == nil {
				trueskillLastUpdated = parsed
			} else {
				trueskillLastUpdated = now
			}
		} else {
			trueskillLastUpdated = now
		}
	} else {
		trueskillLastUpdated = now
	}

	return &models.User{
		Name:                 safeString(rowData[0]),
		DiscordID:            safeString(rowData[1]),
		Active:               safeBool(rowData[2]),
		Banned:               safeBool(rowData[3]),
		MMR:                  s.safeParseNumber(rowData[4], 0),
		TrueSkillMu:          safeFloat(rowData[5], 1000.0),
		TrueSkillSigma:       safeFloat(rowData[6], 8.333),
		TrueSkillLastUpdated: trueskillLastUpdated,
		CreatedAt:            createdAt,
		UpdatedAt:            updatedAt,
	}, nil
}

// ValidateTrackerData validates tracker data structure integrity
// Exact port of JavaScript validateTrackerData() function
func (s *DataTransformationService) ValidateTrackerData(data *TrackerData) error {
	if data == nil {
		return fmt.Errorf("tracker data is nil")
	}

	if data.DiscordID == "" {
		return fmt.Errorf("discord ID is required")
	}

	// Validate MMR values are non-negative
	mmrFields := []struct {
		value int
		name  string
	}{
		{data.OnesCurrentPeak, "onesCurrentPeak"},
		{data.OnesPreviousPeak, "onesPreviousPeak"},
		{data.OnesAllTimePeak, "onesAllTimePeak"},
		{data.TwosCurrentPeak, "twosCurrentPeak"},
		{data.TwosPreviousPeak, "twosPreviousPeak"},
		{data.TwosAllTimePeak, "twosAllTimePeak"},
		{data.ThreesCurrentPeak, "threesCurrentPeak"},
		{data.ThreesPreviousPeak, "threesPreviousPeak"},
		{data.ThreesAllTimePeak, "threesAllTimePeak"},
	}

	for _, field := range mmrFields {
		if field.value < 0 {
			return fmt.Errorf("%s cannot be negative: %d", field.name, field.value)
		}
	}

	// Validate games values are non-negative
	gameFields := []struct {
		value int
		name  string
	}{
		{data.OnesCurrentGames, "onesCurrentGames"},
		{data.OnesPreviousGames, "onesPreviousGames"},
		{data.TwosCurrentGames, "twosCurrentGames"},
		{data.TwosPreviousGames, "twosPreviousGames"},
		{data.ThreesCurrentGames, "threesCurrentGames"},
		{data.ThreesPreviousGames, "threesPreviousGames"},
	}

	for _, field := range gameFields {
		if field.value < 0 {
			return fmt.Errorf("%s cannot be negative: %d", field.name, field.value)
		}
	}

	return nil
}

// GetTrackerDataStats returns statistics about tracker data
// Exact port of JavaScript getTrackerDataStats() function
func (s *DataTransformationService) GetTrackerDataStats(data *TrackerData) map[string]interface{} {
	if data == nil {
		return map[string]interface{}{
			"error": "tracker data is nil",
		}
	}

	totalGames := data.OnesCurrentGames + data.OnesPreviousGames +
		data.TwosCurrentGames + data.TwosPreviousGames +
		data.ThreesCurrentGames + data.ThreesPreviousGames

	maxCurrentPeak := data.OnesCurrentPeak
	if data.TwosCurrentPeak > maxCurrentPeak {
		maxCurrentPeak = data.TwosCurrentPeak
	}
	if data.ThreesCurrentPeak > maxCurrentPeak {
		maxCurrentPeak = data.ThreesCurrentPeak
	}

	maxAllTimePeak := data.OnesAllTimePeak
	if data.TwosAllTimePeak > maxAllTimePeak {
		maxAllTimePeak = data.TwosAllTimePeak
	}
	if data.ThreesAllTimePeak > maxAllTimePeak {
		maxAllTimePeak = data.ThreesAllTimePeak
	}

	activePlaylistsCount := 0
	if data.OnesCurrentGames > 0 || data.OnesPreviousGames > 0 {
		activePlaylistsCount++
	}
	if data.TwosCurrentGames > 0 || data.TwosPreviousGames > 0 {
		activePlaylistsCount++
	}
	if data.ThreesCurrentGames > 0 || data.ThreesPreviousGames > 0 {
		activePlaylistsCount++
	}

	return map[string]interface{}{
		"totalGames":           totalGames,
		"maxCurrentPeak":       maxCurrentPeak,
		"maxAllTimePeak":       maxAllTimePeak,
		"activePlaylistsCount": activePlaylistsCount,
		"hasData":              totalGames > 0,
		"lastUpdated":          data.LastUpdated,
	}
}

// safeParseNumber safely parses interface{} to int with fallback
// Exact port of JavaScript _safeParseNumber() function
func (s *DataTransformationService) safeParseNumber(val interface{}, defaultVal int) int {
	if val == nil {
		return defaultVal
	}

	switch v := val.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float32:
		return int(v)
	case float64:
		return int(v)
	case string:
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return int(f)
		}
	}

	return defaultVal
}
