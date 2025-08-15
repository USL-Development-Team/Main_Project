package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
	"usl-server/internal/config"
	"usl-server/internal/services"
	"usl-server/internal/usl"
)

// USL-specific constants for type safety
const (
	USLDiscordGuildID = "1390537743385231451" // USL Discord Guild ID
)

// FormField represents typed form field names
type FormField string

const (
	FormFieldDiscordID           FormField = "discord_id"
	FormFieldURL                 FormField = "url"
	FormFieldName                FormField = "name"
	FormFieldActive              FormField = "active"
	FormFieldBanned              FormField = "banned"
	FormFieldValid               FormField = "valid"
	FormFieldID                  FormField = "id"
	FormFieldOnesCurrentPeak     FormField = "ones_current_peak"
	FormFieldOnesPreviousPeak    FormField = "ones_previous_peak"
	FormFieldOnesAllTimePeak     FormField = "ones_all_time_peak"
	FormFieldOnesCurrentGames    FormField = "ones_current_games"
	FormFieldOnesPreviousGames   FormField = "ones_previous_games"
	FormFieldTwosCurrentPeak     FormField = "twos_current_peak"
	FormFieldTwosPreviousPeak    FormField = "twos_previous_peak"
	FormFieldTwosAllTimePeak     FormField = "twos_all_time_peak"
	FormFieldTwosCurrentGames    FormField = "twos_current_games"
	FormFieldTwosPreviousGames   FormField = "twos_previous_games"
	FormFieldThreesCurrentPeak   FormField = "threes_current_peak"
	FormFieldThreesPreviousPeak  FormField = "threes_previous_peak"
	FormFieldThreesAllTimePeak   FormField = "threes_all_time_peak"
	FormFieldThreesCurrentGames  FormField = "threes_current_games"
	FormFieldThreesPreviousGames FormField = "threes_previous_games"
)

// TemplateName represents typed template names
type TemplateName string

const (
	TemplateUSLUsers          TemplateName = "usl_users.html"
	TemplateUSLUserForm       TemplateName = "usl_user_form.html"
	TemplateUSLUserEditForm   TemplateName = "usl_user_edit_form.html"
	TemplateUSLTrackers       TemplateName = "usl_trackers.html"
	TemplateUSLTrackerForm    TemplateName = "usl_tracker_form.html"
	TemplateUSLAdminDashboard TemplateName = "usl_admin_dashboard.html"
	TemplateUSLImport         TemplateName = "usl_import.html"
)

// parseIntField safely converts a form field to int, returning 0 for empty or invalid values
func parseIntField(value string) int {
	if value == "" {
		return 0
	}
	if i, err := strconv.Atoi(value); err == nil {
		return i
	}
	return 0
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

// validateTrackerRequired validates required tracker fields
func (h *MigrationHandler) validateTrackerRequired(w http.ResponseWriter, tracker *usl.USLUserTracker) bool {
	if tracker.DiscordID == "" || tracker.URL == "" {
		log.Printf("[USL-HANDLER] Tracker validation failed: Discord=%s, URL=%s", tracker.DiscordID, tracker.URL)
		http.Error(w, "Discord ID and URL are required", http.StatusBadRequest)
		return false
	}
	log.Printf("[USL-HANDLER] Tracker validation passed: Discord=%s", tracker.DiscordID)
	return true
}

// logTrueSkillUpdateFailure logs structured failure information
func (h *MigrationHandler) logTrueSkillUpdateFailure(tracker *usl.USLUserTracker, errorMsg string) {
	log.Printf("[USL-TRUESKILL] WARNING: Update failed - Discord=%s, TrackerID=%d, Error=%s",
		tracker.DiscordID, tracker.ID, errorMsg)
	log.Printf("[USL-TRUESKILL] Tracker creation succeeded, TrueSkill update failed - manual intervention may be needed")
}

// logTrueSkillUpdateSuccess logs structured success information
func (h *MigrationHandler) logTrueSkillUpdateSuccess(tracker *usl.USLUserTracker, mu, sigma float64) {
	log.Printf("[USL-TRUESKILL] SUCCESS: Calculated - Discord=%s, TrackerID=%d, μ=%.1f, σ=%.2f",
		tracker.DiscordID, tracker.ID, mu, sigma)
}

// logUSLSyncFailure logs structured USL sync failure information
func (h *MigrationHandler) logUSLSyncFailure(discordID string, err error) {
	log.Printf("[USL-TRUESKILL] WARNING: USL table sync failed - Discord=%s, Error=%v", discordID, err)
	log.Printf("[USL-TRUESKILL] Core tables updated successfully, USL tables inconsistent - manual sync required")
}

// logUSLSyncSuccess logs structured USL sync success information
func (h *MigrationHandler) logUSLSyncSuccess(discordID string) {
	log.Printf("[USL-TRUESKILL] SUCCESS: Full integration completed - Discord=%s (Core ✓, USL ✓)", discordID)
}

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
		Title string
		Users []*usl.USLUser
		Query string
	}{
		Title: "USL Users",
		Users: users,
		Query: "", // Empty query for full user list
	}

	h.renderTemplate(w, TemplateUSLUsers, data)
}

func (h *MigrationHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Redirect(w, r, "/usl/users", http.StatusSeeOther)
		return
	}

	users, err := h.uslRepo.SearchUsers(query)
	if err != nil {
		h.handleDatabaseError(w, "search users", err)
		return
	}

	data := struct {
		Title string
		Users []*usl.USLUser
		Query string
	}{
		Title: "USL User Search Results",
		Users: users,
		Query: query,
	}

	h.renderTemplate(w, TemplateUSLUsers, data)
}

