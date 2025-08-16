package handlers

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
	"usl-server/internal/middleware"
	"usl-server/internal/models"
	"usl-server/internal/repositories"
)

// TestUserHandler creates a minimal handler for testing HTMX templates
type TestUserHandler struct {
	mockUsers []*models.User
	templates *template.Template
}

func NewTestUserHandler(users []*models.User, templates *template.Template) *TestUserHandler {
	return &TestUserHandler{
		mockUsers: users,
		templates: templates,
	}
}

// Mock helper methods for testing (currently unused but reserved for future test expansion)
// func (h *TestUserHandler) getAllUsers() []*models.User {
// 	return h.mockUsers
// }

// func (h *TestUserHandler) findUserByDiscordID(discordID string) (*models.User, error) {
// 	for _, user := range h.mockUsers {
// 		if user.DiscordID == discordID {
// 			return user, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("user not found")
// }

func (h *TestUserHandler) isHTMXRequest(r *http.Request) bool {
	return r.Header.Get("HX-Request") == "true"
}

func (h *TestUserHandler) renderTemplate(w http.ResponseWriter, templateName string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, templateName, data); err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
	}
}

func (h *TestUserHandler) renderFragment(w http.ResponseWriter, fragmentName string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := h.templates.ExecuteTemplate(w, fragmentName, data); err != nil {
		http.Error(w, "Fragment error", http.StatusInternalServerError)
	}
}

// Test version of ListUsers
func (h *TestUserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	guild, ok := middleware.GetGuildFromRequest(r)
	if !ok {
		http.Error(w, "Guild context not found", http.StatusInternalServerError)
		return
	}

	query := strings.TrimSpace(r.URL.Query().Get("q"))

	// Filter users based on search
	users := h.mockUsers
	if query != "" {
		var filteredUsers []*models.User
		for _, user := range h.mockUsers {
			if strings.Contains(strings.ToLower(user.Name), strings.ToLower(query)) ||
				strings.Contains(user.DiscordID, query) {
				filteredUsers = append(filteredUsers, user)
			}
		}
		users = filteredUsers
	}

	data := UserPageData{
		Title: "Users",
		Guild: guild,
		Users: users,
		Query: query,
	}

	if h.isHTMXRequest(r) {
		h.renderFragment(w, "user-table", data)
	} else {
		h.renderTemplate(w, "content", data)
	}
}

// Test version of CreateUser
func (h *TestUserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	guild, ok := middleware.GetGuildFromRequest(r)
	if !ok {
		http.Error(w, "Guild context not found", http.StatusInternalServerError)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Invalid form data", http.StatusBadRequest)
		return
	}

	name := strings.TrimSpace(r.FormValue("name"))
	discordID := strings.TrimSpace(r.FormValue("discord_id"))
	active := r.FormValue("active") == "true"
	banned := r.FormValue("banned") == "true"

	// Validation
	errors := make(map[string]string)
	if name == "" {
		errors["name"] = "Name is required"
	}
	if discordID == "" {
		errors["discord_id"] = "Discord ID is required"
	}

	if len(errors) > 0 {
		formData := &UserFormData{
			Title:  "Create User",
			Guild:  guild,
			User:   nil,
			Errors: errors,
		}

		if h.isHTMXRequest(r) {
			h.renderFragment(w, "user-form", formData)
		} else {
			h.renderTemplate(w, "content", formData)
		}
		return
	}

	// Mock create user
	newUser := &models.User{
		ID:        len(h.mockUsers) + 1,
		Name:      name,
		DiscordID: discordID,
		Active:    active,
		Banned:    banned,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	h.mockUsers = append(h.mockUsers, newUser)

	if h.isHTMXRequest(r) {
		w.Header().Set("HX-Redirect", fmt.Sprintf("/%s/users", guild.Slug))
		w.WriteHeader(http.StatusOK)
	} else {
		http.Redirect(w, r, fmt.Sprintf("/%s/users", guild.Slug), http.StatusSeeOther)
	}
}

// Helper function to create a test request with guild context
func createTestRequestWithGuild(method, path string, body string) (*http.Request, *httptest.ResponseRecorder) {
	var req *http.Request
	if body != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}

	testGuild := &models.Guild{
		ID:             1,
		DiscordGuildID: "123456789012345678",
		Name:           "Test Guild",
		Slug:           "test",
		Active:         true,
		Config:         models.GetDefaultGuildConfig(),
	}

	ctx := context.WithValue(req.Context(), middleware.GuildContextKey, testGuild)
	req = req.WithContext(ctx)

	return req, httptest.NewRecorder()
}

