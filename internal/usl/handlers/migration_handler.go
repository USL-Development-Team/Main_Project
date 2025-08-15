package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
	"usl-server/internal/config"
	"usl-server/internal/models"
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
	TemplateUSLUsers          TemplateName = "users-list-page"
	TemplateUSLUsersTable     TemplateName = "users-table-fragment"
	TemplateUSLUserDetail     TemplateName = "user-detail-page"
	TemplateUSLTrackers       TemplateName = "trackers-list-page"
	TemplateUSLTrackerDetail  TemplateName = "tracker-detail-page"
	TemplateUSLTrackerNew     TemplateName = "tracker-new-page"
	TemplateUSLTrackerEdit    TemplateName = "tracker-edit-page"
	TemplateUSLAdminDashboard TemplateName = "admin-dashboard-page"
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
			SearchTarget:      "#users-table",
			ClearURL:          "/usl/users",
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
			SearchTarget:      "#users-table",
			ClearURL:          "/usl/users",
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
			ClearURL:          "/usl/trackers",
			ShowFilters:       false, // Trackers don't have status filters like users
			Query:             "",
			StatusFilter:      "",
		},
	}

	h.renderTemplate(w, TemplateUSLTrackers, data)
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
	}{
		Title:       "New Tracker",
		CurrentPage: "trackers",
	}

	h.renderTemplate(w, TemplateUSLTrackerNew, data)
}

func (h *MigrationHandler) CreateTracker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.handleMethodNotAllowed(w, r)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		h.handleInvalidFormData(w, err)
		return
	}

	// Extract and validate required fields
	discordID := r.FormValue("discord_id")
	if discordID == "" {
		http.Error(w, "Discord ID is required", http.StatusBadRequest)
		return
	}

	// Parse integer fields with helper function
	parseInt := func(value string) int {
		if value == "" {
			return 0
		}
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
		return 0
	}

	// Parse boolean field
	parseBool := func(value string) bool {
		return value == "true"
	}

	// Create tracker with all MMR fields
	tracker := &usl.USLUserTracker{
		DiscordID:                       discordID,
		URL:                             r.FormValue("url"),
		OnesCurrentSeasonPeak:           parseInt(r.FormValue("ones_current_peak")),
		OnesPreviousSeasonPeak:          parseInt(r.FormValue("ones_previous_peak")),
		OnesAllTimePeak:                 parseInt(r.FormValue("ones_all_time_peak")),
		OnesCurrentSeasonGamesPlayed:    parseInt(r.FormValue("ones_current_games")),
		OnesPreviousSeasonGamesPlayed:   parseInt(r.FormValue("ones_previous_games")),
		TwosCurrentSeasonPeak:           parseInt(r.FormValue("twos_current_peak")),
		TwosPreviousSeasonPeak:          parseInt(r.FormValue("twos_previous_peak")),
		TwosAllTimePeak:                 parseInt(r.FormValue("twos_all_time_peak")),
		TwosCurrentSeasonGamesPlayed:    parseInt(r.FormValue("twos_current_games")),
		TwosPreviousSeasonGamesPlayed:   parseInt(r.FormValue("twos_previous_games")),
		ThreesCurrentSeasonPeak:         parseInt(r.FormValue("threes_current_peak")),
		ThreesPreviousSeasonPeak:        parseInt(r.FormValue("threes_previous_peak")),
		ThreesAllTimePeak:               parseInt(r.FormValue("threes_all_time_peak")),
		ThreesCurrentSeasonGamesPlayed:  parseInt(r.FormValue("threes_current_games")),
		ThreesPreviousSeasonGamesPlayed: parseInt(r.FormValue("threes_previous_games")),
		Valid:                           parseBool(r.FormValue("valid")),
	}

	// Calculate effective MMR using the same logic as the JavaScript form
	weightedSum := (tracker.OnesCurrentSeasonPeak * tracker.OnesCurrentSeasonGamesPlayed) +
		(tracker.TwosCurrentSeasonPeak * tracker.TwosCurrentSeasonGamesPlayed) +
		(tracker.ThreesCurrentSeasonPeak * tracker.ThreesCurrentSeasonGamesPlayed)

	totalGames := tracker.OnesCurrentSeasonGamesPlayed + tracker.TwosCurrentSeasonGamesPlayed + tracker.ThreesCurrentSeasonGamesPlayed

	if totalGames > 0 {
		tracker.MMR = weightedSum / totalGames
	} else {
		tracker.MMR = 0
	}

	createdTracker, err := h.uslRepo.CreateTracker(tracker)
	if err != nil {
		h.handleDatabaseError(w, "create tracker", err)
		return
	}

	// Redirect to tracker detail
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

	data := struct {
		Title       string
		CurrentPage string
		Tracker     *usl.USLUserTracker
	}{
		Title:       "Edit Tracker",
		CurrentPage: "trackers",
		Tracker:     tracker,
	}

	h.renderTemplate(w, TemplateUSLTrackerEdit, data)
}

