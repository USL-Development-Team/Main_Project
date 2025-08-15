package test

import (
	"fmt"
	"html/template"
	"strings"
	"testing"
	"usl-server/internal/models"
)

// TestTemplateRendering tests that our templates can be loaded and rendered without errors
func TestTemplateRendering(t *testing.T) {
	t.Run("Admin_Dashboard_Template_Loads", func(t *testing.T) {
		// Load templates like the server does
		tmpl := template.New("app")

		// Add template functions (same as server)
		tmpl = tmpl.Funcs(template.FuncMap{
			"dict": func(values ...interface{}) map[string]interface{} {
				dict := make(map[string]interface{})
				for i := 0; i < len(values); i += 2 {
					if i+1 < len(values) {
						dict[values[i].(string)] = values[i+1]
					}
				}
				return dict
			},
			"slice": func(values ...interface{}) []interface{} {
				return values
			},
			"add": func(a, b int) int {
				return a + b
			},
			"sub": func(a, b float64) float64 {
				return a - b
			},
			"mul": func(a, b float64) float64 {
				return a * b
			},
			"printf": func(format string, args ...interface{}) string {
				return fmt.Sprintf(format, args...)
			},
			"lt": func(a, b float64) bool {
				return a < b
			},
		})

		// Parse specific templates needed for admin dashboard
		requiredTemplates := []string{
			"../templates/layouts/admin.html",
			"../templates/partials/navigation.html",
			"../templates/pages/admin-dashboard.html",
		}

		for _, templateFile := range requiredTemplates {
			tmpl = template.Must(tmpl.ParseFiles(templateFile))
		}

		// Test data like USL migration handler provides
		testGuild := &models.Guild{
			ID:             1,
			DiscordGuildID: "1390537743385231451",
			Name:           "USL",
			Slug:           "usl",
			Active:         true,
			Config:         models.GetDefaultGuildConfig(),
			Theme:          models.GetDefaultTheme(),
		}

		data := struct {
			Title        string
			Guild        *models.Guild
			Stats        map[string]interface{}
			CurrentPage  string
			User         interface{}
			FlashMessage string
			FlashType    string
		}{
			Title: "USL Admin Dashboard",
			Guild: testGuild,
			Stats: map[string]interface{}{
				"total_users":  10,
				"active_users": 8,
			},
			CurrentPage:  "admin",
			User:         nil, // Not logged in for test
			FlashMessage: "",
			FlashType:    "",
		}

		// Test that admin-layout template renders without errors
		var buf strings.Builder
		err := tmpl.ExecuteTemplate(&buf, "admin-layout", data)
		if err != nil {
			t.Fatalf("Failed to execute admin-layout template: %v", err)
		}

		output := buf.String()

		// Basic checks that the template rendered correctly
		if !strings.Contains(output, "<!DOCTYPE html>") {
			t.Error("Template output should contain DOCTYPE")
		}

		if !strings.Contains(output, "USL Admin Dashboard") {
			t.Error("Template output should contain title")
		}

		if !strings.Contains(output, "USL") {
			t.Error("Template output should contain guild name")
		}

		if !strings.Contains(output, "/static/dist/output.css") {
			t.Error("Template output should contain Tailwind CSS link")
		}

		if !strings.Contains(output, "/static/htmx.min.js") {
			t.Error("Template output should contain HTMX script")
		}

		if len(output) < 1000 {
			t.Errorf("Template output seems too short (%d chars), might be broken", len(output))
		}
	})

	t.Run("Admin_Page_Template_Exists", func(t *testing.T) {
		// Verify the admin-page template is properly defined
		tmpl := template.New("test")

		// Parse the admin dashboard page
		tmpl, err := tmpl.ParseFiles("../templates/pages/admin-dashboard.html")
		if err != nil {
			t.Fatalf("Failed to parse admin-dashboard.html: %v", err)
		}

		// Test data
		testGuild := &models.Guild{
			ID:   1,
			Name: "USL",
			Slug: "usl",
		}

		data := struct {
			Guild *models.Guild
			Stats map[string]interface{}
		}{
			Guild: testGuild,
			Stats: map[string]interface{}{
				"total_users":    10,
				"active_users":   8,
				"total_trackers": 5,
				"valid_trackers": 4,
			},
		}

		// Test that admin-page template renders
		var buf strings.Builder
		err = tmpl.ExecuteTemplate(&buf, "admin-page", data)
		if err != nil {
			t.Fatalf("Failed to execute admin-page template: %v", err)
		}

		output := buf.String()

		if !strings.Contains(output, "USL Dashboard") {
			t.Error("Admin-page template should contain dashboard header")
		}
	})
}