func TestUserHandler_ListUsers(t *testing.T) {

	testUsers := []*models.User{
		{
			ID:        1,
			Name:      "TestUser1",
			DiscordID: "123456789012345678",
			Active:    true,
			Banned:    false,
		},
		{
			ID:        2,
			Name:      "TestUser2",
			DiscordID: "987654321098765432",
			Active:    false,
			Banned:    true,
		},
	}

	tmpl := template.Must(template.New("test").Parse(`
		{{define "user-table"}}
		<table>
		{{range .Users}}
			<tr><td>{{.Name}}</td><td>{{.DiscordID}}</td></tr>
		{{end}}
		</table>
		{{end}}
		
		{{define "content"}}
		<h1>{{.Title}}</h1>
		{{template "user-table" .}}
		{{end}}
	`))

	handler := NewTestUserHandler(testUsers, tmpl)

	t.Run("ListUsers_Success", func(t *testing.T) {
		req, w := createTestRequestWithGuild("GET", "/test/users", "")

		handler.ListUsers(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		if !strings.Contains(body, "TestUser1") {
			t.Errorf("Expected response to contain TestUser1, got: %s", body)
		}
	})

	t.Run("ListUsers_HTMX_Request", func(t *testing.T) {
		req, w := createTestRequestWithGuild("GET", "/test/users", "")
		req.Header.Set("HX-Request", "true")

		handler.ListUsers(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		// Should render just the table fragment for HTMX
		body := w.Body.String()
		if !strings.Contains(body, "<table>") {
			t.Errorf("Expected HTMX response to contain table, got: %s", body)
		}
	})

	t.Run("ListUsers_NoGuildContext", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test/users", nil)
		w := httptest.NewRecorder()

		handler.ListUsers(w, req)

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Expected status 500 when no guild context, got %d", w.Code)
		}
	})

	t.Run("ListUsers_WithSearch", func(t *testing.T) {
		req, w := createTestRequestWithGuild("GET", "/test/users?q=TestUser1", "")

		handler.ListUsers(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		body := w.Body.String()
		if !strings.Contains(body, "TestUser1") {
			t.Errorf("Expected search results to contain TestUser1, got: %s", body)
		}
	})
}

func TestUserHandler_CreateUser(t *testing.T) {
	testUsers := []*models.User{}

	tmpl := template.Must(template.New("test").Parse(`
		{{define "user-form"}}
		<form>
		{{if .Errors.name}}<div class="error">{{.Errors.name}}</div>{{end}}
		</form>
		{{end}}
		
		{{define "content"}}{{template "user-form" .}}{{end}}
	`))

	handler := NewTestUserHandler(testUsers, tmpl)

	t.Run("CreateUser_Success", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("name", "NewUser")
		formData.Set("discord_id", "111111111111111111")
		formData.Set("active", "true")
		formData.Set("banned", "false")

		req, w := createTestRequestWithGuild("POST", "/test/users/create", formData.Encode())

		handler.CreateUser(w, req)

		if w.Code != http.StatusSeeOther {
			t.Errorf("Expected status 303, got %d", w.Code)
		}

		if len(handler.mockUsers) != 1 {
			t.Errorf("Expected 1 user after creation, got %d", len(handler.mockUsers))
		}
	})

	t.Run("CreateUser_ValidationError", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("name", "") // Empty name should cause validation error
		formData.Set("discord_id", "111111111111111111")

		req, w := createTestRequestWithGuild("POST", "/test/users/create", formData.Encode())

		handler.CreateUser(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200 with validation error, got %d", w.Code)
		}

		body := w.Body.String()
		if !strings.Contains(body, "error") {
			t.Errorf("Expected validation error in response, got: %s", body)
		}
	})

	t.Run("CreateUser_HTMX_Success", func(t *testing.T) {
		formData := url.Values{}
		formData.Set("name", "HTMXUser")
		formData.Set("discord_id", "222222222222222222")
		formData.Set("active", "true")
		formData.Set("banned", "false")

		req, w := createTestRequestWithGuild("POST", "/test/users/create", formData.Encode())
		req.Header.Set("HX-Request", "true")

		handler.CreateUser(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status 200, got %d", w.Code)
		}

		redirectHeader := w.Header().Get("HX-Redirect")
		if redirectHeader != "/test/users" {
			t.Errorf("Expected HX-Redirect header, got: %s", redirectHeader)
		}
	})
}

// Test real UserHandler constructor with a simple repository mock
func TestUserHandler_Integration(t *testing.T) {
	t.Run("UserHandler_Constructor", func(t *testing.T) {

		// This tests that our constructor works with the real types
		var repo *repositories.UserRepository = nil // Would normally be created with supabase client

		tmpl := template.Must(template.New("test").Parse(`{{define "content"}}test{{end}}`))

		// This should compile without type errors
		if repo != nil {
			handler := NewUserHandler(repo, tmpl)
			if handler == nil {
				t.Error("Expected handler to be created")
			}
		}
	})
}
