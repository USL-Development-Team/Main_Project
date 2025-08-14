package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"usl-server/internal/services"
)

const (
	// Template names
	TrueSkillResultTemplate     = "trueskill_result.html"
	UserTrueSkillResultTemplate = "user_trueskill_result.html"
	TrueSkillStatsTemplate      = "trueskill_stats.html"

	// Page titles
	TrueSkillUpdateCompleteTitle = "TrueSkill Update Complete"
	TrueSkillRecalculationTitle  = "TrueSkill Recalculation Complete"
	UserTrueSkillUpdateTitle     = "User TrueSkill Update"
	TrueSkillStatsTitle          = "TrueSkill Statistics"
)

// TrueSkillHandler handles HTTP requests for TrueSkill calculations
type TrueSkillHandler struct {
	trueSkillService *services.UserTrueSkillService
	templates        *template.Template
}

// NewTrueSkillHandler creates a new TrueSkill handler
func NewTrueSkillHandler(trueSkillService *services.UserTrueSkillService, templates *template.Template) *TrueSkillHandler {
	return &TrueSkillHandler{
		trueSkillService: trueSkillService,
		templates:        templates,
	}
}

// UpdateAllUserTrueSkill handles batch TrueSkill updates for all users (HTML response)
func (h *TrueSkillHandler) UpdateAllUserTrueSkill(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	result, err := h.trueSkillService.UpdateAllUserTrueSkill()
	if err != nil {
		log.Printf("Error updating TrueSkill for all users: %v", err)
		http.Error(w, "Failed to update TrueSkill values: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("Batch TrueSkill update completed: %d users processed, %d calculated from trackers, %d given defaults",
		result.ProcessedCount, result.TrackerBasedCount, result.DefaultCount)

	data := struct {
		Title  string
		Result *services.BatchUpdateResult
	}{
		Title:  TrueSkillUpdateCompleteTitle,
		Result: result,
	}

	if err := h.templates.ExecuteTemplate(w, TrueSkillResultTemplate, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// UpdateUserTrueSkill handles TrueSkill update for a single user
func (h *TrueSkillHandler) UpdateUserTrueSkill(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	discordID := r.URL.Query().Get("discord_id")
	if discordID == "" {
		http.Error(w, "Discord ID is required", http.StatusBadRequest)
		return
	}

	result := h.trueSkillService.UpdateUserTrueSkillFromTrackers(discordID)

	data := struct {
		Title     string
		DiscordID string
		Result    *services.TrueSkillUpdateResult
	}{
		Title:     "User TrueSkill Update",
		DiscordID: discordID,
		Result:    result,
	}

	if err := h.templates.ExecuteTemplate(w, UserTrueSkillResultTemplate, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// RecalculateAllUserTrueSkill handles recalculation of TrueSkill for all users
func (h *TrueSkillHandler) RecalculateAllUserTrueSkill(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	result, err := h.trueSkillService.RecalculateAllUserTrueSkill()
	if err != nil {
		log.Printf("Error recalculating TrueSkill for all users: %v", err)
		http.Error(w, "Failed to recalculate TrueSkill values: "+err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("TrueSkill recalculation completed: %d users processed", result.ProcessedCount)

	data := struct {
		Title  string
		Result *services.BatchUpdateResult
	}{
		Title:  "TrueSkill Recalculation Complete",
		Result: result,
	}

	if err := h.templates.ExecuteTemplate(w, TrueSkillResultTemplate, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// GetTrueSkillStats displays TrueSkill service statistics
func (h *TrueSkillHandler) GetTrueSkillStats(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	stats := h.trueSkillService.GetTrueSkillStats()

	data := struct {
		Title string
		Stats map[string]interface{}
	}{
		Title: "TrueSkill Service Statistics",
		Stats: stats,
	}

	if err := h.templates.ExecuteTemplate(w, TrueSkillStatsTemplate, data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// UpdateAllUserTrueSkillAPI handles batch TrueSkill updates via API (JSON response)
func (h *TrueSkillHandler) UpdateAllUserTrueSkillAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	result, err := h.trueSkillService.UpdateAllUserTrueSkill()
	if err != nil {
		log.Printf("API: Error updating TrueSkill for all users: %v", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		if encodeErr := json.NewEncoder(w).Encode(map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}); encodeErr != nil {
			log.Printf("Failed to encode error response: %v", encodeErr)
		}
		return
	}

	log.Printf("API: Batch TrueSkill update completed: %d users processed", result.ProcessedCount)

	response := map[string]interface{}{
		"success": true,
		"result":  result,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
