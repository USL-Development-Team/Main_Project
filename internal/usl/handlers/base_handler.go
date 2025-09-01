package handlers

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"usl-server/internal/config"
	"usl-server/internal/services"
	"usl-server/internal/usl"
)

// BaseHandler contains shared dependencies and utilities used by all USL handlers
type BaseHandler struct {
	uslRepo          *usl.USLRepository
	templates        *template.Template
	trueskillService *services.UserTrueSkillService
	config           *config.Config
}

// NewBaseHandler creates a new BaseHandler with shared dependencies
func NewBaseHandler(
	uslRepo *usl.USLRepository,
	templates *template.Template,
	trueskillService *services.UserTrueSkillService,
	config *config.Config,
) *BaseHandler {
	return &BaseHandler{
		uslRepo:          uslRepo,
		templates:        templates,
		trueskillService: trueskillService,
		config:           config,
	}
}

// Shared utility methods

func (h *BaseHandler) parseUserID(userIDStr string) (int64, error) {
	if userIDStr == "" {
		return 0, ErrUserIDRequired
	}
	return strconv.ParseInt(userIDStr, 10, 64)
}

func (h *BaseHandler) getFormValue(r *http.Request, field FormField) string {
	return r.FormValue(string(field))
}

func (h *BaseHandler) getFormBoolValue(r *http.Request, field FormField) bool {
	return r.FormValue(string(field)) == "on" || r.FormValue(string(field)) == "true"
}

func (h *BaseHandler) getFormIntValue(r *http.Request, field FormField) int {
	return parseIntField(r.FormValue(string(field)))
}

// Shared error handlers

func (h *BaseHandler) handleMethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
}

func (h *BaseHandler) handleInvalidFormData(w http.ResponseWriter, err error) {
	http.Error(w, "Invalid form data: "+err.Error(), http.StatusBadRequest)
}

func (h *BaseHandler) handleDatabaseError(w http.ResponseWriter, operation string, err error) {
	http.Error(w, "Database error during "+operation+": "+err.Error(), http.StatusInternalServerError)
}

func (h *BaseHandler) handleInvalidID(w http.ResponseWriter, idType string) {
	http.Error(w, "Invalid "+idType+" ID", http.StatusBadRequest)
}

func (h *BaseHandler) handleParseError(w http.ResponseWriter, fieldName string) {
	http.Error(w, "Invalid "+fieldName, http.StatusBadRequest)
}

// Template rendering

func (h *BaseHandler) renderTemplate(w http.ResponseWriter, templateName TemplateName, data any) {
	if err := h.templates.ExecuteTemplate(w, string(templateName), data); err != nil {
		http.Error(w, "Template error: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// Helper functions (can be package-level since they don't need handler state)

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

// Common validation helpers

func isValidDiscordID(discordID string) bool {
	if len(discordID) < MinDiscordIDLength || len(discordID) > MaxDiscordIDLength {
		return false
	}
	for _, char := range discordID {
		if char < '0' || char > '9' {
			return false
		}
	}
	return true
}

// Common errors
var (
	ErrUserIDRequired = errors.New("User ID is required")
)
