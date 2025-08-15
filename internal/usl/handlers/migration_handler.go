package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"usl-server/internal/usl"
)

// USL-specific constants
const (
	USL_DISCORD_GUILD_ID = "1390537743385231451" // USL Discord Guild ID
)

// MigrationHandler provides simplified handlers for USL-only operations
// This is a temporary migration solution - no multi-guild complexity
// AUTH NOTE: This handler no longer manages auth - that's handled by unified Discord OAuth in main.go
type MigrationHandler struct {
	uslRepo   *usl.USLRepository
	templates *template.Template
}

func NewMigrationHandler(uslRepo *usl.USLRepository, templates *template.Template) *MigrationHandler {
	// NOTE: auth is now handled by unified Discord OAuth system in main.go
	return &MigrationHandler{
		uslRepo:   uslRepo,
		templates: templates,
	}
}

// User Management - USL Only

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

// Tracker Management - USL Only

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

	// Helper function to parse int from form
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
	http.Redirect(w, r, "/usl/trackers", http.StatusSeeOther)
}

// Admin Tools - USL Only

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

// API endpoints for USL (simple JSON responses)

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

// Helper methods

// NOTE: Authentication methods removed - now handled by unified Discord OAuth system in main.go

func (h *MigrationHandler) renderTemplate(w http.ResponseWriter, templateName string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, templateName, data); err != nil {
		log.Printf("Template rendering error (%s): %v", templateName, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
