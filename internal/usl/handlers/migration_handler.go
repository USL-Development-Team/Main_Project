package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"
	"usl-server/internal/config"
	"usl-server/internal/logger"
	"usl-server/internal/models"
	"usl-server/internal/services"
	"usl-server/internal/usl"
)

const (
	// USL Configuration
	USLDiscordGuildID = "1390537743385231451"

	// Rocket League Business Rules
	MinMMR             = 0
	MaxMMR             = 3000  // SSL is around 1900-2000, allow buffer for edge cases
	MaxGames           = 10000 // Reasonable season game limit
	MinDiscordIDLength = 17
	MaxDiscordIDLength = 19

	// Validation Error Codes
	ValidationCodeRequired      = "required"
	ValidationCodeInvalidFormat = "invalid_format"
	ValidationCodeOutOfRange    = "out_of_range"
	ValidationCodeLogicalError  = "logical_error"
	ValidationCodeInvalidURL    = "invalid_url"
	ValidationCodeNoData        = "no_data"
)

type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

type ValidationResult struct {
	IsValid bool              `json:"is_valid"`
	Errors  []ValidationError `json:"errors"`
}

type FormField string

const (
	// Basic Entity Fields
	FormFieldID        FormField = "id"
	FormFieldDiscordID FormField = "discord_id"
	FormFieldURL       FormField = "url"
	FormFieldName      FormField = "name"
	FormFieldActive    FormField = "active"
	FormFieldBanned    FormField = "banned"
	FormFieldValid     FormField = "valid"

	// 1v1 MMR Fields
	FormFieldOnesCurrentPeak   FormField = "ones_current_peak"
	FormFieldOnesPreviousPeak  FormField = "ones_previous_peak"
	FormFieldOnesAllTimePeak   FormField = "ones_all_time_peak"
	FormFieldOnesCurrentGames  FormField = "ones_current_games"
	FormFieldOnesPreviousGames FormField = "ones_previous_games"

	// 2v2 MMR Fields
	FormFieldTwosCurrentPeak   FormField = "twos_current_peak"
	FormFieldTwosPreviousPeak  FormField = "twos_previous_peak"
	FormFieldTwosAllTimePeak   FormField = "twos_all_time_peak"
	FormFieldTwosCurrentGames  FormField = "twos_current_games"
	FormFieldTwosPreviousGames FormField = "twos_previous_games"

	// 3v3 MMR Fields
	FormFieldThreesCurrentPeak   FormField = "threes_current_peak"
	FormFieldThreesPreviousPeak  FormField = "threes_previous_peak"
	FormFieldThreesAllTimePeak   FormField = "threes_all_time_peak"
	FormFieldThreesCurrentGames  FormField = "threes_current_games"
	FormFieldThreesPreviousGames FormField = "threes_previous_games"
)

// TemplateName represents typed template names
type TemplateName string

const (
	TemplateUSLUsers          TemplateName = "users-list-page"
	TemplateUSLUsersTable     TemplateName = "users-table-fragment"
	TemplateUSLUserDetail     TemplateName = "user-detail-page"
	TemplateUSLTrackers       TemplateName = "trackers-list-page"
	TemplateUSLTrackerDetail  TemplateName = "tracker-detail-page"
	TemplateUSLTrackerNew     TemplateName = "tracker-new-page"
	TemplateUSLTrackerEdit    TemplateName = "tracker-edit-page"
	TemplateUSLAdminDashboard TemplateName = "admin-dashboard-page"
)

// Validation metrics and monitoring structures
type ValidationMetrics struct {
	TotalValidations      int64            `json:"total_validations"`
	SuccessfulValidations int64            `json:"successful_validations"`
	FailedValidations     int64            `json:"failed_validations"`
	ErrorsByType          map[string]int64 `json:"errors_by_type"`
	ErrorsByField         map[string]int64 `json:"errors_by_field"`
	SecurityIncidents     int64            `json:"security_incidents"`
	LastReset             time.Time        `json:"last_reset"`
}

// ValidationEvent represents a validation event for structured logging
type ValidationEvent struct {
	Type           string                 `json:"type"` // "success", "failure", "security_incident"
	DiscordID      string                 `json:"discord_id,omitempty"`
	URL            string                 `json:"url,omitempty"`
	Errors         []ValidationError      `json:"errors,omitempty"`
	Duration       time.Duration          `json:"duration"`
	Timestamp      time.Time              `json:"timestamp"`
	UserAgent      string                 `json:"user_agent,omitempty"`
	RemoteAddr     string                 `json:"remote_addr,omitempty"`
	FormFields     map[string]interface{} `json:"form_fields,omitempty"`
	SecurityReason string                 `json:"security_reason,omitempty"`
}

// Global validation metrics (in production, this would be in a proper metrics store)
var (
	validationMetrics = &ValidationMetrics{
		ErrorsByType:  make(map[string]int64),
		ErrorsByField: make(map[string]int64),
		LastReset:     time.Now(),
	}
	metricsMutex     = &sync.RWMutex{}
	validationLogger = logger.NewLogger("validation")
)

func parseIntField(value string) int {
	if value == "" {
		return 0
	}
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}
	return 0
}

func recordValidationMetrics(event *ValidationEvent) {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()

	validationMetrics.TotalValidations++

	switch event.Type {
	case "success":
		validationMetrics.SuccessfulValidations++
	case "failure":
		validationMetrics.FailedValidations++
		for _, err := range event.Errors {
			validationMetrics.ErrorsByType[err.Code]++
			validationMetrics.ErrorsByField[err.Field]++
		}
	case "security_incident":
		validationMetrics.SecurityIncidents++
		validationMetrics.FailedValidations++
	}
}

func logValidationEvent(event *ValidationEvent) {
	baseArgs := []any{
		"event_type", event.Type,
		"duration", event.Duration,
		"timestamp", event.Timestamp,
	}

	if event.UserAgent != "" {
		baseArgs = append(baseArgs, "user_agent", event.UserAgent)
	}
	if event.RemoteAddr != "" {
		baseArgs = append(baseArgs, "remote_addr", event.RemoteAddr)
	}

	if len(event.Errors) > 0 {
		baseArgs = append(baseArgs, "error_count", len(event.Errors))
		errorCodes := make([]string, len(event.Errors))
		errorFields := make([]string, len(event.Errors))
		for i, err := range event.Errors {
			errorCodes[i] = err.Code
			errorFields[i] = err.Field
		}
		baseArgs = append(baseArgs, "error_codes", errorCodes)
		baseArgs = append(baseArgs, "error_fields", errorFields)
	}

	// Security incident logging
	if event.Type == "security_incident" {
		securityArgs := append(baseArgs, "security_reason", event.SecurityReason)
		// Log with higher severity for security incidents
		validationLogger.Warn("Security incident detected during validation", securityArgs...)
	} else if event.Type == "failure" {
		validationLogger.Info("Validation failed", baseArgs...)
	} else {
		validationLogger.Debug("Validation completed", baseArgs...)
	}
}

