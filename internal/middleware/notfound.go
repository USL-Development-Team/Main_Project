package middleware

import (
	"html/template"
	"log/slog"
	"net/http"
)

// NotFoundHandler creates a 404 handler that serves a custom 404 page
func NotFoundHandler(logger *slog.Logger, templates *template.Template) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		logger.Warn("404 Not Found",
			"method", r.Method,
			"path", r.URL.Path,
			"query", r.URL.RawQuery,
			"user_agent", r.UserAgent(),
			"remote_addr", r.RemoteAddr,
		)

		w.WriteHeader(http.StatusNotFound)
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		data := struct {
			Title string
			Path  string
		}{
			Title: "Page Not Found",
			Path:  r.URL.Path,
		}

		if err := templates.ExecuteTemplate(w, "404.html", data); err != nil {
			logger.Error("Failed to render 404 template", "error", err)
			// Fallback to plain text 404
			http.Error(w, "404 - Page Not Found", http.StatusNotFound)
		}
	}
}