func (h *MigrationHandler) NewUserForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Add USL User",
	}

	h.renderTemplate(w, TemplateUSLUserForm, data)
}

func (h *MigrationHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := h.getFormValue(r, FormFieldName)
	discordID := h.getFormValue(r, FormFieldDiscordID)
	active := h.getFormBoolValue(r, FormFieldActive)
	banned := h.getFormBoolValue(r, FormFieldBanned)

	if name == "" || discordID == "" {
		http.Error(w, "Name and Discord ID are required", http.StatusBadRequest)
		return
	}

	user, err := h.uslRepo.CreateUser(name, discordID, active, banned)
	if err != nil {
		h.handleDatabaseError(w, "create user", err)
		return
	}

	log.Printf("[USL-HANDLER] Created user: %s (ID: %d, Discord: %s)", user.Name, user.ID, user.DiscordID)
	http.Redirect(w, r, "/usl/users", http.StatusSeeOther)
}

func (h *MigrationHandler) EditUserForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	userIDStr := r.URL.Query().Get(string(FormFieldID))
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
		log.Printf("[USL-HANDLER] User lookup failed: ID=%d, error=%v", userID, err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	data := struct {
		Title string
		User  *usl.USLUser
	}{
		Title: "Edit USL User",
		User:  user,
	}

	h.renderTemplate(w, TemplateUSLUserEditForm, data)
}

func (h *MigrationHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.handleMethodNotAllowed(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.handleInvalidFormData(w, err)
		return
	}

	userIDStr := h.getFormValue(r, FormFieldID)
	userID, err := h.parseUserID(userIDStr)
	if err != nil {
		h.handleParseError(w, "user ID")
		return
	}

	name := h.getFormValue(r, FormFieldName)
	active := h.getFormBoolValue(r, FormFieldActive)
	banned := h.getFormBoolValue(r, FormFieldBanned)

	user, err := h.uslRepo.UpdateUser(userID, name, active, banned)
	if err != nil {
		h.handleDatabaseError(w, "update user", err)
		return
	}

	log.Printf("[USL-HANDLER] Updated user: %s (ID: %d, Discord: %s)", user.Name, user.ID, user.DiscordID)
	http.Redirect(w, r, "/usl/users", http.StatusSeeOther)
}

func (h *MigrationHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.handleMethodNotAllowed(w, r)
		return
	}

	userIDStr := h.getFormValue(r, FormFieldID)
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

	log.Printf("[USL-HANDLER] Deleted user: ID=%d", userID)
	http.Redirect(w, r, "/usl/users", http.StatusSeeOther)
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
		Title    string
		Trackers []*usl.USLUserTracker
	}{
		Title:    "USL Trackers",
		Trackers: trackers,
	}

	h.renderTemplate(w, TemplateUSLTrackers, data)
}

func (h *MigrationHandler) NewTrackerForm(w http.ResponseWriter, r *http.Request) {
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
		Title string
		Users []*usl.USLUser
	}{
		Title: "Add USL Tracker",
		Users: users,
	}

	h.renderTemplate(w, TemplateUSLTrackerForm, data)
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

	if !h.validateTrackerRequired(w, tracker) {
		return
	}

	created, err := h.uslRepo.CreateTracker(tracker)
	if err != nil {
		h.handleDatabaseError(w, "create tracker", err)
		return
	}

	log.Printf("[USL-HANDLER] Created tracker: Discord=%s, ID=%d, URL=%s", created.DiscordID, created.ID, created.URL)

	h.performTrueSkillUpdate(created)

	http.Redirect(w, r, "/usl/trackers", http.StatusSeeOther)
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

	data := struct {
		Title string
		Stats map[string]interface{}
	}{
		Title: "USL Admin Dashboard",
		Stats: stats,
	}

	h.renderTemplate(w, TemplateUSLAdminDashboard, data)
}

func (h *MigrationHandler) ImportData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Import USL Data",
	}

	h.renderTemplate(w, TemplateUSLImport, data)
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
func (h *MigrationHandler) performTrueSkillUpdate(tracker *usl.USLUserTracker) {
	if tracker == nil {
		log.Printf("[USL-TRUESKILL] ERROR: Cannot perform update - tracker is nil")
		return
	}

	trackerData := h.mapUSLTrackerToTrackerData(tracker)
	log.Printf("[USL-TRUESKILL] Starting update - Discord=%s, TrackerID=%d", tracker.DiscordID, tracker.ID)

	result := h.trueskillService.UpdateUserTrueSkillFromTrackerData(trackerData)

	if !result.Success {
		h.logTrueSkillUpdateFailure(tracker, result.Error)
		return
	}

	h.logTrueSkillUpdateSuccess(tracker, result.TrueSkillResult.Mu, result.TrueSkillResult.Sigma)

	err := h.uslRepo.UpdateUserTrueSkill(
		tracker.DiscordID,
		result.TrueSkillResult.Mu,
		result.TrueSkillResult.Sigma,
	)

	if err != nil {
		h.logUSLSyncFailure(tracker.DiscordID, err)
	} else {
		h.logUSLSyncSuccess(tracker.DiscordID)
	}
}

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

// NOTE: Authentication methods removed - now handled by unified Discord OAuth system in main.go

func (h *MigrationHandler) renderTemplate(w http.ResponseWriter, templateName TemplateName, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, string(templateName), data); err != nil {
		log.Printf("[USL-HANDLER] Template rendering error: template=%s, error=%v", templateName, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