func getValidationMetrics() ValidationMetrics {
	metricsMutex.RLock()
	defer metricsMutex.RUnlock()

	metrics := ValidationMetrics{
		TotalValidations:      validationMetrics.TotalValidations,
		SuccessfulValidations: validationMetrics.SuccessfulValidations,
		FailedValidations:     validationMetrics.FailedValidations,
		SecurityIncidents:     validationMetrics.SecurityIncidents,
		LastReset:             validationMetrics.LastReset,
		ErrorsByType:          make(map[string]int64),
		ErrorsByField:         make(map[string]int64),
	}

	for k, v := range validationMetrics.ErrorsByType {
		metrics.ErrorsByType[k] = v
	}
	for k, v := range validationMetrics.ErrorsByField {
		metrics.ErrorsByField[k] = v
	}

	return metrics
}

// detectSecurityIncident checks if a validation failure represents a security threat
func detectSecurityIncident(r *http.Request, validation *ValidationResult) *string {
	// Check for common attack patterns
	formData := r.Form

	// SQL injection patterns
	for _, values := range formData {
		for _, value := range values {
			if strings.Contains(strings.ToLower(value), "drop table") ||
				strings.Contains(strings.ToLower(value), "union select") ||
				strings.Contains(strings.ToLower(value), "insert into") ||
				strings.Contains(strings.ToLower(value), "delete from") {
				reason := "SQL injection attempt detected"
				return &reason
			}
		}
	}

	// XSS patterns
	for _, values := range formData {
		for _, value := range values {
			if strings.Contains(value, "<script") ||
				strings.Contains(value, "javascript:") ||
				strings.Contains(value, "onload=") ||
				strings.Contains(value, "onerror=") {
				reason := "XSS attempt detected"
				return &reason
			}
		}
	}

	// Buffer overflow attempts (extremely long inputs)
	for _, values := range formData {
		for _, value := range values {
			if len(value) > 1000 {
				reason := "Buffer overflow attempt detected"
				return &reason
			}
		}
	}

	return nil
}

// parseUserID safely converts a string to a user ID with proper error handling
func (h *MigrationHandler) parseUserID(userIDStr string) (int64, error) {
	if userIDStr == "" {
		return 0, fmt.Errorf("user ID cannot be empty")
	}
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid user ID format: %v", err)
	}
	return userID, nil
}

// getFormValue safely extracts a typed form value
func (h *MigrationHandler) getFormValue(r *http.Request, field FormField) string {
	return r.FormValue(string(field))
}

// getFormBoolValue safely extracts a boolean form value
func (h *MigrationHandler) getFormBoolValue(r *http.Request, field FormField) bool {
	return h.getFormValue(r, field) == "true"
}

// getFormIntValue safely extracts an integer form value
func (h *MigrationHandler) getFormIntValue(r *http.Request, field FormField) int {
	return parseIntField(h.getFormValue(r, field))
}

// handleMethodNotAllowed returns a standardized method not allowed error
func (h *MigrationHandler) handleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

// handleInvalidFormData returns a standardized form data error
func (h *MigrationHandler) handleInvalidFormData(w http.ResponseWriter, err error) {
	log.Printf("[USL-HANDLER] Form parsing error: %v", err)
	http.Error(w, "Invalid form data", http.StatusBadRequest)
}

// handleDatabaseError logs and returns a standardized database error
func (h *MigrationHandler) handleDatabaseError(w http.ResponseWriter, operation string, err error) {
	log.Printf("[USL-HANDLER] Database error during %s: %v", operation, err)
	http.Error(w, fmt.Sprintf("Failed to %s", operation), http.StatusInternalServerError)
}

// handleInvalidID returns a standardized invalid ID error
func (h *MigrationHandler) handleInvalidID(w http.ResponseWriter, idType string) {
	http.Error(w, fmt.Sprintf("%s is required", idType), http.StatusBadRequest)
}

// handleParseError returns a standardized parsing error
func (h *MigrationHandler) handleParseError(w http.ResponseWriter, fieldName string) {
	http.Error(w, fmt.Sprintf("Invalid %s", fieldName), http.StatusBadRequest)
}

// buildTrackerFromForm creates a USLUserTracker from form data
func (h *MigrationHandler) buildTrackerFromForm(r *http.Request) *usl.USLUserTracker {
	discordID := h.getFormValue(r, FormFieldDiscordID)
	url := h.getFormValue(r, FormFieldURL)
	log.Printf("[USL-HANDLER] Building tracker from form: Discord=%s, URL=%s", discordID, url)

	return &usl.USLUserTracker{
		DiscordID:                       discordID,
		URL:                             url,
		OnesCurrentSeasonPeak:           h.getFormIntValue(r, FormFieldOnesCurrentPeak),
		OnesPreviousSeasonPeak:          h.getFormIntValue(r, FormFieldOnesPreviousPeak),
		OnesAllTimePeak:                 h.getFormIntValue(r, FormFieldOnesAllTimePeak),
		OnesCurrentSeasonGamesPlayed:    h.getFormIntValue(r, FormFieldOnesCurrentGames),
		OnesPreviousSeasonGamesPlayed:   h.getFormIntValue(r, FormFieldOnesPreviousGames),
		TwosCurrentSeasonPeak:           h.getFormIntValue(r, FormFieldTwosCurrentPeak),
		TwosPreviousSeasonPeak:          h.getFormIntValue(r, FormFieldTwosPreviousPeak),
		TwosAllTimePeak:                 h.getFormIntValue(r, FormFieldTwosAllTimePeak),
		TwosCurrentSeasonGamesPlayed:    h.getFormIntValue(r, FormFieldTwosCurrentGames),
		TwosPreviousSeasonGamesPlayed:   h.getFormIntValue(r, FormFieldTwosPreviousGames),
		ThreesCurrentSeasonPeak:         h.getFormIntValue(r, FormFieldThreesCurrentPeak),
		ThreesPreviousSeasonPeak:        h.getFormIntValue(r, FormFieldThreesPreviousPeak),
		ThreesAllTimePeak:               h.getFormIntValue(r, FormFieldThreesAllTimePeak),
		ThreesCurrentSeasonGamesPlayed:  h.getFormIntValue(r, FormFieldThreesCurrentGames),
		ThreesPreviousSeasonGamesPlayed: h.getFormIntValue(r, FormFieldThreesPreviousGames),
		Valid:                           h.getFormBoolValue(r, FormFieldValid),
	}
}

// Comprehensive validation system

// validateTracker performs comprehensive validation on a tracker
// validateTrackerWithMetrics performs validation with metrics collection and security monitoring
func (h *MigrationHandler) validateTrackerWithMetrics(r *http.Request, tracker *usl.USLUserTracker) ValidationResult {
	startTime := time.Now()

	// Perform core validation
	validation := h.validateTracker(tracker)

	// Create validation event
	event := &ValidationEvent{
		Timestamp:  startTime,
		Duration:   time.Since(startTime),
		UserAgent:  r.Header.Get("User-Agent"),
		RemoteAddr: r.RemoteAddr,
	}

	// Extract form fields for security analysis (non-sensitive data only)
	if r.Form != nil {
		event.FormFields = make(map[string]interface{})
		for key, values := range r.Form {
			// Only include field names and lengths, not actual values for privacy
			if len(values) > 0 {
				event.FormFields[key] = map[string]interface{}{
					"length": len(values[0]),
					"count":  len(values),
				}
			}
		}
	}

	if validation.IsValid {
		event.Type = "success"
	} else {
		event.Errors = validation.Errors

		// Check for security incidents
		if securityReason := detectSecurityIncident(r, &validation); securityReason != nil {
			event.Type = "security_incident"
			event.SecurityReason = *securityReason
		} else {
			event.Type = "failure"
		}
	}

	// Record metrics and log event
	recordValidationMetrics(event)
	logValidationEvent(event)

	return validation
}

