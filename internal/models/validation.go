package models

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ValidationError represents a validation error with field details
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Value   string `json:"value,omitempty"`
}

// ValidationErrors represents multiple validation errors
type ValidationErrors []ValidationError

func (ve ValidationErrors) Error() string {
	if len(ve) == 0 {
		return "validation failed"
	}

	var messages []string
	for _, err := range ve {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}

	return strings.Join(messages, "; ")
}

// RequestParser handles parsing and validation of HTTP request parameters
type RequestParser struct {
	errors ValidationErrors
}

// NewRequestParser creates a new request parser
func NewRequestParser() *RequestParser {
	return &RequestParser{
		errors: make(ValidationErrors, 0),
	}
}

// ParsePaginationParams parses pagination parameters from HTTP request
func (rp *RequestParser) ParsePaginationParams(r *http.Request) *PaginationParams {
	params := &PaginationParams{
		Page:  1,
		Limit: 20,
		Sort:  "created_at",
		Order: "desc",
	}

	// Parse page
	if pageStr := r.URL.Query().Get("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err != nil {
			rp.AddError("page", "must be a valid integer", pageStr)
		} else if page < 1 {
			rp.AddError("page", "must be greater than 0", pageStr)
		} else {
			params.Page = page
		}
	}

	// Parse limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		if limit, err := strconv.Atoi(limitStr); err != nil {
			rp.AddError("limit", "must be a valid integer", limitStr)
		} else if limit < 1 {
			rp.AddError("limit", "must be greater than 0", limitStr)
		} else if limit > 100 {
			rp.AddError("limit", "must not exceed 100", limitStr)
		} else {
			params.Limit = limit
		}
	}

	// Parse sort
	if sort := r.URL.Query().Get("sort"); sort != "" {
		params.Sort = sort
	}

	// Parse order
	if order := r.URL.Query().Get("order"); order != "" {
		if order != "asc" && order != "desc" {
			rp.AddError("order", "must be 'asc' or 'desc'", order)
		} else {
			params.Order = order
		}
	}

	// Parse cursor
	if cursor := r.URL.Query().Get("cursor"); cursor != "" {
		params.Cursor = cursor
	}

	return params
}

// ParseUserFilters parses user filter parameters from HTTP request
func (rp *RequestParser) ParseUserFilters(r *http.Request) *UserFilters {
	filters := &UserFilters{}

	// Parse search
	if search := r.URL.Query().Get("search"); search != "" {
		filters.Search = strings.TrimSpace(search)
	}

	// Parse status
	if status := r.URL.Query().Get("status"); status != "" {
		if !isValidUserStatus(status) {
			rp.AddError("status", "must be 'active', 'inactive', or 'banned'", status)
		} else {
			filters.Status = status
		}
	}

	// Parse MMR range
	if mmrMinStr := r.URL.Query().Get("mmr_min"); mmrMinStr != "" {
		if mmrMin, err := strconv.Atoi(mmrMinStr); err != nil {
			rp.AddError("mmr_min", "must be a valid integer", mmrMinStr)
		} else if mmrMin < 0 {
			rp.AddError("mmr_min", "must be non-negative", mmrMinStr)
		} else {
			filters.MMRMin = &mmrMin
		}
	}

	if mmrMaxStr := r.URL.Query().Get("mmr_max"); mmrMaxStr != "" {
		if mmrMax, err := strconv.Atoi(mmrMaxStr); err != nil {
			rp.AddError("mmr_max", "must be a valid integer", mmrMaxStr)
		} else if mmrMax < 0 {
			rp.AddError("mmr_max", "must be non-negative", mmrMaxStr)
		} else {
			filters.MMRMax = &mmrMax
		}
	}

	// Validate MMR range
	if isInvalidRange(filters.MMRMin, filters.MMRMax) {
		rp.AddError("mmr_range", "mmr_min must be less than or equal to mmr_max",
			fmt.Sprintf("min=%d, max=%d", *filters.MMRMin, *filters.MMRMax))
	}

	// Parse date range
	if createdAfterStr := r.URL.Query().Get("created_after"); createdAfterStr != "" {
		if createdAfter, err := time.Parse(time.RFC3339, createdAfterStr); err != nil {
			// Try parsing date-only format
			if createdAfter, err := time.Parse("2006-01-02", createdAfterStr); err != nil {
				rp.AddError("created_after", "must be a valid date (YYYY-MM-DD or RFC3339)", createdAfterStr)
			} else {
				filters.CreatedAfter = &createdAfter
			}
		} else {
			filters.CreatedAfter = &createdAfter
		}
	}

	if createdBeforeStr := r.URL.Query().Get("created_before"); createdBeforeStr != "" {
		if createdBefore, err := time.Parse(time.RFC3339, createdBeforeStr); err != nil {
			// Try parsing date-only format
			if createdBefore, err := time.Parse("2006-01-02", createdBeforeStr); err != nil {
				rp.AddError("created_before", "must be a valid date (YYYY-MM-DD or RFC3339)", createdBeforeStr)
			} else {
				filters.CreatedBefore = &createdBefore
			}
		} else {
			filters.CreatedBefore = &createdBefore
		}
	}

	// Parse has_trackers
	if hasTrackersStr := r.URL.Query().Get("has_trackers"); hasTrackersStr != "" {
		if hasTrackersStr != "true" && hasTrackersStr != "false" {
			rp.AddError("has_trackers", "must be 'true' or 'false'", hasTrackersStr)
		} else {
			hasTrackers := hasTrackersStr == "true"
			filters.HasTrackers = &hasTrackers
		}
	}

	return filters
}

