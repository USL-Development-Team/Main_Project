package handlers

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"strings"
	"usl-server/internal/middleware"
	"usl-server/internal/models"
	"usl-server/internal/repositories"
)

type UserHandler struct {
	repo      *repositories.UserRepository
	templates *template.Template
}

func NewUserHandler(userRepo *repositories.UserRepository, templates *template.Template) *UserHandler {
	return &UserHandler{
		repo:      userRepo,
		templates: templates,
	}
}

func (h *UserHandler) validateMethod(w http.ResponseWriter, r *http.Request, allowedMethod string) bool {
	if r.Method != allowedMethod {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func (h *UserHandler) getGuild(w http.ResponseWriter, r *http.Request) (*models.Guild, bool) {
	guild, ok := middleware.GetGuildFromRequest(r)
	if !ok {
		http.Error(w, "Guild context not found", http.StatusInternalServerError)
		return nil, false
	}
	return guild, true
}

func (h *UserHandler) renderFormWithErrors(w http.ResponseWriter, r *http.Request, title string, guild *models.Guild, user *models.User, errors map[string]string) {
	formData := &UserFormData{
		Title:  title,
		Guild:  guild,
		User:   user,
		Errors: errors,
	}

	if h.isHTMXRequest(r) {
		h.renderFragment(w, "user-form", formData)
	} else {
		h.renderTemplate(w, "content", formData)
	}
}

func (h *UserHandler) redirectToUserList(w http.ResponseWriter, r *http.Request, guild *models.Guild) {
	redirectURL := "/" + guild.Slug + "/users"

	if h.isHTMXRequest(r) {
		w.Header().Set("HX-Redirect", redirectURL)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Redirect(w, r, redirectURL, http.StatusSeeOther)
	}
}

func (h *UserHandler) validateUserForm(name, discordID string) map[string]string {
	errors := make(map[string]string)
	if strings.TrimSpace(name) == "" {
		errors["Name"] = "Display name is required"
	}
	if strings.TrimSpace(discordID) == "" {
		errors["DiscordID"] = "Discord ID is required"
	}
	return errors
}

func (h *UserHandler) renderTemplate(w http.ResponseWriter, templateName string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, templateName, data); err != nil {
		log.Printf("Template rendering error (%s): %v", templateName, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *UserHandler) renderFragment(w http.ResponseWriter, fragmentName string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, fragmentName, data); err != nil {
		log.Printf("Fragment rendering error (%s): %v", fragmentName, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *UserHandler) isHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func (h *UserHandler) renderErrorPage(w http.ResponseWriter, title, message string, statusCode int) {
	w.WriteHeader(statusCode)
	errorData := struct {
		Title   string
		Message string
	}{
		Title:   title,
		Message: message,
	}
	h.renderTemplate(w, "content", errorData)
}

type UserPageData struct {
	Title   string
	Guild   *models.Guild
	Users   []*models.User
	Query   string
	Page    int
	HasMore bool
}

type UserFormData struct {
	Title  string
	Guild  *models.Guild
	User   *models.User
	Errors map[string]string
}

// ListUsers displays all users in HTML format
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if !h.validateMethod(w, r, http.MethodGet) {
		return
	}

	guild, ok := h.getGuild(w, r)
	if !ok {
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))

	var users []*models.User
	var err error

	if query != "" {
		users, err = h.repo.SearchUsers(query, 50)
	} else {
		users, err = h.repo.GetAllUsers(false)
	}

	if err != nil {
		log.Printf("Failed to retrieve users: %v", err)
		h.renderErrorPage(w, "Error", "Unable to load users", http.StatusInternalServerError)
		return
	}

	pageData := &UserPageData{
		Title: "User Management",
		Guild: guild,
		Users: users,
		Query: query,
	}

	if h.isHTMXRequest(r) {
		h.renderFragment(w, "user-table", pageData)
	} else {
		h.renderTemplate(w, "content", pageData)
	}
}

// NewUserForm displays the form for creating a new user
func (h *UserHandler) NewUserForm(w http.ResponseWriter, r *http.Request) {
	if !h.validateMethod(w, r, http.MethodGet) {
		return
	}

	guild, ok := h.getGuild(w, r)
	if !ok {
		return
	}

	h.renderFormWithErrors(w, r, "Add New User", guild, nil, make(map[string]string))
}

// CreateUser handles the creation of a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if !h.validateMethod(w, r, http.MethodPost) {
		return
	}

	guild, ok := h.getGuild(w, r)
	if !ok {
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	userData := models.UserCreateRequest{
		Name:      strings.TrimSpace(r.FormValue("name")),
		DiscordID: strings.TrimSpace(r.FormValue("discord_id")),
		Active:    r.FormValue("active") == "true",
		Banned:    r.FormValue("banned") == "true",
		MMR:       0, // MMR will be calculated via TrueSkill
	}

	if errors := h.validateUserForm(userData.Name, userData.DiscordID); len(errors) > 0 {
		h.renderFormWithErrors(w, r, "Add New User", guild, nil, errors)
		return
	}

	user, err := h.repo.CreateUser(userData)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		errors := map[string]string{"general": "Failed to create user: " + err.Error()}
		h.renderFormWithErrors(w, r, "Add New User", guild, nil, errors)
		return
	}

	log.Printf("Created user: %s (%s)", user.Name, user.DiscordID)
	h.redirectToUserList(w, r, guild)
}

