package handlers

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
	"usl-server/internal/logger"
	"usl-server/internal/services"
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

// recordValidationMetrics updates validation metrics in a thread-safe manner
func recordValidationMetrics(event *ValidationEvent) {
	metricsMutex.Lock()
	defer metricsMutex.Unlock()

	validationMetrics.TotalValidations++

	if event.Type == "success" {
		validationMetrics.SuccessfulValidations++
	} else if event.Type == "failure" || event.Type == "security_incident" {
		validationMetrics.FailedValidations++

		for _, err := range event.Errors {
			validationMetrics.ErrorsByType[err.Code]++
			validationMetrics.ErrorsByField[err.Field]++
		}

		if event.Type == "security_incident" {
			validationMetrics.SecurityIncidents++
		}
	}
}

// getValidationMetrics returns a copy of the current validation metrics
func getValidationMetrics() *ValidationMetrics {
	metricsMutex.RLock()
	defer metricsMutex.RUnlock()

	// Return a copy to avoid race conditions
	copy := &ValidationMetrics{
		TotalValidations:      validationMetrics.TotalValidations,
		SuccessfulValidations: validationMetrics.SuccessfulValidations,
		FailedValidations:     validationMetrics.FailedValidations,
		SecurityIncidents:     validationMetrics.SecurityIncidents,
		LastReset:             validationMetrics.LastReset,
		ErrorsByType:          make(map[string]int64),
		ErrorsByField:         make(map[string]int64),
	}

	for k, v := range validationMetrics.ErrorsByType {
		copy.ErrorsByType[k] = v
	}
	for k, v := range validationMetrics.ErrorsByField {
		copy.ErrorsByField[k] = v
	}

	return copy
}

// detectSecurityIncident analyzes a request for potential security threats
func detectSecurityIncident(r *http.Request, validation *ValidationResult) *string {
	// Check for SQL injection patterns
	for key, values := range r.Form {
		for _, value := range values {
			if strings.Contains(strings.ToLower(value), "drop table") ||
				strings.Contains(strings.ToLower(value), "'; drop") ||
				strings.Contains(strings.ToLower(value), "union select") ||
				strings.Contains(strings.ToLower(value), "' or '1'='1") {
				reason := "SQL injection attempt detected"
				return &reason
			}

			// Check for XSS patterns
			if strings.Contains(value, "<script>") ||
				strings.Contains(value, "javascript:") ||
				strings.Contains(value, "onerror=") ||
				strings.Contains(value, "onload=") {
				reason := "XSS attempt detected"
				return &reason
			}

			// Check for potential buffer overflow (extremely long inputs)
			if len(value) > 1000 {
				reason := "Buffer overflow attempt detected"
				return &reason
			}

			// Check for path traversal attempts
			if strings.Contains(value, "../") ||
				strings.Contains(value, "..\\") {
				reason := "Path traversal attempt detected"
				return &reason
			}
		}

		// Log the field name for monitoring
		_ = key
	}

	return nil
}

// Helper function for tests - wrapper around BaseHandler.isValidTrackerURL
func isValidTrackerURL(url string) bool {
	h := &BaseHandler{}
	return h.isValidTrackerURL(url)
}

// buildTrackerFromForm constructs a USLUserTracker from form data
func (h *BaseHandler) buildTrackerFromForm(r *http.Request) *USLUserTracker {
	tracker := &USLUserTracker{
		DiscordID: strings.TrimSpace(h.getFormValue(r, FormFieldDiscordID)),
		URL:       strings.TrimSpace(h.getFormValue(r, FormFieldURL)),
		Valid:     h.getFormBoolValue(r, FormFieldValid),

		OnesCurrentSeasonPeak:         h.getFormIntValue(r, FormFieldOnesCurrentPeak),
		OnesPreviousSeasonPeak:        h.getFormIntValue(r, FormFieldOnesPreviousPeak),
		OnesAllTimePeak:               h.getFormIntValue(r, FormFieldOnesAllTimePeak),
		OnesCurrentSeasonGamesPlayed:  h.getFormIntValue(r, FormFieldOnesCurrentGames),
		OnesPreviousSeasonGamesPlayed: h.getFormIntValue(r, FormFieldOnesPreviousGames),

		TwosCurrentSeasonPeak:         h.getFormIntValue(r, FormFieldTwosCurrentPeak),
		TwosPreviousSeasonPeak:        h.getFormIntValue(r, FormFieldTwosPreviousPeak),
		TwosAllTimePeak:               h.getFormIntValue(r, FormFieldTwosAllTimePeak),
		TwosCurrentSeasonGamesPlayed:  h.getFormIntValue(r, FormFieldTwosCurrentGames),
		TwosPreviousSeasonGamesPlayed: h.getFormIntValue(r, FormFieldTwosPreviousGames),

		ThreesCurrentSeasonPeak:         h.getFormIntValue(r, FormFieldThreesCurrentPeak),
		ThreesPreviousSeasonPeak:        h.getFormIntValue(r, FormFieldThreesPreviousPeak),
		ThreesAllTimePeak:               h.getFormIntValue(r, FormFieldThreesAllTimePeak),
		ThreesCurrentSeasonGamesPlayed:  h.getFormIntValue(r, FormFieldThreesCurrentGames),
		ThreesPreviousSeasonGamesPlayed: h.getFormIntValue(r, FormFieldThreesPreviousGames),
	}

	log.Printf("[USL-HANDLER] Building tracker from form: Discord=%s, URL=%s", tracker.DiscordID, tracker.URL)
	return tracker
}

