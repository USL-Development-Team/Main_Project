package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
	"usl-server/internal/config"
	"usl-server/internal/services"
	"usl-server/internal/usl"
)

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
	// NOTE: auth is now handled by unified Discord OAuth system in main.go
	return &MigrationHandler{
		uslRepo:          uslRepo,
		templates:        templates,
		trueskillService: trueskillService,
		config:           config,
	}
}

func (h *MigrationHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.uslRepo.GetAllUsers()
	if err != nil {
		log.Printf("Error fetching USL users: %v", err)
		http.Error(w, "Failed to load users", http.StatusInternalServerError)
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

	h.renderTemplate(w, "usl_users.html", data)
}

func (h *MigrationHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Redirect(w, r, "/usl/users", http.StatusSeeOther)
		return
	}

	users, err := h.uslRepo.SearchUsers(query)
	if err != nil {
		log.Printf("Error searching USL users: %v", err)
		http.Error(w, "Failed to search users", http.StatusInternalServerError)
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

	h.renderTemplate(w, "usl_users.html", data)
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

	h.renderTemplate(w, "usl_user_form.html", data)
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

	name := r.FormValue("name")
	discordID := r.FormValue("discord_id")
	active := r.FormValue("active") == "true"
	banned := r.FormValue("banned") == "true"

	if name == "" || discordID == "" {
		http.Error(w, "Name and Discord ID are required", http.StatusBadRequest)
		return
	}

	user, err := h.uslRepo.CreateUser(name, discordID, active, banned)
	if err != nil {
		log.Printf("Error creating USL user: %v", err)
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}

	log.Printf("Created USL user: %s (ID: %d)", user.Name, user.ID)
	http.Redirect(w, r, "/usl/users", http.StatusSeeOther)
}

func (h *MigrationHandler) EditUserForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		http.Error(w, "User ID is required", http.StatusBadRequest)
		return
	}

	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.uslRepo.GetUserByID(userID)
	if err != nil {
		log.Printf("Error fetching USL user for edit: %v", err)
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

	h.renderTemplate(w, "usl_user_edit_form.html", data)
}

