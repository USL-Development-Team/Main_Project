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

func NewV2UsersHandler(userRepo *repositories.UserRepository) *V2UsersHandler {
	return &V2UsersHandler{
		userRepo: userRepo,
	}
}

func (h *V2UsersHandler) HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.handleGetUsers(w, r)
	case http.MethodPost:
		h.handleCreateUser(w, r)
	default:
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, msgMethodNotAllowed, nil)
	}
}

// HandleUsersBulk handles bulk operations on users
func (h *V2UsersHandler) HandleUsersBulk(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.writeErrorResponse(w, http.StatusMethodNotAllowed, msgMethodNotAllowed, nil)
		return
	}

	bulkOperation, err := h.parseBulkOperationRequest(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, msgInvalidRequestBody, map[string]string{"error": err.Error()})
		return
	}

	response, err := h.userRepo.BulkUpdateUsers(bulkOperation)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, msgBulkOperationFailed, map[string]string{"error": err.Error()})
		return
	}

	h.writeJSONResponse(w, http.StatusOK, response)
}

// handleGetUsers handles GET /api/v2/users with pagination and filtering
func (h *V2UsersHandler) handleGetUsers(w http.ResponseWriter, r *http.Request) {
	requestStartTime := time.Now()

	// Parse and validate request parameters
	paginationParams, userFilters, validationErr := h.parseAndValidateGetUsersRequest(r)
	if validationErr != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, msgValidationFailed, validationErr)
		return
	}

	// Retrieve paginated users from repository
	users, paginationMetadata, err := h.userRepo.GetUsersPaginated(paginationParams, userFilters)
	if err != nil {
		h.writeErrorResponse(w, http.StatusInternalServerError, msgFailedToGetUsers, map[string]string{"error": err.Error()})
		return
	}

	// Build and send response
	response := h.buildPaginatedUsersResponse(users, paginationMetadata, paginationParams, userFilters, requestStartTime)
	h.writeJSONResponse(w, http.StatusOK, response)
}

// handleCreateUser handles POST /api/v2/users
func (h *V2UsersHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) {
	userCreateRequest, err := h.parseUserCreateRequest(r)
	if err != nil {
		h.writeErrorResponse(w, http.StatusBadRequest, msgInvalidRequestBody, map[string]string{"error": err.Error()})
		return
	}

	createdUser, err := h.userRepo.CreateUser(*userCreateRequest)
	if err != nil {
		h.handleUserCreationError(w, err, userCreateRequest.DiscordID)
		return
	}

	h.writeJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"user":    createdUser,
		"message": msgUserCreatedSuccessfully,
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

// Helper methods for improved code organization

// parseBulkOperationRequest parses a bulk operation request from HTTP body
func (h *V2UsersHandler) parseBulkOperationRequest(r *http.Request) (*models.BulkOperation, error) {
	var bulkOperation models.BulkOperation
	if err := json.NewDecoder(r.Body).Decode(&bulkOperation); err != nil {
		return nil, err
	}
	return &bulkOperation, nil
}

// parseUserCreateRequest parses a user creation request from HTTP body
func (h *V2UsersHandler) parseUserCreateRequest(r *http.Request) (*models.UserCreateRequest, error) {
	var userCreateRequest models.UserCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&userCreateRequest); err != nil {
		return nil, err
	}
	return &userCreateRequest, nil
}

// parseAndValidateGetUsersRequest parses and validates parameters for GET users request
func (h *V2UsersHandler) parseAndValidateGetUsersRequest(r *http.Request) (*models.PaginationParams, *models.UserFilters, interface{}) {
	requestParser := models.NewRequestParser()
	paginationParams := requestParser.ParsePaginationParams(r)
	userFilters := requestParser.ParseUserFilters(r)

	// Check for validation errors
	if requestParser.HasErrors() {
		return nil, nil, map[string]interface{}{
			"errors": requestParser.GetErrors(),
		}
	}

	// Validate sort field
	if paginationParams.Sort != "" && !models.ValidateUserSortField(paginationParams.Sort) {
		return nil, nil, map[string]string{
			"field":   paginationParams.Sort,
			"allowed": allowedUserSortFields,
		}
	}

	return paginationParams, userFilters, nil
}

// buildPaginatedUsersResponse constructs a paginated response for users
func (h *V2UsersHandler) buildPaginatedUsersResponse(
	users []*models.User,
	paginationMetadata *models.PaginationMetadata,
	paginationParams *models.PaginationParams,
	userFilters *models.UserFilters,
	requestStartTime time.Time,
) models.PaginatedResponse[*models.User] {
	responseMetadata := models.ResponseMetadata{
		Sort:           paginationParams.Sort,
		Order:          paginationParams.Order,
		FiltersApplied: userFilters.ToFiltersApplied(),
		QueryTime:      time.Since(requestStartTime).String(),
	}

	return models.PaginatedResponse[*models.User]{
		Data:       users,
		Pagination: *paginationMetadata,
		Meta:       responseMetadata,
	}
}

// handleUserCreationError handles errors that occur during user creation
func (h *V2UsersHandler) handleUserCreationError(w http.ResponseWriter, err error, discordID string) {
	errorMessage := err.Error()
	expectedConflictMessage := fmt.Sprintf("user with Discord ID %s already exists", discordID)

	if errorMessage == expectedConflictMessage {
		h.writeErrorResponse(w, http.StatusConflict, msgUserAlreadyExists, map[string]string{"discord_id": discordID})
		return
	}

	h.writeErrorResponse(w, http.StatusInternalServerError, msgFailedToCreateUser, map[string]string{"error": errorMessage})
}