// validateTrackerWithMetrics performs validation with comprehensive metrics and security monitoring
func (h *BaseHandler) validateTrackerWithMetrics(r *http.Request, tracker *USLUserTracker) ValidationResult {
	start := time.Now()

	validation := h.validateTracker(tracker)
	duration := time.Since(start)

	event := &ValidationEvent{
		DiscordID:  tracker.DiscordID,
		URL:        tracker.URL,
		Duration:   duration,
		Timestamp:  time.Now(),
		UserAgent:  r.Header.Get("User-Agent"),
		RemoteAddr: r.RemoteAddr,
		FormFields: make(map[string]interface{}),
	}

	// Populate form fields (excluding sensitive data)
	for field, value := range r.Form {
		if !strings.Contains(strings.ToLower(field), "password") && !strings.Contains(strings.ToLower(field), "token") {
			if len(value) > 0 {
				event.FormFields[field] = value[0]
			}
		}
	}

	// Check for security incidents
	if h.detectSecurityThreats(tracker) {
		event.Type = "security_incident"
		event.SecurityReason = h.getSecurityReason(tracker)
		event.Errors = validation.Errors

		validationLogger.Warn("Security incident detected during validation",
			"component", "validation",
			"event_type", "security_incident",
			"duration", duration,
			"timestamp", event.Timestamp,
			"remote_addr", event.RemoteAddr,
			"error_count", len(validation.Errors),
			"error_codes", h.getErrorCodes(validation.Errors),
			"error_fields", h.getErrorFields(validation.Errors),
			"security_reason", event.SecurityReason)
	} else if validation.IsValid {
		event.Type = "success"
		validationLogger.Debug("Validation completed",
			"component", "validation",
			"event_type", "success",
			"duration", duration,
			"timestamp", event.Timestamp,
			"remote_addr", event.RemoteAddr)
	} else {
		event.Type = "failure"
		event.Errors = validation.Errors
		validationLogger.Info("Validation failed",
			"component", "validation",
			"event_type", "failure",
			"duration", duration,
			"timestamp", event.Timestamp,
			"remote_addr", event.RemoteAddr,
			"error_count", len(validation.Errors),
			"error_codes", h.getErrorCodes(validation.Errors),
			"error_fields", h.getErrorFields(validation.Errors))
	}

	recordValidationMetrics(event)
	return validation
}