// validateTracker performs comprehensive validation of tracker data with improved structured errors
func (h *MigrationHandler) validateTracker(tracker *usl.USLUserTracker) ValidationResult {
	var errors []ValidationError

	if tracker.DiscordID == "" {
		errors = append(errors, ValidationError{
			Field:   "discord_id",
			Message: "Discord ID is required",
			Code:    ValidationCodeRequired,
		})
	} else if !isValidDiscordID(tracker.DiscordID) {
		errors = append(errors, ValidationError{
			Field:   "discord_id",
			Message: "Discord ID must be 17-19 digits",
			Code:    ValidationCodeInvalidFormat,
		})
	}

	if tracker.URL == "" {
		errors = append(errors, ValidationError{
			Field:   "url",
			Message: "Tracker URL is required",
			Code:    ValidationCodeRequired,
		})
	} else if !isValidTrackerURL(tracker.URL) {
		errors = append(errors, ValidationError{
			Field:   "url",
			Message: "Invalid tracker URL format",
			Code:    ValidationCodeInvalidURL,
		})
	}

	errors = append(errors, h.validatePlaylistMMR("1v1", tracker.OnesCurrentSeasonPeak, tracker.OnesPreviousSeasonPeak, tracker.OnesAllTimePeak)...)
	errors = append(errors, h.validatePlaylistMMR("2v2", tracker.TwosCurrentSeasonPeak, tracker.TwosPreviousSeasonPeak, tracker.TwosAllTimePeak)...)
	errors = append(errors, h.validatePlaylistMMR("3v3", tracker.ThreesCurrentSeasonPeak, tracker.ThreesPreviousSeasonPeak, tracker.ThreesAllTimePeak)...)

	errors = append(errors, h.validateGamesPlayed("1v1", tracker.OnesCurrentSeasonGamesPlayed, tracker.OnesPreviousSeasonGamesPlayed)...)
	errors = append(errors, h.validateGamesPlayed("2v2", tracker.TwosCurrentSeasonGamesPlayed, tracker.TwosPreviousSeasonGamesPlayed)...)
	errors = append(errors, h.validateGamesPlayed("3v3", tracker.ThreesCurrentSeasonGamesPlayed, tracker.ThreesPreviousSeasonGamesPlayed)...)
	if h.hasNoPlaylistData(tracker) {
		errors = append(errors, ValidationError{
			Field:   "general",
			Message: "Tracker must have data for at least one playlist",
			Code:    ValidationCodeNoData,
		})
	}

	return ValidationResult{
		IsValid: len(errors) == 0,
		Errors:  errors,
	}
}

// validatePlaylistMMR validates MMR values for a specific playlist with improved error messages
func (h *MigrationHandler) validatePlaylistMMR(playlist string, current, previous, allTime int) []ValidationError {
	var errors []ValidationError
	fieldPrefix := strings.ToLower(playlist)

	if current > 0 && (current < MinMMR || current > MaxMMR) {
		errors = append(errors, ValidationError{
			Field:   fieldPrefix + "_current_peak",
			Message: fmt.Sprintf("%s current season MMR must be between %d and %d", playlist, MinMMR, MaxMMR),
			Code:    ValidationCodeOutOfRange,
		})
	}

	if previous > 0 && (previous < MinMMR || previous > MaxMMR) {
		errors = append(errors, ValidationError{
			Field:   fieldPrefix + "_previous_peak",
			Message: fmt.Sprintf("%s previous season MMR must be between %d and %d", playlist, MinMMR, MaxMMR),
			Code:    ValidationCodeOutOfRange,
		})
	}

	if allTime > 0 && (allTime < MinMMR || allTime > MaxMMR) {
		errors = append(errors, ValidationError{
			Field:   fieldPrefix + "_all_time_peak",
			Message: fmt.Sprintf("%s all-time MMR must be between %d and %d", playlist, MinMMR, MaxMMR),
			Code:    ValidationCodeOutOfRange,
		})
	}

	if allTime > 0 && current > 0 && allTime < current {
		errors = append(errors, ValidationError{
			Field:   fieldPrefix + "_all_time_peak",
			Message: fmt.Sprintf("%s all-time MMR should not be less than current season MMR", playlist),
			Code:    ValidationCodeLogicalError,
		})
	}

	if allTime > 0 && previous > 0 && allTime < previous {
		errors = append(errors, ValidationError{
			Field:   fieldPrefix + "_all_time_peak",
			Message: fmt.Sprintf("%s all-time MMR should not be less than previous season MMR", playlist),
			Code:    ValidationCodeLogicalError,
		})
	}

	return errors
}

// validateGamesPlayed validates games played values for a specific playlist with improved error messages
func (h *MigrationHandler) validateGamesPlayed(playlist string, current, previous int) []ValidationError {
	var errors []ValidationError
	fieldPrefix := strings.ToLower(playlist)

	if current < 0 || current > MaxGames {
		errors = append(errors, ValidationError{
			Field:   fieldPrefix + "_current_games",
			Message: fmt.Sprintf("%s current season games must be between 0 and %d", playlist, MaxGames),
			Code:    ValidationCodeOutOfRange,
		})
	}

	if previous < 0 || previous > MaxGames {
		errors = append(errors, ValidationError{
			Field:   fieldPrefix + "_previous_games",
			Message: fmt.Sprintf("%s previous season games must be between 0 and %d", playlist, MaxGames),
			Code:    ValidationCodeOutOfRange,
		})
	}

	return errors
}

// hasNoPlaylistData checks if the tracker has any meaningful playlist data
// Considers both MMR and games played. Negative values are treated as "no data" (equivalent to 0)
func (h *MigrationHandler) hasNoPlaylistData(tracker *usl.USLUserTracker) bool {
	hasOnesData := tracker.OnesCurrentSeasonPeak > 0 || tracker.OnesPreviousSeasonPeak > 0 || tracker.OnesAllTimePeak > 0 ||
		tracker.OnesCurrentSeasonGamesPlayed > 0 || tracker.OnesPreviousSeasonGamesPlayed > 0
	hasTwosData := tracker.TwosCurrentSeasonPeak > 0 || tracker.TwosPreviousSeasonPeak > 0 || tracker.TwosAllTimePeak > 0 ||
		tracker.TwosCurrentSeasonGamesPlayed > 0 || tracker.TwosPreviousSeasonGamesPlayed > 0
	hasThreesData := tracker.ThreesCurrentSeasonPeak > 0 || tracker.ThreesPreviousSeasonPeak > 0 || tracker.ThreesAllTimePeak > 0 ||
		tracker.ThreesCurrentSeasonGamesPlayed > 0 || tracker.ThreesPreviousSeasonGamesPlayed > 0

	return !hasOnesData && !hasTwosData && !hasThreesData
}

// isValidDiscordID validates Discord ID format (17-19 digit snowflake)
func isValidDiscordID(id string) bool {
	if len(id) < MinDiscordIDLength || len(id) > MaxDiscordIDLength {
		return false
	}
	_, err := strconv.ParseUint(id, 10, 64)
	return err == nil
}

