package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"strings"
	"usl-server/internal/models"
	"usl-server/internal/repositories"
)

type contextKey string

const (
	GuildContextKey contextKey = "guild"
)

// GuildContextMiddleware extracts guild information from URL slug and adds it to request context
type GuildContextMiddleware struct {
	guildRepo *repositories.GuildRepository
	logger    *slog.Logger
}

// NewGuildContextMiddleware creates a new guild context middleware
func NewGuildContextMiddleware(guildRepo *repositories.GuildRepository, logger *slog.Logger) *GuildContextMiddleware {
	return &GuildContextMiddleware{
		guildRepo: guildRepo,
		logger:    logger,
	}
}

// GuildContext returns a middleware that extracts guild from URL slug
func (m *GuildContextMiddleware) GuildContext() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract guild slug from URL path
			// Expected format: /{guild-slug}/... or /api/{guild-slug}/...
			path := strings.TrimPrefix(r.URL.Path, "/")
			pathParts := strings.Split(path, "/")

			var guildSlug string
			if len(pathParts) > 0 {
				// Handle /api/{guild-slug}/... format
				if pathParts[0] == "api" && len(pathParts) > 1 {
					guildSlug = pathParts[1]
				} else {
					// Handle /{guild-slug}/... format
					guildSlug = pathParts[0]
				}
			}

			if guildSlug == "" {
				// No guild slug in URL, continue without guild context
				next.ServeHTTP(w, r)
				return
			}

			// Skip guild context for certain routes
			if shouldSkipGuildContext(guildSlug) {
				next.ServeHTTP(w, r)
				return
			}

			// Fetch guild from database
			guild, err := m.guildRepo.FindGuildBySlug(guildSlug)
			if err != nil {
				m.logger.Warn("Guild not found", "slug", guildSlug, "error", err)
				// For now, create a mock guild for USL compatibility
				if guildSlug == "usl" {
					guild = &models.Guild{
						ID:             1,
						DiscordGuildID: "123456789012345678",
						Name:           "USL",
						Slug:           "usl",
						Active:         true,
						Config:         models.GetDefaultGuildConfig(),
					}
				} else {
					http.NotFound(w, r)
					return
				}
			}

			if !guild.Active {
				http.Error(w, "Guild is not active", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), GuildContextKey, guild)
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}

// shouldSkipGuildContext determines if guild context should be skipped for certain routes
func shouldSkipGuildContext(slug string) bool {
	skipRoutes := []string{
		"static",
		"auth",
		"login",
		"logout",
		"health",
		"favicon.ico",
		".well-known",
		"robots.txt",
		"sitemap.xml",
	}

	for _, skipRoute := range skipRoutes {
		if slug == skipRoute {
			return true
		}
	}

	return false
}

// GetGuildFromContext extracts guild from request context
func GetGuildFromContext(ctx context.Context) (*models.Guild, bool) {
	guild, ok := ctx.Value(GuildContextKey).(*models.Guild)
	return guild, ok
}

// GetGuildFromRequest is a convenience function to get guild from request context
func GetGuildFromRequest(r *http.Request) (*models.Guild, bool) {
	return GetGuildFromContext(r.Context())
}