// validateTracker performs comprehensive validation of tracker data
func (h *BaseHandler) validateTracker(tracker *USLUserTracker) ValidationResult {
	var errors []ValidationError

	// Validate Discord ID
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

	// Validate URL
	if tracker.URL == "" {
		errors = append(errors, ValidationError{
			Field:   "url",
			Message: "Tracker URL is required",
			Code:    ValidationCodeRequired,
		})
	} else if !h.isValidTrackerURL(tracker.URL) {
		errors = append(errors, ValidationError{
			Field:   "url",
			Message: "Invalid tracker URL format",
			Code:    ValidationCodeInvalidURL,
		})
	}

	// Validate MMR values for each playlist
	errors = append(errors, h.validatePlaylistMMR("1v1", tracker.OnesCurrentSeasonPeak, tracker.OnesPreviousSeasonPeak, tracker.OnesAllTimePeak)...)
	errors = append(errors, h.validatePlaylistMMR("2v2", tracker.TwosCurrentSeasonPeak, tracker.TwosPreviousSeasonPeak, tracker.TwosAllTimePeak)...)
	errors = append(errors, h.validatePlaylistMMR("3v3", tracker.ThreesCurrentSeasonPeak, tracker.ThreesPreviousSeasonPeak, tracker.ThreesAllTimePeak)...)

	// Validate games played for each playlist
	errors = append(errors, h.validateGamesPlayed("1v1", tracker.OnesCurrentSeasonGamesPlayed, tracker.OnesPreviousSeasonGamesPlayed)...)
	errors = append(errors, h.validateGamesPlayed("2v2", tracker.TwosCurrentSeasonGamesPlayed, tracker.TwosPreviousSeasonGamesPlayed)...)
	errors = append(errors, h.validateGamesPlayed("3v3", tracker.ThreesCurrentSeasonGamesPlayed, tracker.ThreesPreviousSeasonGamesPlayed)...)

	// Check if tracker has meaningful data
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

// validatePlaylistMMR validates MMR values for a specific playlist
func (h *BaseHandler) validatePlaylistMMR(playlist string, current, previous, allTime int) []ValidationError {
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

	// Logical validation: all-time should be >= current and previous
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

// validateGamesPlayed validates games played values for a specific playlist
func (h *BaseHandler) validateGamesPlayed(playlist string, current, previous int) []ValidationError {
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
func (h *BaseHandler) hasNoPlaylistData(tracker *USLUserTracker) bool {
	return tracker.OnesCurrentSeasonPeak == 0 && tracker.OnesPreviousSeasonPeak == 0 && tracker.OnesAllTimePeak == 0 &&
		tracker.TwosCurrentSeasonPeak == 0 && tracker.TwosPreviousSeasonPeak == 0 && tracker.TwosAllTimePeak == 0 &&
		tracker.ThreesCurrentSeasonPeak == 0 && tracker.ThreesPreviousSeasonPeak == 0 && tracker.ThreesAllTimePeak == 0
}

// renderFormWithErrors renders a form template with validation errors displayed
func (h *BaseHandler) renderFormWithErrors(w http.ResponseWriter, templateName TemplateName, tracker *USLUserTracker, errors []ValidationError) {
	data := TrackerFormData{
		BasePageData: BasePageData{
			Title:       "Tracker Form",
			CurrentPage: "trackers",
		},
		Tracker: tracker,
		Errors:  h.buildErrorMap(errors),
	}

	h.renderTemplate(w, templateName, data)
}

// buildErrorMap creates a map for easy error lookup in templates
func (h *BaseHandler) buildErrorMap(errors []ValidationError) map[string]string {
	errorMap := make(map[string]string)
	for _, err := range errors {
		if existing, exists := errorMap[err.Field]; exists {
			errorMap[err.Field] = existing + "; " + err.Message
		} else {
			errorMap[err.Field] = err.Message
		}
	}
	return errorMap
}

// calculateEffectiveMMR calculates the effective MMR for a tracker (currently just logs it)
func (h *BaseHandler) calculateEffectiveMMR(tracker *USLUserTracker) {
	maxMMR := 0

	if tracker.OnesCurrentSeasonPeak > maxMMR {
		maxMMR = tracker.OnesCurrentSeasonPeak
	}
	if tracker.TwosCurrentSeasonPeak > maxMMR {
		maxMMR = tracker.TwosCurrentSeasonPeak
	}
	if tracker.ThreesCurrentSeasonPeak > maxMMR {
		maxMMR = tracker.ThreesCurrentSeasonPeak
	}

	// Note: USLUserTracker doesn't have a CalculatedMMR field, so we just log the result
	log.Printf("[USL-HANDLER] Calculated effective MMR for tracker: %d", maxMMR)
}

// Security validation helpers

func (h *BaseHandler) detectSecurityThreats(tracker *USLUserTracker) bool {
	return h.containsSQLInjection(tracker.DiscordID) ||
		h.containsXSS(tracker.DiscordID) ||
		h.containsXSS(tracker.URL) ||
		h.containsPathTraversal(tracker.URL)
}

func (h *BaseHandler) getSecurityReason(tracker *USLUserTracker) string {
	if h.containsSQLInjection(tracker.DiscordID) {
		return "SQL injection attempt detected"
	}
	if h.containsXSS(tracker.DiscordID) || h.containsXSS(tracker.URL) {
		return "XSS attempt detected"
	}
	if h.containsPathTraversal(tracker.URL) {
		return "Path traversal attempt detected"
	}
	return "Generic security threat detected"
}

func (h *BaseHandler) containsSQLInjection(input string) bool {
	sqlPatterns := []string{
		"'", "\"", ";", "--", "/*", "*/", "xp_", "sp_",
		"DROP", "DELETE", "INSERT", "UPDATE", "UNION", "SELECT",
		"drop", "delete", "insert", "update", "union", "select",
	}

	lowInput := strings.ToLower(input)
	for _, pattern := range sqlPatterns {
		if strings.Contains(lowInput, strings.ToLower(pattern)) {
			return true
		}
	}
	return false
}

func (h *BaseHandler) containsXSS(input string) bool {
	xssPatterns := []string{
		"<script", "</script", "javascript:", "vbscript:", "onload=", "onerror=",
		"onclick=", "onmouseover=", "alert(", "eval(", "document.cookie",
	}

	lowInput := strings.ToLower(input)
	for _, pattern := range xssPatterns {
		if strings.Contains(lowInput, pattern) {
			return true
		}
	}
	return false
}

func (h *BaseHandler) containsPathTraversal(input string) bool {
	pathTraversalPatterns := []string{
		"../", "..\\", "..", "/etc/passwd", "/windows/system32", "cmd.exe",
	}

	lowInput := strings.ToLower(input)
	for _, pattern := range pathTraversalPatterns {
		if strings.Contains(lowInput, pattern) {
			return true
		}
	}
	return false
}

func (h *BaseHandler) isValidTrackerURL(trackerURL string) bool {
	if trackerURL == "" {
		return false
	}

	parsedURL, err := url.Parse(trackerURL)
	if err != nil {
		return false
	}

	return parsedURL.Scheme == "http" || parsedURL.Scheme == "https"
}

func (h *BaseHandler) getErrorCodes(errors []ValidationError) []string {
	codes := make([]string, len(errors))
	for i, err := range errors {
		codes[i] = err.Code
	}
	return codes
}

func (h *BaseHandler) getErrorFields(errors []ValidationError) []string {
	fields := make([]string, len(errors))
	for i, err := range errors {
		fields[i] = err.Field
	}
	return fields
}

// updateUSLUserTrueSkillFromTrackers updates TrueSkill for a USL user from their tracker data
func (h *BaseHandler) updateUSLUserTrueSkillFromTrackers(discordID string) *services.TrueSkillUpdateResult {
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

		user, err := h.uslRepo.GetUserByDiscordID(discordID)
		if err != nil || user == nil {
			return &services.TrueSkillUpdateResult{
				Success:     false,
				HadTrackers: false,
				Error:       "user not found for default TrueSkill assignment",
			}
		}

		user.TrueSkillMu = defaultMu
		user.TrueSkillSigma = defaultSigma
		timeStr := time.Now().Format(time.RFC3339)
		user.TrueSkillLastUpdated = &timeStr

		if _, err := h.uslRepo.UpdateUser(user.ID, user.Name, user.Active, user.Banned); err != nil {
			return &services.TrueSkillUpdateResult{
				Success:     false,
				HadTrackers: false,
				Error:       fmt.Sprintf("failed to save default TrueSkill: %v", err),
			}
		}

		return &services.TrueSkillUpdateResult{
			Success:     true,
			HadTrackers: false,
			TrueSkillResult: &services.TrueSkillCalculation{
				Mu:    defaultMu,
				Sigma: defaultSigma,
			},
		}
	}

	// Call TrueSkill service (it handles tracker data retrieval internally)
	result := h.trueskillService.UpdateUserTrueSkillFromTrackers(discordID)

	// Log the result for monitoring
	if result.Success && result.TrueSkillResult != nil {
		validationLogger.Info("TrueSkill updated successfully",
			"discord_id", discordID,
			"mu", result.TrueSkillResult.Mu,
			"sigma", result.TrueSkillResult.Sigma,
			"tracker_count", len(userTrackers))
	} else {
		validationLogger.Error("TrueSkill update failed",
			"discord_id", discordID,
			"error", result.Error,
			"tracker_count", len(userTrackers))
	}

	return result
}

// transformUSLTrackerToTrackerData converts USL tracker format to TrueSkill service input format
func (h *BaseHandler) transformUSLTrackerToTrackerData(uslTracker *USLUserTracker) *services.TrackerData {
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