// isValidTrackerURL validates tracker URL format and domain
func isValidTrackerURL(urlStr string) bool {
	// Check for common tracker sites
	validHosts := []string{
		"rocketleague.tracker.network",
		"ballchasing.com",
		"rltracker.pro",
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	// Must have a host
	if parsedURL.Host == "" {
		return false
	}

	// Check against whitelist
	for _, host := range validHosts {
		if strings.Contains(parsedURL.Host, host) {
			return true
		}
	}
	return false
}

// Legacy validation function - replaced by comprehensive validation system
// NOTE: validateTrackerRequired is deprecated, use validateTracker() instead

// Error display system for validation feedback

// renderFormWithErrors renders a form template with validation errors displayed
func (h *MigrationHandler) renderFormWithErrors(w http.ResponseWriter, templateName TemplateName, tracker *usl.USLUserTracker, errors []ValidationError) {
	data := struct {
		Title       string
		CurrentPage string
		Tracker     *usl.USLUserTracker
		Errors      map[string]string // For template .Errors.field_name access
	}{
		Title:       "Tracker Form",
		CurrentPage: "trackers",
		Tracker:     tracker,
		Errors:      h.buildErrorMap(errors),
	}

	h.renderTemplate(w, templateName, data)
}

func (h *MigrationHandler) buildErrorMap(errors []ValidationError) map[string]string {
	errorMap := make(map[string]string)
	for _, err := range errors {
		errorMap[err.Field] = err.Message
	}
	return errorMap
}

func (h *MigrationHandler) calculateEffectiveMMR(tracker *usl.USLUserTracker) {
	weightedSum := (tracker.OnesCurrentSeasonPeak * tracker.OnesCurrentSeasonGamesPlayed) +
		(tracker.TwosCurrentSeasonPeak * tracker.TwosCurrentSeasonGamesPlayed) +
		(tracker.ThreesCurrentSeasonPeak * tracker.ThreesCurrentSeasonGamesPlayed)

	totalGames := tracker.OnesCurrentSeasonGamesPlayed + tracker.TwosCurrentSeasonGamesPlayed + tracker.ThreesCurrentSeasonGamesPlayed

	if totalGames > 0 {
		tracker.MMR = weightedSum / totalGames
	} else {
		tracker.MMR = 0
	}
}

// TrueSkill integration functions (currently unused in validation-focused CRUD)
// These functions are kept for future TrueSkill integration when async ranking updates are implemented

// logTrueSkillUpdateFailure logs structured failure information
// Reserved for future TrueSkill integration - currently unused but kept for planned implementation
// func (h *MigrationHandler) logTrueSkillUpdateFailure(tracker *usl.USLUserTracker, errorMsg string) {
// 	log.Printf("[USL-TRUESKILL] WARNING: Update failed - Discord=%s, TrackerID=%d, Error=%s",
// 		tracker.DiscordID, tracker.ID, errorMsg)
// 	log.Printf("[USL-TRUESKILL] Tracker creation succeeded, TrueSkill update failed - manual intervention may be needed")
// }

// logTrueSkillUpdateSuccess logs structured success information
// Reserved for future TrueSkill integration - currently unused but kept for planned implementation
// func (h *MigrationHandler) logTrueSkillUpdateSuccess(tracker *usl.USLUserTracker, mu, sigma float64) {
// 	log.Printf("[USL-TRUESKILL] SUCCESS: Calculated - Discord=%s, TrackerID=%d, μ=%.1f, σ=%.2f",
// 		tracker.DiscordID, tracker.ID, mu, sigma)
// }

// logUSLSyncFailure logs structured USL sync failure information
// Reserved for future TrueSkill integration - currently unused but kept for planned implementation
// func (h *MigrationHandler) logUSLSyncFailure(discordID string, err error) {
// 	log.Printf("[USL-TRUESKILL] WARNING: USL table sync failed - Discord=%s, Error=%v", discordID, err)
// 	log.Printf("[USL-TRUESKILL] Core tables updated successfully, USL tables inconsistent - manual sync required")
// }

// logUSLSyncSuccess logs structured USL sync success information
// Reserved for future TrueSkill integration - currently unused but kept for planned implementation
// func (h *MigrationHandler) logUSLSyncSuccess(discordID string) {
// 	log.Printf("[USL-TRUESKILL] SUCCESS: Full integration completed - Discord=%s (Core ✓, USL ✓)", discordID)
// }

const (
	USL_DISCORD_GUILD_ID = "1390537743385231451" // USL Discord Guild ID
)

// MigrationHandler provides simplified handlers for USL-only operations
// This is a temporary migration solution - no multi-guild complexity
// AUTH NOTE: This handler no longer manages auth - that's handled by unified Discord OAuth in main.go
type MigrationHandler struct {
	uslRepo          *usl.USLRepository
	templates        *template.Template
	trueskillService *services.UserTrueSkillService
	config           *config.Config
}

func NewMigrationHandler(
	uslRepo *usl.USLRepository,
	templates *template.Template,
	trueskillService *services.UserTrueSkillService,
	config *config.Config,
) *MigrationHandler {
	return &MigrationHandler{
		uslRepo:          uslRepo,
		templates:        templates,
		trueskillService: trueskillService,
		config:           config,
	}
}

func (h *MigrationHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	users, err := h.uslRepo.GetAllUsers()
	if err != nil {
		h.handleDatabaseError(w, "load users", err)
		return
	}

	data := struct {
		Title        string
		CurrentPage  string
		Users        []*usl.USLUser
		SearchConfig struct {
			SearchPlaceholder string
			SearchURL         string
			SearchTarget      string
			ClearURL          string
			ShowFilters       bool
			Query             string
			StatusFilter      string
		}
	}{
		Title:       "Users",
		CurrentPage: "users",
		Users:       users,
		SearchConfig: struct {
			SearchPlaceholder string
			SearchURL         string
			SearchTarget      string
			ClearURL          string
			ShowFilters       bool
			Query             string
			StatusFilter      string
		}{
			SearchPlaceholder: "Search by name or Discord ID...",
			SearchURL:         "/usl/users/search",
			SearchTarget:      "#users-tbody",
			ClearURL:          "/usl/users/search",
			ShowFilters:       true,
			Query:             "",
			StatusFilter:      "",
		},
	}

	h.renderTemplate(w, TemplateUSLUsers, data)
}

func (h *MigrationHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	query := r.URL.Query().Get("q")

	var users []*usl.USLUser
	var err error

	if query == "" {

		users, err = h.uslRepo.GetAllUsers()
	} else {
		// Search for users matching the query
		users, err = h.uslRepo.SearchUsers(query)
	}
	if err != nil {
		h.handleDatabaseError(w, "search users", err)
		return
	}

	data := struct {
		Title        string
		Users        []*usl.USLUser
		SearchConfig struct {
			SearchPlaceholder string
			SearchURL         string
			SearchTarget      string
			ClearURL          string
			ShowFilters       bool
			Query             string
			StatusFilter      string
		}
	}{
		Title: "Users",
		Users: users,
		SearchConfig: struct {
			SearchPlaceholder string
			SearchURL         string
			SearchTarget      string
			ClearURL          string
			ShowFilters       bool
			Query             string
			StatusFilter      string
		}{
			SearchPlaceholder: "Search by name or Discord ID...",
			SearchURL:         "/usl/users/search",
			SearchTarget:      "#users-tbody",
			ClearURL:          "/usl/users/search",
			ShowFilters:       true,
			Query:             query,
			StatusFilter:      "",
		},
	}

	h.renderTemplate(w, TemplateUSLUsersTable, data)
}

func (h *MigrationHandler) ListTrackers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	trackers, err := h.uslRepo.GetAllTrackers()
	if err != nil {
		h.handleDatabaseError(w, "load trackers", err)
		return
	}

	data := struct {
		Title        string
		CurrentPage  string
		Trackers     []*usl.USLUserTracker
		SearchConfig struct {
			SearchPlaceholder string
			SearchURL         string
			SearchTarget      string
			ClearURL          string
			ShowFilters       bool
			Query             string
			StatusFilter      string
		}
	}{
		Title:       "Trackers",
		CurrentPage: "trackers",
		Trackers:    trackers,
		SearchConfig: struct {
			SearchPlaceholder string
			SearchURL         string
			SearchTarget      string
			ClearURL          string
			ShowFilters       bool
			Query             string
			StatusFilter      string
		}{
			SearchPlaceholder: "Search by URL or Discord ID...",
			SearchURL:         "/usl/trackers/search",
			SearchTarget:      "#trackers-table",
			ClearURL:          "/usl/trackers/search",
			ShowFilters:       false, // Trackers don't have status filters like users
			Query:             "",
			StatusFilter:      "",
		},
	}

	// Populate user data for each tracker
	h.populateTrackerUsers(trackers)

	h.renderTemplate(w, TemplateUSLTrackers, data)
}

