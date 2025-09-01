package test

import (
	"fmt"
	"html/template"
	"strings"
	"testing"
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
			"../templates/base-layout.html",
			"../templates/navigation.html",
			"../templates/admin-dashboard.html",
		}

		for _, templateFile := range requiredTemplates {
			tmpl = template.Must(tmpl.ParseFiles(templateFile))
		}

		// Test data like USL migration handler provides
		data := struct {
			Title       string
			CurrentPage string
			Stats       struct {
				TotalUsers    int `json:"total_users"`
				ActiveUsers   int `json:"active_users"`
				TotalTrackers int `json:"total_trackers"`
				ValidTrackers int `json:"valid_trackers"`
			}
		}{
			Title:       "USL Admin Dashboard",
			CurrentPage: "admin",
			Stats: struct {
				TotalUsers    int `json:"total_users"`
				ActiveUsers   int `json:"active_users"`
				TotalTrackers int `json:"total_trackers"`
				ValidTrackers int `json:"valid_trackers"`
			}{
				TotalUsers:    10,
				ActiveUsers:   8,
				TotalTrackers: 5,
				ValidTrackers: 4,
			},
		}

		// Test that admin-dashboard-page template renders without errors
		var buf strings.Builder
		err := tmpl.ExecuteTemplate(&buf, "admin-dashboard-page", data)
		if err != nil {
			t.Fatalf("Failed to execute admin-dashboard-page template: %v", err)
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

	t.Run("Admin_Dashboard_Template_Exists", func(t *testing.T) {
		// Verify the admin-dashboard-page template is properly defined
		tmpl := template.New("test")

		// Parse the admin dashboard page and its dependencies
		tmpl, err := tmpl.ParseFiles("../templates/admin-dashboard.html", "../templates/navigation.html")
		if err != nil {
			t.Fatalf("Failed to parse admin-dashboard.html: %v", err)
		}

		// Test data
		data := struct {
			Title       string
			CurrentPage string
			Stats       struct {
				TotalUsers    int `json:"total_users"`
				ActiveUsers   int `json:"active_users"`
				TotalTrackers int `json:"total_trackers"`
				ValidTrackers int `json:"valid_trackers"`
			}
		}{
			Title:       "Dashboard",
			CurrentPage: "admin",
			Stats: struct {
				TotalUsers    int `json:"total_users"`
				ActiveUsers   int `json:"active_users"`
				TotalTrackers int `json:"total_trackers"`
				ValidTrackers int `json:"valid_trackers"`
			}{
				TotalUsers:    10,
				ActiveUsers:   8,
				TotalTrackers: 5,
				ValidTrackers: 4,
			},
		}

		// Test that admin-dashboard-page template renders
		var buf strings.Builder
		err = tmpl.ExecuteTemplate(&buf, "admin-dashboard-page", data)
		if err != nil {
			t.Fatalf("Failed to execute admin-dashboard-page template: %v", err)
		}

		output := buf.String()

		if !strings.Contains(output, "Dashboard") {
			t.Error("Admin-dashboard-page template should contain dashboard header")
		}
	})
}
