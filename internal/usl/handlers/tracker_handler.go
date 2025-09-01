package handlers

import (
	"fmt"
	"log"
	"net/http"
)

// TrackerHandler handles all tracker-related operations
type TrackerHandler struct {
	*BaseHandler
}

// NewTrackerHandler creates a new TrackerHandler
func NewTrackerHandler(baseHandler *BaseHandler) *TrackerHandler {
	return &TrackerHandler{BaseHandler: baseHandler}
}

// ListTrackers displays all trackers in a paginated list
func (h *TrackerHandler) ListTrackers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	trackers, err := h.uslRepo.GetAllTrackers()
	if err != nil {
		h.handleDatabaseError(w, "load trackers", err)
		return
	}

	data := TrackersListData{
		BasePageData: BasePageData{
			Title:       "Trackers",
			CurrentPage: "trackers",
		},
		Trackers: trackers,
		SearchConfig: SearchConfig{
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

// populateTrackerUsers fetches and attaches user data to trackers for display
func (h *TrackerHandler) populateTrackerUsers(trackers []*USLUserTracker) {
	for _, tracker := range trackers {
		if user, err := h.uslRepo.GetUserByDiscordID(tracker.DiscordID); err == nil && user != nil {
			tracker.User = user
		}
	}
}

// SearchTrackers handles tracker search requests and returns filtered results
func (h *TrackerHandler) SearchTrackers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	query := r.URL.Query().Get("q")

	var trackers []*USLUserTracker
	var err error
	if query == "" {
		trackers, err = h.uslRepo.GetAllTrackers()
	} else {
		trackers, err = h.uslRepo.SearchTrackers(query)
	}
	if err != nil {
		h.handleDatabaseError(w, "search trackers", err)
		return
	}

	// For fragments like table updates, we may need different template data
	if r.Header.Get("HX-Request") == "true" {
		// This is an HTMX request - return just the table fragment
		data := struct {
			Title        string
			Trackers     []*USLUserTracker
			SearchConfig SearchConfig
		}{
			Title:    "Trackers",
			Trackers: trackers,
			SearchConfig: SearchConfig{
				SearchPlaceholder: "Search by URL or Discord ID...",
				SearchURL:         "/usl/trackers/search",
				SearchTarget:      "#trackers-table",
				ClearURL:          "/usl/trackers/search",
				ShowFilters:       false,
				Query:             query,
				StatusFilter:      "",
			},
		}

		h.populateTrackerUsers(trackers)
		h.renderTemplate(w, "trackers-table-fragment", data)
	} else {
		// Regular request - return full page
		data := TrackersListData{
			BasePageData: BasePageData{
				Title:       "Trackers",
				CurrentPage: "trackers",
			},
			Trackers: trackers,
			SearchConfig: SearchConfig{
				SearchPlaceholder: "Search by URL or Discord ID...",
				SearchURL:         "/usl/trackers/search",
				SearchTarget:      "#trackers-table",
				ClearURL:          "/usl/trackers/search",
				ShowFilters:       false,
				Query:             query,
				StatusFilter:      "",
			},
		}

		h.populateTrackerUsers(trackers)
		h.renderTemplate(w, TemplateUSLTrackers, data)
	}
}

// TrackerDetail displays detailed information about a specific tracker
func (h *TrackerHandler) TrackerDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	trackerID, err := h.parseUserID(r.URL.Query().Get("id"))
	if err != nil {
		h.handleInvalidID(w, "tracker")
		return
	}

	tracker, err := h.uslRepo.GetTrackerByID(trackerID)
	if err != nil {
		h.handleDatabaseError(w, "load tracker", err)
		return
	}
	if tracker == nil {
		http.NotFound(w, r)
		return
	}

	// Load associated user
	user, err := h.uslRepo.GetUserByDiscordID(tracker.DiscordID)
	if err != nil {
		h.handleDatabaseError(w, "load associated user", err)
		return
	}

	data := TrackerDetailData{
		BasePageData: BasePageData{
			Title:       "Tracker Details",
			CurrentPage: "trackers",
		},
		Tracker: tracker,
		User:    user,
	}

	h.renderTemplate(w, TemplateUSLTrackerDetail, data)
}

// NewTrackerForm displays the form for creating a new tracker
func (h *TrackerHandler) NewTrackerForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	// Pre-fill Discord ID from URL parameter if provided
	tracker := &USLUserTracker{}
	if discordID := r.URL.Query().Get("discord_id"); discordID != "" {
		if isValidDiscordID(discordID) {
			tracker.DiscordID = discordID
			log.Printf("[USL-HANDLER] Pre-filled Discord ID from URL parameter: %s", discordID)
		} else {
			log.Printf("[USL-HANDLER] Invalid Discord ID in URL parameter, ignoring: %s", discordID)
		}
	}

	data := TrackerFormData{
		BasePageData: BasePageData{
			Title:       "New Tracker",
			CurrentPage: "trackers",
		},
		Tracker: tracker,
		Errors:  make(map[string]string),
	}

	h.renderTemplate(w, TemplateUSLTrackerNew, data)
}

