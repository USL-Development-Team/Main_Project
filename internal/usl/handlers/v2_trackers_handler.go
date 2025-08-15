package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"usl-server/internal/models"
	"usl-server/internal/repositories"
)

// V2TrackersHandler handles modern API requests for trackers with pagination, filtering, and bulk operations
type V2TrackersHandler struct {
	trackerRepo *repositories.TrackerRepository
}

// NewV2TrackersHandler creates a new V2 trackers handler
func NewV2TrackersHandler(trackerRepo *repositories.TrackerRepository) *V2TrackersHandler {
	return &V2TrackersHandler{
		trackerRepo: trackerRepo,
	}
}

// HandleTrackers handles the main /api/v2/trackers endpoint
func (h *V2TrackersHandler) HandleTrackers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getTrackers(w, r)
	case http.MethodPost:
		h.createTracker(w, r)
	default:
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed", nil)
	}
}

// HandleTrackersBulk handles bulk operations on trackers
func (h *V2TrackersHandler) HandleTrackersBulk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	var operation models.BulkOperation
	if err := json.NewDecoder(r.Body).Decode(&operation); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "invalid request body", map[string]string{"error": err.Error()})
		return
	}

	response, err := h.trackerRepo.BulkUpdateTrackers(&operation)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "bulk operation failed", map[string]string{"error": err.Error()})
		return
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// getTrackers handles GET /api/v2/trackers with pagination and filtering
func (h *V2TrackersHandler) getTrackers(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Parse request parameters
	parser := models.NewRequestParser()
	params := parser.ParsePaginationParams(r)
	filters := parser.ParseTrackerFilters(r)

	// Check for validation errors
	if parser.HasErrors() {
		h.writeErrorResponse(w, http.StatusBadRequest, "validation failed", map[string]interface{}{
			"errors": parser.GetErrors(),
		})
		return
	}

	// Validate sort field
	if params.Sort != "" && !models.ValidateTrackerSortField(params.Sort) {
		h.writeErrorResponse(w, http.StatusBadRequest, "invalid sort field", map[string]string{
			"field":   params.Sort,
			"allowed": "id, discord_id, calculated_mmr, valid, created_at, updated_at, last_updated, ones_current_season_peak, twos_current_season_peak, threes_current_season_peak",
		})
		return
	}

	// Get paginated trackers
	trackers, pagination, err := h.trackerRepo.GetTrackersPaginated(params, filters)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "failed to get trackers", map[string]string{"error": err.Error()})
		return
	}

	// Build response metadata
	meta := models.ResponseMetadata{
		Sort:           params.Sort,
		Order:          params.Order,
		FiltersApplied: filters.ToFiltersApplied(),
		QueryTime:      time.Since(startTime).String(),
	}

	// Build paginated response
	response := models.PaginatedResponse[*models.UserTracker]{
		Data:       trackers,
		Pagination: *pagination,
		Meta:       meta,
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// createTracker handles POST /api/v2/trackers
func (h *V2TrackersHandler) createTracker(w http.ResponseWriter, r *http.Request) {
	var trackerData models.TrackerCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&trackerData); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "invalid request body", map[string]string{"error": err.Error()})
		return
	}

	// Create tracker
	tracker, err := h.trackerRepo.CreateTracker(trackerData)
	if err != nil {
		if err.Error() == fmt.Sprintf("tracker with Discord ID %s already exists", trackerData.DiscordID) {
			h.writeErrorResponse(w, http.StatusConflict, "tracker already exists", map[string]string{"discord_id": trackerData.DiscordID})
			return
		}
		h.writeErrorResponse(w, http.StatusInternalServerError, "failed to create tracker", map[string]string{"error": err.Error()})
		return
	}

	h.writeJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"tracker": tracker,
		"message": "tracker created successfully",
	})
}

// writeJSONResponse writes a JSON response with proper headers
func (h *V2TrackersHandler) writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If encoding fails, log error but don't try to write another response
		fmt.Printf("[ERROR] Failed to encode JSON response: %v\n", err)
	}
}

// writeErrorResponse writes a standardized error response
func (h *V2TrackersHandler) writeErrorResponse(w http.ResponseWriter, status int, message string, details interface{}) {
	errorResponse := map[string]interface{}{
		"error":     message,
		"status":    status,
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}

	if details != nil {
		errorResponse["details"] = details
	}

	h.writeJSONResponse(w, status, errorResponse)
}
