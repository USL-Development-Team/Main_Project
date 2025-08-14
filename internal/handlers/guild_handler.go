package handlers

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"usl-server/internal/models"
	"usl-server/internal/repositories"
	"usl-server/internal/services"
)

const (
	// Number parsing constants
	ParseIntBase    = 10
	ParseIntBitSize = 64
)

type GuildHandler struct {
	guildRepository   *repositories.GuildRepository
	permissionService *services.PermissionService
	templates         *template.Template
}

func NewGuildHandler(guildRepo *repositories.GuildRepository, permService *services.PermissionService, templates *template.Template) *GuildHandler {
	return &GuildHandler{
		guildRepository:   guildRepo,
		permissionService: permService,
		templates:         templates,
	}
}

// Helper methods for common operations

func (h *GuildHandler) validateHTTPMethod(w http.ResponseWriter, r *http.Request, allowedMethod string) bool {
	if r.Method != allowedMethod {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func (h *GuildHandler) renderTemplate(w http.ResponseWriter, templateName string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, templateName, data); err != nil {
		log.Printf("Template rendering error (%s): %v", templateName, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *GuildHandler) renderErrorPage(w http.ResponseWriter, title, message string, statusCode int) {
	w.WriteHeader(statusCode)
	errorData := struct {
		Title   string
		Message string
	}{
		Title:   title,
		Message: message,
	}
	h.renderTemplate(w, "error.html", errorData)
}

func (h *GuildHandler) renderJSONResponse(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("JSON encoding error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// parseGuildIDFromQuery extracts and validates guild ID from query parameters
func (h *GuildHandler) parseGuildIDFromQuery(r *http.Request) (int64, error) {
	guildIDStr := r.URL.Query().Get("id")
	if guildIDStr == "" {
		return 0, fmt.Errorf("Guild ID is required")
	}

	guildID, err := strconv.ParseInt(guildIDStr, ParseIntBase, ParseIntBitSize)
	if err != nil {
		return 0, fmt.Errorf("Invalid guild ID")
	}

	return guildID, nil
}

// Guild management endpoints

// GetGuildConfig handles GET /guilds/{id}/config
func (h *GuildHandler) GetGuildConfig(w http.ResponseWriter, r *http.Request) {
	if !h.validateHTTPMethod(w, r, http.MethodGet) {
		return
	}

	guildID, err := h.parseGuildIDFromQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	config, err := h.guildRepository.GetConfig(guildID)
	if err != nil {
		log.Printf("Failed to get guild config: %v", err)
		h.renderErrorPage(w, "Configuration Error", "Failed to load guild configuration", http.StatusInternalServerError)
		return
	}

	h.renderJSONResponse(w, config, http.StatusOK)
}

// UpdateGuildConfig handles PUT /guilds/{id}/config
func (h *GuildHandler) UpdateGuildConfig(w http.ResponseWriter, r *http.Request) {
	if !h.validateHTTPMethod(w, r, http.MethodPut) {
		return
	}

	guildID, err := h.parseGuildIDFromQuery(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var config models.GuildConfig
	if err := json.NewDecoder(r.Body).Decode(&config); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate configuration immediately
	if err := config.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Update the configuration
	if err := h.guildRepository.UpdateConfig(guildID, &config); err != nil {
		log.Printf("Failed to update guild config: %v", err)
		h.renderErrorPage(w, "Update Error", "Failed to update guild configuration", http.StatusInternalServerError)
		return
	}

	h.renderJSONResponse(w, map[string]string{"status": "success"}, http.StatusOK)
}

// CreateGuild handles POST /guilds
func (h *GuildHandler) CreateGuild(w http.ResponseWriter, r *http.Request) {
	if !h.validateHTTPMethod(w, r, http.MethodPost) {
		return
	}

	var guildData models.GuildCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&guildData); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Set default config if not provided
	if guildData.Config.Discord.BotCommandPrefix == "" {
		guildData.Config = models.GetDefaultGuildConfig()
	}

	// Validate the guild data
	if err := guildData.Config.Validate(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	guild, err := h.guildRepository.CreateGuild(guildData)
	if err != nil {
		log.Printf("Failed to create guild: %v", err)
		h.renderErrorPage(w, "Creation Error", "Failed to create guild", http.StatusInternalServerError)
		return
	}

	h.renderJSONResponse(w, guild, http.StatusCreated)
}

// GetAllGuilds handles GET /guilds
func (h *GuildHandler) GetAllGuilds(w http.ResponseWriter, r *http.Request) {
	if !h.validateHTTPMethod(w, r, http.MethodGet) {
		return
	}

	activeOnlyStr := r.URL.Query().Get("active_only")
	activeOnly := activeOnlyStr == "true"

	guilds, err := h.guildRepository.GetAllGuilds(activeOnly)
	if err != nil {
		log.Printf("Failed to get guilds: %v", err)
		h.renderErrorPage(w, "Retrieval Error", "Failed to load guilds", http.StatusInternalServerError)
		return
	}

	h.renderJSONResponse(w, guilds, http.StatusOK)
}

// GetGuild handles GET /guilds/{id}
func (h *GuildHandler) GetGuild(w http.ResponseWriter, r *http.Request) {
	if !h.validateHTTPMethod(w, r, http.MethodGet) {
		return
	}

	guildIDStr := r.URL.Query().Get("id")
	if guildIDStr == "" {
		http.Error(w, "Guild ID is required", http.StatusBadRequest)
		return
	}

	guildID, err := strconv.ParseInt(guildIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid guild ID", http.StatusBadRequest)
		return
	}

	guild, err := h.guildRepository.FindGuildByID(guildID)
	if err != nil {
		log.Printf("Failed to get guild: %v", err)
		h.renderErrorPage(w, "Not Found", "Guild not found", http.StatusNotFound)
		return
	}

	h.renderJSONResponse(w, guild, http.StatusOK)
}

// GetGuildByDiscordID handles GET /guilds/discord/{discord_id}
func (h *GuildHandler) GetGuildByDiscordID(w http.ResponseWriter, r *http.Request) {
	if !h.validateHTTPMethod(w, r, http.MethodGet) {
		return
	}

	discordGuildID := r.URL.Query().Get("discord_id")
	if discordGuildID == "" {
		http.Error(w, "Discord Guild ID is required", http.StatusBadRequest)
		return
	}

	guild, err := h.guildRepository.FindGuildByDiscordID(discordGuildID)
	if err != nil {
		log.Printf("Failed to get guild by Discord ID: %v", err)
		h.renderErrorPage(w, "Not Found", "Guild not found", http.StatusNotFound)
		return
	}

	h.renderJSONResponse(w, guild, http.StatusOK)
}

// DeactivateGuild handles DELETE /guilds/{id}
func (h *GuildHandler) DeactivateGuild(w http.ResponseWriter, r *http.Request) {
	if !h.validateHTTPMethod(w, r, http.MethodDelete) {
		return
	}

	guildIDStr := r.URL.Query().Get("id")
	if guildIDStr == "" {
		http.Error(w, "Guild ID is required", http.StatusBadRequest)
		return
	}

	guildID, err := strconv.ParseInt(guildIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid guild ID", http.StatusBadRequest)
		return
	}

	guild, err := h.guildRepository.DeactivateGuild(guildID)
	if err != nil {
		log.Printf("Failed to deactivate guild: %v", err)
		h.renderErrorPage(w, "Deactivation Error", "Failed to deactivate guild", http.StatusInternalServerError)
		return
	}

	h.renderJSONResponse(w, guild, http.StatusOK)
}

// CheckPermissions handles GET /guilds/{id}/permissions
func (h *GuildHandler) CheckPermissions(w http.ResponseWriter, r *http.Request) {
	if !h.validateHTTPMethod(w, r, http.MethodGet) {
		return
	}

	guildIDStr := r.URL.Query().Get("id")
	if guildIDStr == "" {
		http.Error(w, "Guild ID is required", http.StatusBadRequest)
		return
	}

	guildID, err := strconv.ParseInt(guildIDStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid guild ID", http.StatusBadRequest)
		return
	}

	// For now, we'll use a placeholder for user roles
	// In a real implementation, this would come from Discord OAuth or session
	userRoles := r.URL.Query()["role"]

	permissions := h.permissionService.GetUserPermissions(userRoles, guildID)
	h.renderJSONResponse(w, permissions, http.StatusOK)
}
