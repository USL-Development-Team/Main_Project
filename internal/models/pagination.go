package models

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

// PaginationParams represents pagination request parameters
type PaginationParams struct {
	Page   int    `json:"page" validate:"min=1"`
	Limit  int    `json:"limit" validate:"min=1,max=100"`
	Sort   string `json:"sort"`
	Order  string `json:"order" validate:"oneof=asc desc"`
	Cursor string `json:"cursor,omitempty"`
}

// PaginationMetadata represents pagination response metadata
type PaginationMetadata struct {
	Page       int    `json:"page"`
	Limit      int    `json:"limit"`
	Total      int64  `json:"total"`
	TotalPages int    `json:"total_pages"`
	HasNext    bool   `json:"has_next"`
	HasPrev    bool   `json:"has_prev"`
	NextCursor string `json:"next_cursor,omitempty"`
	PrevCursor string `json:"prev_cursor,omitempty"`
}

// ResponseMetadata represents general response metadata
type ResponseMetadata struct {
	Sort           string   `json:"sort"`
	Order          string   `json:"order"`
	FiltersApplied []string `json:"filters_applied"`
	QueryTime      string   `json:"query_time,omitempty"`
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse[T any] struct {
	Data       []T                `json:"data"`
	Pagination PaginationMetadata `json:"pagination"`
	Meta       ResponseMetadata   `json:"meta"`
}

// UserFilters represents filtering parameters for users
type UserFilters struct {
	Search        string     `json:"search,omitempty"`
	Status        string     `json:"status,omitempty" validate:"omeof=active inactive banned"`
	MMRMin        *int       `json:"mmr_min,omitempty" validate:"omitempty,min=0"`
	MMRMax        *int       `json:"mmr_max,omitempty" validate:"omitempty,min=0"`
	CreatedAfter  *time.Time `json:"created_after,omitempty"`
	CreatedBefore *time.Time `json:"created_before,omitempty"`
	HasTrackers   *bool      `json:"has_trackers,omitempty"`
}

// TrackerFilters represents filtering parameters for trackers
type TrackerFilters struct {
	Valid         *bool      `json:"valid,omitempty"`
	Playlist      string     `json:"playlist,omitempty" validate:"omeof=ones twos threes"`
	PeakMin       *int       `json:"peak_min,omitempty" validate:"omitempty,min=0"`
	PeakMax       *int       `json:"peak_max,omitempty" validate:"omitempty,min=0"`
	DiscordID     string     `json:"discord_id,omitempty"`
	GamesMin      *int       `json:"games_min,omitempty" validate:"omitempty,min=0"`
	CreatedAfter  *time.Time `json:"created_after,omitempty"`
	CreatedBefore *time.Time `json:"created_before,omitempty"`
}

// BulkOperation represents a bulk operation request
type BulkOperation struct {
	Operation string                 `json:"operation" validate:"required,oneof=create update delete"`
	Filters   map[string]interface{} `json:"filters,omitempty"`
	Updates   map[string]interface{} `json:"updates,omitempty"`
	UserIDs   []string               `json:"user_ids,omitempty"`
	Data      []interface{}          `json:"data,omitempty"`
}

// BulkOperationResult represents the result of a single item in a bulk operation
type BulkOperationResult struct {
	ID     string      `json:"id"`
	Status string      `json:"status" validate:"oneof=success failed skipped"`
	Data   interface{} `json:"data,omitempty"`
	Error  string      `json:"error,omitempty"`
}

// BulkOperationResponse represents the response for a bulk operation
type BulkOperationResponse struct {
	TotalRequested int                   `json:"total_requested"`
	Successful     int                   `json:"successful"`
	Failed         int                   `json:"failed"`
	Skipped        int                   `json:"skipped"`
	Results        []BulkOperationResult `json:"results"`
	Errors         []string              `json:"errors,omitempty"`
	ProcessingTime string                `json:"processing_time"`
}

// NewPaginationParams creates PaginationParams with defaults
func NewPaginationParams(page, limit int, sort, order string) *PaginationParams {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}
	if sort == "" {
		sort = "created_at"
	}
	if order == "" {
		order = "desc"
	}

	return &PaginationParams{
		Page:  page,
		Limit: limit,
		Sort:  sort,
		Order: order,
	}
}

// CalculateOffset calculates the database offset for pagination
func (p *PaginationParams) CalculateOffset() int {
	return (p.Page - 1) * p.Limit
}

// ValidateSort validates and sanitizes sort field names
func (p *PaginationParams) ValidateSort(allowedFields []string) error {
	if p.Sort == "" {
		return nil
	}

	for _, field := range allowedFields {
		if p.Sort == field {
			return nil
		}
	}

	return fmt.Errorf("invalid sort field: %s", p.Sort)
}

// CalculatePagination calculates pagination metadata
func CalculatePagination(params *PaginationParams, total int64) PaginationMetadata {
	totalPages := int((total + int64(params.Limit) - 1) / int64(params.Limit))

	return PaginationMetadata{
		Page:       params.Page,
		Limit:      params.Limit,
		Total:      total,
		TotalPages: totalPages,
		HasNext:    params.Page < totalPages,
		HasPrev:    params.Page > 1,
	}
}

// ToFiltersApplied converts filter struct to string slice for metadata
func (uf *UserFilters) ToFiltersApplied() []string {
	var filters []string

	if uf.Search != "" {
		filters = append(filters, fmt.Sprintf("search=%s", uf.Search))
	}
	if uf.Status != "" {
		filters = append(filters, fmt.Sprintf("status=%s", uf.Status))
	}
	if uf.MMRMin != nil {
		filters = append(filters, fmt.Sprintf("mmr_min=%d", *uf.MMRMin))
	}
	if uf.MMRMax != nil {
		filters = append(filters, fmt.Sprintf("mmr_max=%d", *uf.MMRMax))
	}
	if uf.HasTrackers != nil {
		filters = append(filters, fmt.Sprintf("has_trackers=%t", *uf.HasTrackers))
	}

	return filters
}

// ToFiltersApplied converts filter struct to string slice for metadata
func (tf *TrackerFilters) ToFiltersApplied() []string {
	var filters []string

	if tf.Valid != nil {
		filters = append(filters, fmt.Sprintf("valid=%t", *tf.Valid))
	}
	if tf.Playlist != "" {
		filters = append(filters, fmt.Sprintf("playlist=%s", tf.Playlist))
	}
	if tf.PeakMin != nil {
		filters = append(filters, fmt.Sprintf("peak_min=%d", *tf.PeakMin))
	}
	if tf.PeakMax != nil {
		filters = append(filters, fmt.Sprintf("peak_max=%d", *tf.PeakMax))
	}
	if tf.DiscordID != "" {
		filters = append(filters, fmt.Sprintf("discord_id=%s", tf.DiscordID))
	}

	return filters
}

// ParseIntPointer safely parses a string to *int
func ParseIntPointer(s string) *int {
	if s == "" {
		return nil
	}
	if val, err := strconv.Atoi(s); err == nil {
		return &val
	}
	return nil
}

// ParseBoolPointer safely parses a string to *bool
func ParseBoolPointer(s string) *bool {
	if s == "" {
		return nil
	}
	val := strings.ToLower(s) == "true"
	return &val
}

// ParseTimePointer safely parses a string to *time.Time
func ParseTimePointer(s string) *time.Time {
	if s == "" {
		return nil
	}
	if val, err := time.Parse(time.RFC3339, s); err == nil {
		return &val
	}
	// Try parsing without timezone
	if val, err := time.Parse("2006-01-02", s); err == nil {
		return &val
	}
	return nil
}
