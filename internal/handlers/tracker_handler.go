package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"usl-server/internal/models"
	"usl-server/internal/repositories"
)

// TrackerHandler handles HTTP requests for tracker management
type TrackerHandler struct {
	trackerRepo *repositories.TrackerRepository
	templates   *template.Template
}

// NewTrackerHandler creates a new tracker handler
func NewTrackerHandler(trackerRepo *repositories.TrackerRepository, templates *template.Template) *TrackerHandler {
	return &TrackerHandler{
		trackerRepo: trackerRepo,
		templates:   templates,
	}
}

// ListTrackers displays all trackers in HTML format
func (h *TrackerHandler) ListTrackers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	validOnly := r.URL.Query().Get("valid_only") == "true"
	trackers, err := h.trackerRepo.GetAllTrackers(validOnly)
	if err != nil {
		log.Printf("Error getting trackers: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title     string
		Trackers  []*models.UserTracker
		ValidOnly bool
	}{
		Title:     "Tracker Management",
		Trackers:  trackers,
		ValidOnly: validOnly,
	}

	if err := h.templates.ExecuteTemplate(w, "trackers.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// NewTrackerForm displays the form for creating a new tracker
func (h *TrackerHandler) NewTrackerForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Create New Tracker",
	}

	if err := h.templates.ExecuteTemplate(w, "tracker_form.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// CreateTracker handles the creation of a new tracker
func (h *TrackerHandler) CreateTracker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse form values to integers
	parseInt := func(value string) int {
		if value == "" {
			return 0
		}
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
		return 0
	}

	trackerData := models.TrackerCreateRequest{
		DiscordID:                 r.FormValue("discord_id"),
		URL:                       r.FormValue("url"),
		OnesCurrentSeasonPeak:     parseInt(r.FormValue("ones_current_peak")),
		OnesPreviousSeasonPeak:    parseInt(r.FormValue("ones_previous_peak")),
		OnesAllTimePeak:           parseInt(r.FormValue("ones_all_time_peak")),
		OnesCurrentSeasonGames:    parseInt(r.FormValue("ones_current_games")),
		OnesPreviousSeasonGames:   parseInt(r.FormValue("ones_previous_games")),
		TwosCurrentSeasonPeak:     parseInt(r.FormValue("twos_current_peak")),
		TwosPreviousSeasonPeak:    parseInt(r.FormValue("twos_previous_peak")),
		TwosAllTimePeak:           parseInt(r.FormValue("twos_all_time_peak")),
		TwosCurrentSeasonGames:    parseInt(r.FormValue("twos_current_games")),
		TwosPreviousSeasonGames:   parseInt(r.FormValue("twos_previous_games")),
		ThreesCurrentSeasonPeak:   parseInt(r.FormValue("threes_current_peak")),
		ThreesPreviousSeasonPeak:  parseInt(r.FormValue("threes_previous_peak")),
		ThreesAllTimePeak:         parseInt(r.FormValue("threes_all_time_peak")),
		ThreesCurrentSeasonGames:  parseInt(r.FormValue("threes_current_games")),
		ThreesPreviousSeasonGames: parseInt(r.FormValue("threes_previous_games")),
		Valid:                     r.FormValue("valid") == "true",
	}

	tracker, err := h.trackerRepo.CreateTracker(trackerData)
	if err != nil {
		log.Printf("Error creating tracker: %v", err)
		http.Error(w, "Failed to create tracker: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Created tracker for user: %s", tracker.DiscordID)
	http.Redirect(w, r, "/trackers", http.StatusSeeOther)
}

// EditTrackerForm displays the form for editing an existing tracker
func (h *TrackerHandler) EditTrackerForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	trackerIDStr := r.URL.Query().Get("id")
	if trackerIDStr == "" {
		http.Error(w, "Tracker ID is required", http.StatusBadRequest)
		return
	}

	trackerID, err := strconv.Atoi(trackerIDStr)
	if err != nil {
		http.Error(w, "Invalid tracker ID", http.StatusBadRequest)
		return
	}

	// Get tracker by finding all trackers and filtering by ID (simplified approach)
	allTrackers, err := h.trackerRepo.GetAllTrackers(false)
	if err != nil {
		log.Printf("Error getting trackers: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	var tracker *models.UserTracker
	for _, t := range allTrackers {
		if t.ID == trackerID {
			tracker = t
			break
		}
	}

	if tracker == nil {
		http.Error(w, "Tracker not found", http.StatusNotFound)
		return
	}

	data := struct {
		Title   string
		Tracker *models.UserTracker
	}{
		Title:   "Edit Tracker",
		Tracker: tracker,
	}

	if err := h.templates.ExecuteTemplate(w, "tracker_edit_form.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// UpdateTracker handles updating an existing tracker
func (h *TrackerHandler) UpdateTracker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	trackerIDStr := r.FormValue("tracker_id")
	if trackerIDStr == "" {
		http.Error(w, "Tracker ID is required", http.StatusBadRequest)
		return
	}

	trackerID, err := strconv.Atoi(trackerIDStr)
	if err != nil {
		http.Error(w, "Invalid tracker ID", http.StatusBadRequest)
		return
	}

	// Parse form values to integers
	parseInt := func(value string) int {
		if value == "" {
			return 0
		}
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
		return 0
	}

	trackerData := models.TrackerUpdateRequest{
		URL:                       r.FormValue("url"),
		OnesCurrentSeasonPeak:     parseInt(r.FormValue("ones_current_peak")),
		OnesPreviousSeasonPeak:    parseInt(r.FormValue("ones_previous_peak")),
		OnesAllTimePeak:           parseInt(r.FormValue("ones_all_time_peak")),
		OnesCurrentSeasonGames:    parseInt(r.FormValue("ones_current_games")),
		OnesPreviousSeasonGames:   parseInt(r.FormValue("ones_previous_games")),
		TwosCurrentSeasonPeak:     parseInt(r.FormValue("twos_current_peak")),
		TwosPreviousSeasonPeak:    parseInt(r.FormValue("twos_previous_peak")),
		TwosAllTimePeak:           parseInt(r.FormValue("twos_all_time_peak")),
		TwosCurrentSeasonGames:    parseInt(r.FormValue("twos_current_games")),
		TwosPreviousSeasonGames:   parseInt(r.FormValue("twos_previous_games")),
		ThreesCurrentSeasonPeak:   parseInt(r.FormValue("threes_current_peak")),
		ThreesPreviousSeasonPeak:  parseInt(r.FormValue("threes_previous_peak")),
		ThreesAllTimePeak:         parseInt(r.FormValue("threes_all_time_peak")),
		ThreesCurrentSeasonGames:  parseInt(r.FormValue("threes_current_games")),
		ThreesPreviousSeasonGames: parseInt(r.FormValue("threes_previous_games")),
		Valid:                     r.FormValue("valid") == "true",
	}

	tracker, err := h.trackerRepo.UpdateTracker(trackerID, trackerData)
	if err != nil {
		log.Printf("Error updating tracker: %v", err)
		http.Error(w, "Failed to update tracker: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Updated tracker for user: %s", tracker.DiscordID)
	http.Redirect(w, r, "/trackers", http.StatusSeeOther)
}

// DeleteTracker handles deleting (marking invalid) a tracker
func (h *TrackerHandler) DeleteTracker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	trackerIDStr := r.URL.Query().Get("id")
	if trackerIDStr == "" {
		http.Error(w, "Tracker ID is required", http.StatusBadRequest)
		return
	}

	trackerID, err := strconv.Atoi(trackerIDStr)
	if err != nil {
		http.Error(w, "Invalid tracker ID", http.StatusBadRequest)
		return
	}

	tracker, err := h.trackerRepo.DeleteTracker(trackerID)
	if err != nil {
		log.Printf("Error deleting tracker: %v", err)
		http.Error(w, "Failed to delete tracker: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Deleted tracker for user: %s", tracker.DiscordID)
	http.Redirect(w, r, "/trackers", http.StatusSeeOther)
}

// SearchTrackers handles searching for trackers
func (h *TrackerHandler) SearchTrackers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Redirect(w, r, "/trackers", http.StatusSeeOther)
		return
	}

	trackers, err := h.trackerRepo.SearchTrackers(query, 50)
	if err != nil {
		log.Printf("Error searching trackers: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title    string
		Trackers []*models.UserTracker
		Query    string
	}{
		Title:    "Tracker Search Results",
		Trackers: trackers,
		Query:    query,
	}

	if err := h.templates.ExecuteTemplate(w, "trackers.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *TrackerHandler) ListTrackersAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	validOnly := r.URL.Query().Get("valid_only") == "true"
	trackers, err := h.trackerRepo.GetAllTrackers(validOnly)
	if err != nil {
		log.Printf("Error getting trackers: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(trackers); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