// ParseTrackerFilters parses tracker filter parameters from HTTP request
func (rp *RequestParser) ParseTrackerFilters(r *http.Request) *TrackerFilters {
	filters := &TrackerFilters{}

	// Parse valid
	if validStr := r.URL.Query().Get("valid"); validStr != "" {
		if validStr != "true" && validStr != "false" {
			rp.AddError("valid", "must be 'true' or 'false'", validStr)
		} else {
			valid := validStr == "true"
			filters.Valid = &valid
		}
	}

	// Parse playlist
	if playlist := r.URL.Query().Get("playlist"); playlist != "" {
		if !isValidPlaylist(playlist) {
			rp.AddError("playlist", "must be 'ones', 'twos', or 'threes'", playlist)
		} else {
			filters.Playlist = playlist
		}
	}

	// Parse peak range
	if peakMinStr := r.URL.Query().Get("peak_min"); peakMinStr != "" {
		if peakMin, err := strconv.Atoi(peakMinStr); err != nil {
			rp.AddError("peak_min", "must be a valid integer", peakMinStr)
		} else if peakMin < 0 {
			rp.AddError("peak_min", "must be non-negative", peakMinStr)
		} else {
			filters.PeakMin = &peakMin
		}
	}

	if peakMaxStr := r.URL.Query().Get("peak_max"); peakMaxStr != "" {
		if peakMax, err := strconv.Atoi(peakMaxStr); err != nil {
			rp.AddError("peak_max", "must be a valid integer", peakMaxStr)
		} else if peakMax < 0 {
			rp.AddError("peak_max", "must be non-negative", peakMaxStr)
		} else {
			filters.PeakMax = &peakMax
		}
	}

	// Validate peak range
	if isInvalidRange(filters.PeakMin, filters.PeakMax) {
		rp.AddError("peak_range", "peak_min must be less than or equal to peak_max",
			fmt.Sprintf("min=%d, max=%d", *filters.PeakMin, *filters.PeakMax))
	}

	// Parse discord_id
	if discordID := r.URL.Query().Get("discord_id"); discordID != "" {
		if len(discordID) < 17 || len(discordID) > 19 {
			rp.AddError("discord_id", "must be 17-19 characters long", discordID)
		} else {
			filters.DiscordID = discordID
		}
	}

	// Parse games range
	if gamesMinStr := r.URL.Query().Get("games_min"); gamesMinStr != "" {
		if gamesMin, err := strconv.Atoi(gamesMinStr); err != nil {
			rp.AddError("games_min", "must be a valid integer", gamesMinStr)
		} else if gamesMin < 0 {
			rp.AddError("games_min", "must be non-negative", gamesMinStr)
		} else {
			filters.GamesMin = &gamesMin
		}
	}

	// Parse date range
	if createdAfterStr := r.URL.Query().Get("created_after"); createdAfterStr != "" {
		if createdAfter, err := time.Parse(time.RFC3339, createdAfterStr); err != nil {
			if createdAfter, err := time.Parse("2006-01-02", createdAfterStr); err != nil {
				rp.AddError("created_after", "must be a valid date (YYYY-MM-DD or RFC3339)", createdAfterStr)
			} else {
				filters.CreatedAfter = &createdAfter
			}
		} else {
			filters.CreatedAfter = &createdAfter
		}
	}

	if createdBeforeStr := r.URL.Query().Get("created_before"); createdBeforeStr != "" {
		if createdBefore, err := time.Parse(time.RFC3339, createdBeforeStr); err != nil {
			if createdBefore, err := time.Parse("2006-01-02", createdBeforeStr); err != nil {
				rp.AddError("created_before", "must be a valid date (YYYY-MM-DD or RFC3339)", createdBeforeStr)
			} else {
				filters.CreatedBefore = &createdBefore
			}
		} else {
			filters.CreatedBefore = &createdBefore
		}
	}

	return filters
}

// AddError adds a validation error
func (rp *RequestParser) AddError(field, message, value string) {
	rp.errors = append(rp.errors, ValidationError{
		Field:   field,
		Message: message,
		Value:   value,
	})
}

// HasErrors returns true if there are validation errors
func (rp *RequestParser) HasErrors() bool {
	return len(rp.errors) > 0
}

// GetErrors returns all validation errors
func (rp *RequestParser) GetErrors() ValidationErrors {
	return rp.errors
}

// ValidateUserSortField validates sort fields for users
func ValidateUserSortField(sort string) bool {
	allowedFields := []string{
		"id", "name", "discord_id", "mmr", "trueskill_mu", "trueskill_sigma",
		"created_at", "updated_at", "trueskill_last_updated",
	}

	for _, field := range allowedFields {
		if sort == field {
			return true
		}
	}

	return false
}

// ValidateTrackerSortField validates sort fields for trackers
func ValidateTrackerSortField(sort string) bool {
	allowedFields := []string{
		"id", "discord_id", "calculated_mmr", "valid", "created_at", "updated_at", "last_updated",
		"ones_current_season_peak", "twos_current_season_peak", "threes_current_season_peak",
	}

	for _, field := range allowedFields {
		if sort == field {
			return true
		}
	}

	return false
}

// Helper functions for validation simplification

// isValidUserStatus checks if the given status is valid for users
func isValidUserStatus(status string) bool {
	validStatuses := []string{"active", "inactive", "banned"}
	for _, validStatus := range validStatuses {
		if status == validStatus {
			return true
		}
	}
	return false
}

// isValidPlaylist checks if the given playlist is valid
func isValidPlaylist(playlist string) bool {
	validPlaylists := []string{"ones", "twos", "threes"}
	for _, validPlaylist := range validPlaylists {
		if playlist == validPlaylist {
			return true
		}
	}
	return false
}

// isInvalidRange checks if min value is greater than max value
func isInvalidRange(min, max *int) bool {
	return min != nil && max != nil && *min > *max
}
