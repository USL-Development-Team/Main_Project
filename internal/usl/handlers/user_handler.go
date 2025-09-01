package handlers

import (
	"log"
	"net/http"
	"strings"
)

// UserHandler handles all user-related operations
type UserHandler struct {
	*BaseHandler
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(baseHandler *BaseHandler) *UserHandler {
	return &UserHandler{BaseHandler: baseHandler}
}

// ListUsers displays all users in a paginated list
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	users, err := h.uslRepo.GetAllUsers()
	if err != nil {
		h.handleDatabaseError(w, "load users", err)
		return
	}

	data := UsersListData{
		BasePageData: BasePageData{
			Title:       "Users",
			CurrentPage: "users",
		},
		Users: users,
		SearchConfig: SearchConfig{
			SearchPlaceholder: "Search by name or Discord ID...",
			SearchURL:         "/usl/users/search",
			SearchTarget:      "#users-tbody",
			ClearURL:          "/usl/users/search",
			ShowFilters:       true,
			Query:             "",
			StatusFilter:      "",
		},
	}

	h.renderTemplate(w, TemplateUSLUsers, data)
}

// SearchUsers handles user search requests and returns filtered results
func (h *UserHandler) SearchUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	query := r.URL.Query().Get("q")

	var users []*USLUser
	var err error
	if query == "" {
		users, err = h.uslRepo.GetAllUsers()
	} else {
		users, err = h.uslRepo.SearchUsers(query)
	}
	if err != nil {
		h.handleDatabaseError(w, "search users", err)
		return
	}

	data := UsersListData{
		BasePageData: BasePageData{
			Title:       "Users",
			CurrentPage: "",
		},
		Users: users,
		SearchConfig: SearchConfig{
			SearchPlaceholder: "Search by name or Discord ID...",
			SearchURL:         "/usl/users/search",
			SearchTarget:      "#users-tbody",
			ClearURL:          "/usl/users/search",
			ShowFilters:       true,
			Query:             query,
			StatusFilter:      "",
		},
	}

	h.renderTemplate(w, TemplateUSLUsersTable, data)
}

// UserDetail displays detailed information about a specific user
func (h *UserHandler) UserDetail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	userID, err := h.parseUserID(r.URL.Query().Get("id"))
	if err != nil {
		h.handleInvalidID(w, "user")
		return
	}

	user, err := h.uslRepo.GetUserByID(userID)
	if err != nil {
		h.handleDatabaseError(w, "load user", err)
		return
	}
	if user == nil {
		http.NotFound(w, r)
		return
	}

	userTrackers, err := h.uslRepo.GetTrackersByDiscordID(user.DiscordID)
	if err != nil {
		h.handleDatabaseError(w, "load user trackers", err)
		return
	}

	data := UserDetailData{
		BasePageData: BasePageData{
			Title:       user.Name,
			CurrentPage: "users",
		},
		User:         user,
		UserTrackers: userTrackers,
	}

	h.renderTemplate(w, TemplateUSLUserDetail, data)
}

// NewUserForm displays the form for creating a new user
func (h *UserHandler) NewUserForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	data := struct {
		Title       string
		User        *USLUser
		CurrentPage string
	}{
		Title:       "Add New User",
		User:        &USLUser{}, // Empty user for form
		CurrentPage: "users",
	}

	h.renderTemplate(w, "user-new-page", data)
}

// CreateUser handles the creation of a new user
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.handleMethodNotAllowed(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.handleInvalidFormData(w, err)
		return
	}

	// Extract user data from form
	user := &USLUser{
		Name:      strings.TrimSpace(r.FormValue("name")),
		DiscordID: strings.TrimSpace(r.FormValue("discord_id")),
		Active:    r.FormValue("active") == "on",
		Banned:    r.FormValue("banned") == "on",
	}

	// Validate required fields
	if user.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}
	if user.DiscordID == "" {
		http.Error(w, "Discord ID is required", http.StatusBadRequest)
		return
	}

	// Create the user
	createdUser, err := h.uslRepo.CreateUser(user.Name, user.DiscordID, user.Active, user.Banned)
	if err != nil {
		h.handleDatabaseError(w, "create user", err)
		return
	}

	log.Printf("[USL-HANDLER] Created user: %s (Discord: %s)", createdUser.Name, createdUser.DiscordID)

	// Redirect to user detail page
	http.Redirect(w, r, "/usl/users", http.StatusSeeOther)
}

// EditUserForm displays the form for editing an existing user
func (h *UserHandler) EditUserForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		h.handleMethodNotAllowed(w, r)
		return
	}

	userID, err := h.parseUserID(r.URL.Query().Get("id"))
	if err != nil {
		h.handleInvalidID(w, "user")
		return
	}

	user, err := h.uslRepo.GetUserByID(userID)
	if err != nil {
		h.handleDatabaseError(w, "load user", err)
		return
	}
	if user == nil {
		http.NotFound(w, r)
		return
	}

	data := struct {
		Title       string
		User        *USLUser
		CurrentPage string
	}{
		Title:       "Edit User",
		User:        user,
		CurrentPage: "users",
	}

	h.renderTemplate(w, "user-edit-page", data)
}

// UpdateUser handles updating an existing user
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		h.handleMethodNotAllowed(w, r)
		return
	}

	if err := r.ParseForm(); err != nil {
		h.handleInvalidFormData(w, err)
		return
	}

	userID, err := h.parseUserID(r.FormValue("id"))
	if err != nil {
		h.handleInvalidID(w, "user")
		return
	}

	// Get existing user
	existingUser, err := h.uslRepo.GetUserByID(userID)
	if err != nil {
		h.handleDatabaseError(w, "load user", err)
		return
	}
	if existingUser == nil {
		http.NotFound(w, r)
		return
	}

	// Update user fields
	existingUser.Name = strings.TrimSpace(r.FormValue("name"))
	existingUser.Active = r.FormValue("active") == "on"
	existingUser.Banned = r.FormValue("banned") == "on"

	// Validate required fields
	if existingUser.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		return
	}

	// Update the user
	updatedUser, err := h.uslRepo.UpdateUser(existingUser.ID, existingUser.Name, existingUser.Active, existingUser.Banned)
	if err != nil {
		h.handleDatabaseError(w, "update user", err)
		return
	}

	log.Printf("[USL-HANDLER] Updated user: %s", updatedUser.Name)

	// Redirect to user detail page
	http.Redirect(w, r, "/usl/users", http.StatusSeeOther)
}

// DeleteUser handles deleting a user
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		h.handleMethodNotAllowed(w, r)
		return
	}

	userID, err := h.parseUserID(r.URL.Query().Get("id"))
	if err != nil {
		h.handleInvalidID(w, "user")
		return
	}

	// Get user for logging
	user, err := h.uslRepo.GetUserByID(userID)
	if err != nil {
		h.handleDatabaseError(w, "load user", err)
		return
	}
	if user == nil {
		http.NotFound(w, r)
		return
	}

	// Delete the user
	if err := h.uslRepo.DeleteUser(userID); err != nil {
		h.handleDatabaseError(w, "delete user", err)
		return
	}

	log.Printf("[USL-HANDLER] Deleted user: %s", user.Name)

	// Redirect to users list
	http.Redirect(w, r, "/usl/users", http.StatusSeeOther)
}