// populateTrackerUsers fetches and populates the User field for each tracker
func (h *MigrationHandler) populateTrackerUsers(trackers []*usl.USLUserTracker) {
	for _, tracker := range trackers {
		if user, err := h.uslRepo.GetUserByDiscordID(tracker.DiscordID); err == nil {
			tracker.User = user
		}
		// If error, User remains nil and template will show Discord ID instead
	}
}

func (h *MigrationHandler) SearchTrackers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {

		trackers, err := h.uslRepo.GetAllTrackers()
		if err != nil {
			h.handleDatabaseError(w, "load trackers", err)
			return
		}

		data := struct {
			Title        string
			Trackers     []*usl.USLUserTracker
			SearchConfig struct {
				SearchPlaceholder string
				SearchURL         string
				SearchTarget      string
				ClearURL          string
				ShowFilters       bool
				Query             string
				StatusFilter      string
			}
		}{
			Title:    "Trackers",
			Trackers: trackers,
			SearchConfig: struct {
				SearchPlaceholder string
				SearchURL         string
				SearchTarget      string
				ClearURL          string
				ShowFilters       bool
				Query             string
				StatusFilter      string
			}{
				SearchPlaceholder: "Search by URL or Discord ID...",
				SearchURL:         "/usl/trackers/search",
				SearchTarget:      "#trackers-table",
				ClearURL:          "/usl/trackers/search",
				ShowFilters:       false,
				Query:             "",
				StatusFilter:      "",
			},
		}

		// Populate user data for each tracker
		h.populateTrackerUsers(trackers)

		h.renderTemplate(w, "trackers-table-fragment", data)
		return
	}

	// Search for trackers matching the query
	trackers, err := h.uslRepo.SearchTrackers(query)
	if err != nil {
		h.handleDatabaseError(w, "search trackers", err)
		return
	}

	data := struct {
		Title        string
		Trackers     []*usl.USLUserTracker
		SearchConfig struct {
			SearchPlaceholder string
			SearchURL         string
			SearchTarget      string
			ClearURL          string
			ShowFilters       bool
			Query             string
			StatusFilter      string
		}
	}{
		Title:    "Trackers",
		Trackers: trackers,
		SearchConfig: struct {
			SearchPlaceholder string
			SearchURL         string
			SearchTarget      string
			ClearURL          string
			ShowFilters       bool
			Query             string
			StatusFilter      string
		}{
			SearchPlaceholder: "Search by URL or Discord ID...",
			SearchURL:         "/usl/trackers/search",
			SearchTarget:      "#trackers-table",
			ClearURL:          "/usl/trackers/search",
			ShowFilters:       false,
			Query:             query,
			StatusFilter:      "",
		},
	}

	// Populate user data for each tracker
	h.populateTrackerUsers(trackers)

	h.renderTemplate(w, "trackers-table-fragment", data)
}

func (h *MigrationHandler) UserDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		h.handleInvalidID(w, "User ID")
		return
	}

	userID, err := h.parseUserID(userIDStr)
	if err != nil {
		h.handleParseError(w, "user ID")
		return
	}

	user, err := h.uslRepo.GetUserByID(userID)
	if err != nil {
		h.handleDatabaseError(w, "load user", err)
		return
	}

	// Get user's trackers
	userTrackers, err := h.uslRepo.GetTrackersByDiscordID(user.DiscordID)
	if err != nil {
		h.handleDatabaseError(w, "load user trackers", err)
		return
	}

	data := struct {
		Title        string
		CurrentPage  string
		User         *usl.USLUser
		UserTrackers []*usl.USLUserTracker
	}{
		Title:        user.Name,
		CurrentPage:  "users",
		User:         user,
		UserTrackers: userTrackers,
	}

	h.renderTemplate(w, TemplateUSLUserDetail, data)
}

// updateUSLUserTrueSkillFromTrackers updates TrueSkill for a USL user from their tracker data
// This function manages USL data access and delegates calculation to the TrueSkill service
func (h *MigrationHandler) updateUSLUserTrueSkillFromTrackers(discordID string) *services.TrueSkillUpdateResult {
	// Get USL tracker data
	userTrackers, err := h.uslRepo.GetTrackersByDiscordID(discordID)
	if err != nil {
		validationLogger.Error("Failed to get USL trackers for TrueSkill calculation",
			"discord_id", discordID,
			"error", err)
		return &services.TrueSkillUpdateResult{
			Success:     false,
			HadTrackers: false,
			Error:       fmt.Sprintf("failed to get trackers: %v", err),
		}
	}

	// If no trackers, assign default values
	if len(userTrackers) == 0 {
		defaultMu := 1500.0
		defaultSigma := 8.333

		err = h.uslRepo.UpdateUserTrueSkill(discordID, defaultMu, defaultSigma)
		if err != nil {
			validationLogger.Error("Failed to update USL user with default TrueSkill values",
				"discord_id", discordID,
				"error", err)
			return &services.TrueSkillUpdateResult{
				Success:     false,
				HadTrackers: false,
				Error:       fmt.Sprintf("failed to update user with defaults: %v", err),
			}
		}

		validationLogger.Info("Assigned default TrueSkill values to USL user",
			"discord_id", discordID,
			"mu", defaultMu,
			"sigma", defaultSigma)

		return &services.TrueSkillUpdateResult{
			Success:     true,
			HadTrackers: false,
			TrueSkillResult: &services.TrueSkillCalculation{
				Mu:    defaultMu,
				Sigma: defaultSigma,
			},
		}
	}

	// Transform USL tracker data to TrueSkill service format
	trackerData := h.transformUSLTrackerToTrackerData(userTrackers[0])

	// Use TrueSkill service to calculate values (no database access in service)
	result := h.trueskillService.CalculateTrueSkillFromTrackerData(trackerData)
	if !result.Success {
		validationLogger.Error("TrueSkill calculation failed for USL user",
			"discord_id", discordID,
			"error", result.Error)
		return result
	}

	// Update USL user with calculated values
	err = h.uslRepo.UpdateUserTrueSkill(discordID, result.TrueSkillResult.Mu, result.TrueSkillResult.Sigma)
	if err != nil {
		validationLogger.Error("Failed to update USL user with calculated TrueSkill values",
			"discord_id", discordID,
			"mu", result.TrueSkillResult.Mu,
			"sigma", result.TrueSkillResult.Sigma,
			"error", err)
		return &services.TrueSkillUpdateResult{
			Success:     false,
			HadTrackers: true,
			Error:       fmt.Sprintf("failed to update user: %v", err),
		}
	}

	validationLogger.Info("Successfully calculated and updated TrueSkill for USL user",
		"discord_id", discordID,
		"mu", result.TrueSkillResult.Mu,
		"sigma", result.TrueSkillResult.Sigma)

	return result
}

