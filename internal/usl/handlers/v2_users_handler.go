package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"usl-server/internal/models"
	"usl-server/internal/repositories"
)

// V2UsersHandler handles modern API requests for users with pagination, filtering, and bulk operations
type V2UsersHandler struct {
	userRepo *repositories.UserRepository
}

// NewV2UsersHandler creates a new V2 users handler
func NewV2UsersHandler(userRepo *repositories.UserRepository) *V2UsersHandler {
	return &V2UsersHandler{
		userRepo: userRepo,
	}
}

// HandleUsers handles the main /api/v2/users endpoint
func (h *V2UsersHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getUsers(w, r)
	case http.MethodPost:
		h.createUser(w, r)
	default:
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed", nil)
	}
}

// HandleUsersBulk handles bulk operations on users
func (h *V2UsersHandler) HandleUsersBulk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed", nil)
		return
	}

	var operation models.BulkOperation
	if err := json.NewDecoder(r.Body).Decode(&operation); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "invalid request body", map[string]string{"error": err.Error()})
		return
	}

	response, err := h.userRepo.BulkUpdateUsers(&operation)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "bulk operation failed", map[string]string{"error": err.Error()})
		return
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// getUsers handles GET /api/v2/users with pagination and filtering
func (h *V2UsersHandler) getUsers(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()

	// Parse request parameters
	parser := models.NewRequestParser()
	params := parser.ParsePaginationParams(r)
	filters := parser.ParseUserFilters(r)

	// Check for validation errors
	if parser.HasErrors() {
		h.writeErrorResponse(w, http.StatusBadRequest, "validation failed", map[string]interface{}{
			"errors": parser.GetErrors(),
		})
		return
	}

	// Validate sort field
	if params.Sort != "" && !models.ValidateUserSortField(params.Sort) {
		h.writeErrorResponse(w, http.StatusBadRequest, "invalid sort field", map[string]string{
			"field":   params.Sort,
			"allowed": "id, name, discord_id, mmr, trueskill_mu, trueskill_sigma, created_at, updated_at, trueskill_last_updated",
		})
		return
	}

	// Get paginated users
	users, pagination, err := h.userRepo.GetUsersPaginated(params, filters)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, "failed to get users", map[string]string{"error": err.Error()})
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
	response := models.PaginatedResponse[*models.User]{
		Data:       users,
		Pagination: *pagination,
		Meta:       meta,
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// createUser handles POST /api/v2/users
func (h *V2UsersHandler) createUser(w http.ResponseWriter, r *http.Request) {
	var userData models.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&userData); err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, "invalid request body", map[string]string{"error": err.Error()})
		return
	}

	// Create user
	user, err := h.userRepo.CreateUser(userData)
	if err != nil {
		if err.Error() == fmt.Sprintf("user with Discord ID %s already exists", userData.DiscordID) {
			h.writeErrorResponse(w, http.StatusConflict, "user already exists", map[string]string{"discord_id": userData.DiscordID})
			return
		}
		h.writeErrorResponse(w, http.StatusInternalServerError, "failed to create user", map[string]string{"error": err.Error()})
		return
	}

	h.writeJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"user":    user,
		"message": "user created successfully",
	})
}

// writeJSONResponse writes a JSON response with proper headers
func (h *V2UsersHandler) writeJSONResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		// If encoding fails, log error but don't try to write another response
		fmt.Printf("[ERROR] Failed to encode JSON response: %v\n", err)
	}
}

// writeErrorResponse writes a standardized error response
func (h *V2UsersHandler) writeErrorResponse(w http.ResponseWriter, status int, message string, details interface{}) {
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
