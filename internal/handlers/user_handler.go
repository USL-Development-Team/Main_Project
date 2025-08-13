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

type UserHandler struct {
	userRepository *repositories.UserRepository
	templates      *template.Template
}

func NewUserHandler(userRepo *repositories.UserRepository, templates *template.Template) *UserHandler {
	return &UserHandler{
		userRepository: userRepo,
		templates:      templates,
	}
}

// Helper methods for common operations

func (h *UserHandler) validateHTTPMethod(w http.ResponseWriter, r *http.Request, allowedMethod string) bool {
	if r.Method != allowedMethod {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return false
	}
	return true
}

func (h *UserHandler) renderTemplate(w http.ResponseWriter, templateName string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, templateName, data); err != nil {
		log.Printf("Template rendering error (%s): %v", templateName, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
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
	h.renderTemplate(w, "error.html", errorData)
}

type PageData struct {
	Title string
	Users []*models.User
}

type FormData struct {
	Title string
	User  *models.User
}

// ListUsers displays all users in HTML format
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if !h.validateHTTPMethod(w, r, http.MethodGet) {
		return
	}

	users, err := h.userRepository.GetAllUsers(false)
	if err != nil {
		log.Printf("Failed to retrieve users: %v", err)
		h.renderErrorPage(w, "Error", "Unable to load users", http.StatusInternalServerError)
		return
	}

	pageData := &PageData{
		Title: "User Management - USL Server",
		Users: users,
	}

	h.renderTemplate(w, "users.html", pageData)
}

// NewUserForm displays the form for creating a new user
func (h *UserHandler) NewUserForm(w http.ResponseWriter, r *http.Request) {
	if !h.validateHTTPMethod(w, r, http.MethodGet) {
		return
	}

	formData := &FormData{
		Title: "Create New User - USL Server",
		User:  nil,
	}

	h.renderTemplate(w, "user_form.html", formData)
}

// CreateUser handles the creation of a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	// Parse MMR
	mmr := 0
	if mmrStr := r.FormValue("mmr"); mmrStr != "" {
		if parsedMMR, err := strconv.Atoi(mmrStr); err == nil {
			mmr = parsedMMR
		}
	}

	userData := models.UserCreateRequest{
		Name:      r.FormValue("name"),
		DiscordID: r.FormValue("discord_id"),
		Active:    r.FormValue("active") == "true",
		Banned:    r.FormValue("banned") == "true",
		MMR:       mmr,
	}

	user, err := h.userRepository.CreateUser(userData)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Failed to create user: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Created user: %s (%s)", user.Name, user.DiscordID)
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// EditUserForm displays the form for editing an existing user
func (h *UserHandler) EditUserForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	discordID := r.URL.Query().Get("discord_id")
	if discordID == "" {
		http.Error(w, "Discord ID is required", http.StatusBadRequest)
		return
	}

	user, err := h.userRepository.FindUserByDiscordID(discordID)
	if err != nil {
		log.Printf("Error finding user: %v", err)
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	data := struct {
		Title string
		User  *models.User
	}{
		Title: "Edit User",
		User:  user,
	}

	if err := h.templates.ExecuteTemplate(w, "user_edit_form.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// UpdateUser handles updating an existing user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	originalDiscordID := r.FormValue("original_discord_id")
	if originalDiscordID == "" {
		http.Error(w, "Original Discord ID is required", http.StatusBadRequest)
		return
	}

	// Parse MMR
	mmr := 0
	if mmrStr := r.FormValue("mmr"); mmrStr != "" {
		if parsedMMR, err := strconv.Atoi(mmrStr); err == nil {
			mmr = parsedMMR
		}
	}

	userData := models.UserUpdateRequest{
		Name:      r.FormValue("name"),
		DiscordID: r.FormValue("discord_id"),
		Active:    r.FormValue("active") == "true",
		Banned:    r.FormValue("banned") == "true",
		MMR:       mmr,
	}

	user, err := h.userRepository.UpdateUser(originalDiscordID, userData)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		http.Error(w, "Failed to update user: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Updated user: %s (%s)", user.Name, user.DiscordID)
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// DeleteUser handles deleting (marking inactive) a user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	discordID := r.URL.Query().Get("discord_id")
	if discordID == "" {
		http.Error(w, "Discord ID is required", http.StatusBadRequest)
		return
	}

	user, err := h.userRepository.DeleteUser(discordID)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		http.Error(w, "Failed to delete user: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Deleted user: %s (%s)", user.Name, user.DiscordID)
	http.Redirect(w, r, "/users", http.StatusSeeOther)
}

// SearchUsers handles searching for users
func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Redirect(w, r, "/users", http.StatusSeeOther)
		return
	}

	users, err := h.userRepository.SearchUsers(query, 50)
	if err != nil {
		log.Printf("Error searching users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
		Users []*models.User
		Query string
	}{
		Title: "User Search Results",
		Users: users,
		Query: query,
	}

	if err := h.templates.ExecuteTemplate(w, "users.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// ListUsersAPI returns users in JSON format for API calls
func (h *UserHandler) ListUsersAPI(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.userRepository.GetAllUsers(false)
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