// transformUSLTrackerToTrackerData converts USL tracker data to TrueSkill service format
func (h *MigrationHandler) transformUSLTrackerToTrackerData(uslTracker *usl.USLUserTracker) *services.TrackerData {
	return &services.TrackerData{
		DiscordID:           uslTracker.DiscordID,
		URL:                 uslTracker.URL,
		OnesCurrentPeak:     uslTracker.OnesCurrentSeasonPeak,
		OnesPreviousPeak:    uslTracker.OnesPreviousSeasonPeak,
		OnesAllTimePeak:     uslTracker.OnesAllTimePeak,
		OnesCurrentGames:    uslTracker.OnesCurrentSeasonGamesPlayed,
		OnesPreviousGames:   uslTracker.OnesPreviousSeasonGamesPlayed,
		TwosCurrentPeak:     uslTracker.TwosCurrentSeasonPeak,
		TwosPreviousPeak:    uslTracker.TwosPreviousSeasonPeak,
		TwosAllTimePeak:     uslTracker.TwosAllTimePeak,
		TwosCurrentGames:    uslTracker.TwosCurrentSeasonGamesPlayed,
		TwosPreviousGames:   uslTracker.TwosPreviousSeasonGamesPlayed,
		ThreesCurrentPeak:   uslTracker.ThreesCurrentSeasonPeak,
		ThreesPreviousPeak:  uslTracker.ThreesPreviousSeasonPeak,
		ThreesAllTimePeak:   uslTracker.ThreesAllTimePeak,
		ThreesCurrentGames:  uslTracker.ThreesCurrentSeasonGamesPlayed,
		ThreesPreviousGames: uslTracker.ThreesPreviousSeasonGamesPlayed,
	}
}

// UpdateUserTrueSkill recalculates TrueSkill for a specific user
func (h *MigrationHandler) UpdateUserTrueSkill(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.handleMethodNotAllowed(w, r)
		return
	}

	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		h.handleInvalidID(w, "User ID")
		return
	}

	userID, err := h.parseUserID(userIDStr)
	if err != nil {
		h.handleParseError(w, "user ID")
		return
	}

	user, err := h.uslRepo.GetUserByID(userID)
	if err != nil {
		h.handleDatabaseError(w, "load user", err)
		return
	}

	// Update TrueSkill using USL-specific function
	result := h.updateUSLUserTrueSkillFromTrackers(user.DiscordID)

	// Log the manual TrueSkill update for audit purposes
	validationLogger.Info("Manual TrueSkill update triggered",
		"admin_action", "update_trueskill",
		"target_user_id", userID,
		"target_discord_id", user.DiscordID,
		"update_success", result.Success,
		"had_trackers", result.HadTrackers,
		"remote_addr", r.RemoteAddr,
		"user_agent", r.Header.Get("User-Agent"),
	)

	// Return HTMX-friendly response with update results
	h.renderTrueSkillUpdateResult(w, r, result, user)
}

func (h *MigrationHandler) TrackerDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	trackerIDStr := r.URL.Query().Get("id")
	if trackerIDStr == "" {
		h.handleInvalidID(w, "Tracker ID")
		return
	}

	trackerID, err := strconv.ParseInt(trackerIDStr, 10, 64)
	if err != nil {
		h.handleParseError(w, "tracker ID")
		return
	}

	tracker, err := h.uslRepo.GetTrackerByID(trackerID)
	if err != nil {
		h.handleDatabaseError(w, "load tracker", err)
		return
	}

	// Get user associated with this tracker for display name
	user, err := h.uslRepo.GetUserByDiscordID(tracker.DiscordID)
	if err != nil {
		h.handleDatabaseError(w, "load associated user", err)
		return
	}

	data := struct {
		Title       string
		CurrentPage string
		Tracker     *usl.USLUserTracker
		User        *usl.USLUser
	}{
		Title:       "Tracker Details",
		CurrentPage: "trackers",
		Tracker:     tracker,
		User:        user,
	}

	h.renderTemplate(w, TemplateUSLTrackerDetail, data)
}

func (h *MigrationHandler) NewTrackerForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	data := struct {
		Title       string
		CurrentPage string
		Tracker     *usl.USLUserTracker
		Errors      map[string]string
	}{
		Title:       "New Tracker",
		CurrentPage: "trackers",
		Tracker:     &usl.USLUserTracker{},   // Empty tracker for new forms
		Errors:      make(map[string]string), // Empty errors for initial form load
	}

	h.renderTemplate(w, TemplateUSLTrackerNew, data)
}

func (h *MigrationHandler) CreateTracker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.handleMethodNotAllowed(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.handleInvalidFormData(w, err)
		return
	}

	tracker := h.buildTrackerFromForm(r)

	// Comprehensive validation with metrics and security monitoring
	validation := h.validateTrackerWithMetrics(r, tracker)
	if !validation.IsValid {
		h.renderFormWithErrors(w, TemplateUSLTrackerNew, tracker, validation.Errors)
		return
	}

	// Calculate MMR (using extracted function)
	h.calculateEffectiveMMR(tracker)

	createdTracker, err := h.uslRepo.CreateTracker(tracker)
	if err != nil {
		h.handleDatabaseError(w, "create tracker", err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/usl/trackers/detail?id=%d", createdTracker.ID), http.StatusSeeOther)
}

func (h *MigrationHandler) EditTrackerForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	trackerIDStr := r.URL.Query().Get("id")
	if trackerIDStr == "" {
		h.handleInvalidID(w, "Tracker ID")
		return
	}

	trackerID, err := strconv.ParseInt(trackerIDStr, 10, 64)
	if err != nil {
		h.handleParseError(w, "tracker ID")
		return
	}

	tracker, err := h.uslRepo.GetTrackerByID(trackerID)
	if err != nil {
		h.handleDatabaseError(w, "load tracker", err)
		return
	}

	// Fetch user information for the tracker
	user, err := h.uslRepo.GetUserByDiscordID(tracker.DiscordID)
	if err != nil {
		// If user not found, we'll still show the form but without name
		user = nil
	}

	data := struct {
		Title       string
		CurrentPage string
		Tracker     *usl.USLUserTracker
		User        *usl.USLUser
		Errors      map[string]string
	}{
		Title:       "Edit Tracker",
		CurrentPage: "trackers",
		Tracker:     tracker,
		User:        user,
		Errors:      make(map[string]string), // Empty errors for initial form load
	}

	h.renderTemplate(w, TemplateUSLTrackerEdit, data)
}