// EditUserForm displays the form for editing an existing user
func (h *UserHandler) EditUserForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	guild, ok := middleware.GetGuildFromRequest(r)
	if !ok {
		http.Error(w, "Guild context not found", http.StatusInternalServerError)
		return
	}

	userID := r.URL.Query().Get("id")
	if userID == "" {
		discordID := r.URL.Query().Get("discord_id")
		if discordID == "" {
			http.Error(w, "User ID or Discord ID is required", http.StatusBadRequest)
			return
		}

		user, err := h.repo.FindUserByDiscordID(discordID)
		if err != nil {
			log.Printf("Error finding user: %v", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		formData := &UserFormData{
			Title:  "Edit User",
			Guild:  guild,
			User:   user,
			Errors: make(map[string]string),
		}

		h.renderTemplate(w, "content", formData)
		return
	}

	// Handle numeric user ID
	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	user, err := h.repo.FindUserByID(id)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	formData := &UserFormData{
		Title:  "Edit User",
		Guild:  guild,
		User:   user,
		Errors: make(map[string]string),
	}

	h.renderTemplate(w, "content", formData)
}

// UpdateUser handles updating an existing user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodPut {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	guild, ok := middleware.GetGuildFromRequest(r)
	if !ok {
		http.Error(w, "Guild context not found", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	userID := r.FormValue("user_id")
	if userID == "" {
		originalDiscordID := r.FormValue("original_discord_id")
		if originalDiscordID == "" {
			http.Error(w, "User ID or Original Discord ID is required", http.StatusBadRequest)
			return
		}

		// Find user by Discord ID first
		user, err := h.repo.FindUserByDiscordID(originalDiscordID)
		if err != nil {
			log.Printf("Error finding user: %v", err)
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		userID = strconv.FormatInt(int64(user.ID), 10)
	}

	id, err := strconv.ParseInt(userID, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Restricted user update: only name, active, and banned status
	userData := models.UserUpdateRequest{
		Name:   strings.TrimSpace(r.FormValue("name")),
		Active: r.FormValue("active") == "true",
		Banned: r.FormValue("banned") == "true",
	}

	// Validation
	errors := make(map[string]string)
	if userData.Name == "" {
		errors["Name"] = "Display name is required"
	}

	if len(errors) > 0 {
		user, _ := h.repo.FindUserByID(id)
		formData := &UserFormData{
			Title:  "Edit User",
			Guild:  guild,
			User:   user,
			Errors: errors,
		}

		if h.isHTMXRequest(r) {
			h.renderFragment(w, "user-form", formData)
		} else {
			h.renderTemplate(w, "content", formData)
		}
		return
	}

	user, err := h.repo.UpdateUser(strconv.FormatInt(id, 10), userData)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		errors["general"] = "Failed to update user: " + err.Error()

		user, _ := h.repo.FindUserByID(id)
		formData := &UserFormData{
			Title:  "Edit User",
			Guild:  guild,
			User:   user,
			Errors: errors,
		}

		if h.isHTMXRequest(r) {
			h.renderFragment(w, "user-form", formData)
		} else {
			h.renderTemplate(w, "content", formData)
		}
		return
	}

	log.Printf("Updated user: %s (%s)", user.Name, user.DiscordID)

	if h.isHTMXRequest(r) {
		// Redirect via HTMX
		w.Header().Set("HX-Redirect", "/"+guild.Slug+"/users")
		w.WriteHeader(http.StatusOK)
	} else {
		http.Redirect(w, r, "/"+guild.Slug+"/users", http.StatusSeeOther)
	}
}

// DeleteUser handles deleting (marking inactive) a user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost && r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	guild, ok := middleware.GetGuildFromRequest(r)
	if !ok {
		http.Error(w, "Guild context not found", http.StatusInternalServerError)
		return
	}

	userID := r.URL.Query().Get("id")
	discordID := r.URL.Query().Get("discord_id")

	if userID == "" && discordID == "" {
		http.Error(w, "User ID or Discord ID is required", http.StatusBadRequest)
		return
	}

	var user *models.User
	var err error

	if discordID != "" {
		user, err = h.repo.DeleteUser(discordID)
	} else {
		// Convert userID to find by Discord ID first (legacy compatibility)
		id, parseErr := strconv.ParseInt(userID, 10, 64)
		if parseErr != nil {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}

		userObj, findErr := h.repo.FindUserByID(id)
		if findErr != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		user, err = h.repo.DeleteUser(userObj.DiscordID)
	}

	if err != nil {
		log.Printf("Error deleting user: %v", err)
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Deleted user: %s (%s)", user.Name, user.DiscordID)

	if h.isHTMXRequest(r) {
		// Return success response for HTMX
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("User deleted successfully")); err != nil {
			log.Printf("Failed to write delete response: %v", err)
		}
	} else {
		http.Redirect(w, r, "/"+guild.Slug+"/users", http.StatusSeeOther)
	}
}

// SearchUsers handles searching for users
func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// This method is now handled by ListUsers with query parameter
	h.ListUsers(w, r)
}

// ListUsersAPI returns users in JSON format
func (h *UserHandler) ListUsersAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.repo.GetAllUsers(false)
	if err != nil {
		log.Printf("Error getting users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Printf("Error encoding JSON: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
