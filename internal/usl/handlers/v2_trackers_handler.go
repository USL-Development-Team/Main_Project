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

func NewV2TrackersHandler(trackerRepo *repositories.TrackerRepository) *V2TrackersHandler {
	return &V2TrackersHandler{
		trackerRepo: trackerRepo,
	}
}

func (h *V2TrackersHandler) HandleTrackers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetTrackers(w, r)
	case http.MethodPost:
		h.handleCreateTracker(w, r)
	default:
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, msgMethodNotAllowed, nil)
	}
}

// HandleTrackersBulk handles bulk operations on trackers
func (h *V2TrackersHandler) HandleTrackersBulk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, msgMethodNotAllowed, nil)
		return
	}

	bulkOperation, err := h.parseBulkOperationRequest(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, msgInvalidRequestBody, map[string]string{"error": err.Error()})
		return
	}

	response, err := h.trackerRepo.BulkUpdateTrackers(bulkOperation)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, msgBulkOperationFailed, map[string]string{"error": err.Error()})
		return
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// handleGetTrackers handles GET /api/v2/trackers with pagination and filtering
func (h *V2TrackersHandler) handleGetTrackers(w http.ResponseWriter, r *http.Request) {
	requestStartTime := time.Now()

	// Parse and validate request parameters
	paginationParams, trackerFilters, validationErr := h.parseAndValidateGetTrackersRequest(r)
	if validationErr != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, msgValidationFailed, validationErr)
		return
	}

	// Retrieve paginated trackers from repository
	trackers, paginationMetadata, err := h.trackerRepo.GetTrackersPaginated(paginationParams, trackerFilters)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, msgFailedToGetTrackers, map[string]string{"error": err.Error()})
		return
	}

	// Build and send response
	response := h.buildPaginatedTrackersResponse(trackers, paginationMetadata, paginationParams, trackerFilters, requestStartTime)
	h.writeJSONResponse(w, http.StatusOK, response)
}

// handleCreateTracker handles POST /api/v2/trackers
func (h *V2TrackersHandler) handleCreateTracker(w http.ResponseWriter, r *http.Request) {
	trackerCreateRequest, err := h.parseTrackerCreateRequest(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, msgInvalidRequestBody, map[string]string{"error": err.Error()})
		return
	}

	createdTracker, err := h.trackerRepo.CreateTracker(*trackerCreateRequest)
	if err != nil {
		h.handleTrackerCreationError(w, err, trackerCreateRequest.DiscordID)
		return
	}

	h.writeJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"tracker": createdTracker,
		"message": msgTrackerCreatedSuccessfully,
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

// Helper methods for improved code organization

// parseBulkOperationRequest parses a bulk operation request from HTTP body
func (h *V2TrackersHandler) parseBulkOperationRequest(r *http.Request) (*models.BulkOperation, error) {
	var bulkOperation models.BulkOperation
	if err := json.NewDecoder(r.Body).Decode(&bulkOperation); err != nil {
		return nil, err
	}
	return &bulkOperation, nil
}

// parseTrackerCreateRequest parses a tracker creation request from HTTP body
func (h *V2TrackersHandler) parseTrackerCreateRequest(r *http.Request) (*models.TrackerCreateRequest, error) {
	var trackerCreateRequest models.TrackerCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&trackerCreateRequest); err != nil {
		return nil, err
	}
	return &trackerCreateRequest, nil
}

// parseAndValidateGetTrackersRequest parses and validates parameters for GET trackers request
func (h *V2TrackersHandler) parseAndValidateGetTrackersRequest(r *http.Request) (*models.PaginationParams, *models.TrackerFilters, interface{}) {
	requestParser := models.NewRequestParser()
	paginationParams := requestParser.ParsePaginationParams(r)
	trackerFilters := requestParser.ParseTrackerFilters(r)

	// Check for validation errors
	if requestParser.HasErrors() {
		return nil, nil, map[string]interface{}{
			"errors": requestParser.GetErrors(),
		}
	}

	// Validate sort field
	if paginationParams.Sort != "" && !models.ValidateTrackerSortField(paginationParams.Sort) {
		return nil, nil, map[string]string{
			"field":   paginationParams.Sort,
			"allowed": allowedTrackerSortFields,
		}
	}

	return paginationParams, trackerFilters, nil
}

// buildPaginatedTrackersResponse constructs a paginated response for trackers
func (h *V2TrackersHandler) buildPaginatedTrackersResponse(
	trackers []*models.Tracker,
	paginationMetadata *models.PaginationMetadata,
	paginationParams *models.PaginationParams,
	trackerFilters *models.TrackerFilters,
	requestStartTime time.Time,
) models.PaginatedResponse[*models.Tracker] {
	responseMetadata := models.ResponseMetadata{
		Sort:           paginationParams.Sort,
		Order:          paginationParams.Order,
		FiltersApplied: trackerFilters.ToFiltersApplied(),
		QueryTime:      time.Since(requestStartTime).String(),
	}

	return models.PaginatedResponse[*models.Tracker]{
		Data:       trackers,
		Pagination: *paginationMetadata,
		Meta:       responseMetadata,
	}
}

// handleTrackerCreationError handles errors that occur during tracker creation
func (h *V2TrackersHandler) handleTrackerCreationError(w http.ResponseWriter, err error, discordID string) {
	errorMessage := err.Error()
	expectedConflictMessage := fmt.Sprintf("tracker with Discord ID %s already exists", discordID)

	if errorMessage == expectedConflictMessage {
		h.writeErrorResponse(w, http.StatusConflict, msgTrackerAlreadyExists, map[string]string{"discord_id": discordID})
		return
	}

	h.writeErrorResponse(w, http.StatusInternalServerError, msgFailedToCreateTracker, map[string]string{"error": errorMessage})
}
