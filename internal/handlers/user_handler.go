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

// UserHandler handles HTTP requests for user management
type UserHandler struct {
	userRepo  *repositories.UserRepository
	templates *template.Template
}

// NewUserHandler creates a new user handler
func NewUserHandler(userRepo *repositories.UserRepository, templates *template.Template) *UserHandler {
	return &UserHandler{
		userRepo:  userRepo,
		templates: templates,
	}
}

// ListUsers displays all users in HTML format
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	users, err := h.userRepo.GetAllUsers(false) // Get all users, including inactive
	if err != nil {
		log.Printf("Error getting users: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		Title string
		Users []*models.User
	}{
		Title: "User Management",
		Users: users,
	}

	if err := h.templates.ExecuteTemplate(w, "users.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

// NewUserForm displays the form for creating a new user
func (h *UserHandler) NewUserForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	data := struct {
		Title string
	}{
		Title: "Create New User",
	}

	if err := h.templates.ExecuteTemplate(w, "user_form.html", data); err != nil {
		log.Printf("Error executing template: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
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

	user, err := h.userRepo.CreateUser(userData)
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

	user, err := h.userRepo.FindUserByDiscordID(discordID)
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

	user, err := h.userRepo.UpdateUser(originalDiscordID, userData)
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

	user, err := h.userRepo.DeleteUser(discordID)
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

	users, err := h.userRepo.SearchUsers(query, 50)
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

	users, err := h.userRepo.GetAllUsers(false)
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
