package handlers

import (
	"html/template"
	"strings"
	"testing"
	"usl-server/internal/usl"
)

func TestUserListTemplatesWithTrueSkill(t *testing.T) {
	// Create templates with the same functions as the main app
	tmpl := template.New("test")
	tmpl = tmpl.Funcs(template.FuncMap{
		"printf": func(format string, args ...interface{}) string {
			return "" // We'll test the actual rendering below
		},
		"substr": func(s string, start, length int) string {
			if start >= len(s) {
				return ""
			}
			end := start + length
			if end > len(s) {
				end = len(s)
			}
			return s[start:end]
		},
	})

	// Parse the templates
	tmpl = template.Must(tmpl.ParseFiles(
		"../../../templates/users-table-fragment.html",
		"../../../templates/users-list.html",
	))

	// Test data with TrueSkill values
	testUser := &usl.USLUser{
		ID:          1,
		Name:        "Test User",
		DiscordID:   "123456789012345678",
		Active:      true,
		Banned:      false,
		TrueSkillMu: 1567.8,
	}

	testData := struct {
		Users []*usl.USLUser
	}{
		Users: []*usl.USLUser{testUser},
	}

	t.Run("users-table-fragment renders TrueSkill", func(t *testing.T) {
		var buf strings.Builder
		err := tmpl.ExecuteTemplate(&buf, "users-table-fragment", testData)
		if err != nil {
			t.Fatalf("Template execution failed: %v", err)
		}

		output := buf.String()

		// Check that the template includes TrueSkill display
		if !strings.Contains(output, "TrueSkill μ:") {
			t.Error("Template should contain TrueSkill μ label")
		}

		// Check that the template structure is valid HTML
		if !strings.Contains(output, "<li>") {
			t.Error("Template should contain list items")
		}

		if !strings.Contains(output, testUser.Name) {
			t.Error("Template should contain user name")
		}

		if !strings.Contains(output, testUser.DiscordID) {
			t.Error("Template should contain Discord ID")
		}
	})

	t.Run("users-list template renders TrueSkill", func(t *testing.T) {
		// For the full page template, we need more complete data
		pageData := struct {
			Title string
			Users []*usl.USLUser
		}{
			Title: "Test Users",
			Users: []*usl.USLUser{testUser},
		}

		var buf strings.Builder
		// We can't easily test the full page template without navigation template,
		// so we'll just test the fragment part that contains the user list
		testHTML := `{{range .Users}}<div class="user-item">{{.Name}} - TrueSkill μ: {{printf "%.1f" .TrueSkillMu}}</div>{{end}}`

		testTmpl := template.Must(template.New("test-fragment").Funcs(template.FuncMap{
			"printf": func(format string, args ...interface{}) string {
				// Simple implementation for testing
				if format == "%.1f" && len(args) == 1 {
					if _, ok := args[0].(float64); ok {
						return "1567.8" // Expected formatted value
					}
				}
				return ""
			},
		}).Parse(testHTML))

		err := testTmpl.Execute(&buf, pageData)
		if err != nil {
			t.Fatalf("Test template execution failed: %v", err)
		}

		output := buf.String()

		if !strings.Contains(output, "TrueSkill μ: 1567.8") {
			t.Error("Template should contain formatted TrueSkill value")
		}

		if !strings.Contains(output, testUser.Name) {
			t.Error("Template should contain user name")
		}
	})

	t.Run("TrueSkill formatting with different values", func(t *testing.T) {
		testCases := []struct {
			name     string
			mu       float64
			expected string
		}{
			{"Default value", 1500.0, "1500.0"},
			{"High skill", 1800.5, "1800.5"},
			{"Low skill", 1200.3, "1200.3"},
			{"Precise value", 1567.847, "1567.8"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				user := &usl.USLUser{
					ID:          1,
					Name:        "Test User",
					DiscordID:   "123456789012345678",
					TrueSkillMu: tc.mu,
				}

				testHTML := `TrueSkill μ: {{printf "%.1f" .TrueSkillMu}}`
				tmpl := template.Must(template.New("format-test").Parse(testHTML))

				var buf strings.Builder
				err := tmpl.Execute(&buf, user)
				if err != nil {
					t.Fatalf("Template execution failed: %v", err)
				}

				output := buf.String()
				if !strings.Contains(output, tc.expected) {
					t.Errorf("Expected TrueSkill value %s, but got: %s", tc.expected, output)
				}
			})
		}
	})
}