func (h *MigrationHandler) UpdateTracker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.handleMethodNotAllowed(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.handleInvalidFormData(w, err)
		return
	}

	trackerIDStr := r.FormValue("id")
	if trackerIDStr == "" {
		h.handleInvalidID(w, "Tracker ID")
		return
	}

	trackerID, err := strconv.ParseInt(trackerIDStr, 10, 64)
	if err != nil {
		h.handleParseError(w, "tracker ID")
		return
	}

	tracker := h.buildTrackerFromForm(r)
	tracker.ID = trackerID // Set ID for update operation

	// Comprehensive validation with metrics and security monitoring
	validation := h.validateTrackerWithMetrics(r, tracker)
	if !validation.IsValid {
		h.renderFormWithErrors(w, TemplateUSLTrackerEdit, tracker, validation.Errors)
		return
	}

	// Calculate MMR (using extracted function)
	h.calculateEffectiveMMR(tracker)

	err = h.uslRepo.UpdateTracker(tracker)
	if err != nil {
		h.handleDatabaseError(w, "update tracker", err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/usl/trackers/detail?id=%d", trackerID), http.StatusSeeOther)
}

func (h *MigrationHandler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	stats, err := h.uslRepo.GetStats()
	if err != nil {
		h.handleDatabaseError(w, "load dashboard", err)
		return
	}

	// Create a mock guild for USL compatibility
	_ = &models.Guild{
		ID:             1,
		DiscordGuildID: USLDiscordGuildID,
		Name:           "USL",
		Slug:           "usl",
		Active:         true,
		Config:         models.GetDefaultGuildConfig(),
		Theme:          models.GetDefaultTheme(),
	}

	data := struct {
		Title       string
		CurrentPage string
		Stats       struct {
			TotalUsers    int `json:"total_users"`
			ActiveUsers   int `json:"active_users"`
			TotalTrackers int `json:"total_trackers"`
			ValidTrackers int `json:"valid_trackers"`
		}
	}{
		Title:       "Dashboard",
		CurrentPage: "admin",
		Stats: struct {
			TotalUsers    int `json:"total_users"`
			ActiveUsers   int `json:"active_users"`
			TotalTrackers int `json:"total_trackers"`
			ValidTrackers int `json:"valid_trackers"`
		}{
			TotalUsers:    stats["total_users"].(int),
			ActiveUsers:   stats["active_users"].(int),
			TotalTrackers: stats["total_trackers"].(int),
			ValidTrackers: stats["valid_trackers"].(int),
		},
	}

	h.renderTemplate(w, TemplateUSLAdminDashboard, data)
}

func (h *MigrationHandler) ListUsersAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	users, err := h.uslRepo.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to load users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("[USL-HANDLER] JSON encoding error for users API: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *MigrationHandler) ListTrackersAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	trackers, err := h.uslRepo.GetAllTrackers()
	if err != nil {
		http.Error(w, "Failed to load trackers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(trackers); err != nil {
		log.Printf("[USL-HANDLER] JSON encoding error for trackers API: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

func (h *MigrationHandler) GetLeaderboardAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	users, err := h.uslRepo.GetLeaderboard()
	if err != nil {
		http.Error(w, "Failed to load leaderboard", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("[USL-HANDLER] JSON encoding error for leaderboard API: %v", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// performTrueSkillUpdate handles TrueSkill calculation and synchronization with comprehensive error handling
// Reserved for future TrueSkill integration - currently unused but kept for planned implementation
// func (h *MigrationHandler) performTrueSkillUpdate(tracker *usl.USLUserTracker) {
// 	if tracker == nil {
// 		log.Printf("[USL-TRUESKILL] ERROR: Cannot perform update - tracker is nil")
// 		return
// 	}

// 	trackerData := h.mapUSLTrackerToTrackerData(tracker)
// 	log.Printf("[USL-TRUESKILL] Starting update - Discord=%s, TrackerID=%d", tracker.DiscordID, tracker.ID)

// 	result := h.trueskillService.UpdateUserTrueSkillFromTrackerData(trackerData)

// 	if !result.Success {
// 		h.logTrueSkillUpdateFailure(tracker, result.Error)
// 		return
// 	}

// 	h.logTrueSkillUpdateSuccess(tracker, result.TrueSkillResult.Mu, result.TrueSkillResult.Sigma)

// 	err := h.uslRepo.UpdateUserTrueSkill(
// 		tracker.DiscordID,
// 		result.TrueSkillResult.Mu,
// 		result.TrueSkillResult.Sigma,
// 	)

// 	if err != nil {
// 		h.logUSLSyncFailure(tracker.DiscordID, err)
// 	} else {
// 		h.logUSLSyncSuccess(tracker.DiscordID)
// 	}
// }

// mapUSLTrackerToTrackerData converts USL tracker format to TrueSkill service input format
func (h *MigrationHandler) mapUSLTrackerToTrackerData(uslTracker *usl.USLUserTracker) *services.TrackerData {
	var lastUpdated time.Time
	if uslTracker.LastUpdated != nil && *uslTracker.LastUpdated != "" {
		if parsed, err := time.Parse(time.RFC3339, *uslTracker.LastUpdated); err == nil {
			lastUpdated = parsed
		} else {
			lastUpdated = time.Now()
		}
	} else {
		lastUpdated = time.Now()
	}

	return &services.TrackerData{
		DiscordID:           uslTracker.DiscordID,
		URL:                 uslTracker.URL,
		OnesCurrentPeak:     uslTracker.OnesCurrentSeasonPeak,
		OnesCurrentGames:    uslTracker.OnesCurrentSeasonGamesPlayed,
		OnesPreviousPeak:    uslTracker.OnesPreviousSeasonPeak,
		OnesPreviousGames:   uslTracker.OnesPreviousSeasonGamesPlayed,
		TwosCurrentPeak:     uslTracker.TwosCurrentSeasonPeak,
		TwosCurrentGames:    uslTracker.TwosCurrentSeasonGamesPlayed,
		TwosPreviousPeak:    uslTracker.TwosPreviousSeasonPeak,
		TwosPreviousGames:   uslTracker.TwosPreviousSeasonGamesPlayed,
		ThreesCurrentPeak:   uslTracker.ThreesCurrentSeasonPeak,
		ThreesCurrentGames:  uslTracker.ThreesCurrentSeasonGamesPlayed,
		ThreesPreviousPeak:  uslTracker.ThreesPreviousSeasonPeak,
		ThreesPreviousGames: uslTracker.ThreesPreviousSeasonGamesPlayed,
		LastUpdated:         lastUpdated,
	}
}

// ValidationMetricsAPI returns current validation metrics for monitoring
func (h *MigrationHandler) ValidationMetricsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	metrics := getValidationMetrics()

	// Add calculated statistics
	response := struct {
		ValidationMetrics
		SuccessRate   float64 `json:"success_rate"`
		FailureRate   float64 `json:"failure_rate"`
		SecurityRate  float64 `json:"security_incident_rate"`
		TopErrorTypes []struct {
			Type  string `json:"type"`
			Count int64  `json:"count"`
		} `json:"top_error_types"`
		TopErrorFields []struct {
			Field string `json:"field"`
			Count int64  `json:"count"`
		} `json:"top_error_fields"`
	}{
		ValidationMetrics: metrics,
	}

	// Calculate rates
	if metrics.TotalValidations > 0 {
		response.SuccessRate = float64(metrics.SuccessfulValidations) / float64(metrics.TotalValidations) * 100
		response.FailureRate = float64(metrics.FailedValidations) / float64(metrics.TotalValidations) * 100
		response.SecurityRate = float64(metrics.SecurityIncidents) / float64(metrics.TotalValidations) * 100
	}

	// Get top error types (up to 5)
	for errorType, count := range metrics.ErrorsByType {
		response.TopErrorTypes = append(response.TopErrorTypes, struct {
			Type  string `json:"type"`
			Count int64  `json:"count"`
		}{Type: errorType, Count: count})
	}

	// Get top error fields (up to 5)
	for errorField, count := range metrics.ErrorsByField {
		response.TopErrorFields = append(response.TopErrorFields, struct {
			Field string `json:"field"`
			Count int64  `json:"count"`
		}{Field: errorField, Count: count})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		validationLogger.Error("Failed to encode validation metrics response", "error", err)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// renderTrueSkillUpdateResult renders the TrueSkill update result as an HTMX-friendly response
func (h *MigrationHandler) renderTrueSkillUpdateResult(w http.ResponseWriter, r *http.Request, result *services.TrueSkillUpdateResult, user *usl.USLUser) {
	// Create response data for the template
	data := struct {
		Success         bool   `json:"success"`
		HadTrackers     bool   `json:"hadTrackers"`
		Error           string `json:"error,omitempty"`
		TrueSkillResult *struct {
			Mu    float64 `json:"mu"`
			Sigma float64 `json:"sigma"`
		} `json:"trueSkillResult,omitempty"`
		UserName  string `json:"userName"`
		DiscordID string `json:"discordId"`
	}{
		Success:     result.Success,
		HadTrackers: result.HadTrackers,
		Error:       result.Error,
		UserName:    user.Name,
		DiscordID:   user.DiscordID,
	}

	// Add TrueSkill result if successful
	if result.Success && result.TrueSkillResult != nil {
		data.TrueSkillResult = &struct {
			Mu    float64 `json:"mu"`
			Sigma float64 `json:"sigma"`
		}{
			Mu:    result.TrueSkillResult.Mu,
			Sigma: result.TrueSkillResult.Sigma,
		}
	}

	// For HTMX requests, return just the result fragment
	if r.Header.Get("HX-Request") != "" {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := h.templates.ExecuteTemplate(w, "trueskill-update-result", data)
		if err != nil {
			validationLogger.Error("Failed to render TrueSkill update result template", "error", err)
			http.Error(w, "Template rendering error", http.StatusInternalServerError)
		}
		return
	}

	// For non-HTMX requests, redirect back to user detail page
	http.Redirect(w, r, fmt.Sprintf("/usl/users/detail?id=%d", user.ID), http.StatusSeeOther)
}

// NOTE: Authentication methods removed - now handled by unified Discord OAuth system in main.go

func (h *MigrationHandler) renderTemplate(w http.ResponseWriter, templateName TemplateName, data any) {
	// Use buffer-based rendering to prevent partial output on errors (2025 best practice)
	var buf bytes.Buffer
	err := h.templates.ExecuteTemplate(&buf, string(templateName), data)
	if err != nil {
		log.Printf("[USL-HANDLER] Template rendering error: template=%s, error=%v", templateName, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Only write to ResponseWriter after successful rendering
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := buf.WriteTo(w); err != nil {
		log.Printf("[USL-HANDLER] Failed to write template output: %v", err)
	}
}

// NewUserForm displays the form for creating a new user
func (h *MigrationHandler) NewUserForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	data := struct {
		Title       string
		User        *usl.USLUser
		CurrentPage string
	}{
		Title:       "Add New User",
		User:        &usl.USLUser{}, // Empty user for form
		CurrentPage: "users",
	}

	h.renderTemplate(w, "user-new-page", data)
}

// CreateUser handles the creation of a new user
func (h *MigrationHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.handleMethodNotAllowed(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.handleInvalidFormData(w, err)
		return
	}

	user := &usl.USLUser{
		Name:      strings.TrimSpace(r.FormValue("name")),
		DiscordID: strings.TrimSpace(r.FormValue("discord_id")),
		Active:    r.FormValue("active") == "on",
		Banned:    r.FormValue("banned") == "on",
	}

	// Validate required fields
	if user.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	if user.DiscordID == "" {
		http.Error(w, "Discord ID is required", http.StatusBadRequest)
		return
	}

	createdUser, err := h.uslRepo.CreateUser(user.Name, user.DiscordID, user.Active, user.Banned)
	if err != nil {
		h.handleDatabaseError(w, "create user", err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/usl/users/detail?id=%d", createdUser.ID), http.StatusSeeOther)
}

// EditUserForm displays the form for editing an existing user
func (h *MigrationHandler) EditUserForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		h.handleInvalidID(w, "user")
		return
	}

	userID, err := h.parseUserID(userIDStr)
	if err != nil {
		h.handleParseError(w, "user ID")
		return
	}

	user, err := h.uslRepo.GetUserByID(userID)
	if err != nil {
		h.handleDatabaseError(w, "fetch user", err)
		return
	}

	data := struct {
		Title       string
		User        *usl.USLUser
		CurrentPage string
	}{
		Title:       "Edit User",
		User:        user,
		CurrentPage: "users",
	}

	h.renderTemplate(w, "user-edit-page", data)
}

// UpdateUser handles the update of an existing user
func (h *MigrationHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.handleMethodNotAllowed(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.handleInvalidFormData(w, err)
		return
	}

	userIDStr := r.FormValue("id")
	userID, err := h.parseUserID(userIDStr)
	if err != nil {
		h.handleParseError(w, "user ID")
		return
	}

	user := &usl.USLUser{
		ID:        userID,
		Name:      strings.TrimSpace(r.FormValue("name")),
		DiscordID: strings.TrimSpace(r.FormValue("discord_id")),
		Active:    r.FormValue("active") == "on",
		Banned:    r.FormValue("banned") == "on",
	}

	// Validate required fields
	if user.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	if user.DiscordID == "" {
		http.Error(w, "Discord ID is required", http.StatusBadRequest)
		return
	}

	_, err = h.uslRepo.UpdateUser(userID, user.Name, user.Active, user.Banned)
	if err != nil {
		h.handleDatabaseError(w, "update user", err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/usl/users/detail?id=%d", userID), http.StatusSeeOther)
}

// DeleteUser handles the deletion of a user
func (h *MigrationHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.handleMethodNotAllowed(w, r)
		return
	}

	userIDStr := r.URL.Path[len("/usl/users/"):]
	userID, err := h.parseUserID(userIDStr)
	if err != nil {
		h.handleParseError(w, "user ID")
		return
	}

	err = h.uslRepo.DeleteUser(userID)
	if err != nil {
		h.handleDatabaseError(w, "delete user", err)
		return
	}

	// For HTMX requests, return success
	if r.Header.Get("HX-Request") != "" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// For regular requests, redirect to users list
	http.Redirect(w, r, "/usl/users", http.StatusSeeOther)
}