func (h *MigrationHandler) UpdateTracker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.handleMethodNotAllowed(w, r)
		return
	}

	// Parse form data
	if err := r.ParseForm(); err != nil {
		h.handleInvalidFormData(w, err)
		return
	}

	// Extract tracker ID
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

	// Extract and validate required fields
	discordID := r.FormValue("discord_id")
	if discordID == "" {
		http.Error(w, "Discord ID is required", http.StatusBadRequest)
		return
	}

	// Parse integer fields with helper function
	parseInt := func(value string) int {
		if value == "" {
			return 0
		}
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
		return 0
	}

	// Parse boolean field
	parseBool := func(value string) bool {
		return value == "true"
	}

	// Update tracker with all MMR fields
	tracker := &usl.USLUserTracker{
		ID:                              trackerID,
		DiscordID:                       discordID,
		URL:                             r.FormValue("url"),
		OnesCurrentSeasonPeak:           parseInt(r.FormValue("ones_current_peak")),
		OnesPreviousSeasonPeak:          parseInt(r.FormValue("ones_previous_peak")),
		OnesAllTimePeak:                 parseInt(r.FormValue("ones_all_time_peak")),
		OnesCurrentSeasonGamesPlayed:    parseInt(r.FormValue("ones_current_games")),
		OnesPreviousSeasonGamesPlayed:   parseInt(r.FormValue("ones_previous_games")),
		TwosCurrentSeasonPeak:           parseInt(r.FormValue("twos_current_peak")),
		TwosPreviousSeasonPeak:          parseInt(r.FormValue("twos_previous_peak")),
		TwosAllTimePeak:                 parseInt(r.FormValue("twos_all_time_peak")),
		TwosCurrentSeasonGamesPlayed:    parseInt(r.FormValue("twos_current_games")),
		TwosPreviousSeasonGamesPlayed:   parseInt(r.FormValue("twos_previous_games")),
		ThreesCurrentSeasonPeak:         parseInt(r.FormValue("threes_current_peak")),
		ThreesPreviousSeasonPeak:        parseInt(r.FormValue("threes_previous_peak")),
		ThreesAllTimePeak:               parseInt(r.FormValue("threes_all_time_peak")),
		ThreesCurrentSeasonGamesPlayed:  parseInt(r.FormValue("threes_current_games")),
		ThreesPreviousSeasonGamesPlayed: parseInt(r.FormValue("threes_previous_games")),
		Valid:                           parseBool(r.FormValue("valid")),
	}

	// Calculate effective MMR using the same logic as the JavaScript form
	weightedSum := (tracker.OnesCurrentSeasonPeak * tracker.OnesCurrentSeasonGamesPlayed) +
		(tracker.TwosCurrentSeasonPeak * tracker.TwosCurrentSeasonGamesPlayed) +
		(tracker.ThreesCurrentSeasonPeak * tracker.ThreesCurrentSeasonGamesPlayed)

	totalGames := tracker.OnesCurrentSeasonGamesPlayed + tracker.TwosCurrentSeasonGamesPlayed + tracker.ThreesCurrentSeasonGamesPlayed

	if totalGames > 0 {
		tracker.MMR = weightedSum / totalGames
	} else {
		tracker.MMR = 0
	}

	err = h.uslRepo.UpdateTracker(tracker)
	if err != nil {
		h.handleDatabaseError(w, "update tracker", err)
		return
	}

	// Redirect to tracker detail
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
		Stats       map[string]interface{}
	}{
		Title:       "Dashboard",
		CurrentPage: "admin",
		Stats:       stats,
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
	buf.WriteTo(w)
}
