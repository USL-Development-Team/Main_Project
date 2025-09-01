package handlers

import (
	"html/template"
	"usl-server/internal/config"
	"usl-server/internal/services"
	"usl-server/internal/usl"
)

// MigrationHandler provides a unified interface to all USL handlers for backward compatibility
// This maintains the existing API while using the new modular structure internally
type MigrationHandler struct {
	*UserHandler
	*TrackerHandler
	*AdminHandler
	baseHandler *BaseHandler
}

// NewMigrationHandler creates a new MigrationHandler that combines all specialized handlers
func NewMigrationHandler(
	uslRepo *usl.USLRepository,
	templates *template.Template,
	trueskillService *services.UserTrueSkillService,
	config *config.Config,
) *MigrationHandler {
	baseHandler := NewBaseHandler(uslRepo, templates, trueskillService, config)

	return &MigrationHandler{
		UserHandler:    NewUserHandler(baseHandler),
		TrackerHandler: NewTrackerHandler(baseHandler),
		AdminHandler:   NewAdminHandler(baseHandler),
		baseHandler:    baseHandler,
	}
}
