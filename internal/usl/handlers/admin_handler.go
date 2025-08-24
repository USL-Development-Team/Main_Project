package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"usl-server/internal/models"
)

// AdminHandler handles admin dashboard and API operations
type AdminHandler struct {
	*BaseHandler
}

// NewAdminHandler creates a new AdminHandler
func NewAdminHandler(baseHandler *BaseHandler) *AdminHandler {
	return &AdminHandler{BaseHandler: baseHandler}
}

// AdminDashboard displays the admin dashboard with statistics
func (h *AdminHandler) AdminDashboard(w http.ResponseWriter, r *http.Request) {
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

	data := AdminDashboardData{
		BasePageData: BasePageData{
			Title:       "Dashboard",
			CurrentPage: "admin",
		},
	}

	// Populate stats
	data.Stats.TotalUsers = stats["total_users"].(int)
	data.Stats.ActiveUsers = stats["active_users"].(int)
	data.Stats.TotalTrackers = stats["total_trackers"].(int)
	data.Stats.ValidTrackers = stats["valid_trackers"].(int)

	h.renderTemplate(w, TemplateUSLAdminDashboard, data)
}

// ListUsersAPI returns all users in JSON format for API consumption
func (h *AdminHandler) ListUsersAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	users, err := h.uslRepo.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}

// ListTrackersAPI returns all trackers in JSON format for API consumption
func (h *AdminHandler) ListTrackersAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	trackers, err := h.uslRepo.GetAllTrackers()
	if err != nil {
		http.Error(w, "Failed to fetch trackers", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(trackers); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}

// GetLeaderboardAPI returns leaderboard data in JSON format
func (h *AdminHandler) GetLeaderboardAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	leaderboard, err := h.uslRepo.GetLeaderboard()
	if err != nil {
		http.Error(w, "Failed to fetch leaderboard", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(leaderboard); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}

// UpdateUserTrueSkill handles TrueSkill update requests via HTMX
func (h *AdminHandler) UpdateUserTrueSkill(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.handleMethodNotAllowed(w, r)
		return
	}

	userID, err := h.parseUserID(r.URL.Query().Get("id"))
	if err != nil {
		h.handleInvalidID(w, "user")
		return
	}

	// Load the user
	user, err := h.uslRepo.GetUserByID(userID)
	if err != nil {
		h.handleDatabaseError(w, "load user", err)
		return
	}
	if user == nil {
		http.NotFound(w, r)
		return
	}

	// Update TrueSkill from trackers
	result := h.updateUSLUserTrueSkillFromTrackers(user.DiscordID)

	// Convert service result to our template data
	templateResult := &TrueSkillUpdateResult{
		Success:     result.Success,
		HadTrackers: result.HadTrackers,
		Error:       result.Error,
		UserName:    user.Name,
	}

	if result.TrueSkillResult != nil {
		templateResult.TrueSkillResult = &struct {
			Mu    float64 `json:"mu"`
			Sigma float64 `json:"sigma"`
		}{
			Mu:    result.TrueSkillResult.Mu,
			Sigma: result.TrueSkillResult.Sigma,
		}
	}

	// Render the update result template fragment
	h.renderTrueSkillUpdateResult(w, r, templateResult, user)
}

// ValidationMetricsAPI returns current validation metrics for monitoring
func (h *AdminHandler) ValidationMetricsAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	metricsMutex.RLock()
	defer metricsMutex.RUnlock()

	response := struct {
		Metrics       *ValidationMetrics `json:"metrics"`
		TopErrorTypes []struct {
			Type  string `json:"type"`
			Count int64  `json:"count"`
		} `json:"top_error_types"`
		TopErrorFields []struct {
			Field string `json:"field"`
			Count int64  `json:"count"`
		} `json:"top_error_fields"`
	}{
		Metrics: validationMetrics,
	}

	// Get top error types
	for errorType, count := range validationMetrics.ErrorsByType {
		response.TopErrorTypes = append(response.TopErrorTypes, struct {
			Type  string `json:"type"`
			Count int64  `json:"count"`
		}{
			Type:  errorType,
			Count: count,
		})
	}

	// Get top error fields
	for errorField, count := range validationMetrics.ErrorsByField {
		response.TopErrorFields = append(response.TopErrorFields, struct {
			Field string `json:"field"`
			Count int64  `json:"count"`
		}{
			Field: errorField,
			Count: count,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
		return
	}
}

// renderTrueSkillUpdateResult renders the TrueSkill update result as an HTMX fragment
func (h *AdminHandler) renderTrueSkillUpdateResult(w http.ResponseWriter, r *http.Request, result *TrueSkillUpdateResult, user *USLUser) {
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

	// Convert service result to our template data
	if result.TrueSkillResult != nil {
		data.TrueSkillResult = &struct {
			Mu    float64 `json:"mu"`
			Sigma float64 `json:"sigma"`
		}{
			Mu:    result.TrueSkillResult.Mu,
			Sigma: result.TrueSkillResult.Sigma,
		}
	}

	if result.Success {
		log.Printf("[USL-HANDLER] TrueSkill updated for user %s: μ=%.3f, σ=%.3f",
			user.Name, data.TrueSkillResult.Mu, data.TrueSkillResult.Sigma)
	} else {
		log.Printf("[USL-HANDLER] TrueSkill update failed for user %s: %s",
			user.Name, result.Error)
	}

	// For HTMX requests, render the update result fragment
	h.renderTemplate(w, "trueskill-update-result", data)
}