// CreateTracker handles the creation of a new tracker with TrueSkill auto-update
func (h *TrackerHandler) CreateTracker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.handleMethodNotAllowed(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.handleInvalidFormData(w, err)
		return
	}

	// Build tracker from form (using existing helper)
	tracker := h.buildTrackerFromForm(r)

	// Comprehensive validation with metrics and security monitoring
	validation := h.validateTrackerWithMetrics(r, tracker)
	if !validation.IsValid {
		h.renderFormWithErrors(w, TemplateUSLTrackerNew, tracker, validation.Errors)
		return
	}

	// Calculate MMR (using extracted function)
	h.calculateEffectiveMMR(tracker)

	// Save to database
	createdTracker, err := h.uslRepo.CreateTracker(tracker)
	if err != nil {
		h.handleDatabaseError(w, "create tracker", err)
		return
	}

	// Auto-update TrueSkill after tracker creation (following TrackerHandler pattern)
	result := h.updateUSLUserTrueSkillFromTrackers(createdTracker.DiscordID)
	if !result.Success {
		log.Printf("[USL-HANDLER] TrueSkill auto-update failed for %s: %s", createdTracker.DiscordID, result.Error)
		// Continue anyway - don't fail the tracker operation (graceful degradation)
	} else if result.TrueSkillResult != nil {
		log.Printf("[USL-HANDLER] Auto-updated TrueSkill for %s: μ=%.1f",
			createdTracker.DiscordID, result.TrueSkillResult.Mu)
	}

	log.Printf("[USL-HANDLER] Created tracker for user: %s", createdTracker.DiscordID)

	http.Redirect(w, r, fmt.Sprintf("/usl/trackers/detail?id=%d", createdTracker.ID), http.StatusSeeOther)
}

// EditTrackerForm displays the form for editing an existing tracker
func (h *TrackerHandler) EditTrackerForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	trackerID, err := h.parseUserID(r.URL.Query().Get("id"))
	if err != nil {
		h.handleInvalidID(w, "tracker")
		return
	}

	tracker, err := h.uslRepo.GetTrackerByID(trackerID)
	if err != nil {
		h.handleDatabaseError(w, "load tracker", err)
		return
	}
	if tracker == nil {
		http.NotFound(w, r)
		return
	}

	// Load associated user for display
	user, err := h.uslRepo.GetUserByDiscordID(tracker.DiscordID)
	if err != nil {
		h.handleDatabaseError(w, "load associated user", err)
		return
	}

	data := TrackerFormData{
		BasePageData: BasePageData{
			Title:       "Edit Tracker",
			CurrentPage: "trackers",
		},
		Tracker: tracker,
		User:    user,
		Errors:  make(map[string]string),
	}

	h.renderTemplate(w, TemplateUSLTrackerEdit, data)
}

// UpdateTracker handles updating an existing tracker
func (h *TrackerHandler) UpdateTracker(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.handleMethodNotAllowed(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.handleInvalidFormData(w, err)
		return
	}

	trackerID, err := h.parseUserID(r.FormValue("id"))
	if err != nil {
		h.handleInvalidID(w, "tracker")
		return
	}

	// Get existing tracker
	existingTracker, err := h.uslRepo.GetTrackerByID(trackerID)
	if err != nil {
		h.handleDatabaseError(w, "load tracker", err)
		return
	}
	if existingTracker == nil {
		http.NotFound(w, r)
		return
	}

	// Build updated tracker from form
	updatedTracker := h.buildTrackerFromForm(r)
	updatedTracker.ID = existingTracker.ID // Preserve ID

	// Validate the updated tracker
	validation := h.validateTrackerWithMetrics(r, updatedTracker)
	if !validation.IsValid {
		// Load user for error display
		user, _ := h.uslRepo.GetUserByDiscordID(updatedTracker.DiscordID)
		data := TrackerFormData{
			BasePageData: BasePageData{
				Title:       "Edit Tracker",
				CurrentPage: "trackers",
			},
			Tracker: updatedTracker,
			User:    user,
			Errors:  h.buildErrorMap(validation.Errors),
		}
		h.renderTemplate(w, TemplateUSLTrackerEdit, data)
		return
	}

	// Calculate MMR
	h.calculateEffectiveMMR(updatedTracker)

	// Update the tracker
	err = h.uslRepo.UpdateTracker(updatedTracker)
	if err != nil {
		h.handleDatabaseError(w, "update tracker", err)
		return
	}

	log.Printf("[USL-HANDLER] Updated tracker for user: %s", updatedTracker.DiscordID)

	// Redirect to tracker detail page
	http.Redirect(w, r, fmt.Sprintf("/usl/trackers/detail?id=%d", updatedTracker.ID), http.StatusSeeOther)
}