func (h *MigrationHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	userIDStr := r.FormValue("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	name := r.FormValue("name")
	active := r.FormValue("active") == "true"
	banned := r.FormValue("banned") == "true"

	user, err := h.uslRepo.UpdateUser(userID, name, active, banned)
	if err != nil {
		log.Printf("Error updating USL user: %v", err)
		http.Error(w, "Failed to update user", http.StatusInternalServerError)
		return
	}

	log.Printf("Updated USL user: %s (ID: %d)", user.Name, user.ID)
	http.Redirect(w, r, "/usl/users", http.StatusSeeOther)
}

func (h *MigrationHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userIDStr := r.FormValue("id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	err = h.uslRepo.DeleteUser(userID)
	if err != nil {
		log.Printf("Error deleting USL user: %v", err)
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	log.Printf("Deleted USL user ID: %d", userID)
	http.Redirect(w, r, "/usl/users", http.StatusSeeOther)
}

func (h *MigrationHandler) ListTrackers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	trackers, err := h.uslRepo.GetAllTrackers()
	if err != nil {
		log.Printf("Error fetching USL trackers: %v", err)
		http.Error(w, "Failed to load trackers", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title    string
		Trackers []*usl.USLUserTracker
	}{
		Title:    "USL Trackers",
		Trackers: trackers,
	}

	h.renderTemplate(w, "usl_trackers.html", data)
}

func (h *MigrationHandler) NewTrackerForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.uslRepo.GetAllUsers()
	if err != nil {
		log.Printf("Error fetching users for tracker form: %v", err)
		http.Error(w, "Failed to load users", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
		Users []*usl.USLUser
	}{
		Title: "Add USL Tracker",
		Users: users,
	}

	h.renderTemplate(w, "usl_tracker_form.html", data)
}

func (h *MigrationHandler) CreateTracker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	parseInt := func(value string) int {
		if value == "" {
			return 0
		}
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
		return 0
	}

	tracker := &usl.USLUserTracker{
		DiscordID:                       r.FormValue("discord_id"),
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
		Valid:                           r.FormValue("valid") == "true",
	}

	if tracker.DiscordID == "" || tracker.URL == "" {
		http.Error(w, "Discord ID and URL are required", http.StatusBadRequest)
		return
	}

	created, err := h.uslRepo.CreateTracker(tracker)
	if err != nil {
		log.Printf("Error creating USL tracker: %v", err)
		http.Error(w, "Failed to create tracker", http.StatusInternalServerError)
		return
	}

	log.Printf("Created USL tracker for Discord ID: %s (ID: %d)", created.DiscordID, created.ID)

	// NEW: Trigger TrueSkill update with comprehensive error handling
	h.performTrueSkillUpdate(created)

	http.Redirect(w, r, "/usl/trackers", http.StatusSeeOther)
}

func (h *MigrationHandler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats, err := h.uslRepo.GetStats()
	if err != nil {
		log.Printf("Error fetching USL stats: %v", err)
		http.Error(w, "Failed to load dashboard", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
		Stats map[string]interface{}
	}{
		Title: "USL Admin Dashboard",
		Stats: stats,
	}

	h.renderTemplate(w, "usl_admin_dashboard.html", data)
}

func (h *MigrationHandler) ImportData(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Import USL Data",
	}

	h.renderTemplate(w, "usl_import.html", data)
}

func (h *MigrationHandler) ListUsersAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.uslRepo.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to load users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func (h *MigrationHandler) ListTrackersAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	trackers, err := h.uslRepo.GetAllTrackers()
	if err != nil {
		http.Error(w, "Failed to load trackers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(trackers)
}

func (h *MigrationHandler) GetLeaderboardAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.uslRepo.GetLeaderboard()
	if err != nil {
		http.Error(w, "Failed to load leaderboard", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// performTrueSkillUpdate handles TrueSkill calculation and synchronization with comprehensive error handling
func (h *MigrationHandler) performTrueSkillUpdate(tracker *usl.USLUserTracker) {
	if tracker == nil {
		log.Printf("ERROR: Cannot perform TrueSkill update - tracker is nil")
		return
	}

	// Step 1: Transform tracker data
	trackerData := h.mapUSLTrackerToTrackerData(tracker)
	log.Printf("Starting TrueSkill update for %s (tracker ID: %d)", tracker.DiscordID, tracker.ID)

	// Step 2: Attempt TrueSkill calculation
	result := h.trueskillService.UpdateUserTrueSkillFromTrackerData(trackerData)

	if !result.Success {
		log.Printf("WARNING: TrueSkill update failed for %s: %s", tracker.DiscordID, result.Error)
		log.Printf("  - Tracker was still created successfully (ID: %d)", tracker.ID)
		log.Printf("  - Manual TrueSkill update may be required later")
		return
	}

	// Step 3: TrueSkill calculation succeeded - log results
	log.Printf("SUCCESS: TrueSkill calculated for %s: μ=%.1f, σ=%.2f",
		tracker.DiscordID, result.TrueSkillResult.Mu, result.TrueSkillResult.Sigma)

	// Step 4: Attempt to sync results to USL user table
	err := h.uslRepo.UpdateUserTrueSkill(
		tracker.DiscordID,
		result.TrueSkillResult.Mu,
		result.TrueSkillResult.Sigma,
	)

	if err != nil {
		log.Printf("WARNING: Failed to sync TrueSkill to USL user table for %s: %v", tracker.DiscordID, err)
		log.Printf("  - TrueSkill was calculated and stored in core tables")
		log.Printf("  - USL table sync failed - data inconsistency possible")
		log.Printf("  - Manual sync may be required")
	} else {
		log.Printf("SUCCESS: Full TrueSkill integration completed for %s", tracker.DiscordID)
		log.Printf("  - Core tables updated: ✓")
		log.Printf("  - USL tables updated: ✓")
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

func (h *MigrationHandler) renderTemplate(w http.ResponseWriter, templateName string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, templateName, data); err != nil {
		log.Printf("Template rendering error (%s): %v", templateName, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
